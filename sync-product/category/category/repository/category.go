package repository

import (
	"context"
	"time"

	"github.com/sing3demons/product.product.sync/category/category/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoryRepository struct {
	db *mongo.Collection
}

func NewCategory(db *mongo.Collection) *CategoryRepository {
	return &CategoryRepository{db}
}

func (r *CategoryRepository) CreateCategory(document model.Category) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := r.db.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *CategoryRepository) FindCategory(id string) (*model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var category model.Category
	filter := bson.M{"id": id}
	err := r.db.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) FindCategoryById(_id primitive.ObjectID) (*model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var category model.Category
	err := r.db.FindOne(ctx, bson.M{"_id": _id}).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}
func (r *CategoryRepository) FindAllCategory(filter bson.M) ([]model.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var categories []model.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}