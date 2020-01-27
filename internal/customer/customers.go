package customer

import (
	"context"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
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

// Find gets all the customers from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, _ auth.Claims, req FindRequest) (Customers, error) {
	var queries = []QueryMod{
		Load(models.CustomerRels.SalesRep),
		Load(models.CustomerRels.Branch),
	}

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
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

	customerSlice, err := models.Customers(queries...).All(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	var result Customers
	for _, rec := range customerSlice {
		result = append(result, FromModel(rec))
	}

	return result, nil
}

// ReadByID gets the specified branch by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, claims auth.Claims, id string) (*Customer, error) {
	branchModel, err := models.FindCustomer(ctx, repo.DbConn, id)
	if err != nil {
		return nil, err
	}

	return FromModel(branchModel), nil
}

// Create inserts a new customer into the database.
func (repo *Repository) Create(ctx context.Context, claims auth.Claims, req CreateRequest, now time.Time) (*Customer, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.customer.Create")
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

	// If now empty set it to the current time.
	if now.IsZero() {
		now = time.Now()
	}

	// Always store the time as UTC.
	now = now.UTC()
	// Postgres truncates times to milliseconds when storing. We and do the same
	// here so the value we return is consistent with what we store.
	now = now.Truncate(time.Millisecond)
	m := models.Customer{
		ID:          uuid.NewRandom().String(),
		Email:       req.Email,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		SalesRepID:  req.SalesRepID,
		CreatedAt:   now,
		BranchID:    req.BranchID,
		UpdatedAt:   now,
	}

	if err := m.Insert(ctx, repo.DbConn, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "Insert customer failed")
	}

	return &Customer{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
		Address:     m.Address,
		SalesRepID:  m.SalesRepID,
		BranchID:    m.BranchID,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		ArchivedAt:  nil,
	}, nil
}

// Update replaces an customer in the database.
func (repo *Repository) Update(ctx context.Context, claims auth.Claims, req UpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.customer.Update")
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
	if req.Name != nil {
		cols[models.CustomerColumns.Name] = *req.Name
	}

	if req.Email != nil {
		cols[models.CustomerColumns.Email] = *req.Email
	}

	if req.Address != nil {
		cols[models.CustomerColumns.Address] = *req.Address
	}

	if req.PhoneNumber != nil {
		cols[models.CustomerColumns.PhoneNumber] = *req.PhoneNumber
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

	_,err = models.Customers(models.CustomerWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols)

	return nil
}

// Archive soft deleted the customer from the database.
func (repo *Repository) Archive(ctx context.Context, claims auth.Claims, req ArchiveRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.customer.Archive")
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

	_,err = models.Customers(models.CustomerWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, models.M{models.CustomerColumns.ArchivedAt: now})

	return nil
}

// Delete removes an customer from the database.
func (repo *Repository) Delete(ctx context.Context, claims auth.Claims, req DeleteRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.customer.Delete")
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

	_, err = models.Customers(models.CustomerWhere.ID.EQ(req.ID)).DeleteAll(ctx, repo.DbConn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
