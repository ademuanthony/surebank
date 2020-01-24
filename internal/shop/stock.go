package shop

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
	"github.com/volatiletech/null"
	"merryworld/surebank/internal/platform/auth"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/platform/web/webcontext"
	"merryworld/surebank/internal/postgres/models"
)

// Stock
type Stock struct {
	ID               string     `json:"id"`
	BranchID		 string		`json:"branch_id"`
	BatchNumber      string     `json:"batch_number"`
	ProductID        string     `json:"product_id"`
	ProductName      string     `json:"product_name"`
	UnitCostPrice    float64    `json:"unit_cost_price"`
	Quantity         int        `json:"quantity"`
	DeductedQuantity int        `json:"deducted_quantity"`
	ManufactureDate  *time.Time `json:"manufacture_date"`
	ExpiryDate       *time.Time `json:"expiry_date"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	ArchivedAt       *time.Time `json:"archived_at"`
	CreatedByID      string     `json:"created_by_id"`
	UpdatedByID      string     `json:"updated_by_id"`
	ArchivedByID     *string    `json:"archived_by"`
}

func (s Stock) model() models.Stock {
	m := models.Stock{
		ID:               s.ID,
		BranchID: 		  s.BranchID,
		BatchNumber:      s.BatchNumber,
		ProductID:        s.ProductID,
		UnitCostPrice:    s.UnitCostPrice,
		Quantity:         s.Quantity,
		DeductedQuantity: s.DeductedQuantity,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
		CreatedByID:      s.CreatedByID,
		UpdatedByID:      s.UpdatedByID,
	}

	if s.ManufactureDate != nil {
		m.ManufactureDate = null.TimeFrom(*s.ManufactureDate)
	}

	if s.ExpiryDate != nil {
		m.ExpiryDate = null.TimeFrom(*s.ExpiryDate)
	}

	if s.ArchivedAt != nil {
		m.ArchivedAt = null.TimeFrom(*s.ArchivedAt)
	}

	if s.ArchivedByID != nil {
		m.ArchivedByID = null.StringFrom(*s.ArchivedByID)
	}

	return m
}

func stockFromModel(stock *models.Stock) *Stock {
	s := &Stock{
		ID:               stock.ID,
		BatchNumber:      stock.BatchNumber,
		ProductID:        stock.ProductID,
		UnitCostPrice:    stock.UnitCostPrice,
		Quantity:         stock.Quantity,
		DeductedQuantity: stock.DeductedQuantity,
		CreatedAt:        stock.CreatedAt,
		UpdatedAt:        stock.UpdatedAt,
		CreatedByID:      stock.CreatedByID,
		UpdatedByID:      stock.UpdatedByID,
	}

	if stock.ManufactureDate.Valid {
		s.ManufactureDate = &stock.ManufactureDate.Time
	}

	if stock.ExpiryDate.Valid {
		s.ExpiryDate = &stock.ExpiryDate.Time
	}

	if stock.ArchivedByID.Valid {
		s.ArchivedByID = &stock.ArchivedByID.String
	}

	if stock.ArchivedAt.Valid {
		s.ArchivedAt = &stock.ArchivedAt.Time
	}

	if stock.R != nil {
		if stock.R.Product != nil {
			s.ProductName = stock.R.Product.Name
		}
	}
	
	return s
}

// StockResponse represent a stock that is returned for display
type StockResponse struct {
	ID              string            `json:"id"`
	BatchNumber     string            `json:"batch_number"`
	ProductID       string            `json:"product_id"`
	ProductName     string            `json:"product_name"`
	Quantity		int 			  `json:"quantity"`
	UnitCostPrice   float64           `json:"unit_cost_price"`
	Balance         int               `json:"quantity"`
	ManufactureDate *web.TimeResponse `json:"manufacture_date"`
	ExpiryDate      *web.TimeResponse `json:"expiry_date"`
	CreatedAt       web.TimeResponse  `json:"created_at"`
	UpdatedAt       web.TimeResponse  `json:"updated_at"`
	CreatedByID     string            `json:"created_by_id"`
	UpdatedByID     string            `json:"updated_by_id"`
}

// Response transforms Stock to StockResponse that is used for display.
func (s *Stock) Response(ctx context.Context) *StockResponse {
	if s == nil {
		return nil
	}

	r := &StockResponse{
		ID:            s.ID,
		BatchNumber:   s.BatchNumber,
		ProductID:     s.ProductID,
		ProductName:   s.ProductName,
		Quantity:	   s.Quantity,
		UnitCostPrice: s.UnitCostPrice,
		Balance:       s.Quantity - s.DeductedQuantity,
		CreatedAt:     web.NewTimeResponse(ctx, s.CreatedAt),
		UpdatedAt:     web.NewTimeResponse(ctx, s.UpdatedAt),
		CreatedByID:   s.CreatedByID,
		UpdatedByID:   s.UpdatedByID,
	}

	if s.ManufactureDate != nil {
		manufactureDate := web.NewTimeResponse(ctx, *s.ManufactureDate)
		r.ManufactureDate = &manufactureDate
	}

	if s.ExpiryDate != nil {
		expiryDate := web.NewTimeResponse(ctx, *s.ExpiryDate)
		r.ExpiryDate = &expiryDate
	}

	return r
}

// Stocks a list of Stocks.
type Stocks []*Stock

// Response transforms a list of Stocks to a list of StockResponse.
func (m *Stocks) Response(ctx context.Context) []*StockResponse {
	var l []*StockResponse
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l 
}

// StockCreateRequest contains information needed to create a new Stock.
type StockCreateRequest struct {
	BranchID         string     `json:"branch_id" validate:"required"`
	BatchNumber      string     `json:"batch_number" validate:"required"`
	ProductID        string     `json:"product_id" validate:"required"`
	UnitCostPrice    float64    `json:"unit_cost_price" validate:"required"`
	Quantity         int        `json:"quantity" validate:"required"`
	ManufactureDate  *time.Time `json:"manufacture_date,omitempty" schema:"omitempty"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty" schema:"omitempty"`
}

