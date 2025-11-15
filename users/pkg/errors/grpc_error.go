package errors

import (
	"encoding/json"

	pb "github.com/MamangRust/simple_microservice_ecommerce_pb"
)

func GrpcErrorToJson(err *pb.ErrorResponse) string {
	jsonData, _ := json.Marshal(err)
	return string(jsonData)
}
