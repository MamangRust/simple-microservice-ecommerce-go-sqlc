package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/domain/response"
	orderserviceerrors "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/errors/order_errors/service"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type orderQueryError struct {
	logger logger.LoggerInterface
}

func NewOrderQueryError(logger logger.LoggerInterface) *orderQueryError {
	return &orderQueryError{
		logger: logger,
	}
}

func (o *orderQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.OrderResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.OrderResponse](o.logger, err, method, tracePrefix, span, status, orderserviceerrors.ErrFailedFindAllOrders, fields...)
}

func (o *orderQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.OrderResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.OrderResponseDeleteAt](o.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (e *orderQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	defaultErr *response.ErrorResponse,
	fields ...zap.Field,
) (*response.OrderResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.OrderResponse](e.logger, err, method, tracePrefix, span, status, defaultErr, fields...)
}
