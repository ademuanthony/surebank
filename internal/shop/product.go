package shop

import (
	"context"
	"database/sql"
	"time"

	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/postgres/models"

	"github.com/pborman/uuid"
	"github.com/pkg/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	. "github.com/volatiletech/sqlboiler/queries/qm"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// Product
type Product struct {
	ID           string      `json:"id"`
	BrandID      string      `json:"brand_id"`
	CategoryID   string      `json:"category_id"`
	Brand        string      `json:"brand"`
	Category     string      `json:"category"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Sku          string      `json:"sku"`
	Barcode      string      `json:"barcode"`
	Price        float64     `json:"price"`
	StockBalance int         `json:"stock_balance"`
	ReorderLevel int         `json:"reorder_level"`
	Image        string      `json:"image"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	ArchivedAt   null.Time   `json:"archived_at"`
	CreatedByID  string      `json:"created_by_id"`
	CreatedBy    string      `json:"created_by"`
	UpdatedByID  string      `json:"updated_by_id"`
	UpdatedBy    string      `json:"updated_by"`
	ArchivedByID null.String `json:"archived_by_id"`
}

func (m Product) ToModel() models.Product {
	return models.Product{
		ID:           m.ID,
		BrandID:      null.StringFrom(m.BrandID),
		CategoryID:   m.CategoryID,
		Name:         m.Name,
		Description:  m.Description,
		Sku:          m.Sku,
		Barcode:      m.Barcode,
		Price:        m.Price,
		ReorderLevel: m.ReorderLevel,
		Image:        null.StringFrom(m.Image),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		ArchivedAt:   m.ArchivedAt,
		CreatedByID:  m.CreatedByID,
		UpdatedByID:  m.UpdatedByID,
		ArchivedByID: m.ArchivedByID,
	}
}

func ProductFromModel(product *models.Product) *Product {
	p := &Product{
		ID:           product.ID,
		BrandID:      product.BrandID.String,
		CategoryID:   product.CategoryID,
		Name:         product.Name,
		Description:  product.Description,
		Sku:          product.Sku,
		Barcode:      product.Barcode,
		Price:        product.Price,
		StockBalance: product.StockBalance,
		ReorderLevel: product.ReorderLevel,
		Image:        product.Image.String,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
		ArchivedAt:   product.ArchivedAt,
		CreatedByID:  product.CreatedByID,
		UpdatedByID:  product.UpdatedByID,
	}

	if product.R != nil {
		if product.R.Brand != nil {
			p.Brand = product.R.Brand.Name
		}

		if product.R.Category != nil {
			p.Category = product.R.Category.Name
		}
	}

	return p
}

// ProductResponse represent a Product that is returned for display
type ProductResponse struct {
	ID           string  `json:"id"`
	BrandID      string  `json:"brand_id"`
	Brand        string  `json:"brand"`
	CategoryID   string  `json:"category_id"`
	Category     string  `json:"category"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Sku          string  `json:"sku"`
	Barcode      string  `json:"barcode"`
	Price        float64 `json:"price"`
	StockBalance int     `json:"stock_balance"`
	ReorderLevel int     `json:"reorder_level"`
	Image        string  `json:"image"`

	CreatedAt   web.TimeResponse `json:"created_at"`
	UpdatedAt   web.TimeResponse `json:"updated_at"`
	CreatedByID string           `json:"created_by_id"`
	UpdatedByID string           `json:"updated_by_id"`
}

// Response transforms Product to ProductResponse that is used for display.
func (m *Product) Response(ctx context.Context) *ProductResponse {
	if m == nil {
		return nil
	}

	r := &ProductResponse{
		ID:           m.ID,
		BrandID:      m.BrandID,
		Brand:        m.Brand,
		CategoryID:   m.CategoryID,
		Category:     m.Category,
		Name:         m.Name,
		Description:  m.Description,
		Sku:          m.Sku,
		Barcode:      m.Barcode,
		Price:        m.Price,
		StockBalance: m.StockBalance,
		ReorderLevel: m.ReorderLevel,
		Image:        m.Image,
		CreatedByID:  m.CreatedByID,
		UpdatedByID:  m.UpdatedByID,
		CreatedAt:    web.NewTimeResponse(ctx, m.CreatedAt),
		UpdatedAt:    web.NewTimeResponse(ctx, m.UpdatedAt),
	}

	return r
}

// Products a list of Products.
type Products []*Product

// Response transforms a list of Products to a list of ProductResponse.
func (m *Products) Response(ctx context.Context) []*ProductResponse {
	var l []*ProductResponse
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// ProductCreateRequest contains information needed to create a new Product.
type ProductCreateRequest struct {
	Name         string   `json:"name" validate:"required,unique"  example:"Bread"`
	BrandID      string   `json:"brand_id" validate:"required"`
	CategoryID   string   `json:"category_id" validate:"required"`
	Description  string   `json:"description"`
	Sku          string   `json:"sku" validate:"required,unique"`
	Barcode      string   `json:"barcode" validate:"required,unique"`
	Price        float64  `json:"price" validate:"required"`
	ReorderLevel int      `json:"reorder_level"`
	Image        string   `json:"image"`
	Categories   []string `json:"categories"`
}

// ProductReadRequest defines the information needed to read a Product.
type ProductReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// ProductUpdateRequest defines what information may be provided to modify an existing
// Product. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank.
type ProductUpdateRequest struct {
	ID           string   `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name         *string  `json:"name" example:"Bread"`
	BrandID      *string  `json:"brand_id"`
	CategoryID   *string  `json:"category_id"`
	Description  *string  `json:"description"`
	Sku          *string  `json:"sku"`
	Barcode      *string  `json:"barcode"`
	Price        *float64 `json:"price"`
	ReorderLevel *int     `json:"reorder_level"`
	Image        *string  `json:"image"`
	Categories   []string `json:"categories"`
}

// AddProductToCategoryRequest contains the information needed to link a product to a category
type AddProductToCategoryRequest struct {
	ProductID  string `json:"product_id" validate:"required"`
	CategoryID string `json:"category_id" validate:"required"`
}

// RemoveProductFromCategoryRequest contains the information needed to remove a product from a category
type RemoveProductFromCategoryRequest struct {
	ProductID  string `json:"product_id" validate:"required"`
	CategoryID string `json:"category_id" validate:"required"`
}

// ProductArchiveRequest defines the information needed to archive a Product. This will archive (soft-delete) the
// existing database entry.
type ProductArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// ProductDeleteRequest defines the information needed to delete a Product.
type ProductDeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// ServiceTypeFindRequest defines the possible options to search for ServiceTypes. By default
// archived ServiceType will be excluded from response.
type ProductFindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}

