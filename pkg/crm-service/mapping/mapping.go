package mapping

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	pbcm "etop.vn/backend/pb/common"
	"etop.vn/backend/pb/services/crmservice"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/k0kubun/pp"
)

// MappingConfigInfo contain infomation file_mapping.json
type MappingConfigInfo struct {
	FieldMap *ConfigMap
}

func NewMappingConfigInfo(configMap *ConfigMap) *MappingConfigInfo {
	return &MappingConfigInfo{
		FieldMap: configMap,
	}
}

var configText = `
{
  "Tickets": {
    "note": "cf_884",
    "old_value": "cf_908",
    "new_value": "cf_910"
  }
}
`

// 2019-07-17 10:13:53
const timeLayout = "2006-01-02 15:04:05"

// ConfigMap map config [lead, contact, tickets]
type ConfigMap map[string]GroupConfig

// GroupConfig map vtiger and etop for each group
type GroupConfig map[string]string

type VtigerContact map[string]string

// 2019-07-17 10:13:53
const timeLayoutVtiger = "2006-01-02 15:04:05"

func (f *MappingConfigInfo) mapVtigerToEtop(e interface{}, mappingFieldMap map[string]string, viterValueMap map[string]string) error {
	v := reflect.Indirect(reflect.ValueOf(e))
	t := v.Type()

	for n, i := v.NumField(), 0; i < n; i++ {
		field := v.Field(i)
		fieldType := field.Type().String()
		fieldTag := string(t.Field(i).Tag)
		fieldName := t.Field(i).Name

		// process Tag Name to compate
		arrayTag := strings.Split(string(fieldTag), " ")
		tagNameVtiger := strings.Replace(arrayTag[len(arrayTag)-1], `json:"`, "", 1)
		tagNameVtiger = strings.Replace(tagNameVtiger, `"`, "", 1)

		key := tagNameVtiger
		key = strings.ReplaceAll(key, ",", "")
		key = strings.ReplaceAll(key, "omitempty", "")

		rawValue := viterValueMap[key]
		if rawValue == "" {
			keyVtiger := mappingFieldMap[key]
			rawValue = viterValueMap[keyVtiger]
		}

		if rawValue == "" {
			continue
		}
		if fieldType != "string" && fieldType != "time.Time" && fieldType != "int" {
			pp.Printf("------------------------------------------------fieldType value :: %v :: \n", fieldType)
		}
		switch fieldType {
		case "string":
			field.Set(reflect.ValueOf(rawValue))

		case "int64":
			rawValue = strings.ReplaceAll(rawValue, ".", "")
			intValue, err := strconv.ParseInt(rawValue, 10, 64)
			if err != nil {
				return fmt.Errorf("can not read %#v into field %v of type %v (%v)", rawValue, fieldName, fieldType, err)
			}
			field.Set(reflect.ValueOf(intValue))
		case "int32":
			rawValue = strings.ReplaceAll(rawValue, ".", "")
			intValue, err := strconv.ParseInt(rawValue, 10, 32)
			if err != nil {
				return fmt.Errorf("can not read %#v into field %v of type %v (%v)", rawValue, fieldName, fieldType, err)
			}
			field.Set(reflect.ValueOf(intValue))
		case "timestamp.Timestamp":
			timeValue, err := time.ParseInLocation(timeLayout, rawValue, time.Local)
			if err != nil {
				return fmt.Errorf("can not read %#v into field %v of type %v (%v)", rawValue, fieldName, fieldType, err)
			}
			field.Set(reflect.ValueOf(*pbcm.PbTime(timeValue)))

		default:
			return fmt.Errorf("can not read into field %v of type %v", fieldName, fieldType)
		}
	}
	return nil
}

