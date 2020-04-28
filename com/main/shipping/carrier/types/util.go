package types

import (
	"strings"

	shipmodel "o.o/backend/com/main/shipping/model"
)

func GetShippingProviderNote(ffm *shipmodel.Fulfillment) string {
	noteB := strings.Builder{}
	if note := ffm.AddressFrom.Notes.GetFullNote(); note != "" {
		noteB.WriteString("Lấy hàng: ")
		noteB.WriteString(note)
		noteB.WriteString("\n")
	}
	if note := ffm.AddressTo.Notes.GetFullNote(); note != "" || ffm.ShippingNote != "" {
		noteB.WriteString("Giao hàng: ")
		if ffm.ShippingNote != "" {
			noteB.WriteString(ffm.ShippingNote)
			noteB.WriteString(". \n")
		}
		noteB.WriteString(note)
		noteB.WriteString("\n")
	}
	noteB.WriteString("Giao hàng không thành công hoặc giao một phần, xin gọi lại cho shop. KHÔNG ĐƯỢC TỰ Ý HOÀN HÀNG khi chưa thông báo cho shop.")
	return noteB.String()
}

func Blacklist(current byte, arbitraryCharacter byte, blacks ...byte) byte {
	for _, b := range blacks {
		if current == b {
			// return an arbitrary character which does not collide with blacklist values
			return arbitraryCharacter
		}
	}
	return current
}
