package account

import (
	"context"
	"merryworld/surebank/internal/branch"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/user"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
)

// Repository defines the required dependencies for Account.
type Repository struct {
	DbConn    *sqlx.DB
	accNumMtx sync.Mutex
}

// NewRepository creates a new Repository that defines dependencies for Account.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DbConn: db,
	}
}

// Account represents a customer account.
type Account struct {
	ID         string     `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	CustomerID string     `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Number     string     `json:"number"  validate:"required" example:"Rocket Launch"`
	Type       string     `json:"type" truss:"api-read"`
	Balance    float64    `json:"balance" truss:"api-read"`
	Target     float64    `json:"target" truss:"api-read"`
	TargetInfo string     `json:"target_info" truss:"api-read"`
	SalesRepID string     `json:"sales_rep_id" truss:"api-read"`
	BranchID   string     `json:"branch_id" truss:"api-read"`
	CreatedAt  time.Time  `json:"created_at" truss:"api-read"`
	UpdatedAt  time.Time  `json:"updated_at" truss:"api-read"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" truss:"api-hide"`

	Customer *customer.Customer `json:"customer"`
	SalesRep *user.User         `json:"sales_rep" truss:"api-read"`
	Branch   *branch.Branch     `json:"branch" truss:"api-read"`
}

func FromModel(rec *models.Account) *Account {
	a := &Account{
		ID:         rec.ID,
		CustomerID: rec.CustomerID,
		Number:     rec.Number,
		Type:       rec.AccountType,
		Balance:    rec.Balance,
		Target:     rec.Target,
		TargetInfo: rec.TargetInfo,
		SalesRepID: rec.SalesRepID,
		BranchID:   rec.BranchID,
		CreatedAt:  rec.CreatedAt,
		UpdatedAt:  rec.UpdatedAt,
	}

	if rec.R != nil {
		if rec.R.Branch != nil {
			a.Branch = branch.FromModel(rec.R.Branch)
		}

		if rec.R.Customer != nil {
			a.Customer = customer.FromModel(rec.R.Customer)
		}

		if rec.R.SalesRep != nil {
			a.SalesRep = user.FromModel(rec.R.SalesRep)
		}
	}

	if rec.ArchivedAt.Valid {
		a.ArchivedAt = &rec.ArchivedAt.Time
	}

	return a
}

// Response represents a customer account that is returned for display.
type Response struct {
	ID         string             `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	CustomerID string             `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Customer   *customer.Response `json:"customer,omitempty"`
	Number     string             `json:"number" example:"Rocket Launch"`
	Type       string             `json:"type" truss:"api-read"`
	Balance    float64            `json:"balance" truss:"api-read"`
	Target     float64            `json:"target" truss:"api-read"`
	TargetInfo string             `json:"target_info" truss:"api-read"`
	SalesRepID string             `json:"sales_rep_id" truss:"api-read"`
	BranchID   string             `json:"branch_id" truss:"api-read"`
	SalesRep   *user.UserResponse `json:"sales_rep,omitempty" truss:"api-read"`
	Branch     *branch.Response   `json:"branch,omitempty" truss:"api-read"`
	CreatedAt  web.TimeResponse   `json:"created_at"`            // CreatedAt contains multiple format options for display.
	UpdatedAt  web.TimeResponse   `json:"updated_at"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt *web.TimeResponse  `json:"archived_at,omitempty"` // ArchivedAt contains multiple format options for display.
}

// Response transforms Account to the Response that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *Account) Response(ctx context.Context) *Response {
	if m == nil {
		return nil
	}

	r := &Response{
		ID:         m.ID,
		CustomerID: m.CustomerID,
		Customer:   m.Customer.Response(ctx),
		Number:     m.Number,
		Type:       m.Type,
		Balance:    m.Balance,
		Target:     m.Target,
		TargetInfo: m.TargetInfo,
		SalesRepID: m.SalesRepID,
		BranchID:   m.BranchID,
		SalesRep:   m.SalesRep.Response(ctx),
		Branch:     m.Branch.Response(ctx),
		CreatedAt:  web.NewTimeResponse(ctx, m.CreatedAt),
		UpdatedAt:  web.NewTimeResponse(ctx, m.UpdatedAt),
	}

	if m.ArchivedAt != nil && !m.ArchivedAt.IsZero() {
		at := web.NewTimeResponse(ctx, *m.ArchivedAt)
		r.ArchivedAt = &at
	}

	return r
}

// Accounts a list of Accounts.
type Accounts []*Account

// Response transforms a list of Accounts to a list of Responses.
func (m *Accounts) Response(ctx context.Context) []*Response {
	var l = make([]*Response, 0)
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// CreateRequest contains information needed to create a new Account.
type CreateRequest struct {
	CustomerID string  `json:"customer_id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Type       string  `json:"type" validate:"required"`
	Target     float64 `json:"target"`
	TargetInfo string  `json:"target_info"`
	BranchID   string  `json:"branch_id"`
}

// ReadRequest defines the information needed to read a customer account.
type ReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// UpdateRequest defines what information may be provided to modify an existing
// Account. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank.
type UpdateRequest struct {
	ID         string   `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Type       *string  `json:"type" validate:"required"`
	Target     *float64 `json:"target"`
	TargetInfo *string  `json:"target_info"`
}

// ArchiveRequest defines the information needed to archive a customer account. This will archive (soft-delete) the
// existing database entry.
type ArchiveRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// DeleteRequest defines the information needed to delete a customer account.
type DeleteRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// FindRequest defines the possible options to search for accounts. By default
// archived checklist will be excluded from response.
type FindRequest struct {
	Where           string        `json:"where" example:"name = ? and status = ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
	IncludeCustomer bool          `json:"include_customer" example:"false"`
	IncludeBranch   bool          `json:"include_branch" example:"false"`
	IncludeSalesRep bool          `json:"include_sales_rep" example:"false"`
}
