package shop

import (
	"context"
	"fmt"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/postgres/models"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"
	. "github.com/volatiletech/sqlboiler/v4/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Category
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (m Category) ToModel() models.Category {
	return models.Category{
		ID:   m.ID,
		Name: m.Name,
	}
}

func CategoryFromModel(category *models.Category) *Category {
	return &Category{
		ID:   category.ID,
		Name: category.Name,
	}
}

// CategoryResponse represent a Category that is returned for display
type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Response transforms Category to CategoryResponse that is used for display.
func (m *Category) Response(_ context.Context) *CategoryResponse {
	if m == nil {
		return nil
	}

	r := &CategoryResponse{
		ID:   m.ID,
		Name: m.Name,
	}

	return r
}

// Categories a list of Categories.
type Categories []*Category

// Response transforms a list of Categories to a list of CategoryResponse.
func (m *Categories) Response(ctx context.Context) []*CategoryResponse {
	var l []*CategoryResponse
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// CategoryCreateRequest contains information needed to create a new Category.
type CategoryCreateRequest struct {
	Name string `json:"name" validate:"required,unique"  example:"Registration"`
}

// CategoryReadRequest defines the information needed to read a Category.
type CategoryReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// CategoryUpdateRequest defines what information may be provided to modify an existing
// Category. All fields are optional so clients can send just the fields they want
// changed.
type CategoryUpdateRequest struct {
	ID   string  `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name *string `json:"name" validate:"omitempty,unique"  example:"Registration"`
}

// CategoryDeleteRequest defines the information needed to delete a Category.
type CategoryDeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// ServiceTypeFindRequest defines the possible options to search for ServiceTypes. By default
// archived ServiceType will be excluded from response.
type CategoryFindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}

func (repo Repository) FindCategory(ctx context.Context, req CategoryFindRequest) (Categories, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.FindCategory")
	defer span.Finish()

	var queries []QueryMod

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if req.Limit != nil {
		queries = append(queries, Limit(int(*req.Limit)))
	}

	if req.Offset != nil {
		queries = append(queries, Offset(int(*req.Offset)))
	}

	CategorySlice, err := models.Categories(queries...).All(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	var result Categories
	for _, rec := range CategorySlice {
		result = append(result, CategoryFromModel(rec))
	}

	return result, nil
}

// ReadCategoryByID gets the specified Category by ID from the database.
func (repo *Repository) ReadCategoryByID(ctx context.Context, _ auth.Claims, id string) (*Category, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.ReadCategoryByID")
	defer span.Finish()

	categoryModel, err := models.FindCategory(ctx, repo.DbConn, id)
	if err != nil {
		return nil, fmt.Errorf("%s %s", err.Error(), id)
	}

	return &Category{
		ID:   categoryModel.ID,
		Name: categoryModel.Name,
	}, nil
}

func (repo *Repository) CreateCategory(ctx context.Context, claims auth.Claims, req CategoryCreateRequest) (*Category, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.CreateCategory")
	defer span.Finish()
	if claims.Audience != "" {
		// Admin users can update Categories they have access to.
		if !claims.HasRole(auth.RoleAdmin) {
			return nil, errors.WithStack(ErrForbidden)
		}
	}

	exists, err := models.Categories(models.CategoryWhere.Name.EQ(req.Name)).Exists(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	ctx = webcontext.ContextAddUniqueValue(ctx, req, "Name", !exists)

	// Validate the request.
	v := webcontext.Validator()
	err = v.StructCtx(ctx, req)
	if err != nil {
		return nil, err
	}

	s := Category{
		ID:   uuid.NewRandom().String(),
		Name: req.Name,
	}

	catModel := s.ToModel()
	if err := catModel.Insert(ctx, repo.DbConn, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "create Category failed")
	}

	return &s, nil
}

// UpdateCategory replaces an project in the database.
func (repo *Repository) UpdateCategory(ctx context.Context, claims auth.Claims, req CategoryUpdateRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.UpdateCategory")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Categories they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	unique := true
	if req.Name != nil {
		exists, err := models.Categories(models.CategoryWhere.Name.EQ(*req.Name), models.CategoryWhere.ID.NEQ(req.ID)).Exists(ctx, repo.DbConn)
		if err != nil {
			return err
		}

		unique = !exists
	}

	ctx = webcontext.ContextAddUniqueValue(ctx, req, "Name", unique)

	// Validate the request.
	v := webcontext.Validator()
	err := v.StructCtx(ctx, req)
	if err != nil {
		return err
	}

	cols := models.M{}
	if req.Name != nil {
		cols[models.CategoryColumns.Name] = *req.Name
	}

	if len(cols) == 0 {
		return nil
	}

	_, err = models.Categories(models.CategoryWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols)

	return nil
}

// Delete removes an Category from the database.
func (repo *Repository) DeleteCategory(ctx context.Context, claims auth.Claims, req CategoryDeleteRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.DeleteCategory")
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

	_, err = models.Categories(models.CategoryWhere.ID.EQ(req.ID)).DeleteAll(ctx, repo.DbConn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
