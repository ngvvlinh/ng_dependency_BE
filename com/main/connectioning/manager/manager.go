package manager

import (
	"context"
	"encoding/json"
	"fmt"

	"o.o/api/main/connectioning"
	cm "o.o/backend/pkg/common"
	"o.o/backend/pkg/common/cipherx"
	"o.o/backend/pkg/common/redis"
	"o.o/capi/dot"
	"o.o/common/l"
)

const (
	DefaultTTl     = 2 * 60 * 60
	SecretKey      = "connectionsecretkey"
	VersionCaching = "1.1"
)

var ll = l.New()

type ConnectionManager struct {
	redisStore   redis.Store
	connectionQS connectioning.QueryBus
	cipherx      *cipherx.Cipherx
}

func NewConnectionManager(redisStore redis.Store, connectionQS connectioning.QueryBus) *ConnectionManager {
	_cipherx, _ := cipherx.NewCipherx(SecretKey)
	return &ConnectionManager{
		redisStore:   redisStore,
		connectionQS: connectionQS,
		cipherx:      _cipherx,
	}
}

func (m *ConnectionManager) GetConnectionByID(ctx context.Context, connID dot.ID) (*connectioning.Connection, error) {
	connKey := GetRedisConnectionKeyByID(connID)
	var connection connectioning.Connection
	err := m.loadRedis(connKey, &connection)
	if err != nil {
		query := &connectioning.GetConnectionByIDQuery{
			ID: connID,
		}
		if err := m.connectionQS.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).Wrap(cm.NotFound, "Connection not found").Throw()
		}
		connection = *query.Result
		connKeyCode := getRedisConnectionKeyByCode(connection.Code)
		m.setRedis(connKey, connection)
		m.setRedis(connKeyCode, connection)
	}
	return &connection, nil
}

func (m *ConnectionManager) GetConnectionByCode(ctx context.Context, connCode string) (*connectioning.Connection, error) {
	connKey := getRedisConnectionKeyByCode(connCode)
	var connection connectioning.Connection
	err := m.loadRedis(connKey, &connection)
	if err != nil {
		query := &connectioning.GetConnectionByCodeQuery{
			Code: connCode,
		}
		if err := m.connectionQS.Dispatch(ctx, query); err != nil {
			return nil, cm.MapError(err).Wrap(cm.NotFound, "Connection not found").Throw()
		}
		connection = *query.Result
		connKeyID := GetRedisConnectionKeyByID(connection.ID)
		m.setRedis(connKey, connection)
		m.setRedis(connKeyID, connection)
	}
	return &connection, nil
}

func (m *ConnectionManager) GetShopConnection(ctx context.Context, connID dot.ID, shopID dot.ID) (*connectioning.ShopConnection, error) {
	shopConnKey := GetRedisShopConnectionKey(connID, shopID)
	var shopConnection connectioning.ShopConnection
	err := m.loadRedis(shopConnKey, &shopConnection)
	if err == nil {
		return &shopConnection, nil
	}
	query2 := &connectioning.GetShopConnectionByIDQuery{
		ConnectionID: connID,
		ShopID:       shopID,
	}
	if err := m.connectionQS.Dispatch(ctx, query2); err != nil {
		return nil, err
	}
	shopConnection = *query2.Result
	m.setRedis(shopConnKey, shopConnection)
	return &shopConnection, nil
}

func GetRedisShopConnectionKey(connID dot.ID, shopID dot.ID) string {
	return fmt.Sprintf("shopConn:%v:%v%v", VersionCaching, shopID.String(), connID.String())
}

func GetRedisConnectionKeyByID(connID dot.ID) string {
	return fmt.Sprintf("conn:id:%v:%v", VersionCaching, connID.String())
}

func getRedisConnectionKeyByCode(code string) string {
	return fmt.Sprintf("conn:code:%v:%v", VersionCaching, code)
}

func (m *ConnectionManager) loadRedis(key string, v interface{}) error {
	if m.redisStore == nil {
		return cm.Errorf(cm.Internal, nil, "Redis service nil")
	}
	value, err := m.redisStore.GetString(key)
	if err != nil {
		return err
	}

	data, err := m.cipherx.Decrypt([]byte(value))
	if err != nil {
		ll.Error("Fail to decrypt from redis", l.Error(err))
		return err
	}

	if err := json.Unmarshal(data, &v); err != nil {
		ll.Error("Fail to unmarshal from redis", l.Error(err))
		return err
	}
	return nil
}

func (m *ConnectionManager) setRedis(key string, data interface{}) {
	if m.redisStore == nil {
		return
	}
	xData, err := json.Marshal(data)
	if err != nil {
		return
	}
	dataEncrypt, err := m.cipherx.Encrypt(xData)
	if err != nil {
		return
	}
	value := string(dataEncrypt)
	if err := m.redisStore.SetStringWithTTL(key, value, DefaultTTl); err != nil {
		ll.Error("Can not store to redis", l.Error(err))
	}
	return
}
