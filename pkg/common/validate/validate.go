package validate

import (
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unsafe"

	"github.com/asaskevich/govalidator"
	"golang.org/x/text/unicode/norm"
)

const (
	phoneChars = "0123456789"
	vneseChars = "đĐ" +
		"àáạảãâầấậẩẫăằắặẳẵ" +
		"ÀÁẠẢÃÂẦẤẬẨẪĂẰẮẶẲẴ" +
		"èéẹẻẽêềếệểễ" +
		"ÈÉẸẺẼÊỀẾỆỂỄ" +
		"òóọỏõôồốộổỗơờớợởỡ" +
		"ÒÓỌỎÕÔỒỐỘỔỖƠỜỚỢỞỠ" +
		"ùúụủũưừứựửữ" +
		"ÙÚỤỦŨƯỪỨỰỬỮ" +
		"ìíịỉĩ" + "ỳýỵỷỹ" +
		"ÌÍỊỈĨ" + "ỲÝỴỶỸ"

	numChars   = "0123456789"
	upperChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars = "abcdefghijklmnopqrstuvwxyz"
	signChars  = ` .,/\"'_-+=@#%*()[]{}<>!?$&`
	nameChars  = signChars + numChars + upperChars + lowerChars + vneseChars
)

var (
	phoneRegexp    = regexp.MustCompile(`^0[0-9]{4,14}$`)
	phoneWhiteList = regexp.MustCompile(`[^0-9]+`)

	tagRegexp         = regexp.MustCompile(`^[\-\[\]/\\ .,"'_+=@#%*(){}<>!?&$:;|` + numChars + upperChars + lowerChars + vneseChars + `]{1,100}$`)
	nameRegexp        = regexp.MustCompile(`^[\-\[\]/\\ .,"'_+=@#%*(){}<>!?&$:;|` + numChars + upperChars + lowerChars + vneseChars + `]{2,200}$`)
	genericNameRegexp = regexp.MustCompile(`^[\-\[\]/\\ .,"'_+=@#%*(){}<>!?&$:;|` + numChars + upperChars + lowerChars + vneseChars + `]{2,400}$`)
	nameWhiteList     = regexp.MustCompile(`[^\-\[\]/\\ .,"'_+=@#%*(){}<>!?&$:;|` + numChars + upperChars + lowerChars + vneseChars + `]+`)
	idRegexp          = regexp.MustCompile(`^[0-9A-z]{1,100}$`)
	lowercaseIdRegexp = regexp.MustCompile(`^[0-9a-z_]{1,100}$`)

	// reject \, single and double quotes for preventing potential conflict with
	// some query languages
	externalCodeRegexp = regexp.MustCompile(`^[\-._$#:\w]{2,100}$`)

	spaceWhiteList = regexp.MustCompile(`\s\s+`)

	emailLocalRegexp = regexp.MustCompile(`[0-9a-z._-]{1,128}`)
	subdomainRegexp  = regexp.MustCompile(`[0-9A-z]{1,200}`)
	slugRegexp       = regexp.MustCompile(`^[0-9a-z]([0-9a-z-]{0,62}[0-9a-z])?$`)

	reHost = regexp.MustCompile(`^https?://[a-z0-9.]+$`)
	reTest = regexp.MustCompile(`(-[a-z0-9]+)?-test$`)

	reDomainFromURL = regexp.MustCompile(`^(https?://)?([a-z0-9.]+)`)

	vneseMap map[rune]byte

	// reject \, single and double quotes for preventing potential conflict with
	// some query languages
	specialChars = []bool{
		'(': true,
		')': true,
		'[': true,
		']': true,
		'{': true,
		'}': true,
		'<': true,
		'>': true,
		'/': true,

		'!': true,
		'@': true,
		'#': true,
		'$': true,
		'%': true,
		'^': true,
		'&': true,
		'*': true,
		'-': true,
		'_': true,
		'+': true,
		'=': true,

		'.': true,
		',': true,
		':': true,
		';': true,
		'?': true,
		'|': true,
	}
)

type NormalizedPhone string
type NormalizedEmail string

func (s NormalizedPhone) String() string { return string(s) }
func (s NormalizedEmail) String() string { return string(s) }

