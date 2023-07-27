package handler

import (
	"context"

	"github.com/redis/go-redis/v9"

	"intern2023/database"
	"intern2023/share"
)
type AuthHandler struct {
	redisClient *redis.Client
	
}
func NewService() *AuthHandler {
	redisClient, _ := database.ConnectRedisDatabase()
	return &AuthHandler{redisClient: redisClient}
}

func (ah *AuthHandler) CheckUser(UserId string, SessionId string) (bool, error) {
	if UserId == "" || SessionId == "" {
		return false, nil
	}
	// Check if the user exists in the Redis database
	IdSessionDB, err := ah.redisClient.Get(context.Background(), share.UserPatternSession(UserId)).Result()
	if err != nil || IdSessionDB == "" {
		return false, nil
	}
	if IdSessionDB != SessionId {
		return false, nil
	}
	return true, nil
}