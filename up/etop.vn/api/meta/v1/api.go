package v1

import (
	"bytes"
	"errors"
	"time"

	"github.com/golang/protobuf/jsonpb"
	pbtypes "github.com/golang/protobuf/ptypes"
	pbtimestamp "github.com/golang/protobuf/ptypes/timestamp"
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

func (id UUID) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	return id.MarshalJSON()
}

func (id *UUID) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, src []byte) error {
	return id.UnmarshalJSON(src)
}

func (t *Timestamp) ToTime() time.Time {
	if t == nil {
		return time.Time{}
	}
	pbt := &pbtimestamp.Timestamp{
		Seconds: t.Seconds,
		Nanos:   t.Nanos,
	}
	result, _ := pbtypes.Timestamp(pbt)
	return result
}

func PbTime(t time.Time) *Timestamp {
	if t.Year() < 1990 {
		return nil
	}
	ts, _ := pbtypes.TimestampProto(t)
	if ts == nil {
		return nil
	}
	return &Timestamp{Seconds: ts.Seconds, Nanos: ts.Nanos}
}

func (t Timestamp) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	return []byte(t.ToTime().Format(time.RFC3339Nano)), nil
}

func (t *Timestamp) UnmarshalJSONPB(m *jsonpb.Unmarshaler, src []byte) error {
	var ts pbtimestamp.Timestamp
	err := m.Unmarshal(bytes.NewReader(src), &ts)
	if err != nil {
		return err
	}
	*t = Timestamp{Seconds: ts.Seconds, Nanos: ts.Nanos}
	return nil
}
