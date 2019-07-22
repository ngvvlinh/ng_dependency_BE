package authkey

import (
	"strings"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/backend/pkg/common/validate"
	"etop.vn/common/l"
)

var ll = l.New()

type KeyType int

const (
	TypeAPIKey KeyType = iota + 1
	TypePartnerShopKey
)

const DefaultKeyLength = 64 // chars

type KeyInfo struct {
	Type      KeyType
	AccountID int64
	Subkey    string
}

// auth_key can have the following format
//
//     123456789:<key>       : API Key
//     shop123456789:<key>   : Partner-Shop Key
//
func ValidateAuthKey(key string) (info KeyInfo, ok bool) {
	parts := strings.Split(key, ":")
	if len(parts) != 2 {
		return
	}
	if len(parts[1]) == 0 {
		return
	}

	tag := parts[0]
	info.Subkey = parts[1]
	switch {
	case strings.HasPrefix(tag, "shop"):
		info.Type = TypePartnerShopKey
		info.AccountID, ok = validate.ParseInt64ID(tag[len("shop"):])

	default:
		info.Type = TypeAPIKey
		info.AccountID, ok = validate.ParseInt64ID(tag)
	}
	return
}

func ValidateAuthKeyWithType(typ KeyType, key string) (KeyInfo, bool) {
	info, ok := ValidateAuthKey(key)
	return info, ok && info.Type == typ
}

func GenerateAuthKey(typ KeyType, accountID int64) string {
	if accountID == 0 {
		ll.Panic("Invalid id")
	}

	subkey := gencode.GenerateCode(gencode.Alphabet54, DefaultKeyLength)
	switch typ {
	case TypeAPIKey:
		return cm.IDToDec(accountID) + ":" + subkey

	case TypePartnerShopKey:
		return "shop" + cm.IDToDec(accountID) + ":" + subkey

	default:
		ll.Panic("Invalid key type")
		return ""
	}
}
