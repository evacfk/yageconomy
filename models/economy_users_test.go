// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

func testEconomyUsers(t *testing.T) {
	t.Parallel()

	query := EconomyUsers()

	if query.Query == nil {
		t.Error("expected a query, got nothing")
	}
}

func testEconomyUsersDelete(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
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

	count, err := EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testEconomyUsersQueryDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if rowsAff, err := EconomyUsers().DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testEconomyUsersSliceDeleteAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := EconomyUserSlice{o}

	if rowsAff, err := slice.DeleteAll(ctx, tx); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only have deleted one row, but affected:", rowsAff)
	}

	count, err := EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("want zero records, got:", count)
	}
}

func testEconomyUsersExists(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	e, err := EconomyUserExists(ctx, tx, o.GuildID, o.UserID)
	if err != nil {
		t.Errorf("Unable to check if EconomyUser exists: %s", err)
	}
	if !e {
		t.Errorf("Expected EconomyUserExists to return true, but got false.")
	}
}

func testEconomyUsersFind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	economyUserFound, err := FindEconomyUser(ctx, tx, o.GuildID, o.UserID)
	if err != nil {
		t.Error(err)
	}

	if economyUserFound == nil {
		t.Error("want a record, got nil")
	}
}

func testEconomyUsersBind(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if err = EconomyUsers().Bind(ctx, tx, o); err != nil {
		t.Error(err)
	}
}

func testEconomyUsersOne(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	if x, err := EconomyUsers().One(ctx, tx); err != nil {
		t.Error(err)
	} else if x == nil {
		t.Error("expected to get a non nil record")
	}
}

func testEconomyUsersAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	economyUserOne := &EconomyUser{}
	economyUserTwo := &EconomyUser{}
	if err = randomize.Struct(seed, economyUserOne, economyUserDBTypes, false, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}
	if err = randomize.Struct(seed, economyUserTwo, economyUserDBTypes, false, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = economyUserOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = economyUserTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := EconomyUsers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 2 {
		t.Error("want 2 records, got:", len(slice))
	}
}

func testEconomyUsersCount(t *testing.T) {
	t.Parallel()

	var err error
	seed := randomize.NewSeed()
	economyUserOne := &EconomyUser{}
	economyUserTwo := &EconomyUser{}
	if err = randomize.Struct(seed, economyUserOne, economyUserDBTypes, false, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}
	if err = randomize.Struct(seed, economyUserTwo, economyUserDBTypes, false, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = economyUserOne.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}
	if err = economyUserTwo.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 2 {
		t.Error("want 2 records, got:", count)
	}
}

func testEconomyUsersInsert(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testEconomyUsersInsertWhitelist(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Whitelist(economyUserColumnsWithoutDefault...)); err != nil {
		t.Error(err)
	}

	count, err := EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}
}

func testEconomyUsersReload(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
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

func testEconomyUsersReloadAll(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice := EconomyUserSlice{o}

	if err = slice.ReloadAll(ctx, tx); err != nil {
		t.Error(err)
	}
}

func testEconomyUsersSelect(t *testing.T) {
	t.Parallel()

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	slice, err := EconomyUsers().All(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if len(slice) != 1 {
		t.Error("want one record, got:", len(slice))
	}
}

var (
	economyUserDBTypes = map[string]string{`GuildID`: `bigint`, `UserID`: `bigint`, `MoneyBank`: `bigint`, `MoneyWallet`: `bigint`, `LastDailyClaim`: `timestamp with time zone`, `LastChatmoneyClaim`: `timestamp with time zone`, `LastFishing`: `timestamp with time zone`, `WaifudBy`: `bigint`, `Waifus`: `ARRAYbigint`, `WaifuItemWorth`: `bigint`, `WaifuLastClaimAmount`: `bigint`, `WaifuExtraWorth`: `bigint`, `WaifuAffinityTowards`: `bigint`, `WaifuDivorces`: `integer`, `WaifuAffinityChanges`: `integer`, `FishCaugth`: `bigint`, `GamblingBoostPercentage`: `integer`, `LastInterestUpdate`: `timestamp with time zone`, `LastRobAttempt`: `timestamp with time zone`}
	_                  = bytes.MinRead
)

func testEconomyUsersUpdate(t *testing.T) {
	t.Parallel()

	if 0 == len(economyUserPrimaryKeyColumns) {
		t.Skip("Skipping table with no primary key columns")
	}
	if len(economyUserColumns) == len(economyUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	if rowsAff, err := o.Update(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("should only affect one row but affected", rowsAff)
	}
}

func testEconomyUsersSliceUpdateAll(t *testing.T) {
	t.Parallel()

	if len(economyUserColumns) == len(economyUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	o := &EconomyUser{}
	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserColumnsWithDefault...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Insert(ctx, tx, boil.Infer()); err != nil {
		t.Error(err)
	}

	count, err := EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}

	if count != 1 {
		t.Error("want one record, got:", count)
	}

	if err = randomize.Struct(seed, o, economyUserDBTypes, true, economyUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	// Remove Primary keys and unique columns from what we plan to update
	var fields []string
	if strmangle.StringSliceMatch(economyUserColumns, economyUserPrimaryKeyColumns) {
		fields = economyUserColumns
	} else {
		fields = strmangle.SetComplement(
			economyUserColumns,
			economyUserPrimaryKeyColumns,
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

	slice := EconomyUserSlice{o}
	if rowsAff, err := slice.UpdateAll(ctx, tx, updateMap); err != nil {
		t.Error(err)
	} else if rowsAff != 1 {
		t.Error("wanted one record updated but got", rowsAff)
	}
}

func testEconomyUsersUpsert(t *testing.T) {
	t.Parallel()

	if len(economyUserColumns) == len(economyUserPrimaryKeyColumns) {
		t.Skip("Skipping table with only primary key columns")
	}

	seed := randomize.NewSeed()
	var err error
	// Attempt the INSERT side of an UPSERT
	o := EconomyUser{}
	if err = randomize.Struct(seed, &o, economyUserDBTypes, true); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	ctx := context.Background()
	tx := MustTx(boil.BeginTx(ctx, nil))
	defer func() { _ = tx.Rollback() }()
	if err = o.Upsert(ctx, tx, false, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert EconomyUser: %s", err)
	}

	count, err := EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}

	// Attempt the UPDATE side of an UPSERT
	if err = randomize.Struct(seed, &o, economyUserDBTypes, false, economyUserPrimaryKeyColumns...); err != nil {
		t.Errorf("Unable to randomize EconomyUser struct: %s", err)
	}

	if err = o.Upsert(ctx, tx, true, nil, boil.Infer(), boil.Infer()); err != nil {
		t.Errorf("Unable to upsert EconomyUser: %s", err)
	}

	count, err = EconomyUsers().Count(ctx, tx)
	if err != nil {
		t.Error(err)
	}
	if count != 1 {
		t.Error("want one record, got:", count)
	}
}
