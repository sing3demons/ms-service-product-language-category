package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductPriceLanguageRepository struct {
	db *mongo.Collection
}

func NewProductPriceLanguageRepository(db *mongo.Collection) *ProductPriceLanguageRepository {
	return &ProductPriceLanguageRepository{db}
}

func (r *ProductPriceLanguageRepository) FindAll(doc bson.D) ([]dto.ProductPriceLanguage, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	filter := bson.D{}
	filter = append(filter, bson.E{Key: "deleteDate", Value: nil})
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
	// cur, err := r.db.Find(ctx, filter)
	products, err := utils.GetMultiple[dto.ProductPriceLanguage](r.db, filter)
	if err != nil {
		return nil, err
	}
	// defer cur.Close(ctx)
	// products := []dto.ProductPriceLanguage{}
	// for cur.Next(ctx) {
	// 	var product dto.ProductPriceLanguage
	// 	err := cur.Decode(&product)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	products = append(products, product)
	// }
	return products, nil
}

func (r *ProductPriceLanguageRepository) FindOne(id string) (*dto.ProductPriceLanguage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var product dto.ProductPriceLanguage
	err := r.db.FindOne(ctx, bson.M{"id": id, "deleteDate": nil}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}
