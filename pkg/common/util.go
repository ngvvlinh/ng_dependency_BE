package cm

import (
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/dgrijalva/jwt-go"

	"o.o/capi/dot"
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

func CoalesceFloat(fs ...float64) float64 {
	for _, f := range fs {
		if f != 0 {
			return f
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

func CoalesceTime(times ...time.Time) time.Time {
	for _, _time := range times {
		if !_time.IsZero() {
			return _time
		}
	}
	return time.Time{}
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

var llHigh = ll.WithChannel("high")

// RecoverAndLog captures panic in goroutine, prevents the process from crashing
// and writes the error to logger. It also receives an optional error pointer
// for storing the recovered error.
//
// Note that the goroutine is still stopped. You still have to check the error
// and fix the real bug. Usage:
//
//     go func() { defer cm.RecoverAndLog(); doSomething() }()
//
//     go func() (_err error) {
//         defer cm.RecoverAndLog(&_err);
//         doSomething()
//     }()
func RecoverAndLog(errs ...*error) {
	r := recover()
	if r == nil {
		return
	}
	llHigh.SendMessagef("???? [panic+stopped] @thangtran268 %v\n%s", r, debug.Stack())
	for _, err := range errs {
		if err != nil {
			*err = Errorf(Internal, nil, "panic+stopped: %v", err)
		}
	}
}
