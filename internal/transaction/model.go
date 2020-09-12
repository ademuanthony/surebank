package transaction

import (
	"context"
	"database/sql/driver"
	"errors"
	"merryworld/surebank/internal/customer"
	"merryworld/surebank/internal/dscommission"
	"merryworld/surebank/internal/platform/notify"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/go-playground/validator.v9"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null"

	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
)

// Repository defines the required dependencies for Transaction.
type Repository struct {
	DbConn         *sqlx.DB
	CommissionRepo *dscommission.Repository
	AccountRepo    *customer.AccountRepository
	mongoDb        *mongo.Database
	notifySMS      notify.SMS
	accNumMtx      sync.Mutex
}

// NewRepository creates a new Repository that defines dependencies for Transaction.
func NewRepository(db *sqlx.DB, mongoDb *mongo.Database, commissionRepo *dscommission.Repository,
	accountRepo *customer.AccountRepository, notifySMS notify.SMS) *Repository {

	return &Repository{
		DbConn:         db,
		AccountRepo:    accountRepo,
		mongoDb:        mongoDb,
		CommissionRepo: commissionRepo,
		notifySMS:      notifySMS,
	}
}

// Transaction represents a financial transaction.
type Transaction struct {
	ID             string          `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	AccountID      string          `json:"account_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	AccountNumber  string          `json:"account_number"`
	CustomerID     string          `json:"customer_id"`
	Customer       string          `json:"customer"`
	Type           TransactionType `json:"type" example:"deposit"`
	OpeningBalance float64         `json:"opening_balance" exmaple:"34500.01"`
	Amount         float64         `json:"amount" truss:"api-read"`
	Narration      string          `json:"narration" truss:"api-read"`
	PaymentMethod  string          `json:"payment_method" truss:"api-read"`
	SalesRepID     string          `json:"sales_rep_id" truss:"api-read"`
	SalesRep       string          `json:"sales_rep" truss:"api-read"`
	ReceiptNo      string          `json:"receipt_no"`
	EffectiveDate  int64           `json:"effective_date"`
	CreatedAt      int64           `json:"created_at" truss:"api-read"`
	UpdatedAt      int64           `json:"updated_at" truss:"api-read"`
	ArchivedAt     *time.Time      `json:"archived_at,omitempty" truss:"api-hide"`
}

func FromModel(rec *models.Transaction) *Transaction {
	a := &Transaction{
		ID:             rec.ID,
		AccountID:      rec.AccountID,
		Type:           TransactionType(rec.TXType),
		OpeningBalance: rec.OpeningBalance,
		Amount:         rec.Amount,
		Narration:      rec.Narration,
		PaymentMethod:  rec.PaymentMethod,
		SalesRepID:     rec.SalesRepID,
		ReceiptNo:      rec.ReceiptNo,
		EffectiveDate:  rec.EffectiveDate,
		CreatedAt:      rec.CreatedAt,
		UpdatedAt:      rec.UpdatedAt,
	}

	if rec.R != nil {
		if rec.R.Account != nil {
			a.AccountNumber = rec.R.Account.Number
		}

		if rec.R.SalesRep != nil {
			a.SalesRep = rec.R.SalesRep.FirstName + " " + rec.R.SalesRep.LastName
		}
	}

	if rec.ArchivedAt.Valid {
		archivedAt := time.Unix(rec.ArchivedAt.Int64, 0)
		a.ArchivedAt = &archivedAt
	}

	return a
}

