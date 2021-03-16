package user

import (
	"github.com/gomodule/redigo/redis"
	"gorm.io/gorm"
)

type repository struct {
	db        *gorm.DB
	redisPool *redis.Pool
}

type Repository interface {
}

func NewRepository(db *gorm.DB, redis *redis.Pool) Repository {
	return &repository{
		db:        db,
		redisPool: redis,
	}
}
