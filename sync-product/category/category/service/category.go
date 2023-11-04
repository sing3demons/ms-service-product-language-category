package service

import (
	"fmt"
	"time"

	"github.com/sing3demons/product.product.sync/category/category/model"
	"github.com/sing3demons/product.product.sync/category/category/repository"
	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/producer"
	"github.com/sing3demons/product.product.sync/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(req model.Category) error {
	req.ID, _ = utils.RandomNanoID(11)
	if req.ID == "" {
		return fmt.Errorf("id is empty")
	}
	loc, _ := time.LoadLocation("Asia/Bangkok")
	document := model.Category{
		Type:            "Category",
		ID:              req.ID,
		Name:            req.Name,
		Version:         req.Version,
		LastUpdate:      time.Now().In(loc).Format("2006-01-02T15:04:05Z07:00"),
		ValidFor:        req.ValidFor,
		Products:        req.Products,
		LifecycleStatus: req.LifecycleStatus,
	}

	// servers := "localhost:9092"
	produce := producer.NewProducer()
	if err := produce.SendMessage("category.createCategory", "", document); err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) UpdateCategory(id string, req model.UpdateCategory) error {
	document := model.UpdateCategory{}

	if req.ID == id {
		document.ID = req.ID
	}
	if req.Name != "" {
		document.Name = req.Name
	}

	if req.LastUpdate != "" {
		document.LastUpdate = utils.ConvertTimeBangkok(req.LastUpdate)
	}

	if req.Version != "" {
		document.Version = req.Version
	}

	if req.LifecycleStatus != "" {
		document.LifecycleStatus = req.LifecycleStatus
	}

	if req.Products != nil {
		productsReq := req.Products
		var products []model.ProductRef
		for _, v := range productsReq {
			var product model.ProductRef
			product.ID = v.ID
			if v.Name != "" {
				product.Name = v.Name
			}
			if v.Version != "" {
				product.Version = v.Version
			}
			if v.LastUpdate != "" {
				product.LastUpdate = utils.ConvertTimeBangkok(v.LastUpdate)
			}
			product.Type = "Product"

			products = append(products, product)
		}
		document.Products = products
	}

	// servers := "localhost:9092"
	produce := producer.NewProducer()
	if err := produce.SendMessage("category.updateCategory", "", document); err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) DeleteCategory(id string, req model.UpdateCategory) error {
	document := model.UpdateCategory{}

	// servers := "localhost:9092"
	produce := producer.NewProducer()
	if err := produce.SendMessage("category.deleteCategory", "", document); err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) FindCategory(id string) (*model.Category, error) {
	category, err := s.repo.FindCategory(id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) FindAllCategory(query model.Query) ([]dto.Category, error) {
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

	category, err := s.repo.FindAllCategory(filter)
	if err != nil {
		return nil, err
	}

	return category, nil
}
