package generator

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

type Comment struct {
	Doc *ast.CommentGroup

	Comment *ast.CommentGroup

	Text string
}

// Directive comment has one of following formats
//
//   // +etop:valid=required,optional
//   // +etop:valid=null +gen=foo
//   // +etop:pkg=foo,bar
//   // +etop:valid: 0 < $ && $ <= 10
//
// For example "// +etop:pkg=foo,bar" will be parsed as
//
//   Command: "etop:pkg"
//   Arg:     "foo,bar"
//
// Directives must start at the begin of a line, after "//" and a space (the
// same as "// +build"). Multiple directives can appear in one line.
//
// Directive ending with "=" can not have space in argument and can have
// multiple directives. Directive ending with ":" can have space in argument,
// therefore it will be parsed as a single directive.
type Directive struct {
	Raw string // +etop:pkg:foo this is a string

	Cmd string // etop:pkg

	Arg string // foo,bar
}

type declaration struct {
	Pkg *packages.Package

	Comment Comment

	Directives []Directive
}

type extendedInfo struct {
	// FileSet
	Fset *token.FileSet

	// Map from Ident to declaration
	Declarations map[*ast.Ident]*declaration

	// Map from token.Pos to Ident
	Positions map[token.Pos]*ast.Ident

	// SortedIdents is only available after calling Finalize()
	SortedIdents []*ast.Ident
}

func newExtendedInfo(fset *token.FileSet) *extendedInfo {
	return &extendedInfo{
		Fset:         fset,
		Declarations: make(map[*ast.Ident]*declaration),
		Positions:    make(map[token.Pos]*ast.Ident),
	}
}

func (x *extendedInfo) AddPackage(pkg *packages.Package) error {
	for _, file := range pkg.Syntax {
		err := x.addFile(pkg, file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *extendedInfo) addFile(pkg *packages.Package, file *ast.File) error {
	var genDoc *ast.CommentGroup
	var errs []error
	processDoc := func(doc, cmt *ast.CommentGroup) *declaration {
		if doc == nil {
			doc = genDoc
		}
		comment, directives, err := processDoc(doc, cmt)
		if err != nil {
			errs = append(errs, err)
		}
		return &declaration{
			Pkg:        pkg,
			Comment:    comment,
			Directives: directives,
		}
	}

	declarations := x.Declarations
	positions := x.Positions
	ast.Inspect(file, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.FuncDecl:
			ident := node.Name
			declarations[ident] = processDoc(node.Doc, nil)
			positions[ident.NamePos] = ident

		case *ast.GenDecl:
			// if the declaration has only 1 spec, we treat the doc
			// comment as the spec comment
			if len(node.Specs) == 1 {
				genDoc = node.Doc
			} else {
				genDoc = nil
			}

		case *ast.TypeSpec:
			ident := node.Name
			declarations[ident] = processDoc(node.Doc, node.Comment)
			positions[ident.NamePos] = ident

		case *ast.ValueSpec:
			for _, ident := range node.Names {
				declarations[ident] = processDoc(node.Doc, node.Comment)
				positions[ident.NamePos] = ident
			}

		case *ast.Field:
			for _, ident := range node.Names {
				declarations[ident] = processDoc(node.Doc, node.Comment)
				positions[ident.NamePos] = ident
			}
		}
		return true
	})

	if len(errs) != 0 {
		return newErrors("", errs)
	}
	return nil
}

func (x *extendedInfo) GetDef(ident *ast.Ident) types.Object {
	decl := x.Declarations[ident]
	if decl == nil {
		return nil
	}
	return decl.Pkg.TypesInfo.Defs[ident]
}

func (x *extendedInfo) GetComment(ident *ast.Ident) *Comment {
	decl := x.Declarations[ident]
	if decl == nil {
		return nil
	}
	return &decl.Comment
}
