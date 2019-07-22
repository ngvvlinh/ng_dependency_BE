package cm

import (
	"crypto/rand"
	"encoding/binary"
	"strconv"
	"time"

	"etop.vn/backend/pkg/common/gencode"
	"etop.vn/common/l"
)

const magicNumber = int64(1012345678909e6>>24 + 1)

var (
	zeroDate  = time.Date(2018, 02, 15, 17, 0, 0, 0, time.UTC)
	startDate = FromMillis(zeroDate.UnixNano()/1e6 - magicNumber*10)
)

// NewIDWithTime ...
func NewIDWithTime(t time.Time) int64 {
	var b [3]byte
	_, err := rand.Read(b[:])
	if err != nil {
		ll.Panic("Unable to generate ID", l.Error(err))
	}

	d := int64(t.Sub(startDate)) / 1e7
	id := d<<24 | int64(b[0])<<16 | int64(b[1])<<8 | int64(b[2])
	return id & ^int64(1<<24)
}

// NewID create new int64 ID
func NewID() int64 {
	return NewIDWithTime(time.Now())
}

// NewIDWithTimeAndTag create new int64 ID with 8-bit tag
func NewIDWithTimeAndTag(t time.Time, tag byte) int64 {
	id := NewIDWithTime(t)
	id = id & ^int64(255<<24)
	return id | int64(tag|1)<<24
}

// NewIDWithTag create new int64 ID with 8 bit tag
func NewIDWithTag(tag byte) int64 {
	return NewIDWithTimeAndTag(time.Now(), tag)
}

// GetTag ...
func GetTag(id int64) int64 {
	tag := (id >> 24) & 255
	if tag%2 == 0 {
		return 0
	}
	return tag
}

func IDToDec(id int64) string {
	return strconv.FormatInt(id, 10)
}

func DecToID(s string) (int64, bool) {
	id, _ := strconv.ParseInt(s, 10, 64)
	return id, id != 0
}

func IDToHex(id int64) string {
	return strconv.FormatInt(id, 16)
}

func HexToID(s string) (int64, bool) {
	id, _ := strconv.ParseInt(s, 16, 64)
	return id, id != 0
}

// Base54IDLength ...
const Base54IDLength = 8 + 11 + 1

const checksumMagicNumber = 9

// NewBase54IDWithEntropy ...
func NewBase54IDWithEntropy(ts uint64, e uint64) string {
	var buf [Base54IDLength]byte
	s1 := gencode.Alphabet54.Encode(ts, 8)
	s2 := gencode.Alphabet54.EncodeReverse(e, 11)

	b := buf[:0]
	b = append(b, s1[len(s1)-8:]...)
	b = append(b, s2[:11]...)
	c := gencode.CalcChecksumChar(gencode.Alphabet54, b, checksumMagicNumber)
	b = append(b, c)
	return string(b)
}

// NewBase54ID ...
func NewBase54ID() string {
	ts := uint64(time.Now().UnixNano() / 1e5)
	var b [8]byte
	_, err := rand.Read(b[:])
	if err != nil {
		ll.Panic("Unable to generate ID", l.Error(err))
	}
	e := binary.BigEndian.Uint64(b[:])

	id := NewBase54IDWithEntropy(ts, e)
	return id
}
