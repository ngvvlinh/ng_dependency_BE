// Only run full test with `go test -tags fulltest .`
//
//+build fulltest

package list

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationList(t *testing.T) {
	file, err := os.Open("list_test.csv")
	if err != nil {
		panic(err)
	}

	reader := csv.NewReader(file)

	// skip the first line
	row := 1
	record, err := reader.Read()
	assert.NoError(t, err, "Read the first line")
	fmt.Printf("% 4d: %#v (skip the first line)\n", row, record)

	for {
		row++
		record, err := reader.Read()
		fmt.Printf("% 4d: %#v\n", row, record)
		if err == io.EOF {
			break
		}
		assert.NoError(t, err, "Row %v: Read csv", row)
		assert.True(t, len(record) >= 3, "Row %v: Invalid length %#v", row, record)

		testLocation(t, row, record)
	}
}

const maxFail = 10
const debug = true

var fail = 0

func testLocation(t *testing.T, row int, record []string) (ok bool) {
	defer func() {
		if !ok {
			fail++
			if fail >= maxFail {
				t.Error("Too many errors")
				t.FailNow()
			}
		}
	}()

	province, district, ward := record[0], record[1], record[2]
	loc := findLocation(province, district, ward)

	if loc == nil {
		t.Errorf("Row % 4d: Can not parse province (%v,%v,%v) - %v \n-- %v",
			row, province, district, ward, locToString(loc), printInfo(debugInfo))
		return false
	}
	if loc.Province == nil {
		t.Errorf("Row % 4d: Can not parse province (%v,%v,%v) - %v \n-- %v",
			row, province, district, ward, locToString(loc), printInfo(debugInfo))
		return false
	}
	if loc.District == nil {
		t.Errorf("Row % 4d: Can not parse district (%v,%v,%v) - %v \n-- %v",
			row, province, district, ward, locToString(loc), printInfo(debugInfo))
		return false
	}
	//if loc.Ward == nil {
	//	t.Errorf("Row % 4d: Can not parse ward     (%v,%v,%v) - %v \n-- %v",
	//		row, province, district, ward, locToString(loc), printInfo(debugInfo))
	//	return false
	//}
	return true
}

func printInfo(info debugInfoStruct) string {
	return printStep(0, info.items[0]) + " - " +
		printStep(1, info.items[1]) + " - " +
		printStep(2, info.items[2])
}

func printStep(typ int, step []debugInfoStep) string {
	var b strings.Builder
	for _, s := range step {
		fmt.Fprintf(&b, "`%v` (%v) ", s.String, s.Result)
	}
	return b.String()
}
