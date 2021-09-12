// Code generated by SQLBoiler 4.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"github.com/volatiletech/sqlboiler/v4/queries/qmhelper"
	"github.com/volatiletech/strmangle"
)

// Profit is an object representing the database table.
type Profit struct {
	ID         string     `boil:"id" json:"id" toml:"id" yaml:"id"`
	Amount     float64    `boil:"amount" json:"amount" toml:"amount" yaml:"amount"`
	Narration  string     `boil:"narration" json:"narration" toml:"narration" yaml:"narration"`
	CreatedAt  int64      `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt  int64      `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	ArchivedAt null.Int64 `boil:"archived_at" json:"archived_at,omitempty" toml:"archived_at" yaml:"archived_at,omitempty"`

	R *profitR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L profitL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var ProfitColumns = struct {
	ID         string
	Amount     string
	Narration  string
	CreatedAt  string
	UpdatedAt  string
	ArchivedAt string
}{
	ID:         "id",
	Amount:     "amount",
	Narration:  "narration",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	ArchivedAt: "archived_at",
}

var ProfitTableColumns = struct {
	ID         string
	Amount     string
	Narration  string
	CreatedAt  string
	UpdatedAt  string
	ArchivedAt string
}{
	ID:         "profit.id",
	Amount:     "profit.amount",
	Narration:  "profit.narration",
	CreatedAt:  "profit.created_at",
	UpdatedAt:  "profit.updated_at",
	ArchivedAt: "profit.archived_at",
}

// Generated where

var ProfitWhere = struct {
	ID         whereHelperstring
	Amount     whereHelperfloat64
	Narration  whereHelperstring
	CreatedAt  whereHelperint64
	UpdatedAt  whereHelperint64
	ArchivedAt whereHelpernull_Int64
}{
	ID:         whereHelperstring{field: "\"profit\".\"id\""},
	Amount:     whereHelperfloat64{field: "\"profit\".\"amount\""},
	Narration:  whereHelperstring{field: "\"profit\".\"narration\""},
	CreatedAt:  whereHelperint64{field: "\"profit\".\"created_at\""},
	UpdatedAt:  whereHelperint64{field: "\"profit\".\"updated_at\""},
	ArchivedAt: whereHelpernull_Int64{field: "\"profit\".\"archived_at\""},
}

// ProfitRels is where relationship names are stored.
var ProfitRels = struct {
}{}

// profitR is where relationships are stored.
type profitR struct {
}

// NewStruct creates a new relationship struct
func (*profitR) NewStruct() *profitR {
	return &profitR{}
}

// profitL is where Load methods for each relationship are stored.
type profitL struct{}

var (
	profitAllColumns            = []string{"id", "amount", "narration", "created_at", "updated_at", "archived_at"}
	profitColumnsWithoutDefault = []string{"id", "amount", "created_at", "updated_at", "archived_at"}
	profitColumnsWithDefault    = []string{"narration"}
	profitPrimaryKeyColumns     = []string{"id"}
)

type (
	// ProfitSlice is an alias for a slice of pointers to Profit.
	// This should almost always be used instead of []Profit.
	ProfitSlice []*Profit

	profitQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	profitType                 = reflect.TypeOf(&Profit{})
	profitMapping              = queries.MakeStructMapping(profitType)
	profitPrimaryKeyMapping, _ = queries.BindMapping(profitType, profitMapping, profitPrimaryKeyColumns)
	profitInsertCacheMut       sync.RWMutex
	profitInsertCache          = make(map[string]insertCache)
	profitUpdateCacheMut       sync.RWMutex
	profitUpdateCache          = make(map[string]updateCache)
	profitUpsertCacheMut       sync.RWMutex
	profitUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single profit record from the query.
func (q profitQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Profit, error) {
	o := &Profit{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for profit")
	}

	return o, nil
}

// All returns all Profit records from the query.
func (q profitQuery) All(ctx context.Context, exec boil.ContextExecutor) (ProfitSlice, error) {
	var o []*Profit

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Profit slice")
	}

	return o, nil
}

// Count returns the count of all Profit records in the query.
func (q profitQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count profit rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q profitQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if profit exists")
	}

	return count > 0, nil
}

// Profits retrieves all the records using an executor.
func Profits(mods ...qm.QueryMod) profitQuery {
	mods = append(mods, qm.From("\"profit\""))
	return profitQuery{NewQuery(mods...)}
}

// FindProfit retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindProfit(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Profit, error) {
	profitObj := &Profit{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"profit\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, profitObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from profit")
	}

	return profitObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Profit) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no profit provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(profitColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	profitInsertCacheMut.RLock()
	cache, cached := profitInsertCache[key]
	profitInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			profitAllColumns,
			profitColumnsWithDefault,
			profitColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(profitType, profitMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(profitType, profitMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"profit\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"profit\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into profit")
	}

	if !cached {
		profitInsertCacheMut.Lock()
		profitInsertCache[key] = cache
		profitInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Profit.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Profit) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	profitUpdateCacheMut.RLock()
	cache, cached := profitUpdateCache[key]
	profitUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			profitAllColumns,
			profitPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update profit, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"profit\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, profitPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(profitType, profitMapping, append(wl, profitPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update profit row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for profit")
	}

	if !cached {
		profitUpdateCacheMut.Lock()
		profitUpdateCache[key] = cache
		profitUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q profitQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for profit")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for profit")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o ProfitSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), profitPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"profit\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, profitPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in profit slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all profit")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Profit) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no profit provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(profitColumnsWithDefault, o)

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

	profitUpsertCacheMut.RLock()
	cache, cached := profitUpsertCache[key]
	profitUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			profitAllColumns,
			profitColumnsWithDefault,
			profitColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			profitAllColumns,
			profitPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert profit, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(profitPrimaryKeyColumns))
			copy(conflict, profitPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"profit\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(profitType, profitMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(profitType, profitMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert profit")
	}

	if !cached {
		profitUpsertCacheMut.Lock()
		profitUpsertCache[key] = cache
		profitUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Profit record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Profit) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Profit provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), profitPrimaryKeyMapping)
	sql := "DELETE FROM \"profit\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from profit")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for profit")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q profitQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no profitQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from profit")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for profit")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o ProfitSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), profitPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"profit\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, profitPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from profit slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for profit")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Profit) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindProfit(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *ProfitSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := ProfitSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), profitPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"profit\".* FROM \"profit\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, profitPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in ProfitSlice")
	}

	*o = slice

	return nil
}

// ProfitExists checks if the Profit row exists.
func ProfitExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"profit\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if profit exists")
	}

	return exists, nil
}
