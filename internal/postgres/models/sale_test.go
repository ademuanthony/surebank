// Code generated by SQLBoiler 3.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/randomize"
	"github.com/volatiletech/sqlboiler/strmangle"
)

var (
	// Relationships sometimes use the reflection helper queries.Equal/queries.Assign
	// so force a package dependency in case they don't.
	_ = queries.Equal
)

func testSales(t *testing.T) {
	t.Parallel()

	query := Sales()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testSalesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
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

	count, err := Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSalesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Sales().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSalesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := SaleSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testSalesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := SaleExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Sale exists: %s", err)
	}
	if !e {
		t.Errorf("Expected SaleExists to return true, but got false.")
	}
}

func testSalesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	saleFound, err := FindSale(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if saleFound == nil {
		t.Error("want a record, got nil")
	}
}

func testSalesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Sales().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testSalesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Sales().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testSalesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	saleOne := &Sale{}
	saleTwo := &Sale{}
	if err = randomize.Struct(seed, saleOne, saleDBTypes, false, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}
	if err = randomize.Struct(seed, saleTwo, saleDBTypes, false, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = saleOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = saleTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Sales().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testSalesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	saleOne := &Sale{}
	saleTwo := &Sale{}
	if err = randomize.Struct(seed, saleOne, saleDBTypes, false, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}
	if err = randomize.Struct(seed, saleTwo, saleDBTypes, false, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = saleOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = saleTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testSalesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSalesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(saleColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testSaleToManyPayments(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b, c Payment

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, paymentDBTypes, false, paymentColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, paymentDBTypes, false, paymentColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.SaleID = a.ID
	c.SaleID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Payments().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.SaleID == b.SaleID {
			bFound = true
		}
		if v.SaleID == c.SaleID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := SaleSlice{&a}
	if err = a.L.LoadPayments(ctx, tx, false, (*[]*Sale)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Payments); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Payments = nil
	if err = a.L.LoadPayments(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Payments); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testSaleToManySaleItems(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b, c SaleItem

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, saleItemDBTypes, false, saleItemColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, saleItemDBTypes, false, saleItemColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.SaleID = a.ID
	c.SaleID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.SaleItems().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.SaleID == b.SaleID {
			bFound = true
		}
		if v.SaleID == c.SaleID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := SaleSlice{&a}
	if err = a.L.LoadSaleItems(ctx, tx, false, (*[]*Sale)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.SaleItems); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.SaleItems = nil
	if err = a.L.LoadSaleItems(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.SaleItems); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testSaleToManyAddOpPayments(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b, c, d, e Payment

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, false, strmangle.SetComplement(salePrimaryKeyColumns, saleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Payment{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, paymentDBTypes, false, strmangle.SetComplement(paymentPrimaryKeyColumns, paymentColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*Payment{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddPayments(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.SaleID {
			t.Error("foreign key was wrong value", a.ID, first.SaleID)
		}
		if a.ID != second.SaleID {
			t.Error("foreign key was wrong value", a.ID, second.SaleID)
		}

		if first.R.Sale != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Sale != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Payments[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Payments[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Payments().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testSaleToManyAddOpSaleItems(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b, c, d, e SaleItem

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, false, strmangle.SetComplement(salePrimaryKeyColumns, saleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*SaleItem{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, saleItemDBTypes, false, strmangle.SetComplement(saleItemPrimaryKeyColumns, saleItemColumnsWithoutDefault)...); err != nil {
			t.Fatal(err)
		}
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	foreignersSplitByInsertion := [][]*SaleItem{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddSaleItems(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.SaleID {
			t.Error("foreign key was wrong value", a.ID, first.SaleID)
		}
		if a.ID != second.SaleID {
			t.Error("foreign key was wrong value", a.ID, second.SaleID)
		}

		if first.R.Sale != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Sale != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.SaleItems[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.SaleItems[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.SaleItems().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testSaleToOneUserUsingArchivedBy(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Sale
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.ArchivedByID, foreign.ID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.ArchivedBy().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := SaleSlice{&local}
	if err = local.L.LoadArchivedBy(ctx, tx, false, (*[]*Sale)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.ArchivedBy == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.ArchivedBy = nil
	if err = local.L.LoadArchivedBy(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.ArchivedBy == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testSaleToOneBranchUsingBranch(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Sale
	var foreign Branch

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, saleDBTypes, false, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
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

	slice := SaleSlice{&local}
	if err = local.L.LoadBranch(ctx, tx, false, (*[]*Sale)(&slice), nil); err != nil {
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

func testSaleToOneUserUsingCreatedBy(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Sale
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, saleDBTypes, false, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.CreatedByID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.CreatedBy().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := SaleSlice{&local}
	if err = local.L.LoadCreatedBy(ctx, tx, false, (*[]*Sale)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.CreatedBy == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.CreatedBy = nil
	if err = local.L.LoadCreatedBy(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.CreatedBy == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testSaleToOneUserUsingUpdatedBy(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Sale
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize User struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	queries.Assign(&local.UpdatedByID, foreign.ID)
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.UpdatedBy().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if !queries.Equal(check.ID, foreign.ID) {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := SaleSlice{&local}
	if err = local.L.LoadUpdatedBy(ctx, tx, false, (*[]*Sale)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.UpdatedBy == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.UpdatedBy = nil
	if err = local.L.LoadUpdatedBy(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.UpdatedBy == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testSaleToOneSetOpUserUsingArchivedBy(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, false, strmangle.SetComplement(salePrimaryKeyColumns, saleColumnsWithoutDefault)...); err != nil {
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
		err = a.SetArchivedBy(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.ArchivedBy != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ArchivedBySales[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.ArchivedByID, x.ID) {
			t.Error("foreign key was wrong value", a.ArchivedByID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.ArchivedByID))
		reflect.Indirect(reflect.ValueOf(&a.ArchivedByID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.ArchivedByID, x.ID) {
			t.Error("foreign key was wrong value", a.ArchivedByID, x.ID)
		}
	}
}

func testSaleToOneRemoveOpUserUsingArchivedBy(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, false, strmangle.SetComplement(salePrimaryKeyColumns, saleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetArchivedBy(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveArchivedBy(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.ArchivedBy().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.ArchivedBy != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.ArchivedByID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.ArchivedBySales) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testSaleToOneSetOpBranchUsingBranch(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b, c Branch

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, false, strmangle.SetComplement(salePrimaryKeyColumns, saleColumnsWithoutDefault)...); err != nil {
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

		if x.R.Sales[0] != &a {
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
func testSaleToOneSetOpUserUsingCreatedBy(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, false, strmangle.SetComplement(salePrimaryKeyColumns, saleColumnsWithoutDefault)...); err != nil {
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
		err = a.SetCreatedBy(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.CreatedBy != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.CreatedBySales[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.CreatedByID != x.ID {
			t.Error("foreign key was wrong value", a.CreatedByID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.CreatedByID))
		reflect.Indirect(reflect.ValueOf(&a.CreatedByID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.CreatedByID != x.ID {
			t.Error("foreign key was wrong value", a.CreatedByID, x.ID)
		}
	}
}
func testSaleToOneSetOpUserUsingUpdatedBy(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, false, strmangle.SetComplement(salePrimaryKeyColumns, saleColumnsWithoutDefault)...); err != nil {
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
		err = a.SetUpdatedBy(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.UpdatedBy != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.UpdatedBySales[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if !queries.Equal(a.UpdatedByID, x.ID) {
			t.Error("foreign key was wrong value", a.UpdatedByID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.UpdatedByID))
		reflect.Indirect(reflect.ValueOf(&a.UpdatedByID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if !queries.Equal(a.UpdatedByID, x.ID) {
			t.Error("foreign key was wrong value", a.UpdatedByID, x.ID)
		}
	}
}

func testSaleToOneRemoveOpUserUsingUpdatedBy(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Sale
	var b User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, saleDBTypes, false, strmangle.SetComplement(salePrimaryKeyColumns, saleColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err = a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = a.SetUpdatedBy(ctx, tx, true, &b); err != nil {
		t.Fatal(err)
	}

	if err = a.RemoveUpdatedBy(ctx, tx, &b); err != nil {
		t.Error("failed to remove relationship")
	}

	count, err := a.UpdatedBy().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 0 {
		t.Error("want no relationships remaining")
	}

	if a.R.UpdatedBy != nil {
		t.Error("R struct entry should be nil")
	}

	if !queries.IsValuerNil(a.UpdatedByID) {
		t.Error("foreign key value should be nil")
	}

	if len(b.R.UpdatedBySales) != 0 {
		t.Error("failed to remove a from b's relationships")
	}
}

func testSalesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
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

func testSalesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := SaleSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testSalesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Sales().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	saleDBTypes = map[string]string{`ID`: `character`, `ReceiptNumber`: `character varying`, `Amount`: `double precision`, `AmountTender`: `double precision`, `Balance`: `double precision`, `CustomerName`: `character varying`, `PhoneNumber`: `character varying`, `CreatedAt`: `timestamp with time zone`, `UpdatedAt`: `timestamp with time zone`, `ArchivedAt`: `timestamp with time zone`, `CreatedByID`: `character`, `UpdatedByID`: `character`, `ArchivedByID`: `character`, `BranchID`: `character`}
	_           = bytes.MinRead
)

func testSalesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(salePrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(saleAllColumns) == len(salePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, saleDBTypes, true, salePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testSalesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(saleAllColumns) == len(salePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Sale{}
	if err = randomize.Struct(seed, o, saleDBTypes, true, saleColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, saleDBTypes, true, salePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(saleAllColumns, salePrimaryKeyColumns) {
		fields = saleAllColumns
	} else {
		fields = strmangle.SetComplement(
			saleAllColumns,
			salePrimaryKeyColumns,
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

	slice := SaleSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testSalesUpsert(t *testing.T) {
	t.Parallel()

	if len(saleAllColumns) == len(salePrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Sale{}
	if err = randomize.Struct(seed, &o, saleDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Sale: %s", err)
	}

	count, err := Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, saleDBTypes, false, salePrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Sale struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Sale: %s", err)
	}

	count, err = Sales().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}