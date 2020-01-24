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

// SaleItem is an object representing the database table.
type SaleItem struct {
	ID            string  `boil:"id" json:"id" toml:"id" yaml:"id"`
	SaleID        string  `boil:"sale_id" json:"sale_id" toml:"sale_id" yaml:"sale_id"`
	ProductID     string  `boil:"product_id" json:"product_id" toml:"product_id" yaml:"product_id"`
	UnitPrice     float64 `boil:"unit_price" json:"unit_price" toml:"unit_price" yaml:"unit_price"`
	UnitCostPrice float64 `boil:"unit_cost_price" json:"unit_cost_price" toml:"unit_cost_price" yaml:"unit_cost_price"`
	StockIds      string  `boil:"stock_ids" json:"stock_ids" toml:"stock_ids" yaml:"stock_ids"`

	R *saleItemR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L saleItemL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var SaleItemColumns = struct {
	ID            string
	SaleID        string
	ProductID     string
	UnitPrice     string
	UnitCostPrice string
	StockIds      string
}{
	ID:            "id",
	SaleID:        "sale_id",
	ProductID:     "product_id",
	UnitPrice:     "unit_price",
	UnitCostPrice: "unit_cost_price",
	StockIds:      "stock_ids",
}

// Generated where

var SaleItemWhere = struct {
	ID            whereHelperstring
	SaleID        whereHelperstring
	ProductID     whereHelperstring
	UnitPrice     whereHelperfloat64
	UnitCostPrice whereHelperfloat64
	StockIds      whereHelperstring
}{
	ID:            whereHelperstring{field: "\"sale_item\".\"id\""},
	SaleID:        whereHelperstring{field: "\"sale_item\".\"sale_id\""},
	ProductID:     whereHelperstring{field: "\"sale_item\".\"product_id\""},
	UnitPrice:     whereHelperfloat64{field: "\"sale_item\".\"unit_price\""},
	UnitCostPrice: whereHelperfloat64{field: "\"sale_item\".\"unit_cost_price\""},
	StockIds:      whereHelperstring{field: "\"sale_item\".\"stock_ids\""},
}

// SaleItemRels is where relationship names are stored.
var SaleItemRels = struct {
	Product string
	Sale    string
}{
	Product: "Product",
	Sale:    "Sale",
}

// saleItemR is where relationships are stored.
type saleItemR struct {
	Product *Product
	Sale    *Sale
}

// NewStruct creates a new relationship struct
func (*saleItemR) NewStruct() *saleItemR {
	return &saleItemR{}
}

// saleItemL is where Load methods for each relationship are stored.
type saleItemL struct{}

var (
	saleItemAllColumns            = []string{"id", "sale_id", "product_id", "unit_price", "unit_cost_price", "stock_ids"}
	saleItemColumnsWithoutDefault = []string{"id", "sale_id", "product_id", "unit_price", "unit_cost_price", "stock_ids"}
	saleItemColumnsWithDefault    = []string{}
	saleItemPrimaryKeyColumns     = []string{"id"}
)

type (
	// SaleItemSlice is an alias for a slice of pointers to SaleItem.
	// This should generally be used opposed to []SaleItem.
	SaleItemSlice []*SaleItem

	saleItemQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	saleItemType                 = reflect.TypeOf(&SaleItem{})
	saleItemMapping              = queries.MakeStructMapping(saleItemType)
	saleItemPrimaryKeyMapping, _ = queries.BindMapping(saleItemType, saleItemMapping, saleItemPrimaryKeyColumns)
	saleItemInsertCacheMut       sync.RWMutex
	saleItemInsertCache          = make(map[string]insertCache)
	saleItemUpdateCacheMut       sync.RWMutex
	saleItemUpdateCache          = make(map[string]updateCache)
	saleItemUpsertCacheMut       sync.RWMutex
	saleItemUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single saleItem record from the query.
func (q saleItemQuery) One(ctx context.Context, exec boil.ContextExecutor) (*SaleItem, error) {
	o := &SaleItem{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for sale_item")
	}

	return o, nil
}

// All returns all SaleItem records from the query.
func (q saleItemQuery) All(ctx context.Context, exec boil.ContextExecutor) (SaleItemSlice, error) {
	var o []*SaleItem

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to SaleItem slice")
	}

	return o, nil
}

// Count returns the count of all SaleItem records in the query.
func (q saleItemQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count sale_item rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q saleItemQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if sale_item exists")
	}

	return count > 0, nil
}

