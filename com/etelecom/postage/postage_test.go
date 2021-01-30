package postage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"o.o/api/etelecom/call_direction"
)

func TestCalcPostage(t *testing.T) {
	args := CalcPostageArgs{
		Direction: call_direction.Out,
	}
	t.Run("Mobiphone", func(t *testing.T) {
		args.Phone = "0931231230"

		args.DurationSecond = 20
		fee := CalcPostage(args)
		assert.Equal(t, fee, 150)

		args.DurationSecond = 13
		fee = CalcPostage(args)
		assert.Equal(t, fee, 98)
	})

	t.Run("Vinaphone", func(t *testing.T) {
		args.Phone = "0881231230"

		args.DurationSecond = 20
		fee := CalcPostage(args)
		assert.Equal(t, fee, 150)

		args.DurationSecond = 13
		fee = CalcPostage(args)
		assert.Equal(t, fee, 98)
	})

	t.Run("Order", func(t *testing.T) {
		args.Phone = "0111231230"

		args.DurationSecond = 20
		fee := CalcPostage(args)
		assert.Equal(t, fee, 268)

		args.DurationSecond = 13
		fee = CalcPostage(args)
		assert.Equal(t, fee, 175)
	})
}
