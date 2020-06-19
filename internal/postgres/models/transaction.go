// Code generated by SQLBoiler 3.7.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// Transaction is an object representing the database table.
type Transaction struct {
	ID             string     `boil:"id" json:"id" toml:"id" yaml:"id"`
	AccountID      string     `boil:"account_id" json:"account_id" toml:"account_id" yaml:"account_id"`
	TXType         string     `boil:"tx_type" json:"tx_type" toml:"tx_type" yaml:"tx_type"`
	OpeningBalance float64    `boil:"opening_balance" json:"opening_balance" toml:"opening_balance" yaml:"opening_balance"`
	Amount         float64    `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Narration      string     `boil:"narration" json:"narration" toml:"narration" yaml:"narration"`
	SalesRepID     string     `boil:"sales_rep_id" json:"sales_rep_id" toml:"sales_rep_id" yaml:"sales_rep_id"`
	CreatedAt      int64      `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt      int64      `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	ArchivedAt     null.Int64 `boil:"archived_at" json:"archived_at,omitempty" toml:"archived_at" yaml:"archived_at,omitempty"`
	ReceiptNo      string     `boil:"receipt_no" json:"receipt_no" toml:"receipt_no" yaml:"receipt_no"`
	EffectiveDate  int64      `boil:"effective_date" json:"effective_date" toml:"effective_date" yaml:"effective_date"`

	R *transactionR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L transactionL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var TransactionColumns = struct {
	ID             string
	AccountID      string
	TXType         string
	OpeningBalance string
	Amount         string
	Narration      string
	SalesRepID     string
	CreatedAt      string
	UpdatedAt      string
	ArchivedAt     string
	ReceiptNo      string
	EffectiveDate  string
}{
	ID:             "id",
	AccountID:      "account_id",
	TXType:         "tx_type",
	OpeningBalance: "opening_balance",
	Amount:         "amount",
	Narration:      "narration",
	SalesRepID:     "sales_rep_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	ArchivedAt:     "archived_at",
	ReceiptNo:      "receipt_no",
	EffectiveDate:  "effective_date",
}

// Generated where

var TransactionWhere = struct {
	ID             whereHelperstring
	AccountID      whereHelperstring
	TXType         whereHelperstring
	OpeningBalance whereHelperfloat64
	Amount         whereHelperfloat64
	Narration      whereHelperstring
	SalesRepID     whereHelperstring
	CreatedAt      whereHelperint64
	UpdatedAt      whereHelperint64
	ArchivedAt     whereHelpernull_Int64
	ReceiptNo      whereHelperstring
	EffectiveDate  whereHelperint64
}{
	ID:             whereHelperstring{field: "\"transaction\".\"id\""},
	AccountID:      whereHelperstring{field: "\"transaction\".\"account_id\""},
	TXType:         whereHelperstring{field: "\"transaction\".\"tx_type\""},
	OpeningBalance: whereHelperfloat64{field: "\"transaction\".\"opening_balance\""},
	Amount:         whereHelperfloat64{field: "\"transaction\".\"amount\""},
	Narration:      whereHelperstring{field: "\"transaction\".\"narration\""},
	SalesRepID:     whereHelperstring{field: "\"transaction\".\"sales_rep_id\""},
	CreatedAt:      whereHelperint64{field: "\"transaction\".\"created_at\""},
	UpdatedAt:      whereHelperint64{field: "\"transaction\".\"updated_at\""},
	ArchivedAt:     whereHelpernull_Int64{field: "\"transaction\".\"archived_at\""},
	ReceiptNo:      whereHelperstring{field: "\"transaction\".\"receipt_no\""},
	EffectiveDate:  whereHelperint64{field: "\"transaction\".\"effective_date\""},
}

// TransactionRels is where relationship names are stored.
var TransactionRels = struct {
	Account  string
	SalesRep string
}{
	Account:  "Account",
	SalesRep: "SalesRep",
}

// transactionR is where relationships are stored.
type transactionR struct {
	Account  *Account
	SalesRep *User
}

// NewStruct creates a new relationship struct
func (*transactionR) NewStruct() *transactionR {
	return &transactionR{}
}

// transactionL is where Load methods for each relationship are stored.
type transactionL struct{}

var (
	transactionAllColumns            = []string{"id", "account_id", "tx_type", "opening_balance", "amount", "narration", "sales_rep_id", "created_at", "updated_at", "archived_at", "receipt_no", "effective_date"}
	transactionColumnsWithoutDefault = []string{"id", "tx_type", "opening_balance", "sales_rep_id", "created_at", "updated_at", "archived_at"}
	transactionColumnsWithDefault    = []string{"account_id", "amount", "narration", "receipt_no", "effective_date"}
	transactionPrimaryKeyColumns     = []string{"id"}
)

type (
	// TransactionSlice is an alias for a slice of pointers to Transaction.
	// This should generally be used opposed to []Transaction.
	TransactionSlice []*Transaction

	transactionQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	transactionType                 = reflect.TypeOf(&Transaction{})
	transactionMapping              = queries.MakeStructMapping(transactionType)
	transactionPrimaryKeyMapping, _ = queries.BindMapping(transactionType, transactionMapping, transactionPrimaryKeyColumns)
	transactionInsertCacheMut       sync.RWMutex
	transactionInsertCache          = make(map[string]insertCache)
	transactionUpdateCacheMut       sync.RWMutex
	transactionUpdateCache          = make(map[string]updateCache)
	transactionUpsertCacheMut       sync.RWMutex
	transactionUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single transaction record from the query.
func (q transactionQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Transaction, error) {
	o := &Transaction{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for transaction")
	}

	return o, nil
}

// All returns all Transaction records from the query.
func (q transactionQuery) All(ctx context.Context, exec boil.ContextExecutor) (TransactionSlice, error) {
	var o []*Transaction

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Transaction slice")
	}

	return o, nil
}

// Count returns the count of all Transaction records in the query.
func (q transactionQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count transaction rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q transactionQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if transaction exists")
	}

	return count > 0, nil
}

// Account pointed to by the foreign key.
func (o *Transaction) Account(mods ...qm.QueryMod) accountQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.AccountID),
	}

	queryMods = append(queryMods, mods...)

	query := Accounts(queryMods...)
	queries.SetFrom(query.Query, "\"account\"")

	return query
}

// SalesRep pointed to by the foreign key.
func (o *Transaction) SalesRep(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SalesRepID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadAccount allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (transactionL) LoadAccount(ctx context.Context, e boil.ContextExecutor, singular bool, maybeTransaction interface{}, mods queries.Applicator) error {
	var slice []*Transaction
	var object *Transaction

	if singular {
		object = maybeTransaction.(*Transaction)
	} else {
		slice = *maybeTransaction.(*[]*Transaction)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &transactionR{}
		}
		args = append(args, object.AccountID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &transactionR{}
			}

			for _, a := range args {
				if a == obj.AccountID {
					continue Outer
				}
			}

			args = append(args, obj.AccountID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`account`), qm.WhereIn(`account.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Account")
	}

	var resultSlice []*Account
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Account")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for account")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for account")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Account = foreign
		if foreign.R == nil {
			foreign.R = &accountR{}
		}
		foreign.R.Transactions = append(foreign.R.Transactions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.AccountID == foreign.ID {
				local.R.Account = foreign
				if foreign.R == nil {
					foreign.R = &accountR{}
				}
				foreign.R.Transactions = append(foreign.R.Transactions, local)
				break
			}
		}
	}

	return nil
}

