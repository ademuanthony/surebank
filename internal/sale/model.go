package sale

import (
	"context"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null"

	"merryworld/surebank/internal/branch"
	"merryworld/surebank/internal/inventory"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/shop"
	"merryworld/surebank/internal/user"
)

// Repository defines the required dependencies for Branch.
type Repository struct {
	DbConn        *sqlx.DB
	ShopRepo      *shop.Repository
	InventoryRepo *inventory.Repository
	mutex         sync.Mutex
}

// NewRepository creates a new Repository that defines dependencies for Branch.
func NewRepository(db *sqlx.DB, shopRepo *shop.Repository, inventoryRepo *inventory.Repository) *Repository {
	return &Repository{
		DbConn:        db,
		ShopRepo:      shopRepo,
		InventoryRepo: inventoryRepo,
	}
}

// Sale
type Sale struct {
	ID            string     `json:"id"`
	ReceiptNumber string     `json:"receipt_number"`
	Amount        float64    `json:"amount"`
	AmountTender  float64    `json:"amount_tender"`
	Balance       float64    `json:"balance"`
	CustomerName  string     `json:"customer_name"`
	PhoneNumber   string     `json:"phone_number"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	ArchivedAt    *time.Time `json:"archived_at"`
	CreatedByID   string     `json:"created_by_id"`
	UpdatedByID   string     `json:"updated_by_id"`
	ArchivedByID  *string    `json:"archived_by_id"`
	BranchID      string     `json:"branch_id"`

	Items      []*Item        `json:"items"`
	Branch     *branch.Branch `json:"branch"`
	CreatedBy  *user.User     `json:"created_by"`
	UpdatedBy  *user.User     `json:"updated_by"`
	ArchivedBy *user.User     `json:"archived_by"`
}

func (s Sale) model() models.Sale {
	m := models.Sale{
		ID:            s.ID,
		ReceiptNumber: s.ReceiptNumber,
		Amount:        s.Amount,
		AmountTender:  s.AmountTender,
		Balance:       s.Balance,
		CustomerName:  null.StringFrom(s.CustomerName),
		PhoneNumber:   null.StringFrom(s.PhoneNumber),
		CreatedAt:     s.CreatedAt.Unix(),
		UpdatedAt:     s.UpdatedAt.Unix(),
		CreatedByID:   s.CreatedByID,
		UpdatedByID:   null.StringFrom(s.UpdatedByID),
		BranchID:      s.BranchID,
	}

	if s.ArchivedAt != nil {
		m.ArchivedAt = null.Int64From(s.ArchivedAt.Unix())
	}

	if s.ArchivedByID != nil {
		m.ArchivedByID = null.StringFrom(*s.ArchivedByID)
	}

	return m
}

// FromModel builds a Sale obj from a models.Sale obj
func FromModel(m *models.Sale) *Sale {
	s := &Sale{
		ID:            m.ID,
		ReceiptNumber: m.ReceiptNumber,
		Amount:        m.Amount,
		AmountTender:  m.AmountTender,
		Balance:       m.Balance,
		CustomerName:  m.CustomerName.String,
		PhoneNumber:   m.PhoneNumber.String,
		CreatedAt:     time.Unix(m.CreatedAt, 0),
		UpdatedAt:     time.Unix(m.UpdatedAt, 0),
		CreatedByID:   m.CreatedByID,
		UpdatedByID:   m.UpdatedByID.String,
		BranchID:      m.BranchID,
	}

	if m.ArchivedByID.Valid {
		s.ArchivedByID = &m.ArchivedByID.String
	}

	if m.ArchivedAt.Valid {
		archivedAt := time.Unix(m.ArchivedAt.Int64, 0)
		s.ArchivedAt = &archivedAt
	}

	if m.R != nil {
		if m.R.SaleItems != nil {
			for _, item := range m.R.SaleItems {
				s.Items = append(s.Items, ItemFromModel(item))
			}
		}

		if m.R.Branch != nil {
			s.Branch = branch.FromModel(m.R.Branch)
		}

		if m.R.CreatedBy != nil {
			s.CreatedBy = user.FromModel(m.R.CreatedBy)
		}

		if m.R.UpdatedBy != nil {
			s.UpdatedBy = user.FromModel(m.R.UpdatedBy)
		}

		if m.R.ArchivedBy != nil {
			s.ArchivedBy = user.FromModel(m.R.ArchivedBy)
		}
	}

	return s
}

// Response represent a sale obj for display
type Response struct {
	ID            string            `json:"id"`
	ReceiptNumber string            `json:"receipt_number"`
	Amount        float64           `json:"amount"`
	AmountTender  float64           `json:"amount_tender"`
	Balance       float64           `json:"balance"`
	CustomerName  string            `json:"customer_name"`
	PhoneNumber   string            `json:"phone_number"`
	CreatedAt     web.TimeResponse  `json:"created_at"`
	UpdatedAt     web.TimeResponse  `json:"updated_at"`
	ArchivedAt    *web.TimeResponse `json:"archived_at,omitempty"`
	CreatedByID   string            `json:"created_by_id"`
	UpdatedByID   string            `json:"updated_by_id"`
	ArchivedByID  *string           `json:"archived_by_id,omitempty"`
	BranchID      string            `json:"branch_id"`

	Items      []*ItemResponse `json:"items,omitempty"`
	Branch     *string         `json:"branch,omitempty"`
	CreatedBy  *string         `json:"created_by,omitempty"`
	UpdatedBy  *string         `json:"updated_by,omitempty"`
	ArchivedBy *string         `json:"archived_by,omitempty"`
}

// Response transforms Sale to Response that is used for display.
func (s *Sale) Response(ctx context.Context) *Response {
	r := &Response{
		ID:            s.ID,
		ReceiptNumber: s.ReceiptNumber,
		Amount:        s.AmountTender,
		AmountTender:  s.AmountTender,
		Balance:       s.Balance,
		CustomerName:  s.CustomerName,
		PhoneNumber:   s.PhoneNumber,
		CreatedAt:     web.NewTimeResponse(ctx, s.CreatedAt),
		UpdatedAt:     web.NewTimeResponse(ctx, s.UpdatedAt),
		CreatedByID:   s.CreatedByID,
		UpdatedByID:   s.UpdatedByID,
		BranchID:      s.BranchID,
	}

	if s.ArchivedAt != nil {
		t := web.NewTimeResponse(ctx, *s.ArchivedAt)
		r.ArchivedAt = &t
	}

	if s.ArchivedByID != nil {
		r.ArchivedByID = s.ArchivedByID
	}

	if s.Items != nil {
		for _, i := range s.Items {
			r.Items = append(r.Items, i.Response(ctx))
		}
	}

	if s.Branch != nil {
		r.Branch = &s.Branch.Name
	}

	if s.CreatedBy != nil {
		name := s.CreatedBy.FullName()
		r.CreatedBy = &name
	}

	if s.UpdatedBy != nil {
		name := s.UpdatedBy.FullName()
		r.UpdatedBy = &name
	}

	if s.ArchivedBy != nil {
		name := s.ArchivedBy.FullName()
		r.ArchivedBy = &name
	}

	return r
}

// Customers a list of Customers.
type Sales []*Sale

// Response transforms a list of Customers to a list of Responses.
func (m *Sales) Response(ctx context.Context) []*Response {
	var responses []*Response
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			responses = append(responses, n.Response(ctx))
		}
	}

	return responses
}

// Item represents sales Item
type Item struct {
	ID            string  `json:"id"`
	SaleID        string  `json:"sale_id"`
	ProductID     string  `json:"product_id"`
	Quantity      int     `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	UnitCostPrice float64 `json:"unit_cost_price"`
	StockIds      string  `json:"stock_ids"`

	Product *shop.Product `json:"product,omitempty"`
}

