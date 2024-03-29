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

func testProductCategories(t *testing.T) {
	t.Parallel()

	query := ProductCategories()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testProductCategoriesDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
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

	count, err := ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProductCategoriesQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := ProductCategories().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProductCategoriesSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ProductCategorySlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testProductCategoriesExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := ProductCategoryExists(ctx, tx, o.ID)
	if err != nil {
		t.Errorf("Unable to check if ProductCategory exists: %s", err)
	}
	if !e {
		t.Errorf("Expected ProductCategoryExists to return true, but got false.")
	}
}

func testProductCategoriesFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	productCategoryFound, err := FindProductCategory(ctx, tx, o.ID)
	if err != nil {
		t.Error(err)
	}

	if productCategoryFound == nil {
		t.Error("want a record, got nil")
	}
}

func testProductCategoriesBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = ProductCategories().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testProductCategoriesOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := ProductCategories().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testProductCategoriesAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	productCategoryOne := &ProductCategory{}
	productCategoryTwo := &ProductCategory{}
	if err = randomize.Struct(seed, productCategoryOne, productCategoryDBTypes, false, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}
	if err = randomize.Struct(seed, productCategoryTwo, productCategoryDBTypes, false, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = productCategoryOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = productCategoryTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ProductCategories().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testProductCategoriesCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	productCategoryOne := &ProductCategory{}
	productCategoryTwo := &ProductCategory{}
	if err = randomize.Struct(seed, productCategoryOne, productCategoryDBTypes, false, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}
	if err = randomize.Struct(seed, productCategoryTwo, productCategoryDBTypes, false, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = productCategoryOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = productCategoryTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testProductCategoriesInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testProductCategoriesInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(productCategoryColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testProductCategoryToOneCategoryUsingCategory(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local ProductCategory
	var foreign Category

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, productCategoryDBTypes, false, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}
	if err := randomize.Struct(seed, &foreign, categoryDBTypes, false, categoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize Category struct: %s", err)
	}

	if err := foreign.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	local.CategoryID = foreign.ID
	if err := local.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	check, err := local.Category().One(ctx, tx)
	if err != nil {
		t.Fatal(err)
	}

	if check.ID != foreign.ID {
		t.Errorf("want: %v, got %v", foreign.ID, check.ID)
	}

	slice := ProductCategorySlice{&local}
	if err = local.L.LoadCategory(ctx, tx, false, (*[]*ProductCategory)(&slice), nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Category == nil {
		t.Error("struct should have been eager loaded")
	}

	local.R.Category = nil
	if err = local.L.LoadCategory(ctx, tx, true, &local, nil); err != nil {
		t.Fatal(err)
	}
	if local.R.Category == nil {
		t.Error("struct should have been eager loaded")
	}
}

func testProductCategoryToOneProductUsingProduct(t *testing.T) {
	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var local ProductCategory
	var foreign Product

	seed := randomize.NewSeed()
	if err := randomize.Struct(seed, &local, productCategoryDBTypes, false, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
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

	slice := ProductCategorySlice{&local}
	if err = local.L.LoadProduct(ctx, tx, false, (*[]*ProductCategory)(&slice), nil); err != nil {
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

func testProductCategoryToOneSetOpCategoryUsingCategory(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ProductCategory
	var b, c Category

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, productCategoryDBTypes, false, strmangle.SetComplement(productCategoryPrimaryKeyColumns, productCategoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &b, categoryDBTypes, false, strmangle.SetComplement(categoryPrimaryKeyColumns, categoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}
	if err = randomize.Struct(seed, &c, categoryDBTypes, false, strmangle.SetComplement(categoryPrimaryKeyColumns, categoryColumnsWithoutDefault)...); err != nil {
		t.Fatal(err)
	}

	if err := a.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}
	if err = b.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Fatal(err)
	}

	for i, x := range []*Category{&b, &c} {
		err = a.SetCategory(ctx, tx, i != 0, x)
		if err != nil {
			t.Fatal(err)
		}

		if a.R.Category != x {
			t.Error("relationship struct not set to correct value")
		}

		if x.R.ProductCategories[0] != &a {
			t.Error("failed to append to foreign relationship struct")
		}
		if a.CategoryID != x.ID {
			t.Error("foreign key was wrong value", a.CategoryID)
		}

		zero := reflect.Zero(reflect.TypeOf(a.CategoryID))
		reflect.Indirect(reflect.ValueOf(&a.CategoryID)).Set(zero)

		if err = a.Reload(ctx, tx); err != nil {
			t.Fatal("failed to reload", err)
		}

		if a.CategoryID != x.ID {
			t.Error("foreign key was wrong value", a.CategoryID, x.ID)
		}
	}
}
func testProductCategoryToOneSetOpProductUsingProduct(t *testing.T) {
	var err error

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()

	var a ProductCategory
	var b, c Product

	seed := randomize.NewSeed()
	if err = randomize.Struct(seed, &a, productCategoryDBTypes, false, strmangle.SetComplement(productCategoryPrimaryKeyColumns, productCategoryColumnsWithoutDefault)...); err != nil {
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

		if x.R.ProductCategories[0] != &a {
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

func testProductCategoriesReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
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

func testProductCategoriesReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := ProductCategorySlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testProductCategoriesSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := ProductCategories().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	productCategoryDBTypes = map[string]string{`ID`: `character`, `ProductID`: `character`, `CategoryID`: `character`}
	_                      = bytes.MinRead
)

func testProductCategoriesUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(productCategoryPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(productCategoryAllColumns) == len(productCategoryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testProductCategoriesSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(productCategoryAllColumns) == len(productCategoryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &ProductCategory{}
	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, productCategoryDBTypes, true, productCategoryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(productCategoryAllColumns, productCategoryPrimaryKeyColumns) {
		fields = productCategoryAllColumns
	} else {
		fields = strmangle.SetComplement(
			productCategoryAllColumns,
			productCategoryPrimaryKeyColumns,
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

	slice := ProductCategorySlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testProductCategoriesUpsert(t *testing.T) {
	t.Parallel()

	if len(productCategoryAllColumns) == len(productCategoryPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := ProductCategory{}
	if err = randomize.Struct(seed, &o, productCategoryDBTypes, true); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ProductCategory: %s", err)
	}

	count, err := ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, productCategoryDBTypes, false, productCategoryPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize ProductCategory struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert ProductCategory: %s", err)
	}

	count, err = ProductCategories().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
