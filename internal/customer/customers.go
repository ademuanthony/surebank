package customer

import (
	"context"
	"strings"
	"time"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"merryworld/surebank/internal/branch"
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

// Find gets all the customers from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, _ auth.Claims, req FindRequest) (*PagedResponseList, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.customer.Find")
	defer span.Finish()

	collection := repo.mongoDb.Collection(CollectionName)

	var queries = bson.D{}
	if !req.IncludeArchived {
		queries = append(queries, primitive.E{Key: Columns.ArchivedAt, Value: nil})
	}

	totalCount, err := collection.CountDocuments(ctx, queries)
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer count")
	}

	findOptions := options.Find()

	sort := bson.D{}
	if len(req.Order) > 0 {
		for _, s := range req.Order {
			sortInfo := strings.Split(s, " ")
			if len(sortInfo) != 2 {
				continue
			}
			sort = append(sort, primitive.E{Key: sortInfo[0], Value: sortInfo[1]})
		}
	}
	findOptions.SetSort(sort)

	if req.Limit != nil {
		findOptions.SetLimit(int64(*req.Limit))
	}

	if req.Offset != nil {
		findOptions.Skip = req.Offset
	}

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
	}
	defer cursor.Close(ctx)
	var result Customers
	for cursor.Next(ctx) {
		var c Customer
		cursor.Decode(&c)
		result = append(result, c)
	}
	if err := cursor.Err(); err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
	}

	if len(result) == 0 {
		return &PagedResponseList{}, nil
	}

	return &PagedResponseList{
		Customers:  result.Response(ctx),
		TotalCount: totalCount,
	}, nil
}

func (repo *Repository) ReadByID(ctx context.Context, id string) (*Customer, error) {
	var rec Customer
	collection := repo.mongoDb.Collection(CollectionName)
	err := collection.FindOne(ctx, bson.M{Columns.ID: id}).Decode(&rec)
	return &rec, err
}

func (repo *Repository) CustomersCount(ctx context.Context, claims auth.Claims) (int64, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.customer.CustomersCount")
	defer span.Finish()

	query := bson.M{}
	if !claims.HasRole(auth.RoleAdmin) {
		query[Columns.SalesRepID] = claims.Subject
	}

	collection := repo.mongoDb.Collection(CollectionName)
	return collection.CountDocuments(ctx, query)
}

// Create inserts a new customer into the database.
func (repo *Repository) Create(ctx context.Context, claims auth.Claims, req CreateRequest, now time.Time) (*Customer, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.customer.Create")
	defer span.Finish()
	if claims.Audience == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	salesRep, err := models.Users(models.UserWhere.ID.EQ(claims.Subject)).One(ctx, repo.DbConn)
	if err != nil {
		return nil, weberror.NewErrorMessage(ctx, err, 400, "Something went wrong. Are you logged in?")
	}
	req.SalesRepID = salesRep.ID
	req.BranchID = salesRep.BranchID

	branch, err := branch.ReadByID(ctx, repo.mongoDb, req.BranchID)
	if err != nil {
		return nil, err
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
	m := Customer{
		ID:          uuid.NewRandom().String(),
		Email:       req.Email,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		SalesRepID:  req.SalesRepID,
		SalesRep:    salesRep.FirstName + " " + salesRep.LastName,
		CreatedAt:   now,
		BranchID:    req.BranchID,
		Branch:      branch.Name,
		UpdatedAt:   now,
	}

	if _, err := repo.mongoDb.Collection(CollectionName).InsertOne(ctx, m); err != nil {
		return nil, weberror.WithMessage(ctx, err, "Insert customer failed")
	}

	return &m, nil
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

	cols := bson.M{}
	if req.Name != nil {
		cols[Columns.Name] = *req.Name
	}

	if req.Email != nil {
		cols[Columns.Email] = *req.Email
	}

	if req.Address != nil {
		cols[Columns.Address] = *req.Address
	}

	if req.PhoneNumber != nil {
		cols[Columns.PhoneNumber] = *req.PhoneNumber
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

	cols[Columns.UpdatedAt] = now.Unix()

	collection := repo.mongoDb.Collection(CollectionName)
	_, err = collection.UpdateOne(ctx, bson.M{Columns.ID: req.ID}, bson.M{"$set": cols})

	return err
}

// Archive soft deleted the customer from the database.
func (repo *Repository) Archive(ctx context.Context, claims auth.Claims, req ArchiveRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.customer.Archive")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update customer they have access to.
	if !claims.HasRole(auth.RoleSuperAdmin) {
		return errors.WithStack(ErrForbidden)
	}
	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	collection := repo.mongoDb.Collection(CollectionName)
	_, err = collection.DeleteOne(ctx, bson.M{Columns.ID: req.ID})
	return err
}
