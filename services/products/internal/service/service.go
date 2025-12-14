package service

import (
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/errorhandler"
	productresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/response"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/product/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
)

type Service interface {
	ProductQueryService
	ProductCommandService
}

type service struct {
	ProductQueryService
	ProductCommandService
}

type Deps struct {
	Repositories repository.Repositories
	Logger       logger.LoggerInterface
	Mencache     *mencache.Mencache
}

func NewService(deps *Deps) Service {
	mapperQuery := productresponsemapper.NewProductQueryResponseMapper()
	mapperCommand := productresponsemapper.NewProductCommandResponseMapper()
	errorhandler := errorhandler.NewErrorHandler(deps.Logger)

	return &service{
		ProductQueryService:   newProductQueryService(deps.Repositories.ProductQueryRepo(), deps.Logger, mapperQuery, deps.Mencache.ProductQuery, errorhandler.ProductQueryError),
		ProductCommandService: newProductCommandService(deps.Repositories, deps.Logger, mapperCommand, deps.Mencache.ProductCommand, errorhandler.ProductCommandError),
	}
}

func newProductQueryService(
	repository repository.ProductQueryRepository,
	logger logger.LoggerInterface,
	mapper productresponsemapper.ProductQueryResponseMapper,
	mencache mencache.ProductQueryCache,
	errorhandler errorhandler.ProductQueryError,
) ProductQueryService {
	return NewProductQueryService(&productQueryDeps{
		repository:   repository,
		logger:       logger,
		mapper:       mapper,
		errorhandler: errorhandler,
		mencache:     mencache,
	})
}

func newProductCommandService(
	repository repository.Repositories,
	logger logger.LoggerInterface,
	mapper productresponsemapper.ProductCommandResponseMapper,
	mencache mencache.ProductCommandCache,
	errorhandler errorhandler.ProductCommandError,
) ProductCommandService {
	return NewProductCommandService(&productCommandDeps{
		productQueryRepository:   repository.ProductQueryRepo(),
		productCommandRepository: repository.ProductCommandRepo(),
		logger:                   logger,
		mapper:                   mapper,
		mencache:                 mencache,
		errorhandler:             errorhandler,
	})
}
