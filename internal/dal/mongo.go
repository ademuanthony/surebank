package dal

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var C = struct {
	Account         string
	BankAccount     string
	BankDeposit     string
	Branch          string
	Brand           string
	Category        string
	Customer        string
	DailySummary    string
	DSCommission    string
	Expenditure     string
	Inventory       string
	Payment         string
	Product         string
	ProductCategory string
	RepsExpense     string
	Sale            string
	SaleItem        string
	Transaction     string
	Users           string
}{
	Account:         "account",
	BankAccount:     "bankaccount",
	BankDeposit:     "bankdeposit",
	Branch:          "branch",
	Brand:           "brand",
	Category:        "category",
	Customer:        "customer",
	DailySummary:    "dailysummary",
	DSCommission:    "dscommission",
	Expenditure:     "expenditure",
	Inventory:       "inventory",
	Payment:         "payment",
	Product:         "product",
	ProductCategory: "productcategory",
	RepsExpense:     "repsexpense",
	Sale:            "sale",
	SaleItem:        "saleitem",
	Transaction:     "transaction",
	Users:           "users",
}

var (
	AccountColumns = struct {
		ID              string
		BranchID        string
		Number          string
		CustomerID      string
		Type            string
		Target          string
		TargetInfo      string
		SalesRepID      string
		CreatedAt       string
		UpdatedAt       string
		ArchivedAt      string
		Balance         string
		LastPaymentDate string
	}{
		ID:              "id",
		BranchID:        "branchid",
		Number:          "number",
		CustomerID:      "customerid",
		Type:            "type",
		Target:          "target",
		TargetInfo:      "targetinfo",
		SalesRepID:      "salesrepid",
		CreatedAt:       "createdat",
		UpdatedAt:       "updatedat",
		ArchivedAt:      "archivedat",
		Balance:         "balance",
		LastPaymentDate: "lastpaymentdate",
	}

	BranchColumns = struct {
		ID         string
		Name       string
		CreatedAt  string
		UpdatedAt  string
		ArchivedAt string
	}{
		ID:         "id",
		Name:       "name",
		CreatedAt:  "createdat",
		UpdatedAt:  "updatedat",
		ArchivedAt: "archivedat",
	}

	CustomerColumns = struct {
		ID          string
		BranchID    string
		Email       string
		Name        string
		PhoneNumber string
		Address     string
		SalesRepID  string
		CreatedAt   string
		UpdatedAt   string
		ArchivedAt  string
		Accounts    string
	}{
		ID:          "id",
		BranchID:    "branchid",
		Email:       "email",
		Name:        "name",
		PhoneNumber: "phonenumber",
		Address:     "address",
		SalesRepID:  "salesrepid",
		CreatedAt:   "createdat",
		UpdatedAt:   "updatedat",
		ArchivedAt:  "archivedat",
		Accounts:    "accounts",
	}

	DailySummaryColumns = struct {
		Income      string
		Expenditure string
		BankDeposit string
		Date        string
	}{
		Income:      "income",
		Expenditure: "expenditure",
		BankDeposit: "bankdeposit",
		Date:        "date",
	}

	DSCommissionColumns = struct {
		ID            string
		AccountID     string
		CustomerID    string
		Amount        string
		Date          string
		EffectiveDate string
	}{
		ID:            "id",
		AccountID:     "accountid",
		CustomerID:    "customerid",
		Amount:        "amount",
		Date:          "date",
		EffectiveDate: "effectivedate",
	}

	TransactionColumns = struct {
		ID             string
		AccountID      string
		AccountNumber  string
		CustomerID     string
		Type           string
		OpeningBalance string
		Amount         string
		Narration      string
		SalesRepID     string
		CreatedAt      string
		UpdatedAt      string
		ArchivedAt     string
		ReceiptNo      string
		EffectiveDate  string
		PaymentMethod  string
	}{
		ID:             "id",
		AccountID:      "accountid",
		AccountNumber:  "accountnumber",
		CustomerID:     "customerid",
		Type:           "type",
		OpeningBalance: "openingbalance",
		Amount:         "amount",
		Narration:      "narration",
		SalesRepID:     "salesrepid",
		CreatedAt:      "createdat",
		UpdatedAt:      "updatedat",
		ArchivedAt:     "archivedat",
		ReceiptNo:      "receiptno",
		EffectiveDate:  "effectivedate",
		PaymentMethod:  "paymentmethod",
	}
)

type FindInput struct {
	FilteringQuery interface{}
	Offset         int
	Limit          int
	SortFields     []string
}

var Client *mongo.Client

func Connect(ctx context.Context, server string) (err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", server))
	Client, err = mongo.Connect(ctx, clientOptions)
	return
}

func NewDb(databaseName ...string) *mongo.Database {
	dbName := "main"
	if len(databaseName) > 0 {
		dbName = databaseName[0]
	}
	return Client.Database(dbName)
}
