package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/product.product.sync/category/constants"
	"github.com/sing3demons/product.product.sync/category/model"
	"github.com/sing3demons/product.product.sync/category/repository"
	"github.com/sing3demons/product.product.sync/common"
	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/producer"
	"github.com/sing3demons/product.product.sync/utils"

	"go.mongodb.org/mongo-driver/bson"
)

type CategoryService struct {
	repo    *repository.CategoryRepository
	produce *producer.Producer
}

func NewCategoryService(repo *repository.CategoryRepository, produce *producer.Producer) *CategoryService {
	return &CategoryService{repo: repo, produce: produce}
}

func (s *CategoryService) CreateCategory(c *gin.Context, req model.Category) error {
	req.ID, _ = utils.RandomNanoID(11)
	if req.ID == "" {
		return fmt.Errorf("id is empty")
	}

	document := model.Category{
		Type:            "Category",
		ID:              req.ID,
		Name:            req.Name,
		Version:         req.Version,
		LastUpdate:      utils.ConvertTimeBangkok(time.Now().String()),
		ValidFor:        req.ValidFor,
		Products:        req.Products,
		LifecycleStatus: req.LifecycleStatus,
	}

	if err := s.produce.SendMessage(c, constants.CREATE_CATEGORY, "", document); err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) UpdateCategory(c *gin.Context, id string, req model.UpdateCategory) error {
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

	if err := s.produce.SendMessage(c, constants.UPDATE_CATEGORY, "", document); err != nil {
		return err
	}

	return nil
}

func (s *CategoryService) DeleteCategory(c *gin.Context, id string, req model.UpdateCategory) error {
	document := model.UpdateCategory{}

	if err := s.produce.SendMessage(c, constants.DELETE_CATEGORY, "", document); err != nil {
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
					products = s.GetProductFromCategory(category)
					// products, err = common.GetProduct(category)
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
								Href:    utils.Href(v.Type, v.ID),
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

func (s *CategoryService) GetProductFromCategory(category model.Category) []dto.Products {
	start := time.Now()
	var wg sync.WaitGroup
	productCh := make(chan dto.Products, len(category.Products))
	poolSize := 10
	semaphore := make(chan struct{}, poolSize)
	for _, v := range category.Products {
		wg.Add(1)
		go func(productID string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			b, err := common.HttpGET(utils.GetHost() + "/products/" + productID)
			if err != nil {
				fmt.Println("error call product for productID:", productID)
				fmt.Println(err)
				return
			}

			var product dto.Products
			if err := json.Unmarshal(b, &product); err != nil {
				fmt.Println("error Unmarshal for productID:", productID)
				fmt.Println(err)
				return
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
					} else {
						v.LastUpdate = utils.ConvertTimeBangkok(category.LastUpdate)
					}
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

			productCh <- productRef
		}(v.ID)
	}

	go func() {
		wg.Wait()
		close(productCh)
	}()

	var products []dto.Products
	for productRef := range productCh {
		products = append(products, productRef)
	}
	end := time.Now()
	fmt.Printf("Time: %s\n", end.Sub(start))

	return products
}
