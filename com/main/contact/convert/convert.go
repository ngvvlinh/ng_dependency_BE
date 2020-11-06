package convert

import (
	"o.o/api/main/contact"
	cm "o.o/backend/pkg/common"
)

// +gen:convert: o.o/backend/com/main/contact/model  -> o.o/api/main/contact
// +gen:convert: o.o/api/main/contact

func createContact(args *contact.CreateContactArgs, out *contact.Contact) {
	apply_contact_CreateContactArgs_contact_Contact(args, out)
	out.ID = cm.NewID()
}