// Response represents a transaction that is returned for display.
type Response struct {
	ID             string            `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	AccountID      string            `json:"account_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	Type           TransactionType   `json:"type,omitempty" example:"deposit"`
	AccountNumber  string            `json:"account_number" example:"SB10003001" truss:"api-read"`
	CustomerID     string            `json:"customer_id" truss:"api-read"`
	Customer       string            `json:"customer"`
	OpeningBalance float64           `json:"opening_balance" truss:"api-read"`
	Amount         float64           `json:"amount" truss:"api-read"`
	Narration      string            `json:"narration" truss:"api-read"`
	PaymentMethod  string            `json:"payment_method" truss:"api-read"`
	SalesRepID     string            `json:"sales_rep_id" truss:"api-read"`
	SalesRep       string            `json:"sales_rep,omitempty" truss:"api-read"`
	ReceiptNo      string            `json:"receipt_no"`
	EffectiveDate  web.TimeResponse  `json:"effective_date" truss:"api-read"`
	CreatedAt      web.TimeResponse  `json:"created_at" truss:"api-read"`            // CreatedAt contains multiple format options for display.
	UpdatedAt      web.TimeResponse  `json:"updated_at" truss:"api-read"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt     *web.TimeResponse `json:"archived_at,omitempty" truss:"api-read"` // ArchivedAt contains multiple format options for display.
}

// Response transforms Transaction to the Response that is used for display.
// Additional filtering by context values or translations could be applied.
func (m *Transaction) Response(ctx context.Context) Response {
	if m == nil {
		return Response{}
	}

	r := Response{
		ID:             m.ID,
		AccountID:      m.AccountID,
		AccountNumber:  m.AccountNumber,
		CustomerID:     m.CustomerID,
		Customer:       m.Customer,
		Type:           m.Type,
		OpeningBalance: m.OpeningBalance,
		Amount:         m.Amount,
		Narration:      m.Narration,
		PaymentMethod:  m.PaymentMethod,
		ReceiptNo:      m.ReceiptNo,
		SalesRepID:     m.SalesRepID,
		SalesRep:       m.SalesRep,
		EffectiveDate:  web.NewTimeResponse(ctx, time.Unix(m.EffectiveDate, 0)),
		CreatedAt:      web.NewTimeResponse(ctx, time.Unix(m.CreatedAt, 0)),
		UpdatedAt:      web.NewTimeResponse(ctx, time.Unix(m.UpdatedAt, 0)),
	}

	if m.ArchivedAt != nil && !m.ArchivedAt.IsZero() {
		at := web.NewTimeResponse(ctx, *m.ArchivedAt)
		r.ArchivedAt = &at
	}

	return r
}

type TxReportResponse struct {
	ID            string          `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	AccountID     string          `json:"account_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86" truss:"api-read"`
	Type          TransactionType `json:"tx_type,omitempty" example:"deposit"`
	AccountNumber string          `json:"account_number" example:"SB10003001" truss:"api-read"`
	CustomerID    string          `json:"customer_id" truss:"api-read"`
	CustomerName  string          `json:"customer_name" truss:"api-read"`
	Amount        float64         `json:"amount" truss:"api-read"`
	Narration     string          `json:"narration" truss:"api-read"`
	PaymentMethod string          `json:"payment_method" truss:"api-read"`
	SalesRepID    string          `json:"sales_rep_id" truss:"api-read"`
	SalesRep      string          `json:"sales_rep,omitempty" truss:"api-read"`
	ReceiptNo     string          `json:"receipt_no"`
	EffectiveDate int64           `json:"effective_date" truss:"api-read"`
	CreatedAt     int64           `json:"created_at" truss:"api-read"`            // CreatedAt contains multiple format options for display.
	UpdatedAt     int64           `json:"updated_at" truss:"api-read"`            // UpdatedAt contains multiple format options for display.
	ArchivedAt    null.Int64      `json:"archived_at,omitempty" truss:"api-read"` // ArchivedAt contains multiple format options for display.
}

// Transactions a list of Transactions.
type Transactions []Transaction

