package generator

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"

	"etop.vn/common/l"
)

var buildFlags = strings.Split("-tags generator", " ")

type GenerateFileNameInput struct {
	PluginName string
}

type Config struct {
	// default to "zz_generated.{{.Name}}.go"
	GenerateFileName func(GenerateFileNameInput) string

	EnabledPlugins []string

	CleanOnly bool

	Namespace string
}

func Start(cfg Config, patterns ...string) error {
	return theEngine.clone().start(cfg, patterns...)
}

func (ng *engine) start(cfg Config, patterns ...string) (_err error) {
	{
		if len(patterns) == 0 {
			return errorf(nil, "no patterns")
		}
		if len(ng.plugins) == 0 {
			return errorf(nil, "no registed plugins")
		}
		if err := ng.validateConfig(&cfg); err != nil {
			return err
		}
		ng.xcfg = cfg
	}
	{
		mode := packages.NeedName | packages.NeedImports | packages.NeedDeps |
			packages.NeedFiles | packages.NeedCompiledGoFiles
		ng.pkgcfg = packages.Config{Mode: mode}
		pkgs, err := packages.Load(&ng.pkgcfg, patterns...)
		if err != nil {
			return errorf(err, "can not load package: %v", err)
		}

		// populate cleanedFileNames
		cleanedFileNames := make(map[string]bool)
		for _, pl := range ng.enabledPlugins {
			input := GenerateFileNameInput{PluginName: pl.name}
			filename := ng.genFilename(input)
			cleanedFileNames[filename] = true
		}
		for _, pkg := range pkgs {
			for _, file := range pkg.GoFiles {
				filename := filepath.Base(file)
				if cleanedFileNames[filename] {
					if err := os.Remove(file); err != nil {
						return errorf(err, "can not remove file %v: %v", file, err)
					}
				}
			}
		}
		ng.cleanedFileNames = cleanedFileNames
		if cfg.CleanOnly {
			return nil
		}

		// populate collectedPackages, includes, srcMap
		ps := parsePatterns(patterns)
		if err = ng.collectPackages(ps, pkgs); err != nil {
			return err
		}

		if ll.Verbosed(4) {
			for _, pkg := range ng.collectedPackages {
				ll.V(4).Debugf("collected package: %v", pkg.PkgPath)
			}
		}
	}
	{
		var pkgPatterns []string
		for i, pkg := range ng.collectedPackages {
			if ng.includes[i] {
				pkgPatterns = append(pkgPatterns, pkg.PkgPath)
			}
		}
		ll.V(3).Debug("load all syntax from", l.Any("patterns", pkgPatterns))
		if len(pkgPatterns) == 0 {
			fmt.Println("no packages for generating")
			return nil
		}

		ng.pkgcfg = packages.Config{
			Mode:       packages.LoadAllSyntax,
			BuildFlags: buildFlags,
			Overlay:    ng.srcMap,
		}
		pkgs, err := packages.Load(&ng.pkgcfg, pkgPatterns...)
		if err != nil {
			return errorf(err, "can not load package: %v", err)
		}

		// populate xinfo
		ng.xinfo = newExtendedInfo(ng.pkgcfg.Fset)
		packages.Visit(pkgs,
			func(pkg *packages.Package) bool {
				if cfg.Namespace != "" && !strings.HasPrefix(pkg.PkgPath, cfg.Namespace) {
					return true
				}
				if err := ng.xinfo.AddPackage(pkg); err != nil {
					_err = err
					return false
				}
				return true
			}, nil)
		if _err != nil {
			return _err
		}

		// populate pkgMap
		packages.Visit(pkgs,
			func(pkg *packages.Package) bool {
				ng.pkgMap[pkg.PkgPath] = pkg
				return true
			}, nil)
	}
	{
		// populate generatedFile
		for _, pl := range ng.enabledPlugins {
			wrapNg := &wrapEngine{engine: ng, plugin: pl}
			if err := pl.plugin.Generate(wrapNg); err != nil {
				return errorf(err, "%v: %v", pl.name, err)
			}
			for _, gpkg := range wrapNg.pkgs {
				printer := gpkg.printer
				if printer != nil {
					// close the printer for writing to file, but only if there
					// are any bytes written
					if err := printer.Close(); err != nil {
						return err
					}
				}
			}
		}
	}
	{
		sort.Strings(ng.generatedFile)
		fmt.Println("Generated file:")
		pwd, err := os.Getwd()
		must(err)
		for _, filename := range ng.generatedFile {
			filename, err = filepath.Rel(pwd, filename)
			must(err)
			fmt.Printf("\t./%v\n", filename)
		}
		if err := execGoimport(ng.generatedFile); err != nil {
			return err
		}
	}
	return nil
}

func (ng *engine) collectPackages(ps patternsStruct, pkgs []*packages.Package) error {
	collectedPackages, fileContents, err := collectPackages(ps, pkgs, ng.cleanedFileNames)
	if err != nil {
		return err
	}
	sort.Slice(collectedPackages, func(i, j int) bool {
		return collectedPackages[i].PkgPath < collectedPackages[j].PkgPath
	})
	ng.includes = make([]bool, len(collectedPackages))
	for _, pl := range ng.enabledPlugins {
		pl.includes = make([]bool, len(collectedPackages))
		for i := range collectedPackages {
			ppkg := collectedPackages[i] // capture the value
			include, err := pl.plugin.FilterPackage(&ppkg)
			if err != nil {
				return errorf(err, "%v: filter %v: %v", pl.name, ppkg.PkgPath, err)
			}
			ng.includes[i] = ng.includes[i] || include
			pl.includes[i] = include
			if include {
				pl.includesN++
			}
		}
	}
	ng.collectedPackages = collectedPackages

	srcMap := make(map[string][]byte)
	for _, fileContent := range fileContents {
		srcMap[fileContent.Path] = fileContent.Body
	}
	ng.srcMap = srcMap
	return nil
}

