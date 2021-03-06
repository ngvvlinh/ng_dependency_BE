// +build !generator

// Code generated by generator sqlsel. DO NOT EDIT.

package model

import (
	"database/sql"

	core "o.o/backend/pkg/common/sql/sq/core"
)

type FbExternalPostFtTotalComments []*FbExternalPostFtTotalComment

func (m *FbExternalPostFtTotalComment) SQLTableName() string  { return "" }
func (m *FbExternalPostFtTotalComments) SQLTableName() string { return "" }

func (m *FbExternalPostFtTotalComment) SQLScanArgs(opts core.Opts) []interface{} {
	return []interface{}{
		(*core.String)(&m.ExternalPostID),
		(*core.Int)(&m.Count),
	}
}

func (m *FbExternalPostFtTotalComment) SQLScan(opts core.Opts, row *sql.Row) error {
	return row.Scan(m.SQLScanArgs(opts)...)
}

func (ms *FbExternalPostFtTotalComments) SQLScan(opts core.Opts, rows *sql.Rows) error {
	res := make(FbExternalPostFtTotalComments, 0, 128)
	for rows.Next() {
		m := new(FbExternalPostFtTotalComment)
		args := m.SQLScanArgs(opts)
		if err := rows.Scan(args...); err != nil {
			return err
		}
		res = append(res, m)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	*ms = res
	return nil
}

func (_ *FbExternalPostFtTotalComment) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT external_post_id, count(id)`)
	return nil
}

func (_ *FbExternalPostFtTotalComments) SQLSelect(w core.SQLWriter) error {
	w.WriteRawString(`SELECT external_post_id, count(id)`)
	return nil
}
