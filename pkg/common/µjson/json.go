// Minimal JSON parser which works with correct input only.
// Usecase:
//   1. Walk through unstructured json
//   2. Transform unstructured json
// without fully unmarshalling it into a map[string]interface{}
//
// Caution: Behaviour is undefined on invalid json. Use on trusted input only.

package µjson

import (
	"fmt"
	"strings"

	"etop.vn/backend/pkg/common/cmenv"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

var ll = l.New()

func parseError(i int, c byte, msg string) error {
	return fmt.Errorf("json error at %v '%c' 0x%2x: %v", i, c, c, msg)
}

func ShouldAddComma(value string, lastChar byte) bool {
	return value != "}" && value != "]" &&
		lastChar != ',' && lastChar != '{' && lastChar != '['
}

func Reconstruct(s []byte) ([]byte, error) {
	b := make([]byte, 0, 1024)
	err := Walk(s, 0, func(_, st int, key, value string) bool {
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
	if cmenv.IsDev() {
		lenb := len(b)
		defer func() {
			var v interface{}
			out := b[lenb:]
			err1 := jsonx.Unmarshal(input, &v)
			err2 := jsonx.Unmarshal(out, &v)
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

	err := Walk(input, 0, func(pos, st int, key, value string) bool {

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
