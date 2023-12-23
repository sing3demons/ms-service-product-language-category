package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/product.product.sync/productLanguage/model"
	"github.com/sing3demons/product.product.sync/productLanguage/service"
	"github.com/sing3demons/product.product.sync/utils"
)

type ProductLanguage struct {
	service *service.ProductLanguageService
}

func NewProductLanguage(service *service.ProductLanguageService) *ProductLanguage {
	return &ProductLanguage{service: service}
}

func (h *ProductLanguage) CreateProductLanguage(c *gin.Context) {
	var req model.ProductLanguage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateProductLanguage(c, req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (h *ProductLanguage) FindProductLanguage(c *gin.Context) {
	id := c.Param("id")
	product, err := h.service.FindProductLanguage(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	product.Href = utils.GetHost() + "/" + product.Type + "/" + product.ID

	c.JSON(http.StatusOK, product)
}

func (h *ProductLanguage) FindAllProductLanguage(c *gin.Context) {
	query := model.Query{}
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
	category, err := h.service.FindAllCategory(query)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}
