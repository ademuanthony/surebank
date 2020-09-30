package sale

import (
	"context"
	"fmt"
	"math/rand"
	"merryworld/surebank/internal/transaction"
	"strconv"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"merryworld/surebank/internal/inventory"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/postgres/models"
)

var (
	// ErrNotFound abstracts the postgres not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// Find gets all the sales from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, _ auth.Claims, req FindRequest) (*PagedResponseList, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.Sale.Find")
	defer span.Finish()

	var queries []QueryMod

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
	}

	if req.IncludeBranch {
		queries = append(queries, Load(models.SaleRels.Branch))
	}

	if req.IncludeArchivedBy {
		queries = append(queries, Load(models.SaleRels.ArchivedBy))
	}

	if req.IncludeCreatedBy {
		queries = append(queries, Load(models.SaleRels.CreatedBy))
	}

	if req.IncludeUpdatedBy {
		queries = append(queries, Load(models.SaleRels.UpdatedBy))
	}

	totalCount, err := models.Sales(queries...).Count(ctx, repo.DbConn)

	if len(req.Order) > 0 {
		for _, s := range req.Order {
			queries = append(queries, OrderBy(s))
		}
	} else {
		queries = append(queries, OrderBy("created_at desc"))
	}

	if req.Limit != nil {
		queries = append(queries, Limit(int(*req.Limit)))
	}

	if req.Offset != nil {
		queries = append(queries, Offset(int(*req.Offset)))
	}

	saleSlice, err := models.Sales(queries...).All(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	var result = Sales{}
	for _, rec := range saleSlice {
		result = append(result, FromModel(rec))
	}

	return &PagedResponseList{
		Sales:      result.Response(ctx),
		TotalCount: totalCount,
	}, nil
}

