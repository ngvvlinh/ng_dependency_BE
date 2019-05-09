package v1

type EventType = EventDataEnum

type IsEventData interface {
	isEventData_Data
	GetEnumTag() EventType
}
