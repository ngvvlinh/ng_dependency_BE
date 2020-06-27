package dot

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

func (id ID) Value() (driver.Value, error) {
	if id == 0 {
		return nil, nil
	}
	return int64(id), nil
}

func (id *ID) Scan(src interface{}) error {
	var ni sql.NullInt64
	err := ni.Scan(src)
	*id = ID(ni.Int64)
	return err
}

func (id NullID) Value() (driver.Value, error) {
	if !id.Valid || id.ID == 0 {
		return nil, nil
	}
	return int64(id.ID), nil
}

func (id *NullID) Scan(src interface{}) error {
	var ni sql.NullInt64
	err := ni.Scan(src)
	id.Valid, id.ID = ni.Int64 != 0, ID(ni.Int64)
	return err
}

func (n NullBool) Value() (driver.Value, error) {
	return sql.NullBool(n).Value()
}

func (n *NullBool) Scan(src interface{}) error {
	return (*sql.NullBool)(n).Scan(src)
}

func (n NullString) Value() (driver.Value, error) {
	return sql.NullString(n).Value()
}

func (n *NullString) Scan(src interface{}) error {
	return (*sql.NullString)(n).Scan(src)
}

func (n NullInt32) Value() (driver.Value, error) {
	return sql.NullInt32(n).Value()
}

func (n *NullInt32) Scan(src interface{}) error {
	return (*sql.NullInt32)(n).Scan(src)
}

func (n NullInt64) Value() (driver.Value, error) {
	return sql.NullInt64(n).Value()
}

func (n *NullInt64) Scan(src interface{}) error {
	return (*sql.NullInt64)(n).Scan(src)
}

func (n NullInt) Value() (driver.Value, error) {
	_n := NullInt64{Int64: int64(n.Int), Valid: n.Valid}
	return sql.NullInt64(_n).Value()
}

func (n *NullInt) Scan(src interface{}) error {
	var _n NullInt64
	if err := _n.Scan(src); err != nil {
		*n = NullInt{}
		return err
	}
	n.Valid = _n.Valid
	n.Int = int(_n.Int64)
	return nil
}

func (n NullFloat64) Value() (driver.Value, error) {
	return sql.NullFloat64(n).Value()
}

func (n *NullFloat64) Scan(src interface{}) error {
	return (*sql.NullFloat64)(n).Scan(src)
}

func (t Time) Value() (driver.Value, error) {
	return time.Time(t), nil
}

func (t *Time) Scan(src interface{}) error {
	var _t sql.NullTime
	if err := _t.Scan(src); err != nil {
		*t = Time{}
		return err
	}
	if _t.Valid {
		*t = Time(_t.Time)
	} else {
		*t = Time{}
	}
	return nil
}
