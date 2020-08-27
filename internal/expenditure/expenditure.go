package expenditure

import (
	"context"
	"database/sql"
	"fmt"
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

// Find gets all the expenditures from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, claims auth.Claims, req FindRequest) (*PagedResponseList, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.expenditure.Find")
	defer span.Finish()

	var queries = []QueryMod{}

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	// if the current sales resp is not an admin, show only his transactions
	if !claims.HasRole(auth.RoleAdmin) {
		queries = append(queries, And(fmt.Sprintf("%s = '%s", models.RepsExpenseColumns.SalesRepID, claims.Subject)))
	}

	totalCount, err := models.RepsExpenses(queries...).Count(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get expenditure count")
	}

	if req.IncludeSalesRep {
		queries = append(queries, Load(models.RepsExpenseRels.SalesRep))
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

	slice, err := models.RepsExpenses(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return &PagedResponseList{}, nil
		}
		return nil, weberror.NewError(ctx, err, 500)
	}

	var result Expenditures
	for _, rec := range slice {
		result = append(result, FromModel(rec))
	}

	if len(result) == 0 {
		return &PagedResponseList{}, nil
	}

	return &PagedResponseList{
		Expenditures: result.Response(ctx),
		TotalCount:   totalCount,
	}, nil
}

// ReadByID gets the specified expenditure by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, _ auth.Claims, id string) (*Expenditure, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.expenditure.ReadByID")
	defer span.Finish()

	queries := []QueryMod{
		models.RepsExpenseWhere.ID.EQ(id),
		Load(models.RepsExpenseRels.SalesRep),
	}
	model, err := models.RepsExpenses(queries...).One(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	return FromModel(model), nil
}

// Create inserts a new expenditure into the database.
func (repo *Repository) Create(ctx context.Context, claims auth.Claims, req CreateRequest, now time.Time) (*Expenditure, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.expenditure.Create")
	defer span.Finish()
	if claims.Audience == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	// Admin users can update branch they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return nil, errors.WithStack(ErrForbidden)
	}

	// Validate the request.
	v := webcontext.Validator()
	err := v.StructCtx(ctx, req)
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

	salesRep, err := models.Users(models.UserWhere.PhoneNumber.EQ(req.SalesRepPhoneNumber)).One(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return nil, errors.New("Invalid phone number")
		}
		return nil, err
	}
	m := models.RepsExpense{
		ID:         uuid.NewRandom().String(),
		SalesRepID: salesRep.ID,
		Date:       now.Unix(),
		Amount:     req.Amount,
		Reason:     req.Reason,
	}

	if err := m.Insert(ctx, repo.DbConn, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "Insert expenditure failed")
	}

	return &Expenditure{
		ID:         m.ID,
		SalesRepID: req.SalesRepPhoneNumber,
		Date:       now,
		Amount:     req.Amount,
		Reason:     req.Reason,
	}, nil
}

// Update replaces an branch in the database.
func (repo *Repository) Update(ctx context.Context, claims auth.Claims, req UpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.expenditure.Update")
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
	if req.Amount != nil {
		cols[models.RepsExpenseColumns.Amount] = *req.Amount
	}
	if req.Reason != nil {
		cols[models.RepsExpenseColumns.Reason] = *req.Reason
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

	cols[models.BranchColumns.UpdatedAt] = now

	_, err = models.RepsExpenses(models.RepsExpenseWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols)

	return nil
}

// Delete removes an expenditure from the database.
func (repo *Repository) Delete(ctx context.Context, claims auth.Claims, req DeleteRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.expenditure.Delete")
	defer span.Finish()

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	// Ensure the claims can modify the project specified in the request.
	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Categories they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	_, err = models.RepsExpenses(models.RepsExpenseWhere.ID.EQ(req.ID)).DeleteAll(ctx, repo.DbConn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
