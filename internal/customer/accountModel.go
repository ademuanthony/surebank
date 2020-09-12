package customer

import (
	"context"
	"sync"
	"time"

	"merryworld/surebank/internal/branch"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

// AccountRepository defines the required dependencies for Account.
type AccountRepository struct {
	DbConn       *sqlx.DB
	accNumMtx    sync.Mutex
	mongoDb      *mongo.Database
	customerRepo *Repository
	branchRepo   *branch.Repository
}

// NewAccountRepository creates a new Repository that defines dependencies for Account.
func NewAccountRepository(db *sqlx.DB, mongoDb *mongo.Database, customerRepo *Repository,
	branchRepo *branch.Repository) *AccountRepository {
	return &AccountRepository{
		DbConn:       db,
		mongoDb:      mongoDb,
		customerRepo: customerRepo,
		branchRepo:   branchRepo,
	}
}

// Account represents a customer account.
type Account struct {
	ID              string  `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	CustomerID      string  `json:"customer_id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Number          string  `json:"number"  validate:"required" example:"Rocket Launch"`
	Type            string  `json:"type" truss:"api-read"`
	Balance         float64 `json:"balance" truss:"api-read"`
	Target          float64 `json:"target" truss:"api-read"`
	TargetInfo      string  `json:"target_info" truss:"api-read"`
	SalesRepID      string  `json:"sales_rep_id" truss:"api-read"`
	BranchID        string  `json:"branch_id" truss:"api-read"`
	LastPaymentDate int64   `json:"last_payment_date"`
	CreatedAt       int64   `json:"created_at" truss:"api-read"`
	UpdatedAt       int64   `json:"updated_at" truss:"api-read"`
	ArchivedAt      int64   `json:"archived_at,omitempty" truss:"api-hide"`
	Customer        string  `json:"customer"`
	PhoneNumber     string  `json:"phone_number"`
	SalesRep        string  `json:"sales_rep" truss:"api-read"`
	Branch          string  `json:"branch" truss:"api-read"`
}

func AccountFromModel(rec *models.Account) *Account {
	a := &Account{
		ID:              rec.ID,
		CustomerID:      rec.CustomerID,
		Number:          rec.Number,
		Type:            rec.AccountType,
		Balance:         rec.Balance,
		Target:          rec.Target,
		TargetInfo:      rec.TargetInfo,
		SalesRepID:      rec.SalesRepID,
		BranchID:        rec.BranchID,
		LastPaymentDate: rec.LastPaymentDate,
		CreatedAt:       rec.CreatedAt,
		UpdatedAt:       rec.UpdatedAt,
	}

	if rec.R != nil {
		if rec.R.Branch != nil {
			a.Branch = rec.R.Branch.Name
		}

		if rec.R.Customer != nil {
			a.Customer = rec.R.Customer.Name
		}

		if rec.R.SalesRep != nil {
			a.SalesRep = rec.R.SalesRep.FirstName + " " + rec.R.SalesRep.LastName
		}
	}

	if rec.ArchivedAt.Valid {
		a.ArchivedAt = rec.ArchivedAt.Int64
	}

	return a
}

// AccountResponse represents a customer account that is returned for display.
type AccountResponse struct {
	ID              string            `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	CustomerID      string            `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	Customer        string            `json:"customer,omitempty" truss:"api-read"`
	PhoneNumber     string            `json:"phone_number"`
	Number          string            `json:"number" example:"Rocket Launch" truss:"api-read"`
	Type            string            `json:"type" truss:"api-read"`
	Balance         float64           `json:"balance" truss:"api-read"`
	Target          float64           `json:"target" truss:"api-read"`
	TargetInfo      string            `json:"target_info" truss:"api-read"`
	SalesRepID      string            `json:"sales_rep_id" truss:"api-read"`
	BranchID        string            `json:"branch_id" truss:"api-read"`
	SalesRep        string            `json:"sales_rep,omitempty" truss:"api-read"`
	Branch          string            `json:"branch,omitempty" truss:"api-read"`
	LastPaymentDate web.TimeResponse  `json:"created_at"`            // CreatedAt contains multiple format options for display.
	CreatedAt       web.TimeResponse  `json:"created_at"`            // CreatedAt contains multiple format options for display.
	UpdatedAt       web.TimeResponse  `json:"updated_at"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt      *web.TimeResponse `json:"archived_at,omitempty"` // ArchivedAt contains multiple format options for display.
}

// Response transforms Account to the Response that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *Account) Response(ctx context.Context) *AccountResponse {
	if m == nil {
		return nil
	}

	r := &AccountResponse{
		ID:              m.ID,
		CustomerID:      m.CustomerID,
		Customer:        m.Customer,
		PhoneNumber:     m.PhoneNumber,
		Number:          m.Number,
		Type:            m.Type,
		Balance:         m.Balance,
		Target:          m.Target,
		TargetInfo:      m.TargetInfo,
		SalesRepID:      m.SalesRepID,
		BranchID:        m.BranchID,
		LastPaymentDate: web.NewTimeResponse(ctx, time.Unix(m.LastPaymentDate, 0)),
		CreatedAt:       web.NewTimeResponse(ctx, time.Unix(m.CreatedAt, 0)),
		UpdatedAt:       web.NewTimeResponse(ctx, time.Unix(m.UpdatedAt, 0)),
	}

	if m.ArchivedAt > 0 {
		at := web.NewTimeResponse(ctx, time.Unix(m.ArchivedAt, 0))
		r.ArchivedAt = &at
	}

	if m.SalesRep != "" {
		r.SalesRep = m.SalesRep
	}

	if m.Branch != "" {
		r.Branch = m.Branch
	}

	return r
}

// Accounts a list of Accounts.
type Accounts []Account

// Response transforms a list of Accounts to a list of Responses.
func (m *Accounts) Response(ctx context.Context) []*AccountResponse {
	var l = make([]*AccountResponse, 0)
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// AccountPagedResponseList hold a list of accounts and total count
type AccountPagedResponseList struct {
	Accounts   []*AccountResponse `json:"accounts"`
	TotalCount int64              `json:"total_count"`
}

// CreateAccountRequest contains information needed to create a new Account.
type CreateAccountRequest struct {
	CustomerID string  `json:"customer_id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Type       string  `json:"type" validate:"required"`
	Target     float64 `json:"target"`
	TargetInfo string  `json:"target_info"`
	BranchID   string  `json:"branch_id"`
}

// ReadAccountRequest defines the information needed to read a customer account.
type ReadAccountRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// UpdateAccountRequest defines what information may be provided to modify an existing
// Account. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank.
type UpdateAccountRequest struct {
	ID         string   `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Type       *string  `json:"type" validate:"required"`
	Target     *float64 `json:"target"`
	TargetInfo *string  `json:"target_info"`
}

// ArchiveAccountRequest defines the information needed to archive a customer account. This will archive (soft-delete) the
// existing database entry.
type ArchiveAccountRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// DeleteAccountRequest defines the information needed to delete a customer account.
type DeleteAccountRequest struct {
	ID string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
}

// FindAccountRequest defines the possible options to search for accounts. By default
// archived checklist will be excluded from response.
type FindAccountRequest struct {
	CustomerID      string   `json:"where" example:"name = ? and status = ?"`
	Order           []string `json:"order" example:"created_at desc"`
	Limit           *uint    `json:"limit" example:"10"`
	Offset          *int64   `json:"offset" example:"20"`
	IncludeArchived bool     `json:"include-archived" example:"false"`
	IncludeCustomer bool     `json:"include_customer" example:"false"`
	IncludeBranch   bool     `json:"include_branch" example:"false"`
	IncludeSalesRep bool     `json:"include_sales_rep" example:"false"`
}
