package repository

import (
	"context"
	"time"

	"github.com/sing3demons/product.product.sync/product/product/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepository struct {
	db *mongo.Collection
}

func NewProduct(db *mongo.Collection) *ProductRepository {
	return &ProductRepository{db}
}

func (r *ProductRepository) CreateProduct(document model.Products) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := r.db.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (r *ProductRepository) FindProduct(id string) (*model.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var product model.Products
	filter := bson.M{"id": id}
	err := r.db.FindOne(ctx, filter).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) FindProductById(_id primitive.ObjectID) (*model.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var product model.Products
	err := r.db.FindOne(ctx, bson.M{"_id": _id}).Decode(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}
func (r *ProductRepository) FindAllProduct(filter bson.M) ([]model.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var products []model.Products
	for cursor.Next(ctx) {
		var product model.Products
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	// if err := cursor.All(ctx, &products); err != nil {
	// 	return nil, err
	// }

	return products, nil
}

func (r *ProductRepository) GetTotalProduct(filter bson.M, ch chan int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	count, err := r.db.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	ch <- count

	return nil
}
