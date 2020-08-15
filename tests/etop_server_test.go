package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"o.o/backend/cmd/etop-server/config"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/sql/cmsql"
	e2e "o.o/backend/pkg/common/testing"
	"o.o/backend/tools/pkg/gen"
	"o.o/common/l"
)

var ll = l.New()

type M map[string]interface{}

var (
	httpServerMain *httpreq.Resty
)

const (
	routerCheckUserRegistration string = "/api/etop.User/CheckUserRegistration"
	routerInitSession           string = "/api/etop.User/InitSession"
	routerSendPhoneVerification string = "/api/etop.User/SendPhoneVerification"
	routerUserRegistration      string = "/api/etop.User/Register"
	routerUserLogin             string = "/api/etop.User/Login"

	routerGetAddresses  string = "/api/etop.Address/GetAddresses"
	routerCreateAddress string = "/api/etop.Address/CreateAddress"
	routerUpdateAddress string = "/api/etop.Address/UpdateAddress"
	routerRemoveAddress string = "/api/etop.Address/RemoveAddress"
)

func TestMain(m *testing.M) {
	exitCode := runTest(m)
	defer os.Exit(exitCode)

	e2e.StopTestProcess("counter")
}

func runTest(m *testing.M) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// load config
	cfg := config.Default()
	cfg.Databases.Postgres.Database = "test"
	cfg.Databases.PostgresAffiliate.Database = "test"
	cfg.Databases.PostgresNotifier.Database = "test"
	cfg.Databases.PostgresWebServer.Database = "test"
	cfg.Databases.PostgresLogs.Database = "test"
	db := cmsql.MustConnect(cfg.Databases.Postgres)

	_, _ = db.Exec(`DROP DATABASE IF EXISTS etop_dev_test;`)
	_, _ = db.Exec(`CREATE DATABASE etop_dev_test;`)

	pathDB := filepath.Join(gen.ProjectPath(), "/db/main/")
	contents := e2e.LoadContentPath(pathDB)

	cfg.Databases.Postgres.Database = "etop_dev_test"

	db = cmsql.MustConnect(cfg.Databases.Postgres)

	pathDBMigration := filepath.Join(gen.ProjectPath(), "/tests/main/shop/")

	contents = append(contents, e2e.LoadContentPath(pathDBMigration)...)

	err := e2e.LoadDataWithContents(db, contents)

	if err != nil {
		ll.Fatal(err.Error())
	}
	serverHost := cfg.SharedConfig.HTTP.Host
	if serverHost == "" {
		serverHost = "127.0.0.1"
	}
	httpAddress := fmt.Sprintf("http://%s:%d", serverHost, cfg.SharedConfig.HTTP.Port)
	httpServerMain = httpreq.NewResty(httpreq.RestyConfig{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	})
	httpServerMain.SetHostURL(httpAddress)
	httpServerMain.SetHeader("Content-Type", "application/json")

	configStr, _ := yaml.Marshal(cfg)

	// run server
	server, err := e2e.New(httpAddress, "etop-server", "", string(configStr))
	if err != nil {
		panic(err)
	}

	readyCh, err := server.StartServerTest(ctx, dropDatabase)
	if err != nil {
		panic(err)
	}
	err = <-readyCh
	if err != nil {
		ll.Error("can not start server", l.Error(err))
		return -1
	}

	return m.Run()
}

func dropDatabase() error {
	cfg := config.Default()
	cfg.Databases.Postgres.Database = "test"
	cfg.Databases.PostgresAffiliate.Database = "test"
	cfg.Databases.PostgresNotifier.Database = "test"
	cfg.Databases.PostgresWebServer.Database = "test"
	cfg.Databases.PostgresLogs.Database = "test"
	db := cmsql.MustConnect(cfg.Databases.Postgres)
	_, err := db.Exec(`DROP DATABASE IF EXISTS etop_dev_test;`)
	if err != nil {
		ll.Info("Drop database ", l.Error(err))
	}
	return err
}

