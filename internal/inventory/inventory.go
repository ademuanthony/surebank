package inventory

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/transaction"
)

var (
	// ErrNotFound abstracts the postgres not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// Find gets all the inventory from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, _ auth.Claims, req FindRequest) (*PagedResponseList, error) {
	var queries []QueryMod

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
	}

	totalCount, err := models.Inventories(queries...).Count(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get inventory count")
	}

	if len(req.Order) > 0 {
		for _, s := range req.Order {
			queries = append(queries, OrderBy(s))
		}
	}

	if req.IncludeBranch {
		queries = append(queries, Load(models.InventoryRels.Branch))
	}

	if req.IncludeProduct {
		queries = append(queries, Load(models.InventoryRels.Product))
	}

	if req.IncludeSalesRep {
		queries = append(queries, Load(models.InventoryRels.SalesRep))
	}

	if req.Limit != nil {
		queries = append(queries, Limit(int(*req.Limit)))
	}

	if req.Offset != nil {
		queries = append(queries, Offset(int(*req.Offset)))
	}

	slice, err := models.Inventories(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &PagedResponseList{}, nil
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	var result Inventories
	for _, rec := range slice {
		result = append(result, FromModel(rec))
	}

	if len(result) == 0 {
		return &PagedResponseList{}, nil
	}

	return &PagedResponseList{
		Transactions: result.Response(ctx),
		TotalCount:   totalCount,
	}, nil
}

// ReadByID gets the specified inventory by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, _ auth.Claims, id string) (*Inventory, error) {
	queries := []QueryMod{
		models.InventoryWhere.ID.EQ(id),
		Load(models.InventoryRels.Branch),
		Load(models.InventoryRels.Product),
		Load(models.InventoryRels.SalesRep),
	}
	model, err := models.Inventories(queries...).One(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	return FromModel(model), nil
}

// Balance gets the balance of the specified product from the database.
func (repo *Repository) Balance(ctx context.Context, claim auth.Claims, productID string, branchID string) (balance int64, err error) {
	lastTx, err := repo.lastTransaction(ctx, productID, branchID)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			return 0, err
		}
	}

	return accountBalanceAtTx(lastTx), nil
}

// AddStock inserts a new inventory transaction into the database.
func (repo *Repository) AddStock(ctx context.Context, claims auth.Claims, req AddStockRequest, now time.Time) (*Inventory, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.inventory.addStock")
	defer span.Finish()
	if claims.Subject == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	salesRep, err := models.FindUser(ctx, repo.DbConn, claims.Subject)
	if err != nil {
		return nil, errors.WithStack(ErrForbidden)
	}

	// Validate the request.
	v := webcontext.Validator()
	err = v.Struct(req)
	if err != nil {
		return nil, err
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

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return nil, err
	}

	lastTransaction, err := repo.lastTransaction(ctx, req.ProductID, salesRep.BranchID)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			_ = tx.Rollback()
			return nil, err
		}
	}

	m := models.Inventory{
		ID:             uuid.NewRandom().String(),
		ProductID:      req.ProductID,
		BranchID:       salesRep.BranchID,
		TXType:         transaction.TransactionType_Deposit.String(),
		OpeningBalance: float64(accountBalanceAtTx(lastTransaction)),
		Quantity:       req.Quantity,
		SalesRepID:     claims.Subject,
		CreatedAt:      now.Unix(),
		UpdatedAt:      now.Unix(),
	}

	if err := m.Insert(ctx, tx, boil.Infer()); err != nil {
		_ = tx.Rollback()
		return nil, errors.WithMessage(err, "Insert deposit failed")
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Inventory{
		ID:             m.ID,
		ProductID:      m.ProductID,
		Quantity:       m.Quantity,
		OpeningBalance: m.OpeningBalance,
		Narration:      m.Narration,
		TXType:         m.TXType,
		SalesRepID:     m.SalesRepID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}, nil
}

// RemoveStock deducts from an inventory in the database.
func (repo *Repository) RemoveStock(ctx context.Context, claims auth.Claims, req RemoveStockRequest, now time.Time) (*Inventory, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.inventory.removeStock")
	defer span.Finish()
	if claims.Subject == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	if !claims.HasRole(auth.RoleSuperAdmin) {
		return nil, ErrForbidden
	}

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "starting DB transaction")
	}
	inv, err := repo.MakeStockDeduction(ctx, claims, MakeStockDeductionRequest{ProductID: req.ProductID,
		Quantity: req.Quantity, Ref: req.Reason}, now, tx)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}
	_ = tx.Commit()
	return inv, nil
}