func mapEtopToVtiger(e interface{}, mappingFieldMap map[string]string) (map[string]string, error) {

	mapVtiger := make(map[string]string)
	//pp.Print(mappingFieldMap)

	v := reflect.Indirect(reflect.ValueOf(e))
	t := v.Type()

	for n, i := v.NumField(), 0; i < n; i++ {
		field := v.Field(i)
		//fieldName := t.Field(i).Name
		fieldType := field.Type().String()

		fieldTag := t.Field(i).Tag
		arrayTag := strings.Split(string(fieldTag), " ")
		tagNameVtiger := strings.Replace(arrayTag[len(arrayTag)-1], `json:"`, "", 1)
		tagNameVtiger = strings.Replace(tagNameVtiger, `"`, "", 1)

		key := mappingFieldMap[tagNameVtiger]

		if key == "" {
			key = tagNameVtiger
		}
		key = strings.ReplaceAll(key, ",", "")
		key = strings.ReplaceAll(key, "omitempty", "")

		if key == "modifiedtime" || key == "createdtime" || key == "-" {
			continue
		}
		pp.Println("field type : ", fieldType)
		switch fieldType {
		case "int32", "int64":
			if fieldType == "int64" && field.Interface().(int64) == 0 {
				continue
			}
			if fieldType == "int32" && field.Interface().(int32) == 0 {
				continue
			}
			mapVtiger[key] = fmt.Sprint(field)
		case "*timestamp.Timestamp":
			time := field.Interface().(*timestamp.Timestamp)
			t := pbcm.PbTimeToModel(time)
			if t.IsZero() {
				continue
			}
			mapVtiger[key] = t.Format(timeLayoutVtiger)
		default:
			value := field.String()
			if value == "" {
				continue
			}
			mapVtiger[key] = value
		}
	}
	return mapVtiger, nil

}

func (f *MappingConfigInfo) MapingContactEtop2Vtiger(contact *crmservice.Contact) (map[string]string, error) {

	configText := f.FieldMap
	var config ConfigMap
	config = *configText

	var mappingFieldMap = config["Contacts"]

	mapReturn, err := mapEtopToVtiger(contact, mappingFieldMap)
	return mapReturn, err
}

func (f *MappingConfigInfo) MapingContactVtiger2Etop(v map[string]string) (*crmservice.Contact, error) {

	configText := f.FieldMap
	var config ConfigMap
	config = *configText

	var mappingFieldMap = config["Contacts"]

	contact := &crmservice.Contact{}
	err := f.mapVtigerToEtop(contact, mappingFieldMap, v)
	return contact, err
}

func (f *MappingConfigInfo) MappingTicketEtop2Vtiger(tiket *crmservice.Ticket) (map[string]string, error) {
	configText := f.FieldMap
	var config ConfigMap
	config = *configText

	var mappingFieldMap = config["HelpDesk"]

	mapReturn, err := mapEtopToVtiger(tiket, mappingFieldMap)
	return mapReturn, err
}

func (f *MappingConfigInfo) MappingTicketVtiger2Etop(v map[string]string) (*crmservice.Ticket, error) {
	configText := f.FieldMap
	var config ConfigMap
	config = *configText

	var mappingFieldMap = config["HelpDesk"]

	ticket := &crmservice.Ticket{}
	err := f.mapVtigerToEtop(ticket, mappingFieldMap, v)
	return ticket, err
}

func (f *MappingConfigInfo) MappingLeadEtop2Vtiger(lead *crmservice.Lead) (map[string]string, error) {
	configText := f.FieldMap
	var config ConfigMap
	config = *configText

	var mappingFieldMap = config["Lead"]

	mapReturn, err := mapEtopToVtiger(lead, mappingFieldMap)
	return mapReturn, err
}

func (f *MappingConfigInfo) MappingLeadVtiger2Etop(v map[string]string) (*crmservice.Lead, error) {
	configText := f.FieldMap
	var config ConfigMap
	config = *configText

	var mappingFieldMap = config["Lead"]

	lead := &crmservice.Lead{}
	err := f.mapVtigerToEtop(lead, mappingFieldMap, v)
	return lead, err
}
