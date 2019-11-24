package cmsql

import (
	"context"
	"database/sql"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sq"
	"etop.vn/backend/pkg/common/sq/core"
	"etop.vn/capi/dot"
)

// QueryInterface wraps common/sql.QueryInterface
type QueryInterface interface {
	Get(obj core.IGet, preds ...interface{}) (bool, error)
	Find(objs core.IFind, preds ...interface{}) error
	Insert(objs ...core.IInsert) (int64, error)
	Update(objs ...core.IUpdate) (int64, error)
	UpdateMap(m map[string]interface{}) (int64, error)
	Delete(obj core.ITableName) (int64, error)
	Count(obj core.ITableName, preds ...interface{}) (uint64, error)
	FindRows(objs core.IFind, preds ...interface{}) (core.Opts, *sql.Rows, error)

	ShouldGet(obj core.IGet, preds ...interface{}) error
	ShouldInsert(objs ...core.IInsert) error
	ShouldUpdate(objs ...core.IUpdate) error
	ShouldUpdateMap(m map[string]interface{}) error
	ShouldDelete(obj core.ITableName) error

	Table(name string) Query
	Prefix(sql string, args ...interface{}) Query
	Select(cols ...string) Query
	From(table string) Query
	SQL(args ...interface{}) Query
	Where(args ...interface{}) Query
	OrderBy(orderBys ...string) Query
	GroupBy(groupBys ...string) Query
	Limit(limit uint64) Query
	Offset(offset uint64) Query
	Suffix(sql string, args ...interface{}) Query
	UpdateAll() Query
	In(column string, args ...interface{}) Query
	NotIn(column string, args ...interface{}) Query
	Exists(column string, exists bool) Query
	IsNull(column string, null bool) Query
	InOrEqIDs(column string, args []dot.ID) Query

	Preload(table string, preds ...interface{}) Query
	Apply(funcs ...func(core.CommonQuery)) Query
}

type DBX interface {
	NewQuery() Query
	QueryInterface
}

type TxKey struct{ dbID int64 }

type Transactioner interface {
	InTransaction(ctx context.Context, callback func(QueryInterface) error) (err error)
}

// Tx rexports common/sql.Tx
type Tx interface {
	Commit() error
	Rollback() error
	NewQuery() Query

	sq.DBInterface
	QueryInterface
}

var _ DBX = Database{}
var _ DBX = tx{}
var _ QueryInterface = Query{}
var _ QueryInterface = tx{}

//-- Query --//

// Get ...
func (q Query) Get(obj core.IGet, preds ...interface{}) (bool, error) {
	return q.query.Get(obj, preds...)
}

// Find ...
func (q Query) Find(objs core.IFind, preds ...interface{}) error {
	return q.query.Find(objs, preds...)
}

// FindRows ...
func (q Query) FindRows(objs core.IFind, preds ...interface{}) (core.Opts, *sql.Rows, error) {
	return q.query.FindRows(objs, preds...)
}

// Insert ...
func (q Query) Insert(objs ...core.IInsert) (int64, error) {
	return q.query.Insert(objs...)
}

// UpdateInfo ...
func (q Query) Update(objs ...core.IUpdate) (int64, error) {
	return q.query.Update(objs...)
}

// UpdateMap ...
func (q Query) UpdateMap(obj map[string]interface{}) (int64, error) {
	return q.query.UpdateMap(obj)
}

// Delete ...
func (q Query) Delete(obj core.ITableName) (int64, error) {
	return q.query.Delete(obj)
}

// Count ...
func (q Query) Count(obj core.ITableName, preds ...interface{}) (uint64, error) {
	return q.query.Count(obj, preds...)
}

// Table ...
func (q Query) Table(name string) Query {
	return Query{q.query.Table(name)}
}

// Prefix ...
func (q Query) Prefix(sql string, args ...interface{}) Query {
	return Query{q.query.Prefix(sql, args...)}
}

// Select ...
func (q Query) Select(cols ...string) Query {
	return Query{q.query.Select(cols...)}
}

// From ...
func (q Query) From(table string) Query {
	return Query{q.query.From(table)}
}

// SQL ...
func (q Query) SQL(args ...interface{}) Query {
	return Query{q.query.SQL(args...)}
}

// Where ...
func (q Query) Where(args ...interface{}) Query {
	return Query{q.query.Where(args...)}
}

