package httprpc

import "o.o/common/xerrors"

func ServerHTTPStatusFromErrorCode(code xerrors.Code) int {
	switch code {
	case xerrors.Canceled:
		return 408 // RequestTimeout
	case xerrors.Unknown:
		return 500 // Internal Server Error
	case xerrors.InvalidArgument:
		return 400 // BadRequest
	case xerrors.Malformed:
		return 400 // BadRequest
	case xerrors.DeadlineExceeded:
		return 408 // RequestTimeout
	case xerrors.NotFound:
		return 404 // Not Found
	case xerrors.BadRoute:
		return 404 // Not Found
	case xerrors.AlreadyExists:
		return 409 // Conflict
	case xerrors.PermissionDenied:
		return 403 // Forbidden
	case xerrors.Unauthenticated:
		return 401 // Unauthorized
	case xerrors.ResourceExhausted:
		return 403 // Forbidden
	case xerrors.FailedPrecondition:
		return 412 // Precondition Failed
	case xerrors.Aborted:
		return 409 // Conflict
	case xerrors.OutOfRange:
		return 400 // Bad Request
	case xerrors.Unimplemented:
		return 501 // Not Implemented
	case xerrors.Internal:
		return 500 // Internal Server Error
	case xerrors.Unavailable:
		return 503 // Service Unavailable
	case xerrors.DataLoss:
		return 500 // Internal Server Error
	case xerrors.NoError:
		return 200 // OK
	default:
		return 0 // Invalid!
	}
}