// FindProduct gets all the products from the database based on the request params.
func (repo Repository) FindProduct(ctx context.Context, req ProductFindRequest) (Products, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.FindProduct")
	defer span.Finish()
	var queries []QueryMod

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
	}

	if req.Limit != nil {
		queries = append(queries, Limit(int(*req.Limit)))
	}

	if req.Offset != nil {
		queries = append(queries, Offset(int(*req.Offset)))
	}

	ProductSlice, err := models.Products(queries...).All(ctx, repo.DbConn)
	if err != nil {
		if err.Error() == sql.ErrNoRows.Error() {
			return Products{}, nil
		}
		return nil, err
	}

	var result Products
	for _, rec := range ProductSlice {
		result = append(result, ProductFromModel(rec))
	}

	return result, nil
}

// ReadProductByID gets the specified Product by ID from the database.
func (repo *Repository) ReadProductByID(ctx context.Context, _ auth.Claims, id string) (*Product, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.ReadProductByID")
	defer span.Finish()
	queries := []QueryMod{
		models.ProductWhere.ID.EQ(id),
		Load(models.ProductRels.Brand),
		Load(models.ProductRels.Category),
	}
	productModel, err := models.Products(queries...).One(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	return ProductFromModel(productModel), nil
}

func (repo *Repository) ReadProductByIDTx(ctx context.Context, _ auth.Claims, id string, tx *sql.Tx) (*Product, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.ReadProductByIDTx")
	defer span.Finish()
	queries := []QueryMod{
		models.ProductWhere.ID.EQ(id),
		Load(models.ProductRels.Brand),
		Load(models.ProductRels.Category),
	}
	productModel, err := models.Products(queries...).One(ctx, tx)
	if err != nil {
		return nil, err
	}

	return ProductFromModel(productModel), nil
}

func (repo *Repository) CreateProduct(ctx context.Context, claims auth.Claims, req ProductCreateRequest, now time.Time) (*Product, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.CreateProduct")
	defer span.Finish()
	if claims.Audience == "" {
		return nil, errors.WithStack(ErrForbidden)
	}

	if !claims.HasRole(auth.RoleAdmin) {
		return nil, errors.WithStack(ErrForbidden)
	}

	exist, err := models.Products(models.ProductWhere.Name.EQ(req.Name)).Exists(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}
	ctx = webcontext.ContextAddUniqueValue(ctx, req, "Name", !exist)

	exist, err = models.Products(models.ProductWhere.Sku.EQ(req.Sku)).Exists(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}
	ctx = webcontext.ContextAddUniqueValue(ctx, req, "Sku", !exist)

	exist, err = models.Products(models.ProductWhere.Barcode.EQ(req.Barcode)).Exists(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}
	ctx = webcontext.ContextAddUniqueValue(ctx, req, "Barcode", !exist)

	// Validate the request.
	if err := webcontext.Validator().StructCtx(ctx, req); err != nil {
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

	s := Product{
		ID:           uuid.NewRandom().String(),
		BrandID:      req.BrandID,
		CategoryID:   req.CategoryID,
		Name:         req.Name,
		Description:  req.Description,
		Sku:          req.Sku,
		Barcode:      req.Barcode,
		Price:        req.Price,
		ReorderLevel: req.ReorderLevel,
		Image:        req.Image,
		CreatedAt:    now,
		UpdatedAt:    now,
		CreatedByID:  claims.Subject,
		UpdatedByID:  claims.Subject,
	}

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "create product failed, cannot start db transaction"))
	}

	prodModel := s.ToModel()
	if err := prodModel.Insert(ctx, tx, boil.Infer()); err != nil {
		_ = tx.Rollback()
		return nil, errors.WithMessage(err, "create Product failed")
	}

	for _, categoryId := range req.Categories {
		pCat := models.ProductCategory{
			ID:         uuid.NewRandom().String(),
			ProductID:  s.ID,
			CategoryID: categoryId,
		}
		if err = pCat.Insert(ctx, tx, boil.Infer()); err != nil {
			_ = tx.Rollback()
			return nil, errors.WithStack(errors.WithMessage(err, "create product failed, cannot link product category"))
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "create product failed, cannot commit db transaction"))
	}

	return &s, nil
}

