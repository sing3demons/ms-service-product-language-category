package common

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sing3demons/product.product.sync/category/model"
	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/product/db"
	"github.com/sing3demons/product.product.sync/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func GetProductFromCategory(category model.Category) []dto.Products {
	start := time.Now()
	products := []dto.Products{}
	for _, v := range category.Products {
		b, err := HttpGET(utils.GetHost() + "/products/" + v.ID)
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
		} else {
			product.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
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
				if v.LastUpdate != "" {
					v.LastUpdate = utils.ConvertTimeBangkok(v.LastUpdate)
				}
				v.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
			}
		}

		if product.Category != nil {
			for _, v := range product.Category {
				if v.LastUpdate == "" {
					v.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
				}
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
	end := time.Now()
	fmt.Printf("Time: %s\n", end.Sub(start))
	return products
}

func GetProduct(category model.Category) ([]dto.Products, error) {
	col, err := db.ConnectMonoDB()
	if err != nil {
		return nil, err
	}

	products := []dto.Products{}
	for _, v := range category.Products {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		product := dto.Products{}
		if err := col.FindOne(ctx, bson.M{"id": v.ID}).Decode(&product); err != nil {
			return nil, err
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
		} else {
			product.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
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
				if v.LastUpdate != "" {
					v.LastUpdate = utils.ConvertTimeBangkok(v.LastUpdate)
				}
				v.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
			}
		}

		if product.Category != nil {
			for _, v := range product.Category {
				if v.LastUpdate == "" {
					v.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
				}
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
	return products, nil
}
