package gencode

import (
	"crypto/rand"
	"math/big"
	"strconv"
	"strings"
	"time"

	"etop.vn/common/l"
)

var (
	startDate time.Time
	startCode int

	ll = l.New()
)

func init() {
	var err error
	startCode, err = Alphabet32.Parse("00")
	if err != nil {
		ll.Panic("Unexpected", l.Error(err))
	}

	startDate = time.Date(2018, time.January, 1, 0, 0, 0, 0, time.Local)
	code := GetOrderPrefixByDate(startDate.Add(time.Hour * 24))
	if code != "01" {
		ll.Panic("Unexpected code")
	}
}

// GetOrderPrefixByDate ...
func GetOrderPrefixByDate(t time.Time) string {
	days := int(t.Sub(startDate) / (24 * time.Hour))
	s := Alphabet32.Encode(uint64(startCode+days), 2)
	return string(s)
}

// GenerateMerchantOrderCodeFromTime ...
func GenerateMerchantOrderCodeFromTime(merchant string, length int, t time.Time) string {
	// Remove prefix M- from merchant code
	parts := strings.Split(merchant, "-")
	if len(parts) > 1 {
		merchant = parts[1]
	}

	prefix := GetOrderPrefixByDate(t)
	mercha := merchant[:len(merchant)-1]
	code := prefix + mercha + GenerateCode(Alphabet32, length)
	return code + CalcChecksumCharStr(Alphabet32Checksum, code, 1)
}

// GenerateOrderItemCode ...
func GenerateOrderItemCode(mocode string, numbers ...int) func(gs, i int) string {
	var ndigits int
	var gscode []string

	switch len(numbers) {
	case 0:
		ll.Panic("GenerateOrderItemCode: Unexpected args")

	case 1:
		ndigits = nDigits(numbers[0])

	default:
		gscode = make([]string, len(numbers))
		ngs := nChars(len(numbers))
		for gs := range gscode {
			gscode[gs] = string(Alphabet22.Encode(uint64(gs), ngs))
		}

		max := 0
		for _, n := range numbers {
			if n <= 0 {
				ll.Panic("Invalid number of item", l.String("mcode", mocode), l.Object("numbers", numbers))
			}
			if n > max {
				max = n
			}
		}
		ndigits = nDigits(max)
	}

	return func(gs, i int) string {
		if gs < 0 || gs >= len(numbers) {
			ll.Panic("Invalid gs", l.Int("gs", gs), l.Int("i", i), l.Int("expected gs", len(numbers)))
		}
		if i < 0 || i >= numbers[gs] {
			ll.Panic("Invalid i", l.Int("gs", gs), l.Int("i", i), l.Int("expected i", numbers[gs]))
		}

		icode := toStrWithPad(i+1, ndigits)
		if len(numbers) == 1 {
			if numbers[0] == 1 {
				return mocode
			}
			return mocode + "-" + icode
		}
		return mocode + "-" + gscode[gs] + icode
	}
}

// GenerateMerchantCode ...
func GenerateMerchantCode(length int) string {
	code := GenerateCode(Alphabet32, length)
	return code + CalcChecksumCharStr(Alphabet32Checksum, code, 1)
}

// GeneratePAgentCode ...
func GeneratePAgentCode(length int) string {
	code := GenerateCode(Alphabet32, length)
	return code + CalcChecksumCharStr(Alphabet32Checksum, code, 2)
}

// GenerateCategoryCode ...
func GenerateCategoryCode(length int) string {
	code := GenerateCode(Alphabet32, length)
	return code + CalcChecksumCharStr(Alphabet32Checksum, code, 3)
}

// CalcChecksumCharStr ...
func CalcChecksumCharStr(alphabet Alphabet, s string, checksumNumber int) string {
	var sum int
	for i := range s {
		sum += alphabet.Index(s[i])
	}

	sum = sum % len(alphabet)
	checksum := (len(alphabet) - sum + checksumNumber) % len(alphabet)
	return string(alphabet[checksum])
}

// CalcChecksumChar ...
func CalcChecksumChar(alphabet Alphabet, s []byte, checksumNumber int) byte {
	var sum int
	for i := range s {
		sum += alphabet.Index(s[i])
	}

	sum = sum % len(alphabet)
	checksum := (len(alphabet) - sum + checksumNumber) % len(alphabet)
	return alphabet[checksum]
}

