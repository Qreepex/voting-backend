package service

import (
	"github.com/qreepex/voting-backend/internal/redis"
)

type CheckVoteEligibilityService struct {
	redis *redis.Redis
}

func NewCheckVoteEligibilityService(redis *redis.Redis) *CheckVoteEligibilityService {
	return &CheckVoteEligibilityService{redis}
}
