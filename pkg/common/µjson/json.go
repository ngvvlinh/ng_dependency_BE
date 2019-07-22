// Minimal JSON parser which works with correct input only.
// Usecase:
//   1. Walk through unstructured json
//   2. Transform unstructured json
// without fully unmarshalling it into a map[string]interface{}
//
// Caution: Behaviour is undefined on invalid json. Use on trusted input only.

package µjson

import (
	"encoding/json"
	"fmt"
	"strings"

	cm "etop.vn/backend/pkg/common"
	"etop.vn/common/l"
)

var ll = l.New()

func Walk(s []byte, i int, fn func(st int, key, value string) bool) error {
	var si, ei, st int
	var key string

	// fn returns false to skip a whole array or object
	sst := 1024

	// Trim the last newline
	if len(s) > 0 && s[len(s)-1] == '\n' {
		s = s[:len(s)-1]
	}

value:
	si = i
	switch s[i] {
	case 'n', 't': // null, true
		i += 4
		ei = i
		if st <= sst {
			fn(st, key, string(s[si:i]))
		}
		key = ""
		goto closing
	case 'f': // false
		i += 5
		ei = i
		if st <= sst {
			fn(st, key, string(s[si:i]))
		}
		key = ""
		goto closing
	case '{', '[':
		if st <= sst && !fn(st, key, string(s[i])) {
			sst = st
		}
		key = ""
		st++
		i++
		if s[i] == '}' || s[i] == ']' {
			goto closing
		}
		goto value
	case '"': // scan string
		for {
			i++
			switch s[i] {
			case '\\': // \. - skip 2
				i++
			case '"': // end of string
				i++
				ei = i // space, ignore
				for s[i] == ' ' ||
					s[i] == '\t' ||
					s[i] == '\n' ||
					s[i] == '\r' {
					i++
				}
				if s[i] != ':' {
					if st <= sst {
						fn(st, key, string(s[si:ei]))
					}
					key = ""
				}
				goto closing
			}
		}
	case ' ', '\t', '\n', '\r': // space, ignore
		i++
		goto value
	default: // scan number
		for i < len(s) {
			switch s[i] {
			case ',', '}', ']', ' ', '\t', '\n', '\r':
				ei = i
				for s[i] == ' ' ||
					s[i] == '\t' ||
					s[i] == '\n' ||
					s[i] == '\r' {
					i++
				}
				if st <= sst {
					fn(st, key, string(s[si:ei]))
				}
				key = ""
				goto closing
			}
			i++
		}
	}

closing:
	if i >= len(s) {
		return nil
	}
	switch s[i] {
	case ':':
		key = string(s[si:ei])
		i++
		goto value
	case ',':
		i++
		goto value
	case ']', '}':
		st--
		if st == sst {
			sst = 1024
		} else {
			fn(st, "", string(s[i]))
		}
		if st <= 0 {
			return nil
		}
		i++
		goto closing
	case ' ', '\t', '\n', '\r':
		i++ // space, ignore
		goto closing
	default:
		return parseError(i, s[i], `expect ']', '}' or ','`)
	}
}

func parseError(i int, c byte, msg string) error {
	return fmt.Errorf("json error at %v '%c' 0x%2x: %v", i, c, c, msg)
}

func ShouldAddComma(value string, lastChar byte) bool {
	return value != "}" && value != "]" &&
		lastChar != ',' && lastChar != '{' && lastChar != '['
}

func Reconstruct(s []byte) ([]byte, error) {
	b := make([]byte, 0, 1024)
	err := Walk(s, 0, func(st int, key, value string) bool {
		if len(b) != 0 && ShouldAddComma(value, b[len(b)-1]) {
			b = append(b, ',')
		}
		if key != "" {
			b = append(b, key...)
			b = append(b, ':')
		}
		b = append(b, value...)
		return true
	})
	return b, err
}

func FilterAndRename(b []byte, input []byte) (output []byte, _ error) {
	// Validate json input and output
	if cm.IsDev() {
		lenb := len(b)
		defer func() {
			var v interface{}
			out := b[lenb:]
			err1 := json.Unmarshal(input, &v)
			err2 := json.Unmarshal(out, &v)
			if err1 != nil || err2 != nil {
				if err1 != nil {
					ll.Error("Invalid json input")
				}
				if err2 != nil {
					ll.Error("Invalid json output")
				}
			}
		}()
	}

	err := Walk(input, 0, func(st int, key, value string) bool {

		// Ignore fields with null value
		if value == "null" {
			return true
		}

		wrap := false
		if key != "" {
			// Remove quotes
			key = key[1 : len(key)-1]

			// Skip _ids
			if strings.HasSuffix(key, "_ids") {
				return false
			}

			// Rename external_ to x_
			if strings.HasPrefix(key, "external_") {
				key = "x_" + key[len("external_"):]

			} else if (key == "id" || strings.HasSuffix(key, "_id")) &&
				value[0] >= '0' && value[0] <= '9' {
				wrap = true
			}
		}

		if len(b) != 0 && ShouldAddComma(value, b[len(b)-1]) {
			b = append(b, ',')
		}
		if key != "" {
			b = append(b, '"')
			b = append(b, key...)
			b = append(b, '"')
			b = append(b, ':')
		}
		if wrap {
			b = append(b, '"')
			b = append(b, value...)
			b = append(b, '"')
		} else {
			b = append(b, value...)
		}
		return true
	})
	return b, err
}
