package protomapper

import pb "github.com/MamangRust/simple_microservice_ecommerce_pb"

func MapPaginationMeta(s *pb.Pagination) *pb.Pagination {
	return &pb.Pagination{
		CurrentPage:  int32(s.CurrentPage),
		PageSize:     int32(s.PageSize),
		TotalPages:   int32(s.TotalPages),
		TotalRecords: int32(s.TotalRecords),
	}
}
