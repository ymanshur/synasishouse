package util

import (
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCCodeToHTTPStatus converts a gRPC codes.Code into a corresponding HTTP status code.
// This mapping follows the general guidelines recommended for gRPC-to-HTTP transcoding.
func GRPCCodeToHTTPStatus(code codes.Code) (int, bool) {
	switch code {
	case codes.OK:
		return http.StatusOK, true // 200: Everything is fine
	case codes.Canceled:
		return 499, true // 499: Client Closed Request (Non-standard but common for Canceled)
	case codes.Unknown:
		return http.StatusInternalServerError, true // 500: Generic server error
	case codes.InvalidArgument:
		return http.StatusBadRequest, true // 400: The client specified an invalid argument
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout, true // 504: The request took too long
	case codes.NotFound:
		return http.StatusNotFound, true // 404: Resource not found
	case codes.AlreadyExists:
		return http.StatusConflict, true // 409: Resource already exists (often used for unique constraint violations)
	case codes.PermissionDenied:
		return http.StatusForbidden, true // 403: The client is authenticated but lacks permission
	case codes.Unauthenticated:
		return http.StatusUnauthorized, true // 401: Authentication credentials are missing or invalid
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests, true // 429: Rate limit or quota exceeded
	case codes.FailedPrecondition:
		return http.StatusBadRequest, true // 400: State change failed (e.g., trying to delete a non-empty directory)
	case codes.Aborted:
		return http.StatusConflict, true // 409: Concurrency error (e.g., transaction conflict)
	case codes.OutOfRange:
		return http.StatusBadRequest, true // 400: An index was out of bounds
	case codes.Unimplemented:
		return http.StatusNotImplemented, true // 501: Method is not implemented by server
	case codes.Internal:
		return http.StatusInternalServerError, true // 500: Server-side crash or serious error
	case codes.Unavailable:
		return http.StatusServiceUnavailable, true // 503: Service is temporarily unavailable (e.g., overload)
	case codes.DataLoss:
		return http.StatusInternalServerError, true // 500: Unrecoverable data loss
	default:
		return http.StatusInternalServerError, true
	}
}

// TranslateGRPCError is a utility function that extracts the gRPC status code
// from a standard Go error (if it's a gRPC status.Status) and converts it
// to an HTTP status code.
// It defaults to 500 for non-gRPC errors or unknown gRPC errors.
func TranslateGRPCError(err error) (int, bool) {
	if err == nil {
		return http.StatusOK, false
	}

	s, ok := status.FromError(err)
	if ok {
		return GRPCCodeToHTTPStatus(s.Code())
	}

	return http.StatusInternalServerError, false
}
