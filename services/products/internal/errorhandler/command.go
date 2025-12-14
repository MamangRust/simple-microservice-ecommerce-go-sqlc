package errorhandler

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/domain/response"
	productserviceerror "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/errors/service"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type productCommandError struct {
	logger logger.LoggerInterface
}

func NewProductCommandError(logger logger.LoggerInterface) *productCommandError {
	return &productCommandError{
		logger: logger,
	}
}

func (p *productCommandError) HandleRepositorySingleError(
	err error,
	method, tracePrefix string,
	span trace.Span,
	status *string,
	errResp *response.ErrorResponse,
	fields ...zap.Field,
) (*response.ProductResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponse](p.logger, err, method, tracePrefix, span, status, errResp, fields...)
}

func (p *productCommandError) HandleCreateProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponse](p.logger, err, method, tracePrefix, span, status, productserviceerror.ErrFailedCreateProduct, fields...)
}

func (p *productCommandError) HandleUpdateProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponse, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponse](p.logger, err, method, tracePrefix, span, status, productserviceerror.ErrFailedUpdateProduct, fields...)
}

func (p *productCommandError) HandleTrashedProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponseDeleteAt](p.logger, err, method, tracePrefix, span, status, productserviceerror.ErrFailedTrashProduct, fields...)
}

func (p *productCommandError) HandleRestoreProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (*response.ProductResponseDeleteAt, *response.ErrorResponse) {
	return handleErrorRepository[*response.ProductResponseDeleteAt](p.logger, err, method, tracePrefix, span, status, productserviceerror.ErrFailedRestoreProduct, fields...)
}

func (p *productCommandError) HandleDeleteProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](p.logger, err, method, tracePrefix, span, status, productserviceerror.ErrFailedDeleteProductPermanent, fields...)
}

func (p *productCommandError) HandleRestoreAllProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](p.logger, err, method, tracePrefix, span, status, productserviceerror.ErrFailedRestoreAllProducts, fields...)
}

func (p *productCommandError) HandleDeleteAllProductError(err error, method, tracePrefix string, span trace.Span, status *string, fields ...zap.Field) (bool, *response.ErrorResponse) {
	return handleErrorRepository[bool](p.logger, err, method, tracePrefix, span, status, productserviceerror.ErrFailedDeleteAllProductsPermanent, fields...)
}
