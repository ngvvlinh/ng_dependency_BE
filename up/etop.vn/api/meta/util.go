package meta

type AutoGeneratable interface {
	autoGenerate()
}

func AutoFill(id *UUID) {
	if id.IsZero() {
		*id = NewUUID()
	}
}
