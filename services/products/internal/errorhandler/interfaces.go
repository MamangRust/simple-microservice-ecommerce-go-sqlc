package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ProductQueryError interface {
	HandleRepositoryPaginationError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		fields ...zap.Field,
	) ([]*response.ProductResponse, *int, *response.ErrorResponse)
	HandleRepositoryPaginationDeleteAtError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse)
	HandleRepositorySingleError(
		err error,
		method, tracePrefix string,
		span trace.Span,
		status *string,
		errResp *response.ErrorResponse,
		fields ...zap.Field,
	) (*response.ProductResponse, *response.ErrorResponse)
}

type ProductCommandError interface {
	HandleCreateProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponse, *response.ErrorResponse)
	HandleUpdateProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponse, *response.ErrorResponse)
	HandleTrashedProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	HandleRestoreProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponseDeleteAt, *response.ErrorResponse)
	HandleDeleteProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleRestoreAllProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
	HandleDeleteAllProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse)
}
