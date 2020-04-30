package convert

import (
	_ "o.o/api/main/connectioning"
	_ "o.o/backend/com/main/connectioning/model"
)

// +gen:convert: o.o/backend/com/main/connectioning/model -> o.o/api/main/connectioning
// +gen:convert: o.o/api/main/connectioning
