package notifier_entity

func PbEntity(s string) NotifierEntity {
	return NotifierEntity(NotifierEntity_value[s])
}

func PbEntityFromInt(s int) NotifierEntity {
	return NotifierEntity(s)
}

func (s *NotifierEntity) ToModel() string {
	if s == nil {
		return ""
	}
	return NotifierEntity_name[int(*s)]
}

func (x NotifierEntity) MarshalJSON() ([]byte, error) {
	return []byte(`"` + x.String() + `"`), nil
}
