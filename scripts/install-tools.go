// The purpose of this file is making go.mod import required versions. It should be kept in sync with install-tools.sh
package scripts

import (
	_ "github.com/elliots/protoc-gen-twirp_swagger/genswagger"
	_ "github.com/gogo/protobuf/proto"
	_ "github.com/jteeuwen/go-bindata"
	_ "github.com/twitchtv/twirp"
	_ "golang.org/x/tools/imports"
)
