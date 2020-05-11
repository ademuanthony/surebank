// Code generated by SQLBoiler 4.1.1 (https://github.com/volatiletech/sqlboiler/v4). DO NOT EDIT.
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
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// BankDeposit is an object representing the database table.
type BankDeposit struct {
	ID            string  `boil:"id" json:"id" toml:"id" yaml:"id"`
	BankAccountID string  `boil:"bank_account_id" json:"bank_account_id" toml:"bank_account_id" yaml:"bank_account_id"`
	Amount        float64 `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Date          int64   `boil:"date" json:"date" toml:"date" yaml:"date"`

	R *bankDepositR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L bankDepositL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var BankDepositColumns = struct {
	ID            string
	BankAccountID string
	Amount        string
	Date          string
}{
	ID:            "id",
	BankAccountID: "bank_account_id",
	Amount:        "amount",
	Date:          "date",
}

// Generated where

var BankDepositWhere = struct {
	ID            whereHelperstring
	BankAccountID whereHelperstring
	Amount        whereHelperfloat64
	Date          whereHelperint64
}{
	ID:            whereHelperstring{field: "\"bank_deposit\".\"id\""},
	BankAccountID: whereHelperstring{field: "\"bank_deposit\".\"bank_account_id\""},
	Amount:        whereHelperfloat64{field: "\"bank_deposit\".\"amount\""},
	Date:          whereHelperint64{field: "\"bank_deposit\".\"date\""},
}

// BankDepositRels is where relationship names are stored.
var BankDepositRels = struct {
	BankAccount string
}{
	BankAccount: "BankAccount",
}

// bankDepositR is where relationships are stored.
type bankDepositR struct {
	BankAccount *BankAccount `boil:"BankAccount" json:"BankAccount" toml:"BankAccount" yaml:"BankAccount"`
}

// NewStruct creates a new relationship struct
func (*bankDepositR) NewStruct() *bankDepositR {
	return &bankDepositR{}
}

// bankDepositL is where Load methods for each relationship are stored.
type bankDepositL struct{}

var (
	bankDepositAllColumns            = []string{"id", "bank_account_id", "amount", "date"}
	bankDepositColumnsWithoutDefault = []string{"id", "bank_account_id", "amount", "date"}
	bankDepositColumnsWithDefault    = []string{}
	bankDepositPrimaryKeyColumns     = []string{"id"}
)

type (
	// BankDepositSlice is an alias for a slice of pointers to BankDeposit.
	// This should generally be used opposed to []BankDeposit.
	BankDepositSlice []*BankDeposit

	bankDepositQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	bankDepositType                 = reflect.TypeOf(&BankDeposit{})
	bankDepositMapping              = queries.MakeStructMapping(bankDepositType)
	bankDepositPrimaryKeyMapping, _ = queries.BindMapping(bankDepositType, bankDepositMapping, bankDepositPrimaryKeyColumns)
	bankDepositInsertCacheMut       sync.RWMutex
	bankDepositInsertCache          = make(map[string]insertCache)
	bankDepositUpdateCacheMut       sync.RWMutex
	bankDepositUpdateCache          = make(map[string]updateCache)
	bankDepositUpsertCacheMut       sync.RWMutex
	bankDepositUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single bankDeposit record from the query.
func (q bankDepositQuery) One(ctx context.Context, exec boil.ContextExecutor) (*BankDeposit, error) {
	o := &BankDeposit{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for bank_deposit")
	}

	return o, nil
}

// All returns all BankDeposit records from the query.
func (q bankDepositQuery) All(ctx context.Context, exec boil.ContextExecutor) (BankDepositSlice, error) {
	var o []*BankDeposit

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to BankDeposit slice")
	}

	return o, nil
}

// Count returns the count of all BankDeposit records in the query.
func (q bankDepositQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count bank_deposit rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q bankDepositQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if bank_deposit exists")
	}

	return count > 0, nil
}

// BankAccount pointed to by the foreign key.
func (o *BankDeposit) BankAccount(mods ...qm.QueryMod) bankAccountQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.BankAccountID),
	}

	queryMods = append(queryMods, mods...)

	query := BankAccounts(queryMods...)
	queries.SetFrom(query.Query, "\"bank_account\"")

	return query
}

// LoadBankAccount allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (bankDepositL) LoadBankAccount(ctx context.Context, e boil.ContextExecutor, singular bool, maybeBankDeposit interface{}, mods queries.Applicator) error {
	var slice []*BankDeposit
	var object *BankDeposit

	if singular {
		object = maybeBankDeposit.(*BankDeposit)
	} else {
		slice = *maybeBankDeposit.(*[]*BankDeposit)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &bankDepositR{}
		}
		args = append(args, object.BankAccountID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &bankDepositR{}
			}

			for _, a := range args {
				if a == obj.BankAccountID {
					continue Outer
				}
			}

			args = append(args, obj.BankAccountID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(
		qm.From(`bank_account`),
		qm.WhereIn(`bank_account.id in ?`, args...),
	)
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load BankAccount")
	}

	var resultSlice []*BankAccount
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice BankAccount")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for bank_account")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for bank_account")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.BankAccount = foreign
		if foreign.R == nil {
			foreign.R = &bankAccountR{}
		}
		foreign.R.BankDeposits = append(foreign.R.BankDeposits, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.BankAccountID == foreign.ID {
				local.R.BankAccount = foreign
				if foreign.R == nil {
					foreign.R = &bankAccountR{}
				}
				foreign.R.BankDeposits = append(foreign.R.BankDeposits, local)
				break
			}
		}
	}

	return nil
}

// SetBankAccount of the bankDeposit to the related item.
// Sets o.R.BankAccount to related.
// Adds o to related.R.BankDeposits.
func (o *BankDeposit) SetBankAccount(ctx context.Context, exec boil.ContextExecutor, insert bool, related *BankAccount) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"bank_deposit\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"bank_account_id"}),
		strmangle.WhereClause("\"", "\"", 2, bankDepositPrimaryKeyColumns),
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

	o.BankAccountID = related.ID
	if o.R == nil {
		o.R = &bankDepositR{
			BankAccount: related,
		}
	} else {
		o.R.BankAccount = related
	}

	if related.R == nil {
		related.R = &bankAccountR{
			BankDeposits: BankDepositSlice{o},
		}
	} else {
		related.R.BankDeposits = append(related.R.BankDeposits, o)
	}

	return nil
}

// BankDeposits retrieves all the records using an executor.
func BankDeposits(mods ...qm.QueryMod) bankDepositQuery {
	mods = append(mods, qm.From("\"bank_deposit\""))
	return bankDepositQuery{NewQuery(mods...)}
}

// FindBankDeposit retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindBankDeposit(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*BankDeposit, error) {
	bankDepositObj := &BankDeposit{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"bank_deposit\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, bankDepositObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from bank_deposit")
	}

	return bankDepositObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *BankDeposit) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no bank_deposit provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(bankDepositColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	bankDepositInsertCacheMut.RLock()
	cache, cached := bankDepositInsertCache[key]
	bankDepositInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			bankDepositAllColumns,
			bankDepositColumnsWithDefault,
			bankDepositColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(bankDepositType, bankDepositMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(bankDepositType, bankDepositMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"bank_deposit\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"bank_deposit\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into bank_deposit")
	}

	if !cached {
		bankDepositInsertCacheMut.Lock()
		bankDepositInsertCache[key] = cache
		bankDepositInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the BankDeposit.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *BankDeposit) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	bankDepositUpdateCacheMut.RLock()
	cache, cached := bankDepositUpdateCache[key]
	bankDepositUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			bankDepositAllColumns,
			bankDepositPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update bank_deposit, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"bank_deposit\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, bankDepositPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(bankDepositType, bankDepositMapping, append(wl, bankDepositPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update bank_deposit row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for bank_deposit")
	}

	if !cached {
		bankDepositUpdateCacheMut.Lock()
		bankDepositUpdateCache[key] = cache
		bankDepositUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q bankDepositQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for bank_deposit")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for bank_deposit")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o BankDepositSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), bankDepositPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"bank_deposit\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, bankDepositPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in bankDeposit slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all bankDeposit")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *BankDeposit) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no bank_deposit provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(bankDepositColumnsWithDefault, o)

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

	bankDepositUpsertCacheMut.RLock()
	cache, cached := bankDepositUpsertCache[key]
	bankDepositUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			bankDepositAllColumns,
			bankDepositColumnsWithDefault,
			bankDepositColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			bankDepositAllColumns,
			bankDepositPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert bank_deposit, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(bankDepositPrimaryKeyColumns))
			copy(conflict, bankDepositPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"bank_deposit\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(bankDepositType, bankDepositMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(bankDepositType, bankDepositMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert bank_deposit")
	}

	if !cached {
		bankDepositUpsertCacheMut.Lock()
		bankDepositUpsertCache[key] = cache
		bankDepositUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single BankDeposit record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *BankDeposit) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no BankDeposit provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), bankDepositPrimaryKeyMapping)
	sql := "DELETE FROM \"bank_deposit\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from bank_deposit")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for bank_deposit")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q bankDepositQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no bankDepositQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from bank_deposit")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for bank_deposit")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o BankDepositSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), bankDepositPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"bank_deposit\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, bankDepositPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from bankDeposit slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for bank_deposit")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *BankDeposit) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindBankDeposit(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BankDepositSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := BankDepositSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), bankDepositPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"bank_deposit\".* FROM \"bank_deposit\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, bankDepositPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in BankDepositSlice")
	}

	*o = slice

	return nil
}

// BankDepositExists checks if the BankDeposit row exists.
func BankDepositExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"bank_deposit\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if bank_deposit exists")
	}

	return exists, nil
}
