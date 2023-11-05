package service

import (
	"fmt"
	"time"

	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/producer"
	"github.com/sing3demons/product.product.sync/product/model"
	"github.com/sing3demons/product.product.sync/product/repository"
	"github.com/sing3demons/product.product.sync/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductService struct {
	repo    *repository.ProductRepository
	produce *producer.Producer
}

func NewProductService(repo *repository.ProductRepository, produce *producer.Producer) *ProductService {
	return &ProductService{repo: repo, produce: produce}
}

func (s *ProductService) CreateProduct(req model.Products) error {
	if req.ID == "" {
		return fmt.Errorf("id is empty")
	}

	document := model.Products{
		Name:               req.Name,
		Version:            req.Version,
		LastUpdate:         utils.ConvertTimeBangkok(time.Now().String()),
		ValidFor:           req.ValidFor,
		ID:                 req.ID,
		LifecycleStatus:    req.LifecycleStatus,
		Category:           req.Category,
		SupportingLanguage: req.SupportingLanguage,
	}

	if err := s.produce.SendMessage("product.createProduct", "", document); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) FindProduct(id string) (*model.Products, error) {
	product, err := s.repo.FindProduct(id)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) FindAllProducts(query dto.Query) (*model.ResponseDataWithTOtal, error) {
	fmt.Print(query)
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
	if query.Category == "product.Category" {
		filter = append(filter, bson.E{Key: "category", Value: query.Category})
	}
	if query.ProductLanguage == "product.productLanguage" {
		filter = append(filter, bson.E{Key: "productLanguage", Value: query.ProductLanguage})
	}

	products, err := s.repo.FindAllProduct(filter)
	if err != nil {
		return nil, err
	}

	total := make(chan int64)

	go s.repo.GetTotalProduct(filter, total)

	responses := model.ResponseDataWithTOtal{
		Total:    <-total,
		Products: products,
	}

	return &responses, nil
}
