// Code generated by SQLBoiler 4.1.1 (https://github.com/volatiletech/sqlboiler/v4). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/randomize"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries"
	"github.com/volatiletech/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testInventories(t *testing.T) {
	t.Parallel()

	query := Inventories()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testInventoriesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := o.Delete(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testInventoriesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Inventories().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testInventoriesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := InventorySlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testInventoriesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := InventoryExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Inventory exists: %s", err)
	}
	if !e {
		t.Errorf("Expected InventoryExists to return true, but got false.")
	}
}

func testInventoriesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	inventoryFound, err := FindInventory(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if inventoryFound == nil {
		t.Error("want a record, got nil")
	}
}

func testInventoriesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Inventories().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testInventoriesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Inventories().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testInventoriesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	inventoryOne := &Inventory{}
	inventoryTwo := &Inventory{}
	if err = randomize.Struct(seed, inventoryOne, inventoryDBTypes, false, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}
	if err = randomize.Struct(seed, inventoryTwo, inventoryDBTypes, false, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = inventoryOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = inventoryTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Inventories().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testInventoriesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	inventoryOne := &Inventory{}
	inventoryTwo := &Inventory{}
	if err = randomize.Struct(seed, inventoryOne, inventoryDBTypes, false, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}
	if err = randomize.Struct(seed, inventoryTwo, inventoryDBTypes, false, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = inventoryOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = inventoryTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testInventoriesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testInventoriesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(inventoryColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testInventoryToOneBranchUsingBranch(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Inventory
	var foreign Branch

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, inventoryDBTypes, false, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, branchDBTypes, false, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.BranchID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Branch().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := InventorySlice{&local}
	if err = local.L.LoadBranch(ctx, tx, false, (*[]*Inventory)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Branch == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Branch = nil
	if err = local.L.LoadBranch(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Branch == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testInventoryToOneProductUsingProduct(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Inventory
	var foreign Product

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, inventoryDBTypes, false, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, productDBTypes, false, productColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Product struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.ProductID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Product().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := InventorySlice{&local}
	if err = local.L.LoadProduct(ctx, tx, false, (*[]*Inventory)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Product == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Product = nil
	if err = local.L.LoadProduct(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Product == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testInventoryToOneUserUsingSalesRep(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Inventory
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, inventoryDBTypes, false, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.SalesRepID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.SalesRep().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := InventorySlice{&local}
	if err = local.L.LoadSalesRep(ctx, tx, false, (*[]*Inventory)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.SalesRep == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.SalesRep = nil
	if err = local.L.LoadSalesRep(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.SalesRep == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testInventoryToOneSetOpBranchUsingBranch(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Inventory
	var b, c Branch

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, inventoryDBTypes, false, strmangle.SetComplement(inventoryPrimaryKeyColumns, inventoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, branchDBTypes, false, strmangle.SetComplement(branchPrimaryKeyColumns, branchColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, branchDBTypes, false, strmangle.SetComplement(branchPrimaryKeyColumns, branchColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Branch{&b, &c} {
		err = a.SetBranch(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Branch != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Inventories[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.BranchID != x.ID {
			t.Error("foreign key was wrong value", a.BranchID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.BranchID))
		reflect.Indirect(reflect.ValueOf(&a.BranchID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.BranchID != x.ID {
			t.Error("foreign key was wrong value", a.BranchID, x.ID)
		}
	}
}
func testInventoryToOneSetOpProductUsingProduct(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Inventory
	var b, c Product

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, inventoryDBTypes, false, strmangle.SetComplement(inventoryPrimaryKeyColumns, inventoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, productDBTypes, false, strmangle.SetComplement(productPrimaryKeyColumns, productColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, productDBTypes, false, strmangle.SetComplement(productPrimaryKeyColumns, productColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Product{&b, &c} {
		err = a.SetProduct(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Product != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Inventories[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.ProductID != x.ID {
			t.Error("foreign key was wrong value", a.ProductID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ProductID))
		reflect.Indirect(reflect.ValueOf(&a.ProductID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.ProductID != x.ID {
			t.Error("foreign key was wrong value", a.ProductID, x.ID)
		}
	}
}
func testInventoryToOneSetOpUserUsingSalesRep(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Inventory
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, inventoryDBTypes, false, strmangle.SetComplement(inventoryPrimaryKeyColumns, inventoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*User{&b, &c} {
		err = a.SetSalesRep(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.SalesRep != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.SalesRepInventories[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.SalesRepID != x.ID {
			t.Error("foreign key was wrong value", a.SalesRepID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.SalesRepID))
		reflect.Indirect(reflect.ValueOf(&a.SalesRepID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.SalesRepID != x.ID {
			t.Error("foreign key was wrong value", a.SalesRepID, x.ID)
		}
	}
}

func testInventoriesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = o.Reload(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testInventoriesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := InventorySlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testInventoriesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Inventories().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	inventoryDBTypes = map[string]string{`ID`: `character`, `ProductID`: `character`, `BranchID`: `character`, `TXType`: `character varying`, `OpeningBalance`: `double precision`, `Quantity`: `double precision`, `Narration`: `character varying`, `SalesRepID`: `character`, `CreatedAt`: `bigint`, `UpdatedAt`: `bigint`, `ArchivedAt`: `bigint`}
	_                = bytes.MinRead
)

func testInventoriesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(inventoryPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(inventoryAllColumns) == len(inventoryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testInventoriesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(inventoryAllColumns) == len(inventoryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Inventory{}
	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, inventoryDBTypes, true, inventoryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(inventoryAllColumns, inventoryPrimaryKeyColumns) {
		fields = inventoryAllColumns
	} else {
		fields = strmangle.SetComplement(
			inventoryAllColumns,
			inventoryPrimaryKeyColumns,
		)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	typ := reflect.TypeOf(o).Elem()
	n := typ.NumField()

	updateMap := M{}
	for _, col := range fields {
		for i := 0; i < n; i++ {
			f := typ.Field(i)
			if f.Tag.Get("boil") == col {
				updateMap[col] = value.Field(i).Interface()
			}
		}
	}

	slice := InventorySlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testInventoriesUpsert(t *testing.T) {
	t.Parallel()

	if len(inventoryAllColumns) == len(inventoryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Inventory{}
	if err = randomize.Struct(seed, &o, inventoryDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Inventory: %s", err)
	}

	count, err := Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, inventoryDBTypes, false, inventoryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Inventory struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Inventory: %s", err)
	}

	count, err = Inventories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
