// +build !release

package fabo

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"o.o/backend/tools/pkg/gen"
)

func Asset(name string) ([]byte, error) {
	base := filepath.Join(gen.ProjectPath(), "res/dl/fabo")
	if strings.Contains(name, "..") {
		panic(fmt.Sprintf("invalid name (%v)", name))
	}
	return ioutil.ReadFile(filepath.Join(base, name))
}