// Response transforms a list of Transactions to a list of Responses.
func (m *Transactions) Response(ctx context.Context) []Response {
	var l = make([]Response, 0)
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// PagedResponseList holds list of transaction and total count for pagination
type PagedResponseList struct {
	Transactions []Response `json:"transactions"`
	TotalCount   int64      `json:"total_count"`
}

// CreateRequest contains information needed to make a new Transaction.
type CreateRequest struct {
	Type          TransactionType `json:"type" validate:"required,oneof=deposit withdrawal"`
	AccountNumber string          `json:"account_number" validate:"required"`
	Amount        float64         `json:"amount" validate:"required,gt=0"`
	Narration     string          `json:"narration"`
	PaymentMethod string          `json:"payment_method"`
}

// WithdrawRequest contains information needed to make a new Transaction.
type WithdrawRequest struct {
	Type              TransactionType `json:"type" validate:"required,oneof=deposit withdrawal"`
	AccountNumber     string          `json:"account_number" validate:"required"`
	Amount            float64         `json:"amount" validate:"required,gt=0"`
	PaymentMethod     string          `json:"payment_method" validate:"required"`
	Bank              string          `json:"bank"`
	BankAccountNumber string          `json:"bank_account_number"`
	Narration         string          `json:"narration"`
}

type MakeDeductionRequest struct {
	AccountNumber string  `json:"account_number" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	Narration     string  `json:"narration"`
}

// CreateDepositRequest contains information needed to add a new Transaction of type, deposit.
type CreateDepositRequest struct {
	AccountNumber string  `json:"account_number" validate:"required"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
	PaymentMethod string  `json:"payment_method" validate:"required" truss:"api-read"`
	Narration     string  `json:"narration"`
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
	ID        string   `json:"id" validate:"required,uuid" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
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
	CustomerID      string          `json:"customer_id"`
	AccountID       string          `json:"account_id"`
	AccountNumber   string          `json:"account_number"`
	SalesRepID      string          `json:"sales_rep_id"`
	PaymentMethod   string          `json:"payment_method"`
	StartDate       int64           `json:"start_date"`
	EndDate         int64           `json:"end_date"`
	Type            TransactionType `json:"type"`
	Order           []string        `json:"order" example:"created_at desc"`
	Limit           *uint           `json:"limit" example:"10"`
	Offset          *int64          `json:"offset" example:"20"`
	IncludeArchived bool            `json:"include-archived" example:"false"`
}

// ChecklistStatus represents the status of checklist.
type TransactionType string

// ChecklistStatus values define the status field of checklist.
const (
	// TransactionType_Deposit defines the type of deposit transaction.
	TransactionType_Deposit TransactionType = "deposit"
	// TransactionType_Withdrawal defines the type of withdrawal transaction.
	TransactionType_Withdrawal TransactionType = "withdrawal"

	PaymentMethod_Cash string = "cash"
	PaymentMethod_Bank string = "bank_deposit"
)

// TransactionType_Values provides list of valid TransactionType values.
var TransactionType_Values = []TransactionType{
	TransactionType_Deposit,
	TransactionType_Withdrawal,
}

// TransactionType_ValuesInterface returns the TransactionType options as a slice interface.
func TransactionType_ValuesInterface() []interface{} {
	var l = make([]interface{}, len(TransactionType_Values))
	for i, v := range TransactionType_Values {
		l[i] = v.String()
	}
	return l
}

// Scan supports reading the TransactionType value from the database.
func (s *TransactionType) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Scan source is not []byte")
	}

	*s = TransactionType(string(asBytes))
	return nil
}

// Value converts the TransactionType value to be stored in the database.
func (s TransactionType) Value() (driver.Value, error) {
	v := validator.New()
	errs := v.Var(s, "required,oneof=deposit withdrawal")
	if errs != nil {
		return nil, errs
	}

	return string(s), nil
}

// String converts the TransactionType value to a string.
func (s TransactionType) String() string {
	return string(s)
}

var PaymentMethods = []string{
	PaymentMethod_Bank, PaymentMethod_Cash,
}

// DailySummary is an object representing the database table.
type DailySummary struct {
	Income      float64 `boil:"income" json:"income" toml:"income" yaml:"income"`
	Expenditure float64 `boil:"expenditure" json:"expenditure" toml:"expenditure" yaml:"expenditure"`
	BankDeposit float64 `boil:"bank_deposit" json:"bank_deposit" toml:"bank_deposit" yaml:"bank_deposit"`
	Date        int64   `boil:"date" json:"date" toml:"date" yaml:"date"`
}
