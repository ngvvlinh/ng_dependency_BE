package cm

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/dgrijalva/jwt-go"

	"etop.vn/capi/dot"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

const SaltSize = 16

func PNonZeroString(s string) dot.NullString {
	if s == "" {
		return dot.NullString{}
	}
	return dot.String(s)
}

func BoolDefault(b *bool, def bool) bool {
	if b == nil {
		return def
	}
	return *b
}

func JSON(v interface{}) []byte {
	data, err := jsonx.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func JSONString(v interface{}) string {
	data, err := jsonx.Marshal(v)
	if err != nil {
		panic(err)
	}
	return UnsafeBytesToString(data)
}

func Coalesce(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func CoalesceStrings(sss ...[]string) []string {
	for _, ss := range sss {
		if len(ss) != 0 {
			return ss
		}
	}
	return nil
}

func CoalesceInt32(is ...int32) int32 {
	for _, i := range is {
		if i != 0 {
			return i
		}
	}
	return 0
}

func CoalesceInt(is ...int) int {
	for _, i := range is {
		if i != 0 {
			return i
		}
	}
	return 0
}

func CoalesceInt64(is ...int64) int64 {
	for _, i := range is {
		if i != 0 {
			return i
		}
	}
	return 0
}

func StringsContain(ss []string, s string) bool {
	for _, item := range ss {
		if item == s {
			return true
		}
	}
	return false
}

func IDsContain(list []dot.ID, i dot.ID) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}

func IntsContain(list []int, i int) bool {
	for _, v := range list {
		if v == i {
			return true
		}
	}
	return false
}

func ToJSON(v interface{}) []byte {
	data, err := jsonx.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func GetJWTExpires(tokenStr string) time.Time {
	var claims jwt.StandardClaims
	jwt.ParseWithClaims(tokenStr, &claims, nil)
	if claims.ExpiresAt != 0 {
		return time.Unix(claims.ExpiresAt, 0)
	}
	return time.Time{}
}

func hexa(data []byte) string {
	return hex.EncodeToString(data)
}

func dehexa(s string) []byte {
	b, _ := hex.DecodeString(s)
	return b
}

func EncodePassword(password string) string {
	return hexa(saltedHashPassword([]byte(password)))
}

func DecodePassword(hashpw string) string {
	return string(dehexa(hashpw))
}

func saltedHashPassword(secret []byte) []byte {
	buf := make([]byte, SaltSize, SaltSize+sha1.Size)
	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		ll.Panic("Unable to read from rand.Reader", l.Error(err))
		panic(err)
	}

	h := sha1.New()
	_, err = h.Write(buf)
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	_, err = h.Write(secret)
	if err != nil {
		ll.Error("Write to buffer", l.Error(err))
	}

	return h.Sum(buf)
}

func CountConds(conds ...bool) int {
	count := 0
	for _, cond := range conds {
		if cond {
			count++
		}
	}
	return count
}

func UnsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func SortStrings(a []string) []string {
	sort.Strings(a)
	return a
}

func GetFormValue(ss []string) string {
	if ss == nil {
		return ""
	}
	return ss[0]
}

func URL(baseUrl, path string) string {
	return baseUrl + path
}

func Abs(num int) int {
	if num >= 0 {
		return num
	}
	return -num
}

func FormatCurrency(num int) string {
	sign := ""
	if num < 0 {
		sign += "-"
	}
	num = Abs(num)
	str := strconv.Itoa(num)
	var res []string
	for {
		if len(str) <= 0 {
			break
		}
		index := len(str) - 3
		if index < 0 {
			res = append([]string{str}, res...)
			break
		}
		res = append([]string{str[index:]}, res...)
		str = str[:index]
	}
	return sign + strings.Join(res, ".")
}

func ConvertStructToMapStringString(data interface{}) map[string]string {
	_data, _ := jsonx.Marshal(data)
	var metaX map[string]interface{}
	_ = jsonx.Unmarshal(_data, &metaX)
	meta := make(map[string]string)
	for k, v := range metaX {
		meta[k] = fmt.Sprint(v)
	}
	return meta
}
