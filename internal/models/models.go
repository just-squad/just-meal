package models

import (
	"github.com/google/uuid"
)

type MealType string
type UnitOfMeasurement string
type RecipeItem string
type Tag string

// Тип приема пищи
const (
	Breakfast MealType = "breakfast"
	Brunch    MealType = "branch"
	Lunch     MealType = "lunch"
	Dinner    MealType = "dinner"
	Supper    MealType = "supper"
	Snack     MealType = "snack"
)

const (
	Item  UnitOfMeasurement = "unit"
	Gram  UnitOfMeasurement = "gram"
	Liter UnitOfMeasurement = "liter"
)

type Dish struct {
	Id          uuid.UUID    `json:"id" bson:"id"`                     // Уникальный идентификатор
	Name        string       `json:"name" bson:"name"`                 // Название блюда
	Recipe      []RecipeItem `json:"recipe" bson:"recipe"`             // Текст рецепта
	Ingredients []Ingredient `json:"ingredients" bson:"ingredients"`   // Список ингредиентов
	Nutrition   Nutrition    `json:"nutrition" bson:"nutrition"`       // КБЖУ
	MealType    MealType     `json:"meal_type" bson:"meal_type"`       // Тип блюда
	CookingTime int32        `json:"cooking_time" bson:"cooking_time"` // Время готовки в минутах
	Servings    int32        `json:"servings" bson:"servings"`         // Количество порций
	Tag         []Tag        `json:"tag" bson:"tag"`                   // теги для блюда
}

// Пищевая ценность КБЖУ
type Nutrition struct {
	Calories float64 `json:"calories"` // Калории
	Protein  float64 `json:"protein"`  // Белки (г)
	Fat      float64 `json:"fat"`      // Жиры (г)
	Carbs    float64 `json:"carbs"`    // Углеводы (г)
}

// Ингредиент блюда
type Ingredient struct {
	Name     string  `json:"name"`     // Название ингредиента
	Quantity float64 `json:"quantity"` // Количество (например: "200г", "1 шт.")
}
