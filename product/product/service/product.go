package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sing3demons/product/producer"
	"github.com/sing3demons/product/product/client"
	"github.com/sing3demons/product/product/model"
	"github.com/sing3demons/product/product/repository"
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

func (s *ProductService) FindAllProducts(query model.Query) ([]model.Products, error) {
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

	var queryExpand []string
	if query.Expand != "" {
		queryExpand = append(queryExpand, strings.Split(query.Expand, ",")...)
	}
	var (
		queryCategory           bool
		querySupportingLanguage bool
	)
	for _, expand := range queryExpand {
		if expand == "category" {
			queryCategory = true
		} else if expand == "supportingLanguage" {
			querySupportingLanguage = true
		}
	}

	products, err := s.repo.FindAllProduct(filter)
	if err != nil {
		return nil, err
	}

	responses := []model.Products{}

	for _, product := range products {
		categories := []model.Category{}
		if product.Category != nil {
			for _, category := range product.Category {
				if category.ID != "" {
					if queryCategory {
						b, err := client.RequestHttpGET("http://localhost:8080/category/" + category.ID)
						if err != nil {
							return nil, err
						}
						var cate model.Category
						json.Unmarshal(b, &cate)

						categories = append(categories, model.Category{
							Type:       cate.Type,
							ID:         cate.ID,
							Href:       cate.Href,
							Name:       cate.Name,
							Version:    cate.Version,
							LastUpdate: cate.LastUpdate,
							ValidFor:   cate.ValidFor,
						})
					} else {
						categories = append(categories, model.Category{
							Type: category.Type,
							ID:   category.ID,
							Name: category.Name,
						})
					}

				} else {
					categories = append(categories, model.Category{
						Type: category.Type,
						ID:   category.ID,
						Name: category.Name,
					})
				}

				product.Category = categories
			}
		}
		productLanguages := []model.ProductLanguage{}
		for _, language := range product.SupportingLanguage {
			if language.ID != "" {
				if querySupportingLanguage {
					b, err := client.RequestHttpGET("http://localhost:9090/productLanguage/" + language.ID)
					if err != nil {
						return nil, err
					}
					var productLanguage model.ProductLanguage
					json.Unmarshal(b, &productLanguage)

					productLanguages = append(productLanguages, model.ProductLanguage{
						Type:         productLanguage.Type,
						ID:           productLanguage.ID,
						Href:         productLanguage.Href,
						LanguageCode: productLanguage.LanguageCode,
						Name:         productLanguage.Name,
						Version:      productLanguage.Version,
						LastUpdate:   productLanguage.LastUpdate,
						ValidFor:     productLanguage.ValidFor,
						Attachment:   productLanguage.Attachment,
					})
				} else {
					productLanguages = append(productLanguages, model.ProductLanguage{
						Type:         language.Type,
						ID:           language.ID,
						Href:         language.Href,
						LanguageCode: language.LanguageCode,
					})
				}

			}
		}
		product.SupportingLanguage = productLanguages

		responses = append(responses, product)

	}

	return responses, nil
}
