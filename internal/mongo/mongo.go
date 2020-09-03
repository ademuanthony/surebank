package mongo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Collections = struct {
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

var CustomerColumns = struct {
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
}{
	ID:          "_id",
	BranchID:    "branch_id",
	Email:       "email",
	Name:        "name",
	PhoneNumber: "phone_number",
	Address:     "address",
	SalesRepID:  "sales_rep_id",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	ArchivedAt:  "archived_at",
}

type FindInput struct {
	FilteringQuery interface{}
	Offset         int
	Limit          int
	SortFields     []string
}

var session *mgo.Session

func Connect(server string) (err error) {
	session, err = mgo.Dial(server)
	return
}

func NewDb(databaseName ...string) *mgo.Database {
	newSession := session.Copy()
	dbName := "main"
	if len(databaseName) > 0 {
		dbName = databaseName[0]
	}
	db := newSession.DB(dbName)
	return db
}

func Insert(collection *mgo.Collection, obj ...interface{}) error {
	err := collection.Insert(obj...)
	return err
}

func Save(collection *mgo.Collection, obj interface{}) error {
	err := collection.Insert(obj)
	return err
}

func Update(collection *mgo.Collection, selector interface{}, obj interface{}) error {
	err := collection.Update(selector, obj)
	return err
}

func Patch(collection *mgo.Collection, selector interface{}, changes interface{}) error {
	_, err := collection.UpdateAll(selector, changes)
	return err
}

func FindAll(collection *mgo.Collection, input FindInput, receiver interface{}) (err error) {
	query := collection.Find(input.FilteringQuery)
	if len(input.SortFields) > 0 {
		query = query.Sort(input.SortFields...)
	}
	if input.Limit > 0 {
		query = query.Limit(input.Limit)
	}
	if input.Offset != 0 {
		query = query.Skip(input.Offset)
	}
	err = query.All(receiver)
	return
}

func FindOne(collection *mgo.Collection, filteringQuery interface{}, receiver interface{}) (err error) {
	query := collection.Find(filteringQuery)
	err = query.One(receiver)
	return
}

func FindById(collection *mgo.Collection, id bson.ObjectId, receiver interface{}) (err error) {
	err = collection.FindId(id).One(receiver)
	return
}

func Exists(collection *mgo.Collection, filteringQuery interface{}) (bool, error) {
	count, err := collection.Find(filteringQuery).Count()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func Count(collection *mgo.Collection, filteringQuery interface{}) (int, error) {
	return collection.Find(filteringQuery).Count()
}
