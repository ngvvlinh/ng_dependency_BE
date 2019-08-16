package main

import (
	"flag"
	"fmt"

	cm "etop.vn/backend/pkg/common"
)

var flN = flag.Int("N", 1, "Number of id to generate")

func main() {
	for i, n := 0, *flN; i < n; i++ {
		fmt.Println(cm.NewID())
	}
}
