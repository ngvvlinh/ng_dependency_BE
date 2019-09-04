package mapping

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"etop.vn/api/supporting/crm/vtiger"
)

// ConfigMap maps config [lead, contact, tickets]
type ConfigMap map[string]ConfigGroup

// ConfigGroup maps vtiger and etop for each group
type ConfigGroup map[string]string

// Mapper contains infomation file_mapping.json
type Mapper struct {
	FieldMap ConfigMap
}

func NewMappingConfigInfo(configMap ConfigMap) *Mapper {
	return &Mapper{
		FieldMap: configMap,
	}
}

const TimeLayout = "2006-01-02 15:04:05"

func mapVtigerToEtop(e interface{}, mappingFieldMap map[string]string, viterValueMap map[string]string) error {
	v := reflect.Indirect(reflect.ValueOf(e))
	t := v.Type()

	for n, i := v.NumField(), 0; i < n; i++ {
		field := v.Field(i)
		fieldType := field.Type().String()
		fieldTag := string(t.Field(i).Tag)
		fieldName := t.Field(i).Name

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

		case "time.Time":
			timeValue, err := time.ParseInLocation(TimeLayout, rawValue, time.Local)
			if err != nil {
				return fmt.Errorf("can not read %#v into field %v of type %v (%v)", rawValue, fieldName, fieldType, err)
			}
			field.Set(reflect.ValueOf(timeValue))

		default:
			return fmt.Errorf("can not read into field %v of type %v", fieldName, fieldType)
		}
	}
	return nil
}

func mapEtopToVtiger(e interface{}, mappingFieldMap map[string]string) (map[string]string, error) {
	mapVtiger := make(map[string]string)
	v := reflect.Indirect(reflect.ValueOf(e))
	t := v.Type()

	for n, i := v.NumField(), 0; i < n; i++ {
		field := v.Field(i)
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

		switch fieldType {
		case "int32", "int64":
			if fieldType == "int64" && field.Interface().(int64) == 0 {
				continue
			}
			if fieldType == "int32" && field.Interface().(int32) == 0 {
				continue
			}
			mapVtiger[key] = fmt.Sprint(field)

		case "time.Time":
			ts := field.Interface().(time.Time)
			if ts.IsZero() {
				continue
			}
			mapVtiger[key] = ts.Format(TimeLayout)

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

func (f *Mapper) MapingContactEtop2Vtiger(contact *vtiger.Contact) (map[string]string, error) {
	var mappingFieldMap = f.FieldMap["Contacts"]
	mapReturn, err := mapEtopToVtiger(contact, mappingFieldMap)
	return mapReturn, err
}

func (f *Mapper) MapingContactVtiger2Etop(v map[string]string) (*vtiger.Contact, error) {
	var mappingFieldMap = f.FieldMap["Contacts"]
	var contact vtiger.Contact
	err := mapVtigerToEtop(&contact, mappingFieldMap, v)
	return &contact, err
}

func (f *Mapper) MappingTicketEtop2Vtiger(ticket *vtiger.Ticket) (map[string]string, error) {
	var mappingFieldMap = f.FieldMap["HelpDesk"]
	mapReturn, err := mapEtopToVtiger(ticket, mappingFieldMap)
	return mapReturn, err
}

func (f *Mapper) MappingTicketVtiger2Etop(v map[string]string) (*vtiger.Ticket, error) {
	var mappingFieldMap = f.FieldMap["HelpDesk"]
	var ticket vtiger.Ticket
	err := mapVtigerToEtop(&ticket, mappingFieldMap, v)
	return &ticket, err
}

func (f *Mapper) MappingLeadEtop2Vtiger(lead *vtiger.Lead) (map[string]string, error) {
	var mappingFieldMap = f.FieldMap["Leads"]
	mapReturn, err := mapEtopToVtiger(lead, mappingFieldMap)
	return mapReturn, err
}

func (f *Mapper) MappingLeadVtiger2Etop(v map[string]string) (*vtiger.Lead, error) {
	var mappingFieldMap = f.FieldMap["Leads"]
	var lead vtiger.Lead
	err := mapVtigerToEtop(&lead, mappingFieldMap, v)
	return &lead, err
}
