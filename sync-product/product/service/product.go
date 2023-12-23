package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/product.product.sync/common"
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

func (s *ProductService) CreateProduct(c *gin.Context, req model.Products) error {
	if req.ID == "" {
		return fmt.Errorf("id is empty")
	}

	var document model.Products
	document.ID = req.ID
	document.Name = req.Name
	document.Version = req.Version
	document.LastUpdate = req.LastUpdate
	document.ValidFor = req.ValidFor
	document.LifecycleStatus = req.LifecycleStatus
	document.Category = req.Category
	document.SupportingLanguage = req.SupportingLanguage
	document.Type = "Products"

	if err := s.produce.SendMessage(c, "product.createProduct", "", document); err != nil {
		return err
	}

	return nil
}

func (s *ProductService) FindProduct(id string) (*model.Products, error) {
	product, err := s.repo.FindProduct(id)

	if err != nil {
		return nil, err
	}

	productLanguage := []model.ProductLanguage{}
	var wg sync.WaitGroup
	for _, v := range product.SupportingLanguage {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			url := utils.GetHost() + "/productLanguage/" + id
			result, err := common.RequestHttpGet[model.ProductLanguage](url)
			if err != nil {
				fmt.Printf("error call product for productID: %s, error : %v\n", id, err)
				return
			}

			productLanguage = append(productLanguage, *result)

		}(v.ID)
	}
	wg.Wait()

	if len(productLanguage) != 0 {
		product.SupportingLanguage = productLanguage
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

	result, err := s.repo.FindAllProduct(filter)
	if err != nil {
		return nil, err
	}

	products := []model.Products{}

	var categoryProducts bool
	if query.Expand != "" {
		expand := strings.Split(query.Expand, ",")
		for i := 0; i < len(expand); i++ {
			if expand[i] == "product.category" {
				categoryProducts = true
			}
		}
	}

	for _, item := range result {
		categories := []model.Category{}
		if item.Category != nil {
			if categoryProducts {
				result := s.GetCategoryFromProducts(item.Category)
				for _, v := range result {
					p := []model.Products{}
					if v.Products != nil {
						for _, vv := range v.Products {
							p = append(p, model.Products{
								ID:              vv.ID,
								Name:            vv.Name,
								Type:            vv.Type,
								Href:            utils.Href(vv.Type, vv.ID),
								LifecycleStatus: vv.LifecycleStatus,
								LastUpdate:      vv.LastUpdate,
								ValidFor:        vv.ValidFor,
								Version:         vv.Version,
							})
						}

					}
					categories = append(categories, model.Category{
						Type:            v.Type,
						ID:              v.ID,
						Href:            utils.Href(v.Type, v.ID),
						Name:            v.Name,
						Version:         v.Version,
						LastUpdate:      v.LastUpdate,
						ValidFor:        v.ValidFor,
						LifecycleStatus: v.LifecycleStatus,
						Products:        p,
					})
				}

			} else {
				for _, v := range item.Category {
					if v.LastUpdate != "" {
						v.LastUpdate = utils.ConvertTimeBangkok(v.LastUpdate)
					}
					categories = append(categories, model.Category{
						Type:            v.Type,
						ID:              v.ID,
						Href:            utils.Href(v.Type, v.ID),
						Name:            v.Name,
						Version:         v.Version,
						LastUpdate:      v.LastUpdate,
						ValidFor:        v.ValidFor,
						LifecycleStatus: v.LifecycleStatus,
					})
				}
			}
		}

		product := model.Products{
			ID:                 item.ID,
			Name:               item.Name,
			Version:            item.Version,
			Href:               utils.Href(item.Type, item.ID),
			LastUpdate:         item.LastUpdate,
			ValidFor:           item.ValidFor,
			LifecycleStatus:    item.LifecycleStatus,
			Category:           categories,
			SupportingLanguage: item.SupportingLanguage,
		}

		products = append(products, product)
	}

	total := make(chan int64)

	go s.repo.GetTotalProduct(filter, total)

	responses := model.ResponseDataWithTOtal{
		Total:    <-total,
		Products: products,
	}

	return &responses, nil
}

func (s *ProductService) GetCategoryFromProducts(categories []model.Category) []model.Category {
	start := time.Now()
	var wg sync.WaitGroup
	resultCh := make(chan model.Category, len(categories))
	poolSize := 10
	semaphore := make(chan struct{}, poolSize)
	for _, v := range categories {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			b, err := common.HttpGET(utils.GetHost() + "/category/" + id)
			if err != nil {
				fmt.Println("error call product for productID:", id)
				fmt.Println(err)
				return
			}

			var result model.Category
			if err := json.Unmarshal(b, &result); err != nil {
				fmt.Println("error Unmarshal for productID:", id)
				fmt.Println(err)
				return
			}
			products := []model.Products{}
			for _, v := range result.Products {
				productByte, err := common.HttpGET(utils.GetHost() + "/products/" + v.ID)
				if err != nil {
					fmt.Println("error call product for productID:", id)
					fmt.Println(err)
					return
				}

				var product model.Products
				if err := json.Unmarshal(productByte, &product); err != nil {
					fmt.Println("error Unmarshal for productID:", id)
					fmt.Println(err)
					return
				}

				products = append(products, model.Products{
					Type:       product.Type,
					ID:         product.ID,
					Name:       product.Name,
					Version:    product.Version,
					Href:       utils.Href(product.Type, product.ID),
					ValidFor:   product.ValidFor,
					LastUpdate: product.LastUpdate,
				})
			}

			result.Products = products

			// var validFor *dto.ValidFor
			// if result.ValidFor != nil {
			// 	validFor = &model.ValidFor{
			// 		StartDateTime: utils.ConvertTimeBangkok(result.ValidFor.StartDateTime),
			// 		EndDateTime:   utils.ConvertTimeBangkok(result.ValidFor.EndDateTime),
			// 	}
			// 	result.ValidFor = validFor
			// }

			if result.LastUpdate != "" {
				result.LastUpdate = utils.ConvertTimeBangkok(result.LastUpdate)
			}

			resultCh <- result
		}(v.ID)
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	var categoriesDto []model.Category
	for r := range resultCh {
		categoriesDto = append(categoriesDto, r)
	}
	end := time.Now()
	fmt.Printf("Time: %s\n", end.Sub(start))

	return categoriesDto
}
