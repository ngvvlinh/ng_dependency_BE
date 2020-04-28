package randgenerator

import (
	"math/rand"

	"o.o/backend/pkg/common/code/gencode"
)

type RandGenerator struct {
	rd *rand.Rand
}

func NewGenerator(seed int64) *RandGenerator {
	src := rand.NewSource(seed)
	rd := rand.New(src)
	return &RandGenerator{rd}
}

func (g *RandGenerator) RandomAlphabet32(len int) []byte {
	n := g.rd.Uint64()
	res := gencode.Alphabet32.EncodeReverse(n, len)
	return res[:len]
}
