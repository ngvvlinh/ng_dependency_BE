package meta

import (
	"bytes"
	"errors"

	uuid "github.com/satori/go.uuid"
)

var zeroUUID = uuid.UUID{}

func (id UUID) IsZero() bool {
	return len(id.Data) == 0 || bytes.Equal(id.Data, zeroUUID[:])
}

func (id UUID) MarshalJSON() ([]byte, error) {
	if id.IsZero() {
		return []byte("null"), nil
	}
	u, err := uuid.FromBytes(id.Data)
	if err != nil {
		return nil, err
	}
	return []byte(`"` + u.String() + `"`), nil
}

func (id *UUID) UnmarshalJSON(src []byte) error {
	if len(src) == 0 ||
		len(src) == 2 && string(src) == `""` ||
		len(src) == 4 && string(src) == "null" {
		id.Data = nil
		return nil
	}
	if len(src) >= 2 && src[0] == '"' && src[len(src)-1] == '"' {
		var u uuid.UUID
		if err := u.UnmarshalText(src[1 : len(src)-1]); err != nil {
			return err
		}
		if u == uuid.Nil {
			id.Data = nil
		} else {
			id.Data = u[:]
		}
		return nil
	}
	return errors.New("invalid uuid format")
}
