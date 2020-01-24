// Code generated by SQLBoiler 3.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Customer is an object representing the database table.
type Customer struct {
	ID          string    `boil:"id" json:"id" toml:"id" yaml:"id"`
	Email       string    `boil:"email" json:"email" toml:"email" yaml:"email"`
	Name        string    `boil:"name" json:"name" toml:"name" yaml:"name"`
	PhoneNumber string    `boil:"phone_number" json:"phone_number" toml:"phone_number" yaml:"phone_number"`
	Address     string    `boil:"address" json:"address" toml:"address" yaml:"address"`
	SalesRepID  string    `boil:"sales_rep_id" json:"sales_rep_id" toml:"sales_rep_id" yaml:"sales_rep_id"`
	CreatedAt   time.Time `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	ArchivedAt  null.Time `boil:"archived_at" json:"archived_at,omitempty" toml:"archived_at" yaml:"archived_at,omitempty"`
	BranchID    string    `boil:"branch_id" json:"branch_id" toml:"branch_id" yaml:"branch_id"`
	UpdatedAt   time.Time `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`

	R *customerR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L customerL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var CustomerColumns = struct {
	ID          string
	Email       string
	Name        string
	PhoneNumber string
	Address     string
	SalesRepID  string
	CreatedAt   string
	ArchivedAt  string
	BranchID    string
	UpdatedAt   string
}{
	ID:          "id",
	Email:       "email",
	Name:        "name",
	PhoneNumber: "phone_number",
	Address:     "address",
	SalesRepID:  "sales_rep_id",
	CreatedAt:   "created_at",
	ArchivedAt:  "archived_at",
	BranchID:    "branch_id",
	UpdatedAt:   "updated_at",
}

// Generated where

var CustomerWhere = struct {
	ID          whereHelperstring
	Email       whereHelperstring
	Name        whereHelperstring
	PhoneNumber whereHelperstring
	Address     whereHelperstring
	SalesRepID  whereHelperstring
	CreatedAt   whereHelpertime_Time
	ArchivedAt  whereHelpernull_Time
	BranchID    whereHelperstring
	UpdatedAt   whereHelpertime_Time
}{
	ID:          whereHelperstring{field: "\"customer\".\"id\""},
	Email:       whereHelperstring{field: "\"customer\".\"email\""},
	Name:        whereHelperstring{field: "\"customer\".\"name\""},
	PhoneNumber: whereHelperstring{field: "\"customer\".\"phone_number\""},
	Address:     whereHelperstring{field: "\"customer\".\"address\""},
	SalesRepID:  whereHelperstring{field: "\"customer\".\"sales_rep_id\""},
	CreatedAt:   whereHelpertime_Time{field: "\"customer\".\"created_at\""},
	ArchivedAt:  whereHelpernull_Time{field: "\"customer\".\"archived_at\""},
	BranchID:    whereHelperstring{field: "\"customer\".\"branch_id\""},
	UpdatedAt:   whereHelpertime_Time{field: "\"customer\".\"updated_at\""},
}

// CustomerRels is where relationship names are stored.
var CustomerRels = struct {
	Branch   string
	SalesRep string
	Accounts string
}{
	Branch:   "Branch",
	SalesRep: "SalesRep",
	Accounts: "Accounts",
}

// customerR is where relationships are stored.
type customerR struct {
	Branch   *Branch
	SalesRep *User
	Accounts AccountSlice
}

// NewStruct creates a new relationship struct
func (*customerR) NewStruct() *customerR {
	return &customerR{}
}

// customerL is where Load methods for each relationship are stored.
type customerL struct{}

var (
	customerAllColumns            = []string{"id", "email", "name", "phone_number", "address", "sales_rep_id", "created_at", "archived_at", "branch_id", "updated_at"}
	customerColumnsWithoutDefault = []string{"id", "email", "phone_number", "address", "sales_rep_id", "created_at", "archived_at", "updated_at"}
	customerColumnsWithDefault    = []string{"name", "branch_id"}
	customerPrimaryKeyColumns     = []string{"id"}
)

type (
	// CustomerSlice is an alias for a slice of pointers to Customer.
	// This should generally be used opposed to []Customer.
	CustomerSlice []*Customer

	customerQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	customerType                 = reflect.TypeOf(&Customer{})
	customerMapping              = queries.MakeStructMapping(customerType)
	customerPrimaryKeyMapping, _ = queries.BindMapping(customerType, customerMapping, customerPrimaryKeyColumns)
	customerInsertCacheMut       sync.RWMutex
	customerInsertCache          = make(map[string]insertCache)
	customerUpdateCacheMut       sync.RWMutex
	customerUpdateCache          = make(map[string]updateCache)
	customerUpsertCacheMut       sync.RWMutex
	customerUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single customer record from the query.
func (q customerQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Customer, error) {
	o := &Customer{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for customer")
	}

	return o, nil
}

// All returns all Customer records from the query.
func (q customerQuery) All(ctx context.Context, exec boil.ContextExecutor) (CustomerSlice, error) {
	var o []*Customer

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Customer slice")
	}

	return o, nil
}

// Count returns the count of all Customer records in the query.
func (q customerQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count customer rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q customerQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if customer exists")
	}

	return count > 0, nil
}

// Branch pointed to by the foreign key.
func (o *Customer) Branch(mods ...qm.QueryMod) branchQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.BranchID),
	}

	queryMods = append(queryMods, mods...)

	query := Branches(queryMods...)
	queries.SetFrom(query.Query, "\"branch\"")

	return query
}

// SalesRep pointed to by the foreign key.
func (o *Customer) SalesRep(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SalesRepID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// Accounts retrieves all the account's Accounts with an executor.
func (o *Customer) Accounts(mods ...qm.QueryMod) accountQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"account\".\"customer_id\"=?", o.ID),
	)

	query := Accounts(queryMods...)
	queries.SetFrom(query.Query, "\"account\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"account\".*"})
	}

	return query
}

// LoadBranch allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (customerL) LoadBranch(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCustomer interface{}, mods queries.Applicator) error {
	var slice []*Customer
	var object *Customer

	if singular {
		object = maybeCustomer.(*Customer)
	} else {
		slice = *maybeCustomer.(*[]*Customer)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &customerR{}
		}
		args = append(args, object.BranchID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &customerR{}
			}

			for _, a := range args {
				if a == obj.BranchID {
					continue Outer
				}
			}

			args = append(args, obj.BranchID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`branch`), qm.WhereIn(`branch.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Branch")
	}

	var resultSlice []*Branch
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Branch")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for branch")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for branch")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Branch = foreign
		if foreign.R == nil {
			foreign.R = &branchR{}
		}
		foreign.R.Customers = append(foreign.R.Customers, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.BranchID == foreign.ID {
				local.R.Branch = foreign
				if foreign.R == nil {
					foreign.R = &branchR{}
				}
				foreign.R.Customers = append(foreign.R.Customers, local)
				break
			}
		}
	}

	return nil
}