// VerifyChecksum ...
func VerifyChecksum(alphabet Alphabet, s string, checksumNumber int) bool {
	var sum int
	for i := range s {
		index := Alphabet32.Index(s[i])
		sum += index
	}
	return sum == checksumNumber
}

// GenerateMerchantOrderCodeWithRandom ...
func GenerateMerchantOrderCodeWithRandom(merchant string, length int) string {
	return GenerateMerchantOrderCodeFromTime(merchant, length, time.Now())
}

// GenerateMerchantOrderCode ...
func GenerateMerchantOrderCode(merchant string, length int) string {

	// Remove prefix M- from merchant code
	parts := strings.Split(merchant, "-")
	if len(parts) > 1 {
		merchant = parts[1]
	}

	t := time.Now()
	prefix := GetOrderPrefixByDate(t)
	return prefix + "." + merchant + "."
}

func maxN(numbers ...int) int {
	max := numbers[0]
	for _, n := range numbers {
		if n > max {
			max = n
		}
	}
	return max
}

func toStrWithPad(i, minLength int) string {
	s := strconv.Itoa(i)
	for len(s) < minLength {
		s = "0" + s
	}
	return s
}

func nDigits(i int) int {
	switch {
	case i <= 0:
		ll.Panic("Must be positive number")
	case i < 10:
		return 1
	case i < 100:
		return 2
	case i < 1000:
		return 3
	}

	ll.Panic("Unsupported too many items (i)")
	return -1
}

func nChars(i int) int {
	const L = len(Alphabet22)

	switch {
	case i <= 0:
		ll.Panic("Must be positive number")
	case i < L:
		return 1
	case i < L*L:
		return 2
	}

	ll.Panic("Unsupported too many items (p)")
	return -1
}

func GenerateShopCode() string {
	length := 4
	code := GenerateCode(Alphabet32, length)
	return code
}

func GenerateOrderCode(shopCode string, t time.Time) string {
	if shopCode == "" {
		panic("empty shop code")
	}

	orderCodeLength := 4
	prefix := GetOrderPrefixByDate(t)

	orderCode := GenerateCode(Alphabet32, orderCodeLength)
	code := prefix + "." + shopCode + "." + orderCode
	return code + CalcChecksumCharStr(Alphabet32Checksum, code, 1)
}

func GenerateCodeWithType(char string, shopCode string, t time.Time) string {
	if shopCode == "" {
		panic("empty shop code")
	}
	if char == "" {
		panic("empty char code")
	}

	codeLength := 4
	prefix := GetOrderPrefixByDate(t)

	lastPart := GenerateCode(Alphabet32, codeLength)
	code := string(char) + "-" + prefix + "." + shopCode + "." + lastPart
	return code + CalcChecksumCharStr(Alphabet32Checksum, code, 1)
}

// GenerateOrderLineCode ...
func GenerateLineCode(mocode string, number int) func(i int) string {
	ngs := nChars(number)
	return func(gs int) string {
		if gs < 0 || gs >= number {
			ll.Panic("Invalid gs", l.Int("gs", gs), l.Int("expected gs", number))
		}
		return mocode + "-" + string(Alphabet22.Encode(uint64(gs), ngs))
	}
}

var max6digits = big.NewInt(1e6)
var max8digits = big.NewInt(1e8)

func Random6Digits() (string, error) {
	return randomDigits(6, max6digits)
}

func Random8Digits() (string, error) {
	return randomDigits(8, max8digits)
}

func randomDigits(length int, max *big.Int) (string, error) {
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	s := strconv.FormatInt(n.Int64(), 10)
	for len(s) < length {
		s = "0" + s
	}
	return s, nil
}

// CheckSumDigitUPC calculates sum digit using UPC method
//
// https://en.wikipedia.org/wiki/Check_digit#UPC
func CheckSumDigitUPC(s string) string {
	if len(s) == 0 {
		panic("invalid input")
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			panic("invalid input")
		}
	}
	odd, even := 0, 0
	for i := 0; i < len(s); i++ {
		num := int(s[i] - '0')
		if i%2 == 0 { // i+1 is odd
			odd += num
		} else {
			even += num
		}
	}
	sum := (odd*3 + even) % 10
	if sum != 0 {
		sum = 10 - sum
	}
	return s + string(sum+'0')
}
