package tests_test

import (
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"

	"etop.vn/backend/tools/pkg/generator"
)

const testPath = "etop.vn/backend/tools/pkg/generator/tests"
const testPatterns = testPath + "/..."

type mockPlugin struct {
	ng generator.Engine

	filter   func(*generator.PreparsedPackage) (bool, error)
	generate func(generator.Engine) error
	imAlias  func(string, string) string
}

func (m *mockPlugin) Name() string    { return "mock" }
func (m *mockPlugin) Command() string { return "gen:mock" }

func (m *mockPlugin) FilterPackage(p *generator.PreparsedPackage) (bool, error) {
	if m.filter != nil {
		return m.filter(p)
	}
	return true, nil
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

func TestFilter(t *testing.T) {
	reset()
	{
		var pkgs []*generator.PreparsedPackage
		mock.filter = func(pkg *generator.PreparsedPackage) (bool, error) {
			pkgs = append(pkgs, pkg)
			// skip package "one"
			if pkg.PkgPath == testPath+"/one" {
				return false, nil
			}
			return true, nil
		}

		cfg := generator.Config{}
		err := generator.Start(cfg, testPatterns)
		require.NoError(t, err)

		// verify preparsed packages
		expecteds := []struct {
			pkgPath    string
			directives []generator.Directive
		}{
			{
				pkgPath: testPath,
			},
			{
				pkgPath: testPath + "/one",
				directives: []generator.Directive{
					{Raw: "+gen", Cmd: "gen", Arg: ""},
					{Raw: "+gen:sample=10", Cmd: "gen:sample", Arg: "10"},
					{Raw: "+gen:last: 20: number:int * x", Cmd: "gen:last", Arg: "20: number:int * x"},
				},
			},
			{
				pkgPath: testPath + "/one/one-and-a-half",
			},
			{
				pkgPath: testPath + "/two",
			},
		}
		require.Len(t, pkgs, len(expecteds))
		for i, pkg := range pkgs {
			expected := expecteds[i]
			require.Equal(t, expected.pkgPath, pkg.PkgPath)
			require.EqualValues(t, expected.directives, pkg.Directives)
		}
	}
	{
		// verify that package "one" is skipped
		expecteds := []string{
			testPath,
			testPath + "/one/one-and-a-half",
			testPath + "/two",
		}

		pkgs := mock.ng.GeneratingPackages()
		for _, p := range pkgs {
			t.Logf("generating package %v", p.Package.PkgPath)
		}
		require.Equal(t, len(expecteds), len(pkgs))
		for i, p := range pkgs {
			require.Equal(t, expecteds[i], p.Package.PkgPath)
		}
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
				continue
			}

			p := pkg.Generate()
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

func mustWrite(w io.Writer, p []byte) {
	if _, err := w.Write(p); err != nil {
		panic(err)
	}
}
