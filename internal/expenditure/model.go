package expenditure

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/user"
)

// Repository defines the required dependencies for Transaction.
type Repository struct {
	DbConn *sqlx.DB
}

// NewRepository creates a new Repository that defines dependencies for Transaction.
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		DbConn: db,
	}
}

// Expenditure represents a financial expense by a rep.
type Expenditure struct {
	ID         string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	SalesRepID string    `boil:"sales_rep_id" json:"sales_rep_id" toml:"sales_rep_id" yaml:"sales_rep_id"`
	Amount     float64   `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Reason     string    `boil:"reason" json:"reason" toml:"reason" yaml:"reason"`
	Date       time.Time `boil:"date" json:"date" toml:"date" yaml:"date"`

	SalesRep *user.User
}

func FromModel(rec *models.RepsExpense) *Expenditure {
	a := &Expenditure{
		ID:         rec.ID,
		SalesRepID: rec.SalesRepID,
		Amount:     rec.Amount,
		Reason:     rec.Reason,
		Date:       time.Unix(rec.Date, 0).UTC(),
	}

	if rec.R != nil {
		if rec.R.SalesRep != nil {
			a.SalesRep = user.FromModel(rec.R.SalesRep)
		}
	}
	return a
}

// Response represents a expenditure that is returned for display.
type Response struct {
	ID         string           `boil:"id" json:"id" toml:"id" yaml:"id"`
	SalesRepID string           `boil:"sales_rep_id" json:"sales_rep_id" toml:"sales_rep_id" yaml:"sales_rep_id"`
	Amount     float64          `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Reason     string           `boil:"reason" json:"reason" toml:"reason" yaml:"reason"`
	Date       web.TimeResponse `json:"effective_date" truss:"api-read"`

	SalesRep string `json:"sales_rep"`
}

// Response transforms Transaction to the Response that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *Expenditure) Response(ctx context.Context) *Response {
	if m == nil {
		return nil
	}

	r := &Response{
		ID:         m.ID,
		SalesRepID: m.SalesRepID,
		Amount:     m.Amount,
		Reason:     m.Reason,
		Date:       web.NewTimeResponse(ctx, m.Date),
	}

	if m.SalesRep != nil {
		r.SalesRep = m.SalesRep.LastName + " " + m.SalesRep.FirstName
	}

	return r
}

// Expenditures a list of Expenditures.
type Expenditures []*Expenditure

// Response transforms a list of Expenditures to a list of Responses.
func (m *Expenditures) Response(ctx context.Context) []*Response {
	var l = make([]*Response, 0)
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// PagedResponseList holds list of expenditures and total count for pagination
type PagedResponseList struct {
	Expenditures []*Response `json:"expenditures"`
	TotalCount   int64       `json:"total_count"`
}

// CreateRequest contains information needed to make a new Transaction.
type CreateRequest struct {
	SalesRepPhoneNumber string  `validate:"required" json:"sales_rep_phone_number" toml:"sales_rep_phone_number" yaml:"sales_rep_phone_number"`
	Amount              float64 `validate:"required" json:"amount" toml:"amount" yaml:"amount"`
	Reason              string  `validate:"required" json:"reason" toml:"reason" yaml:"reason"`
}

// ReadRequest defines the information needed to read a deposit from the system.
type ReadRequest struct {
	ID              string `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	IncludeArchived bool   `json:"include-archived" example:"false"`
}

// UpdateRequest defines what information may be provided to modify an existing
// Transaction. All fields are optional so clients can send just the fields they want
// changed. It uses pointer fields so we can differentiate between a field that
// was not provided and a field that was provided as explicitly blank.
type UpdateRequest struct {
	ID     string   `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Amount *float64 `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Reason *string  `boil:"reason" json:"reason" toml:"reason" yaml:"reason"`
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
	Where           string        `json:"where" example:"type = deposit and account_id = ? and created_at > ? and created_at < ?"`
	Args            []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order           []string      `json:"order" example:"created_at desc"`
	Limit           *uint         `json:"limit" example:"10"`
	Offset          *uint         `json:"offset" example:"20"`
	IncludeSalesRep bool          `json:"include_sales_rep" example:"false"`
}
