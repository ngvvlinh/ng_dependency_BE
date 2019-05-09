package cmstr

import (
	"reflect"
	"strings"
)

// ToGoDefinition ...
func ToGoDefinition(v interface{}) string {
	switch v := v.(type) {
	case string:
		return `"` + v + `"`
	case []string:
		var result []byte
		result = append(result, `[]string{`...)
		for i, s := range v {
			if i > 0 {
				result = append(result, ", "...)
			}
			result = append(result, '"')
			result = append(result, s...)
			result = append(result, '"')
		}
		result = append(result, '}')
		return string(result)
	}

	panic("common/str: Unsupported type " + reflect.TypeOf(v).String())
}

// ToTitle ...
func ToTitle(input string) string {
	var output string
	ss := strings.Split(input, "_")
	for _, s := range ss {
		if s == "" {
			continue
		}
		output += strings.ToUpper(string(s[0])) + s[1:]
	}
	return output
}

// ToTitleNorm ...
func ToTitleNorm(input string) string {
	var output []byte
	var upperCount int
	for i, c := range input {
		switch {
		case c >= 'A' && c <= 'Z':
			if upperCount == 0 || nextIsLower(input, i) {
				output = append(output, byte(c))
			} else {
				output = append(output, byte(c-'A'+'a'))
			}
			upperCount++

		case c >= 'a' && c <= 'z':
			output = append(output, byte(c))
			upperCount = 0

		case c >= '0' && c <= '9':
			if i == 0 {
				panic("common/str: Identifier must start with a character: `" + input + "`")
			}
			output = append(output, byte(c))
			upperCount = 0
		}
	}
	return string(output)
}

// ToSnake ...
func ToSnake(input string) string {
	var output []byte
	var upperCount int
	for i, c := range input {
		switch {
		case c >= 'A' && c <= 'Z':
			if i > 0 && (upperCount == 0 || nextIsLower(input, i)) {
				output = append(output, '_')
			}
			output = append(output, byte(c-'A'+'a'))
			upperCount++

		case c >= 'a' && c <= 'z':
			output = append(output, byte(c))
			upperCount = 0

		case c >= '0' && c <= '9':
			if i == 0 {
				panic("common/str: Identifier must start with a character: `" + input + "`")
			}
			output = append(output, byte(c))
			// prevIsLower = true

		default:
			panic("common/str: Invalid identifier: `" + input + "`")
		}
	}
	return string(output)
}

// MapToSnake ...
func MapToSnake(A []string) []string {
	B := make([]string, len(A))
	for i, s := range A {
		B[i] = ToSnake(s)
	}
	return B
}

// The next character is lower case, but not the last 's'.
//
//     HTMLFile -> html_file
//     URLs     -> urls
func nextIsLower(input string, i int) bool {
	i++
	if i >= len(input) {
		return false
	}
	c := input[i]
	if c == 's' && i == len(input)-1 {
		return false
	}
	return c >= 'a' && c <= 'z'
}

// Abbr ...
func Abbr(s string) string {
	var res []byte
	for _, c := range ToTitleNorm(s) {
		if c >= 'A' && c <= 'Z' {
			res = append(res, byte(c)-'A'+'a')
		}
	}
	return string(res)
}

func TrimLastPunctuation(s string) string {
	if s != "" {
		switch s[len(s)-1] {
		case '.', '!':
			return s[:len(s)-1]
		}
	}
	return s
}

func TrimMax(s string, max int) string {
	if len(s) >= max {
		return s[:max]
	}
	return s
}
