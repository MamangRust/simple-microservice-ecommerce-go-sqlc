package mencache

import (
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type Mencache struct {
	UserCommandCache UserCommandCache
	UserQueryCache   UserQueryCache
}

type Deps struct {
	Redis  *redis.Client
	Logger logger.LoggerInterface
}

func NewMencache(deps *Deps) *Mencache {
	cacheStore := NewCacheStore(deps.Redis, deps.Logger)

	return &Mencache{
		UserCommandCache: NewUserCommandCache(cacheStore),
		UserQueryCache:   NewUserQueryCache(cacheStore),
	}
}