// LoadSalesRep allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (customerL) LoadSalesRep(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCustomer interface{}, mods queries.Applicator) error {
	var slice []*Customer
	var object *Customer

	if singular {
		object = maybeCustomer.(*Customer)
	} else {
		slice = *maybeCustomer.(*[]*Customer)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &customerR{}
		}
		args = append(args, object.SalesRepID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &customerR{}
			}

			for _, a := range args {
				if a == obj.SalesRepID {
					continue Outer
				}
			}

			args = append(args, obj.SalesRepID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`users`), qm.WhereIn(`users.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load User")
	}

	var resultSlice []*User
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice User")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for users")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for users")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.SalesRep = foreign
		if foreign.R == nil {
			foreign.R = &userR{}
		}
		foreign.R.SalesRepCustomers = append(foreign.R.SalesRepCustomers, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SalesRepID == foreign.ID {
				local.R.SalesRep = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.SalesRepCustomers = append(foreign.R.SalesRepCustomers, local)
				break
			}
		}
	}

	return nil
}

// LoadAccounts allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (customerL) LoadAccounts(ctx context.Context, e boil.ContextExecutor, singular bool, maybeCustomer interface{}, mods queries.Applicator) error {
	var slice []*Customer
	var object *Customer

	if singular {
		object = maybeCustomer.(*Customer)
	} else {
		slice = *maybeCustomer.(*[]*Customer)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &customerR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &customerR{}
			}

			for _, a := range args {
				if a == obj.ID {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`account`), qm.WhereIn(`account.customer_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load account")
	}

	var resultSlice []*Account
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice account")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on account")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for account")
	}

	if singular {
		object.R.Accounts = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &accountR{}
			}
			foreign.R.Customer = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if local.ID == foreign.CustomerID {
				local.R.Accounts = append(local.R.Accounts, foreign)
				if foreign.R == nil {
					foreign.R = &accountR{}
				}
				foreign.R.Customer = local
				break
			}
		}
	}

	return nil
}

