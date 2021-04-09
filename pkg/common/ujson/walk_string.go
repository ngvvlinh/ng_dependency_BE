package ujson

func WalkString(s string, i int, fn func(pos, st int, key, value string) bool) error {
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
		pos := i
		i += 4
		ei = i
		if st <= sst {
			fn(pos, st, key, string(s[si:i]))
		}
		key = ""
		goto closing
	case 'f': // false
		pos := i
		i += 5
		ei = i
		if st <= sst {
			fn(pos, st, key, string(s[si:i]))
		}
		key = ""
		goto closing
	case '{', '[':
		pos := i
		if st <= sst && !fn(pos, st, key, string(s[i:i+1])) {
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
		pos := i
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
						fn(pos, st, key, string(s[si:ei]))
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
		pos := i
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
					fn(pos, st, key, string(s[si:ei]))
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
		pos := i
		st--
		if st == sst {
			sst = 1024
		} else {
			fn(pos, st, "", string(s[i:i+1]))
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