func (item *Item) model() *models.SaleItem {
	return &models.SaleItem{
		ID:            item.ID,
		SaleID:        item.SaleID,
		ProductID:     item.ProductID,
		UnitPrice:     item.UnitPrice,
		UnitCostPrice: item.UnitCostPrice,
		StockIds:      item.StockIds,
	}
}

func ItemFromModel(m *models.SaleItem) *Item {
	it := &Item{
		ID:            m.ID,
		SaleID:        m.SaleID,
		ProductID:     m.ProductID,
		Quantity:      m.Quantity,
		UnitPrice:     m.UnitPrice,
		UnitCostPrice: m.UnitCostPrice,
		StockIds:      m.StockIds,
	}

	if m.R != nil {
		if m.R.Product != nil {
			it.Product = shop.ProductFromModel(m.R.Product)
		}
	}

	return it
}

type ItemResponse struct {
	ID            string  `json:"id"`
	SaleID        string  `json:"sale_id"`
	ProductID     string  `json:"product_id"`
	Quantity      int     `json:"quantity"`
	UnitPrice     float64 `json:"unit_price"`
	UnitCostPrice float64 `json:"unit_cost_price"`
	StockIds      string  `json:"stock_ids"`

	Product *string `json:"product,omitempty"`
}

func (item Item) Response(ctx context.Context) *ItemResponse {
	r := &ItemResponse{
		ID:            item.ID,
		SaleID:        item.SaleID,
		ProductID:     item.ProductID,
		Quantity:      item.Quantity,
		UnitPrice:     item.UnitPrice,
		UnitCostPrice: item.UnitCostPrice,
		StockIds:      item.StockIds,
	}

	if item.Product != nil {
		r.Product = &item.Product.Name
	}

	return r
}

// PagedResponseList holds a list of sales and total count
type PagedResponseList struct {
	Sales      []*Response `json:"sales"`
	TotalCount int64       `json:"total_count"`
}

// MakeSalesRequest contains the payload for capturing a new sale
type MakeSalesRequest struct {
	AmountTender float64 `json:"amount_tender" validate:"required"`
	CustomerName string  `json:"customer_name"`
	PhoneNumber  string  `json:"phone_number"`

	Items []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

// ReadRequest defines the information needed to read a sale.
type ReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// ArchiveRequest defines the information needed to archive a sale. This will archive (soft-delete) the
// existing database entry.
type ArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// DeleteRequest defines the information needed to delete a sale.
type DeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// FindRequest defines the possible options to search for sales. By default
// archived checklist will be excluded from response.
type FindRequest struct {
	Where             string        `json:"where" example:"name = ? and status = ?"`
	Args              []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order             []string      `json:"order" example:"created_at desc"`
	Limit             *uint         `json:"limit" example:"10"`
	Offset            *uint         `json:"offset" example:"20"`
	IncludeArchived   bool          `json:"include-archived" example:"false"`
	IncludeItems      bool          `json:"include-items"`
	IncludeBranch     bool          `json:"include-branch"`
	IncludeCreatedBy  bool          `json:"include-created-by"`
	IncludeUpdatedBy  bool          `json:"include-updated-by"`
	IncludeArchivedBy bool          `json:"include-archived-by"`
}
