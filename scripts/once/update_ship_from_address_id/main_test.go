package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	cc "o.o/backend/pkg/common/config"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/capi/dot"
)

func init() {
	db = cmsql.MustConnect(cc.DefaultPostgres())
}

func TestSQLAddress(t *testing.T) {
	sql, args, err := db.SQL(sqlScanAddress, dot.ID(123), "shipfrom", false).Build()
	require.NoError(t, err)

	expected := `
SELECT "id","full_name","first_name","last_name","phone","position","email","country","city","province","district","ward","zip","is_default","district_code","province_code","ward_code","company","address1","address2","type","account_id","notes","created_at","updated_at","coordinates","rid" FROM (
	SELECT DISTINCT ON (account_id) * FROM address
	WHERE (id > $1 AND type = $2 AND is_default = $3)
	ORDER BY account_id ASC, created_at ASC LIMIT 500
) t ORDER BY created_at ASC
`
	expectedArgs := []interface{}{dot.ID(123), "shipfrom", false}
	assert.Equal(t, expected, sql)
	assert.EqualValues(t, expectedArgs, args)
}
