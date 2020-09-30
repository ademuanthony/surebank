package transaction

import (
	"context"
	"database/sql/driver"
	"errors"
	"merryworld/surebank/internal/dscommission"
	"merryworld/surebank/internal/platform/notify"
	"sync"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null"

	"merryworld/surebank/internal/account"
	"merryworld/surebank/internal/platform/web"
	"merryworld/surebank/internal/postgres/models"
	"merryworld/surebank/internal/user"
)

// Repository defines the required dependencies for Transaction.
type Repository struct {
	DbConn         *sqlx.DB
	CommissionRepo *dscommission.Repository
	notifySMS      notify.SMS
	accNumMtx      sync.Mutex
	creatDB        func() (*sqlx.DB, error)
}

// NewRepository creates a new Repository that defines dependencies for Transaction.
func NewRepository(db *sqlx.DB, commissionRepo *dscommission.Repository, notifySMS notify.SMS, creatDB func() (*sqlx.DB, error)) *Repository {
	return &Repository{
		DbConn:         db,
		CommissionRepo: commissionRepo,
		notifySMS:      notifySMS,
		creatDB:        creatDB,
	}
}

// Transaction represents a financial transaction.
type Transaction struct {
	ID             string          `json:"id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	AccountID      string          `json:"account_id" example:"985f1746-1d9f-459f-a2d9-fc53ece5ae86"`
	Type           TransactionType `json:"type" example:"deposit"`
	OpeningBalance float64         `json:"opening_balance" exmaple:"34500.01"`
	Amount         float64         `json:"amount" truss:"api-read"`
	Narration      string          `json:"narration" truss:"api-read"`
	PaymentMethod  string          `json:"payment_method" truss:"api-read"`
	SalesRepID     string          `json:"sales_rep_id" truss:"api-read"`
	ReceiptNo      string          `json:"receipt_no"`
	EffectiveDate  time.Time       `json:"effective_date"`
	CreatedAt      time.Time       `json:"created_at" truss:"api-read"`
	UpdatedAt      time.Time       `json:"updated_at" truss:"api-read"`
	ArchivedAt     *time.Time      `json:"archived_at,omitempty" truss:"api-hide"`

	SalesRep *user.User       `json:"sales_rep" truss:"api-read"`
	Account  *account.Account `json:"account" truss:"api-read"`
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
		EffectiveDate:  time.Unix(rec.EffectiveDate, 0).UTC(),
		CreatedAt:      time.Unix(rec.CreatedAt, 0).UTC(),
		UpdatedAt:      time.Unix(rec.UpdatedAt, 0).UTC(),
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
func (m *Transaction) Response(ctx context.Context) *Response {
	if m == nil {
		return nil
	}

	r := &Response{
		ID:             m.ID,
		AccountID:      m.AccountID,
		Type:           m.Type,
		OpeningBalance: m.OpeningBalance,
		Amount:         m.Amount,
		Narration:      m.Narration,
		PaymentMethod:  m.PaymentMethod,
		ReceiptNo:      m.ReceiptNo,
		SalesRepID:     m.SalesRepID,
		EffectiveDate:  web.NewTimeResponse(ctx, m.EffectiveDate),
		CreatedAt:      web.NewTimeResponse(ctx, m.CreatedAt),
		UpdatedAt:      web.NewTimeResponse(ctx, m.UpdatedAt),
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
type Transactions []*Transaction

// Response transforms a list of Transactions to a list of Responses.
func (m *Transactions) Response(ctx context.Context) []*Response {
	var l = make([]*Response, 0)
	if m != nil && len(*m) > 0 {
		for _, n := range *m {
			l = append(l, n.Response(ctx))
		}
	}

	return l
}

// PagedResponseList holds list of transaction and total count for pagination
type PagedResponseList struct {
	Transactions []*Response `json:"transactions"`
	TotalCount   int64       `json:"total_count"`
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
	Where            string        `json:"where" example:"type = deposit and account_id = ? and created_at > ? and created_at < ?"`
	Args             []interface{} `json:"args" swaggertype:"array,string" example:"Moon Launch,active"`
	Order            []string      `json:"order" example:"created_at desc"`
	Limit            *uint         `json:"limit" example:"10"`
	Offset           *uint         `json:"offset" example:"20"`
	IncludeArchived  bool          `json:"include-archived" example:"false"`
	IncludeAccount   bool          `json:"include_account" example:"false"`
	IncludeCustomer  bool          `json:"include_customer"`
	IncludeAccountNo bool          `json:"include_account_no" example:"false"`
	IncludeSalesRep  bool          `json:"include_sales_rep" example:"false"`
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
	var l []interface{}
	for _, v := range TransactionType_Values {
		l = append(l, v.String())
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
