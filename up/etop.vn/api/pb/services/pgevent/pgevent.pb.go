package pgevent

import (
	"etop.vn/common/jsonx"
)

type GenerateEventsRequest struct {
	// Be careful, as invalid message can cause errors in other services!
	//
	// Format: `table:rid:op:id`
	// Example: `fulfillment:81460:INSERT:1052886332676874435`
	RawEvents []string `json:"raw_events"`
	// Example: `{fulfillment:81460:INSERT:1052886332676874435,fulfillment:81461:UPDATE:1052886332676874435}`
	RawEventsPg string `json:"raw_events_pg"`
	// Control how many events are dispatched asynchronously as a group.
	// We don't want to dispatch too many events together.
	// Example: 100. Default: 0 - events are dispatched synchronously.
	ItemsPerBatch int32 `json:"items_per_batch"`
}

func (m *GenerateEventsRequest) Reset()         { *m = GenerateEventsRequest{} }
func (m *GenerateEventsRequest) String() string { return jsonx.MustMarshalToString(m) }