func init() {
	SetupDefault()

	vneseMap = make(map[rune]byte)
	initVneseMap("đĐ", 'd')
	initVneseMap("àáạảãâầấậẩẫăằắặẳẵ", 'a')
	initVneseMap("ÀÁẠẢÃÂẦẤẬẨẪĂẰẮẶẲẴ", 'a')
	initVneseMap("èéẹẻẽêềếệểễ", 'e')
	initVneseMap("ÈÉẸẺẼÊỀẾỆỂỄ", 'e')
	initVneseMap("òóọỏõôồốộổỗơờớợởỡ", 'o')
	initVneseMap("ÒÓỌỎÕÔỒỐỘỔỖƠỜỚỢỞỠ", 'o')
	initVneseMap("ùúụủũưừứựửữ", 'u')
	initVneseMap("ÙÚỤỦŨƯỪỨỰỬỮ", 'u')
	initVneseMap("ìíịỉĩÌÍỊỈĨ", 'i')
	initVneseMap("ỳýỵỷỹỲÝỴỶỸ", 'y')
}

func initVneseMap(s string, c byte) {
	for _, src := range s {
		vneseMap[src] = c
	}
}

func IsAcceptedSpecialChar(c rune) bool {
	return int(c) < len(specialChars) && specialChars[c]
}

// IsEmail is a simple function to distinct between email address and phone number.
// Must continue use ValidateEmail or ValidatePhone for more accurrate validateion.
func IsEmail(s string) bool {
	return strings.Contains(s, "@")
}

// SetupDefault ...
func SetupDefault() {
	govalidator.CustomTypeTagMap.Set("phone",
		func(v interface{}, ctx interface{}) bool {
			if s, ok := assertString(v); ok {
				return phoneRegexp.MatchString(s)
			}
			return false
		})

	govalidator.CustomTypeTagMap.Set("code",
		func(v interface{}, ctx interface{}) bool {
			if s, ok := assertString(v); ok {
				return Code(s)
			}
			return false
		})

	govalidator.CustomTypeTagMap.Set("name",
		func(v interface{}, ctx interface{}) bool {
			if s, ok := assertString(v); ok {
				if len(s) < 2 || len(s) > 200 {
					return false
				}
				return nameRegexp.MatchString(s)
			}
			return false
		})
}

func assertString(s interface{}) (string, bool) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.String {
		return "", false
	}
	return v.String(), true
}

// Check ...
func Check(v interface{}) error {
	_, err := govalidator.ValidateStruct(v)
	return err
}

// Code ...
func Code(s string) bool {
	if len(s) < 3 || len(s) > 64 {
		return false
	}
	for i, l := 0, len(s); i < l; i++ {
		// Only allow printable ASCII chars (from `!` to `~`)
		if s[i] < 33 || s[i] > 126 {
			return false
		}
	}
	return true
}

// ID ...
func ID(s string) bool {
	return idRegexp.MatchString(s)
}

func LowercaseID(s string) bool {
	return lowercaseIdRegexp.MatchString(s)
}

// ExternalCode ...
func ExternalCode(s string) bool {
	return externalCodeRegexp.MatchString(s)
}

// If the id starts with ~, we will strip it. Otherwise, verify it.
func NormalizeExternalCode(s string) string {
	if s == "" {
		return ""
	}
	if s[0] == '~' {
		return normalizeSearch(s[0:], "", false, false)
	}
	if !ExternalCode(s) {
		return ""
	}
	return s
}

func ParseInt64ID(s string) (int64, bool) {
	i, _ := strconv.ParseInt(s, 10, 64)

	// our id has 19 characters
	if i <= 1e18 {
		return 0, false
	}
	return i, true
}

// URLSlug ...
func URLSlug(s string) bool {
	return slugRegexp.MatchString(s)
}

// URL ...
func URL(s string) bool {
	return govalidator.IsURL(s)
}

// TrimTest ...
func TrimTest(s string) (string, string, bool) {
	matches := reTest.FindStringSubmatch(s)
	if len(matches) == 0 {
		return s, "", false
	}
	return s[:len(s)-len(matches[0])], matches[0], true
}

// NormalizePhone ...
func NormalizePhone(s string) (res NormalizedPhone, ok bool) {
	if strings.HasSuffix(s, "-test") {
		var suffix string
		s, suffix, _ = TrimTest(s)
		defer func() {
			if ok {
				res += NormalizedPhone(suffix)
			}
		}()
	}

	s = parseSinglePhoneNumber(s)

	// số điện thoại bàn có thể là 11 hoặc 12 số
	// chỉ kiểm tra là số điện thoại di động nếu có đầu số là 09
	return NormalizedPhone(s), len(s) >= 10 && len(s) <= 12
}

// MustNormalizePhone ...
func MustNormalizePhone(s string) string {
	return parseSinglePhoneNumber(s)
}

