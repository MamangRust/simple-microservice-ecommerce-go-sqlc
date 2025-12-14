package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	orderitemserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/errors/orderitem_errors/service"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderItemQueryError struct {
	logger logger.LoggerInterface
}

func NewOrderItemQueryError(logger logger.LoggerInterface) *orderItemQueryError {
	return &orderItemQueryError{
		logger: logger,
	}
}

func (o *orderItemQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.OrderItemResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.OrderItemResponse](o.logger, err, method, tracePrefix, span, status, orderitemserviceerrors.ErrFailedFindAllOrderItems, fields...)
}

func (o *orderItemQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.OrderItemResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.OrderItemResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (o *orderItemQueryError) HandleRepositoryListError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.OrderItemResponse, *response.ErrorResponse) {
	return handleErrorRepository[[]*response.OrderItemResponse](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
