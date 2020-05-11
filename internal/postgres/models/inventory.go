// Code generated by SQLBoiler 4.1.1 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

// Inventory is an object representing the database table.
type Inventory struct {
	ID             string     `boil:"id" json:"id" toml:"id" yaml:"id"`
	ProductID      string     `boil:"product_id" json:"product_id" toml:"product_id" yaml:"product_id"`
	BranchID       string     `boil:"branch_id" json:"branch_id" toml:"branch_id" yaml:"branch_id"`
	TXType         string     `boil:"tx_type" json:"tx_type" toml:"tx_type" yaml:"tx_type"`
	OpeningBalance float64    `boil:"opening_balance" json:"opening_balance" toml:"opening_balance" yaml:"opening_balance"`
	Quantity       float64    `boil:"quantity" json:"quantity" toml:"quantity" yaml:"quantity"`
	Narration      string     `boil:"narration" json:"narration" toml:"narration" yaml:"narration"`
	SalesRepID     string     `boil:"sales_rep_id" json:"sales_rep_id" toml:"sales_rep_id" yaml:"sales_rep_id"`
	CreatedAt      int64      `boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt      int64      `boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
	ArchivedAt     null.Int64 `boil:"archived_at" json:"archived_at,omitempty" toml:"archived_at" yaml:"archived_at,omitempty"`

	R *inventoryR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L inventoryL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var InventoryColumns = struct {
	ID             string
	ProductID      string
	BranchID       string
	TXType         string
	OpeningBalance string
	Quantity       string
	Narration      string
	SalesRepID     string
	CreatedAt      string
	UpdatedAt      string
	ArchivedAt     string
}{
	ID:             "id",
	ProductID:      "product_id",
	BranchID:       "branch_id",
	TXType:         "tx_type",
	OpeningBalance: "opening_balance",
	Quantity:       "quantity",
	Narration:      "narration",
	SalesRepID:     "sales_rep_id",
	CreatedAt:      "created_at",
	UpdatedAt:      "updated_at",
	ArchivedAt:     "archived_at",
}

// Generated where

var InventoryWhere = struct {
	ID             whereHelperstring
	ProductID      whereHelperstring
	BranchID       whereHelperstring
	TXType         whereHelperstring
	OpeningBalance whereHelperfloat64
	Quantity       whereHelperfloat64
	Narration      whereHelperstring
	SalesRepID     whereHelperstring
	CreatedAt      whereHelperint64
	UpdatedAt      whereHelperint64
	ArchivedAt     whereHelpernull_Int64
}{
	ID:             whereHelperstring{field: "\"inventory\".\"id\""},
	ProductID:      whereHelperstring{field: "\"inventory\".\"product_id\""},
	BranchID:       whereHelperstring{field: "\"inventory\".\"branch_id\""},
	TXType:         whereHelperstring{field: "\"inventory\".\"tx_type\""},
	OpeningBalance: whereHelperfloat64{field: "\"inventory\".\"opening_balance\""},
	Quantity:       whereHelperfloat64{field: "\"inventory\".\"quantity\""},
	Narration:      whereHelperstring{field: "\"inventory\".\"narration\""},
	SalesRepID:     whereHelperstring{field: "\"inventory\".\"sales_rep_id\""},
	CreatedAt:      whereHelperint64{field: "\"inventory\".\"created_at\""},
	UpdatedAt:      whereHelperint64{field: "\"inventory\".\"updated_at\""},
	ArchivedAt:     whereHelpernull_Int64{field: "\"inventory\".\"archived_at\""},
}

// InventoryRels is where relationship names are stored.
var InventoryRels = struct {
	Branch   string
	Product  string
	SalesRep string
}{
	Branch:   "Branch",
	Product:  "Product",
	SalesRep: "SalesRep",
}

// inventoryR is where relationships are stored.
type inventoryR struct {
	Branch   *Branch  `boil:"Branch" json:"Branch" toml:"Branch" yaml:"Branch"`
	Product  *Product `boil:"Product" json:"Product" toml:"Product" yaml:"Product"`
	SalesRep *User    `boil:"SalesRep" json:"SalesRep" toml:"SalesRep" yaml:"SalesRep"`
}

// NewStruct creates a new relationship struct
func (*inventoryR) NewStruct() *inventoryR {
	return &inventoryR{}
}

// inventoryL is where Load methods for each relationship are stored.
type inventoryL struct{}

var (
	inventoryAllColumns            = []string{"id", "product_id", "branch_id", "tx_type", "opening_balance", "quantity", "narration", "sales_rep_id", "created_at", "updated_at", "archived_at"}
	inventoryColumnsWithoutDefault = []string{"id", "tx_type", "opening_balance", "sales_rep_id", "created_at", "updated_at", "archived_at"}
	inventoryColumnsWithDefault    = []string{"product_id", "branch_id", "quantity", "narration"}
	inventoryPrimaryKeyColumns     = []string{"id"}
)

type (
	// InventorySlice is an alias for a slice of pointers to Inventory.
	// This should generally be used opposed to []Inventory.
	InventorySlice []*Inventory

	inventoryQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	inventoryType                 = reflect.TypeOf(&Inventory{})
	inventoryMapping              = queries.MakeStructMapping(inventoryType)
	inventoryPrimaryKeyMapping, _ = queries.BindMapping(inventoryType, inventoryMapping, inventoryPrimaryKeyColumns)
	inventoryInsertCacheMut       sync.RWMutex
	inventoryInsertCache          = make(map[string]insertCache)
	inventoryUpdateCacheMut       sync.RWMutex
	inventoryUpdateCache          = make(map[string]updateCache)
	inventoryUpsertCacheMut       sync.RWMutex
	inventoryUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// One returns a single inventory record from the query.
func (q inventoryQuery) One(ctx context.Context, exec boil.ContextExecutor) (*Inventory, error) {
	o := &Inventory{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for inventory")
	}

	return o, nil
}

// All returns all Inventory records from the query.
func (q inventoryQuery) All(ctx context.Context, exec boil.ContextExecutor) (InventorySlice, error) {
	var o []*Inventory

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to Inventory slice")
	}

	return o, nil
}

// Count returns the count of all Inventory records in the query.
func (q inventoryQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count inventory rows")
	}

	return count, nil
}