// UpdateProduct replaces an project in the database.
func (repo *Repository) UpdateProduct(ctx context.Context, claims auth.Claims, req ProductUpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.UpdateProduct")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Products they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}
	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	if req.Name != nil {
		if e, _ := models.Products(models.ProductWhere.Name.EQ(*req.Name), models.ProductWhere.ID.NEQ(req.ID)).Exists(ctx, repo.DbConn); e {
			return errors.WithStack(errors.New("the specified name already exists"))
		}
	}

	if req.Sku != nil {
		if e, _ := models.Products(models.ProductWhere.Sku.EQ(*req.Sku), models.ProductWhere.ID.NEQ(req.ID)).Exists(ctx, repo.DbConn); e {
			return errors.WithStack(errors.New("the specified sku already exists"))
		}
	}

	if req.Barcode != nil {
		if e, _ := models.Products(models.ProductWhere.Barcode.EQ(*req.Barcode), models.ProductWhere.ID.NEQ(req.ID)).Exists(ctx, repo.DbConn); e {
			return errors.WithStack(errors.New("the specified barcode already exists"))
		}
	}

	cols := models.M{}
	if req.Name != nil {
		cols[models.ProductColumns.Name] = *req.Name
	}

	if req.Barcode != nil {
		cols[models.ProductColumns.Barcode] = *req.Barcode
	}

	if req.BrandID != nil {
		cols[models.ProductColumns.BrandID] = *req.BrandID
	}

	if req.CategoryID != nil {
		cols[models.ProductColumns.CategoryID] = *req.CategoryID
	}

	if req.Description != nil {
		cols[models.ProductColumns.Description] = *req.Description
	}

	if req.Price != nil {
		cols[models.ProductColumns.Price] = *req.Price
	}

	if req.ReorderLevel != nil {
		cols[models.ProductColumns.ReorderLevel] = *req.ReorderLevel
	}

	if req.Image != nil {
		cols[models.ProductColumns.Image] = *req.Image
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

	cols[models.ProductColumns.UpdatedByID] = claims.Subject
	cols[models.ProductColumns.UpdatedAt] = now

	if len(cols) == 0 {
		return nil
	}

	_, err = models.Products(models.ProductWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols)

	return nil
}

func (repo *Repository) AddProductToCategory(ctx context.Context, claims auth.Claims, req AddProductToCategoryRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.AddProductToCategory")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Products they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}
	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	if exists, _ := models.Products(models.ProductWhere.ID.EQ(req.ProductID)).Exists(ctx, repo.DbConn); !exists {
		return errors.New("invalid product ID")
	}

	if exists, _ := models.Categories(models.CategoryWhere.ID.EQ(req.CategoryID)).Exists(ctx, repo.DbConn); !exists {
		return errors.New("invalid category ID")
	}

	if exists, _ := models.ProductCategories(
		models.ProductCategoryWhere.ProductID.EQ(req.ProductID),
		models.ProductCategoryWhere.CategoryID.EQ(req.CategoryID),
	).Exists(ctx, repo.DbConn); !exists {
		return errors.New("product already linked to the selected category")
	}

	pCat := models.ProductCategory{
		ID:         uuid.NewRandom().String(),
		ProductID:  req.ProductID,
		CategoryID: req.CategoryID,
	}

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "cannot link product, error in starting db transaction"))
	}

	if err := pCat.Insert(ctx, tx, boil.Infer()); err != nil {
		_ = tx.Rollback()
		return errors.WithStack(err)
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

	cols := models.M{models.ProductColumns.UpdatedAt: now, models.ProductColumns.UpdatedByID: claims.Subject}
	if _, err := models.Products(models.ProductWhere.ID.EQ(req.ProductID)).UpdateAll(ctx, tx, cols); err != nil {
		_ = tx.Rollback()
		return errors.WithStack(errors.WithMessage(err, "cannot link product, cannot update product"))
	}

	if err = tx.Commit(); err != nil {
		return errors.WithStack(errors.WithMessage(err, "cannot link product, cannot commit db transaction"))
	}

	return nil
}

