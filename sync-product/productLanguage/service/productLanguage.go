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
	return &ProductLanguageService{repo: repo}
}

func (s *ProductLanguageService) CreateProductLanguage(req model.ProductLanguage) error {
	if req.ID == "" {
		return fmt.Errorf("id is empty")
	}
	document := model.ProductLanguage{
		ID:           req.ID,
		Name:         req.Name,
		Version:      req.Version,
		LastUpdate:   utils.ConvertTimeBangkok(time.Now().String()),
		ValidFor:     req.ValidFor,
		LanguageCode: req.LanguageCode,
		Attachment:   req.Attachment,
	}

	if err := s.produce.SendMessage("productLanguage.createProductLanguage", "", document); err != nil {
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
	filter := bson.M{}

	if query.Name != "" {
		filter["name"] = query.Name
	}

	if query.ID != "" {
		filter["id"] = query.ID
	}

	if query.Limit != 0 {
		filter["limit"] = query.Limit
	}

	if query.LifecycleStatus != "" {
		filter["lifecycleStatus"] = query.LifecycleStatus
	}

	category, err := s.repo.FindAllProductLanguage(filter)
	if err != nil {
		return nil, err
	}

	return category, nil
}
