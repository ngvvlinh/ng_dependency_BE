// +build !release

package istest

import (
	"os"
	"strings"
)

var isTest bool

func init() {
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, "-test.") {
			isTest = true
			break
		}
	}
}

// IsTest ...
func IsTest() bool {
	return isTest
}