// OrderBy ...
func (q Query) OrderBy(orderBys ...string) Query {
	return Query{q.query.OrderBy(orderBys...)}
}

// GroupBy ...
func (q Query) GroupBy(groupBys ...string) Query {
	return Query{q.query.GroupBy(groupBys...)}
}

// Limit ...
func (q Query) Limit(limit uint64) Query {
	return Query{q.query.Limit(limit)}
}

// Offset ...
func (q Query) Offset(offset uint64) Query {
	return Query{q.query.Offset(offset)}
}

// Suffix ...
func (q Query) Suffix(sql string, args ...interface{}) Query {
	return Query{q.query.Suffix(sql)}
}

// UpdateAll ...
func (q Query) UpdateAll() Query {
	return Query{q.query.UpdateAll()}
}

// In ...
func (q Query) In(column string, args ...interface{}) Query {
	return Query{q.query.In(column, args...)}
}

// NotIn ...
func (q Query) NotIn(column string, args ...interface{}) Query {
	return Query{q.query.NotIn(column, args...)}
}

// Exists ...
func (q Query) Exists(column string, exists bool) Query {
	return Query{q.query.Exists(column, exists)}
}

// IsNull ...
func (q Query) IsNull(column string, null bool) Query {
	return Query{q.query.IsNull(column, null)}
}

// InOrEqIDs ...
func (q Query) InOrEqIDs(column string, ids []dot.ID) Query {
	return inOrEqIDs(q.query, column, ids)
}

func (q Query) Preload(table string, preds ...interface{}) Query {
	return Query{q.query.Preload(table, preds...)}
}

func (q Query) Apply(funcs ...func(core.CommonQuery)) Query {
	return Query{q.query.Apply(funcs...)}
}

//-- Database --//

// Get ...
func (db Database) Get(obj core.IGet, preds ...interface{}) (bool, error) {
	return db.db.Get(obj, preds...)
}

// Find ...
func (db Database) Find(objs core.IFind, preds ...interface{}) error {
	return db.db.Find(objs, preds...)
}

// FindRows ...
func (db Database) FindRows(objs core.IFind, preds ...interface{}) (core.Opts, *sql.Rows, error) {
	return db.db.FindRows(objs, preds...)
}

// Insert ...
func (db Database) Insert(objs ...core.IInsert) (int64, error) {
	return db.db.Insert(objs...)
}

// UpdateInfo ...
func (db Database) Update(objs ...core.IUpdate) (int64, error) {
	return db.db.Update(objs...)
}

// UpdateMap ...
func (db Database) UpdateMap(obj map[string]interface{}) (int64, error) {
	return db.db.UpdateMap(obj)
}

// Delete ...
func (db Database) Delete(obj core.ITableName) (int64, error) {
	return db.db.Delete(obj)
}

// Count ...
func (db Database) Count(obj core.ITableName, preds ...interface{}) (uint64, error) {
	return db.db.Count(obj, preds...)
}

// Table ...
func (db Database) Table(name string) Query {
	return Query{db.db.Table(name)}
}

// Prefix ...
func (db Database) Prefix(sql string, args ...interface{}) Query {
	return Query{db.db.Prefix(sql, args...)}
}

// Select ...
func (db Database) Select(cols ...string) Query {
	return Query{db.db.Select(cols...)}
}

// From ...
func (db Database) From(table string) Query {
	return Query{db.db.From(table)}
}

// SQL ...
func (db Database) SQL(args ...interface{}) Query {
	return Query{db.db.SQL(args...)}
}

// Where ...
func (db Database) Where(args ...interface{}) Query {
	return Query{db.db.Where(args...)}
}

// OrderBy ...
func (db Database) OrderBy(orderBys ...string) Query {
	return Query{db.db.OrderBy(orderBys...)}
}

// GroupBy ...
func (db Database) GroupBy(groupBys ...string) Query {
	return Query{db.db.GroupBy(groupBys...)}
}

// Limit ...
func (db Database) Limit(limit uint64) Query {
	return Query{db.db.Limit(limit)}
}

// Offset ...
func (db Database) Offset(offset uint64) Query {
	return Query{db.db.Offset(offset)}
}

// Suffix ...
func (db Database) Suffix(sql string, args ...interface{}) Query {
	return Query{db.db.Suffix(sql)}
}

