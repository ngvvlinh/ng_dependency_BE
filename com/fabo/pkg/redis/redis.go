package redis

import (
	"encoding/json"
	"fmt"
	"sync"

	"o.o/api/fabo/fbmessaging"
	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	"o.o/backend/pkg/common/redis"
	"o.o/common/l"
)

var ll = l.New()

const (
	PrefixPSID                  = "psid"
	PrefixExternalConversation  = "external_conversation"
	VersionExternalConversation = "v1.6"
	PrefixProfilePSID           = "profile_psid"
	PrefixLockCallAPI           = "lock_call_api"

	page      = "page"
	messenger = "messenger"
)

type FaboRedis struct {
	redisStore redis.Store
	mu         sync.Mutex
	version    string
}

func NewFaboRedis(redisStore redis.Store) *FaboRedis {
	return &FaboRedis{
		redisStore: redisStore,
		version:    "1.0",
	}
}

func (r *FaboRedis) LoadProfilePSID(pageID, PSID string) (*fbclientmodel.Profile, error) {
	profilePSIDstr, err := r.redisStore.GetString(r.GenerateProfilePSIDKey(pageID, PSID))
	if err != nil {
		return nil, err
	}

	var profile fbclientmodel.Profile

	if err := json.Unmarshal([]byte(profilePSIDstr), &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *FaboRedis) SaveProfilePSID(pageID, PSID string, profile *fbclientmodel.Profile) error {
	profileBytes, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	return r.redisStore.SetString(r.GenerateProfilePSIDKey(pageID, PSID), string(profileBytes))
}

func (r *FaboRedis) GenerateProfilePSIDKey(externalPageID, PSID string) string {
	return fmt.Sprintf("%s:%s:%s_%s", PrefixProfilePSID, r.version, externalPageID, PSID)
}

func (r *FaboRedis) SaveExternalConversation(
	externalPageID, externalUserID string,
	fbExternalConversation fbmessaging.FbExternalConversation,
) error {
	key := r.GenerateExternalConversationKey(externalPageID, externalUserID)

	return r.redisStore.Set(key, fbExternalConversation)
}

func (r *FaboRedis) LoadExternalConversation(externalPageID, externalUserID string) (*fbmessaging.FbExternalConversation, error) {
	key := r.GenerateExternalConversationKey(externalPageID, externalUserID)

	var fbExternalConversation fbmessaging.FbExternalConversation

	err := r.redisStore.Get(key, &fbExternalConversation)
	switch err {
	case redis.ErrNil:
		return nil, nil
	case nil:
		return &fbExternalConversation, nil
	default:
		return nil, err
	}
}

func (r *FaboRedis) ClearExternalConversations(externalPageIDs, externalUserIDs []string) error {
	var keys []string

	for idx, externalPageID := range externalPageIDs {
		keys = append(keys, r.GenerateExternalConversationKey(externalPageID, externalUserIDs[idx]))
	}

	return r.redisStore.Del(keys...)
}

func (r *FaboRedis) GenerateExternalConversationKey(externalPageID, externalUserID string) string {
	return fmt.Sprintf("%s:%s:%s_%s", PrefixExternalConversation, VersionExternalConversation, externalPageID, externalUserID)
}

// in minutes
func (r *FaboRedis) LockCallAPI(externalPageID string, TTL int) error {
	if r.IsLockCallAPI(externalPageID) {
		return nil
	}

	ll.SendMessagef("lock call apis with page (%s)", externalPageID)
	key := r.generateKeyLockCallAPI(externalPageID, "")
	return r.redisStore.SetStringWithTTL(key, externalPageID, TTL*60)
}

func (r *FaboRedis) IsLockCallAPI(externalPageID string) bool {
	key := r.generateKeyLockCallAPI(externalPageID, "")
	return r.redisStore.IsExist(key)
}

func (r *FaboRedis) LockCallAPIPage(externalPageID string, TTL int) error {
	if r.IsLockCallAPIPage(externalPageID) {
		return nil
	}

	ll.SendMessagef("lock call apis (page) with page (%s)", externalPageID)
	key := r.generateKeyLockCallAPI(externalPageID, page)
	return r.redisStore.SetStringWithTTL(key, externalPageID, TTL*60)
}

func (r *FaboRedis) IsLockCallAPIPage(externalPageID string) bool {
	key := r.generateKeyLockCallAPI(externalPageID, page)
	return r.redisStore.IsExist(key)
}

func (r *FaboRedis) LockCallAPIMessenger(externalPageID string, TTL int) error {
	if r.IsLockCallAPIMessenger(externalPageID) {
		return nil
	}

	ll.SendMessagef("lock call apis (messenger) with page (%s)", externalPageID)
	key := r.generateKeyLockCallAPI(externalPageID, messenger)
	return r.redisStore.SetStringWithTTL(key, externalPageID, TTL*60)
}

func (r *FaboRedis) IsLockCallAPIMessenger(externalPageID string) bool {
	key := r.generateKeyLockCallAPI(externalPageID, messenger)
	return r.redisStore.IsExist(key)
}

func (r *FaboRedis) generateKeyLockCallAPI(externalPageID, typ string) string {
	if typ == "" {
		return fmt.Sprintf("%s:%s:%s", PrefixLockCallAPI, r.version, externalPageID)
	}
	return fmt.Sprintf("%s:%s:%s:%s", PrefixLockCallAPI, r.version, typ, externalPageID)
}

func (r *FaboRedis) SetKey(key string, val interface{}) error {
	return r.redisStore.Set(key, val)
}

func (r *FaboRedis) SetWithTTL(key string, val interface{}, ttl int) error {
	return r.redisStore.SetWithTTL(key, val, ttl)
}

func (r *FaboRedis) IsExist(key string) bool {
	return r.redisStore.IsExist(key)
}

func (r *FaboRedis) DelKey(key string) error {
	return r.redisStore.Del(key)
}
