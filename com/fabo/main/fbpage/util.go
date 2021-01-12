package fbpage

import (
	"fmt"

	"o.o/api/fabo/fbpaging"
	"o.o/backend/pkg/common/cipherx"
	"o.o/backend/pkg/common/redis"
)

const (
	versionCached           = "v1.0.7"
	ttl                     = 10 * 60
	prefixKeyFbPage         = "fb_page"
	prefixKeyFbPageInternal = "fb_page_internal"
	secretKey               = "facebookpageinternal"
)

type FbPageUtil struct {
	rd      redis.Store
	cipherx *cipherx.Cipherx
}

func NewFbPageUtil(rd redis.Store) *FbPageUtil {
	_cipherx, _ := cipherx.NewCipherx(secretKey)
	return &FbPageUtil{
		rd:      rd,
		cipherx: _cipherx,
	}
}

func (s *FbPageUtil) generateKey(prefixKey, externalPageID string) string {
	return fmt.Sprintf("%s:%s:%s", prefixKey, versionCached, externalPageID)
}

func (s *FbPageUtil) SetFbPage(externalPageID string, fbPage fbpaging.FbExternalPage) error {
	key := s.generateKey(prefixKeyFbPage, externalPageID)

	return s.rd.SetWithTTL(key, fbPage, ttl)
}

func (s *FbPageUtil) GetFbPage(externalPageID string) (*fbpaging.FbExternalPage, error) {
	key := s.generateKey(prefixKeyFbPage, externalPageID)

	var fbPage fbpaging.FbExternalPage

	err := s.rd.Get(key, &fbPage)
	switch err {
	case redis.ErrNil:
		return nil, nil
	case nil:
		return &fbPage, nil
	default:
		return nil, err
	}
}

func (s *FbPageUtil) ClearFbPages(externalPageIDs ...string) error {
	var keys []string

	for _, externalPageID := range externalPageIDs {
		keys = append(keys, s.generateKey(prefixKeyFbPage, externalPageID))
	}

	return s.rd.Del(keys...)
}

func (s *FbPageUtil) ClearFbPage(externalPageID string) error {
	key := s.generateKey(prefixKeyFbPage, externalPageID)

	return s.rd.Del(key)
}

func (s *FbPageUtil) SetFbPageInternal(externalPageID string, fbPageInternal fbpaging.FbExternalPageInternal) error {
	key := s.generateKey(prefixKeyFbPageInternal, externalPageID)

	tokenEncrypted, err := s.cipherx.Encrypt([]byte(fbPageInternal.Token))
	if err != nil {
		return err
	}

	fbPageInternal.Token = string(tokenEncrypted)
	return s.rd.SetWithTTL(key, fbPageInternal, ttl)
}

func (s *FbPageUtil) GetFbPageInternal(externalPageID string) (*fbpaging.FbExternalPageInternal, error) {
	key := s.generateKey(prefixKeyFbPageInternal, externalPageID)

	var fbPageInternal fbpaging.FbExternalPageInternal

	err := s.rd.Get(key, &fbPageInternal)
	switch err {
	case redis.ErrNil:
		return nil, nil
	case nil:
		tokenDecrypted, err := s.cipherx.Decrypt([]byte(fbPageInternal.Token))
		if err != nil {
			return nil, err
		}

		fbPageInternal.Token = string(tokenDecrypted)
		return &fbPageInternal, nil
	default:
		return nil, err
	}
}

func (s *FbPageUtil) ClearFbPageInternal(externalPageID string) error {
	key := s.generateKey(prefixKeyFbPage, externalPageID)

	return s.rd.Del(key)
}

func (s *FbPageUtil) ClearFbPageInternals(externalPageInternalIDs ...string) error {
	var keys []string

	for _, externalPageInternalID := range externalPageInternalIDs {
		keys = append(keys, s.generateKey(prefixKeyFbPage, externalPageInternalID))
	}

	return s.rd.Del(keys...)
}
