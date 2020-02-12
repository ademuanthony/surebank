package sale

import (
	"context"
	"github.com/pborman/uuid"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/shop"
	"time"

	"github.com/pkg/errors"
	. "github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web/webcontext"
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
	defer  repo.mutex.Unlock()

	var itemSlice models.SaleItemSlice
	var amount float64
	for _, item := range req.Items {
		prod, err := repo.ShopRepo.ReadProductByID(ctx, claims, item.ProductID)
		if err != nil {
			return nil, weberror.WithMessagef(ctx, err, "Invalid product ID, %s", item.ProductID)
		}
		bal, err := repo.ShopRepo.StockBalance(ctx, claims, shop.StockBalanceRequest{ProductID: item.ProductID})
		if err != nil {
			return nil, weberror.WithMessagef(ctx, err, "Cannot get stock balance for product, %s", prod.Name)
		}

		if bal < int64(item.Quantity) {
			return nil, weberror.WithMessagef(ctx, errors.New("Low stock balance"),
				"%s is remaining %d, cannot sell %d", prod.Name, bal, item.Quantity)
		}

		itemSlice = append(itemSlice, &models.SaleItem{
			ID:            uuid.NewRandom().String(),
			ProductID:     item.ProductID,
			Quantity:      item.Quantity,
			UnitPrice:     prod.Price,
			UnitCostPrice: prod.Price,
		})
		amount += float64(item.Quantity) * prod.Price
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

	salesRep, err := models.Users(models.UserWhere.ID.EQ(claims.Subject)).One(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, 400, "Something went wrong. Are you logged in?")
	}

	sale := Sale{
		ID:            uuid.NewRandom().String(),
		ReceiptNumber: generateReceiptNumber(),
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

}

func generateReceiptNumber() string {

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

	_,err = models.Sales(models.SaleWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, models.M{models.SaleColumns.ArchivedAt: now})

	return nil
}
