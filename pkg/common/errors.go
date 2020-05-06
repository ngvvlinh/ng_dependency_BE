package cm

import "o.o/common/xerrors"

const (
	Canceled           = xerrors.Canceled
	Unknown            = xerrors.Unknown
	InvalidArgument    = xerrors.InvalidArgument
	DeadlineExceeded   = xerrors.DeadlineExceeded
	NotFound           = xerrors.NotFound
	AlreadyExists      = xerrors.AlreadyExists
	PermissionDenied   = xerrors.PermissionDenied
	Unauthenticated    = xerrors.Unauthenticated
	ResourceExhausted  = xerrors.ResourceExhausted
	FailedPrecondition = xerrors.FailedPrecondition
	Aborted            = xerrors.Aborted
	OutOfRange         = xerrors.OutOfRange
	Unimplemented      = xerrors.Unimplemented
	Internal           = xerrors.Internal
	Unavailable        = xerrors.Unavailable
	DataLoss           = xerrors.DataLoss
	NoError            = xerrors.NoError
	OK                 = xerrors.NoError

	RuntimePanic         = xerrors.RuntimePanic
	WrongPassword        = xerrors.WrongPassword
	ValidationFailed     = xerrors.ValidationFailed
	STokenRequired       = xerrors.STokenRequired
	CaptchaRequired      = xerrors.CaptchaRequired
	CaptchaInvalid       = xerrors.CaptchaInvalid
	RegisterRequired     = xerrors.RegisterRequired
	ExternalServiceError = xerrors.ExternalServiceError
	AccountClosed        = xerrors.AccountClosed

	SkipSync = xerrors.SkipSync

	FacebookPermissionDenied = xerrors.FacebookPermissionDenied
	FacebookError            = xerrors.Facebook
)

var (
	ErrTODO             = xerrors.ErrTODO
	ErrREMOVED          = xerrors.ErrREMOVED
	ErrUnauthenticated  = xerrors.ErrUnauthenticated
	ErrPermissionDenied = xerrors.ErrPermissionDenied
)

var (
	Error       = xerrors.Error
	Errorf      = xerrors.Errorf
	Trace       = xerrors.Trace
	ErrorTracef = xerrors.ErrorTracef
	NSErrorf    = xerrors.NSErrorf
	ErrorCode   = xerrors.ErrorCode
	ToError     = xerrors.ToError
	MapError    = xerrors.MapError
)
