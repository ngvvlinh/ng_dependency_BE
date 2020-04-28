package gencode

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"o.o/common/l"
)

func TestParseCode(t *testing.T) {
	alpha := Alphabet32
	for _, tt := range []struct {
		input  string
		expect int
	}{
		{"00", 0},
		{"09", 9},
		{"00", 0},
		{"0Z", 31},
		{"1Z", 63},
		{"100", 1024},
	} {
		t.Run(tt.input, func(t *testing.T) {
			output, err := alpha.Parse(tt.input)
			require.NoError(t, err)
			require.Equal(t, tt.expect, output)
		})
	}
}

func TestEncode(t *testing.T) {
	alpha := Alphabet32
	for _, tt := range []struct {
		input  uint64
		expect string
	}{
		{0, "00"},
		{9, "09"},
		{0, "00"},
		{31, "0Z"},
		{63, "1Z"},
		{1024, "100"},
	} {
		t.Run(tt.expect, func(t *testing.T) {
			output := alpha.Encode(tt.input, 2)
			require.Equal(t, tt.expect, string(output))
		})
	}
}

func TestAlphabet(t *testing.T) {
	assert.Equal(t, 22, len(Alphabet22))
	assert.Equal(t, 32, len(Alphabet32))
	assert.Equal(t, 32, len(Alphabet32Checksum))
	assert.Equal(t, 54, len(Alphabet54))
}

func testAlphabetSubset(t *testing.T, A, B string) {
	require.True(t, len(A) <= len(B))
	for _, ca := range A {
		ok := false
		for _, cb := range B {
			if ca == cb {
				ok = true
				break
			}
		}
		if !ok {
			t.Errorf("Character `%v` in `%v` does not present in `%v`", ca, A, B)
		}
	}
}

func TestGenerateMerchantOrderCode(t *testing.T) {
	mcode := "M-ABCD"
	now1 := time.Date(2018, time.January, 2, 0, 0, 0, 0, time.Local)
	now2 := time.Date(2018, time.January, 3, 23, 59, 0, 0, time.Local)

	for i := 0; i < 10; i++ {
		output1 := GenerateMerchantOrderCodeFromTime(mcode, 3, now1)
		ll.Info("GenerateMerchantOrderCode", l.String("code", output1))
		require.Len(t, output1, 9)
		require.Equal(t, "01", output1[:2])
		require.Equal(t, 1, sumChar(output1)%len(Alphabet32), "Expect checksum")
		require.Equal(t, "ABC", output1[2:5])

		output2 := GenerateMerchantOrderCodeFromTime(mcode, 3, now2)
		ll.Info("GenerateMerchantOrderCode", l.String("code", output2))
		require.Len(t, output2, 9)
		require.Equal(t, "02", output2[:2])
		Alphabet32.Parse(output2)
		require.Equal(t, 1, sumChar(output2)%len(Alphabet32), "Expect checksum")
		require.Equal(t, "ABC", output1[2:5])
	}

	mcode = "ABCDE"
	output := GenerateMerchantOrderCodeFromTime(mcode, 3, now2)
	ll.Info("GenerateMerchantOrderCode", l.String("code", output))
	require.Len(t, output, 10)
	require.Equal(t, "02", output[:2])
	Alphabet32.Parse(output)
	require.Equal(t, 1, sumChar(output)%len(Alphabet32), "Expect checksum")
	require.Equal(t, "ABCD", output[2:6])
}

func sumChar(s string) int {
	var sum int
	for i := range s {
		index := Alphabet32Checksum.Index(s[i])
		if index < 0 {
			ll.Panic("Invalid character")
		}
		sum += index
	}
	return sum
}

func TestNDigits(t *testing.T) {
	assert.Panics(t, func() {
		nDigits(0)
	})
	assert.Equal(t, 1, nDigits(1))
	assert.Equal(t, 1, nDigits(9))
	assert.Equal(t, 2, nDigits(10))
	assert.Equal(t, 2, nDigits(99))
	assert.Equal(t, 3, nDigits(100))
}

func TestGenerateOrderItemCode(t *testing.T) {
	t.Run("Zero (panic)", func(t *testing.T) {
		assert.Panics(t, func() {
			GenerateOrderItemCode("ABCD")
		})
	})

	t.Run("0 (panic)", func(t *testing.T) {
		assert.Panics(t, func() {
			GenerateOrderItemCode("ABCD", 0)
		})
	})

	t.Run("2 0 3 (panic)", func(t *testing.T) {
		assert.Panics(t, func() {
			GenerateOrderItemCode("ABCD", 2, 0, 3)
		})
	})

	t.Run("1", func(t *testing.T) {
		fn := GenerateOrderItemCode("ABCD", 1)
		assert.Equal(t, "ABCD", fn(0, 0))
		assert.Panics(t, func() {
			fn(1, 0)
		})
		assert.Panics(t, func() {
			fn(0, 1)
		})
	})

	t.Run("1 1", func(t *testing.T) {
		fn := GenerateOrderItemCode("ABCD", 1, 1)
		assert.Equal(t, "ABCD-A1", fn(0, 0))
		assert.Equal(t, "ABCD-B1", fn(1, 0))
		assert.Panics(t, func() {
			fn(2, 0)
		})
		assert.Panics(t, func() {
			fn(0, 1)
		})
	})

	t.Run("9", func(t *testing.T) {
		fn := GenerateOrderItemCode("ABCD", 9)
		assert.Equal(t, "ABCD-1", fn(0, 0))
		assert.Equal(t, "ABCD-9", fn(0, 8))
		assert.Panics(t, func() {
			fn(0, 9)
		})
	})

	t.Run("10", func(t *testing.T) {
		fn := GenerateOrderItemCode("ABCD", 10)
		assert.Equal(t, "ABCD-01", fn(0, 0))
		assert.Equal(t, "ABCD-09", fn(0, 8))
		assert.Equal(t, "ABCD-10", fn(0, 9))
		assert.Panics(t, func() {
			fn(0, 10)
		})
	})

	t.Run("1 1", func(t *testing.T) {
		fn := GenerateOrderItemCode("ABCD", 1, 1)
		assert.Equal(t, "ABCD-A1", fn(0, 0))
		assert.Equal(t, "ABCD-B1", fn(1, 0))
		assert.Panics(t, func() {
			fn(0, 1)
		})
		assert.Panics(t, func() {
			fn(2, 0)
		})
	})

	t.Run("3 10 5", func(t *testing.T) {
		fn := GenerateOrderItemCode("ABCD", 3, 10, 5)
		assert.Equal(t, "ABCD-A01", fn(0, 0))
		assert.Equal(t, "ABCD-B01", fn(1, 0))
		assert.Equal(t, "ABCD-B10", fn(1, 9))
		assert.Equal(t, "ABCD-C01", fn(2, 0))
		assert.Equal(t, "ABCD-C05", fn(2, 4))
		assert.Panics(t, func() {
			fn(0, 3)
		})
		assert.Panics(t, func() {
			fn(1, 10)
		})
		assert.Panics(t, func() {
			fn(2, 5)
		})
	})
}

func TestCheckSumDigitUPC(t *testing.T) {
	for _, tt := range []struct {
		input  string
		expect string
	}{
		{"100001", "1000016"},
		{"100012", "1000122"},
		{"111112", "1111127"},
		{"123456789", "1234567895"},
		{"1", "17"},
		{"03600024145", "036000241457"},
		{"01010101010", "010101010105"},
	} {
		t.Run(tt.expect, func(t *testing.T) {
			assert.Equal(t, tt.expect, CheckSumDigitUPC(tt.input))
		})
	}
}