// Product pointed to by the foreign key.
func (o *SaleItem) Product(mods ...qm.QueryMod) productQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ProductID),
	}

	queryMods = append(queryMods, mods...)

	query := Products(queryMods...)
	queries.SetFrom(query.Query, "\"product\"")

	return query
}

// Sale pointed to by the foreign key.
func (o *SaleItem) Sale(mods ...qm.QueryMod) saleQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SaleID),
	}

	queryMods = append(queryMods, mods...)

	query := Sales(queryMods...)
	queries.SetFrom(query.Query, "\"sale\"")

	return query
}

// LoadProduct allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (saleItemL) LoadProduct(ctx context.Context, e boil.ContextExecutor, singular bool, maybeSaleItem interface{}, mods queries.Applicator) error {
	var slice []*SaleItem
	var object *SaleItem

	if singular {
		object = maybeSaleItem.(*SaleItem)
	} else {
		slice = *maybeSaleItem.(*[]*SaleItem)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &saleItemR{}
		}
		args = append(args, object.ProductID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &saleItemR{}
			}

			for _, a := range args {
				if a == obj.ProductID {
					continue Outer
				}
			}

			args = append(args, obj.ProductID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`product`), qm.WhereIn(`product.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Product")
	}

	var resultSlice []*Product
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Product")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for product")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for product")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Product = foreign
		if foreign.R == nil {
			foreign.R = &productR{}
		}
		foreign.R.SaleItems = append(foreign.R.SaleItems, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ProductID == foreign.ID {
				local.R.Product = foreign
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.SaleItems = append(foreign.R.SaleItems, local)
				break
			}
		}
	}

	return nil
}

// LoadSale allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (saleItemL) LoadSale(ctx context.Context, e boil.ContextExecutor, singular bool, maybeSaleItem interface{}, mods queries.Applicator) error {
	var slice []*SaleItem
	var object *SaleItem

	if singular {
		object = maybeSaleItem.(*SaleItem)
	} else {
		slice = *maybeSaleItem.(*[]*SaleItem)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &saleItemR{}
		}
		args = append(args, object.SaleID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &saleItemR{}
			}

			for _, a := range args {
				if a == obj.SaleID {
					continue Outer
				}
			}

			args = append(args, obj.SaleID)

		}
	}

	if len(args) == 0 {
		return nil
	}

	query := NewQuery(qm.From(`sale`), qm.WhereIn(`sale.id in ?`, args...))
	if mods != nil {
		mods.Apply(query)
	}

	results, err := query.QueryContext(ctx, e)
	if err != nil {
		return errors.Wrap(err, "failed to eager load Sale")
	}

	var resultSlice []*Sale
	if err = queries.Bind(results, &resultSlice); err != nil {
		return errors.Wrap(err, "failed to bind eager loaded slice Sale")
	}

	if err = results.Close(); err != nil {
		return errors.Wrap(err, "failed to close results of eager load for sale")
	}
	if err = results.Err(); err != nil {
		return errors.Wrap(err, "error occurred during iteration of eager loaded relations for sale")
	}

	if len(resultSlice) == 0 {
		return nil
	}

	if singular {
		foreign := resultSlice[0]
		object.R.Sale = foreign
		if foreign.R == nil {
			foreign.R = &saleR{}
		}
		foreign.R.SaleItems = append(foreign.R.SaleItems, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SaleID == foreign.ID {
				local.R.Sale = foreign
				if foreign.R == nil {
					foreign.R = &saleR{}
				}
				foreign.R.SaleItems = append(foreign.R.SaleItems, local)
				break
			}
		}
	}

	return nil
}

// SetProduct of the saleItem to the related item.
// Sets o.R.Product to related.
// Adds o to related.R.SaleItems.
func (o *SaleItem) SetProduct(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Product) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"sale_item\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"product_id"}),
		strmangle.WhereClause("\"", "\"", 2, saleItemPrimaryKeyColumns),
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

	o.ProductID = related.ID
	if o.R == nil {
		o.R = &saleItemR{
			Product: related,
		}
	} else {
		o.R.Product = related
	}

	if related.R == nil {
		related.R = &productR{
			SaleItems: SaleItemSlice{o},
		}
	} else {
		related.R.SaleItems = append(related.R.SaleItems, o)
	}

	return nil
}

// SetSale of the saleItem to the related item.
// Sets o.R.Sale to related.
// Adds o to related.R.SaleItems.
func (o *SaleItem) SetSale(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Sale) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"sale_item\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"sale_id"}),
		strmangle.WhereClause("\"", "\"", 2, saleItemPrimaryKeyColumns),
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

	o.SaleID = related.ID
	if o.R == nil {
		o.R = &saleItemR{
			Sale: related,
		}
	} else {
		o.R.Sale = related
	}

	if related.R == nil {
		related.R = &saleR{
			SaleItems: SaleItemSlice{o},
		}
	} else {
		related.R.SaleItems = append(related.R.SaleItems, o)
	}

	return nil
}

// SaleItems retrieves all the records using an executor.
func SaleItems(mods ...qm.QueryMod) saleItemQuery {
	mods = append(mods, qm.From("\"sale_item\""))
	return saleItemQuery{NewQuery(mods...)}
}

// FindSaleItem retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindSaleItem(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*SaleItem, error) {
	saleItemObj := &SaleItem{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"sale_item\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, saleItemObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from sale_item")
	}

	return saleItemObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *SaleItem) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no sale_item provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(saleItemColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	saleItemInsertCacheMut.RLock()
	cache, cached := saleItemInsertCache[key]
	saleItemInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			saleItemAllColumns,
			saleItemColumnsWithDefault,
			saleItemColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(saleItemType, saleItemMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(saleItemType, saleItemMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"sale_item\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"sale_item\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into sale_item")
	}

	if !cached {
		saleItemInsertCacheMut.Lock()
		saleItemInsertCache[key] = cache
		saleItemInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the SaleItem.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *SaleItem) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	saleItemUpdateCacheMut.RLock()
	cache, cached := saleItemUpdateCache[key]
	saleItemUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			saleItemAllColumns,
			saleItemPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update sale_item, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"sale_item\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, saleItemPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(saleItemType, saleItemMapping, append(wl, saleItemPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update sale_item row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for sale_item")
	}

	if !cached {
		saleItemUpdateCacheMut.Lock()
		saleItemUpdateCache[key] = cache
		saleItemUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q saleItemQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for sale_item")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for sale_item")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o SaleItemSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), saleItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"sale_item\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, saleItemPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in saleItem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all saleItem")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *SaleItem) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no sale_item provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(saleItemColumnsWithDefault, o)

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

	saleItemUpsertCacheMut.RLock()
	cache, cached := saleItemUpsertCache[key]
	saleItemUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			saleItemAllColumns,
			saleItemColumnsWithDefault,
			saleItemColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			saleItemAllColumns,
			saleItemPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert sale_item, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(saleItemPrimaryKeyColumns))
			copy(conflict, saleItemPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"sale_item\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(saleItemType, saleItemMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(saleItemType, saleItemMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert sale_item")
	}

	if !cached {
		saleItemUpsertCacheMut.Lock()
		saleItemUpsertCache[key] = cache
		saleItemUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single SaleItem record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *SaleItem) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no SaleItem provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), saleItemPrimaryKeyMapping)
	sql := "DELETE FROM \"sale_item\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from sale_item")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for sale_item")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q saleItemQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no saleItemQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from sale_item")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for sale_item")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o SaleItemSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), saleItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"sale_item\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, saleItemPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from saleItem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for sale_item")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *SaleItem) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindSaleItem(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *SaleItemSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := SaleItemSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), saleItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"sale_item\".* FROM \"sale_item\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, saleItemPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in SaleItemSlice")
	}

	*o = slice

	return nil
}

// SaleItemExists checks if the SaleItem row exists.
func SaleItemExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"sale_item\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if sale_item exists")
	}

	return exists, nil
}
