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

func testExpenditures(t *testing.T) {
	t.Parallel()

	query := Expenditures()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testExpendituresDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
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

	count, err := Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testExpendituresQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Expenditures().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testExpendituresSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ExpenditureSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testExpendituresExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ExpenditureExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Expenditure exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ExpenditureExists to return true, but got false.")
	}
}

func testExpendituresFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	expenditureFound, err := FindExpenditure(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if expenditureFound == nil {
		t.Error("want a record, got nil")
	}
}

func testExpendituresBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Expenditures().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testExpendituresOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Expenditures().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testExpendituresAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	expenditureOne := &Expenditure{}
	expenditureTwo := &Expenditure{}
	if err = randomize.Struct(seed, expenditureOne, expenditureDBTypes, false, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}
	if err = randomize.Struct(seed, expenditureTwo, expenditureDBTypes, false, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = expenditureOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = expenditureTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Expenditures().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testExpendituresCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	expenditureOne := &Expenditure{}
	expenditureTwo := &Expenditure{}
	if err = randomize.Struct(seed, expenditureOne, expenditureDBTypes, false, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}
	if err = randomize.Struct(seed, expenditureTwo, expenditureDBTypes, false, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = expenditureOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = expenditureTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testExpendituresInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testExpendituresInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(expenditureColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testExpendituresReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
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

func testExpendituresReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ExpenditureSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testExpendituresSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Expenditures().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	expenditureDBTypes = map[string]string{`ID`: `character`, `Amount`: `double precision`, `Date`: `bigint`}
	_                  = bytes.MinRead
)

func testExpendituresUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(expenditurePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(expenditureAllColumns) == len(expenditurePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditurePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testExpendituresSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(expenditureAllColumns) == len(expenditurePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Expenditure{}
	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditureColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, expenditureDBTypes, true, expenditurePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(expenditureAllColumns, expenditurePrimaryKeyColumns) {
		fields = expenditureAllColumns
	} else {
		fields = strmangle.SetComplement(
			expenditureAllColumns,
			expenditurePrimaryKeyColumns,
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

	slice := ExpenditureSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testExpendituresUpsert(t *testing.T) {
	t.Parallel()

	if len(expenditureAllColumns) == len(expenditurePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Expenditure{}
	if err = randomize.Struct(seed, &o, expenditureDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Expenditure: %s", err)
	}

	count, err := Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, expenditureDBTypes, false, expenditurePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Expenditure struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Expenditure: %s", err)
	}

	count, err = Expenditures().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}