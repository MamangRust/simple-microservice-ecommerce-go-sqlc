package mencache

import (
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	UserCache      UserCache
	RoleCache      RoleCache
	ProductCache   ProductCache
	OrderCache     OrderCache
	OrderItemCache OrderItemCache
}

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		UserCache:      NewUserCache(cacheStore),
		RoleCache:      NewRoleCache(cacheStore),
		ProductCache:   NewProductCache(cacheStore),
		OrderCache:     NewOrderCache(cacheStore),
		OrderItemCache: NewOrderItemCache(cacheStore),
	}
}