// NormalizeName ...
func NormalizeName(s string) (string, bool) {
	s = strings.TrimSpace(s)
	s = norm.NFC.String(s)
	s = WhiteList(s, nameWhiteList)
	s = TrimInnerSpace(s)

	if len(s) == 0 {
		return "", false
	}
	return s, nameRegexp.MatchString(s)
}

func NormalizeGenericName(s string) (string, bool) {
	s = strings.TrimSpace(s)
	s = norm.NFC.String(s)
	s = WhiteList(s, nameWhiteList)
	s = TrimInnerSpace(s)
	if len(s) == 0 {
		return "", false
	}
	return s, genericNameRegexp.MatchString(s)
}

// NormalizeTag ...
func NormalizeTag(s string) (string, bool) {
	s = strings.TrimSpace(s)
	s = norm.NFC.String(s)
	s = WhiteList(s, nameWhiteList)
	s = TrimInnerSpace(s)

	if len(s) == 0 {
		return "", false
	}
	return s, tagRegexp.MatchString(s)
}

// NormalizeSubdomain ...
func NormalizeSubdomain(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return "", false
	}
	return s, subdomainRegexp.MatchString(s)
}

// NormalizeEmail ...
func NormalizeEmail(s string) (res NormalizedEmail, ok bool) {
	if strings.HasSuffix(s, "-test") {
		var suffix string
		s, suffix, _ = TrimTest(s)
		defer func() {
			if ok {
				res += NormalizedEmail(suffix)
			}
		}()
	}

	s = strings.TrimSpace(s)
	s = strings.ToLower(s)
	ss := strings.Split(s, "@")
	if len(ss) != 2 {
		return "", false
	}
	if !emailLocalRegexp.MatchString(ss[0]) {
		return "", false
	}
	if ss[1] == "gmail.com" {
		ss[0] = strings.Replace(ss[0], ".", "", -1)
	}
	return NormalizedEmail(ss[0] + "@" + ss[1]), govalidator.IsEmail(s)
}

func NormalizeEmailOrPhone(s string) (email string, phone string, ok bool) {
	if IsEmail(s) {
		var emailNorm NormalizedEmail
		emailNorm, ok = NormalizeEmail(s)
		email = string(emailNorm)
	} else {
		var phoneNorm NormalizedPhone
		phoneNorm, ok = NormalizePhone(s)
		phone = string(phoneNorm)
	}
	return
}

// MustNormalizeName ...
func MustNormalizeName(s string) string {
	s, _ = NormalizeName(s)
	return s
}

// WhiteList remove characters that do not appear in the whitelist.
func WhiteList(s string, r *regexp.Regexp) string {
	return r.ReplaceAllString(s, "")
}

// TrimInnerSpace remove inner space characters
func TrimInnerSpace(s string) string {
	return spaceWhiteList.ReplaceAllString(s, " ")
}

// ParsePhoneInput ...
func ParsePhoneInput(inputPhone string) (string, string, bool) {
	count := 0
	for i := range inputPhone {
		if isNumberic(inputPhone[i]) {
			count++
		}
	}
	if count < 16 {
		s := parseSinglePhoneNumber(inputPhone)
		return s, "", s != ""
	}

	// Detect a single delimiter
	counts := make(map[byte]int)
	for i := range inputPhone {
		c := inputPhone[i]
		switch c {
		case ' ', '-', '_', '.', ',', ';', '|', '&', '\t':
			_, ok := counts[c]
			if !ok {
				// Store the index
				counts[c] = i

			} else {
				// Mark invalid
				counts[c] = -1
			}

		default:
			continue
		}
	}
	for _, index := range counts {
		if index >= 0 && index+1 < len(inputPhone) {
			phone1 := parseSinglePhoneNumber(inputPhone[:index])
			phone2 := parseSinglePhoneNumber(inputPhone[index+1:])
			return phone1, phone2, phone1 != "" && phone2 != ""
		}
	}
	return "", "", false
}

func parseSinglePhoneNumber(input string) string {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, "+84") {
		input = "0" + input[3:]
	} else if strings.HasPrefix(input, "84") {
		input = "0" + input[2:]
	}
	p := make([]byte, 0, 16)

loop:
	for i := range input {
		c := input[i]
		if c >= '0' && c <= '9' {
			p = append(p, c)
			continue
		}
		switch c {
		case ' ', '\t', '\r', '\n', '\v', ')', '-', '+', '.':

		// 01 234 5678 (16)
		case '(':
			if len(p) > 0 {
				break loop
			}

		// 01 234 5678 ext 16
		default:
			break loop
		}
	}

	// Remove prefix 0
	for len(p) > 0 && p[0] == '0' {
		p = p[1:]
	}

	// Too short
	if len(p) < 4 {
		return ""
	}

	// 1900-1234
	if len(p) <= 8 {
		return string(p)
	}
	return "0" + string(p)
}

