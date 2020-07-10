package cm

import (
	"fmt"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/dgrijalva/jwt-go"

	"o.o/capi/dot"
	"o.o/common/jsonx"
)

func Coalesce(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}

func CoalesceInt(is ...int) int {
	for _, i := range is {
		if i != 0 {
			return i
		}
	}
	return 0
}

func CoalesceID(ids ...dot.ID) dot.ID {
	for _, id := range ids {
		if id != 0 {
			return id
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

func GetJWTExpires(tokenStr string) time.Time {
	var claims jwt.StandardClaims
	jwt.ParseWithClaims(tokenStr, &claims, nil)
	if claims.ExpiresAt != 0 {
		return time.Unix(claims.ExpiresAt, 0)
	}
	return time.Time{}
}

func UnsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func SortStrings(a []string) []string {
	sort.Strings(a)
	return a
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

var llHigh = ll.WithChannel("high")

// RecoverAndLog captures panic in goroutine, prevents the process from crashing
// and writes the error to logger. Note that the gorountine is still stopped.
// You still have to check the error and fix the real bug. Usage:
//
//     go func() { defer cm.RecoverAndLog(); doSomething() }()
//
func RecoverAndLog() {
	r := recover()
	if r != nil {
		llHigh.SendMessagef("🔥 [panic+stopped] @thangtran268 %v\n%s", r, debug.Stack())
	}
}
