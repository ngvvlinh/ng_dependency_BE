package tests

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/sql/cmsql"
	e2e "o.o/backend/pkg/common/testing"
	"o.o/backend/tools/pkg/gen"
	"o.o/backend/zexp/sample/counter/config"
	"o.o/common/l"
)

var ll = l.New()
var httpClient *httpreq.Resty

type M map[string]interface{}

const (
	routeUpsertCounter string = "/counter.Counter/Counter"
	routeGetCounter    string = "/counter.Counter/Get"
)

func TestMain(m *testing.M) {
	exitCode := runTest(m)
	defer os.Exit(exitCode)

	e2e.StopTestProcess("counter")
}

func runTest(m *testing.M) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Default()
	db := cmsql.MustConnect(cfg.Postgres)

	_, _ = db.Exec(`DROP DATABASE IF EXISTS counter;`)
	_, _ = db.Exec(`CREATE DATABASE counter;`)

	pathDB := filepath.Join(gen.ProjectPath(), "/zexp/sample/counter/db/")
	contents := e2e.LoadContentPath(pathDB)
	cfg.Postgres.Database = "counter"

	db = cmsql.MustConnect(cfg.Postgres)

	err := e2e.LoadDataWithContents(db, contents)

	if err != nil {
		ll.Fatal(err.Error())
	}

	httpAddress := fmt.Sprintf("http://%s", cfg.HTTP.Address())
	httpClient = httpreq.NewResty(httpreq.RestyConfig{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	})
	httpClient.SetHostURL(httpAddress)
	httpClient.SetHeader("Content-Type", "application/json")

	configStr, _ := yaml.Marshal(cfg)

	// run server
	server, err := e2e.New(httpAddress, "counter", "", string(configStr))
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
	db := cmsql.MustConnect(cfg.Postgres)
	_, err := db.Exec(`DROP DATABASE IF EXISTS counter;`)
	if err != nil {
		ll.Info("Drop database ", l.Error(err))
	}
	return err
}

func TestUpsertCounter(t *testing.T) {
	t.Run("add 10", func(t *testing.T) {
		req := M{}
		req["label"] = "testlabel"
		req["value"] = 10

		var resp map[string]interface{}
		_, err := httpClient.NewRequest().SetBody(req).SetResult(&resp).
			Post(routeUpsertCounter)

		require.NoError(t, err)
		assert.Equal(t, resp["value"], float64(10))
	})
	t.Run("get counter", func(t *testing.T) {
		req := M{}
		req["label"] = "testlabel"

		var resp map[string]interface{}
		_, err := httpClient.NewRequest().SetBody(req).SetResult(&resp).
			Post(routeGetCounter)

		require.NoError(t, err)
		assert.Equal(t, resp["value"], float64(10))
	})
	t.Run("add another", func(t *testing.T) {
		req := M{}
		req["label"] = "testlabel"
		req["value"] = 10

		var resp map[string]interface{}
		_, err := httpClient.NewRequest().SetBody(req).SetResult(&resp).
			Post(routeUpsertCounter)

		require.NoError(t, err)
		assert.Equal(t, resp["value"], float64(20))
	})
	t.Run("get counter", func(t *testing.T) {
		req := M{}
		req["label"] = "testlabel"

		var resp map[string]interface{}
		_, err := httpClient.NewRequest().SetBody(req).SetResult(&resp).
			Post(routeGetCounter)

		require.NoError(t, err)
		assert.Equal(t, resp["value"], float64(20))
	})
}
