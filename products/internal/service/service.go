package service

import (
	productresponsemapper "github.com/MamangRust/simple_microservice_ecommerce/product/internal/mapper/response"
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

func NewService(repository repository.Repositories, logger logger.LoggerInterface) Service {
	mapperQuery := productresponsemapper.NewProductQueryResponseMapper()
	mapperCommand := productresponsemapper.NewProductCommandResponseMapper()

	return &service{
		ProductQueryService:   newProductQueryService(repository.ProductQueryRepo(), logger, mapperQuery),
		ProductCommandService: newProductCommandService(repository, logger, mapperCommand),
	}
}

func newProductQueryService(
	repository repository.ProductQueryRepository,
	logger logger.LoggerInterface,
	mapper productresponsemapper.ProductQueryResponseMapper,
) ProductQueryService {
	return NewProductQueryService(&productQueryDeps{
		repository: repository,
		logger:     logger,
		mapper:     mapper,
	})
}

func newProductCommandService(
	repository repository.Repositories,
	logger logger.LoggerInterface,
	mapper productresponsemapper.ProductCommandResponseMapper,
) ProductCommandService {
	return NewProductCommandService(&productCommandDeps{
		productQueryRepository:   repository.ProductQueryRepo(),
		productCommandRepository: repository.ProductCommandRepo(),
		logger:                   logger,
		mapper:                   mapper,
	})
}
