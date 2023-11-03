package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sing3demons/product.product.sync/category/category/model"
	"github.com/sing3demons/product.product.sync/utils"
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
	category := model.Category{}
	filter := bson.M{"id": id}
	err := r.db.FindOne(ctx, filter).Decode(&category)
	if err != nil {
		return nil, err
	}

	if category.LastUpdate != "" {
		category.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
	}

	if category.ValidFor.StartDateTime != "" {
		category.ValidFor.StartDateTime = utils.ConvertTimeBangkok(category.ValidFor.StartDateTime)
	}
	if category.ValidFor.EndDateTime != "" {
		category.ValidFor.EndDateTime = utils.ConvertTimeBangkok(category.ValidFor.EndDateTime)
	}

	if category.Products != nil {
		var products []model.ProductRef
		for _, v := range category.Products {
			products = append(products, model.ProductRef{
				ID:         v.ID,
				Href:       utils.Href(v.Type, v.ID),
				Type:       v.Type,
				Name:       v.Name,
				Version:    v.Version,
				LastUpdate: utils.ConvertTimeBangkok(v.LastUpdate),
			})
		}
		category.Products = products
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
func (r *CategoryRepository) FindAllCategory(doc bson.D) ([]model.Category, error) {
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
			// filter = bson.D{{Key: "$or", Value: filterOr}}
			filter = append(filter, bson.E{Key: "$or", Value: filterOr})
		}
		if v.Key == "lifecycleStatus" {
			// lifecycleStatus=Active
			lifecycleStatus := fmt.Sprintf("%s", v.Value)
			if lifecycleStatus != "" {
				// filter = bson.D{{Key: "lifecycleStatus", Value: lifecycleStatus}}
				filter = append(filter, bson.E{Key: "lifecycleStatus", Value: lifecycleStatus})
			}

		}
	}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	categories := []model.Category{}

	for cursor.Next(ctx) {
		var category model.Category
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}

		validFor := &model.ValidFor{
			StartDateTime: utils.ConvertTimeBangkok(category.ValidFor.StartDateTime),
			EndDateTime:   utils.ConvertTimeBangkok(category.ValidFor.EndDateTime),
		}

		categories = append(categories, model.Category{
			Type:            category.Type,
			ID:              category.ID,
			Href:            utils.Href(category.Type, category.ID),
			Name:            category.Name,
			Version:         category.Version,
			LastUpdate:      utils.ConvertTimeBangkok(category.LastUpdate),
			ValidFor:        validFor,
			Products:        category.Products,
			LifecycleStatus: category.LifecycleStatus,
		})
	}
	return categories, nil
}
