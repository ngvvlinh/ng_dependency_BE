package subject_referral

// +enum
// +enum:zero=null
type SubjectReferral int

type NullSubjectReferral struct {
	Enum  SubjectReferral
	Valid bool
}

const (
	// +enum=credit
	Credit SubjectReferral = 3

	// +enum=invoice
	Invoice SubjectReferral = 7

	// +enum=subscription
	Subscription SubjectReferral = 9
)
