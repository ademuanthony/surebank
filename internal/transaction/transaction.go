package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
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
)

var (
	// ErrNotFound abstracts the postgres not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// Find gets all the transaction from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, _ auth.Claims, req FindRequest) (*PagedResponseList, error) {
	var queries = []QueryMod {
		Load(models.TransactionRels.SalesRep),
		Load(models.TransactionRels.Account),
	}

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
	}

	totalCount, err := models.Transactions(queries...).Count(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get transaction count")
	}

	if req.IncludeAccount {
		queries = append(queries, Load(models.TransactionRels.Account))
	}

	if req.IncludeSalesRep {
		queries = append(queries, Load(models.TransactionRels.SalesRep))
	}

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

	slice, err := models.Transactions(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &PagedResponseList{}, nil
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	var result Transactions
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

// ReadByID gets the specified transaction by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, _ auth.Claims, id string) (*Transaction, error) {
	queries := []QueryMod{
		models.TransactionWhere.ID.EQ(id),
		Load(models.TransactionRels.Account),
		Load(models.TransactionRels.SalesRep),
	}
	model, err := models.Transactions(queries...).One(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	return FromModel(model), nil
}

// AccountBalance gets the balance of the specified account from the database.
func (repo *Repository) AccountBalance(ctx context.Context, _ auth.Claims, accountNumber string) (balance float64, err error) {
	account, err := models.Accounts(models.AccountWhere.Number.EQ(accountNumber)).One(ctx, repo.DbConn)
	if err != nil {
		return 0, weberror.WithMessage(ctx, err, "Invalid account number")
	}

	lastTx, err := repo.lastTransaction(ctx, account.ID)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			return 0, err
		}
	}

	return accountBalanceAtTx(lastTx), nil
}

// Create inserts a new transaction into the database.
func (repo *Repository) Create(ctx context.Context, claims auth.Claims, req CreateRequest, now time.Time) (*Transaction, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.transaction.Create")
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

	account, err := models.Accounts(models.AccountWhere.Number.EQ(req.AccountNumber)).One(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, http.StatusBadRequest, "Invalid account number")
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

	repo.accNumMtx.Lock()
	defer repo.accNumMtx.Unlock()

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return nil, err
	}

	lastTransaction, err := repo.lastTransaction(ctx, account.ID)
	if err != nil {
		if err.Error() != sql.ErrNoRows.Error() {
			_ = tx.Rollback()
			return nil, err
		}
	}

	m := models.Transaction{
		ID:             uuid.NewRandom().String(),
		AccountID:      account.ID,
		OpeningBalance: accountBalanceAtTx(lastTransaction),
		Amount:         req.Amount,
		Narration:      req.Narration,
		TXType:         req.Type.String(),
		SalesRepID:     claims.Subject,
		CreatedAt:      now.Unix(),
		UpdatedAt:      now.Unix(),
	}

	if err := m.Insert(ctx, tx, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "Insert deposit failed")
	}

	if _, err := models.Accounts(models.AccountWhere.ID.EQ(account.ID)).UpdateAll(ctx, tx, models.M{
		models.AccountColumns.Balance: accountBalanceAtTx(&m),
	}); err != nil {

		_ = tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return &Transaction{
		ID:             m.ID,
		AccountID:      m.AccountID,
		Amount:         m.Amount,
		OpeningBalance: m.OpeningBalance,
		Narration:      m.Narration,
		Type:           TransactionType(m.TXType),
		SalesRepID:     m.SalesRepID,
		CreatedAt:      time.Unix(m.CreatedAt, 0),
		UpdatedAt:      time.Unix(m.UpdatedAt, 0),
	}, nil
}

// lastTransaction returns the last transaction for the specified account
func (repo *Repository) lastTransaction(ctx context.Context, accountID string) (*models.Transaction, error) {
	return models.Transactions(
		models.TransactionWhere.AccountID.EQ(accountID),
		OrderBy(fmt.Sprintf("%s desc", models.TransactionColumns.CreatedAt)),
		Limit(1),
	).One(ctx, repo.DbConn)
}

// accountBalanceAtTx returns the account balance as at the specified tx
func accountBalanceAtTx(tx *models.Transaction) float64 {
	var lastBalance float64
	if tx != nil {
		if tx.TXType == TransactionType_Deposit.String() {
			lastBalance = tx.OpeningBalance + tx.Amount
		} else {
			lastBalance = tx.OpeningBalance - tx.Amount
		}
	}

	return lastBalance
}

// Update replaces an exiting transaction in the database.
func (repo *Repository) Update(ctx context.Context, claims auth.Claims, req UpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.deposit.Update")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update transactions they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	cols := models.M{}
	if req.Narration != nil {
		cols[models.TransactionColumns.Narration] = req.Narration
	}

	if req.Amount != nil {
		cols[models.TransactionColumns.Amount] = req.Amount
	}

	if len(cols) == 0 {
		return nil
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

	cols[models.CustomerColumns.UpdatedAt] = now

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return err
	}

	_,err = models.Transactions(models.CustomerWhere.ID.EQ(req.ID)).UpdateAll(ctx, tx, cols)

	if req.Amount != nil {
		tranx, err := models.FindTransaction(ctx, repo.DbConn, req.ID)
		if err != nil {
			_ = tx.Rollback()
			return err
		}

		diff := *req.Amount - tranx.Amount
		if tranx.TXType == TransactionType_Withdrawal.String() {
			diff *= -1
		}

		if diff != 0 {
			if _, err := models.Accounts(models.AccountWhere.ID.EQ(tranx.AccountID)).UpdateAll(ctx, tx, models.M{
				models.AccountColumns.Balance: fmt.Sprintf("%s + (%f)", models.AccountColumns.Balance, diff),
			}); err != nil {

				_ = tx.Rollback()
				return err
			}

			_, err = models.Transactions(
				models.TransactionWhere.CreatedAt.GT(tranx.CreatedAt)).
				UpdateAll(ctx, tx, models.M{models.TransactionColumns.Amount: fmt.Sprintf("amount + (%f)", diff)})

			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// Archive soft deleted the transaction from the database.
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

	tranx, err := models.Transactions(models.TransactionWhere.ID.EQ(req.ID)).One(ctx, tx)
	if err != nil {
		return err
	}

	_,err = models.Transactions(models.TransactionWhere.ID.EQ(req.ID)).UpdateAll(ctx, tx, models.M{models.TransactionColumns.ArchivedAt: now})

	var txAmount = tranx.Amount
	if tranx.TXType == TransactionType_Withdrawal.String() {
		txAmount *= -1
	}

	// updated all trans after it
	_, err = models.Transactions(
		models.TransactionWhere.CreatedAt.GT(tranx.CreatedAt)).
		UpdateAll(ctx, tx, models.M{models.TransactionColumns.Amount: fmt.Sprintf("amount + (%f)", txAmount)})

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// update all accounts balance
	_, err = models.Accounts(models.AccountWhere.ID.EQ(tranx.AccountID)).UpdateAll(ctx, tx,
		models.M{models.AccountColumns.Balance: fmt.Sprintf("balance + (%f)", txAmount)})

	if err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}