func TestLoginAndRegistration(t *testing.T) {
	var recaptchaToken = "03AGdBq24aw3C3exEMtwydQqcuySIxxIDvILNMMRM8ei3DJ9_CiiuWd8XNwCkXJUjyA6mYwwd26d_MIUcLExeG83u203C0sHrv1voOSKsYhN5MPWxr0RgdUtcHsWLbU4OtnqUst0EGxJjYcqNEo5cM4cWU73XOQ2LUWZSgUvobpDTokoQV-zBJhw7b0uw05rO4PQYgTNp-MssFch8qfHSz7j7zsMKHCXRcqMnLiFS2MhIFg2kHhmKsSvDTLFSbrAA7nZFCVA_WDboRlj8wTSS2CxB1iJwqgw1jPoaWnAjrjQ-vARayUunDq6jiCI0vytMA1QCm_EkWaZwTVQYbKkDiCU6XSlMpz7FvU1RDIKTc9mCUBaDukVk3ZlI"
	var resp map[string]interface{}

	t.Run("register with phone number", func(t *testing.T) {
		// login with phone number
		req := M{}
		req["phone"] = "0973218967-1-test"
		req["recaptcha_token"] = recaptchaToken

		_, err := httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerCheckUserRegistration)

		require.NoError(t, err)
		assert.Equal(t, resp["exists"], false)

		req = M{}

		_, err = httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerInitSession)

		require.NoError(t, err)

		assert.Len(t, resp["access_token"], 43)
		assert.Equal(t, resp["expires_in"], float64(604800))

		var accessToken = resp["access_token"]

		req = M{}
		req["phone"] = "0973218967-1-test"

		_, err = httpServerMain.NewRequest().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", accessToken)).
			SetBody(req).SetResult(&resp).
			Post(routerSendPhoneVerification)

		require.NoError(t, err)
	})

	t.Run("registration user wrong password", func(t *testing.T) {
		req := M{
			"agree_email_info":       true,
			"agree_tos":              true,
			"auto_accept_invitation": true,
			"email":                  "etop_test@gmail.com",
			"full_name":              "etop test",
			"password":               "123456",
			"phone":                  "0987654321",
			"short_name":             "Jeremie",
		}

		_, err := httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerUserRegistration)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Len(t, resp["access_token"], 43)
		assert.Equal(t, resp["expires_in"], float64(604800))
	})

	t.Run("registration user", func(t *testing.T) {
		req := M{
			"agree_email_info":       true,
			"agree_tos":              true,
			"auto_accept_invitation": true,
			"email":                  "etop_test@gmail.com-1-test",
			"full_name":              "etop test",
			"password":               "123456789",
			"phone":                  "0987654321-1-test",
			"short_name":             "Jeremie",
		}

		_, err := httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerUserRegistration)

		require.NoError(t, err)

		require.NotNil(t, resp)
		assert.NotNil(t, resp["user"])
	})

	t.Run("user login with wrong phone number", func(t *testing.T) {
		// login with wrong phone number
		req := M{}
		req["phone"] = "0987654321-1-testtt"
		req["recaptcha_token"] = recaptchaToken

		_, err := httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerCheckUserRegistration)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, resp["exists"], false)
	})

	t.Run("user login with correct phone number and wrong password", func(t *testing.T) {
		// login with phone number
		req := M{}
		req["phone"] = "0987654321-1-test"
		req["recaptcha_token"] = recaptchaToken

		_, err := httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerCheckUserRegistration)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, resp["exists"], true)

		req = M{}
		req["login"] = "0987654321-1-test"
		req["password"] = "0987654321"
		req["account_type"] = "shop"
		_, err = httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerUserLogin)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, resp["exists"], true)
	})

	t.Run("user login with correct phone number and password", func(t *testing.T) {
		// login with phone number
		req := M{}
		req["phone"] = "0987654321-1-test"
		req["recaptcha_token"] = recaptchaToken

		_, err := httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerCheckUserRegistration)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, resp["exists"], true)

		req = M{}
		req["login"] = "0987654321-1-test"
		req["password"] = "123456789"
		req["account_type"] = "shop"
		_, err = httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerUserLogin)

		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, resp["exists"], true)
		assert.Len(t, resp["access_token"], 43)

		assert.Equal(t, resp["expires_in"], float64(604800))
		require.NotNil(t, resp["user"])
		user := resp["user"].(map[string]interface{})
		assert.Equal(t, user["email"], "etop_test@gmail.com-1-test")
		assert.Equal(t, user["phone"], "0987654321-1-test")
	})

	var accessToken string
	t.Run("user login with shop account", func(t *testing.T) {
		// login with phone number
		req := M{}
		req["phone"] = "0987654321"
		req["recaptcha_token"] = recaptchaToken

		_, err := httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerCheckUserRegistration)

		require.NoError(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, resp["exists"], true)

		req = M{}
		req["login"] = "0987654321"
		req["password"] = "123456789"
		req["account_type"] = "shop"
		_, err = httpServerMain.NewRequest().SetBody(req).SetResult(&resp).
			Post(routerUserLogin)

		require.NoError(t, err)
		require.NotNil(t, resp)

		assert.Equal(t, resp["exists"], true)
		assert.Len(t, resp["access_token"], 43)

		assert.Equal(t, resp["expires_in"], float64(604800))
		require.NotNil(t, resp["user"])
		user := resp["user"].(map[string]interface{})
		assert.Equal(t, user["email"], "test@etop.vn")
		assert.Equal(t, user["phone"], "0987654321")

		require.NotNil(t, resp["available_accounts"])

		availableAccounts := resp["available_accounts"].([]interface{})

		assert.Len(t, availableAccounts, 1)
		account := availableAccounts[0].(map[string]interface{})

		assert.Equal(t, account["id"], "1137360087033143265")
		accessToken = account["access_token"].(string)
	})

	t.Run("get address with shop account", func(t *testing.T) {
		req := M{}
		var respListAddresses map[string]interface{}
		_, err := httpServerMain.NewRequest().SetBody(req).SetHeader("authorization", fmt.Sprintf("Bearer %s", accessToken)).SetResult(&respListAddresses).
			Post(routerGetAddresses)

		require.Nil(t, err)
		require.NotNil(t, respListAddresses)

		addresses := respListAddresses["addresses"].([]interface{})

		assert.Len(t, addresses, 0)
	})

	var addressId string
	t.Run("create address with shop account", func(t *testing.T) {
		req := M{}

		req["type"] = "shipfrom"
		req["phone"] = "+840973218967-1-test"
		req["ward_code"] = "27322"
		req["district_code"] = "774"
		req["province_code"] = "79"
		req["ward"] = "Phường 07"
		req["district"] = "Quận 5"
		req["province"] = "Hồ Chí Minh"
		req["address1"] = "Hồ Chí Minh"
		req["full_name"] = "xxxx"
		req["coordinates"] = map[string]interface{}{
			"latitude":  10.234123,
			"longitude": 20.12312321,
		}
		req["notes"] = map[string]interface{}{
			"other":       "no-call",
			"open_time":   "08:00 - 12:00",
			"lunch_break": "",
		}

		var respCreateAddress map[string]interface{}
		_, err := httpServerMain.NewRequest().SetBody(req).SetHeader("authorization", fmt.Sprintf("Bearer %s", accessToken)).SetResult(&respCreateAddress).
			Post(routerCreateAddress)

		require.Nil(t, err)
		require.NotNil(t, respCreateAddress)
		require.NotNil(t, respCreateAddress["notes"])
		require.NotNil(t, respCreateAddress["coordinates"])

		assert.Len(t, respCreateAddress["id"], 19)

		addressId = respCreateAddress["id"].(string)
		assert.Equal(t, respCreateAddress["province"], "Thành phố Hồ Chí Minh")
		assert.Equal(t, respCreateAddress["province_code"], "79")
		assert.Equal(t, respCreateAddress["ward"], "Phường 07")
		assert.Equal(t, respCreateAddress["ward_code"], "27322")
		assert.Equal(t, respCreateAddress["address1"], "Hồ Chí Minh")
		assert.Equal(t, respCreateAddress["phone"], "0973218967-1-test")
		assert.Equal(t, respCreateAddress["district"], "Quận 5")
		assert.Equal(t, respCreateAddress["district_code"], "774")
		assert.Equal(t, respCreateAddress["type"], "shipfrom")
		assert.Equal(t, respCreateAddress["full_name"], "xxxx")
	})

	t.Run("update address (only full name) info with shop account", func(t *testing.T) {
		req := M{}
		req["id"] = addressId
		req["type"] = "shipfrom"
		req["full_name"] = "yyyyyyyyyyyyyyyyyyyy"

		var respUpdateAddress map[string]interface{}
		_, err := httpServerMain.NewRequest().SetBody(req).SetHeader("authorization", fmt.Sprintf("Bearer %s", accessToken)).SetResult(&respUpdateAddress).
			Post(routerUpdateAddress)

		require.Nil(t, err)
		require.NotNil(t, respUpdateAddress)
		require.NotNil(t, respUpdateAddress["notes"])
		require.NotNil(t, respUpdateAddress["coordinates"])

		assert.Equal(t, respUpdateAddress["full_name"], "yyyyyyyyyyyyyyyyyyyy")

		assert.Len(t, respUpdateAddress["id"], 19)
		assert.Equal(t, respUpdateAddress["id"], addressId)
		assert.Equal(t, respUpdateAddress["province"], "Thành phố Hồ Chí Minh")
		assert.Equal(t, respUpdateAddress["province_code"], "79")
		assert.Equal(t, respUpdateAddress["ward"], "Phường 07")
		assert.Equal(t, respUpdateAddress["ward_code"], "27322")
		assert.Equal(t, respUpdateAddress["address1"], "Hồ Chí Minh")
		assert.Equal(t, respUpdateAddress["phone"], "0973218967-1-test")
		assert.Equal(t, respUpdateAddress["district"], "Quận 5")
		assert.Equal(t, respUpdateAddress["district_code"], "774")
		assert.Equal(t, respUpdateAddress["type"], "shipfrom")
	})

	t.Run("update address info with shop account", func(t *testing.T) {
		req := M{}
		req["id"] = addressId
		req["type"] = "shipfrom"
		req["phone"] = "+840973218967-2-test"
		req["full_name"] = "zzzzzzzzzzzzzzzzzzz"

		var respUpdateAddress map[string]interface{}
		_, err := httpServerMain.NewRequest().SetBody(req).SetHeader("authorization", fmt.Sprintf("Bearer %s", accessToken)).SetResult(&respUpdateAddress).
			Post(routerUpdateAddress)

		require.Nil(t, err)
		require.NotNil(t, respUpdateAddress)
		require.NotNil(t, respUpdateAddress["notes"])
		require.NotNil(t, respUpdateAddress["coordinates"])

		assert.Equal(t, respUpdateAddress["full_name"], "zzzzzzzzzzzzzzzzzzz")

		assert.Len(t, respUpdateAddress["id"], 19)
		assert.Equal(t, respUpdateAddress["id"], addressId)
		assert.Equal(t, respUpdateAddress["province"], "Thành phố Hồ Chí Minh")
		assert.Equal(t, respUpdateAddress["province_code"], "79")
		assert.Equal(t, respUpdateAddress["ward"], "Phường 07")
		assert.Equal(t, respUpdateAddress["ward_code"], "27322")
		assert.Equal(t, respUpdateAddress["address1"], "Hồ Chí Minh")
		assert.Equal(t, respUpdateAddress["phone"], "0973218967-2-test")
		assert.Equal(t, respUpdateAddress["district"], "Quận 5")
		assert.Equal(t, respUpdateAddress["district_code"], "774")
		assert.Equal(t, respUpdateAddress["type"], "shipfrom")
	})

	t.Run("get address with shop account", func(t *testing.T) {
		req := M{}
		var respListAddresses map[string]interface{}
		_, err := httpServerMain.NewRequest().SetBody(req).SetHeader("authorization", fmt.Sprintf("Bearer %s", accessToken)).SetResult(&respListAddresses).
			Post(routerGetAddresses)

		require.Nil(t, err)
		require.NotNil(t, respListAddresses)

		addresses := respListAddresses["addresses"].([]interface{})

		assert.Len(t, addresses, 1)
	})

	t.Run("delete address with shop account", func(t *testing.T) {
		req := M{}
		req["id"] = addressId
		var respRemoveAddress map[string]interface{}
		_, err := httpServerMain.NewRequest().SetBody(req).SetHeader("authorization", fmt.Sprintf("Bearer %s", accessToken)).SetResult(&respRemoveAddress).
			Post(routerRemoveAddress)

		require.Nil(t, err)
		require.NotNil(t, respRemoveAddress)
	})

	t.Run("get address with shop account", func(t *testing.T) {
		req := M{}
		var respListAddresses map[string]interface{}
		_, err := httpServerMain.NewRequest().SetBody(req).SetHeader("authorization", fmt.Sprintf("Bearer %s", accessToken)).SetResult(&respListAddresses).
			Post(routerGetAddresses)

		require.Nil(t, err)
		require.NotNil(t, respListAddresses)

		addresses := respListAddresses["addresses"].([]interface{})

		assert.Len(t, addresses, 0)
	})
}
