// +build tools

// The purpose of this file is making go.mod import required versions. It should
// be kept in sync with install-tools.sh
package scripts

import (
	_ "github.com/axw/gocov/gocov"
	_ "github.com/go-bindata/go-bindata/v3/go-bindata"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/matm/gocov-html"
	_ "golang.org/x/tools/cmd/goimports"
)
