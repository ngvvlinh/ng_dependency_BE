// +build !generator

// Code generated by generator event. DO NOT EDIT.

package invoice

func (e *InvoiceDeletedEvent) GetTopic() string { return "event/invoice" }
func (e *InvoicePaidEvent) GetTopic() string    { return "event/invoice" }
func (e *InvoicePayingEvent) GetTopic() string  { return "event/invoice" }