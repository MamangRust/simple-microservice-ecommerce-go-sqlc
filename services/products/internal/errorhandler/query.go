package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	productserviceerror "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/errors/service"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type productQueryError struct {
	logger logger.LoggerInterface
}

func NewProductQueryError(logger logger.LoggerInterface) *productQueryError {
	return &productQueryError{
		logger: logger,
	}
}

func (p *productQueryError) HandleRepositoryPaginationError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	fields ...zap.Field,
) ([]*response.ProductResponse, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.ProductResponse](p.logger, err, method, tracePrefix, span, status, productserviceerror.ErrFailedFindAllProducts, fields...)
}

func (p *productQueryError) HandleRepositoryPaginationDeleteAtError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) ([]*response.ProductResponseDeleteAt, *int, *response.ErrorResponse) {
	return handleErrorPagination[[]*response.ProductResponseDeleteAt](p.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (p *productQueryError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.ProductResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponse](p.logger, err, method, tracePrefix, span, status, errResp, fields...)
}
