// Code generated by SQLBoiler (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
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

	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"github.com/volatiletech/sqlboiler/queries/qmhelper"
	"github.com/volatiletech/sqlboiler/strmangle"
)

// EconomyShopListItem is an object representing the database table.
type EconomyShopListItem struct {
	GuildID     int64  `boil:"guild_id" json:"guild_id" toml:"guild_id" yaml:"guild_id"`
	LocalID     int64  `boil:"local_id" json:"local_id" toml:"local_id" yaml:"local_id"`
	ListID      int64  `boil:"list_id" json:"list_id" toml:"list_id" yaml:"list_id"`
	Value       string `boil:"value" json:"value" toml:"value" yaml:"value"`
	PurchasedBy int64  `boil:"purchased_by" json:"purchased_by" toml:"purchased_by" yaml:"purchased_by"`

	R *economyShopListItemR `boil:"-" json:"-" toml:"-" yaml:"-"`
	L economyShopListItemL  `boil:"-" json:"-" toml:"-" yaml:"-"`
}

var EconomyShopListItemColumns = struct {
	GuildID     string
	LocalID     string
	ListID      string
	Value       string
	PurchasedBy string
}{
	GuildID:     "guild_id",
	LocalID:     "local_id",
	ListID:      "list_id",
	Value:       "value",
	PurchasedBy: "purchased_by",
}

// Generated where

var EconomyShopListItemWhere = struct {
	GuildID     whereHelperint64
	LocalID     whereHelperint64
	ListID      whereHelperint64
	Value       whereHelperstring
	PurchasedBy whereHelperint64
}{
	GuildID:     whereHelperint64{field: "\"economy_shop_list_items\".\"guild_id\""},
	LocalID:     whereHelperint64{field: "\"economy_shop_list_items\".\"local_id\""},
	ListID:      whereHelperint64{field: "\"economy_shop_list_items\".\"list_id\""},
	Value:       whereHelperstring{field: "\"economy_shop_list_items\".\"value\""},
	PurchasedBy: whereHelperint64{field: "\"economy_shop_list_items\".\"purchased_by\""},
}

// EconomyShopListItemRels is where relationship names are stored.
var EconomyShopListItemRels = struct {
}{}

// economyShopListItemR is where relationships are stored.
type economyShopListItemR struct {
}

// NewStruct creates a new relationship struct
func (*economyShopListItemR) NewStruct() *economyShopListItemR {
	return &economyShopListItemR{}
}

// economyShopListItemL is where Load methods for each relationship are stored.
type economyShopListItemL struct{}

var (
	economyShopListItemAllColumns            = []string{"guild_id", "local_id", "list_id", "value", "purchased_by"}
	economyShopListItemColumnsWithoutDefault = []string{"guild_id", "local_id", "list_id", "value", "purchased_by"}
	economyShopListItemColumnsWithDefault    = []string{}
	economyShopListItemPrimaryKeyColumns     = []string{"guild_id", "local_id"}
)

type (
	// EconomyShopListItemSlice is an alias for a slice of pointers to EconomyShopListItem.
	// This should generally be used opposed to []EconomyShopListItem.
	EconomyShopListItemSlice []*EconomyShopListItem

	economyShopListItemQuery struct {
		*queries.Query
	}
)

// Cache for insert, update and upsert
var (
	economyShopListItemType                 = reflect.TypeOf(&EconomyShopListItem{})
	economyShopListItemMapping              = queries.MakeStructMapping(economyShopListItemType)
	economyShopListItemPrimaryKeyMapping, _ = queries.BindMapping(economyShopListItemType, economyShopListItemMapping, economyShopListItemPrimaryKeyColumns)
	economyShopListItemInsertCacheMut       sync.RWMutex
	economyShopListItemInsertCache          = make(map[string]insertCache)
	economyShopListItemUpdateCacheMut       sync.RWMutex
	economyShopListItemUpdateCache          = make(map[string]updateCache)
	economyShopListItemUpsertCacheMut       sync.RWMutex
	economyShopListItemUpsertCache          = make(map[string]insertCache)
)

var (
	// Force time package dependency for automated UpdatedAt/CreatedAt.
	_ = time.Second
	// Force qmhelper dependency for where clause generation (which doesn't
	// always happen)
	_ = qmhelper.Where
)

// OneG returns a single economyShopListItem record from the query using the global executor.
func (q economyShopListItemQuery) OneG(ctx context.Context) (*EconomyShopListItem, error) {
	return q.One(ctx, boil.GetContextDB())
}

