// Code generated by SQLBoiler 4.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

func testRepsExpenses(t *testing.T) {
	t.Parallel()

	query := RepsExpenses()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testRepsExpensesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
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

	count, err := RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRepsExpensesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := RepsExpenses().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRepsExpensesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := RepsExpenseSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testRepsExpensesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := RepsExpenseExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if RepsExpense exists: %s", err)
	}
	if !e {
		t.Errorf("Expected RepsExpenseExists to return true, but got false.")
	}
}

func testRepsExpensesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	repsExpenseFound, err := FindRepsExpense(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if repsExpenseFound == nil {
		t.Error("want a record, got nil")
	}
}

func testRepsExpensesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = RepsExpenses().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testRepsExpensesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := RepsExpenses().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testRepsExpensesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	repsExpenseOne := &RepsExpense{}
	repsExpenseTwo := &RepsExpense{}
	if err = randomize.Struct(seed, repsExpenseOne, repsExpenseDBTypes, false, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}
	if err = randomize.Struct(seed, repsExpenseTwo, repsExpenseDBTypes, false, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = repsExpenseOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = repsExpenseTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := RepsExpenses().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testRepsExpensesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	repsExpenseOne := &RepsExpense{}
	repsExpenseTwo := &RepsExpense{}
	if err = randomize.Struct(seed, repsExpenseOne, repsExpenseDBTypes, false, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}
	if err = randomize.Struct(seed, repsExpenseTwo, repsExpenseDBTypes, false, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = repsExpenseOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = repsExpenseTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testRepsExpensesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRepsExpensesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(repsExpenseColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testRepsExpenseToOneUserUsingSalesRep(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local RepsExpense
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, repsExpenseDBTypes, false, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
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

	slice := RepsExpenseSlice{&local}
	if err = local.L.LoadSalesRep(ctx, tx, false, (*[]*RepsExpense)(&slice), nil); err != nil {
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

func testRepsExpenseToOneSetOpUserUsingSalesRep(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a RepsExpense
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, repsExpenseDBTypes, false, strmangle.SetComplement(repsExpensePrimaryKeyColumns, repsExpenseColumnsWithoutDefault)...); err != nil {
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

		if x.R.SalesRepRepsExpenses[0] != &a {
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

func testRepsExpensesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
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

func testRepsExpensesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := RepsExpenseSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testRepsExpensesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := RepsExpenses().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	repsExpenseDBTypes = map[string]string{`ID`: `character`, `SalesRepID`: `character`, `Amount`: `double precision`, `Reason`: `character varying`, `Date`: `bigint`}
	_                  = bytes.MinRead
)

func testRepsExpensesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(repsExpensePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(repsExpenseAllColumns) == len(repsExpensePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpensePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testRepsExpensesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(repsExpenseAllColumns) == len(repsExpensePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &RepsExpense{}
	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpenseColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, repsExpenseDBTypes, true, repsExpensePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(repsExpenseAllColumns, repsExpensePrimaryKeyColumns) {
		fields = repsExpenseAllColumns
	} else {
		fields = strmangle.SetComplement(
			repsExpenseAllColumns,
			repsExpensePrimaryKeyColumns,
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

	slice := RepsExpenseSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testRepsExpensesUpsert(t *testing.T) {
	t.Parallel()

	if len(repsExpenseAllColumns) == len(repsExpensePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := RepsExpense{}
	if err = randomize.Struct(seed, &o, repsExpenseDBTypes, true); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert RepsExpense: %s", err)
	}

	count, err := RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, repsExpenseDBTypes, false, repsExpensePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize RepsExpense struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert RepsExpense: %s", err)
	}

	count, err = RepsExpenses().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
