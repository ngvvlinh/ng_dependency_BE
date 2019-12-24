// +build !release

package doc

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"etop.vn/backend/tools/pkg/gen"
)

func Asset(name string) ([]byte, error) {
	base := filepath.Join(gen.ProjectPath(), "doc")
	if strings.Contains(name, "..") {
		panic(fmt.Sprintf("invalid name (%v)", name))
	}
	return ioutil.ReadFile(filepath.Join(base, name))
}