// One returns a single economyShopListItem record from the query.
func (q economyShopListItemQuery) One(ctx context.Context, exec boil.ContextExecutor) (*EconomyShopListItem, error) {
	o := &EconomyShopListItem{}

	queries.SetLimit(q.Query, 1)

	err := q.Bind(ctx, exec, o)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: failed to execute a one query for economy_shop_list_items")
	}

	return o, nil
}

// AllG returns all EconomyShopListItem records from the query using the global executor.
func (q economyShopListItemQuery) AllG(ctx context.Context) (EconomyShopListItemSlice, error) {
	return q.All(ctx, boil.GetContextDB())
}

// All returns all EconomyShopListItem records from the query.
func (q economyShopListItemQuery) All(ctx context.Context, exec boil.ContextExecutor) (EconomyShopListItemSlice, error) {
	var o []*EconomyShopListItem

	err := q.Bind(ctx, exec, &o)
	if err != nil {
		return nil, errors.Wrap(err, "models: failed to assign all query results to EconomyShopListItem slice")
	}

	return o, nil
}

// CountG returns the count of all EconomyShopListItem records in the query, and panics on error.
func (q economyShopListItemQuery) CountG(ctx context.Context) (int64, error) {
	return q.Count(ctx, boil.GetContextDB())
}

// Count returns the count of all EconomyShopListItem records in the query.
func (q economyShopListItemQuery) Count(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to count economy_shop_list_items rows")
	}

	return count, nil
}

// ExistsG checks if the row exists in the table, and panics on error.
func (q economyShopListItemQuery) ExistsG(ctx context.Context) (bool, error) {
	return q.Exists(ctx, boil.GetContextDB())
}

// Exists checks if the row exists in the table.
func (q economyShopListItemQuery) Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error) {
	var count int64

	queries.SetSelect(q.Query, nil)
	queries.SetCount(q.Query)
	queries.SetLimit(q.Query, 1)

	err := q.Query.QueryRowContext(ctx, exec).Scan(&count)
	if err != nil {
		return false, errors.Wrap(err, "models: failed to check if economy_shop_list_items exists")
	}

	return count > 0, nil
}

// EconomyShopListItems retrieves all the records using an executor.
func EconomyShopListItems(mods ...qm.QueryMod) economyShopListItemQuery {
	mods = append(mods, qm.From("\"economy_shop_list_items\""))
	return economyShopListItemQuery{NewQuery(mods...)}
}

// FindEconomyShopListItemG retrieves a single record by ID.
func FindEconomyShopListItemG(ctx context.Context, guildID int64, localID int64, selectCols ...string) (*EconomyShopListItem, error) {
	return FindEconomyShopListItem(ctx, boil.GetContextDB(), guildID, localID, selectCols...)
}

// FindEconomyShopListItem retrieves a single record by ID with an executor.
// If selectCols is empty Find will return all columns.
func FindEconomyShopListItem(ctx context.Context, exec boil.ContextExecutor, guildID int64, localID int64, selectCols ...string) (*EconomyShopListItem, error) {
	economyShopListItemObj := &EconomyShopListItem{}

	sel := "*"
	if len(selectCols) > 0 {
		sel = strings.Join(strmangle.IdentQuoteSlice(dialect.LQ, dialect.RQ, selectCols), ",")
	}
	query := fmt.Sprintf(
		"select %s from \"economy_shop_list_items\" where \"guild_id\"=$1 AND \"local_id\"=$2", sel,
	)

	q := queries.Raw(query, guildID, localID)

	err := q.Bind(ctx, exec, economyShopListItemObj)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, errors.Wrap(err, "models: unable to select from economy_shop_list_items")
	}

	return economyShopListItemObj, nil
}

// InsertG a single record. See Insert for whitelist behavior description.
func (o *EconomyShopListItem) InsertG(ctx context.Context, columns boil.Columns) error {
	return o.Insert(ctx, boil.GetContextDB(), columns)
}

