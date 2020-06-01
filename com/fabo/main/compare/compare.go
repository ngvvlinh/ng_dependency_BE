package compare

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"o.o/api/fabo/fbmessaging"
)

var timeType = reflect.TypeOf(time.Time{})
var equalMethod, _ = timeType.MethodByName("Equal")
var equalFunc = equalMethod.Func

func Compare(item, otherItem interface{}) bool {
	itemValue := reflect.ValueOf(item)
	otherItemValue := reflect.ValueOf(otherItem)

	if !itemValue.IsValid() && !otherItemValue.IsValid() {
		return true
	}

	if itemValue.Kind() != otherItemValue.Kind() {
		return false
	}

	itemTyp := reflect.TypeOf(item)
	switch itemValue.Kind() {
	case reflect.Struct:
		if itemValue.NumField() != otherItemValue.NumField() {
			return false
		}
		for i := 0; i < itemValue.NumField(); i++ {
			if itemTyp.Field(i).Tag.Get("compare") == "ignore" {
				continue
			}

			itemType := fmt.Sprintf("%T", item)
			otherItemType := fmt.Sprintf("%T", otherItem)
			if itemType != otherItemType {
				return false
			}

			if strings.Contains(itemType, "time.Time") {
				return equalFunc.Call([]reflect.Value{itemValue, otherItemValue})[0].Bool()
			}

			if isEqual := Compare(itemValue.Field(i).Interface(), otherItemValue.Field(i).Interface()); !isEqual {
				return false
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return itemValue.Int() == otherItemValue.Int()
	case reflect.Bool:
		return itemValue.Bool() == otherItemValue.Bool()
	case reflect.String:
		return itemValue.String() == otherItemValue.String()
	case reflect.Float64, reflect.Float32:
		return itemValue.Float() == otherItemValue.Float()
	case reflect.Slice, reflect.Array:
		if itemValue.Len() != otherItemValue.Len() {
			return false
		}
		for i := 0; i < itemValue.Len(); i++ {
			if isEqual := Compare(itemValue.Index(i).Interface(), otherItemValue.Index(i).Interface()); !isEqual {
				return false
			}
		}
	case reflect.Ptr:
		countNil := 0
		if itemValue.IsNil() {
			countNil += 1
		}
		if otherItemValue.IsNil() {
			countNil += 1
		}
		switch countNil {
		case 0:
			if compare := Compare(itemValue.Elem().Interface(), otherItemValue.Elem().Interface()); !compare {
				return false
			}
		case 1:
			return false
		case 2:
			return true
		}
	case reflect.Map:
		if len(itemValue.MapKeys()) != len(otherItemValue.MapKeys()) {
			return false
		}

		mapItemKey := make(map[interface{}]reflect.Value)
		for _, key := range itemValue.MapKeys() {
			switch key.Interface().(type) {
			case int:
				mapItemKey[key.Int()] = key
			case string:
				mapItemKey[key.String()] = key
			default:
				mapItemKey[key.Interface()] = key
			}
		}

		for _, key := range otherItemValue.MapKeys() {
			switch key.Interface().(type) {
			case int:
				if _, ok := mapItemKey[key.Int()]; !ok {
					return false
				}
			case string:
				if _, ok := mapItemKey[key.String()]; !ok {
					return false
				}
			default:
				if _, ok := mapItemKey[key.Interface()]; !ok {
					return false
				}
			}
		}

		for _, key := range itemValue.MapKeys() {
			if isEqual := Compare(itemValue.MapIndex(key).Interface(), otherItemValue.MapIndex(key).Interface()); !isEqual {
				return false
			}
		}
	default:
		panic(fmt.Sprintf("does not support type %v", itemValue.Kind().String()))
	}
	return true
}

func CompareFbExternalComments(old *fbmessaging.FbExternalComment, new *fbmessaging.FbExternalComment) bool {
	if !Compare(old, new) {
		return false
	}

	oldCommentAttachment := old.ExternalAttachment
	newCommentAttachment := new.ExternalAttachment

	if (oldCommentAttachment != nil && newCommentAttachment == nil) ||
		(oldCommentAttachment == nil && newCommentAttachment != nil) {
		return false
	}

	if oldCommentAttachment == nil && newCommentAttachment == nil {
		return true
	}

	if oldCommentAttachment.Target != nil && newCommentAttachment.Target != nil &&
		oldCommentAttachment.Target.ID != newCommentAttachment.Target.ID {
		return false
	}

	return true
}

func CompareFbExternalMessages(old *fbmessaging.FbExternalMessage, new *fbmessaging.FbExternalMessage) bool {
	if !Compare(old, new) {
		return false
	}

	oldMessageAttachments := old.ExternalAttachments
	newMessageAttachments := new.ExternalAttachments

	if len(oldMessageAttachments) != len(newMessageAttachments) {
		return false
	}

	// key: attachment.ID
	mapOldMessageAttachment := make(map[string]bool)
	mapNewMessageAttachment := make(map[string]bool)

	for _, messageAttachment := range oldMessageAttachments {
		mapOldMessageAttachment[messageAttachment.ID] = true
	}

	for _, messageAttachment := range newMessageAttachments {
		mapNewMessageAttachment[messageAttachment.ID] = true
	}

	for id := range mapOldMessageAttachment {
		if _, ok := mapNewMessageAttachment[id]; !ok {
			return false
		}
	}

	return true
}
