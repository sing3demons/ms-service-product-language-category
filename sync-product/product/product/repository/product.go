package repository

import (
	"context"
	"fmt"
	"strings"
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
func (r *ProductRepository) FindAllProduct(doc bson.D) ([]model.Products, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{}
	for _, v := range doc {
		if v.Key == "name" {
			names := strings.Split(fmt.Sprintf("%s", v.Value), ",")

			var filterOr bson.A
			for _, name := range names {
				filterOr = append(filterOr, bson.D{{Key: "name", Value: name}})
			}
			filter = append(filter, bson.E{Key: "$or", Value: filterOr})
		}
		if v.Key == "lifecycleStatus" {
			lifecycleStatus := fmt.Sprintf("%s", v.Value)
			if lifecycleStatus != "" {
				filter = append(filter, bson.E{Key: "lifecycleStatus", Value: lifecycleStatus})
			}

		}
	}

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

func (r *ProductRepository) GetTotalProduct(doc bson.D, ch chan int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.D{}
	for _, v := range doc {
		if v.Key == "name" {
			names := strings.Split(fmt.Sprintf("%s", v.Value), ",")

			var filterOr bson.A
			for _, name := range names {
				filterOr = append(filterOr, bson.D{{Key: "name", Value: name}})
			}
			filter = append(filter, bson.E{Key: "$or", Value: filterOr})
		}
		if v.Key == "lifecycleStatus" {
			lifecycleStatus := fmt.Sprintf("%s", v.Value)
			if lifecycleStatus != "" {
				filter = append(filter, bson.E{Key: "lifecycleStatus", Value: lifecycleStatus})
			}

		}
	}
	count, err := r.db.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}
	ch <- count

	return nil
}
