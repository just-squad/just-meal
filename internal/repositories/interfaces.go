package repositories

import (
	"context"
	"just-meal/internal/models"

	"github.com/google/uuid"
)

// Интерфейс репозитория для управления сущностью блюдо
type DishRepository interface {
	CreateDish(ctx context.Context, dish *models.Dish) error                               // Создать новое блюдо
	GetDish(ctx context.Context, id uuid.UUID) (*models.Dish, error)                       // Получить блюдо по идентификатору
	UpdateDish(ctx context.Context, id uuid.UUID, dish *models.Dish) error                 // Обновить данные о блюде
	DeleteDish(ctx context.Context, id uuid.UUID) error                                    // Удалить существующее блюдо
	GetDishesByType(ctx context.Context, mealType models.MealType) ([]*models.Dish, error) // Получить список блюд по типам
}
