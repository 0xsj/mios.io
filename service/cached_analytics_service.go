package service

import "github.com/0xsj/gin-sqlc/log"

type CachedAnalyticsService struct {
	baseService AnalyticsService
	cache       cache.CacheService
	keyBuilder  *cache.CacheKeyBuilder
	logger      log.Logger
}