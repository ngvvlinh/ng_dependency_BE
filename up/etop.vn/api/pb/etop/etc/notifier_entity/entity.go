package notifier_entity

func PbEntity(s string) NotifierEntity {
	return NotifierEntity(NotifierEntity_value[s])
}

func PbEntityFromInt(s int) NotifierEntity {
	return NotifierEntity(s)
}

func (x *NotifierEntity) ToModel() string {
	if x == nil {
		return ""
	}
	return NotifierEntity_name[int(*x)]
}

func (x NotifierEntity) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
