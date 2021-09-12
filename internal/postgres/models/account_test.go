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

func testAccounts(t *testing.T) {
	t.Parallel()

	query := Accounts()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testAccountsDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
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

	count, err := Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAccountsQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := Accounts().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAccountsSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AccountSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testAccountsExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := AccountExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if Account exists: %s", err)
	}
	if !e {
		t.Errorf("Expected AccountExists to return true, but got false.")
	}
}

func testAccountsFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	accountFound, err := FindAccount(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if accountFound == nil {
		t.Error("want a record, got nil")
	}
}

func testAccountsBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = Accounts().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testAccountsOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := Accounts().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testAccountsAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	accountOne := &Account{}
	accountTwo := &Account{}
	if err = randomize.Struct(seed, accountOne, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}
	if err = randomize.Struct(seed, accountTwo, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = accountOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = accountTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Accounts().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testAccountsCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	accountOne := &Account{}
	accountTwo := &Account{}
	if err = randomize.Struct(seed, accountOne, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}
	if err = randomize.Struct(seed, accountTwo, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = accountOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = accountTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testAccountsInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAccountsInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(accountColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testAccountToManyDSCommissions(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Account
	var b, c DSCommission

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, dsCommissionDBTypes, false, dsCommissionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, dsCommissionDBTypes, false, dsCommissionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.AccountID = a.ID
	c.AccountID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.DSCommissions().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.AccountID == b.AccountID {
			bFound = true
		}
		if v.AccountID == c.AccountID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := AccountSlice{&a}
	if err = a.L.LoadDSCommissions(ctx, tx, false, (*[]*Account)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DSCommissions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.DSCommissions = nil
	if err = a.L.LoadDSCommissions(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.DSCommissions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testAccountToManyTransactions(t *testing.T) {
	var err error
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Account
	var b, c Transaction

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	if err = randomize.Struct(seed, &b, transactionDBTypes, false, transactionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, transactionDBTypes, false, transactionColumnsWithDefault...); err != nil {
		t.Fatal(err)
	}

	b.AccountID = a.ID
	c.AccountID = a.ID

	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = c.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := a.Transactions().All(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	bFound, cFound := false, false
	for _, v := range check {
		if v.AccountID == b.AccountID {
			bFound = true
		}
		if v.AccountID == c.AccountID {
			cFound = true
		}
	}

	if !bFound {
		t.Error("expected to find b")
	}
	if !cFound {
		t.Error("expected to find c")
	}

	slice := AccountSlice{&a}
	if err = a.L.LoadTransactions(ctx, tx, false, (*[]*Account)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Transactions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	a.R.Transactions = nil
	if err = a.L.LoadTransactions(ctx, tx, true, &a, nil); err != nil {
		t.Fatal(err)
	}
	if got := len(a.R.Transactions); got != 2 {
		t.Error("number of eager loaded records wrong, got:", got)
	}

	if t.Failed() {
		t.Logf("%#v", check)
	}
}

func testAccountToManyAddOpDSCommissions(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Account
	var b, c, d, e DSCommission

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, accountDBTypes, false, strmangle.SetComplement(accountPrimaryKeyColumns, accountColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*DSCommission{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, dsCommissionDBTypes, false, strmangle.SetComplement(dsCommissionPrimaryKeyColumns, dsCommissionColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*DSCommission{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddDSCommissions(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.AccountID {
			t.Error("foreign key was wrong value", a.ID, first.AccountID)
		}
		if a.ID != second.AccountID {
			t.Error("foreign key was wrong value", a.ID, second.AccountID)
		}

		if first.R.Account != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Account != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.DSCommissions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.DSCommissions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.DSCommissions().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testAccountToManyAddOpTransactions(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Account
	var b, c, d, e Transaction

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, accountDBTypes, false, strmangle.SetComplement(accountPrimaryKeyColumns, accountColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	foreigners := []*Transaction{&b, &c, &d, &e}
	for _, x := range foreigners {
		if err = randomize.Struct(seed, x, transactionDBTypes, false, strmangle.SetComplement(transactionPrimaryKeyColumns, transactionColumnsWithoutDefault)...); err != nil {
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

	foreignersSplitByInsertion := [][]*Transaction{
		{&b, &c},
		{&d, &e},
	}

	for i, x := range foreignersSplitByInsertion {
		err = a.AddTransactions(ctx, tx, i != 0, x...)
		if err != nil {
			t.Fatal(err)
		}

		first := x[0]
		second := x[1]

		if a.ID != first.AccountID {
			t.Error("foreign key was wrong value", a.ID, first.AccountID)
		}
		if a.ID != second.AccountID {
			t.Error("foreign key was wrong value", a.ID, second.AccountID)
		}

		if first.R.Account != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}
		if second.R.Account != &a {
			t.Error("relationship was not added properly to the foreign slice")
		}

		if a.R.Transactions[i*2] != first {
			t.Error("relationship struct slice not set to correct value")
		}
		if a.R.Transactions[i*2+1] != second {
			t.Error("relationship struct slice not set to correct value")
		}

		count, err := a.Transactions().Count(ctx, tx)
		if err != nil {
			t.Fatal(err)
		}
		if want := int64((i + 1) * 2); count != want {
			t.Error("want", want, "got", count)
		}
	}
}
func testAccountToOneBranchUsingBranch(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Account
	var foreign Branch

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
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

	slice := AccountSlice{&local}
	if err = local.L.LoadBranch(ctx, tx, false, (*[]*Account)(&slice), nil); err != nil {
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

func testAccountToOneCustomerUsingCustomer(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Account
	var foreign Customer

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, customerDBTypes, false, customerColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Customer struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.CustomerID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Customer().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := AccountSlice{&local}
	if err = local.L.LoadCustomer(ctx, tx, false, (*[]*Account)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Customer == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Customer = nil
	if err = local.L.LoadCustomer(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Customer == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testAccountToOneUserUsingSalesRep(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local Account
	var foreign User

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, accountDBTypes, false, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
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

	slice := AccountSlice{&local}
	if err = local.L.LoadSalesRep(ctx, tx, false, (*[]*Account)(&slice), nil); err != nil {
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

func testAccountToOneSetOpBranchUsingBranch(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Account
	var b, c Branch

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, accountDBTypes, false, strmangle.SetComplement(accountPrimaryKeyColumns, accountColumnsWithoutDefault)...); err != nil {
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

		if x.R.Accounts[0] != &a {
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
func testAccountToOneSetOpCustomerUsingCustomer(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Account
	var b, c Customer

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, accountDBTypes, false, strmangle.SetComplement(accountPrimaryKeyColumns, accountColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, customerDBTypes, false, strmangle.SetComplement(customerPrimaryKeyColumns, customerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, customerDBTypes, false, strmangle.SetComplement(customerPrimaryKeyColumns, customerColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Customer{&b, &c} {
		err = a.SetCustomer(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Customer != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.Accounts[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.CustomerID != x.ID {
			t.Error("foreign key was wrong value", a.CustomerID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.CustomerID))
		reflect.Indirect(reflect.ValueOf(&a.CustomerID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.CustomerID != x.ID {
			t.Error("foreign key was wrong value", a.CustomerID, x.ID)
		}
	}
}
func testAccountToOneSetOpUserUsingSalesRep(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a Account
	var b, c User

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, accountDBTypes, false, strmangle.SetComplement(accountPrimaryKeyColumns, accountColumnsWithoutDefault)...); err != nil {
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

		if x.R.SalesRepAccounts[0] != &a {
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

func testAccountsReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
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

func testAccountsReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := AccountSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testAccountsSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := Accounts().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	accountDBTypes = map[string]string{`ID`: `character`, `BranchID`: `character`, `Number`: `character varying`, `CustomerID`: `character`, `AccountType`: `character varying`, `Target`: `double precision`, `TargetInfo`: `character varying`, `SalesRepID`: `character`, `CreatedAt`: `bigint`, `UpdatedAt`: `bigint`, `ArchivedAt`: `bigint`, `Balance`: `double precision`, `LastPaymentDate`: `bigint`}
	_              = bytes.MinRead
)

func testAccountsUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(accountPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(accountAllColumns) == len(accountPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, accountDBTypes, true, accountPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testAccountsSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(accountAllColumns) == len(accountPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &Account{}
	if err = randomize.Struct(seed, o, accountDBTypes, true, accountColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, accountDBTypes, true, accountPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(accountAllColumns, accountPrimaryKeyColumns) {
		fields = accountAllColumns
	} else {
		fields = strmangle.SetComplement(
			accountAllColumns,
			accountPrimaryKeyColumns,
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

	slice := AccountSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testAccountsUpsert(t *testing.T) {
	t.Parallel()

	if len(accountAllColumns) == len(accountPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := Account{}
	if err = randomize.Struct(seed, &o, accountDBTypes, true); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Account: %s", err)
	}

	count, err := Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, accountDBTypes, false, accountPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize Account struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert Account: %s", err)
	}

	count, err = Accounts().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
