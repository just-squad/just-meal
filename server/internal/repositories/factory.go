package repositories

import (
	"context"
	"fmt"
)

func NewDishRepository(ctx context.Context, cfg DBConfig) (DishRepository, error) {
	switch cfg.Type {
	case Postgres:
		return newDishPostgresRepository(ctx, cfg.Postgres)
	case Mongo:
		return newDishMongoRepository(ctx, cfg.Mongo)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}
}
