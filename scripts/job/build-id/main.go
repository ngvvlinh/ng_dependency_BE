package main

import (
	"flag"
	"fmt"

	"o.o/api/top/types/etc/account_tag"
	cm "o.o/backend/pkg/common"
	"o.o/capi/dot"
)

var IDType string

func main() {
	flag.StringVar(&IDType, "type", "", "Input ID type")
	flag.Parse()
	if IDType == "" {
		fmt.Println("Please input ID type. Use: -type=shop|partner|carrier|normal")
		return
	}
	fmt.Println("type :: ", IDType)
	var IDResult dot.ID
	switch IDType {
	case "shop":
		IDResult = cm.NewIDWithTag(account_tag.TagShop)
	case "partner":
		IDResult = cm.NewIDWithTag(account_tag.TagPartner)
	case "carrier":
		IDResult = cm.NewIDWithTag(account_tag.TagCarrier)
	case "normal":
		IDResult = cm.NewID()
	default:
		fmt.Println("Please input ID type. Use: shop | partner | normal")
		return
	}
	fmt.Println("id :: ", IDResult)
}
