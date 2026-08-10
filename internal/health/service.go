package health

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

type Status struct {
	Service  string            `json:"service"`
	Status   string            `json:"status"`
	Version  string            `json:"version"`
	Services map[string]string `json:"services"`
}

func NewService(db *pgxpool.Pool, redisClient *redis.Client) *Service {
	return &Service{
		db:    db,
		redis: redisClient,
	}
}

func (s *Service) Check(ctx context.Context) Status {
	services := map[string]string{
		"postgres": "ok",
		"redis":    "ok",
	}

	status := "ok"

	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	if err := s.db.Ping(checkCtx); err != nil {
		services["postgres"] = "unavailable"
		status = "degraded"
	}

	if err := s.redis.Ping(checkCtx).Err(); err != nil {
		services["redis"] = "unavailable"
		status = "degraded"
	}

	return Status{
		Service:  "peakcloud-api",
		Status:   status,
		Version:  "0.1.0",
		Services: services,
	}
}
