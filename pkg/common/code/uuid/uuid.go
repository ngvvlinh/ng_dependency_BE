package uuid

import (
	"database/sql/driver"
	"errors"

	uuid "github.com/satori/go.uuid"
	"o.o/api/meta"
)

type UUID meta.UUID

func (id UUID) IsZero() bool {
	return meta.UUID(id).IsZero()
}

func (id *UUID) Scan(src interface{}) error {
	switch src := src.(type) {
	case nil:
		id.Data = nil
		return nil

	case string:
		_id, err := uuid.FromString(src)
		if err != nil {
			return err
		}
		if (_id == uuid.UUID{}) {
			id.Data = nil
		} else {
			id.Data = _id[:]
		}
		return nil

	default:
		return errors.New("invalid uuid src")
	}
}

func (id UUID) Value() (driver.Value, error) {
	if len(id.Data) == 0 {
		return nil, nil
	}
	u, err := uuid.FromBytes(id.Data)
	if err != nil {
		return nil, err
	}
	if (u == uuid.UUID{}) {
		return nil, nil
	}
	return u.String(), nil
}
