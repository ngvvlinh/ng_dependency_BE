package sqlstore

import (
	"context"
	"time"

	"etop.vn/api/main/connectioning"
	"etop.vn/api/top/types/etc/status3"
	"etop.vn/backend/com/main/connectioning/model"
	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/sql/cmsql"
	"etop.vn/backend/pkg/common/sql/sq"
	"etop.vn/backend/pkg/common/sql/sqlstore"
	"etop.vn/capi/dot"
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
	s.preds = append(s.preds, s.ft.ByIsGlobal(isGlobal))
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

func (s *ShopConnectionStore) UpdateShopConnectionToken(args *connectioning.UpdateShopConnectionExternalDataArgs) (*connectioning.ShopConnection, error) {
	var externalData model.ShopConnectionExternalData
	if err := scheme.Convert(args.ExternalData, &externalData); err != nil {
		return nil, err
	}
	update := &model.ShopConnection{
		Token:          args.Token,
		TokenExpiresAt: args.TokenExpiresAt,
		ExternalData:   &externalData,
	}
	query := s.query().Where(s.ft.ByConnectionID(args.ConnectionID))
	if args.ShopID != 0 {
		query = query.Where(s.ft.ByShopID(args.ShopID))
	} else {
		query = query.Where(s.ft.ByIsGlobal(true))
	}
	if err := query.ShouldUpdate(update); err != nil {
		return nil, err
	}
	return s.OptionalShopID(args.ShopID).ConnectionID(args.ConnectionID).GetShopConnection()
}

func (s *ShopConnectionStore) ConfirmShopConnection(shopID dot.ID, connID dot.ID) (updated int, err error) {
	if err := s.query().Table("shop_connection").Where(s.ft.ByShopID(shopID)).Where(s.ft.ByConnectionID(connID)).ShouldUpdateMap(map[string]interface{}{
		"status": status3.P,
	}); err != nil {
		return 0, err
	}
	return 1, nil
}
