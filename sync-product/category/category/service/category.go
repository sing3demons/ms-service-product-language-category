package service

import (
	"fmt"
	"time"

	"github.com/sing3demons/product.product.sync/category/category/model"
	"github.com/sing3demons/product.product.sync/category/category/repository"
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
		Type:       "Category",
		ID:         req.ID,
		Name:       req.Name,
		Version:    req.Version,
		LastUpdate: time.Now().In(loc).Format("2006-01-02T15:04:05Z07:00"),
		ValidFor:   req.ValidFor,
		Products:   req.Products,
	}
	// servers := "localhost:9092"
	produce := producer.NewProducer()
	if err := produce.SendMessage("category.createCategory", "", document); err != nil {
		return err
	}

	// id, err := s.repo.CreateCategory(document)
	// if err != nil {
	// 	return model.Category{}, err
	// }
	// category, err := s.repo.FindCategoryById(id)
	// if err != nil {
	// 	return model.Category{}, err
	// }

	return nil
}

func (s *CategoryService) FindCategory(id string) (*model.Category, error) {
	category, err := s.repo.FindCategory(id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) FindAllCategory(query model.Query) ([]model.Category, error) {
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
