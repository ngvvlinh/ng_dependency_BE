package tests_test

import (
	"io"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"etop.vn/backend/tools/pkg/generator"
)

const testPath = "etop.vn/backend/tools/pkg/generator/tests"
const testPatterns = testPath + "/..."

type mockPlugin struct {
	ng generator.Engine

	filter   func(generator.FilterEngine) error
	generate func(generator.Engine) error
	imAlias  func(string, string) string
}

func (m *mockPlugin) Name() string    { return "mock" }
func (m *mockPlugin) Command() string { return "gen:mock" }

func (m *mockPlugin) Filter(ng generator.FilterEngine) error {
	if m.filter != nil {
		return m.filter(ng)
	}
	for _, p := range ng.ParsingPackages() {
		p.Include()
	}
	return nil
}

func (m *mockPlugin) Generate(ng generator.Engine) error {
	m.ng = ng
	if m.generate != nil {
		return m.generate(ng)
	}
	return nil
}

func (m *mockPlugin) ImportAlias(pkgPath, importPath string) string {
	if m.imAlias != nil {
		return m.imAlias(pkgPath, importPath)
	}
	return ""
}

var registered = false
var mock = &mockPlugin{}

func reset() {
	*mock = mockPlugin{} // reset the plugin
	if registered {
		return
	}
	registered = true
	if err := generator.RegisterPlugin(mock); err != nil {
		panic(err)
	}
}

func TestObjects(t *testing.T) {
	reset()
	cfg := generator.Config{}
	err := generator.Start(cfg, testPatterns)
	require.NoError(t, err)

	ng := mock.ng
	pkg := ng.GetPackageByPath(testPath + "/one")
	require.NotNil(t, pkg)

	objects := ng.GetObjectsByScope(pkg.Types.Scope())
	require.Len(t, objects, 2)
	require.Equal(t, "A", objects[0].Name())
	require.Equal(t, "B", objects[1].Name())
	directives := ng.GetDirectives(objects[1])
	require.Len(t, directives, 1)
	require.Equal(t, "gen:b", directives[0].Cmd)
}

func TestGenerate(t *testing.T) {
	reset()
	var pkgs []*generator.GeneratingPackage
	mock.generate = func(ng generator.Engine) error {
		pkgs = ng.GeneratingPackages()
		for _, pkg := range pkgs {
			// skip package "two"
			if pkg.Package.PkgPath == testPath+"/two" {
				_ = pkg.GetPrinter()
				continue
			}

			p := pkg.GetPrinter()
			mustWrite(p, []byte(" ")) // write a single byte for triggering p.Close()
		}
		return nil
	}

	cfg := generator.Config{}
	err := generator.Start(cfg, testPatterns)
	require.NoError(t, err)

	output, err := exec.Command("sh", "-c", `find . | grep zz | sort`).
		CombinedOutput()
	require.NoError(t, err)

	expected := `
./one/one-and-a-half/zz_generated.mock.go
./one/zz_generated.mock.go
./zz_generated.mock.go
`[1:]
	require.Equal(t, expected, string(output))
}

func TestClean(t *testing.T) {
	reset()
	cfg := generator.Config{CleanOnly: true}
	err := generator.Start(cfg, testPatterns)
	require.NoError(t, err)

	output, err := exec.Command("sh", "-c", `find . | grep zz | sort`).
		CombinedOutput()
	require.NoError(t, err)
	require.Equal(t, "", string(output))
}

func TestInclude(t *testing.T) {
	reset()

	parentPath := filepath.Dir(testPath)
	mock.filter = func(ng generator.FilterEngine) error {
		ng.IncludePackage(testPath + "/two")
		ng.IncludePackage(parentPath) // parentPath is outside of testPatterns
		return nil
	}

	cfg := generator.Config{}
	err := generator.Start(cfg, testPatterns)
	require.NoError(t, err)

	expecteds := []string{
		parentPath,
		testPath + "/two",
	}
	pkgs := mock.ng.GeneratingPackages()
	require.Len(t, pkgs, 2)
	for i, pkg := range pkgs {
		require.Equal(t, expecteds[i], pkg.PkgPath)
	}
}

func mustWrite(w io.Writer, p []byte) {
	if _, err := w.Write(p); err != nil {
		panic(err)
	}
}