func isNumberic(c byte) bool {
	return c >= '0' && c <= '9'
}

func NormalizeUnaccent(s string) string {
	return normalizeSearch(s, " ", false, true)
}

func NormalizeSearch(s string) string {
	return normalizeSearch(s, " ", true, true)
}

func NormalizeSearchSimple(s string) string {
	return normalizeSearchSimple(s, " ")
}

func NormalizeSlug(s string) string {
	return normalizeSearchSimple(s, "-")
}

func NormalizeUnderscore(s string) string {
	return normalizeSearchSimple(s, "_")
}

func NormalizeSearchQueryAnd(s string) string {
	return normalizeSearch(s, " & ", true, true)
}

func NormalizeSearchQueryOr(s string) string {
	return normalizeSearch(s, " | ", true, true)
}

// the old version, which only keep alphanumeric characters
func normalizeSearchSimple(s string, space string) string {
	b := make([]byte, 0, len(s))
	lastSpace := true
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			b = append(b, byte(c))
			lastSpace = false

		case c >= 'A' && c <= 'Z':
			b = append(b, byte(c)+'a'-'A')
			lastSpace = false

		case c >= 'a' && c <= 'z':
			b = append(b, byte(c))
			lastSpace = false

		case vneseMap[c] != 0:
			b = append(b, vneseMap[c])
			lastSpace = false

		default:
			if !lastSpace {
				lastSpace = true
				b = append(b, space...)
			}
		}
	}
	if lastSpace && len(b) != 0 {
		b = b[:len(b)-len(space)]
	}
	return string(b)
}

// Keep alphanumeric and some special characters while ignoring the rest.
//
//     hello@world -> hello @ world
//     hello #@@@ world -> hello # @ @@@ world
//     hello(1) -> hello ( 1 )
//     hello.world -> hello . world
func normalizeSearch(s string, space string, quote bool, lower bool) string {
	var lastChar rune
	lastGroup := 0 // space

	b := make([]byte, 0, len(s))
	for _, c := range s {
		switch {
		case c >= '0' && c <= '9':
			if lastGroup == 2 {
				b = append(b, space...)
			}
			b = append(b, byte(c))
			lastGroup = 1 // numeric

		case c >= 'A' && c <= 'Z':
			if lastGroup == 1 {
				b = append(b, space...)
			}
			if lower {
				b = append(b, byte(c)+'a'-'A')
			} else {
				b = append(b, byte(c))
			}
			lastGroup = 2 // alpha

		case c >= 'a' && c <= 'z':
			if lastGroup == 1 {
				b = append(b, space...)
			}
			b = append(b, byte(c))
			lastGroup = 2 // alpha

		case vneseMap[c] != 0:
			if lastGroup == 1 {
				b = append(b, space...)
			}
			b = append(b, vneseMap[c])
			lastGroup = 2 // alpha

		default:
			if lastGroup != 0 {
				lastGroup = 0 // space
				b = append(b, space...)
			}
			if IsAcceptedSpecialChar(c) {
				if c == lastChar {
					continue
				}
				if quote {
					b = append(b, '\'', byte(c), '\'')
				} else {
					b = append(b, byte(c))
				}
				b = append(b, space...)
				lastGroup = 0 // space
			}
		}
		lastChar = c
	}
	if lastGroup == 0 && len(b) != 0 {
		b = b[:len(b)-len(space)]
	}
	return unsafeBytesToString(b)
}

func ValidateStruct(v interface{}) (bool, error) {
	return govalidator.ValidateStruct(v)
}

func Host(s string) bool {
	return reHost.MatchString(s)
}

func DomainFromURL(s string) string {
	parts := reDomainFromURL.FindStringSubmatch(s)
	if len(parts) == 0 {
		return ""
	}
	return parts[2]
}

func NormalizedSearchToTsVector(s string) string {
	if s == "" {
		return s
	}
	ss := strings.Split(s, " ")
	for i, item := range ss {
		if item[0] == '\'' {
			ss[i] = item[1 : len(item)-1]
		}
	}
	sort.Strings(ss)

	lastItem := ""
	b := make([]byte, 0, len(s)+len(ss)*2)
	for _, item := range ss {
		if item == lastItem {
			continue
		}
		lastItem = item
		if len(b) > 0 {
			b = append(b, ' ')
		}
		b = append(b, '\'')
		b = append(b, item...)
		b = append(b, '\'')
	}
	return unsafeBytesToString(b)
}

func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
