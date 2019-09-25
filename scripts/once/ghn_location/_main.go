package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	locationlist "etop.vn/backend/com/main/location/list"

	"etop.vn/backend/cmd/etop-server/config"
	cc "etop.vn/backend/pkg/common/config"
	"etop.vn/backend/pkg/integration/shipping/vtpost"
	vtpostClient "etop.vn/backend/pkg/integration/shipping/vtpost/client"
	"etop.vn/common/bus"
	"etop.vn/common/jsonx"
	"etop.vn/common/l"
)

var (
	cfg                config.Config
	ll                 = l.New()
	vtpostProvincesMap = make(map[int]*vtpostClient.Province)
	vtpostDistrictsMap = make(map[int]*vtpostClient.District)
	vtpostWardsMap     = make(map[int]*location.VTPostWard)
)

type District struct {
	Code             string
	DistrictID       int
	DistrictName     string
	ProvinceName     string
	IsRepresentative bool
}

func main() {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		panic("No GOPATH")
	}
	cc.InitFlags()
	flag.Parse()

	var err error
	cfg, err = config.Load(true)
	if err != nil {
		ll.Fatal("Unable to load config", l.Error(err))
	}

	ctx := context.Background()
	buildWards(ctx)

	// if cfg.VTPost.AccountDefault.Username != "" {
	// 	if err := vtpost.Init(cfg.VTPost); err != nil {
	// 		ll.Fatal("Unable to connect to VTPost", l.Error(err))
	// 	}
	// 	// 	{
	// 	// 		cmd := &vtpost.GetProvincesCommand{}
	// 	// 		if err := bus.Dispatch(ctx, cmd); err != nil {
	// 	// 			ll.Fatal("VTPost: fail to get provinces", l.Error(err))
	// 	// 		}
	// 	// 		for _, p := range cmd.Result {
	// 	// 			vtpostProvincesMap[p.ProvinceID] = p
	// 	// 		}
	// 	// 		buildProvinces(cmd.Result)
	// 	// 	}
	// 	// 	{
	// 	// 		// uncomment when you want to build District
	// 	// 		cmd := &vtpost.GetDistrictsCommand{}
	// 	// 		if err := bus.Dispatch(ctx, cmd); err != nil {
	// 	// 			ll.Fatal("VTPost: fail to get districts", l.Error(err))
	// 	// 		}
	// 	// 		for _, d := range cmd.Result {
	// 	// 			vtpostDistrictsMap[d.DistrictID] = d
	// 	// 		}
	// 	// 		buildDistricts(gopath, cmd.Result)
	// 	// 	}
	// 	buildVtpostWards(ctx)
	// }
}

func buildDistricts(gopath string, vtpostDistricts []*vtpostClient.District) {
	data, err := ioutil.ReadFile(filepath.Join(gopath, "src/etop.vn/backend/pkg/integration/location/.list_ghn.json"))
	must(err)

	var list struct {
		Data []*District `json:"data"`
	}
	err = jsonx.Unmarshal(data, &list)
	must(err)

	districts := list.Data

	fmt.Printf("// Code generated by `go run ./scripts/ghn_location/main.go` DO NOT EDIT.\n\n")

	{
		count := 0
		for _, d := range districts {
			if d.IsRepresentative {
				continue
			}
			loc := location.FindLocation(d.ProvinceName, d.DistrictName, "")
			if loc.District == nil {
				fmt.Printf("// District not found: %v %v | %v\n", d.DistrictID, d.DistrictName, d.ProvinceName)
				count++
				continue
			}

			loc.District.GhnID = d.DistrictID
		}

		fmt.Printf("// Total GHN District not found: %v/%v\n", count, len(districts))
	}
	{
		count := 0
		for _, d := range vtpostDistricts {
			province := vtpostProvincesMap[d.ProvinceID]
			loc := location.FindLocation(province.ProvinceName, d.DistrictName, "")
			if loc.District == nil {
				fmt.Printf("// VTPost District not found: %v %v | %v\n", d.DistrictID, d.DistrictName, province.ProvinceName)
				count++
				continue
			}

			loc.District.VTPostID = d.DistrictID
		}

		fmt.Printf("// Total VTPost District not found: %v/%v\n", count, len(vtpostDistricts))
	}
	{
		countGHN := 0
		countVTPost := 0
		for _, d := range list.Districts {
			if d.GhnID == 0 {
				fmt.Printf("// District without ghn: %v %v | %v\n", d.Code, d.Name, d.ProvinceCode)
				countGHN++
			}
			if d.VTPostID == 0 {
				fmt.Printf("// District without vtpost: %v %v | %v\n", d.Code, d.Name, d.ProvinceCode)
				countVTPost++
			}
		}
		fmt.Printf("// Total GHN: %v/%v\n\n", countGHN, len(list.Districts))
		fmt.Printf("// Total VTPost: %v/%v\n\n", countVTPost, len(list.Districts))
	}

	fmt.Printf("package location\n\nvar Districts = []*District{\n")
	for _, d := range list.Districts {
		fmt.Printf(`{
	Code: "%v",
	ProvinceCode: "%v",
	Name: "%v",
	GhnID: %v,
	VTPostID: %v,
`, d.Code, d.ProvinceCode, d.Name, d.GhnID, d.VTPostID)
		if len(d.Alias) > 0 {
			fmt.Printf("\tAlias: %#v,\n", d.Alias)
		}
		fmt.Print("},\n")
	}
	fmt.Printf("}\n")
}

