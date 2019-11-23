// +build tools

// The purpose of this file is making go.mod import required versions. It should
// be kept in sync with install-tools.sh
package scripts

import (
	_ "github.com/gogo/protobuf/protoc-gen-gogo"
	_ "github.com/jteeuwen/go-bindata/go-bindata"
	_ "golang.org/x/tools/cmd/goimports"
)
