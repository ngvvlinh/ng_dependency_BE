// +build !generator

// Code generated by generator event. DO NOT EDIT.

package shipping

func (e *FulfillmentCreatingEvent) GetTopic() string           { return "event/shipping" }
func (e *FulfillmentShippingFeeChangedEvent) GetTopic() string { return "event/shipping" }
func (e *FulfillmentUpdatingEvent) GetTopic() string           { return "event/shipping" }