package response

import (
	"encoding/json"

	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/errors"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewErrorResponse(message string, code int) *ErrorResponse {
	return &ErrorResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}
}

func mapHttpToGrpcCode(httpCode int) codes.Code {
	switch httpCode {
	case 400:
		return codes.InvalidArgument
	case 401:
		return codes.Unauthenticated
	case 403:
		return codes.PermissionDenied
	case 404:
		return codes.NotFound
	case 409:
		return codes.Aborted
	case 422:
		return codes.FailedPrecondition
	case 429:
		return codes.ResourceExhausted
	case 500:
		return codes.Internal
	default:
		return codes.Unknown
	}
}

func (e *ErrorResponse) ToGRPCError() error {
	if e == nil {
		return nil
	}

	payload := &pb.ErrorResponse{
		Status:  e.Status,
		Message: e.Message,
		Code:    int32(e.Code),
	}

	jsonMsg, _ := json.Marshal(payload)

	return status.Error(mapHttpToGrpcCode(e.Code), string(jsonMsg))
}

func NewApiErrorResponse(statusText, message string, code int) *ErrorResponse {
	return &ErrorResponse{
		Status:  statusText,
		Message: message,
		Code:    code,
	}
}

func NewGrpcError(statusText string, message string, code int) error {
	return status.Errorf(codes.Code(code),
		"%s", errors.GrpcErrorToJson(&pb.ErrorResponse{
			Status:  statusText,
			Message: message,
			Code:    int32(code),
		}),
	)
}
