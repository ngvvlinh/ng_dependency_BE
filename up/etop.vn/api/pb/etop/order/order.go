package order

import "etop.vn/api/pb/common"

func (m *OrderWithErrorsResponse) HasErrors() []*common.Error {
	return m.FulfillmentErrors
}

func (m *ImportOrdersResponse) HasErrors() []*common.Error {
	if len(m.CellErrors) > 0 {
		return m.CellErrors
	}
	return m.ImportErrors
}
