package main

import (
	"fmt"
	cm "o.o/backend/pkg/common"
	"os"
)

func main() {
	fmt.Println(cm.NewID())
	os.Exit(1)
}