type fileContent struct {
	Path string
	Body []byte
}

func collectPackages(
	ps patternsStruct, pkgs []*packages.Package, cleanedFileNames map[string]bool,
) (collectedPackages []PreparsedPackage, files []fileContent, _err error) {

	var wg0, wg sync.WaitGroup
	wg0.Add(2)

	// collection file contents
	fileCh := make(chan fileContent, 4)
	go func() {
		defer wg0.Done()
		for file := range fileCh {
			files = append(files, file)
		}
	}()

	// collect errors
	errCh := make(chan error, 4)
	var errs []error
	go func() {
		defer wg0.Done()
		for err := range errCh {
			errs = append(errs, err)
		}
	}()

	limit := make(chan struct{}, 16) // limit concurrency
	collectedPackages = make([]PreparsedPackage, len(pkgs))
	for i := range pkgs {
		i, pkg := i, pkgs[i] // capture values for closure
		limit <- struct{}{}  // limit
		wg.Add(1)
		go func() {
			defer func() { wg.Done(); <-limit }() // release limit
			directives, err := parseDirectivesFromPackage(fileCh, pkg, cleanedFileNames)
			if err != nil {
				_err = errorf(err, "parsing %v: %v", pkg.PkgPath, err)
			}
			p := PreparsedPackage{
				PkgPath:    pkg.PkgPath,
				Imports:    pkg.Imports,
				Directives: directives,
			}
			collectedPackages[i] = p
		}()
	}
	wg.Wait()
	close(fileCh)
	close(errCh)
	wg0.Wait()
	if len(errs) != 0 {
		_err = newErrors("can not parse packages", errs)
	}
	return
}

func parseDirectivesFromPackage(fileCh chan<- fileContent, pkg *packages.Package, cleanedFileNames map[string]bool) (directives []Directive, _err error) {
	for _, file := range pkg.CompiledGoFiles {
		if cleanedFileNames[filepath.Base(file)] {
			continue
		}
		body, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		fileCh <- fileContent{Path: file, Body: body}
		ds, errs := parseDirectivesFromBody(directives, body)
		if len(errs) != 0 {
			// ignore unknown directives
			if ll.Verbosed(2) {
				for _, e := range errs {
					ll.V(1).Debugf("ignored %v", e)
				}
			}
			continue
		}
		directives = append(directives, ds...)
	}
	return
}

var startDirective = []byte(startDirectiveStr)

func parseDirectivesFromBody(directives []Directive, body []byte) (_ []Directive, errs []error) {
	// store processing directives, they may be discarded if they are not
	// followed by a blank line
	var tmp []Directive
	lastIdx := -1
	for idx := 1; idx < len(body); idx++ {
		if body[idx] != '\n' {
			continue
		}

		// process the last found directive, remove "// " but keep "+"
		if lastIdx >= 0 {
			line := body[lastIdx+len(startDirective)-1 : idx]
			lastIdx = -1

			ds, err := ParseDirective(string(line))
			if err != nil {
				errs = append(errs, err)
				continue
			}
			tmp = append(tmp, ds...)
		}
		// directives are followed by a blank line, accept them
		if idx+1 < len(body) && body[idx+1] == '\n' {
			directives = append(directives, tmp...)
			tmp = tmp[:0]
		}
		// find the next directive
		if !bytes.HasPrefix(body[idx+1:], startDirective) && idx+1 != len(body) {
			// discard directives not followed by a blank line
			tmp = tmp[:0]
			continue
		}
		lastIdx = idx + 1
	}
	// source file should end with a newline, so we don't process remaining lastIdx
	directives = append(directives, tmp...)
	return directives, errs
}

func (ng *engine) validateConfig(cfg *Config) (_err error) {
	defer func() {
		if _err != nil {
			_err = errorf(_err, "config error: %v", _err)
		}
	}()

	// populate enabledPlugins
	if cfg.EnabledPlugins != nil {
		for _, enabled := range cfg.EnabledPlugins {
			pl := ng.pluginsMap[enabled]
			if pl == nil {
				return errorf(nil, "plugin %v not found", enabled)
			}
			pl.enabled = true
			ng.enabledPlugins = append(ng.enabledPlugins, pl)
		}
	} else {
		for _, pl := range ng.plugins {
			pl.enabled = true
		}
		ng.enabledPlugins = ng.plugins
	}

	if cfg.GenerateFileName == nil {
		cfg.GenerateFileName = defaultGeneratedFileName(defaultGeneratedFileNameTpl)
	}

	if ng.bufPool.New == nil {
		ng.bufPool.New = func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, defaultBufSize))
		}
	}
	return nil
}

func (ng *engine) genFilename(input GenerateFileNameInput) string {
	return ng.xcfg.GenerateFileName(input)
}

func (ng *engine) writeFile(filePath string) (io.WriteCloser, error) {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	ng.generatedFile = append(ng.generatedFile, filePath)
	return f, nil
}

func execGoimport(files []string) error {
	args := []string{"-w"}
	args = append(args, files...)
	cmd := exec.Command("goimports", args...)
	ll.V(4).Debugf("goimports %v", args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errorf(err, "goimports: %s\n\n%s\n", err, out)
	}
	return nil
}
