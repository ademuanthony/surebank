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

func testCustomers(t *testing.T) {
	t.Parallel()

	query := Customers()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testCustomersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
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

	count, err := Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCustomersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Customers().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCustomersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := CustomerSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testCustomersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := CustomerExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Customer exists: %s", err)
	}
	if !e {
		t.Errorf("Expected CustomerExists to return true, but got false.")
	}
}

func testCustomersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	customerFound, err := FindCustomer(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if customerFound == nil {
		t.Error("want a record, got nil")
	}
}

func testCustomersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Customers().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testCustomersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Customers().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testCustomersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	customerOne := &Customer{}
	customerTwo := &Customer{}
	if err = randomize.Struct(seed, customerOne, customerDBTypes, false, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}
	if err = randomize.Struct(seed, customerTwo, customerDBTypes, false, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = customerOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = customerTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Customers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testCustomersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	customerOne := &Customer{}
	customerTwo := &Customer{}
	if err = randomize.Struct(seed, customerOne, customerDBTypes, false, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}
	if err = randomize.Struct(seed, customerTwo, customerDBTypes, false, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = customerOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = customerTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testCustomersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCustomersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(customerColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testCustomerToManyAccounts(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Customer
	var b, c Account

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
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

	b.CustomerID = a.ID
	c.CustomerID = a.ID

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
		if v.CustomerID == b.CustomerID {
			bFound = true
		}
		if v.CustomerID == c.CustomerID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := CustomerSlice{&a}
	if err = a.L.LoadAccounts(ctx, tx, false, (*[]*Customer)(&slice), nil); err != nil {
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

func testCustomerToManyAddOpAccounts(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Customer
	var b, c, d, e Account

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, customerDBTypes, false, strmangle.SetComplement(customerPrimaryKeyColumns, customerColumnsWithoutDefault)...); err != nil {
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

		if a.ID != first.CustomerID {
			t.Error("foreign key was wrong value", a.ID, first.CustomerID)
		}
		if a.ID != second.CustomerID {
			t.Error("foreign key was wrong value", a.ID, second.CustomerID)
		}

		if first.R.Customer != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Customer != &a {
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
func testCustomerToOneBranchUsingBranch(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Customer
	var foreign Branch

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, customerDBTypes, false, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
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

	slice := CustomerSlice{&local}
	if err = local.L.LoadBranch(ctx, tx, false, (*[]*Customer)(&slice), nil); err != nil {
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

func testCustomerToOneUserUsingSalesRep(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Customer
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, customerDBTypes, false, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
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

	slice := CustomerSlice{&local}
	if err = local.L.LoadSalesRep(ctx, tx, false, (*[]*Customer)(&slice), nil); err != nil {
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

func testCustomerToOneSetOpBranchUsingBranch(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Customer
	var b, c Branch

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, customerDBTypes, false, strmangle.SetComplement(customerPrimaryKeyColumns, customerColumnsWithoutDefault)...); err != nil {
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

		if x.R.Customers[0] != &a {
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
func testCustomerToOneSetOpUserUsingSalesRep(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Customer
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, customerDBTypes, false, strmangle.SetComplement(customerPrimaryKeyColumns, customerColumnsWithoutDefault)...); err != nil {
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

		if x.R.SalesRepCustomers[0] != &a {
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

func testCustomersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
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

func testCustomersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := CustomerSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testCustomersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Customers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	customerDBTypes = map[string]string{`ID`: `character`, `Email`: `character varying`, `Name`: `character varying`, `PhoneNumber`: `character varying`, `Address`: `character varying`, `SalesRepID`: `character`, `CreatedAt`: `timestamp with time zone`, `ArchivedAt`: `timestamp with time zone`, `BranchID`: `character`, `UpdatedAt`: `timestamp with time zone`}
	_               = bytes.MinRead
)

func testCustomersUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(customerPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(customerAllColumns) == len(customerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, customerDBTypes, true, customerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testCustomersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(customerAllColumns) == len(customerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Customer{}
	if err = randomize.Struct(seed, o, customerDBTypes, true, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, customerDBTypes, true, customerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(customerAllColumns, customerPrimaryKeyColumns) {
		fields = customerAllColumns
	} else {
		fields = strmangle.SetComplement(
			customerAllColumns,
			customerPrimaryKeyColumns,
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

	slice := CustomerSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testCustomersUpsert(t *testing.T) {
	t.Parallel()

	if len(customerAllColumns) == len(customerPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Customer{}
	if err = randomize.Struct(seed, &o, customerDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Customer: %s", err)
	}

	count, err := Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, customerDBTypes, false, customerPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Customer: %s", err)
	}

	count, err = Customers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
