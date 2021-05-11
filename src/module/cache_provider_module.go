package module

import (
	"awesomeProject/src/cache"
	config2 "awesomeProject/src/config"
	log "awesomeProject/src/logs"
	"sync"
)

var instance *CacheProviderModule
var onceLock sync.Once

func GetCacheProviderInstance() *CacheProviderModule {
	onceLock.Do(func() {
		instance = &CacheProviderModule{}
	})
	return instance
}

type CacheProviderModule struct {
	*cache.RedisCache
}

func (cacheProviderModule *CacheProviderModule) Init() error {
	config, err := config2.GetConfigInstance().GlobalConfig()
	if err != nil {
		log.Logger.WithError(err).Error("get init config error!")
		return err
	}
	cacheProviderModule.RedisCache = cache.NewRedisCache(config.Redis)
	return nil
}

func (cacheProviderModule *CacheProviderModule) Destroy() {
	cacheProviderModule.RedisCache.Close()
}