// LoadSalesRep allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (transactionL) LoadSalesRep(ctx context.Context, e boil.ContextExecutor, singular bool, maybeTransaction interface{}, mods queries.Applicator) error {
	var slice []*Transaction
	var object *Transaction

	if singular {
		object = maybeTransaction.(*Transaction)
	} else {
		slice = *maybeTransaction.(*[]*Transaction)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &transactionR{}
		}
		args = append(args, object.SalesRepID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &transactionR{}
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
		foreign.R.SalesRepTransactions = append(foreign.R.SalesRepTransactions, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SalesRepID == foreign.ID {
				local.R.SalesRep = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.SalesRepTransactions = append(foreign.R.SalesRepTransactions, local)
				break
			}
		}
	}

	return nil
}

// SetAccount of the transaction to the related item.
// Sets o.R.Account to related.
// Adds o to related.R.Transactions.
func (o *Transaction) SetAccount(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Account) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"transaction\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"account_id"}),
		strmangle.WhereClause("\"", "\"", 2, transactionPrimaryKeyColumns),
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

	o.AccountID = related.ID
	if o.R == nil {
		o.R = &transactionR{
			Account: related,
		}
	} else {
		o.R.Account = related
	}

	if related.R == nil {
		related.R = &accountR{
			Transactions: TransactionSlice{o},
		}
	} else {
		related.R.Transactions = append(related.R.Transactions, o)
	}

	return nil
}