// StockReadRequest defines the information needed to read a Stock.
type StockReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// StockUpdateRequest defines what information may be provided to modify an existing
// Stock. All fields are optional so clients can send just the fields they want
// changed.
type StockUpdateRequest struct {
	ID              string     `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	BranchID        *string     `json:"branch_id"`
	ProductID 		*string    `json:"product_id"`
	BatchNumber     *string    `json:"batch_number"`
	UnitCostPrice   *float64   `json:"unit_cost_price"`
	Quantity        *int       `json:"quantity"`
	ManufactureDate *time.Time `json:"manufacture_date" schema:"omitempty"`
	ExpiryDate      *time.Time `json:"expiry_date" schema:"omitempty"`
}

// StockArchiveRequest defines the information needed to archive a Stock. This will archive (soft-delete) the
// existing database entry.
type StockArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// StockDeleteRequest defines the information needed to delete a Stock.
type StockDeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// ServiceTypeFindRequest defines the possible options to search for ServiceTypes. By default
// archived ServiceType will be excluded from response.
type StockFindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
	IncludeProducts bool          `json:"include_products" example:"false"`
}

// StockBalanceRequest defines the information needed to get Stock balance.
type StockBalanceRequest struct {
	ProductID string `json:"product_id" validate:"required"`
}

// StockReportRequest defines the possible options to generate stock summary report. By default
// archived Stock will be excluded from response.
type StockReportRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}

type StockInfo struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int64  `json:"quantity"`
}

func (repo Repository) FindStock(ctx context.Context, req StockFindRequest) (Stocks, error) {
	var queries []QueryMod
	
	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if !req.IncludeArchived {
		queries = append(queries, And("archived_at is null"))
	}

	if req.Where != "" {
		queries = append(queries, Where(req.Where, req.Args...))
	}

	if req.Limit != nil {
		queries = append(queries, Limit(int(*req.Limit)))
	}

	if req.Offset != nil {
		queries = append(queries, Offset(int(*req.Offset)))
	}

	if req.IncludeProducts {
		queries = append(queries, Load(models.StockRels.Product))
	}

	StockSlice, err := models.Stocks(queries...).All(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	var result Stocks
	for _, rec := range StockSlice {
		result = append(result, stockFromModel(rec))
	}

	return result, nil
}

// ReadStockByID gets the specified Stock by ID from the database.
func (repo *Repository) ReadStockByID(ctx context.Context, _ auth.Claims, id string) (*Stock, error) {
	stockModel, err := models.Stocks(Load(models.StockRels.Product), models.StockWhere.ID.EQ(id)).One(ctx, repo.DbConn)
	if err != nil {
		return  nil, err
	}

	return stockFromModel(stockModel), nil
}

func (repo *Repository) CreateStock(ctx context.Context, claims auth.Claims, req StockCreateRequest, now time.Time) (*Stock, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.CreateStock")
	defer span.Finish()
	if claims.Audience != "" {
		// Admin users can update Stocks they have access to.
		if !claims.HasRole(auth.RoleAdmin) {
			return nil, errors.WithStack(ErrForbidden)
		}
	}

	exists, err := models.Stocks(
		models.StockWhere.ProductID.EQ(req.ProductID),
		models.StockWhere.BatchNumber.EQ(req.BatchNumber)).Exists(ctx, repo.DbConn)
	if err != nil {
		return nil, err
	}

	ctx = webcontext.ContextAddUniqueValue(ctx, req, "ProductID", !exists)

	// Validate the request.
	v := webcontext.Validator()
	err = v.Struct(req)
	if err != nil {
		return nil, err
	}

	now = now.UTC()
	now.Truncate(time.Millisecond)

	s := Stock{
		ID:               uuid.NewRandom().String(),
		BranchID:		  req.BranchID,
		BatchNumber:      req.BatchNumber,
		ProductID:        req.ProductID,
		ProductName:      req.ProductID,
		UnitCostPrice:    req.UnitCostPrice,
		Quantity:         req.Quantity,
		CreatedAt:        now,
		UpdatedAt:        now,
		CreatedByID:      claims.Subject,
		UpdatedByID:      claims.Subject,
	}

	if req.ManufactureDate != nil && !req.ManufactureDate.IsZero() {
		manufactureDate := req.ManufactureDate.UTC()
		s.ManufactureDate = &manufactureDate
	}

	if req.ExpiryDate != nil && !req.ExpiryDate.IsZero() {
		expiryDate := req.ExpiryDate.UTC()
		s.ManufactureDate = &expiryDate
	}

	stockModel := s.model()
	if err := stockModel.Insert(ctx, repo.DbConn, boil.Infer()); err != nil {
		return nil, errors.WithMessage(err, "create Stock failed")
	}

	return &s, nil
}

// UpdateStock replaces an project in the database.
func (repo *Repository) UpdateStock(ctx context.Context, claims auth.Claims, req StockUpdateRequest, now time.Time) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.UpdateStock")
	defer span.Finish()

	if claims.Audience == "" {
		return errors.WithStack(ErrForbidden)
	}
	// Admin users can update Stocks they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}
	// Validate the request.
	v := webcontext.Validator()
	err := v.Struct(req)
	if err != nil {
		return err
	}

	now = now.UTC()
	now = now.Truncate(time.Millisecond)

	cols := models.M{
		models.StockColumns.UpdatedAt: now,
	}

	if req.Quantity != nil {
		cols[models.StockColumns.Quantity] = *req.Quantity
	}

	if req.UnitCostPrice != nil {
		cols[models.StockColumns.UnitCostPrice] = *req.UnitCostPrice
	}

	if req.ManufactureDate != nil {
		cols[models.StockColumns.ManufactureDate] = null.TimeFrom(*req.ManufactureDate)
	}

	if req.ExpiryDate != nil {
		cols[models.StockColumns.ExpiryDate] = null.TimeFrom(*req.ExpiryDate)
	}

	if len(cols) == 0 {
		return nil
	}

	_,err = models.Stocks(models.StockWhere.ID.EQ(req.ID)).UpdateAll(ctx, repo.DbConn, cols)

	return nil
}

// Delete removes an Stock from the database.
func (repo *Repository) DeleteStock(ctx context.Context, claims auth.Claims, req StockDeleteRequest) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "internal.shop.DeleteStock")
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
	// Admin users can update Stocks they have access to.
	if !claims.HasRole(auth.RoleAdmin) {
		return errors.WithStack(ErrForbidden)
	}

	_, err = models.Stocks(models.StockWhere.ID.EQ(req.ID)).DeleteAll(ctx, repo.DbConn)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// StockBalance returns the stock quantity of the specified product
func (repo *Repository) StockBalance(ctx context.Context, claims auth.Claims, req StockBalanceRequest) (quantity int64, err error) {
	var stockReportRequest = StockReportRequest{
		Where:           "product.id = ?",
		Args:            []interface{}{req.ProductID},
		Order:           nil,
		IncludeArchived: false,
	}
	report, err := repo.StockReport(ctx, claims, stockReportRequest)

	if err != nil || len(report) == 0 {
		return 0, err
	}

	return report[0].Quantity, nil
}

// StockReport returns a list of stock balance by product
func (repo *Repository) StockReport(ctx context.Context, _ auth.Claims, req StockReportRequest) ([]StockInfo, error) {
	if !req.IncludeArchived {
		if len(req.Where) > 0 {
			req.Where += " AND"
		}
		req.Where += fmt.Sprintf(" %s = null ", models.ProductColumns.ArchivedAt)
	}

	var report []StockInfo

	selectQuery := `SELECT product.id as product_id, product.name as product_name,
	(SUM(stock.deducted_quantity)*-1) + SUM(stock.quantity) AS quantity
	FROM stock INNER JOIN product ON stock.product_id = product.id
	GROUP BY product.id`
	// var query = []QueryMod {
	// 	SQL(selectQuery),
	// }
	// if len(req.Where) > 0 {
	// 	query = append(query, Where(req.Where, req.Args))
	// }

	// if req.Limit != nil {
	// 	query = append(query, Limit(int(*req.Limit)))
	// }

	// if req.Offset != nil {
	// 	query = append(query, Limit(int(*req.Offset)))
	// }

	err := models.NewQuery(SQL(selectQuery)).Bind(ctx, repo.DbConn, &report)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return report, errors.WithMessage(err, selectQuery)
}