func (repo *Repository) RemoveProductFromCategory(ctx context.Context, claims auth.Claims, req RemoveProductFromCategoryRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.RemoveProductFromCategory")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Products they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}
	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	pCat, err := models.ProductCategories(
		models.ProductCategoryWhere.ProductID.EQ(req.ProductID),
		models.ProductCategoryWhere.CategoryID.EQ(req.CategoryID)).One(ctx, repo.DbConn)
	if err != nil {
		return nil
	}

	tx, err := repo.DbConn.Begin()
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "cannot remove category, error in starting db transaction"))
	}

	if _, err = pCat.Delete(ctx, tx); err != nil {
		_ = tx.Rollback()
		return errors.WithStack(errors.WithMessage(err, "cannot remove category"))
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

	cols := models.M{models.ProductColumns.UpdatedAt: now, models.ProductColumns.UpdatedByID: claims.Subject}
	if _, err := models.Products(models.ProductWhere.ID.EQ(req.ProductID)).UpdateAll(ctx, tx, cols); err != nil {
		_ = tx.Rollback()
		return errors.WithStack(errors.WithMessage(err, "cannot link product, cannot update product"))
	}

	if err = tx.Commit(); err != nil {
		return errors.WithStack(errors.WithMessage(err, "cannot link product, cannot commit db transaction"))
	}

	return nil
}

// Archive soft deletes the product from the database.
func (repo *Repository) ArchiveProduct(ctx context.Context, claims auth.Claims, req ProductArchiveRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.ArchiveProduct")
	defer span.Finish()

	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Products they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
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

	cols := models.M{models.ProductColumns.ArchivedByID: claims.Subject, models.ProductColumns.ArchivedAt: now}
	if _, err := models.Products(models.ProductWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// Delete removes an Product from the database.
func (repo *Repository) DeleteProduct(ctx context.Context, claims auth.Claims, req ProductDeleteRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.DeleteProduct")
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
	// Admin users can update Products they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	_, err = models.Products(models.ProductWhere.ID.EQ(req.ID)).DeleteAll(ctx, repo.DbConn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
