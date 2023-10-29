package service

import (
	"fmt"
	"time"

	"github.com/sing3demons/product.product.sync/productLanguage/language/model"
	"github.com/sing3demons/product.product.sync/productLanguage/language/repository"
	"github.com/sing3demons/product.product.sync/productLanguage/producer"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductLanguageService struct {
	repo *repository.ProductLanguageRepository
}

func NewProductLanguageService(repo *repository.ProductLanguageRepository) *ProductLanguageService {
	return &ProductLanguageService{repo: repo}
}

func (s *ProductLanguageService) CreateProductLanguage(req model.ProductLanguage) error {
	if req.ID == "" {
		return fmt.Errorf("id is empty")
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	document := model.ProductLanguage{
		ID:           req.ID,
		Href:         "/productLanguage/" + req.Href,
		Name:         req.Name,
		Version:      req.Version,
		LastUpdate:   time.Now().In(loc).Format("2006-01-02T15:04:05Z07:00"),
		ValidFor:     req.ValidFor,
		LanguageCode: req.LanguageCode,
		Attachment:   req.Attachment,
	}
	// servers := "localhost:9092"
	produce := producer.NewProducer("localhost:9092")
	if err := produce.SendMessage("productLanguage.createProductLanguage", "", document); err != nil {
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
