package customer

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
)

// Repository defines the required dependencies for Customer.
type Repository struct {
	DbConn *sqlx.DB
}

// NewRepository creates a new Repository that defines dependencies for Customer.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DbConn: db,
	}
}

// Customer represents a workflow.
type Customer struct {
	ID          string     `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name        string     `json:"name"  validate:"required" example:"Rocket Launch"`
	Email       string     `json:"email" truss:"api-read"`
	PhoneNumber string     `json:"phone_number" truss:"api-read"`
	Address     string     `json:"address" truss:"api-read"`
	SalesRepID  string     `json:"sales_rep_id" truss:"api-read"`
	BranchID    string     `json:"branch_id" truss:"api-read"`
	SalesRep    string            `json:"sales_rep" truss:"api-read"`
	Branch      string            `json:"branch" truss:"api-read"`
	CreatedAt   time.Time  `json:"created_at" truss:"api-read"`
	UpdatedAt   time.Time  `json:"updated_at" truss:"api-read"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty" truss:"api-hide"`
}

func FromModel(rec *models.Customer) *Customer {
	c := &Customer{
		ID:          rec.ID,
		Name:        rec.Name,
		Email:       rec.Email,
		PhoneNumber: rec.PhoneNumber,
		Address:     rec.Address,
		BranchID:    rec.BranchID,
		SalesRepID:  rec.SalesRepID,
		CreatedAt:   rec.CreatedAt,
		UpdatedAt:   rec.UpdatedAt,
	}

	if rec.R != nil {
		if rec.R.Branch != nil {
			c.Branch = rec.R.Branch.Name
		}

		if rec.R.SalesRep != nil {
			c.SalesRep = rec.R.SalesRep.FirstName + " " + rec.R.SalesRep.LastName
		}
	}

	if rec.ArchivedAt.Valid {
		c.ArchivedAt = &rec.ArchivedAt.Time
	}

	return c
}

// Response represents a workflow that is returned for display.
type Response struct {
	ID          string            `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name        string            `json:"name"  validate:"required" example:"Rocket Launch"`
	Email       string            `json:"email" truss:"api-read"`
	PhoneNumber string            `json:"phone_number" truss:"api-read"`
	Address     string            `json:"address" truss:"api-read"`
	SalesRepID  string            `json:"sales_rep_id" truss:"api-read"`
	BranchID    string            `json:"branch_id" truss:"api-read"`
	SalesRep    string            `json:"sales_rep" truss:"api-read"`
	Branch      string            `json:"branch" truss:"api-read"`
	CreatedAt   web.TimeResponse  `json:"created_at"`            // CreatedAt contains multiple format options for display.
	UpdatedAt   web.TimeResponse  `json:"updated_at"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt  *web.TimeResponse `json:"archived_at,omitempty"` // ArchivedAt contains multiple format options for display.
}

// Response transforms Customer to the Response that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *Customer) Response(ctx context.Context) *Response {
	if m == nil {
		return nil
	}

	r := &Response{
		ID:          m.ID,
		Name:        m.Name,
		Email:       m.Email,
		PhoneNumber: m.PhoneNumber,
		Address:     m.Address,
		SalesRepID:  m.SalesRepID,
		SalesRep:    m.SalesRep,
		BranchID:    m.BranchID,
		Branch:      m.Branch,
		CreatedAt:   web.NewTimeResponse(ctx, m.CreatedAt),
		UpdatedAt:   web.NewTimeResponse(ctx, m.UpdatedAt),
	}

	if m.ArchivedAt != nil && !m.ArchivedAt.IsZero() {
		at := web.NewTimeResponse(ctx, *m.ArchivedAt)
		r.ArchivedAt = &at
	}

	return r
}

// Customers a list of Customers.
type Customers []*Customer

// Response transforms a list of Customers to a list of Responses.
func (m *Customers) Response(ctx context.Context) []*Response {
	var l []*Response
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// PagedResponseList holds a list of customers and total count
type PagedResponseList struct {
	Customers  []*Response `json:"customers"`
	TotalCount int64     `json:"total_count"`
}

// CreateRequest contains information needed to create a new Customer.
type CreateRequest struct {
	Name        string `json:"name" validate:"required"  example:"Oluwafe Dami"`
	Email       string `json:"email" validate:"required" example:"a@b.c"`
	PhoneNumber string `json:"phone_number" validate:"required" example:"0809000000"`
	Address     string `json:"address" example:"No 3 Ab Rd, Agege, Lagos"`
	SalesRepID  string `json:"sales_rep_id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	BranchID    string `json:"branch_id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`

	Type       string  `json:"type" validate:"required"`
	Target     float64 `json:"target"`
	TargetInfo string  `json:"target_info"`
}

// ReadRequest defines the information needed to read a checklist.
type ReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// UpdateRequest defines what information may be provided to modify an existing
// Customer. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank.
type UpdateRequest struct {
	ID     string           `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name        *string `json:"name" example:"Oluwafe Dami"`
	Email       *string `json:"email" example:"a@b.c"`
	PhoneNumber *string `json:"phone_number" example:"0809000000"`
	Address     *string `json:"address" example:"No 3 Ab Rd, Agege, Lagos"`
	SalesRepID  *string `json:"sales_rep_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	BranchID    *string `json:"branch_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// ArchiveRequest defines the information needed to archive a checklist. This will archive (soft-delete) the
// existing database entry.
type ArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// DeleteRequest defines the information needed to delete a branch.
type DeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// FindRequest defines the possible options to search for branches. By default
// archived checklist will be excluded from response.
type FindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
}