// Exists checks if the row exists in the table.
func (q inventoryQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if inventory exists")
	}

	return count > 0, nil
}

// Branch pointed to by the foreign key.
func (o *Inventory) Branch(mods ...qm.QueryMod) branchQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.BranchID),
	}

	queryMods = append(queryMods, mods...)

	query := Branches(queryMods...)
	queries.SetFrom(query.Query, "\"branch\"")

	return query
}

// Product pointed to by the foreign key.
func (o *Inventory) Product(mods ...qm.QueryMod) productQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.ProductID),
	}

	queryMods = append(queryMods, mods...)

	query := Products(queryMods...)
	queries.SetFrom(query.Query, "\"product\"")

	return query
}

// SalesRep pointed to by the foreign key.
func (o *Inventory) SalesRep(mods ...qm.QueryMod) userQuery {
	queryMods := []qm.QueryMod{
		qm.Where("\"id\" = ?", o.SalesRepID),
	}

	queryMods = append(queryMods, mods...)

	query := Users(queryMods...)
	queries.SetFrom(query.Query, "\"users\"")

	return query
}

// LoadBranch allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (inventoryL) LoadBranch(ctx context.Context, e boil.ContextExecutor, singular bool, maybeInventory interface{}, mods queries.Applicator) error {
	var slice []*Inventory
	var object *Inventory

	if singular {
		object = maybeInventory.(*Inventory)
	} else {
		slice = *maybeInventory.(*[]*Inventory)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &inventoryR{}
		}
		args = append(args, object.BranchID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &inventoryR{}
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

	query := NewQuery(
		qm.From(`branch`),
		qm.WhereIn(`branch.id in ?`, args...),
	)
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
		foreign.R.Inventories = append(foreign.R.Inventories, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.BranchID == foreign.ID {
				local.R.Branch = foreign
				if foreign.R == nil {
					foreign.R = &branchR{}
				}
				foreign.R.Inventories = append(foreign.R.Inventories, local)
				break
			}
		}
	}

	return nil
}

