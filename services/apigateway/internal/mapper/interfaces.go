package mapper

import (
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/domain/response"
	pbauth "github.com/MamangRust/simple_microservice_ecommerce_pb/auth"
	pborder "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
	pborder_item "github.com/MamangRust/simple_microservice_ecommerce_pb/order_item"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
)

type AuthResponseMapper interface {
	ToResponseVerifyCode(res *pbauth.ApiResponseVerifyCode) *response.ApiResponseVerifyCode
	ToResponseForgotPassword(res *pbauth.ApiResponseForgotPassword) *response.ApiResponseForgotPassword
	ToResponseResetPassword(res *pbauth.ApiResponseResetPassword) *response.ApiResponseResetPassword
	ToResponseLogin(res *pbauth.ApiResponseLogin) *response.ApiResponseLogin
	ToResponseRegister(res *pbauth.ApiResponseRegister) *response.ApiResponseRegister
	ToResponseRefreshToken(res *pbauth.ApiResponseRefreshToken) *response.ApiResponseRefreshToken
	ToResponseGetMe(res *pbauth.ApiResponseGetMe) *response.ApiResponseGetMe
}

type RoleResponseMapper interface {
	ToApiResponseRoleAll(pbResponse *pbrole.ApiResponseRoleAll) *response.ApiResponseRoleAll
	ToApiResponseRoleDelete(pbResponse *pbrole.ApiResponseRoleDelete) *response.ApiResponseRoleDelete
	ToApiResponseRole(pbResponse *pbrole.ApiResponseRole) *response.ApiResponseRole
	ToApiResponseRoleDeleteAt(pbResponse *pbrole.ApiResponseRoleDeleteAt) *response.ApiResponseRoleDeleteAt
	ToApiResponsesRole(pbResponse *pbrole.ApiResponsesRole) *response.ApiResponsesRole
	ToApiResponsePaginationRole(pbResponse *pbrole.ApiResponsePaginationRole) *response.ApiResponsePaginationRole
	ToApiResponsePaginationRoleDeleteAt(pbResponse *pbrole.ApiResponsePaginationRoleDeleteAt) *response.ApiResponsePaginationRoleDeleteAt
}

type UserResponseMapper interface {
	ToApiResponseUserDeleteAt(pbResponse *pbuser.ApiResponseUserDeleteAt) *response.ApiResponseUserDeleteAt
	ToApiResponseUser(pbResponse *pbuser.ApiResponseUser) *response.ApiResponseUser

	ToApiResponseUserDelete(pbResponse *pbuser.ApiResponseUserDelete) *response.ApiResponseUserDelete
	ToApiResponseUserAll(pbResponse *pbuser.ApiResponseUserAll) *response.ApiResponseUserAll
	ToApiResponsePaginationUserDeleteAt(pbResponse *pbuser.ApiResponsePaginationUserDeleteAt) *response.ApiResponsePaginationUserDeleteAt
	ToApiResponsePaginationUser(pbResponse *pbuser.ApiResponsePaginationUser) *response.ApiResponsePaginationUser
}

type OrderResponseMapper interface {
	ToApiResponseOrder(pbResponse *pborder.ApiResponseOrder) *response.ApiResponseOrder
	ToApiResponseOrderDeleteAt(pbResponse *pborder.ApiResponseOrderDeleteAt) *response.ApiResponseOrderDeleteAt
	ToApiResponseOrderDelete(pbResponse *pborder.ApiResponseOrderDelete) *response.ApiResponseOrderDelete
	ToApiResponseOrderAll(pbResponse *pborder.ApiResponseOrderAll) *response.ApiResponseOrderAll
	ToApiResponsePaginationOrderDeleteAt(pbResponse *pborder.ApiResponsePaginationOrderDeleteAt) *response.ApiResponsePaginationOrderDeleteAt
	ToApiResponsePaginationOrder(pbResponse *pborder.ApiResponsePaginationOrder) *response.ApiResponsePaginationOrder
}

type OrderItemResponseMapper interface {
	ToApiResponseOrderItem(pbResponse *pborder_item.ApiResponseOrderItem) *response.ApiResponseOrderItem
	ToApiResponsesOrderItem(pbResponse *pborder_item.ApiResponsesOrderItem) *response.ApiResponsesOrderItem
	ToApiResponseOrderItemDelete(pbResponse *pborder_item.ApiResponseOrderItemDelete) *response.ApiResponseOrderItemDelete
	ToApiResponseOrderItemAll(pbResponse *pborder_item.ApiResponseOrderItemAll) *response.ApiResponseOrderItemAll
	ToApiResponsePaginationOrderItemDeleteAt(pbResponse *pborder_item.ApiResponsePaginationOrderItemDeleteAt) *response.ApiResponsePaginationOrderItemDeleteAt
	ToApiResponsePaginationOrderItem(pbResponse *pborder_item.ApiResponsePaginationOrderItem) *response.ApiResponsePaginationOrderItem
}

type ProductResponseMapper interface {
	ToApiResponseProduct(pbResponse *pbproduct.ApiResponseProduct) *response.ApiResponseProduct
	ToApiResponsesProductDeleteAt(pbResponse *pbproduct.ApiResponseProductDeleteAt) *response.ApiResponseProductDeleteAt
	ToApiResponseProductDelete(pbResponse *pbproduct.ApiResponseProductDelete) *response.ApiResponseProductDelete
	ToApiResponseProductAll(pbResponse *pbproduct.ApiResponseProductAll) *response.ApiResponseProductAll
	ToApiResponsePaginationProductDeleteAt(pbResponse *pbproduct.ApiResponsePaginationProductDeleteAt) *response.ApiResponsePaginationProductDeleteAt
	ToApiResponsePaginationProduct(pbResponse *pbproduct.ApiResponsePaginationProduct) *response.ApiResponsePaginationProduct
}
