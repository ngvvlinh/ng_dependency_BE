package etl

import . "reflect"

// unwrap converts from Models to []*Model
func unwrap(namedSlice interface{}) interface{} {
	// input           is Models ([]*Model)
	// input.Elem.Elem is *Model
	// output          is []*Model
	elem := TypeOf(namedSlice).Elem()
	dstSlice := SliceOf(elem)
	return ValueOf(namedSlice).Convert(dstSlice).Interface()
}

// unwrapPtr converts from *Models to *[]*Model
func unwrapPtr(ptrNamedSlice interface{}) interface{} {
	// input           is *Models (*[]*Model)
	// input.Elem.Elem is *Model
	// output          is *[]*Model
	elem := TypeOf(ptrNamedSlice).Elem().Elem()
	dstSlice := PtrTo(SliceOf(elem))
	return ValueOf(ptrNamedSlice).Convert(dstSlice).Interface()
}
