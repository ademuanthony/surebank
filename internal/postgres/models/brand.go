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
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// Brand is an object representing the database table.
type Brand struct {
	ID   string `boil:"id" json:"id" toml:"id" yaml:"id"`
	Name string `boil:"name" json:"name" toml:"name" yaml:"name"`
	Code string `boil:"code" json:"code" toml:"code" yaml:"code"`
	Logo string `boil:"logo" json:"logo" toml:"logo" yaml:"logo"`

	R *brandR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L brandL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var BrandColumns = struct {
	ID   string
	Name string
	Code string
	Logo string
}{
	ID:   "id",
	Name: "name",
	Code: "code",
	Logo: "logo",
}

// Generated where

var BrandWhere = struct {
	ID   whereHelperstring
	Name whereHelperstring
	Code whereHelperstring
	Logo whereHelperstring
}{
	ID:   whereHelperstring{field: "\"brand\".\"id\""},
	Name: whereHelperstring{field: "\"brand\".\"name\""},
	Code: whereHelperstring{field: "\"brand\".\"code\""},
	Logo: whereHelperstring{field: "\"brand\".\"logo\""},
}

// BrandRels is where relationship names are stored.
var BrandRels = struct {
	Products string
}{
	Products: "Products",
}

// brandR is where relationships are stored.
type brandR struct {
	Products ProductSlice
}

// NewStruct creates a new relationship struct
func (*brandR) NewStruct() *brandR {
	return &brandR{}
}

// brandL is where Load methods for each relationship are stored.
type brandL struct{}

var (
	brandAllColumns            = []string{"id", "name", "code", "logo"}
	brandColumnsWithoutDefault = []string{"id", "name", "code", "logo"}
	brandColumnsWithDefault    = []string{}
	brandPrimaryKeyColumns     = []string{"id"}
)

type (
	// BrandSlice is an alias for a slice of pointers to Brand.
	// This should generally be used opposed to []Brand.
	BrandSlice []*Brand

	brandQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	brandType                 = reflect.TypeOf(&Brand{})
	brandMapping              = queries.MakeStructMapping(brandType)
	brandPrimaryKeyMapping, _ = queries.BindMapping(brandType, brandMapping, brandPrimaryKeyColumns)
	brandInsertCacheMut       sync.RWMutex
	brandInsertCache          = make(map[string]insertCache)
	brandUpdateCacheMut       sync.RWMutex
	brandUpdateCache          = make(map[string]updateCache)
	brandUpsertCacheMut       sync.RWMutex
	brandUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single brand record from the query.
func (q brandQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Brand, error) {
	o := &Brand{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for brand")
	}

	return o, nil
}

// All returns all Brand records from the query.
func (q brandQuery) All(ctx context.Context, exec boil.ContextExecutor) (BrandSlice, error) {
	var o []*Brand

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Brand slice")
	}

	return o, nil
}

// Count returns the count of all Brand records in the query.
func (q brandQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count brand rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q brandQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if brand exists")
	}

	return count > 0, nil
}

// Products retrieves all the product's Products with an executor.
func (o *Brand) Products(mods ...qm.QueryMod) productQuery {
	var queryMods []qm.QueryMod
	if len(mods) != 0 {
		queryMods = append(queryMods, mods...)
	}

	queryMods = append(queryMods,
		qm.Where("\"product\".\"brand_id\"=?", o.ID),
	)

	query := Products(queryMods...)
	queries.SetFrom(query.Query, "\"product\"")

	if len(queries.GetSelect(query.Query)) == 0 {
		queries.SetSelect(query.Query, []string{"\"product\".*"})
	}

	return query
}

