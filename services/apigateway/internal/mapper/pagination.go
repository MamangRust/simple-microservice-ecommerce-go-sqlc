package mapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
)

func mapPaginationMeta(s *pb.Pagination) *response.PaginationMeta {
	return &response.PaginationMeta{
		CurrentPage:  int(s.CurrentPage),
		PageSize:     int(s.PageSize),
		TotalRecords: int(s.TotalRecords),
		TotalPages:   int(s.TotalPages),
	}
}