// LoadProduct allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (inventoryL) LoadProduct(ctx context.Context, e boil.ContextExecutor, singular bool, maybeInventory interface{}, mods queries.Applicator) error {
	var slice []*Inventory
	var object *Inventory

	if singular {
		object = maybeInventory.(*Inventory)
	} else {
		slice = *maybeInventory.(*[]*Inventory)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &inventoryR{}
		}
		args = append(args, object.ProductID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &inventoryR{}
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

	query := NewQuery(
		qm.From(`product`),
		qm.WhereIn(`product.id in ?`, args...),
	)
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
		foreign.R.Inventories = append(foreign.R.Inventories, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.ProductID == foreign.ID {
				local.R.Product = foreign
				if foreign.R == nil {
					foreign.R = &productR{}
				}
				foreign.R.Inventories = append(foreign.R.Inventories, local)
				break
			}
		}
	}

	return nil
}

// LoadSalesRep allows an eager lookup of values, cached into the
// loaded structs of the objects. This is for an N-1 relationship.
func (inventoryL) LoadSalesRep(ctx context.Context, e boil.ContextExecutor, singular bool, maybeInventory interface{}, mods queries.Applicator) error {
	var slice []*Inventory
	var object *Inventory

	if singular {
		object = maybeInventory.(*Inventory)
	} else {
		slice = *maybeInventory.(*[]*Inventory)
	}

	args := make([]interface{}, 0, 1)
	if singular {
		if object.R == nil {
			object.R = &inventoryR{}
		}
		args = append(args, object.SalesRepID)

	} else {
	Outer:
		for _, obj := range slice {
			if obj.R == nil {
				obj.R = &inventoryR{}
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

	query := NewQuery(
		qm.From(`users`),
		qm.WhereIn(`users.id in ?`, args...),
	)
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
		foreign.R.SalesRepInventories = append(foreign.R.SalesRepInventories, object)
		return nil
	}

	for _, local := range slice {
		for _, foreign := range resultSlice {
			if local.SalesRepID == foreign.ID {
				local.R.SalesRep = foreign
				if foreign.R == nil {
					foreign.R = &userR{}
				}
				foreign.R.SalesRepInventories = append(foreign.R.SalesRepInventories, local)
				break
			}
		}
	}

	return nil
}

// SetBranch of the inventory to the related item.
// Sets o.R.Branch to related.
// Adds o to related.R.Inventories.
func (o *Inventory) SetBranch(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Branch) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"inventory\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"branch_id"}),
		strmangle.WhereClause("\"", "\"", 2, inventoryPrimaryKeyColumns),
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
		o.R = &inventoryR{
			Branch: related,
		}
	} else {
		o.R.Branch = related
	}

	if related.R == nil {
		related.R = &branchR{
			Inventories: InventorySlice{o},
		}
	} else {
		related.R.Inventories = append(related.R.Inventories, o)
	}

	return nil
}

// SetProduct of the inventory to the related item.
// Sets o.R.Product to related.
// Adds o to related.R.Inventories.
func (o *Inventory) SetProduct(ctx context.Context, exec boil.ContextExecutor, insert bool, related *Product) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"inventory\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"product_id"}),
		strmangle.WhereClause("\"", "\"", 2, inventoryPrimaryKeyColumns),
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
		o.R = &inventoryR{
			Product: related,
		}
	} else {
		o.R.Product = related
	}

	if related.R == nil {
		related.R = &productR{
			Inventories: InventorySlice{o},
		}
	} else {
		related.R.Inventories = append(related.R.Inventories, o)
	}

	return nil
}

// SetSalesRep of the inventory to the related item.
// Sets o.R.SalesRep to related.
// Adds o to related.R.SalesRepInventories.
func (o *Inventory) SetSalesRep(ctx context.Context, exec boil.ContextExecutor, insert bool, related *User) error {
	var err error
	if insert {
		if err = related.Insert(ctx, exec, boil.Infer()); err != nil {
			return errors.Wrap(err, "failed to insert into foreign table")
		}
	}

	updateQuery := fmt.Sprintf(
		"UPDATE \"inventory\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, []string{"sales_rep_id"}),
		strmangle.WhereClause("\"", "\"", 2, inventoryPrimaryKeyColumns),
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
		o.R = &inventoryR{
			SalesRep: related,
		}
	} else {
		o.R.SalesRep = related
	}

	if related.R == nil {
		related.R = &userR{
			SalesRepInventories: InventorySlice{o},
		}
	} else {
		related.R.SalesRepInventories = append(related.R.SalesRepInventories, o)
	}

	return nil
}

