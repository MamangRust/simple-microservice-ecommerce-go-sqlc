package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type OrderItemQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.OrderItemResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositoryListError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.OrderItemResponse, *response.ErrorResponse)
}

type OrderQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.OrderResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		defaultErr *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.OrderResponse, *response.ErrorResponse)
}

type OrderCommandError interface {
	HandleErrorInsufficientStockTemplate(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.OrderResponse, *response.ErrorResponse)
	HandleErrorInvalidCountStockTemplate(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.OrderResponse, *response.ErrorResponse)
	HandleCreateOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponse, *response.ErrorResponse)
	HandleUpdateOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponse, *response.ErrorResponse)
	HandleTrashedOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.OrderResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllOrderError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
