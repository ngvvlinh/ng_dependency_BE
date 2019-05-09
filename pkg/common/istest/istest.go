// +build !production

package istest

import "flag"

var isTest bool

func init() {
	isTest = flag.Lookup("test.v") != nil
}

// IsTest ...
func IsTest() bool {
	return isTest
}
