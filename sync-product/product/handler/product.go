package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/product/model"
	"github.com/sing3demons/product.product.sync/product/service"
)

type Product struct {
	service *service.ProductService
}

func NewProduct(service *service.ProductService) *Product {
	return &Product{service: service}
}

func (h *Product) CreateProduct(c *gin.Context) {
	var req model.Products
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateProduct(c, req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (h *Product) FindProduct(c *gin.Context) {
	start := time.Now()
	id := c.Param("id")
	product, err := h.service.FindProduct(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]any{
		"product": product,
		"duration": time.Since(start).String(),
	})
}

func (h *Product) FindAllProduct(c *gin.Context) {
	query := dto.Query{}
	name := c.Query("name")
	if name != "" {
		query.Name = name
	}

	id := c.Query("id")
	if id != "" {
		query.ID = id
	}

	limit := c.Query("limit")
	if limit != "" {
		size, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		query.Limit = size
	}

	query.LifecycleStatus = c.Query("lifecycleStatus")
	query.Expand = c.Query("expand")
	products, err := h.service.FindAllProducts(query)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}
