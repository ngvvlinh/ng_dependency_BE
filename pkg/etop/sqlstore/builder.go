package sqlstore

type sqlUpdateBuilder struct {
	b    []byte
	args []interface{}

	prev  bool
	where bool
}

func sqlUpdate(table string) *sqlUpdateBuilder {
	b := make([]byte, 0, 1024)
	b = append(b, `UPDATE "`...)
	b = append(b, table...)
	b = append(b, `" SET`...)
	return &sqlUpdateBuilder{
		b:    b,
		args: make([]interface{}, 0, 32),
	}
}

func (s *sqlUpdateBuilder) ColIndex(col string, index int, val interface{}) *sqlUpdateBuilder {
	if s.where {
		panic("Unexpected")
	}
	if s.prev {
		s.b = append(s.b, ',')
	}
	s.prev = true
	s.b = append(s.b, ' ')
	s.b = append(s.b, col...)
	s.b = append(s.b, '[')
	s.b = appendInt(s.b, index)
	s.b = append(s.b, `] = ?`...)
	s.args = append(s.args, val)
	return s
}

func (s *sqlUpdateBuilder) Where(cond string, args ...interface{}) *sqlUpdateBuilder {
	if s.where {
		s.b = append(s.b, ` AND `...)
	} else {
		s.b = append(s.b, ` WHERE `...)
	}
	s.where = true
	s.b = append(s.b, cond...)
	s.args = append(s.args, args...)
	return s
}

func (s *sqlUpdateBuilder) Build() (string, []interface{}) {
	return string(s.b), s.args
}
