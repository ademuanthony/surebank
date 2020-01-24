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

func testBranches(t *testing.T) {
	t.Parallel()

	query := Branches()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testBranchesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
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

	count, err := Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBranchesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Branches().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBranchesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BranchSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testBranchesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := BranchExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Branch exists: %s", err)
	}
	if !e {
		t.Errorf("Expected BranchExists to return true, but got false.")
	}
}

func testBranchesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	branchFound, err := FindBranch(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if branchFound == nil {
		t.Error("want a record, got nil")
	}
}

func testBranchesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Branches().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testBranchesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Branches().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testBranchesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	branchOne := &Branch{}
	branchTwo := &Branch{}
	if err = randomize.Struct(seed, branchOne, branchDBTypes, false, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}
	if err = randomize.Struct(seed, branchTwo, branchDBTypes, false, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = branchOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = branchTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Branches().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testBranchesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	branchOne := &Branch{}
	branchTwo := &Branch{}
	if err = randomize.Struct(seed, branchOne, branchDBTypes, false, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}
	if err = randomize.Struct(seed, branchTwo, branchDBTypes, false, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = branchOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = branchTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testBranchesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBranchesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(branchColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testBranchToManyAccounts(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c Account

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.BranchID = a.ID
	c.BranchID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Accounts().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.BranchID == b.BranchID {
			bFound = true
		}
		if v.BranchID == c.BranchID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := BranchSlice{&a}
	if err = a.L.LoadAccounts(ctx, tx, false, (*[]*Branch)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Accounts); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Accounts = nil
	if err = a.L.LoadAccounts(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Accounts); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testBranchToManyCustomers(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c Customer

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, customerDBTypes, false, customerColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, customerDBTypes, false, customerColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.BranchID = a.ID
	c.BranchID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Customers().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.BranchID == b.BranchID {
			bFound = true
		}
		if v.BranchID == c.BranchID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := BranchSlice{&a}
	if err = a.L.LoadCustomers(ctx, tx, false, (*[]*Branch)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Customers); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Customers = nil
	if err = a.L.LoadCustomers(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Customers); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testBranchToManySales(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c Sale

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, saleDBTypes, false, saleColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, saleDBTypes, false, saleColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.BranchID = a.ID
	c.BranchID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Sales().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.BranchID == b.BranchID {
			bFound = true
		}
		if v.BranchID == c.BranchID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := BranchSlice{&a}
	if err = a.L.LoadSales(ctx, tx, false, (*[]*Branch)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Sales); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Sales = nil
	if err = a.L.LoadSales(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Sales); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testBranchToManyStocks(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c Stock

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, stockDBTypes, false, stockColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, stockDBTypes, false, stockColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.BranchID = a.ID
	c.BranchID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Stocks().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.BranchID == b.BranchID {
			bFound = true
		}
		if v.BranchID == c.BranchID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := BranchSlice{&a}
	if err = a.L.LoadStocks(ctx, tx, false, (*[]*Branch)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Stocks); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Stocks = nil
	if err = a.L.LoadStocks(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Stocks); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testBranchToManyUsers(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, userDBTypes, false, userColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.BranchID = a.ID
	c.BranchID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Users().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.BranchID == b.BranchID {
			bFound = true
		}
		if v.BranchID == c.BranchID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := BranchSlice{&a}
	if err = a.L.LoadUsers(ctx, tx, false, (*[]*Branch)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Users); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Users = nil
	if err = a.L.LoadUsers(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Users); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testBranchToManyAddOpAccounts(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c, d, e Account

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, false, strmangle.SetComplement(branchPrimaryKeyColumns, branchColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Account{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, accountDBTypes, false, strmangle.SetComplement(accountPrimaryKeyColumns, accountColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Account{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddAccounts(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.BranchID {
			t.Error("foreign key was wrong value", a.ID, first.BranchID)
		}
		if a.ID != second.BranchID {
			t.Error("foreign key was wrong value", a.ID, second.BranchID)
		}

		if first.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Accounts[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Accounts[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Accounts().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testBranchToManyAddOpCustomers(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c, d, e Customer

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, false, strmangle.SetComplement(branchPrimaryKeyColumns, branchColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Customer{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, customerDBTypes, false, strmangle.SetComplement(customerPrimaryKeyColumns, customerColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Customer{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddCustomers(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.BranchID {
			t.Error("foreign key was wrong value", a.ID, first.BranchID)
		}
		if a.ID != second.BranchID {
			t.Error("foreign key was wrong value", a.ID, second.BranchID)
		}

		if first.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Customers[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Customers[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Customers().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testBranchToManyAddOpSales(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c, d, e Sale

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, false, strmangle.SetComplement(branchPrimaryKeyColumns, branchColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Sale{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, saleDBTypes, false, strmangle.SetComplement(salePrimaryKeyColumns, saleColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Sale{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddSales(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.BranchID {
			t.Error("foreign key was wrong value", a.ID, first.BranchID)
		}
		if a.ID != second.BranchID {
			t.Error("foreign key was wrong value", a.ID, second.BranchID)
		}

		if first.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Sales[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Sales[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Sales().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testBranchToManyAddOpStocks(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c, d, e Stock

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, false, strmangle.SetComplement(branchPrimaryKeyColumns, branchColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Stock{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, stockDBTypes, false, strmangle.SetComplement(stockPrimaryKeyColumns, stockColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Stock{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddStocks(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.BranchID {
			t.Error("foreign key was wrong value", a.ID, first.BranchID)
		}
		if a.ID != second.BranchID {
			t.Error("foreign key was wrong value", a.ID, second.BranchID)
		}

		if first.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Stocks[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Stocks[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Stocks().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testBranchToManyAddOpUsers(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Branch
	var b, c, d, e User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, branchDBTypes, false, strmangle.SetComplement(branchPrimaryKeyColumns, branchColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*User{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, userDBTypes, false, strmangle.SetComplement(userPrimaryKeyColumns, userColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*User{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddUsers(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.BranchID {
			t.Error("foreign key was wrong value", a.ID, first.BranchID)
		}
		if a.ID != second.BranchID {
			t.Error("foreign key was wrong value", a.ID, second.BranchID)
		}

		if first.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Branch != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Users[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Users[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Users().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}

func testBranchesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
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

func testBranchesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := BranchSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testBranchesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Branches().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	branchDBTypes = map[string]string{`ID`: `character`, `Name`: `character varying`, `CreatedAt`: `timestamp with time zone`, `UpdatedAt`: `timestamp with time zone`, `ArchivedAt`: `timestamp with time zone`}
	_             = bytes.MinRead
)

func testBranchesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(branchPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(branchAllColumns) == len(branchPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, branchDBTypes, true, branchPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testBranchesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(branchAllColumns) == len(branchPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Branch{}
	if err = randomize.Struct(seed, o, branchDBTypes, true, branchColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, branchDBTypes, true, branchPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(branchAllColumns, branchPrimaryKeyColumns) {
		fields = branchAllColumns
	} else {
		fields = strmangle.SetComplement(
			branchAllColumns,
			branchPrimaryKeyColumns,
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

	slice := BranchSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testBranchesUpsert(t *testing.T) {
	t.Parallel()

	if len(branchAllColumns) == len(branchPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Branch{}
	if err = randomize.Struct(seed, &o, branchDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Branch: %s", err)
	}

	count, err := Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, branchDBTypes, false, branchPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Branch struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Branch: %s", err)
	}

	count, err = Branches().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
