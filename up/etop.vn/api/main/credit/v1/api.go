package v1

type EventType = EventDataEnum

type IsEventData interface {
	isEventData_Data
	GetEnumTag() EventType
}

func (m *CreditEvent) SetEventData(data IsEventData) {
	if data == nil {
		m.Type = 0
		m.Data = nil
		return
	}
	m.Type = int32(data.GetEnumTag())
	m.Data = &EventData{Data: data}
	if m.Type == 0 {
		panic("unexpected")
	}
}

func (m *CreditEvent) GetEventData() (EventType, IsEventData) {
	data := m.Data.Data.(IsEventData)
	return data.GetEnumTag(), m.Data.GetData().(IsEventData)
}
