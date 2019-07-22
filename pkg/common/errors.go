package cm

import "etop.vn/common/xerrors"

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
)

var (
	ErrTODO             = xerrors.ErrTODO
	ErrREMOVED          = xerrors.ErrREMOVED
	ErrUnauthenticated  = xerrors.ErrUnauthenticated
	ErrPermissionDenied = xerrors.ErrPermissionDenied
)

func Error(code xerrors.Code, message string, err error) *xerrors.APIError {
	return xerrors.Error(code, message, err)
}

func Errorf(code xerrors.Code, err error, message string, args ...interface{}) *xerrors.APIError {
	return xerrors.Errorf(code, err, message, args...)
}

func Trace(err error) *xerrors.APIError {
	return xerrors.Trace(err)
}

func ErrorTracef(code xerrors.Code, err error, message string, args ...interface{}) *xerrors.APIError {
	return xerrors.ErrorTracef(code, err, message, args...)
}

func NSErrorf(code xerrors.Code, err error, message string, args ...interface{}) *xerrors.APIError {
	return xerrors.NSErrorf(code, err, message, args...)
}

func ErrorCode(err error) xerrors.Code {
	return xerrors.ErrorCode(err)
}

func ToError(err error) *xerrors.APIError {
	return xerrors.ToError(err)
}

func MapError(err error) *xerrors.MapErrors {
	return xerrors.MapError(err)
}
