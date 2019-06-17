package gencode

import (
	"crypto/rand"
	"encoding/binary"
	"errors"

	"etop.vn/backend/pkg/common/l"
)

// Alphabet ...
type Alphabet string

// Alphabet
const (
	Alphabet22 Alphabet = `ABCDEFGHJKLMNPQRSTVWYZ`
	Alphabet32 Alphabet = `0123456789ABCDEFGHJKLMNPQRSTVWYZ`
	Alphabet54 Alphabet = `0123456789ABCDEFGHJKLMNPQRSTVWYZabcdefghjklmnpqrstvwyz`

	Alphabet32Checksum Alphabet = `LY0TABZ2W4NGVRHDFQ6M59SP1KCJE738`
)

// Index ...
func (a Alphabet) Index(c byte) int {
	for i := 0; i < len(a); i++ {
		if a[i] == c {
			return i
		}
	}
	return -1
}

// Parse ...
func (a Alphabet) Parse(s string) (int, error) {
	var v int
	for i := 0; i < len(s); i++ {
		index := a.Index(s[i])
		if index < 0 {
			return v, errors.New("invalid character")
		}

		v *= len(a)
		v += index
	}
	return v, nil
}

// Encode ...
func (a Alphabet) Encode(v uint64, minLen int) []byte {
	s := a.EncodeReverse(v, minLen)
	for i, l := 0, len(s); i < l/2; i++ {
		s[i], s[l-i-1] = s[l-i-1], s[i]
	}
	return s
}

// EncodeReverse ...
func (a Alphabet) EncodeReverse(v uint64, minLen int) []byte {
	var s []byte
	for v > 0 || len(s) < minLen {
		l := uint64(len(a))
		s = append(s, a[v%l])
		v = v / l
	}
	return s
}

// GenerateCode ...
func GenerateCode(alphabet Alphabet, length int) string {
	const word = 8
	b := make([]byte, length*word)
	n, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	if n < length*word {
		ll.Panic("Unexpected n < l", l.Int("n", n))
	}

	s := make([]byte, length)
	for i := 0; i < length; i++ {
		bi := b[i*word : (i+1)*word]
		ui := binary.BigEndian.Uint64(bi)
		s[i] = alphabet[ui%uint64(len(alphabet))]
	}
	return string(s)
}
