package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/connection_type"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/connectioning/convert"
	"o.o/backend/com/main/connectioning/model"
	"o.o/backend/pkg/common/apifw/whitelabel/wl"
	"o.o/backend/pkg/common/conversion"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ConnectionStore struct {
	ft    ConnectionFilters
	query func() cmsql.QueryInterface
	preds []interface{}
	ctx   context.Context

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
			ctx: ctx,
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

func (s *ConnectionStore) OriginConnectionID(originConnectionID dot.ID) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByOriginConnectionID(originConnectionID))
	return s
}

func (s *ConnectionStore) Status(status status3.Status) *ConnectionStore {
	s.preds = append(s.preds, s.ft.ByStatus(status))
	return s
}

func (s *ConnectionStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
	_deleted, err := query.Table("connection").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *ConnectionStore) GetConnectionDB() (*model.Connection, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	query = s.ByWhiteLabelPartner(s.ctx, query)
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
	query = s.ByWhiteLabelPartner(s.ctx, query)

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
	connDB.WLPartnerID = wl.GetWLPartnerID(s.ctx)
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
	query := s.query().Where(s.ft.ByID(conn.ID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	if err := query.ShouldUpdate(&connDB); err != nil {
		return nil, err
	}
	return s.ID(conn.ID).GetConnection()
}

func (s *ConnectionStore) ConfirmConnection(connID dot.ID) (updated int, err error) {
	query := s.query().Table("connection").Where(s.ft.ByID(connID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	if err := query.ShouldUpdateMap(map[string]interface{}{
		"status": status3.P,
	}); err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *ConnectionStore) DisableConnection(connID dot.ID) (updated int, err error) {
	query := s.query().Table("connection").Where(s.ft.ByID(connID))
	query = s.ByWhiteLabelPartner(s.ctx, query)
	if err := query.ShouldUpdateMap(map[string]interface{}{
		"status": status3.Z,
	}); err != nil {
		return 0, err
	}
	return 1, nil
}

func (s *ConnectionStore) ByWhiteLabelPartner(ctx context.Context, query cmsql.Query) cmsql.Query {
	partner := wl.X(ctx)
	if partner.IsWhiteLabel() {
		return query.Where(s.ft.ByWLPartnerID(partner.ID))
	}
	return query.Where(s.ft.NotBelongWLPartner())
}