// SetBranch of the customer to the related item.
// Sets o.R.Branch to related.
// Adds o to related.R.Customers.
func (o *Customer) SetBranch(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Branch) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"customer\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"branch_id"}),
		strmangle.WhereClause("\"", "\"", 2, customerPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.BranchID = related.ID
	if o.R == nil {
		o.R = &customerR{
			Branch: related,
		}
	} else {
		o.R.Branch = related
	}

	if related.R == nil {
		related.R = &branchR{
			Customers: CustomerSlice{o},
		}
	} else {
		related.R.Customers = append(related.R.Customers, o)
	}

	return nil
}

// SetSalesRep of the customer to the related item.
// Sets o.R.SalesRep to related.
// Adds o to related.R.SalesRepCustomers.
func (o *Customer) SetSalesRep(ctx context.Context, exec boil.ContextExecutor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"customer\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"sales_rep_id"}),
		strmangle.WhereClause("\"", "\"", 2, customerPrimaryKeyColumns),
	)
	values := []interface{}{related.ID, o.ID}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, updateQuery)
		fmt.Fprintln(writer, values)
	}
	if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
		return errors.Wrap(err, "failed to update local table")
	}

	o.SalesRepID = related.ID
	if o.R == nil {
		o.R = &customerR{
			SalesRep: related,
		}
	} else {
		o.R.SalesRep = related
	}

	if related.R == nil {
		related.R = &userR{
			SalesRepCustomers: CustomerSlice{o},
		}
	} else {
		related.R.SalesRepCustomers = append(related.R.SalesRepCustomers, o)
	}

	return nil
}

// AddAccounts adds the given related objects to the existing relationships
// of the customer, optionally inserting them as new records.
// Appends related to o.R.Accounts.
// Sets related.R.Customer appropriately.
func (o *Customer) AddAccounts(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Account) error {
	var err error
	for _, rel := range related {
		if insert {
			rel.CustomerID = o.ID
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"account\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"customer_id"}),
				strmangle.WhereClause("\"", "\"", 2, accountPrimaryKeyColumns),
			)
			values := []interface{}{o.ID, rel.ID}

			if boil.IsDebug(ctx) {
				writer := boil.DebugWriterFrom(ctx)
				fmt.Fprintln(writer, updateQuery)
				fmt.Fprintln(writer, values)
			}
			if _, err = exec.ExecContext(ctx, updateQuery, values...); err != nil {
				return errors.Wrap(err, "failed to update foreign table")
			}

			rel.CustomerID = o.ID
		}
	}

	if o.R == nil {
		o.R = &customerR{
			Accounts: related,
		}
	} else {
		o.R.Accounts = append(o.R.Accounts, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &accountR{
				Customer: o,
			}
		} else {
			rel.R.Customer = o
		}
	}
	return nil
}

// Customers retrieves all the records using an executor.
func Customers(mods ...qm.QueryMod) customerQuery {
	mods = append(mods, qm.From("\"customer\""))
	return customerQuery{NewQuery(mods...)}
}