// ReadByID gets the specified sale by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, claims auth.Claims, id string) (*Sale, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.Sale.ReadByID")
	defer span.Finish()

	var queries = []QueryMod{
		models.SaleWhere.ID.EQ(id),
		Load(models.SaleRels.UpdatedBy),
		Load(models.SaleRels.CreatedBy),
		Load(models.SaleRels.ArchivedBy),
		Load(models.SaleRels.Branch),
	}
	sale, err := models.Sales(queries...).One(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	sale.R.SaleItems, err = models.SaleItems(models.SaleItemWhere.SaleID.EQ(id), Load(models.SaleItemRels.Product)).All(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	return FromModel(sale), nil
}

func (repo *Repository) MakeSale(ctx context.Context, claims auth.Claims, req MakeSalesRequest, now time.Time) (*Sale, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.sale.makeSale")
	defer span.Finish()
	if claims.Audience == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return nil, err
	}

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot start a db transaction")
	}

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	salesRep, err := models.Users(models.UserWhere.ID.EQ(claims.Subject)).One(ctx, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, weberror.NewErrorMessage(ctx, err, 400, "Something went wrong. Are you logged in?")
	}

	saleID := uuid.NewRandom().String()
	var itemSlice models.SaleItemSlice
	var amount float64
	for _, item := range req.Items {
		prod, err := repo.ShopRepo.ReadProductByIDTx(ctx, claims, item.ProductID, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, weberror.WithMessagef(ctx, err, "Invalid product ID, %s", item.ProductID)
		}
		bal, err := repo.InventoryRepo.Balance(ctx, claims, item.ProductID, salesRep.BranchID, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, weberror.WithMessagef(ctx, err, "Cannot get stock balance for product, %s", prod.Name)
		}

		if bal < int64(item.Quantity) {
			_ = tx.Rollback()
			return nil, weberror.NewError(ctx, weberror.WithMessagef(ctx, errors.New("Low stock balance"),
				"%s is remaining %d, cannot sell %d", prod.Name, bal, item.Quantity), 400)
		}

		if _, err = repo.InventoryRepo.MakeStockDeduction(ctx, claims, inventory.MakeStockDeductionRequest{
			ProductID: item.ProductID,
			Quantity:  int64(item.Quantity),
			Ref:       fmt.Sprintf("Sold out, %s", saleID),
		}, now, tx); err != nil {
			_ = tx.Rollback()
			return nil, weberror.NewError(ctx, weberror.WithMessagef(ctx, err, "Cannot make stock deduction for %s", prod.Name), 400)
		}

		itemSlice = append(itemSlice, &models.SaleItem{
			ID:            uuid.NewRandom().String(),
			SaleID:        saleID,
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			UnitPrice:     prod.Price,
			UnitCostPrice: prod.Price,
		})
		amount += float64(item.Quantity) * prod.Price
	}

	if req.PaymentMethod == "cash" && amount > req.AmountTender {
		_ = tx.Rollback()
		return nil, weberror.NewError(ctx, fmt.Errorf("you must collect %f from the customer to make this sale", amount), 400)
	}

	receiptNumber := repo.generateReceiptNumber(ctx)

	if req.PaymentMethod == "wallet" {
		if req.AccountNumber == "" {
			_ = tx.Rollback()
			return nil, weberror.NewError(ctx, errors.New("You must specify the buyer's account number to use wallet for payment"), 400)
		}
		_, err = repo.TransactionRepo.MakeDeduction(ctx, claims, transaction.MakeDeductionRequest{
			AccountNumber: req.AccountNumber,
			Amount:        amount,
			Narration:     fmt.Sprintf("sale:%s:%s", receiptNumber, saleID),
		}, now, tx)
		if err != nil {
			_ = tx.Rollback()
			return nil, weberror.NewErrorMessage(ctx, err, 500, "cannot make deduction")
		}
	}

	// If now empty set it to the current time.
	if now.IsZero() {
		now = time.Now()
	}

	// Always store the time as UTC.
	now = now.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)

	sale := Sale{
		ID:            saleID,
		ReceiptNumber: receiptNumber,
		Amount:        amount,
		AmountTender:  req.AmountTender,
		Balance:       req.AmountTender - amount,
		CustomerName:  req.CustomerName,
		PhoneNumber:   req.PhoneNumber,
		CreatedAt:     now,
		UpdatedAt:     now,
		CreatedByID:   claims.Subject,
		UpdatedByID:   claims.Subject,
		BranchID:      salesRep.BranchID,
	}

	saleModel := sale.model()
	if err := saleModel.Insert(ctx, tx, boil.Infer()); err != nil {
		_ = tx.Rollback()
		return nil, weberror.WithMessage(ctx, err, "Cannot save sale")
	}

	for _, item := range itemSlice {
		if err = item.Insert(ctx, tx, boil.Infer()); err != nil {
			_ = tx.Rollback()
			return nil, weberror.WithMessage(ctx, err, "Cannot insert sales items")
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, weberror.WithMessage(ctx, err, "Unable to commit DB transaction")
	}

	return &sale, nil
}

func (repo *Repository) generateReceiptNumber(ctx context.Context) string {
	var receipt string
	for receipt == "" || repo.saleExist(ctx, receipt) {
		receipt = "SB"
		rand.Seed(time.Now().UTC().UnixNano())
		for i := 0; i < 6; i++ {
			receipt += strconv.Itoa(rand.Intn(10))
		}
	}
	return receipt
}

func (repo *Repository) saleExist(ctx context.Context, receiptNumber string) bool {
	exist, _ := models.Sales(models.SaleWhere.ReceiptNumber.EQ(receiptNumber)).Exists(ctx, repo.DbConn)
	return exist
}

// Archive soft deleted the sale from the database.
func (repo *Repository) Archive(ctx context.Context, claims auth.Claims, req ArchiveRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.sale.Archive")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update branches they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}
	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	// If now empty set it to the current time.
	if now.IsZero() {
		now = time.Now()
	}

	// Always store the time as UTC.
	now = now.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)

	_, err = models.Sales(models.SaleWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, models.M{models.SaleColumns.ArchivedAt: now})

	return nil
}