func buildProvinces(vtpostProvinces []*vtpostClient.Province) {
	fmt.Printf("// Code generated by `go run ./scripts/ghn_location/main.go` DO NOT EDIT.\n\n")
	{
		count := 0
		for _, p := range vtpostProvinces {
			loc := location.FindLocation(p.ProvinceName, "", "")
			if loc.Province == nil {
				fmt.Printf("// VTPost Province not found: %v | %v\n", p.ProvinceID, p.ProvinceName)
				count++
				continue
			}

			loc.Province.VTPostID = p.ProvinceID
		}

		fmt.Printf("// Total VTPost Province not found: %v/%v\n", count, len(vtpostProvinces))
	}
	{
		countVTPost := 0
		for _, p := range locationlist.Provinces {
			if p.VTPostID == 0 {
				fmt.Printf("// Province without vtpost: %v | %v\n", p.Code, p.Name)
				countVTPost++
			}
		}
		fmt.Printf("// Total VTPost: %v/%v\n\n", countVTPost, len(locationlist.Provinces))
	}

	fmt.Printf("package location\n\nvar Countries = []*Country{{Name: CountryVietnam}}\n\nvar Provinces = []*Province{\n")
	for _, p := range locationlist.Provinces {
		fmt.Printf(`{
	Code: "%v",
	Region: %v,
	Name: "%v",
	VTPostID: %v,
`, p.Code, ProvinceRegionName[p.Region], p.Name, p.VTPostID)
		if p.Special {
			fmt.Printf("\tSpecial: %v,\n", p.Special)
		}
		if len(p.Alias) > 0 {
			fmt.Printf("\tAlias: %#v,\n", p.Alias)
		}
		fmt.Print("},\n")
	}
	fmt.Printf("}\n")
}

func buildVtpostWards(ctx context.Context) {

	for _, district := range locationlist.Districts {
		go func(d *location.District) {
			if d.VTPostID == 0 {
				return
			}
			cmd := &vtpost.GetWardsByDistrictCommand{
				Request: &vtpostClient.GetWardsByDistrictRequest{
					DistrictID: d.VTPostID,
				},
			}
			if err := bus.Dispatch(ctx, cmd); err != nil {
				// ll.Fatal("VTPost: fail to get wards", l.Error(err))
			}
			for _, w := range cmd.Result {
				vtpostWardsMap[w.WardsID] = &location.VTPostWard{
					WardsID:          w.WardsID,
					WardsName:        w.WardsName,
					DistrictID:       w.DistrictID,
					EtopDistrictCode: d.Code,
				}
			}
		}(district)
	}
	fmt.Printf("package location\n\nvar VTPostWards = []*VTPostWard{\n")
	for _, w := range vtpostWardsMap {
		fmt.Printf(`{
			WardsID: %v,
			WardsName: "%v",
			DistrictID: %v,
			EtopDistrictCode: "%v",
		`, w.WardsID, w.WardsName, w.DistrictID, w.EtopDistrictCode)
		fmt.Printf("},\n")
	}
	fmt.Printf("}\n")
}

func buildWards(ctx context.Context) {
	fmt.Printf("// Code generated by `go run ./scripts/ghn_location/main.go` DO NOT EDIT.\n\n")
	{
		count := 0
		for _, vtWard := range locationlist.VTPostWards {
			ward := location.FindWardByDistrictCode(vtWard.WardsName, vtWard.EtopDistrictCode)
			if ward == nil {
				// fmt.Printf("// VTPost Ward not found: %v | %v\n", vtWard.WardsID, vtWard.WardsName)
				count++
				continue
			}
			ward.VTPostID = vtWard.WardsID
		}
		// fmt.Printf("// -- Total VTPost Ward not found: %v/%v\n", count, len(location.VTPostWards))
	}
	{
		countVTPost := 0
		for _, w := range locationlist.Wards {
			if w.VTPostID == 0 {
				fmt.Printf("// Ward without vtpost: %v | %v\n", w.Code, w.Name)
				countVTPost++
			}
		}
		fmt.Printf("// Total Ward not found VTPostID: %v/%v\n\n", countVTPost, len(locationlist.Wards))
	}
	fmt.Printf("package location\n\nvar Wards = []*Ward{\n")
	for _, w := range locationlist.Wards {
		fmt.Printf(`{
		Code: "%v",
		DistrictCode: "%v",
		Name: "%v",
		VTPostID: %v,
	`, w.Code, w.DistrictCode, w.Name, w.VTPostID)
		if len(w.Alias) > 0 {
			fmt.Printf("\tAlias: %#v,\n", w.Alias)
		}
		fmt.Printf("},")
	}
	fmt.Printf("}\n")
}

var ProvinceRegionName = map[location.Region]string{
	1: "North",
	2: "Middl",
	3: "South",
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
