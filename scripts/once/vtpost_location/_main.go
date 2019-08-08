package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	locationlist "etop.vn/backend/com/main/location/list"

	"github.com/360EntSecGroup-Skylar/excelize"

	"etop.vn/backend/cmd/etop-server/config"
	cc "etop.vn/backend/pkg/common/config"
	vtpostClient "etop.vn/backend/pkg/integration/vtpost/client"
	"etop.vn/common/l"
)

var (
	cfg                config.Config
	ll                 = l.New()
	vtpostProvincesMap = make(map[int]*vtpostClient.Province)
	vtpostDistrictsMap = make(map[int]*vtpostClient.District)
)

type District struct {
	Code             string
	DistrictID       int
	DistrictName     string
	ProvinceName     string
	IsRepresentative bool
}

type Location struct {
	ProvinceName string
	DistrictName string
	Area         string
	Region       string
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

	// ctx := context.Background()

	file, err := os.Open(filepath.Join(gopath, "src/etop.vn/backend/scripts/vtpost_location/khu_vuc.xlsx"))
	if err != nil {
		panic(err)
	}

	rawData, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	excelFile, err := excelize.OpenReader(bytes.NewReader(rawData))
	if err != nil {
		panic(err)
	}
	sheetName := excelFile.GetSheetName(1)
	rows := excelFile.GetRows(sheetName)
	var locations []*Location
	regionLabel := "Miền"
	for _, row := range rows {
		province := row[2]
		district := row[3]
		area := row[4]
		region := row[5]
		if checkNull(province, district, area, region) && region != regionLabel {
			locations = append(locations, &Location{
				ProvinceName: province,
				DistrictName: district,
				Area:         area,
				Region:       region,
			})
		}
	}
	count := 0
	for _, loc := range locations {
		_loc := location.FindLocation(loc.ProvinceName, loc.DistrictName, "")
		if _loc.District == nil {
			fmt.Printf("// District not found: %v | %v\n", loc.ProvinceName, loc.DistrictName)
			count++
			continue
		}
		_loc.District.UrbanType = location.GetUrbanType(loc.Area)
	}
	fmt.Printf("// Total District not found: %v/%v\n", count, len(locations))

	notFoundCount := 0
	for _, d := range locationlist.Districts {
		if d.UrbanType == 0 {
			notFoundCount++
		}
	}
	fmt.Printf("// Total Etop District not found: %v/%v\n", notFoundCount, len(locationlist.Districts))
	fmt.Printf("package location\n\nvar Districts = []*District{\n")
	for _, d := range locationlist.Districts {
		fmt.Printf(`{
		Code: "%v",
		ProvinceCode: "%v",
		Name: "%v",
		GhnID: %v,
		VTPostID: %v,
	`, d.Code, d.ProvinceCode, d.Name, d.GhnID, d.VTPostID)
		if d.UrbanType != 0 {
			fmt.Printf("\tUrbanType: %#v,\n", d.UrbanType)
		}
		if len(d.Alias) > 0 {
			fmt.Printf("\tAlias: %#v,\n", d.Alias)
		}
		fmt.Print("},\n")
	}
	fmt.Printf("}\n")
}

func checkNull(ss ...string) bool {
	for _, s := range ss {
		if s == "" {
			return false
		}
	}
	return true
}

var ProvinceRegionName = map[location.Region]string{
	1: "North",
	2: "Middl",
	3: "South",
}