// FindCustomer retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindCustomer(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Customer, error) {
	customerObj := &Customer{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"customer\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, customerObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from customer")
	}

	return customerObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Customer) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no customer provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(customerColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	customerInsertCacheMut.RLock()
	cache, cached := customerInsertCache[key]
	customerInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			customerAllColumns,
			customerColumnsWithDefault,
			customerColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(customerType, customerMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(customerType, customerMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"customer\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"customer\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into customer")
	}

	if !cached {
		customerInsertCacheMut.Lock()
		customerInsertCache[key] = cache
		customerInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Customer.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Customer) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	customerUpdateCacheMut.RLock()
	cache, cached := customerUpdateCache[key]
	customerUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			customerAllColumns,
			customerPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update customer, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"customer\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, customerPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(customerType, customerMapping, append(wl, customerPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, values)
	}
	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update customer row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for customer")
	}

	if !cached {
		customerUpdateCacheMut.Lock()
		customerUpdateCache[key] = cache
		customerUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q customerQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for customer")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for customer")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o CustomerSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	ln := int64(len(o))
	if ln == 0 {
		return 0, nil
	}

	if len(cols) == 0 {
		return 0, errors.New("models: update all requires at least one column argument")
	}

	colNames := make([]string, len(cols))
	args := make([]interface{}, len(cols))

	i := 0
	for name, value := range cols {
		colNames[i] = name
		args[i] = value
		i++
	}

	// Append all of the primary key values for each column
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), customerPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"customer\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, customerPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in customer slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all customer")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Customer) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no customer provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(customerColumnsWithDefault, o)

	// Build cache key in-line uglily - mysql vs psql problems
	buf := strmangle.GetBuffer()
	if updateOnConflict {
		buf.WriteByte('t')
	} else {
		buf.WriteByte('f')
	}
	buf.WriteByte('.')
	for _, c := range conflictColumns {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(updateColumns.Kind))
	for _, c := range updateColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	buf.WriteString(strconv.Itoa(insertColumns.Kind))
	for _, c := range insertColumns.Cols {
		buf.WriteString(c)
	}
	buf.WriteByte('.')
	for _, c := range nzDefaults {
		buf.WriteString(c)
	}
	key := buf.String()
	strmangle.PutBuffer(buf)

	customerUpsertCacheMut.RLock()
	cache, cached := customerUpsertCache[key]
	customerUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			customerAllColumns,
			customerColumnsWithDefault,
			customerColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			customerAllColumns,
			customerPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert customer, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(customerPrimaryKeyColumns))
			copy(conflict, customerPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"customer\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(customerType, customerMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(customerType, customerMapping, ret)
			if err != nil {
				return err
			}
		}
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)
	var returns []interface{}
	if len(cache.retMapping) != 0 {
		returns = queries.PtrsFromMapping(value, cache.retMapping)
	}

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, cache.query)
		fmt.Fprintln(writer, vals)
	}
	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(returns...)
		if err == sql.ErrNoRows {
			err = nil // Postgres doesn't return anything when there's no update
		}
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}
	if err != nil {
		return errors.Wrap(err, "models: unable to upsert customer")
	}

	if !cached {
		customerUpsertCacheMut.Lock()
		customerUpsertCache[key] = cache
		customerUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Customer record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Customer) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Customer provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), customerPrimaryKeyMapping)
	sql := "DELETE FROM \"customer\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from customer")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for customer")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q customerQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no customerQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from customer")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for customer")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o CustomerSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), customerPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"customer\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, customerPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from customer slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for customer")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Customer) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindCustomer(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *CustomerSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := CustomerSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), customerPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"customer\".* FROM \"customer\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, customerPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in CustomerSlice")
	}

	*o = slice

	return nil
}

// CustomerExists checks if the Customer row exists.
func CustomerExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"customer\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if customer exists")
	}

	return exists, nil
}