// Insert a single record using an executor.
// See boil.Columns.InsertColumnSet documentation to understand column list inference for inserts.
func (o *EconomyShopListItem) Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error {
	if o == nil {
		return errors.New("models: no economy_shop_list_items provided for insertion")
	}

	var err error

	nzDefaults := queries.NonZeroDefaultSet(economyShopListItemColumnsWithDefault, o)

	key := makeCacheKey(columns, nzDefaults)
	economyShopListItemInsertCacheMut.RLock()
	cache, cached := economyShopListItemInsertCache[key]
	economyShopListItemInsertCacheMut.RUnlock()

	if !cached {
		wl, returnColumns := columns.InsertColumnSet(
			economyShopListItemAllColumns,
			economyShopListItemColumnsWithDefault,
			economyShopListItemColumnsWithoutDefault,
			nzDefaults,
		)

		cache.valueMapping, err = queries.BindMapping(economyShopListItemType, economyShopListItemMapping, wl)
		if err != nil {
			return err
		}
		cache.retMapping, err = queries.BindMapping(economyShopListItemType, economyShopListItemMapping, returnColumns)
		if err != nil {
			return err
		}
		if len(wl) != 0 {
			cache.query = fmt.Sprintf("INSERT INTO \"economy_shop_list_items\" (\"%s\") %%sVALUES (%s)%%s", strings.Join(wl, "\",\""), strmangle.Placeholders(dialect.UseIndexPlaceholders, len(wl), 1, 1))
		} else {
			cache.query = "INSERT INTO \"economy_shop_list_items\" %sDEFAULT VALUES%s"
		}

		var queryOutput, queryReturning string

		if len(cache.retMapping) != 0 {
			queryReturning = fmt.Sprintf(" RETURNING \"%s\"", strings.Join(returnColumns, "\",\""))
		}

		cache.query = fmt.Sprintf(cache.query, queryOutput, queryReturning)
	}

	value := reflect.Indirect(reflect.ValueOf(o))
	vals := queries.ValuesFromMapping(value, cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
	}

	if len(cache.retMapping) != 0 {
		err = exec.QueryRowContext(ctx, cache.query, vals...).Scan(queries.PtrsFromMapping(value, cache.retMapping)...)
	} else {
		_, err = exec.ExecContext(ctx, cache.query, vals...)
	}

	if err != nil {
		return errors.Wrap(err, "models: unable to insert into economy_shop_list_items")
	}

	if !cached {
		economyShopListItemInsertCacheMut.Lock()
		economyShopListItemInsertCache[key] = cache
		economyShopListItemInsertCacheMut.Unlock()
	}

	return nil
}

// UpdateG a single EconomyShopListItem record using the global executor.
// See Update for more documentation.
func (o *EconomyShopListItem) UpdateG(ctx context.Context, columns boil.Columns) (int64, error) {
	return o.Update(ctx, boil.GetContextDB(), columns)
}

// Update uses an executor to update the EconomyShopListItem.
// See boil.Columns.UpdateColumnSet documentation to understand column list inference for updates.
// Update does not automatically update the record in case of default values. Use .Reload() to refresh the records.
func (o *EconomyShopListItem) Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error) {
	var err error
	key := makeCacheKey(columns, nil)
	economyShopListItemUpdateCacheMut.RLock()
	cache, cached := economyShopListItemUpdateCache[key]
	economyShopListItemUpdateCacheMut.RUnlock()

	if !cached {
		wl := columns.UpdateColumnSet(
			economyShopListItemAllColumns,
			economyShopListItemPrimaryKeyColumns,
		)

		if !columns.IsWhitelist() {
			wl = strmangle.SetComplement(wl, []string{"created_at"})
		}
		if len(wl) == 0 {
			return 0, errors.New("models: unable to update economy_shop_list_items, could not build whitelist")
		}

		cache.query = fmt.Sprintf("UPDATE \"economy_shop_list_items\" SET %s WHERE %s",
			strmangle.SetParamNames("\"", "\"", 1, wl),
			strmangle.WhereClause("\"", "\"", len(wl)+1, economyShopListItemPrimaryKeyColumns),
		)
		cache.valueMapping, err = queries.BindMapping(economyShopListItemType, economyShopListItemMapping, append(wl, economyShopListItemPrimaryKeyColumns...))
		if err != nil {
			return 0, err
		}
	}

	values := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), cache.valueMapping)

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, values)
	}

	var result sql.Result
	result, err = exec.ExecContext(ctx, cache.query, values...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update economy_shop_list_items row")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by update for economy_shop_list_items")
	}

	if !cached {
		economyShopListItemUpdateCacheMut.Lock()
		economyShopListItemUpdateCache[key] = cache
		economyShopListItemUpdateCacheMut.Unlock()
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (q economyShopListItemQuery) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return q.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values.
func (q economyShopListItemQuery) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
	queries.SetUpdate(q.Query, cols)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all for economy_shop_list_items")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected for economy_shop_list_items")
	}

	return rowsAff, nil
}