// UpdateAll ...
func (db Database) UpdateAll() Query {
	return Query{db.db.UpdateAll()}
}

// In ...
func (db Database) In(column string, args ...interface{}) Query {
	return Query{db.db.In(column, args...)}
}

// NotIn ...
func (db Database) NotIn(column string, args ...interface{}) Query {
	return Query{db.db.NotIn(column, args...)}
}

// Exists ...
func (db Database) Exists(column string, exists bool) Query {
	return Query{db.db.Exists(column, exists)}
}

// IsNull ...
func (db Database) IsNull(column string, null bool) Query {
	return Query{db.db.In(column, null)}
}

// InOrEqIDs ...
func (db Database) InOrEqIDs(column string, ids []dot.ID) Query {
	return inOrEqIDs(&db.db, column, ids)
}

func (db Database) Preload(table string, preds ...interface{}) Query {
	return Query{db.db.Preload(table, preds...)}
}

func (db Database) Apply(funcs ...func(core.CommonQuery)) Query {
	return Query{db.db.Apply(funcs...)}
}

//-- Tx --//

func (tx tx) NewQuery() Query {
	return Query{tx.tx.NewQuery()}
}

// Get ...
func (tx tx) Get(obj core.IGet, preds ...interface{}) (bool, error) {
	return tx.tx.Get(obj, preds...)
}

// Find ...
func (tx tx) Find(objs core.IFind, preds ...interface{}) error {
	return tx.tx.Find(objs, preds...)
}

// FindRows ...
func (tx tx) FindRows(objs core.IFind, preds ...interface{}) (core.Opts, *sql.Rows, error) {
	return tx.tx.FindRows(objs, preds...)
}

// Insert ...
func (tx tx) Insert(objs ...core.IInsert) (int64, error) {
	return tx.tx.Insert(objs...)
}

// UpdateInfo ...
func (tx tx) Update(objs ...core.IUpdate) (int64, error) {
	return tx.tx.Update(objs...)
}

// UpdateMap ...
func (tx tx) UpdateMap(obj map[string]interface{}) (int64, error) {
	return tx.tx.UpdateMap(obj)
}

// Delete ...
func (tx tx) Delete(obj core.ITableName) (int64, error) {
	return tx.tx.Delete(obj)
}

// Count ...
func (tx tx) Count(obj core.ITableName, preds ...interface{}) (uint64, error) {
	return tx.tx.Count(obj, preds...)
}

// Table ...
func (tx tx) Table(name string) Query {
	return Query{tx.tx.Table(name)}
}

// Prefix ...
func (tx tx) Prefix(sql string, args ...interface{}) Query {
	return Query{tx.tx.Prefix(sql, args...)}
}

// Select ...
func (tx tx) Select(cols ...string) Query {
	return Query{tx.tx.Select(cols...)}
}

// From ...
func (tx tx) From(table string) Query {
	return Query{tx.tx.From(table)}
}

// SQL ...
func (tx tx) SQL(args ...interface{}) Query {
	return Query{tx.tx.SQL(args...)}
}

// Where ...
func (tx tx) Where(args ...interface{}) Query {
	return Query{tx.tx.Where(args...)}
}

// OrderBy ...
func (tx tx) OrderBy(orderBys ...string) Query {
	return Query{tx.tx.OrderBy(orderBys...)}
}

// GroupBy ...
func (tx tx) GroupBy(groupBys ...string) Query {
	return Query{tx.tx.GroupBy(groupBys...)}
}

// Limit ...
func (tx tx) Limit(limit uint64) Query {
	return Query{tx.tx.Limit(limit)}
}

// Offset ...
func (tx tx) Offset(offset uint64) Query {
	return Query{tx.tx.Offset(offset)}
}

// Suffix ...
func (tx tx) Suffix(sql string, args ...interface{}) Query {
	return Query{tx.tx.Suffix(sql)}
}

// UpdateAll ...
func (tx tx) UpdateAll() Query {
	return Query{tx.tx.UpdateAll()}
}

// In ...
func (tx tx) In(column string, args ...interface{}) Query {
	return Query{tx.tx.In(column, args...)}
}

// NotIn ...
func (tx tx) NotIn(column string, args ...interface{}) Query {
	return Query{tx.tx.NotIn(column, args...)}
}

