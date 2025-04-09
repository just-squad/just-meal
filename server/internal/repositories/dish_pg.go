package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"just-meal-api/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const databaseName string = "dishes"

// Ошибки репозитория
var (
	ErrNotFound = fmt.Errorf("not found")
)

type dishPostgresRepository struct {
	connPool *pgxpool.Pool
}

func newDishPostgresRepository(ctx context.Context, cfg PgConfig) (*dishPostgresRepository, error) {
	config, err := pgxpool.ParseConfig(cfg.GetConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres connection config: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// проверка соединения
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &dishPostgresRepository{connPool: pool}, nil
}

func (r *dishPostgresRepository) CreateDish(ctx context.Context, dish *models.Dish) error {
	ingredientsJson, err := json.Marshal(dish.Ingredients)
	if err != nil {
		return fmt.Errorf("failed to marshal ingredients: %w", err)
	}
	nutritionJson, err := json.Marshal(dish.Nutrition)
	if err != nil {
		return fmt.Errorf("failed to marshal nutrition: %w", err)
	}
	recipeJson, err := json.Marshal(dish.Recipe)
	if err != nil {
		return fmt.Errorf("failed to marshal recipe: %w", err)
	}

	sqlString := fmt.Sprintf(`insert into %s(id, name, recipe, ingredients, nutrition, meal_type, cooking_time, servings)
	values ($1, $2, $3, $4, $5, $6, $7, $8)`, databaseName)

	_, err = r.connPool.Exec(ctx,
		sqlString,
		dish.Id,
		dish.Name,
		recipeJson,
		ingredientsJson,
		nutritionJson,
		dish.MealType,
		dish.CookingTime,
		dish.Servings,
	)
	if err != nil {
		return fmt.Errorf("failed to insert dish: %w", err)
	}

	return nil
}

func (r *dishPostgresRepository) GetDish(ctx context.Context, id uuid.UUID) (*models.Dish, error) {
	var dish models.Dish
	var ingredients, nutrition, recipe []byte

	sqlString := fmt.Sprintf(`select id, name, recipe, ingredients, nutrition, meal_type, cooking_time, servings from %s where id = $1`, databaseName)

	err := r.connPool.QueryRow(
		ctx,
		sqlString,
		id,
	).Scan(
		&dish.Id,
		&dish.Name,
		&recipe,
		&ingredients,
		&nutrition,
		&dish.MealType,
		&dish.CookingTime,
		&dish.Servings,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get dish: %w", err)
	}

	if err := json.Unmarshal(ingredients, &dish.Ingredients); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ingredients: %w", err)
	}
	if err := json.Unmarshal(nutrition, &dish.Nutrition); err != nil {
		return nil, fmt.Errorf("failed to unmarshal nutrition: %w", err)
	}
	if err := json.Unmarshal(recipe, &dish.Recipe); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recipe: %w", err)
	}

	return &dish, nil
}

func (r *dishPostgresRepository) UpdateDish(ctx context.Context, id uuid.UUID, dish *models.Dish) error {
	ingredientsJson, err := json.Marshal(dish.Ingredients)
	if err != nil {
		return fmt.Errorf("failed to marshal ingredients: %w", err)
	}
	nutritionJson, err := json.Marshal(dish.Nutrition)
	if err != nil {
		return fmt.Errorf("failed to marshal nutrition: %w", err)
	}
	recipeJson, err := json.Marshal(dish.Recipe)
	if err != nil {
		return fmt.Errorf("failed to marshal recipe: %w", err)
	}

	sql := fmt.Sprintf(`update %s set 
		name $2,
		recipe = $3,
		ingredients = $4,
		nutrition = $5,
		meal_type = $6,
		cooking_time = $7,
		servings = $8
	where id = $1`, databaseName)

	cmd, err := r.connPool.Exec(ctx,
		sql,
		id,
		dish.Name,
		recipeJson,
		ingredientsJson,
		nutritionJson,
		dish.MealType,
		dish.CookingTime,
		dish.Servings,
	)
	if err != nil {
		return fmt.Errorf("failed to update dish: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *dishPostgresRepository) DeleteDish(ctx context.Context, id uuid.UUID) error {
	sql := fmt.Sprintf(`delete from %s where id = $1`, databaseName)

	cmd, err := r.connPool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete dish: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *dishPostgresRepository) GetDishesByType(ctx context.Context, mealType models.MealType) ([]*models.Dish, error) {
	sql := fmt.Sprintf(`select id, name, recipe, ingredients, nutrition, meal_type, cooking_time, servings from %s
	where meal_type = $1`, databaseName)

	rows, err := r.connPool.Query(ctx,
		sql,
		mealType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query dishes: %w", err)
	}
	defer rows.Close()

	var dishes []*models.Dish
	for rows.Next() {
		var dish models.Dish
		var ingredients, nutrition, recipe []byte

		err := rows.Scan(
			&dish.Id,
			&dish.Name,
			&recipe,
			&ingredients,
			&nutrition,
			&dish.MealType,
			&dish.CookingTime,
			&dish.Servings,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan dish row: %w", err)
		}
		if err := json.Unmarshal(ingredients, &dish.Ingredients); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ingredients: %w", err)
		}
		if err := json.Unmarshal(nutrition, &dish.Nutrition); err != nil {
			return nil, fmt.Errorf("failed to unmarshal nutrition: %w", err)
		}
		if err := json.Unmarshal(recipe, &dish.Recipe); err != nil {
			return nil, fmt.Errorf("failed to unmarshal recipe: %w", err)
		}
		dishes = append(dishes, &dish)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return nil, nil
}
