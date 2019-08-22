package generator

import (
	"bytes"
	"fmt"
	"go/ast"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"

	"etop.vn/common/l"
)

const defaultGeneratedFileNameTpl = "zz_generated.%v.go"
const defaultBufSize = 1024 * 4
const startDirectiveStr = "// +"

var ll = l.New()
var reCommand = regexp.MustCompile(`[a-z]([a-z0-9.:-]*[a-z0-9])?`)

func FilterByCommand(command string) Filter {
	return filterByCommand(command)
}

type filterByCommand string

func (cmd filterByCommand) FilterPackage(p *PreparsedPackage) (bool, error) {
	for _, d := range p.Directives {
		if d.Cmd == string(cmd) {
			return true, nil
		}
	}
	return false, nil
}

func defaultGeneratedFileName(tpl string) func(GenerateFileNameInput) string {
	return func(input GenerateFileNameInput) string {
		return fmt.Sprintf(tpl, input.PluginName)
	}
}

// processDoc splits directive and text comment
func processDoc(doc, cmt *ast.CommentGroup) (Comment, []Directive, error) {
	comment := Comment{
		Doc:     doc,
		Comment: cmt,
	}
	if doc == nil {
		return comment, nil, nil
	}

	processedDoc := make([]*ast.Comment, 0, len(doc.List))
	directives := make([]Directive, 0, 4)
	for _, line := range doc.List {
		if !strings.HasPrefix(line.Text, startDirectiveStr) {
			continue
		}
		// remove "// " but keep "+"
		text := line.Text[len(startDirectiveStr)-1:]
		_directives, err := ParseDirective(text)
		if err != nil {
			return Comment{}, nil, err
		}
		directives = append(directives, _directives...)
	}
	comment.Text = (&ast.CommentGroup{List: processedDoc}).Text()
	return comment, directives, nil
}

func ParseDirective(text string) (result []Directive, err error) {
	return parseDirective(text, result)
}

func parseDirective(text string, result []Directive) ([]Directive, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, nil
	}
	if text[0] != '+' {
		return nil, errorf(nil, "invalid directive")
	}
	cmdIdx := reCommand.FindStringIndex(text)
	if cmdIdx == nil {
		return nil, errorf(nil, "invalid directive")
	}
	if cmdIdx[0] != 1 {
		return nil, errorf(nil, "invalid directive")
	}
	dtext := text[:cmdIdx[1]]
	directive := Directive{
		Cmd: dtext[1:], // remove "+"
	}
	remain := text[len(dtext):]
	if remain == "" {
		directive.Raw = dtext
		return append(result, directive), nil
	}
	if remain[0] == ' ' || remain[0] == '\t' {
		directive.Raw = dtext
		result = append(result, directive)
		return parseDirective(remain, result)
	}
	if remain[0] == ':' {
		remain = remain[1:] // remove ":"
		directive.Raw = text
		directive.Arg = strings.TrimSpace(remain)
		if directive.Arg == "" {
			return nil, errorf(nil, "invalid directive")
		}
		return append(result, directive), nil
	}
	if remain[0] == '=' {
		remain = remain[1:] // remove "="
		idx := strings.IndexAny(text, " \t")
		if idx < 0 {
			directive.Raw = text
			directive.Arg = strings.TrimSpace(remain)
			if directive.Arg == "" {
				return nil, errorf(nil, "invalid directive")
			}
			return append(result, directive), nil
		}
		directive.Raw = text[:idx]
		directive.Arg = strings.TrimSpace(text[len(dtext)+1 : idx])
		if directive.Arg == "" {
			return nil, errorf(nil, "invalid directive")
		}
		result = append(result, directive)
		return parseDirective(text[idx:], result)
	}
	return nil, errorf(nil, "invalid directive")
}

// TODO: handle "unicode..." for being compatible with
// https://golang.org/cmd/go/#hdr-Package_lists_and_patterns

// a trie for quickly match package path
type patternsStruct struct {
	paths map[string]struct{}

	// A path prefix like "example.com/world/water/..." will be splitted to
	//
	//   "example.com":             1
	//   "example.com/world":       1
	//   "example.com/world/water": 2
	prefixes map[string]int
}

func parsePatterns(patterns []string) patternsStruct {
	paths := make(map[string]struct{})
	prefixes := make(map[string]int)

	sort.Strings(patterns)
	var lastPath, lastPattern_ string
	for _, pattern := range patterns {
		if !strings.HasSuffix(pattern, "/...") {
			if lastPattern_ == "" || !strings.HasPrefix(pattern, lastPattern_) {
				paths[pattern] = struct{}{}
				lastPath = pattern
			}
			continue
		}

		pattern = strings.TrimSuffix(pattern, "...")
		lastPattern_ = pattern
		if lastPath == pattern[:len(pattern)-1] {
			delete(paths, lastPath)
		}

		var part string
		var found bool
		for idx := 0; idx < len(pattern); idx++ {
			if pattern[idx] != '/' {
				continue
			}
			part = pattern[:idx+1]
			if prefixes[part] == 2 {
				found = true
				break
			}
			prefixes[part] = 1
		}
		if !found {
			prefixes[pattern] = 2
		}
	}
	result := patternsStruct{
		paths:    paths,
		prefixes: prefixes,
	}
	return result
}

func (ps patternsStruct) match(pkgPath string) bool {
	for idx := 0; idx < len(pkgPath); idx++ {
		if pkgPath[idx] != '/' {
			continue
		}
		part := pkgPath[:idx+1]
		if ps.prefixes[part] == 2 {
			return true
		}
		if ps.prefixes[part] == 0 {
			return false
		}
	}
	return ps.prefixes[pkgPath+"/"] == 2
}

type listErrors struct {
	Msg    string
	Errors []error
}

func newErrors(msg string, errs []error) error {
	return listErrors{Msg: msg, Errors: errs}
}

func (es listErrors) Error() string {
	return fmt.Sprint(es)
}

func (es listErrors) Format(st fmt.State, c rune) {
	if es.Msg == "" && len(es.Errors) == 0 {
		_, _ = st.Write([]byte("<nil>"))
		return
	}

	width, ok := st.Width()
	if !ok {
		width = 8
	}

	verbose := st.Flag('#') || st.Flag('+')
	var b bytes.Buffer
	if es.Msg != "" {
		b.WriteString(es.Msg)
		if len(es.Errors) == 0 {
			return
		}
		if verbose {
			b.WriteString(":\n")
		} else {
			b.WriteString(": ")
		}
	}
	for i, e := range es.Errors {
		if verbose {
			for i := 0; i < width; i++ {
				b.WriteByte(' ')
			}
		}
		_ = fmt.Sprint(b, e)
		if i > 0 {
			if verbose {
				b.WriteString("\n")
			} else {
				b.WriteString("; ")
			}
		}
	}
	_, _ = st.Write(b.Bytes())
}

func errorf(err error, format string, args ...interface{}) error {
	if err != nil {
		return errors.WithMessagef(err, format, args...)
	}
	msg := fmt.Sprintf(format, args...)
	return errors.New(msg)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
