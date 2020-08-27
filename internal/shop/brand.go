package shop

import (
	"context"
	"database/sql"
	"time"

	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/postgres/models"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Brand
type Brand struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

func (m Brand) toModel() models.Brand {
	return models.Brand{
		ID:   m.ID,
		Name: m.Name,
	}
}

func brandFromModel(brand *models.Brand) *Brand {
	return &Brand{
		ID:   brand.ID,
		Name: brand.Name,
		Logo: brand.Logo,
	}
}

// BrandResponse represent a brand that is returned for display
type BrandResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Logo string `json:"logo"`
}

// Response transforms Brand to BrandResponse that is used for display.
func (m *Brand) Response(ctx context.Context) *BrandResponse {
	if m == nil {
		return nil
	}

	r := &BrandResponse{
		ID:   m.ID,
		Name: m.Name,
		Logo: m.Logo,
	}

	return r
}

// Brands a list of Brands.
type Brands []*Brand

// Response transforms a list of Brands to a list of BrandResponse.
func (m *Brands) Response(ctx context.Context) []*BrandResponse {
	var l []*BrandResponse
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// BrandCreateRequest contains information needed to create a new Brand.
type BrandCreateRequest struct {
	Name string `json:"name" validate:"required,unique"  example:"Registration"`
	Logo string `json:"logo"`
}

// BrandReadRequest defines the information needed to read a Brand.
type BrandReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// BrandUpdateRequest defines what information may be provided to modify an existing
// Brand. All fields are optional so clients can send just the fields they want
// changed.
type BrandUpdateRequest struct {
	ID   string  `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name *string `json:"name" validate:"omitempty,unique"  example:"Registration"`
	Logo *string `json:"logo"`
}

// BrandArchiveRequest defines the information needed to archive a Brand. This will archive (soft-delete) the
// existing database entry.
type BrandArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// BrandDeleteRequest defines the information needed to delete a Brand.
type BrandDeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// ServiceTypeFindRequest defines the possible options to search for ServiceTypes. By default
// archived ServiceType will be excluded from response.
type BrandFindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}

func (repo Repository) FindBrand(ctx context.Context, req BrandFindRequest) (Brands, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.FindBranch")
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

	brandSlice, err := models.Brands(queries...).All(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	var result Brands
	for _, rec := range brandSlice {
		result = append(result, brandFromModel(rec))
	}

	return result, nil
}

// ReadBrandByID gets the specified brand by ID from the database.
func (repo *Repository) ReadBrandByID(ctx context.Context, _ auth.Claims, id string) (*Brand, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.ReadBrandByID")
	defer span.Finish()

	brandModel, err := models.FindBrand(ctx, repo.DbConn, id)
	if err != nil {
		return nil, err
	}

	return &Brand{
		ID:   brandModel.ID,
		Name: brandModel.Name,
		Logo: brandModel.Logo,
	}, nil
}

func (repo *Repository) CreateBrand(ctx context.Context, claims auth.Claims, req BrandCreateRequest, now time.Time) (*Brand, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.CreateBrand")
	defer span.Finish()
	if claims.Audience != "" {
		// Admin users can update brands they have access to.
		if !claims.HasRole(auth.RoleAdmin) {
			return nil, errors.WithStack(ErrForbidden)
		}
	}

	exists, err := models.Brands(models.BrandWhere.Name.EQ(req.Name)).Exists(ctx, repo.DbConn)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return nil, err
	}

	ctx = webcontext.ContextAddUniqueValue(ctx, req, "Name", !exists)

	// Validate the request.
	v := webcontext.Validator()
	err = v.StructCtx(ctx, req)
	if err != nil {
		return nil, err
	}

	s := models.Brand{
		ID:   uuid.NewRandom().String(),
		Name: req.Name,
	}

	if err := s.Insert(ctx, repo.DbConn, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "create brand failed")
	}

	return &Brand{
		ID:   s.ID,
		Name: req.Name,
	}, nil
}

// UpdateBrand replaces an project in the database.
func (repo *Repository) UpdateBrand(ctx context.Context, claims auth.Claims, req BrandUpdateRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.UpdateBrand")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update brands they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	uniq := true
	if req.Name != nil {
		exists, err := models.Brands(models.BrandWhere.Name.EQ(*req.Name), models.BrandWhere.ID.NEQ(req.ID)).Exists(ctx, repo.DbConn)
		if err != nil {
			return err
		}

		uniq = !exists
	}

	ctx = webcontext.ContextAddUniqueValue(ctx, req, "Name", uniq)

	// Validate the request.
	v := webcontext.Validator()
	err := v.StructCtx(ctx, req)
	if err != nil {
		return err
	}

	cols := models.M{}
	if req.Name != nil {
		cols[models.BrandColumns.Name] = *req.Name
	}

	if req.Logo != nil {
		cols[models.BrandColumns.Logo] = *req.Logo
	}

	if len(cols) == 0 {
		return nil
	}

	_, err = models.Brands(models.BrandWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols)

	return nil
}

// Delete removes an brand from the database.
func (repo *Repository) DeleteBrand(ctx context.Context, claims auth.Claims, req BrandDeleteRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.DeleteBrand")
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
	// Admin users can update brands they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	_, err = models.Brands(models.BrandWhere.ID.EQ(req.ID)).DeleteAll(ctx, repo.DbConn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
