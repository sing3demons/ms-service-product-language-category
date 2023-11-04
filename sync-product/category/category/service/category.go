package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sing3demons/product.product.sync/category/category/model"
	"github.com/sing3demons/product.product.sync/category/category/repository"
	"github.com/sing3demons/product.product.sync/common"
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

	categories, err := s.repo.FindAllCategory(filter)
	if err != nil {
		return nil, err
	}

	var categoryProducts bool
	if query.Expand != "" {
		expand := strings.Split(query.Expand, ",")
		for i := 0; i < len(expand); i++ {
			if expand[i] == "category.products" {
				categoryProducts = true
			}
		}
	}

	result := []dto.Category{}

	if categories != nil {
		var products []dto.Products
		for i := 0; i < len(categories); i++ {
			category := categories[i]
			validFor := &dto.ValidFor{
				StartDateTime: utils.ConvertTimeBangkok(category.ValidFor.StartDateTime),
				EndDateTime:   utils.ConvertTimeBangkok(category.ValidFor.EndDateTime),
			}
			if category.Products != nil {
				if categoryProducts {
					for _, v := range category.Products {
						b, err := common.HttpGET(utils.GetHost() + "/products/" + v.ID)
						if err != nil {
							fmt.Println("error call product")
							fmt.Println(err)
						}
						var product dto.Products
						if err := json.Unmarshal(b, &product); err != nil {
							fmt.Println("error Unmarshal")
							fmt.Println(err)
						}
						var validFor *dto.ValidFor
						if product.ValidFor != nil {
							validFor = &dto.ValidFor{
								StartDateTime: utils.ConvertTimeBangkok(product.ValidFor.StartDateTime),
								EndDateTime:   utils.ConvertTimeBangkok(product.ValidFor.EndDateTime),
							}
						}

						if product.LastUpdate != "" {
							product.LastUpdate = utils.ConvertTimeBangkok(product.LastUpdate)
						} else {
							product.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
						}

						productRef := dto.Products{
							Type:            product.Type,
							ID:              product.ID,
							Href:            utils.Href(product.Type, product.ID),
							Name:            product.Name,
							Version:         product.Version,
							LastUpdate:      product.LastUpdate,
							ValidFor:        validFor,
							LifecycleStatus: product.LifecycleStatus,
						}

						if product.Category != nil {
							for _, v := range product.Category {
								if v.LastUpdate != "" {
									v.LastUpdate = utils.ConvertTimeBangkok(v.LastUpdate)
								}
								v.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
								fmt.Printf("=== %v-->\n", v.LastUpdate)

							}
						}

						if product.Category != nil {
							for _, v := range product.Category {
								if v.LastUpdate == "" {
									v.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
								}
								productRef.Category = append(productRef.Category, dto.Category{
									Type:            v.Type,
									ID:              v.ID,
									Href:            utils.Href(v.Type, v.ID),
									Name:            v.Name,
									Version:         v.Version,
									LastUpdate:      utils.ConvertTimeBangkok(v.LastUpdate),
									ValidFor:        validFor,
									LifecycleStatus: v.LifecycleStatus,
								})
							}
						}

						products = append(products, productRef)

					}
				} else {
					for _, v := range category.Products {
						b, err := common.HttpGET(utils.GetHost() + "/products/" + v.ID)
						if err != nil {
							fmt.Println("error call product")
							fmt.Println(err)
						}
						var product = new(dto.Products)
						if err := json.Unmarshal(b, &product); err != nil {
							fmt.Println("error Unmarshal")
							fmt.Println(err)
						}
						if product != nil {
							products = append(products, dto.Products{
								Type:    "Product",
								ID:      v.ID,
								Href:    utils.Href("Product", v.ID),
								Name:    v.Name,
								Version: v.Version,
							})
						} else {
							products = append(products, dto.Products{
								ID:   v.ID,
								Name: v.Name,
							})
						}
					}
				}
			}
			result = append(result, dto.Category{
				Type:            category.Type,
				ID:              category.ID,
				Href:            utils.Href(category.Type, category.ID),
				Name:            category.Name,
				Version:         category.Version,
				LastUpdate:      utils.ConvertTimeBangkok(category.LastUpdate),
				ValidFor:        validFor,
				Products:        products,
				LifecycleStatus: category.LifecycleStatus,
			})
		}
	}

	return result, nil
}
