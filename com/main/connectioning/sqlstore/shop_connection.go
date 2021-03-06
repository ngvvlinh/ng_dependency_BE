package sqlstore

import (
	"context"
	"time"

	"o.o/api/main/connectioning"
	"o.o/api/top/types/etc/status3"
	"o.o/backend/com/main/connectioning/model"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/backend/pkg/common/sql/sq"
	"o.o/backend/pkg/common/sql/sqlstore"
	"o.o/capi/dot"
)

type ShopConnectionStore struct {
	ft    ShopConnectionFilters
	query func() cmsql.QueryInterface
	preds []interface{}

	includeDeleted sqlstore.IncludeDeleted
}

type ShopConnectionStoreFactory func(ctx context.Context) *ShopConnectionStore

func NewShopConnectionStore(db *cmsql.Database) ShopConnectionStoreFactory {
	return func(ctx context.Context) *ShopConnectionStore {
		return &ShopConnectionStore{
			query: func() cmsql.QueryInterface {
				return cmsql.GetTxOrNewQuery(ctx, db)
			},
		}

	}
}

func (s *ShopConnectionStore) Clone() *ShopConnectionStore {
	return &ShopConnectionStore{
		query:          s.query,
		preds:          s.preds,
		includeDeleted: s.includeDeleted,
	}
}

func (s *ShopConnectionStore) ShopID(shopID dot.ID) *ShopConnectionStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID))
	return s
}

func (s *ShopConnectionStore) OwnerID(id dot.ID) *ShopConnectionStore {
	s.preds = append(s.preds, s.ft.ByOwnerID(id))
	return s
}

func (s *ShopConnectionStore) OptionalOwnerID(ownerID dot.ID) *ShopConnectionStore {
	s.preds = append(s.preds, s.ft.ByOwnerID(ownerID).Optional())
	return s
}

func (s *ShopConnectionStore) OptionalShopID(shopID dot.ID) *ShopConnectionStore {
	s.preds = append(s.preds, s.ft.ByShopID(shopID).Optional())
	return s
}

func (s *ShopConnectionStore) ConnectionID(connectionID dot.ID) *ShopConnectionStore {
	s.preds = append(s.preds, s.ft.ByConnectionID(connectionID))
	return s
}

func (s *ShopConnectionStore) ConnectionIDs(connectionIDs ...dot.ID) *ShopConnectionStore {
	s.preds = append(s.preds, sq.In("connection_id", connectionIDs))
	return s
}

func (s *ShopConnectionStore) IsGlobal(isGlobal bool) *ShopConnectionStore {
	s.preds = append(s.preds, s.ft.ByIsGlobalPtr(&isGlobal))
	return s
}

func (s *ShopConnectionStore) SoftDelete() (int, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())
	_deleted, err := query.Table("shop_connection").UpdateMap(map[string]interface{}{
		"deleted_at": time.Now(),
	})
	return _deleted, err
}

func (s *ShopConnectionStore) GetShopConnectionDB() (*model.ShopConnection, error) {
	query := s.query().Where(s.preds)
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	var shopConn model.ShopConnection
	err := query.ShouldGet(&shopConn)
	return &shopConn, err
}

func (s *ShopConnectionStore) GetShopConnection() (*connectioning.ShopConnection, error) {
	shopConnDB, err := s.GetShopConnectionDB()
	if err != nil {
		return nil, err
	}
	var res connectioning.ShopConnection
	if err := scheme.Convert(shopConnDB, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *ShopConnectionStore) ListShopConnectionsDB() (res []*model.ShopConnection, err error) {
	query := s.query().Where(s.preds).Where(s.ft.ByStatus(status3.P))
	query = s.includeDeleted.Check(query, s.ft.NotDeleted())

	err = query.Find((*model.ShopConnections)(&res))
	return
}

func (s *ShopConnectionStore) ListShopConnections() (res []*connectioning.ShopConnection, _ error) {
	shopConnsDB, err := s.ListShopConnectionsDB()
	if err != nil {
		return nil, err
	}
	if err := scheme.Convert(shopConnsDB, &res); err != nil {
		return nil, err
	}
	return
}

func (s *ShopConnectionStore) CreateShopConnection(shopConn *connectioning.ShopConnection) (*connectioning.ShopConnection, error) {
	sqlstore.MustNoPreds(s.preds)
	var shopConnDB model.ShopConnection
	if err := scheme.Convert(shopConn, &shopConnDB); err != nil {
		return nil, err
	}
	if err := s.query().ShouldInsert(&shopConnDB); err != nil {
		return nil, err
	}
	return s.OptionalShopID(shopConn.ShopID).ConnectionID(shopConn.ConnectionID).GetShopConnection()
}

func (s *ShopConnectionStore) CreateShopConnectionDB(shopConn *model.ShopConnection) error {
	if shopConn.ConnectionID == 0 {
		return cm.Errorf(cm.InvalidArgument, nil, "Missing ConnectionID")
	}
	if err := s.query().ShouldInsert(shopConn); err != nil {
		return err
	}
	return nil
}

func (s *ShopConnectionStore) UpdateShopConnection(shopConn *connectioning.ShopConnection) error {
	var shopConnDB model.ShopConnection
	if err := scheme.Convert(shopConn, &shopConnDB); err != nil {
		return err
	}
	query := s.query().Where(s.preds)
	return query.ShouldUpdate(&shopConnDB)
}

func (s *ShopConnectionStore) UpdateShopConnectionLastSyncAt(args *connectioning.UpdateShopConnectionLastSyncAtArgs) (*connectioning.ShopConnection, error) {
	update := &model.ShopConnection{
		LastSyncAt: args.LastSyncAt,
	}
	query := s.query().Where(s.ft.ByConnectionID(args.ConnectionID))
	if args.ShopID != 0 {
		query = query.Where(s.ft.ByShopID(args.ShopID))
	} else if args.OwnerID != 0 {
		query = query.Where(s.ft.ByOwnerID(args.OwnerID))
	} else {
		query = query.Where(s.ft.ByIsGlobal(true))
	}
	if err := query.ShouldUpdate(update); err != nil {
		return nil, err
	}
	return s.OptionalShopID(args.ShopID).OptionalOwnerID(args.OwnerID).ConnectionID(args.ConnectionID).GetShopConnection()
}

func (s *ShopConnectionStore) ConfirmShopConnection() (updated int, _ error) {
	if len(s.preds) == 0 {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	update := &model.ShopConnection{
		Status: status3.P,
	}
	return s.query().Where(s.preds).Update(update)
}

func (s *ShopConnectionStore) DisableShopConnection() (updated int, _ error) {
	if len(s.preds) == 0 {
		return 0, cm.Errorf(cm.FailedPrecondition, nil, "must provide preds")
	}
	update := &model.ShopConnection{
		Status: status3.N,
	}
	return s.query().Where(s.preds).Update(update)
}
