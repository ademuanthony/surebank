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
	BankAccount:     "bank_account",
	BankDeposit:     "bank_deposit",
	Branch:          "branch",
	Brand:           "brand",
	Category:        "category",
	Customer:        "customer",
	DailySummary:    "daily_summary",
	DSCommission:    "ds_commission",
	Expenditure:     "expenditure",
	Inventory:       "inventory",
	Payment:         "payment",
	Product:         "product",
	ProductCategory: "product_category",
	RepsExpense:     "reps_expense",
	Sale:            "sale",
	SaleItem:        "sale_item",
	Transaction:     "transaction",
	Users:           "users",
}

var (
	AccountColumns = struct {
		ID              string
		BranchID        string
		Number          string
		CustomerID      string
		AccountType     string
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
		BranchID:        "branch_id",
		Number:          "number",
		CustomerID:      "customer_id",
		AccountType:     "account_type",
		Target:          "target",
		TargetInfo:      "target_info",
		SalesRepID:      "sales_rep_id",
		CreatedAt:       "created_at",
		UpdatedAt:       "updated_at",
		ArchivedAt:      "archived_at",
		Balance:         "balance",
		LastPaymentDate: "last_payment_date",
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
		CreatedAt:  "created_at",
		UpdatedAt:  "updated_at",
		ArchivedAt: "archived_at",
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
		BranchID:    "branch_id",
		Email:       "email",
		Name:        "name",
		PhoneNumber: "phone_number",
		Address:     "address",
		SalesRepID:  "sales_rep_id",
		CreatedAt:   "created_at",
		UpdatedAt:   "updated_at",
		ArchivedAt:  "archived_at",
		Accounts:    "accounts",
	}
)

type FindInput struct {
	FilteringQuery interface{}
	Offset         int
	Limit          int
	SortFields     []string
}

var client *mongo.Client

func Connect(ctx context.Context, server string) (err error) {
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s", server))
	client, err = mongo.Connect(ctx, clientOptions)
	return
}

func NewDb(databaseName ...string) *mongo.Database {
	dbName := "main"
	if len(databaseName) > 0 {
		dbName = databaseName[0]
	}
	return client.Database(dbName)
}
