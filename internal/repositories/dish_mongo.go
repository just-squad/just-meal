package repositories

import (
	"context"
	"just-meal/internal/models"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dishMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func newDishMongoRepository(ctx context.Context, cfg MongoConfig) (*dishMongoRepository, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}

	collecion := client.Database(cfg.Database).Collection("dishes")

	return &dishMongoRepository{
		client:     client,
		collection: collecion,
	}, nil
}

func (r *dishMongoRepository) CreateDish(ctx context.Context, dish *models.Dish) error {
	_, err := r.collection.InsertOne(ctx, dish)
	return err
}

func (r *dishMongoRepository) GetDish(ctx context.Context, id uuid.UUID) (*models.Dish, error) {
	var dish models.Dish

	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&dish)

	return &dish, err
}

func (r *dishMongoRepository) UpdateDish(ctx context.Context, id uuid.UUID, dish *models.Dish) error {
	return nil
}

func (r *dishMongoRepository) DeleteDish(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (r *dishMongoRepository) GetDishesByType(ctx context.Context, mealType models.MealType) ([]*models.Dish, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"meal_type": mealType})
	if err != nil {
		return nil, err
	}

	var dishes []*models.Dish
	return dishes, cursor.All(ctx, &dishes)
}