// SetSalesRep of the transaction to the related item.
// Sets o.R.SalesRep to related.
// Adds o to related.R.SalesRepTransactions.
func (o *Transaction) SetSalesRep(ctx context.Context, exec boil.ContextExecutor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"transaction\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"sales_rep_id"}),
		strmangle.WhereClause("\"", "\"", 2, transactionPrimaryKeyColumns),
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
		o.R = &transactionR{
			SalesRep: related,
		}
	} else {
		o.R.SalesRep = related
	}

	if related.R == nil {
		related.R = &userR{
			SalesRepTransactions: TransactionSlice{o},
		}
	} else {
		related.R.SalesRepTransactions = append(related.R.SalesRepTransactions, o)
	}

	return nil
}

// Transactions retrieves all the records using an executor.
func Transactions(mods ...qm.QueryMod) transactionQuery {
	mods = append(mods, qm.From("\"transaction\""))
	return transactionQuery{NewQuery(mods...)}
}

// FindTransaction retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindTransaction(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Transaction, error) {
	transactionObj := &Transaction{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"transaction\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, transactionObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from transaction")
	}

	return transactionObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Transaction) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no transaction provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(transactionColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	transactionInsertCacheMut.RLock()
	cache, cached := transactionInsertCache[key]
	transactionInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			transactionAllColumns,
			transactionColumnsWithDefault,
			transactionColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(transactionType, transactionMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(transactionType, transactionMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"transaction\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"transaction\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into transaction")
	}

	if !cached {
		transactionInsertCacheMut.Lock()
		transactionInsertCache[key] = cache
		transactionInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Transaction.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Transaction) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	transactionUpdateCacheMut.RLock()
	cache, cached := transactionUpdateCache[key]
	transactionUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			transactionAllColumns,
			transactionPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update transaction, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"transaction\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, transactionPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(transactionType, transactionMapping, append(wl, transactionPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update transaction row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for transaction")
	}

	if !cached {
		transactionUpdateCacheMut.Lock()
		transactionUpdateCache[key] = cache
		transactionUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q transactionQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for transaction")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for transaction")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o TransactionSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), transactionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"transaction\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, transactionPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in transaction slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all transaction")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Transaction) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no transaction provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(transactionColumnsWithDefault, o)

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

	transactionUpsertCacheMut.RLock()
	cache, cached := transactionUpsertCache[key]
	transactionUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			transactionAllColumns,
			transactionColumnsWithDefault,
			transactionColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			transactionAllColumns,
			transactionPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert transaction, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(transactionPrimaryKeyColumns))
			copy(conflict, transactionPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"transaction\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(transactionType, transactionMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(transactionType, transactionMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert transaction")
	}

	if !cached {
		transactionUpsertCacheMut.Lock()
		transactionUpsertCache[key] = cache
		transactionUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Transaction record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Transaction) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Transaction provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), transactionPrimaryKeyMapping)
	sql := "DELETE FROM \"transaction\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from transaction")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for transaction")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q transactionQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no transactionQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from transaction")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for transaction")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o TransactionSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), transactionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"transaction\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, transactionPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from transaction slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for transaction")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Transaction) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindTransaction(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *TransactionSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := TransactionSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), transactionPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"transaction\".* FROM \"transaction\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, transactionPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in TransactionSlice")
	}

	*o = slice

	return nil
}

// TransactionExists checks if the Transaction row exists.
func TransactionExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"transaction\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if transaction exists")
	}

	return exists, nil
}
