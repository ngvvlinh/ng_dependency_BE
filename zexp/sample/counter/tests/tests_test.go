package tests

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/sql/cmsql"
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

	count := 0
	t := time.Tick(100 * time.Millisecond)
	select {
	case <-t:
		count++
		if count > 50 { // 5 seconds
			fmt.Println("timeout waiting for process to exit")
			os.Exit(2)
		}
		output, _ := exec.Command("bash", "-c", "ps c | grep 'counter' | grep -v grep").CombinedOutput()
		if len(output) == 0 {
			ll.Info("counter exited")
			return
		}
	}
}

type serverConfig struct {
	HTTPAddress string
}

func runTest(m *testing.M) int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Default()
	db := cmsql.MustConnect(cfg.Postgres)

	_, _ = db.Exec(`DROP DATABASE IF EXISTS counter;`)
	_, _ = db.Exec(`CREATE DATABASE counter;`)

	pathDB := filepath.Join(gen.ProjectPath(), "/zexp/sample/counter/db/")
	files := loadFileWithPath(pathDB)
	cfg.Postgres.Database = "counter"

	db = cmsql.MustConnect(cfg.Postgres)

	if len(files) > 0 {
		for _, file := range files {
			query, err := ioutil.ReadFile(file)
			if err != nil {
				panic(err)
			}
			if _, err = db.Exec(string(query)); err != nil {
				panic(err)
			}
		}
	}

	readyCh, err := startServer(ctx)
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

func loadFileWithPath(path string) []string {
	var files []string

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}
	if len(files) > 0 {
		return files[1:]
	}
	return files
}

func startServer(ctx context.Context) (<-chan error, error) {
	cfg := config.Default()
	cfg.Postgres.Database = "counter"

	configStr, _ := yaml.Marshal(cfg)

	output, err := exec.Command("go", "install", "o.o/backend/zexp/sample/counter").CombinedOutput()
	if err != nil {
		ll.Info("build counter server", l.String("output", string(output)))
		return nil, err
	}

	pathCounterServer := filepath.Join(gen.GoPath(), "/bin/counter")

	cmd := exec.Command(pathCounterServer, "-config-yaml", string(configStr))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return nil, err
	}

	httpAddress := "http://" + cfg.HTTP.Address()
	httpClient = httpreq.NewResty(httpreq.RestyConfig{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	})
	httpClient.SetHostURL(httpAddress)
	httpClient.SetHeader("Content-Type", "application/json")

	process := cmd.Process
	readyCh := make(chan error)
	go func() {
		state, err2 := process.Wait()
		if err2 != nil {
			trySend(readyCh, err2)
			return
		}
		if state.ExitCode() != 0 {
			trySend(readyCh, fmt.Errorf("process exited with exit code %v", state.ExitCode()))
			return
		}
	}()

	go func() {
		defer func() {
			err2 := process.Kill()
			_ = dropDatabase()
			ll.Info("killed process", l.Error(err2))
		}()

		err2 := waitForReady(ctx, httpAddress+health.DefaultRoute)
		trySend(readyCh, err2)

		select {
		case <-ctx.Done():
			_ = process.Signal(syscall.SIGTERM)
			_, _ = process.Wait()
		}
	}()
	return readyCh, nil
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

func waitForReady(ctx context.Context, healthURL string) error {
	count := 0
	t := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil

		case <-t.C:
			count++
			if count > 60 { // 60s
				return fmt.Errorf("timeout waiting for ready")
			}
			_, err := httpClient.NewRequest().Get(healthURL)
			if err != nil {
				ll.Debug("healcheck fail: ", l.Error(err))
				continue
			}
			t.Stop()
			return nil
		}
	}
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

func trySend(ch chan error, err error) {
	select {
	case ch <- err:
	}
}