// UpdateAllG updates all rows with the specified column values.
func (o EconomyShopListItemSlice) UpdateAllG(ctx context.Context, cols M) (int64, error) {
	return o.UpdateAll(ctx, boil.GetContextDB(), cols)
}

// UpdateAll updates all rows with the specified column values, using an executor.
func (o EconomyShopListItemSlice) UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols M) (int64, error) {
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
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), economyShopListItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := fmt.Sprintf("UPDATE \"economy_shop_list_items\" SET %s WHERE %s",
		strmangle.SetParamNames("\"", "\"", 1, colNames),
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), len(colNames)+1, economyShopListItemPrimaryKeyColumns, len(o)))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to update all in economyShopListItem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to retrieve rows affected all in update all economyShopListItem")
	}
	return rowsAff, nil
}

// UpsertG attempts an insert, and does an update or ignore on conflict.
func (o *EconomyShopListItem) UpsertG(ctx context.Context, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	return o.Upsert(ctx, boil.GetContextDB(), updateOnConflict, conflictColumns, updateColumns, insertColumns)
}

// Upsert attempts an insert using an executor, and does an update or ignore on conflict.
// See boil.Columns documentation for how to properly use updateColumns and insertColumns.
func (o *EconomyShopListItem) Upsert(ctx context.Context, exec boil.ContextExecutor, updateOnConflict bool, conflictColumns []string, updateColumns, insertColumns boil.Columns) error {
	if o == nil {
		return errors.New("models: no economy_shop_list_items provided for upsert")
	}

	nzDefaults := queries.NonZeroDefaultSet(economyShopListItemColumnsWithDefault, o)

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

	economyShopListItemUpsertCacheMut.RLock()
	cache, cached := economyShopListItemUpsertCache[key]
	economyShopListItemUpsertCacheMut.RUnlock()

	var err error

	if !cached {
		insert, ret := insertColumns.InsertColumnSet(
			economyShopListItemAllColumns,
			economyShopListItemColumnsWithDefault,
			economyShopListItemColumnsWithoutDefault,
			nzDefaults,
		)
		update := updateColumns.UpdateColumnSet(
			economyShopListItemAllColumns,
			economyShopListItemPrimaryKeyColumns,
		)

		if updateOnConflict && len(update) == 0 {
			return errors.New("models: unable to upsert economy_shop_list_items, could not build update column list")
		}

		conflict := conflictColumns
		if len(conflict) == 0 {
			conflict = make([]string, len(economyShopListItemPrimaryKeyColumns))
			copy(conflict, economyShopListItemPrimaryKeyColumns)
		}
		cache.query = buildUpsertQueryPostgres(dialect, "\"economy_shop_list_items\"", updateOnConflict, ret, update, conflict, insert)

		cache.valueMapping, err = queries.BindMapping(economyShopListItemType, economyShopListItemMapping, insert)
		if err != nil {
			return err
		}
		if len(ret) != 0 {
			cache.retMapping, err = queries.BindMapping(economyShopListItemType, economyShopListItemMapping, ret)
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

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, cache.query)
		fmt.Fprintln(boil.DebugWriter, vals)
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
		return errors.Wrap(err, "models: unable to upsert economy_shop_list_items")
	}

	if !cached {
		economyShopListItemUpsertCacheMut.Lock()
		economyShopListItemUpsertCache[key] = cache
		economyShopListItemUpsertCacheMut.Unlock()
	}

	return nil
}

// DeleteG deletes a single EconomyShopListItem record.
// DeleteG will match against the primary key column to find the record to delete.
func (o *EconomyShopListItem) DeleteG(ctx context.Context) (int64, error) {
	return o.Delete(ctx, boil.GetContextDB())
}

// Delete deletes a single EconomyShopListItem record with an executor.
// Delete will match against the primary key column to find the record to delete.
func (o *EconomyShopListItem) Delete(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if o == nil {
		return 0, errors.New("models: no EconomyShopListItem provided for delete")
	}

	args := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(o)), economyShopListItemPrimaryKeyMapping)
	sql := "DELETE FROM \"economy_shop_list_items\" WHERE \"guild_id\"=$1 AND \"local_id\"=$2"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args...)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete from economy_shop_list_items")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by delete for economy_shop_list_items")
	}

	return rowsAff, nil
}