// MakeStockDeduction inserts a new inventory transaction into the database.
func (repo *Repository) MakeStockDeduction(ctx context.Context, claims auth.Claims, req MakeStockDeductionRequest, now time.Time, tx *sql.Tx) (*Inventory, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.inventory.makeStockDeduction")
	defer span.Finish()
	if claims.Subject == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	salesRep, err := models.FindUser(ctx, tx, claims.Subject)
	if err != nil {
		return nil, errors.WithStack(ErrForbidden)
	}

	// Validate the request.
	v := webcontext.Validator()
	err = v.Struct(req)
	if err != nil {
		return nil, err
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

	repo.mutex.Lock()
	defer repo.mutex.Unlock()

	lastTransaction, err := repo.lastTransaction(ctx, req.ProductID, salesRep.BranchID)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			return nil, err
		}
	}

	balance := accountBalanceAtTx(lastTransaction)
	if balance < req.Quantity {
		return nil, errors.New("not enough quantity. Aborted")
	}

	m := models.Inventory{
		ID:             uuid.NewRandom().String(),
		ProductID:      req.ProductID,
		BranchID:       salesRep.BranchID,
		TXType:         transaction.TransactionType_Withdrawal.String(),
		OpeningBalance: float64(balance),
		Quantity:       float64(req.Quantity),
		SalesRepID:     claims.Subject,
		Narration:      req.Ref,
		CreatedAt:      now.Unix(),
		UpdatedAt:      now.Unix(),
	}

	if err := m.Insert(ctx, tx, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "Insert deposit failed")
	}

	return &Inventory{
		ID:             m.ID,
		ProductID:      m.ProductID,
		Quantity:       m.Quantity,
		OpeningBalance: m.OpeningBalance,
		Narration:      m.Narration,
		TXType:         m.TXType,
		SalesRepID:     m.SalesRepID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}, nil
}

// lastTransaction returns the last transaction for the specified product
func (repo *Repository) lastTransaction(ctx context.Context, productID string, branchID string) (*models.Inventory, error) {
	return models.Inventories(
		models.InventoryWhere.ProductID.EQ(productID),
		models.InventoryWhere.BranchID.EQ(branchID),
		OrderBy(fmt.Sprintf("%s desc", models.InventoryColumns.CreatedAt)),
		Limit(1),
	).One(ctx, repo.DbConn)
}

// accountBalanceAtTx returns the account balance as at the specified tx
func accountBalanceAtTx(tx *models.Inventory) int64 {
	var lastBalance float64
	if tx != nil {
		if tx.TXType == transaction.TransactionType_Deposit.String() {
			lastBalance = tx.OpeningBalance + tx.Quantity
		} else {
			lastBalance = tx.OpeningBalance - tx.Quantity
		}
	}

	return int64(lastBalance)
}

// Archive soft deleted the stock transaction from the database.
func (repo *Repository) Archive(ctx context.Context, claims auth.Claims, req ArchiveRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.deposit.Archive")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update customer they have access to.
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

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return err
	}

	tranx, err := models.Inventories(models.InventoryWhere.ID.EQ(req.ID)).One(ctx, tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = models.Inventories(models.InventoryWhere.ID.EQ(req.ID)).UpdateAll(ctx, tx, models.M{models.InventoryColumns.ArchivedAt: now})
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	var txAmount = tranx.Quantity
	if tranx.TXType == transaction.TransactionType_Withdrawal.String() {
		txAmount *= -1
	}

	// updated all trans after it
	_, err = models.Inventories(
		models.InventoryWhere.CreatedAt.GT(tranx.CreatedAt)).
		UpdateAll(ctx, tx, models.M{models.InventoryColumns.OpeningBalance: fmt.Sprintf("opening_balance + (%f)", txAmount)})

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// update all product balance
	_, err = models.Products(models.ProductWhere.ID.EQ(tranx.ProductID)).UpdateAll(ctx, tx,
		models.M{models.ProductColumns.StockBalance: fmt.Sprintf("stocke_balance + (%f)", txAmount)})

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func (repo *Repository) Report(ctx context.Context, claims auth.Claims, req ReportRequest) (*PagedStockInfo, error) {
	salesRep, err := models.FindUser(ctx, repo.DbConn, claims.Subject)
	if err != nil {
		return nil, errors.WithStack(ErrForbidden)
	}

	statement := `select p.name as product_name, t.product_id, t.opening_balance, t.quantity, t.tx_type from inventory t
			inner join product p on t.product_id = p.id
			inner join (
				select product_id, max(created_at) as MaxDate
				from inventory
				group by product_id
			) tm on t.product_id = tm.product_id and t.created_at = tm.MaxDate`

	if !claims.HasRole(auth.RoleAdmin) {
		if len(req.Where) > 0 {
			req.Where = fmt.Sprintf("(%s) and branch_id = ?", req.Where)
		} else {
			req.Where = "branch_id = ?"
		}
		req.Args = append(req.Args, salesRep.BranchID)
	}

	if req.Where != "" {
		statement = fmt.Sprintf("%s where %s", statement, req.Where)
	}

	if len(req.Order) > 0 {
		statement = fmt.Sprintf("%s order by %s", statement, strings.Join(req.Order, ","))
	}

	if req.Offset != nil {
		statement = fmt.Sprintf("%s offset %d", statement, req.Offset)
	}

	if req.Limit != nil {
		statement = fmt.Sprintf("%s offset %d", statement, req.Limit)
	}

	var stockInfos []StockInfo
	err = models.Inventories(SQL(statement)).Bind(ctx, repo.DbConn, &stockInfos)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &PagedStockInfo{}, nil
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	var result = PagedStockInfo{
		Inventories: stockInfos,
	}

	if len(stockInfos) == 0 {
		return &PagedStockInfo{}, nil
	}

	return &result, nil
}
