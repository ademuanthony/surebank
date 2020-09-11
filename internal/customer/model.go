package customer

import (
	"context"
	"fmt"
	"strings"
	"time"

	"merryworld/surebank/internal/branch"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

// Repository defines the required dependencies for Customer.
type Repository struct {
	DbConn     *sqlx.DB
	mongoDb    *mongo.Database
	branchRepo *branch.Repository
}

var AccountTypes = []string{
	models.AccountTypeSB,
	models.AccountTypeDS,
}

// NewRepository creates a new Repository that defines dependencies for Customer.
func NewRepository(db *sqlx.DB, mongoDb *mongo.Database, branchRepo *branch.Repository) *Repository {
	return &Repository{
		DbConn:     db,
		mongoDb:    mongoDb,
		branchRepo: branchRepo,
	}
}

// Customer represents a workflow.
type Customer struct {
	ID          string     `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name        string     `json:"name"  validate:"required" example:"Rocket Launch"`
	ShortName   string     `json:"short_name"`
	Email       string     `json:"email" truss:"api-read"`
	PhoneNumber string     `json:"phone_number" truss:"api-read"`
	Address     string     `json:"address" truss:"api-read"`
	SalesRepID  string     `json:"sales_rep_id" truss:"api-read"`
	BranchID    string     `json:"branch_id" truss:"api-read"`
	SalesRep    string     `json:"sales_rep" truss:"api-read"`
	Branch      string     `json:"branch" truss:"api-read"`
	CreatedAt   time.Time  `json:"created_at" truss:"api-read"`
	UpdatedAt   time.Time  `json:"updated_at" truss:"api-read"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty" truss:"api-hide"`

	Accounts Accounts `json:"accounts"`
}

func FromModel(rec *models.Customer) *Customer {
	c := &Customer{
		ID:          rec.ID,
		Name:        rec.Name,
		ShortName:   rec.Name,
		Email:       rec.Email,
		PhoneNumber: rec.PhoneNumber,
		Address:     rec.Address,
		BranchID:    rec.BranchID,
		SalesRepID:  rec.SalesRepID,
		CreatedAt:   time.Unix(rec.CreatedAt, 0),
		UpdatedAt:   time.Unix(rec.UpdatedAt, 0),
	}

	if rec.R != nil {
		if rec.R.Branch != nil {
			c.Branch = rec.R.Branch.Name
		}

		if rec.R.SalesRep != nil {
			c.SalesRep = rec.R.SalesRep.FirstName + " " + rec.R.SalesRep.LastName
		}

		if accs := rec.R.Accounts; accs != nil {
			var numbers []string
			for _, a := range accs {
				numbers = append(numbers, a.Number)
			}

			c.Name = fmt.Sprintf("%s (%s)", c.Name, strings.Join(numbers, ", "))
		}
	}

	if rec.ArchivedAt.Valid {
		archivedAt := time.Unix(rec.ArchivedAt.Int64, 0)
		c.ArchivedAt = &archivedAt
	}

	return c
}

// Response represents a workflow that is returned for display.
type Response struct {
	ID          string            `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name        string            `json:"name"  validate:"required" example:"Rocket Launch"`
	ShortName   string            `json:"short_name"  validate:"required" example:"Rocket Launch"`
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
		ShortName:   m.ShortName,
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
type Customers []Customer

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
	TotalCount int64       `json:"total_count"`
}

// CreateRequest contains information needed to create a new Customer.
type CreateRequest struct {
	Name        string `json:"name" validate:"required"  example:"Oluwafe Dami"`
	Email       string `json:"email" example:"a@b.c"`
	PhoneNumber string `json:"phone_number" validate:"required" example:"0809000000"`
	Address     string `json:"address" example:"No 3 Ab Rd, Agege, Lagos"`
	SalesRepID  string `json:"sales_rep_id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	BranchID    string `json:"branch_id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`

	Type       string  `json:"type" validate:"required"`
	Target     float64 `json:"target"`
	TargetInfo string  `json:"target_info"`
}

// ReadRequest defines the information needed to read a customer.
type ReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// UpdateRequest defines what information may be provided to modify an existing
// Customer. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank.
type UpdateRequest struct {
	ID          string  `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Name        *string `json:"name" example:"Oluwafe Dami"`
	Email       *string `json:"email" example:"a@b.c"`
	PhoneNumber *string `json:"phone_number" example:"0809000000"`
	Address     *string `json:"address" example:"No 3 Ab Rd, Agege, Lagos"`
	SalesRepID  *string `json:"sales_rep_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	BranchID    *string `json:"branch_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// ArchiveRequest defines the information needed to archive a customer. This will archive (soft-delete) the
// existing database entry.
type ArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// DeleteRequest defines the information needed to delete a customer.
type DeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// FindRequest defines the possible options to search for customers. By default
// archived checklist will be excluded from response.
type FindRequest struct {
	Where            string        `json:"where" example:"name = ? and status = ?"`
	Args             []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order            []string      `json:"order" example:"created_at desc"`
	Limit            *uint         `json:"limit" example:"10"`
	Offset           *int64        `json:"offset" example:"20"`
	IncludeArchived  bool          `json:"include-archived" example:"false"`
	IncludeAccountNo bool          `json:"include-account-no" example:"false"`
}