// Exists ...
func (tx tx) Exists(column string, exists bool) Query {
	return Query{tx.tx.Exists(column, exists)}
}

// IsNull ...
func (tx tx) IsNull(column string, null bool) Query {
	return Query{tx.tx.In(column, null)}
}

// InOrEqIDs ...
func (tx tx) InOrEqIDs(column string, ids []dot.ID) Query {
	return inOrEqIDs(tx.tx, column, ids)
}

func (tx tx) Preload(table string, preds ...interface{}) Query {
	return Query{tx.tx.Preload(table, preds...)}
}

func (tx tx) Apply(funcs ...func(core.CommonQuery)) Query {
	return Query{tx.tx.Apply(funcs...)}
}

//-- Should --//

// ShouldGet ...
func (db Database) ShouldGet(obj core.IGet, preds ...interface{}) error {
	return ShouldGet(db.db.Get(obj, preds...))
}

// ShouldUpdate ...
func (db Database) ShouldUpdate(obj ...core.IUpdate) error {
	return Should(db.db.Update(obj...))
}

// ShouldUpdateMap ...
func (db Database) ShouldUpdateMap(obj map[string]interface{}) error {
	return Should(db.db.UpdateMap(obj))
}

// ShouldInsert ...
func (db Database) ShouldInsert(objs ...core.IInsert) error {
	return Should(db.db.Insert(objs...))
}

// ShouldDelete ...
func (db Database) ShouldDelete(obj core.ITableName) error {
	return Should(db.db.Delete(obj))
}

// ShouldGet ...
func (q Query) ShouldGet(obj core.IGet, preds ...interface{}) error {
	return ShouldGet(q.query.Get(obj, preds...))
}

// ShouldUpdate ...
func (q Query) ShouldUpdate(obj ...core.IUpdate) error {
	return Should(q.query.Update(obj...))
}

// ShouldUpdateMap ...
func (q Query) ShouldUpdateMap(obj map[string]interface{}) error {
	return Should(q.query.UpdateMap(obj))
}

// ShouldInsert ...
func (q Query) ShouldInsert(objs ...core.IInsert) error {
	return Should(q.query.Insert(objs...))
}

// ShouldDelete ...
func (q Query) ShouldDelete(obj core.ITableName) error {
	return Should(q.query.Delete(obj))
}

// ShouldGet ...
func (tx tx) ShouldGet(obj core.IGet, preds ...interface{}) error {
	return ShouldGet(tx.tx.Get(obj, preds...))
}

// ShouldUpdate ...
func (tx tx) ShouldUpdate(obj ...core.IUpdate) error {
	return Should(tx.tx.Update(obj...))
}

// ShouldUpdateMap ...
func (tx tx) ShouldUpdateMap(obj map[string]interface{}) error {
	return Should(tx.tx.UpdateMap(obj))
}

// ShouldInsert ...
func (tx tx) ShouldInsert(objs ...core.IInsert) error {
	return Should(tx.tx.Insert(objs...))
}

// ShouldDelete ...
func (tx tx) ShouldDelete(obj core.ITableName) error {
	return Should(tx.tx.Delete(obj))
}

func Should(n int64, err error) error {
	if err != nil {
		return err
	}
	if n == 0 {
		return cm.Error(cm.NotFound, "", nil)
	}
	return nil
}

func ShouldGet(has bool, err error) error {
	if err != nil {
		return err
	}
	if !has {
		return cm.Error(cm.NotFound, "", nil)
	}
	return nil
}

// Build ...
func (q Query) Build() (string, []interface{}, error) {
	return q.query.Build()
}

// BuildCount ...
func (q Query) BuildCount(obj core.ITableName) (string, []interface{}, error) {
	return q.query.BuildCount(obj)
}

// BuildGet ...
func (q Query) BuildGet(obj core.IGet) (string, []interface{}, error) {
	return q.query.BuildGet(obj)
}

// BuildFind ...
func (q Query) BuildFind(objs core.IFind) (string, []interface{}, error) {
	return q.query.BuildFind(objs)
}

func inOrEqIDs(q sq.CommonQuery, column string, ids []dot.ID) Query {
	if len(ids) == 1 {
		return Query{q.Where(column+`=?`, ids[0])}
	}
	return Query{q.In(column, ids)}
}
