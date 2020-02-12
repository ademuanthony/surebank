package inventory

import (
	"context"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"merryworld/surebank/internal/branch"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/shop"
	"merryworld/surebank/internal/user"
)

// Repository defines the required dependencies for Inventory.
type Repository struct {
	DbConn *sqlx.DB
	mutex  sync.Mutex
}

// NewRepository creates a new Repository that defines dependencies for Inventory.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DbConn: db,
	}
}

// Inventory represents a financial transaction.
type Inventory struct {
	ID             string  `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	ProductID      string  `json:"product_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	BranchID       string  `json:"branch_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	TXType         string  `json:"tx_type"`
	OpeningBalance float64 `json:"opening_balance"`
	Quantity       float64 `json:"quantity"`
	Narration      string  `json:"narration"`
	SalesRepID     string  `json:"sales_rep_id"`
	CreatedAt      int64   `json:"created_at"`
	UpdatedAt      int64   `json:"updated_at"`
	ArchivedAt     *int64  `json:"archived_at"`

	Product  *shop.Product  `json:"product,omitempty"`
	Branch   *branch.Branch `json:"branch,omitempty"`
	SalesRep *user.User     `json:"sales_rep,omitempty"`
}

func FromModel(rec *models.Inventory) *Inventory {
	a := &Inventory{
		ID:             rec.ID,
		ProductID:      rec.ProductID,
		BranchID:       rec.BranchID,
		TXType:         rec.TXType,
		OpeningBalance: rec.OpeningBalance,
		Quantity:       rec.Quantity,
		Narration:      rec.Narration,
		SalesRepID:     rec.SalesRepID,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}

	if rec.R != nil {
			if rec.R.Product != nil {
			a.Product = shop.ProductFromModel(rec.R.Product)
		}

		if rec.R.Branch != nil {
			a.Branch = branch.FromModel(rec.R.Branch)
		}

		if rec.R.SalesRep != nil {
			a.SalesRep = user.FromModel(rec.R.SalesRep)
		}
	}

	if rec.ArchivedAt.Valid {
		a.ArchivedAt = &rec.ArchivedAt.Int64
	}

	return a
}

// Response represents a transaction that is returned for display.
type Response struct {
	ID             string            `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	ProductID      string            `json:"product_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	BranchID       string            `json:"branch_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	TXType         string            `json:"tx_type"`
	OpeningBalance int64             `json:"opening_balance"`
	Quantity       int64             `json:"quantity"`
	Narration      string            `json:"narration"`
	SalesRepID     string            `json:"sales_rep_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	CreatedAt      web.TimeResponse  `json:"created_at" truss:"api-read"`            // CreatedAt contains multiple format options for display.
	UpdatedAt      web.TimeResponse  `json:"updated_at" truss:"api-read"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt     *web.TimeResponse `json:"archived_at,omitempty" truss:"api-read"` // ArchivedAt contains multiple format options for display.

	Product  string `json:"product,omitempty" truss:"api-read"`
	Branch   string `json:"branch,omitempty" truss:"api-read"`
	SalesRep string `json:"sales_rep,omitempty" truss:"api-read"`
}

// Response transforms Inventory to the Response that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *Inventory) Response(ctx context.Context) *Response {
	if m == nil {
		return nil
	}

	r := &Response{
		ID:             m.ID,
		ProductID:      m.ProductID,
		BranchID:       m.BranchID,
		TXType:         m.TXType,
		OpeningBalance: int64(m.OpeningBalance),
		Quantity:       int64(m.Quantity),
		Narration:      m.Narration,
		SalesRepID:     m.SalesRepID,
		CreatedAt:      web.NewTimeResponse(ctx, time.Unix(m.CreatedAt, 0)),
		UpdatedAt:      web.NewTimeResponse(ctx, time.Unix(m.UpdatedAt, 0)),
	}

	if m.ArchivedAt != nil {
		at := web.NewTimeResponse(ctx, time.Unix(*m.ArchivedAt, 0))
		r.ArchivedAt = &at
	}

	if m.Branch != nil {
		r.Branch = m.Branch.Name
	}

	if m.Product != nil {
		r.Product = m.Product.Name
	}

	if m.SalesRep != nil {
		r.SalesRep = m.SalesRep.FullName()
	}

	return r
}

// Inventories a list of Inventories.
type Inventories []*Inventory

// Response transforms a list of Inventories to a list of Responses.
func (m *Inventories) Response(ctx context.Context) []*Response {
	var l = make([]*Response, 0)
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// PagedResponseList holds list of inventory and total count for pagination
type PagedResponseList struct {
	Transactions []*Response `json:"transactions"`
	TotalCount   int64       `json:"total_count"`
}

// AddStockRequest contains information needed to add a new Inventory of type, deposit.
type AddStockRequest struct {
	ProductID string  `json:"product_id" validate:"required"`
	Quantity  float64 `json:"amount" validate:"required,gt=0"`
}

// ReadRequest defines the information needed to read a inventory from the system.
type ReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// ArchiveRequest defines the information needed to archive a deposit. This will archive (soft-delete) the
// existing database entry.
type ArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// DeleteRequest defines the information needed to delete a customer account.
type DeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

type MakeStockDeductionRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int64  `json:"quantity"`
	Ref       string `json:"ref"`
}

// FindRequest defines the possible options to search for inventory. By default
// archived inventories will be excluded from response.
type FindRequest struct {
	Where           string        `json:"where" example:"type = deposit and branch_id = ? and created_at > ? and created_at < ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
	IncludeProduct  bool          `json:"include-product" example:"false"`
	IncludeBranch   bool          `json:"include-branch" example:"false"`
	IncludeSalesRep bool          `json:"include-sales-rep" example:"false"`
}

// ReportRequest defines the possible options to search for stock report
type ReportRequest struct {
	Where           string        `json:"where" example:"type = deposit and branch_id = ? and created_at > ? and created_at < ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
	IncludeProduct  bool          `json:"include-product" example:"false"`
	IncludeBranch   bool          `json:"include-branch" example:"false"`
	IncludeSalesRep bool          `json:"include-sales-rep" example:"false"`
}

type StockInfo struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int64  `json:"quantity"`
}
