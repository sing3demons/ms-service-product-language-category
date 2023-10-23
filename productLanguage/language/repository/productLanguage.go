package repository

import (
	"context"
	"time"

	"github.com/sing3demons/productLanguage/language/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductLanguageRepository struct {
	db *mongo.Collection
}

func NewProductLanguage(db *mongo.Collection) *ProductLanguageRepository {
	return &ProductLanguageRepository{db}
}

func (r *ProductLanguageRepository) CreateCategory(document model.ProductLanguage) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := r.db.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *ProductLanguageRepository) FindProductLanguage(id string) (*model.ProductLanguage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var productLanguage model.ProductLanguage
	filter:= bson.M{"id": id}
	err := r.db.FindOne(ctx, filter).Decode(&productLanguage)
	if err != nil {
		return nil, err
	}

	return &productLanguage, nil
}

func (r *ProductLanguageRepository) FindProductLanguageById(_id primitive.ObjectID) (*model.ProductLanguage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var category model.ProductLanguage
	err := r.db.FindOne(ctx, bson.M{"_id": _id}).Decode(&category)
	if err != nil {
		return nil, err
	}

	return &category, nil
}
func (r *ProductLanguageRepository) FindAllProductLanguage(filter bson.M) ([]model.ProductLanguage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var categories []model.ProductLanguage
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	return categories, nil
}
