package dscommission

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"

	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
)

// Repository defines the required dependencies for Transaction.
type Repository struct {
	DbConn  *sqlx.DB
	mongoDb *mongo.Database
}

// NewRepository creates a new Repository that defines dependencies for Transaction.
func NewRepository(db *sqlx.DB, mongoDb *mongo.Database) *Repository {
	return &Repository{
		DbConn:  db,
		mongoDb: mongoDb,
	}
}

type DsCommission struct {
	ID            string  `boil:"id" json:"id" toml:"id" yaml:"id"`
	AccountID     string  `boil:"account_id" json:"account_id" toml:"account_id" yaml:"account_id"`
	AccountNumber string  `json:"account_number"`
	CustomerID    string  `boil:"customer_id" json:"customer_id" toml:"customer_id" yaml:"customer_id"`
	Customer      string  `json:"customer"`
	Amount        float64 `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Date          int64   `boil:"date" json:"date" toml:"date" yaml:"date"`
	EffectiveDate int64   `boil:"effective_date" json:"effective_date" toml:"effective_date" yaml:"effective_date"`
}

func FromModel(rec *models.DSCommission) *DsCommission {
	a := &DsCommission{
		ID:            rec.ID,
		AccountID:     rec.AccountID,
		Amount:        rec.Amount,
		CustomerID:    rec.CustomerID,
		Date:          rec.Date,
		EffectiveDate: rec.EffectiveDate,
	}

	if rec.R != nil {
		if rec.R.Account != nil {
			a.AccountNumber = rec.R.Account.Number
		}

		if rec.R.Customer != nil {
			a.Customer = rec.R.Customer.Name
		}
	}
	return a
}

// Response represents a transaction that is returned for display.
type Response struct {
	ID            string           `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	AccountID     string           `json:"account_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	AccountNumber string           `json:"account_number" example:"SB10003001" truss:"api-read"`
	CustomerID    string           `json:"customer_id" truss:"api-read"`
	CustomerName  string           `json:customer_name" truss:"api-read"`
	Amount        float64          `json:"amount" truss:"api-read"`
	Date          web.TimeResponse `json:"date" truss:"api-read"`
	EffectiveDate web.TimeResponse `json:"effective_date" truss:"api-read"`
}

// Response transforms DsCommission to the Response that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *DsCommission) Response(ctx context.Context) Response {
	if m == nil {
		return Response{}
	}

	r := Response{
		ID:            m.ID,
		AccountID:     m.AccountID,
		AccountNumber: m.AccountNumber,
		CustomerID:    m.CustomerID,
		CustomerName:  m.Customer,
		Amount:        m.Amount,
		Date:          web.NewTimeResponse(ctx, time.Unix(m.Date, 0)),
		EffectiveDate: web.NewTimeResponse(ctx, time.Unix(m.EffectiveDate, 0)),
	}

	return r
}

// DsCommissions a list of DsCommissions.
type DsCommissions []DsCommission

// Response transforms a list of DsCommissions to a list of Responses.
func (m DsCommissions) Response(ctx context.Context) []Response {
	var l = make([]Response, 0)
	if m != nil && len(m) > 0 {
		for _, n := range m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// PagedResponseList holds list of transaction and total count for pagination
type PagedResponseList struct {
	Items      []Response `json:"items"`
	TotalCount int64      `json:"total_count"`
}

// FindRequest defines the possible options to search for commission. By default
// archived checklist will be excluded from response.
type FindRequest struct {
	StateDate       int64    `json:"start_date"`
	EndDate         int64    `json:"end_date"`
	Order           []string `json:"order" example:"created_at desc"`
	Limit           *uint    `json:"limit" example:"10"`
	Offset          *int64   `json:"offset" example:"20"`
	IncludeArchived bool     `json:"include-archived" example:"false"`
	IncludeAccount  bool     `json:"include-account" example:"false"`
	IncludeCustomer bool     `json:"include-customer" example:"false"`
}
