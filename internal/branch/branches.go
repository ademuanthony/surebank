package branch

import (
	"context"
	"strings"
	"time"

	"merryworld/surebank/internal/dal"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/platform/web/weberror"
	"merryworld/surebank/internal/postgres/models"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

var (
	// ErrNotFound abstracts the postgres not found error.
	ErrNotFound = errors.New("Entity not found")

	// ErrForbidden occurs when a user tries to do something that is forbidden to them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

// Find gets all the branches from the database based on the request params.
func (repo *Repository) Find(ctx context.Context, req FindRequest) (Branches, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.branch.Find")
	defer span.Finish()

	queries := bson.M{}

	if !req.IncludeArchived {
		queries[dal.AccountColumns.ArchivedAt] = nil
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

	cursor, err := repo.mongoDb.Collection(dal.C.Branch).Find(ctx, bson.M{})
	if err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
	}
	defer cursor.Close(ctx)
	var result Branches
	for cursor.Next(ctx) {
		var c Branch
		cursor.Decode(&c)
		result = append(result, c)
	}
	if err := cursor.Err(); err != nil {
		return nil, weberror.WithMessage(ctx, err, "Cannot get customer list")
	}

	if len(result) == 0 {
		return Branches{}, nil
	}

	return result, nil
}

// ReadByID gets the specified branch by ID from the database.
func (repo *Repository) ReadByID(ctx context.Context, id string) (*Branch, error) {
	var rec Branch
	collection := repo.mongoDb.Collection(dal.C.Branch)
	err := collection.FindOne(ctx, bson.M{dal.BranchColumns.ID: id}).Decode(&rec)
	// if err != nil {
	// 	if err = repo.migrate(ctx); err != nil {
	// 		return nil, err
	// 	}
	// 	err = collection.FindOne(ctx, bson.M{dal.BranchColumns.ID: id}).Decode(&rec)
	// }
	return &rec, err
}

func (repo *Repository) branchExist(ctx context.Context, name string) bool {
	var rec Branch
	collection := repo.mongoDb.Collection(dal.C.Branch)
	_ = collection.FindOne(ctx, bson.M{dal.BranchColumns.Name: name}).Decode(&rec)
	return rec.ID != ""
}

// Create inserts a new checklist into the database.
func (repo *Repository) Create(ctx context.Context, claims auth.Claims, req CreateRequest, now time.Time) (*Branch, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.branch.Create")
	defer span.Finish()
	if claims.Audience == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	// Admin users can update branch they have access to.
	if !claims.HasRole(auth.RoleSuperAdmin) {
		return nil, errors.WithStack(ErrForbidden)
	}

	ctx = webcontext.ContextAddUniqueValue(ctx, req, "Name", !repo.branchExist(ctx, req.Name))
	// ctx = context.WithValue(ctx, webcontext.KeyTagUnique, !exists)

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
	m := Branch{
		ID:        uuid.NewRandom().String(),
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := repo.mongoDb.Collection(dal.C.Branch).InsertOne(ctx, m); err != nil {
		return nil, errors.WithMessage(err, "Insert branch failed")
	}

	return &m, nil
}

// Update replaces an branch in the database.
func (repo *Repository) Update(ctx context.Context, claims auth.Claims, req UpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.branch.Update")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update branches they have access to.
	if !claims.HasRole(auth.RoleSuperAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	var unique = true
	if req.Name != nil {
		unique = !repo.branchExist(ctx, *req.Name)
	}

	ctx = webcontext.ContextAddUniqueValue(ctx, req, "Name", unique)

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	cols := bson.M{}
	if req.Name != nil {
		cols[models.BrandColumns.Name] = *req.Name
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

	_, err = repo.mongoDb.Collection(dal.C.Branch).UpdateOne(ctx, bson.M{dal.BranchColumns.ID: req.ID}, cols)

	return err
}

// Archive soft deleted the checklist from the database.
func (repo *Repository) Archive(ctx context.Context, claims auth.Claims, req ArchiveRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.branch.Archive")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update branches they have access to.
	if !claims.HasRole(auth.RoleSuperAdmin) {
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

	cols := bson.M{dal.BranchColumns.ArchivedAt: now}

	_, err = repo.mongoDb.Collection(dal.C.Branch).UpdateOne(ctx, bson.M{dal.BranchColumns.ID: req.ID}, cols)

	return err
}

// Delete removes an checklist from the database.
func (repo *Repository) Delete(ctx context.Context, claims auth.Claims, req DeleteRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.branch.Delete")
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
	if !claims.HasRole(auth.RoleSuperAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	_, err = repo.mongoDb.Collection(dal.C.Branch).DeleteOne(ctx, bson.M{dal.BranchColumns.ID: req.ID})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (repo *Repository) Migrate(ctx context.Context) error {
	if c, _ := repo.mongoDb.Collection(dal.C.Branch).CountDocuments(ctx, bson.M{}); c > 0 {
		return nil
	}
	branches, err := models.Branches().All(ctx, repo.DbConn)
	if err != nil {
		return err
	}
	for _, b := range branches {
		if _, err := repo.mongoDb.Collection(dal.C.Branch).InsertOne(ctx, FromModel(b)); err != nil {
			return err
		}
	}
	return nil
}
