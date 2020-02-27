package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/types/etc/connection_type"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/connectioning/convert"
	"etop.vn/backend/com/main/connectioning/model"
	"etop.vn/backend/pkg/common/conversion"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
)

type ConnectionStore struct {
	ft    ConnectionFilters
	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type ConnectionStoreFactory func(ctx context.Context) *ConnectionStore

var scheme = conversion.Build(convert.RegisterConversions)

func NewConnectionStore(db *cmsql.Database) ConnectionStoreFactory {
	return func(ctx context.Context) *ConnectionStore {
		return &ConnectionStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}
	}
}

func (s *ConnectionStore) ID(id dot.ID) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByID(id))
	return s
}

func (s *ConnectionStore) Code(code string) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByCode(code))
	return s
}

func (s *ConnectionStore) PartnerID(partnerID dot.ID) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(partnerID))
	return s
}

func (s *ConnectionStore) OptionalPartnerID(partnerID dot.ID) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByPartnerID(partnerID).Optional())
	return s
}

func (s *ConnectionStore) OptionalConnectionMethod(method connection_type.ConnectionMethod) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByConnectionMethod(method).Optional())
	return s
}

func (s *ConnectionStore) OptionalConnectionType(_type connection_type.ConnectionType) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByConnectionType(_type).Optional())
	return s
}

func (s *ConnectionStore) OptionalConnectionProvider(provider connection_type.ConnectionProvider) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByConnectionProvider(provider).Optional())
	return s
}

func (s *ConnectionStore) Status(status status3.Status) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *ConnectionStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("connection").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *ConnectionStore) GetConnectionDB() (*model.Connection, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	var conn model.Connection
	err := query.ShouldGet(&conn)
	return &conn, err
}

func (s *ConnectionStore) GetConnection() (*connectioning.Connection, error) {
	connDB, err := s.GetConnectionDB()
	if err != nil {
		return nil, err
	}
	var res connectioning.Connection
	if err := scheme.Convert(connDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *ConnectionStore) ListConnectionsDB() (res []*model.Connection, err error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	err = query.Find((*model.Connections)(&res))
	return
}

func (s *ConnectionStore) ListConnections(status status3.NullStatus) (res []*connectioning.Connection, _ error) {
	if status.Valid {
		s = s.Status(status.Enum)
	}
	connsDB, err := s.ListConnectionsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(connsDB, &res); err != nil {
		return nil, err
	}
	return
}

func (s *ConnectionStore) CreateConnection(conn *connectioning.Connection) (*connectioning.Connection, error) {
	sqlstore.MustNoPreds(s.preds)
	var connDB model.Connection
	if err := scheme.Convert(conn, &connDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&connDB); err != nil {
		return nil, err
	}
	return s.ID(conn.ID).GetConnection()
}

func (s *ConnectionStore) UpdateConnection(conn *connectioning.Connection) (*connectioning.Connection, error) {
	sqlstore.MustNoPreds(s.preds)
	var connDB model.Connection
	if err := scheme.Convert(conn, &connDB); err != nil {
		return nil, err
	}
	if err := s.query().Where(s.ft.ByID(conn.ID)).ShouldUpdate(&connDB); err != nil {
		return nil, err
	}
	return s.ID(conn.ID).GetConnection()
}

func (s *ConnectionStore) ConfirmConnection(connID dot.ID) (updated int, err error) {
	if err := s.query().Table("connection").Where(s.ft.ByID(connID)).ShouldUpdateMap(map[string]interface{}{
		"status": status3.P,
	}); err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *ConnectionStore) DisableConnection(connID dot.ID) (updated int, err error) {
	if err := s.query().Table("connection").Where(s.ft.ByID(connID)).ShouldUpdateMap(map[string]interface{}{
		"status": status3.Z,
	}); err != nil {
		return 0, err
	}
	return 1, nil
}
