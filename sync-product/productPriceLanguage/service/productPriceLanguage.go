package service

import (
	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/producer"
	"github.com/sing3demons/product.product.sync/productPriceLanguage/repository"
	"github.com/sing3demons/product.product.sync/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductPriceLanguageService struct {
	repo *repository.ProductPriceLanguageRepository
}

func NewProductPriceLanguageService(repo *repository.ProductPriceLanguageRepository) *ProductPriceLanguageService {
	return &ProductPriceLanguageService{repo}
}

func (s *ProductPriceLanguageService) FindProductPriceLanguages(query dto.Query) ([]dto.ProductPriceLanguage, error) {
	filter := bson.D{}

	if query.Name != "" {
		filter = append(filter, bson.E{Key: "name", Value: query.Name})
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

	if query.Expand != "" {
		filter = append(filter, bson.E{Key: "expand", Value: query.Expand})
	}
	return s.repo.FindAll(filter)
}

func (s *ProductPriceLanguageService) FindProductPriceLanguage(id string) (*dto.ProductPriceLanguage, error) {
	return s.repo.FindOne(id)
}

func (s *ProductPriceLanguageService) CreateProductPriceLanguage(req dto.ProductPriceLanguage) error {
	doc := dto.ProductPriceLanguage{
		Type:         "productPriceLanguage",
		Id:           req.Id,
		LanguageCode: req.LanguageCode,
		Version:      req.Version,
		Price:        req.Price,
		LastUpdate:   utils.ConvertTimeBangkok(req.LastUpdate),
	}

	if req.ValidFor.StartDateTime != "" && req.ValidFor.EndDateTime != "" {
		doc.ValidFor.StartDateTime = utils.ConvertTimeBangkok(req.ValidFor.StartDateTime)
		doc.ValidFor.EndDateTime = utils.ConvertTimeBangkok(req.ValidFor.EndDateTime)
	}

	produce := producer.NewProducer()
	if err := produce.SendMessage("product.createProductPriceLanguage", "", doc); err != nil {
		return err
	}
	return nil
}
