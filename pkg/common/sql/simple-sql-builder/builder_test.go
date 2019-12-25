package builder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	tests := []struct {
		Name     string
		Query    string
		Args     []interface{}
		Error    string
		Expected string
	}{
		{
			"ok (1)",
			"SELECT * FROM table WHERE id = ? AND name = ? AND enabled = ?",
			[]interface{}{10, "foo", true},
			"",
			"SELECT * FROM table WHERE id = 10 AND name = 'foo' AND enabled = true",
		},
		{
			"ok (2)",
			"SELECT * FROM table WHERE id = ? AND status = true",
			[]interface{}{"foo"},
			"",
			"SELECT * FROM table WHERE id = 'foo' AND status = true",
		},
		{
			"argument length does not match (1)",
			"SELECT * FROM table",
			[]interface{}{10},
			"expected 0 arguments but 1 arguments were provided",
			"",
		},
		{
			"argument length does not match (2)",
			"SELECT * FROM table WHERE id = ? AND status = ?",
			[]interface{}{10},
			"not enough argument at index 1",
			"",
		},
		{
			"invalid argument type",
			"SELECT * FROM table WHERE id = ?",
			[]interface{}{complex(10, 20)},
			"unexpected argument of type complex128 ((10+20i))",
			"",
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			var b SimpleSQLBuilder
			b.Printf(tt.Query, tt.Args...)
			query, err := b.String()
			if tt.Error == "" {
				require.NoError(t, err)
				require.Equal(t, tt.Expected, query)
			} else {
				require.EqualError(t, err, tt.Error)
			}
		})
	}
}
