package account

import (
	"context"
	"database/sql"
	"math/rand"
	"strconv"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
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

// Find gets all the accounts from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, _ auth.Claims, req FindRequest) (*PagedResponseList, error) {
	var queries []QueryMod

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
	}

	totalCount, err := models.Accounts(queries...).Count(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get account total count")
	}

	if req.IncludeBranch {
		queries = append(queries, Load(models.AccountRels.Branch))
	}

	if req.IncludeCustomer {
		queries = append(queries, Load(models.AccountRels.Customer))
	}

	if req.IncludeSalesRep {
		queries = append(queries, Load(models.AccountRels.SalesRep))
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

	accountSlice, err := models.Accounts(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &PagedResponseList{}, nil
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	var result Accounts
	for _, rec := range accountSlice {
		result = append(result, FromModel(rec))
	}

	if len(result) == 0 {
		return &PagedResponseList{}, nil
	}

	return &PagedResponseList{
		Accounts:   result.Response(ctx),
		TotalCount: totalCount,
	}, nil
}

// FindDs gets all the accounts from the database that are of DS type and have > 0 balance.
func (repo *Repository) FindDs(ctx context.Context, _ auth.Claims, req FindRequest) (*PagedResponseList, error) {
	var queries []QueryMod

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	queries = append(queries,
		models.AccountWhere.AccountType.EQ(models.AccountTypeDS),
		models.AccountWhere.Balance.GT(0),
	)

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
	}

	totalCount, err := models.Accounts(queries...).Count(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get account total count")
	}

	if req.IncludeBranch {
		queries = append(queries, Load(models.AccountRels.Branch))
	}

	if req.IncludeCustomer {
		queries = append(queries, Load(models.AccountRels.Customer))
	}

	if req.IncludeSalesRep {
		queries = append(queries, Load(models.AccountRels.SalesRep))
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

	accountSlice, err := models.Accounts(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &PagedResponseList{}, nil
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	var result Accounts
	for _, rec := range accountSlice {
		result = append(result, FromModel(rec))
	}

	if len(result) == 0 {
		return &PagedResponseList{}, nil
	}

	return &PagedResponseList{
		Accounts:   result.Response(ctx),
		TotalCount: totalCount,
	}, nil
}

// FindDs gets all the accounts from the database that are of DS type and have > 0 balance.
func (repo *Repository) Debtors(ctx context.Context, _ auth.Claims, req FindRequest, currentDate time.Time) (*PagedResponseList, error) {
	var queries []QueryMod

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	threeDaysAgo := currentDate.Add(-3 * 24 * time.Hour).Unix()
	thirtyDaysAgo := currentDate.Add(-30 * 24 * time.Hour).Unix()

	queries = append(queries,
		models.AccountWhere.LastPaymentDate.GTE(thirtyDaysAgo),
		models.AccountWhere.LastPaymentDate.LTE(threeDaysAgo),
	)

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
	}

	totalCount, err := models.Accounts(queries...).Count(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get account total count")
	}

	if req.IncludeBranch {
		queries = append(queries, Load(models.AccountRels.Branch))
	}

	if req.IncludeCustomer {
		queries = append(queries, Load(models.AccountRels.Customer))
	}

	if req.IncludeSalesRep {
		queries = append(queries, Load(models.AccountRels.SalesRep))
	}

	queries = append(queries, qm.OrderBy(models.AccountColumns.LastPaymentDate))
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

	accountSlice, err := models.Accounts(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &PagedResponseList{}, nil
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	var result Accounts
	for _, rec := range accountSlice {
		result = append(result, FromModel(rec))
	}

	if len(result) == 0 {
		return &PagedResponseList{}, nil
	}

	return &PagedResponseList{
		Accounts:   result.Response(ctx),
		TotalCount: totalCount,
	}, nil
}

// ReadByID gets the specified branch by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, claims auth.Claims, id string) (*Account, error) {
	queries := []QueryMod{
		models.AccountWhere.ID.EQ(id),
		Load(models.AccountRels.Branch),
		Load(models.AccountRels.Customer),
		Load(models.AccountRels.SalesRep),
	}
	branchModel, err := models.Accounts(queries...).One(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, weberror.WithMessage(ctx, err, "Invalid account id")
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	return FromModel(branchModel), nil
}

func (repo *Repository) AccountsCount(ctx context.Context, claims auth.Claims) (int64, error) {
	var queries []QueryMod
	if !claims.HasRole(auth.RoleAdmin) {
		queries = append(queries, models.AccountWhere.SalesRepID.EQ(claims.Subject))
	}

	return models.Accounts(queries...).Count(ctx, repo.DbConn)
}

// Create inserts a new account into the database.
func (repo *Repository) Create(ctx context.Context, claims auth.Claims, req CreateRequest, now time.Time) (*Account, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.Create")
	defer span.Finish()
	if claims.Audience == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	salesRep, err := models.Users(models.UserWhere.ID.EQ(claims.Subject)).One(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, 400, "Something went wrong. Are you logged in?")
	}
	req.BranchID = salesRep.BranchID

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

	repo.accNumMtx.Lock()
	defer repo.accNumMtx.Unlock()

	m := models.Account{
		ID:          uuid.NewRandom().String(),
		Number:      repo.generateAccountNumber(ctx, req.Type),
		CustomerID:  req.CustomerID,
		AccountType: req.Type,
		Target:      req.Target,
		TargetInfo:  req.TargetInfo,
		SalesRepID:  claims.Subject,
		CreatedAt:   now.Unix(),
		BranchID:    req.BranchID,
		UpdatedAt:   now.Unix(),
	}

	if err := m.Insert(ctx, repo.DbConn, boil.Infer()); err != nil {
		return nil, weberror.WithMessage(ctx, err, "Insert account failed")
	}

	return &Account{
		ID:         m.ID,
		CustomerID: m.CustomerID,
		Number:     m.Number,
		Type:       m.AccountType,
		Target:     m.Target,
		TargetInfo: m.TargetInfo,
		SalesRepID: m.SalesRepID,
		BranchID:   m.BranchID,
		CreatedAt:  time.Unix(m.CreatedAt, 0),
		UpdatedAt:  time.Unix(m.UpdatedAt, 0),
	}, nil
}

func (repo *Repository) generateAccountNumber(ctx context.Context, accountType string) string {
	var accountNumber string
	for accountNumber == "" || repo.accountNumberExists(ctx, accountNumber) {
		accountNumber = accountType
		rand.Seed(time.Now().UTC().UnixNano())
		for i := 0; i < 5; i++ {
			accountNumber += strconv.Itoa(rand.Intn(10))
		}
	}
	return accountNumber
}

func (repo *Repository) accountNumberExists(ctx context.Context, number string) bool {
	exists, _ := models.Accounts(models.AccountWhere.Number.EQ(number)).Exists(ctx, repo.DbConn)
	return exists
}

// Update replaces an account in the database.
func (repo *Repository) Update(ctx context.Context, claims auth.Claims, req UpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.Update")
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

	cols := models.M{}
	if req.TargetInfo != nil {
		cols[models.AccountColumns.TargetInfo] = *req.TargetInfo
	}

	if req.Target != nil {
		cols[models.AccountColumns.Target] = *req.Target
	}

	if req.Type != nil {
		cols[models.AccountColumns.AccountType] = *req.Type
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

	cols[models.AccountColumns.UpdatedAt] = now.Unix()

	_, err = models.Accounts(models.AccountWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols)

	if err != nil {
		return weberror.NewError(ctx, err, 500)
	}

	return nil
}

// Archive soft deleted the account from the database.
func (repo *Repository) Archive(ctx context.Context, claims auth.Claims, req ArchiveRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.account.Archive")
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

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return err
	}

	if _, err := models.Transactions(models.TransactionWhere.AccountID.EQ(req.ID)).DeleteAll(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if _, err = models.Accounts(models.AccountWhere.ID.EQ(req.ID)).DeleteAll(ctx, tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (repo *Repository) DbConnCount(ctx context.Context) (int, error) {
	result := struct {
		Number int `json:"number"`
	}{}
	statement := "SELECT sum(numbackends) as number FROM pg_stat_database;"
	if err := models.NewQuery(qm.SQL(statement)).Bind(ctx, repo.DbConn, &result); err != nil {
		return 0, err
	}
	return result.Number, nil
}