// LoadProducts allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for a 1-M or N-M relationship.
func (brandL) LoadProducts(ctx context.Context, e boil.ContextExecutor, singular bool, maybeBrand interface{}, mods queries.Applicator) error {
	var slice []*Brand
	var object *Brand

	if singular {
		object = maybeBrand.(*Brand)
	} else {
		slice = *maybeBrand.(*[]*Brand)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &brandR{}
		}
		args = append(args, object.ID)
	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &brandR{}
			}

			for _, a := range args {
				if queries.Equal(a, obj.ID) {
					continue Outer
				}
			}

			args = append(args, obj.ID)
		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`product`), qm.WhereIn(`product.brand_id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load product")
	}

	var resultSlice []*Product
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice product")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results in eager load on product")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for product")
	}

	if singular {
		object.R.Products = resultSlice
		for _, foreign := range resultSlice {
			if foreign.R == nil {
				foreign.R = &productR{}
			}
			foreign.R.Brand = object
		}
		return nil
	}

	for _, foreign := range resultSlice {
		for _, local := range slice {
			if queries.Equal(local.ID, foreign.BrandID) {
				local.R.Products = append(local.R.Products, foreign)
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.Brand = local
				break
			}
		}
	}

	return nil
}

// AddProducts adds the given related objects to the existing relationships
// of the brand, optionally inserting them as new records.
// Appends related to o.R.Products.
// Sets related.R.Brand appropriately.
func (o *Brand) AddProducts(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Product) error {
	var err error
	for _, rel := range related {
		if insert {
			queries.Assign(&rel.BrandID, o.ID)
			if err = rel.Insert(ctx, exec, boil.Infer()); err != nil {
				return errors.Wrap(err, "failed to insert into foreign table")
			}
		} else {
			updateQuery := fmt.Sprintf(
				"UPDATE \"product\" SET %s WHERE %s",
				strmangle.SetParamNames("\"", "\"", 1, []string{"brand_id"}),
				strmangle.WhereClause("\"", "\"", 2, productPrimaryKeyColumns),
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

			queries.Assign(&rel.BrandID, o.ID)
		}
	}

	if o.R == nil {
		o.R = &brandR{
			Products: related,
		}
	} else {
		o.R.Products = append(o.R.Products, related...)
	}

	for _, rel := range related {
		if rel.R == nil {
			rel.R = &productR{
				Brand: o,
			}
		} else {
			rel.R.Brand = o
		}
	}
	return nil
}

// SetProducts removes all previously related items of the
// brand replacing them completely with the passed
// in related items, optionally inserting them as new records.
// Sets o.R.Brand's Products accordingly.
// Replaces o.R.Products with related.
// Sets related.R.Brand's Products accordingly.
func (o *Brand) SetProducts(ctx context.Context, exec boil.ContextExecutor, insert bool, related ...*Product) error {
	query := "update \"product\" set \"brand_id\" = null where \"brand_id\" = $1"
	values := []interface{}{o.ID}
	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, query)
		fmt.Fprintln(writer, values)
	}
	_, err := exec.ExecContext(ctx, query, values...)
	if err != nil {
		return errors.Wrap(err, "failed to remove relationships before set")
	}

	if o.R != nil {
		for _, rel := range o.R.Products {
			queries.SetScanner(&rel.BrandID, nil)
			if rel.R == nil {
				continue
			}

			rel.R.Brand = nil
		}

		o.R.Products = nil
	}
	return o.AddProducts(ctx, exec, insert, related...)
}

// RemoveProducts relationships from objects passed in.
// Removes related items from R.Products (uses pointer comparison, removal does not keep order)
// Sets related.R.Brand.
func (o *Brand) RemoveProducts(ctx context.Context, exec boil.ContextExecutor, related ...*Product) error {
	var err error
	for _, rel := range related {
		queries.SetScanner(&rel.BrandID, nil)
		if rel.R != nil {
			rel.R.Brand = nil
		}
		if _, err = rel.Update(ctx, exec, boil.Whitelist("brand_id")); err != nil {
			return err
		}
	}
	if o.R == nil {
		return nil
	}

	for _, rel := range related {
		for i, ri := range o.R.Products {
			if rel != ri {
				continue
			}

			ln := len(o.R.Products)
			if ln > 1 && i < ln-1 {
				o.R.Products[i] = o.R.Products[ln-1]
			}
			o.R.Products = o.R.Products[:ln-1]
			break
		}
	}

	return nil
}

// Brands retrieves all the records using an executor.
func Brands(mods ...qm.QueryMod) brandQuery {
	mods = append(mods, qm.From("\"brand\""))
	return brandQuery{NewQuery(mods...)}
}

// FindBrand retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindBrand(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Brand, error) {
	brandObj := &Brand{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"brand\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, brandObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from brand")
	}

	return brandObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Brand) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no brand provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(brandColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	brandInsertCacheMut.RLock()
	cache, cached := brandInsertCache[key]
	brandInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			brandAllColumns,
			brandColumnsWithDefault,
			brandColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(brandType, brandMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(brandType, brandMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"brand\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"brand\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into brand")
	}

	if !cached {
		brandInsertCacheMut.Lock()
		brandInsertCache[key] = cache
		brandInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Brand.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Brand) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	brandUpdateCacheMut.RLock()
	cache, cached := brandUpdateCache[key]
	brandUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			brandAllColumns,
			brandPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update brand, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"brand\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, brandPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(brandType, brandMapping, append(wl, brandPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update brand row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for brand")
	}

	if !cached {
		brandUpdateCacheMut.Lock()
		brandUpdateCache[key] = cache
		brandUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q brandQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for brand")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for brand")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o BrandSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), brandPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"brand\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, brandPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in brand slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all brand")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Brand) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no brand provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(brandColumnsWithDefault, o)

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

	brandUpsertCacheMut.RLock()
	cache, cached := brandUpsertCache[key]
	brandUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			brandAllColumns,
			brandColumnsWithDefault,
			brandColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			brandAllColumns,
			brandPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert brand, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(brandPrimaryKeyColumns))
			copy(conflict, brandPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"brand\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(brandType, brandMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(brandType, brandMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert brand")
	}

	if !cached {
		brandUpsertCacheMut.Lock()
		brandUpsertCache[key] = cache
		brandUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Brand record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Brand) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Brand provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), brandPrimaryKeyMapping)
	sql := "DELETE FROM \"brand\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from brand")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for brand")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q brandQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no brandQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from brand")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for brand")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o BrandSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), brandPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"brand\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, brandPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from brand slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for brand")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Brand) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindBrand(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *BrandSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := BrandSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), brandPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"brand\".* FROM \"brand\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, brandPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in BrandSlice")
	}

	*o = slice

	return nil
}

// BrandExists checks if the Brand row exists.
func BrandExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"brand\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if brand exists")
	}

	return exists, nil
}
