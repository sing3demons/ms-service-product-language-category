package service

import (
	"fmt"
	"time"

	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/producer"
	"github.com/sing3demons/product.product.sync/product/product/model"
	"github.com/sing3demons/product.product.sync/productPrice/repository"
	"github.com/sing3demons/product.product.sync/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductPriceService struct {
	repo *repository.ProductPriceRepository
}

func NewProductPriceService(repo *repository.ProductPriceRepository) *ProductPriceService {
	return &ProductPriceService{repo}
}

func (svc *ProductPriceService) FindProductPrices(query model.Query) ([]dto.ProductPrice, error) {
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

	result, err := svc.repo.FindAll(filter)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (svc *ProductPriceService) FindProductPrice(id string) (*dto.ProductPrice, error) {
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	return svc.repo.FindOne(id)
}

func (svc *ProductPriceService) CreateProductPrice(req dto.ProductPrice) error {

	doc := dto.ProductPrice{
		Id:                 req.Id,
		Type:               req.Type,
		LifecycleStatus:    req.LifecycleStatus,
		Name:               req.Name,
		LastUpdate:         utils.ConvertTimeBangkok(time.Now().String()),
		Version:            req.Version,
		Price:              req.Price,
		PopRelationship:    req.PopRelationship,
		SupportingLanguage: req.SupportingLanguage,
	}

	if req.ValidFor.StartDateTime != "" && req.ValidFor.EndDateTime != "" {
		doc.ValidFor.StartDateTime = utils.ConvertTimeBangkok(req.ValidFor.StartDateTime)
		doc.ValidFor.EndDateTime = utils.ConvertTimeBangkok(req.ValidFor.EndDateTime)
	}

	produce := producer.NewProducer()
	if err := produce.SendMessage("product.createProductPrice", "", doc); err != nil {
		return err
	}
	return nil
}
