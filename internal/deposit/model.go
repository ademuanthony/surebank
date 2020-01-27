package deposit

import (
	"context"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"

	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/user"
)

// Repository defines the required dependencies for Deposit.
type Repository struct {
	DbConn    *sqlx.DB
	accNumMtx sync.Mutex
}

// NewRepository creates a new Repository that defines dependencies for Deposit.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DbConn: db,
	}
}

// Deposit represents a deposit.
type Deposit struct {
	ID         string     `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	AccountID  string     `json:"account_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Amount     float64    `json:"amount" truss:"api-read"`
	Narration  string     `json:"narration" truss:"api-read"`
	SalesRepID string     `json:"sales_rep_id" truss:"api-read"`
	CreatedAt  time.Time  `json:"created_at" truss:"api-read"`
	UpdatedAt  time.Time  `json:"updated_at" truss:"api-read"`
	ArchivedAt *time.Time `json:"archived_at,omitempty" truss:"api-hide"`

	SalesRep *user.User     `json:"sales_rep" truss:"api-read"`
	Account   *account.Account `json:"account" truss:"api-read"`
}

func FromModel(rec *models.Deposit) *Deposit {
	a := &Deposit{
		ID:         rec.ID,
		AccountID: rec.AccountID,
		Amount: rec.Amount,
		Narration: rec.Narration,
		SalesRepID: rec.SalesRepID,
		CreatedAt:  rec.CreatedAt,
		UpdatedAt:  rec.UpdatedAt,
	}

	if rec.R != nil {
		if rec.R.Account != nil {
			a.Account = account.FromModel(rec.R.Account)
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

// Response represents a deposit that is returned for display.
type Response struct {
	ID            string            `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	AccountID     string            `json:"account_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	AccountNumber string            `json:"account_number" example:"SB10003001" truss:"api-read"`
	CustomerID    string            `json:"customer_id" truss:"api-read"`
	Amount        float64           `json:"amount" truss:"api-read"`
	Narration     string            `json:"narration" truss:"api-read"`
	SalesRepID    string            `json:"sales_rep_id" truss:"api-read"`
	SalesRep      string            `json:"sales_rep,omitempty" truss:"api-read"`
	CreatedAt     web.TimeResponse  `json:"created_at" truss:"api-read"`            // CreatedAt contains multiple format options for display.
	UpdatedAt     web.TimeResponse  `json:"updated_at" truss:"api-read"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt    *web.TimeResponse `json:"archived_at,omitempty" truss:"api-read"` // ArchivedAt contains multiple format options for display.
}

// Response transforms Deposit to the Response that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *Deposit) Response(ctx context.Context) *Response {
	if m == nil {
		return nil
	}

	r := &Response{
		ID:            m.ID,
		AccountID:     m.AccountID,
		Amount:        m.Amount,
		Narration:     m.Narration,
		SalesRepID:    m.SalesRepID,
		CreatedAt:     web.NewTimeResponse(ctx, m.CreatedAt),
		UpdatedAt:     web.NewTimeResponse(ctx, m.UpdatedAt),
	}

	if m.ArchivedAt != nil && !m.ArchivedAt.IsZero() {
		at := web.NewTimeResponse(ctx, *m.ArchivedAt)
		r.ArchivedAt = &at
	}

	if m.Account != nil {
		r.AccountNumber = m.Account.Number
		r.CustomerID = m.Account.CustomerID
	}

	if m.SalesRep != nil {
		r.SalesRep = m.SalesRep.LastName + " " + m.SalesRep.FirstName
	}

	return r
}

// Deposits a list of Deposits.
type Deposits []*Deposit

// Response transforms a list of Deposits to a list of Responses.
func (m *Deposits) Response(ctx context.Context) []*Response {
	var l = make([]*Response, 0)
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// CreateRequest contains information needed to make a new Deposit.
type CreateRequest struct {
	AccountNumber string  `json:"account_number" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Narration     string  `json:"narration"`
}

// ReadRequest defines the information needed to read a deposit from the system.
type ReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// UpdateRequest defines what information may be provided to modify an existing
// Deposit. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank.
type UpdateRequest struct {
	ID        string  `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Amount    *float64 `json:"amount" validate:"omitempty,gt=0"`
	Narration *string  `json:"narration"`
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

// FindRequest defines the possible options to search for accounts. By default
// archived checklist will be excluded from response.
type FindRequest struct {
	Where           string        `json:"where" example:"account_id = ? and created_at > ? and created_at < ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeArchived bool          `json:"include-archived" example:"false"`
	IncludeCustomer bool          `json:"include_customer" example:"false"`
	IncludeBranch   bool          `json:"include_branch" example:"false"`
	IncludeSalesRep bool          `json:"include_sales_rep" example:"false"`
}