// DeleteAll deletes all matching rows.
func (q economyShopListItemQuery) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if q.Query == nil {
		return 0, errors.New("models: no economyShopListItemQuery provided for delete all")
	}

	queries.SetDelete(q.Query)

	result, err := q.Query.ExecContext(ctx, exec)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from economy_shop_list_items")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for economy_shop_list_items")
	}

	return rowsAff, nil
}

// DeleteAllG deletes all rows in the slice.
func (o EconomyShopListItemSlice) DeleteAllG(ctx context.Context) (int64, error) {
	return o.DeleteAll(ctx, boil.GetContextDB())
}

// DeleteAll deletes all rows in the slice, using an executor.
func (o EconomyShopListItemSlice) DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error) {
	if len(o) == 0 {
		return 0, nil
	}

	var args []interface{}
	for _, obj := range o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), economyShopListItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "DELETE FROM \"economy_shop_list_items\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, economyShopListItemPrimaryKeyColumns, len(o))

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, args)
	}

	result, err := exec.ExecContext(ctx, sql, args...)
	if err != nil {
		return 0, errors.Wrap(err, "models: unable to delete all from economyShopListItem slice")
	}

	rowsAff, err := result.RowsAffected()
	if err != nil {
		return 0, errors.Wrap(err, "models: failed to get rows affected by deleteall for economy_shop_list_items")
	}

	return rowsAff, nil
}

// ReloadG refetches the object from the database using the primary keys.
func (o *EconomyShopListItem) ReloadG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: no EconomyShopListItem provided for reload")
	}

	return o.Reload(ctx, boil.GetContextDB())
}

// Reload refetches the object from the database
// using the primary keys with an executor.
func (o *EconomyShopListItem) Reload(ctx context.Context, exec boil.ContextExecutor) error {
	ret, err := FindEconomyShopListItem(ctx, exec, o.GuildID, o.LocalID)
	if err != nil {
		return err
	}

	*o = *ret
	return nil
}

// ReloadAllG refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *EconomyShopListItemSlice) ReloadAllG(ctx context.Context) error {
	if o == nil {
		return errors.New("models: empty EconomyShopListItemSlice provided for reload all")
	}

	return o.ReloadAll(ctx, boil.GetContextDB())
}

// ReloadAll refetches every row with matching primary key column values
// and overwrites the original object slice with the newly updated slice.
func (o *EconomyShopListItemSlice) ReloadAll(ctx context.Context, exec boil.ContextExecutor) error {
	if o == nil || len(*o) == 0 {
		return nil
	}

	slice := EconomyShopListItemSlice{}
	var args []interface{}
	for _, obj := range *o {
		pkeyArgs := queries.ValuesFromMapping(reflect.Indirect(reflect.ValueOf(obj)), economyShopListItemPrimaryKeyMapping)
		args = append(args, pkeyArgs...)
	}

	sql := "SELECT \"economy_shop_list_items\".* FROM \"economy_shop_list_items\" WHERE " +
		strmangle.WhereClauseRepeated(string(dialect.LQ), string(dialect.RQ), 1, economyShopListItemPrimaryKeyColumns, len(*o))

	q := queries.Raw(sql, args...)

	err := q.Bind(ctx, exec, &slice)
	if err != nil {
		return errors.Wrap(err, "models: unable to reload all in EconomyShopListItemSlice")
	}

	*o = slice

	return nil
}

// EconomyShopListItemExistsG checks if the EconomyShopListItem row exists.
func EconomyShopListItemExistsG(ctx context.Context, guildID int64, localID int64) (bool, error) {
	return EconomyShopListItemExists(ctx, boil.GetContextDB(), guildID, localID)
}

// EconomyShopListItemExists checks if the EconomyShopListItem row exists.
func EconomyShopListItemExists(ctx context.Context, exec boil.ContextExecutor, guildID int64, localID int64) (bool, error) {
	var exists bool
	sql := "select exists(select 1 from \"economy_shop_list_items\" where \"guild_id\"=$1 AND \"local_id\"=$2 limit 1)"

	if boil.DebugMode {
		fmt.Fprintln(boil.DebugWriter, sql)
		fmt.Fprintln(boil.DebugWriter, guildID, localID)
	}

	row := exec.QueryRowContext(ctx, sql, guildID, localID)

	err := row.Scan(&exists)
	if err != nil {
		return false, errors.Wrap(err, "models: unable to check if economy_shop_list_items exists")
	}

	return exists, nil
}