// Inventories retrieves all the records using an executor.
func Inventories(mods ...qm.QueryMod) inventoryQuery {
	mods = append(mods, qm.From("\"inventory\""))
	return inventoryQuery{NewQuery(mods...)}
}

// FindInventory retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindInventory(ctx context.Context, exec boil.ContextExecutor, iD string, selectCols ...string) (*Inventory, error) {
	inventoryObj := &Inventory{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"inventory\" where \"id\"=$1", sel,
	)

	q := queries.Raw(query, iD)

	err := q.Bind(ctx, exec, inventoryObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from inventory")
	}

	return inventoryObj, nil
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *Inventory) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no inventory provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(inventoryColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	inventoryInsertCacheMut.RLock()
	cache, cached := inventoryInsertCache[key]
	inventoryInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			inventoryAllColumns,
			inventoryColumnsWithDefault,
			inventoryColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(inventoryType, inventoryMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(inventoryType, inventoryMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"inventory\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"inventory\" %sDEFAULT VALUES%s"
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
		return errors.Wrap(err, "models: unable to insert into inventory")
	}

	if !cached {
		inventoryInsertCacheMut.Lock()
		inventoryInsertCache[key] = cache
		inventoryInsertCacheMut.Unlock()
	}

	return nil
}

// Update uses an executor to update the Inventory.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *Inventory) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	inventoryUpdateCacheMut.RLock()
	cache, cached := inventoryUpdateCache[key]
	inventoryUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			inventoryAllColumns,
			inventoryPrimaryKeyColumns,
		)

		if len(wl) == 0 {
			return 0, errors.New("models: unable to update inventory, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"inventory\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, inventoryPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(inventoryType, inventoryMapping, append(wl, inventoryPrimaryKeyColumns...))
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
		return 0, errors.Wrap(err, "models: unable to update inventory row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for inventory")
	}

	if !cached {
		inventoryUpdateCacheMut.Lock()
		inventoryUpdateCache[key] = cache
		inventoryUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values.
func (q inventoryQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for inventory")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for inventory")
	}

	return rowsAff, nil
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o InventorySlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), inventoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"inventory\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, inventoryPrimaryKeyColumns, len(o)))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in inventory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all inventory")
	}
	return rowsAff, nil
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *Inventory) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no inventory provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(inventoryColumnsWithDefault, o)

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

	inventoryUpsertCacheMut.RLock()
	cache, cached := inventoryUpsertCache[key]
	inventoryUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			inventoryAllColumns,
			inventoryColumnsWithDefault,
			inventoryColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			inventoryAllColumns,
			inventoryPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert inventory, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(inventoryPrimaryKeyColumns))
			copy(conflict, inventoryPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"inventory\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(inventoryType, inventoryMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(inventoryType, inventoryMapping, ret)
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
		return errors.Wrap(err, "models: unable to upsert inventory")
	}

	if !cached {
		inventoryUpsertCacheMut.Lock()
		inventoryUpsertCache[key] = cache
		inventoryUpsertCacheMut.Unlock()
	}

	return nil
}

// Delete deletes a single Inventory record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *Inventory) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no Inventory provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), inventoryPrimaryKeyMapping)
	sql := "DELETE FROM \"inventory\" WHERE \"id\"=$1"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args...)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from inventory")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for inventory")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q inventoryQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no inventoryQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from inventory")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for inventory")
	}

	return rowsAff, nil
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o InventorySlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), inventoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"inventory\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, inventoryPrimaryKeyColumns, len(o))

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, args)
	}
	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from inventory slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for inventory")
	}

	return rowsAff, nil
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *Inventory) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindInventory(ctx, exec, o.ID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *InventorySlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := InventorySlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), inventoryPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"inventory\".* FROM \"inventory\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, inventoryPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in InventorySlice")
	}

	*o = slice

	return nil
}

// InventoryExists checks if the Inventory row exists.
func InventoryExists(ctx context.Context, exec boil.ContextExecutor, iD string) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"inventory\" where \"id\"=$1 limit 1)"

	if boil.IsDebug(ctx) {
		writer := boil.DebugWriterFrom(ctx)
		fmt.Fprintln(writer, sql)
		fmt.Fprintln(writer, iD)
	}
	row := exec.QueryRowContext(ctx, sql, iD)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if inventory exists")
	}

	return exists, nil
}
