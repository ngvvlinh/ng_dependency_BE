package testing

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"o.o/backend/pkg/common/apifw/health"
	"o.o/backend/pkg/common/apifw/httpreq"
	"o.o/backend/pkg/common/bus"
	"o.o/backend/pkg/common/projectpath"
	"o.o/backend/pkg/common/sql/cmsql"
	"o.o/common/l"
)

var (
	httpClient *httpreq.Resty
)

type ServerTest struct {
	HttpAddress string
	PackageName string
	ConfigFile  string
	ConfigYml   string
	Client      *httpreq.Resty
}

type Content struct {
	Path string
	Body []byte
}

// default config
var DefaultServerTest = ServerTest{
	HttpAddress: "127.0.0.1",
	Client: httpreq.NewResty(httpreq.RestyConfig{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}),
	PackageName: "etop-server",
}

func New(httpAddress, packageName, configFile, configYml string) (*ServerTest, error) {
	server := DefaultServerTest

	err := mustLoadConfig(httpAddress, packageName, configFile, configYml, &server)

	return &server, err
}

func mustLoadConfig(httpAddress, packageName, configFile, configYml string, cfg *ServerTest) error {
	if httpAddress != "" {
		cfg.HttpAddress = httpAddress
	}
	if packageName != "" {
		cfg.PackageName = packageName
	}

	if configFile != "" && configYml != "" {
		return fmt.Errorf("must provide only -config-file or -config-yaml")
	}

	if configFile != "" {
		cfg.ConfigFile = configFile
	}

	if configYml != "" {
		cfg.ConfigYml = configYml
	}

	return nil
}

func (s *ServerTest) StartServerTest(ctx context.Context, funcClean func() error) (<-chan error, error) {
	output, err := exec.Command("go", "list", "o.o/backend/...").Output()

	if err != nil {
		return nil, fmt.Errorf("can not get list directory: %v", err)
	}

	arrPackagePath := []string{}
	for _, value := range strings.Split(string(output), "\n") {
		if !strings.Contains(value, fmt.Sprintf("/%s", s.PackageName)) {
			continue
		}
		if strings.LastIndex(value, s.PackageName) != (len(value) - len(s.PackageName)) {
			continue
		}
		arrPackagePath = append(arrPackagePath, value)
	}

	if len(arrPackagePath) == 0 {
		return nil, fmt.Errorf("not found package")
	}

	packagePath := arrPackagePath[0]

	//Install compiles and installs the packages named by the import paths
	output, err = exec.Command("go", "install", packagePath).CombinedOutput()
	if err != nil {
		ll.Info(fmt.Sprintf("build %s server fail", s.PackageName), l.String("output", string(output)))
		return nil, err
	}

	// get path server with $GOPATH
	pathServer := filepath.Join(projectpath.GetGoPath(), fmt.Sprintf("/bin/%s", s.PackageName))

	var cmd *exec.Cmd
	if s.ConfigFile != "" {
		cmd = exec.Command(pathServer, "-config-file", s.ConfigFile)
	} else if s.ConfigYml != "" {
		cmd = exec.Command(pathServer, "-config-yaml", s.ConfigYml)
	} else {
		cmd = exec.Command(pathServer)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		return nil, err
	}

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
			_ = funcClean()
			ll.Info("killed process", l.Error(err2))
		}()

		err2 := s.waitForReady(ctx)
		trySend(readyCh, err2)

		select {
		case <-ctx.Done():
			_ = process.Signal(syscall.SIGTERM)
			_, _ = process.Wait()
		}
	}()
	return readyCh, nil
}

func trySend(ch chan error, err error) {
	select {
	case ch <- err:
	}
}

func (s *ServerTest) waitForReady(ctx context.Context) error {
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

			_, err := s.Client.SetHostURL(s.HttpAddress).NewRequest().Get(health.DefaultRoute)
			if err != nil {
				ll.Debug("healcheck fail: ", l.Error(err))
				continue
			}
			t.Stop()
			return nil
		}
	}
}

func StopTestProcess(serverName string) {
	count := 0
	t := time.Tick(100 * time.Millisecond)
	select {
	case <-t:
		count++
		if count > 50 { // 5 seconds
			fmt.Println("timeout waiting for process to exit")
			os.Exit(2)
		}
		output, _ := exec.Command("bash", "-c", fmt.Sprintf("ps c | grep '%s' | grep -v grep", serverName)).CombinedOutput()
		if len(output) == 0 {
			ll.Info("counter exited")
			return
		}
	}
}

func LoadContentPath(sqlPath string) []Content {
	var contents []Content
	err := filepath.Walk(sqlPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			ll.Error("error", l.Error(err))
			return err
		}
		if info == nil {
			ll.Error("unexpected")
			return err
		}

		baseName := filepath.Base(path)
		if strings.HasPrefix(baseName, "_") {
			if info.IsDir() {
				ll.S.Infof("skipped directory %v", baseName)
			} else {
				ll.S.Infof("skipped file %v", baseName)
			}
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(baseName, ".sql") {
			ll.S.Infof("skipped non-sql file %v", baseName)
			return nil
		}
		body, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		contents = append(contents, Content{path, body})
		return nil
	})

	if err != nil {
		ll.Error("Error while executing", l.Error(err))
	}
	return contents
}

func LoadDataWithContents(db *cmsql.Database, contents []Content) error {
	return db.InTransaction(bus.Ctx(), func(tx cmsql.QueryInterface) error {
		for _, content := range contents {
			ll.S.Infof("--- Executing %v", content.Path)
			_, _err := tx.SQL(string(content.Body)).Exec()
			if _err != nil {
				ll.Error("Error while executing", l.String("script", content.Path), l.Error(_err))
				return _err
			}
		}
		return nil
	})
}
