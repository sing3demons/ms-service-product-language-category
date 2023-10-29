package service

import (
	"fmt"
	"time"

	"github.com/sing3demons/product.product.sync/product/producer"
	"github.com/sing3demons/product.product.sync/product/product/model"
	"github.com/sing3demons/product.product.sync/product/product/repository"
	"go.mongodb.org/mongo-driver/bson"
)

type ProductService struct {
	repo *repository.ProductRepository
}

func NewProductService(repo *repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) CreateProduct(req model.Products) error {
	if req.ID == "" {
		return fmt.Errorf("id is empty")
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	document := model.Products{
		Name:               req.Name,
		Version:            req.Version,
		LastUpdate:         time.Now().In(loc).Format("2006-01-02T15:04:05Z07:00"),
		ValidFor:           req.ValidFor,
		ID:                 req.ID,
		LifecycleStatus:    req.LifecycleStatus,
		Category:           req.Category,
		SupportingLanguage: req.SupportingLanguage,
	}
	servers := "localhost:9092"
	produce := producer.NewProducer(servers)
	if err := produce.SendMessage("product.createProduct", "", document); err != nil {
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

func (s *ProductService) FindAllProducts(query model.Query) (*model.ResponseDataWithTOtal, error) {
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
