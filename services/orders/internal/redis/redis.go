package mencache

import (
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	OrderQueryCache     OrderQueryCache
	OrderCommandCache   OrderCommandCache
	OrderItemQueryCache OrderItemQueryCache
}

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		OrderQueryCache:     NewOrderQueryCache(cacheStore),
		OrderCommandCache:   NewOrderCommandCache(cacheStore),
		OrderItemQueryCache: NewOrderItemQueryCache(cacheStore),
	}
}
