package service

import "github.com/0xsj/gin-sqlc/log"

type CachedLinkMetadataService struct {
	baseService LinkMetadataService
	cache       cache.CacheService
	keyBuilder  *cache.CacheKeyBuilder
	logger      log.Logger
}