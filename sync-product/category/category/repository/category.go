package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sing3demons/product.product.sync/category/category/model"
	"github.com/sing3demons/product.product.sync/common"
	"github.com/sing3demons/product.product.sync/common/dto"
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
func (r *CategoryRepository) FindAllCategory(doc bson.D) ([]dto.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{}
	var categoryProducts bool
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

		if v.Key == "expand" {
			expand := strings.Split(fmt.Sprintf("%s", v.Value), ",")
			for i := 0; i < len(expand); i++ {
				if expand[i] == "category.products" {
					categoryProducts = true
				}
			}
		}
	}

	cursor, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	categories := []dto.Category{}

	for cursor.Next(ctx) {
		var category dto.CategoryProducts
		if err := cursor.Decode(&category); err != nil {
			return nil, err
		}

		validFor := &dto.ValidFor{
			StartDateTime: utils.ConvertTimeBangkok(category.ValidFor.StartDateTime),
			EndDateTime:   utils.ConvertTimeBangkok(category.ValidFor.EndDateTime),
		}

		var products []dto.Products
		if category.Products != nil {
			if categoryProducts {
				for _, v := range category.Products {
					b, err := common.HttpGET(utils.GetHost() + "/products/" + v.ID)
					if err != nil {
						fmt.Println("error call product")
						fmt.Println(err)
					}
					var product dto.Products
					if err := json.Unmarshal(b, &product); err != nil {
						fmt.Println("error Unmarshal")
						fmt.Println(err)
					}
					var validFor *dto.ValidFor
					if product.ValidFor != nil {
						validFor = &dto.ValidFor{
							StartDateTime: utils.ConvertTimeBangkok(product.ValidFor.StartDateTime),
							EndDateTime:   utils.ConvertTimeBangkok(product.ValidFor.EndDateTime),
						}
					}

					if product.LastUpdate != "" {
						product.LastUpdate = utils.ConvertTimeBangkok(product.LastUpdate)
					}
					productRef := dto.Products{
						Type:            product.Type,
						ID:              product.ID,
						Href:            utils.Href(product.Type, product.ID),
						Name:            product.Name,
						Version:         product.Version,
						LastUpdate:      product.LastUpdate,
						ValidFor:        validFor,
						LifecycleStatus: product.LifecycleStatus,
					}

					if product.Category != nil {
						for _, v := range product.Category {
							productRef.Category = append(productRef.Category, dto.Category{
								Type:            v.Type,
								ID:              v.ID,
								Href:            utils.Href(v.Type, v.ID),
								Name:            v.Name,
								Version:         v.Version,
								LastUpdate:      utils.ConvertTimeBangkok(v.LastUpdate),
								ValidFor:        validFor,
								LifecycleStatus: v.LifecycleStatus,
							})
						}
					}

					products = append(products, productRef)

				}
			} else {
				for _, v := range category.Products {
					b, err := common.HttpGET(utils.GetHost() + "/products/" + v.ID)
					if err != nil {
						fmt.Println("error call product")
						fmt.Println(err)
					}
					var product = new(dto.Products)
					if err := json.Unmarshal(b, &product); err != nil {
						fmt.Println("error Unmarshal")
						fmt.Println(err)
					}
					if product != nil {
						products = append(products, dto.Products{
							Type:    "Product",
							ID:      v.ID,
							Href:    utils.Href("Product", v.ID),
							Name:    v.Name,
							Version: v.Version,
						})
					} else {
						products = append(products, dto.Products{
							ID:   v.ID,
							Name: v.Name,
						})
					}

				}
			}

		}

		categories = append(categories, dto.Category{
			Type:            category.Type,
			ID:              category.ID,
			Href:            utils.Href(category.Type, category.ID),
			Name:            category.Name,
			Version:         category.Version,
			LastUpdate:      utils.ConvertTimeBangkok(category.LastUpdate),
			ValidFor:        validFor,
			Products:        products,
			LifecycleStatus: category.LifecycleStatus,
		})
	}
	return categories, nil
}
