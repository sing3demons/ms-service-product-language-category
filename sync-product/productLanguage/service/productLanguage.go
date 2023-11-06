package service

import (
	"fmt"
	"time"

	"github.com/sing3demons/product.product.sync/producer"
	"github.com/sing3demons/product.product.sync/productLanguage/model"
	"github.com/sing3demons/product.product.sync/productLanguage/repository"
	"github.com/sing3demons/product.product.sync/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductLanguageService struct {
	repo    *repository.ProductLanguageRepository
	produce *producer.Producer
}

func NewProductLanguageService(repo *repository.ProductLanguageRepository, produce *producer.Producer) *ProductLanguageService {
	return &ProductLanguageService{repo: repo, produce: produce}
}

func (s *ProductLanguageService) CreateProductLanguage(req model.ProductLanguage) error {
	if req.ID == "" {
		return fmt.Errorf("id is empty")
	}

	document := model.ProductLanguage{
		ID:           req.ID,
		Name:         req.Name,
		Version:      req.Version,
		LastUpdate:   time.Now().Format(time.RFC3339),
		ValidFor:     req.ValidFor,
		LanguageCode: req.LanguageCode,
		Attachment:   req.Attachment,
	}

	if err := s.produce.SendMessage("productLanguage.createProductLanguage", "", document); err != nil {
		fmt.Printf("error: %v", err)
		return err
	}

	return nil
}

func (s *ProductLanguageService) FindProductLanguage(id string) (*model.ProductLanguage, error) {
	category, err := s.repo.FindProductLanguage(id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *ProductLanguageService) FindAllCategory(query model.Query) ([]model.ProductLanguage, error) {
	fmt.Print(query)
	filter := bson.D{}

	if query.Name != "" {
		filter = append(filter, bson.E{Key: "name", Value: bson.D{{Key: "$regex", Value: query.Name}}})
	}

	if query.ID != "" {
		filter = append(filter, bson.E{Key: "id", Value: query.ID})
	}

	if query.Limit != 0 {
		filter = append(filter, bson.E{Key: "limit", Value: query.Limit})
	}

	if query.LifecycleStatus != "" {
		filter = append(filter, bson.E{Key: "lifecycleStatus", Value: query.LifecycleStatus})
	}

	result, err := s.repo.FindAllProductLanguage(filter)
	products := []model.ProductLanguage{}
	if err != nil {
		return nil, err
	}

	fmt.Printf("result: %v", result)

	if len(result) != 0 {
		for _, v := range result {
			var product model.ProductLanguage
			product.ID = v.ID
			product.Href = utils.Href(v.Type, v.ID)
			if v.Name != "" {
				product.Name = v.Name
			}

			if v.Version != "" {
				product.Version = v.Version
			}

			if v.LastUpdate != "" {
				product.LastUpdate = utils.ConvertTimeBangkok(v.LastUpdate)
			}

			if v.LanguageCode != "" {
				product.LanguageCode = v.LanguageCode
			}

			if v.Attachment != nil {
				product.Attachment = v.Attachment
			}

			if v.ValidFor != nil {
				if product.ValidFor.EndDateTime != "" && product.ValidFor.StartDateTime != "" {
					product.ValidFor.EndDateTime = utils.ConvertTimeBangkok(product.ValidFor.EndDateTime)
					product.ValidFor.StartDateTime = utils.ConvertTimeBangkok(product.ValidFor.StartDateTime)
				}
			}

			products = append(products, product)
		}
	}

	return products, nil
}
