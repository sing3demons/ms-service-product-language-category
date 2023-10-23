package service

import (
	"fmt"
	"time"

	"github.com/sing3demons/category/category/model"
	"github.com/sing3demons/category/category/repository"
	"github.com/sing3demons/category/customNanoid"
	"github.com/sing3demons/category/producer"
	"go.mongodb.org/mongo-driver/bson"
)

type CategoryService struct {
	repo *repository.CategoryRepository
}

func NewCategoryService(repo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(req model.Category) error {
	req.ID, _ = customNanoid.RandomNanoID(11)
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
		Href:       fmt.Sprintf("/category/%s", req.ID),
	}
	// servers := "localhost:9092"
	produce := producer.NewProducer("localhost:9092")
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

	category, err := s.repo.FindAllCategory(filter)
	if err != nil {
		return nil, err
	}

	return category, nil
}
