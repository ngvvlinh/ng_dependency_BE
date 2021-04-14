package subject_referral

// +enum
// +enum:zero=null
type SubjectReferral int

type NullSubjectReferral struct {
	Enum  SubjectReferral
	Valid bool
}

const (
	// +enum=unknown
	Unknown SubjectReferral = 0

	// +enum=credit
	Credit SubjectReferral = 3

	// +enum=invoice
	Invoice SubjectReferral = 7

	// +enum=subscription
	Subscription SubjectReferral = 9

	// +enum=order
	Order SubjectReferral = 11
)

func (s SubjectReferral) GetCode() string {
	switch s {
	case Credit:
		return "cr"
	case Invoice:
		return "in"
	case Subscription:
		return "su"
	case Order:
		return "or"
	}
	return ""
}
