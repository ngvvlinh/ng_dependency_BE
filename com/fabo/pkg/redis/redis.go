package redis

import (
	"encoding/json"
	"fmt"
	"sync"

	fbclientmodel "o.o/backend/com/fabo/pkg/fbclient/model"
	"o.o/backend/pkg/common/redis"
)

const (
	PrefixPSID                 = "psid"
	PrefixExternalConversation = "external_conversation"
	PrefixProfilePSID          = "profile_psid"
)

type FaboRedis struct {
	redisStore redis.Store
	mu         sync.Mutex
}

func NewFaboRedis(redisStore redis.Store) *FaboRedis {
	return &FaboRedis{
		redisStore: redisStore,
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
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.redisStore.SetString(r.GenerateProfilePSIDKey(pageID, PSID), string(profileBytes))
}

func (r *FaboRedis) GenerateProfilePSIDKey(externalPageID, PSID string) string {
	return fmt.Sprintf("%s:%s_%s", PrefixProfilePSID, externalPageID, PSID)
}

func (r *FaboRedis) LoadPSID(pageID, PSID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.redisStore.GetString(r.GeneratePSIDKey(pageID, PSID))
}

func (r *FaboRedis) SavePSID(pageID, PSID, externalUserID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.redisStore.SetString(r.GeneratePSIDKey(pageID, PSID), externalUserID)
}

func (r *FaboRedis) GeneratePSIDKey(externalPageID, psid string) string {
	return fmt.Sprintf("%s:%s_%s", PrefixPSID, externalPageID, psid)
}

func (r *FaboRedis) LoadExternalConversationID(externalPageID, externalUserID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.redisStore.GetString(r.GenerateExternalConversationKey(externalPageID, externalUserID))
}

func (r *FaboRedis) SaveExternalConversationID(externalPageID, externalUserID, externalConversationID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.redisStore.SetString(r.GenerateExternalConversationKey(externalPageID, externalUserID), externalConversationID)
}

func (r *FaboRedis) GenerateExternalConversationKey(externalPageID, externalUserID string) string {
	return fmt.Sprintf("%s:%s_%s", PrefixExternalConversation, externalPageID, externalUserID)
}
