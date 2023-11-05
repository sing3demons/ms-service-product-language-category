package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/product/product/model"
	"github.com/sing3demons/product.product.sync/productPriceLanguage/service"
)

type ProductPriceLanguageHandler struct {
	svc *service.ProductPriceLanguageService
}

func NewProductPriceLanguageHandler(svc *service.ProductPriceLanguageService) *ProductPriceLanguageHandler {
	return &ProductPriceLanguageHandler{svc}
}

func (h *ProductPriceLanguageHandler) FindProductPriceLanguages(c *gin.Context) {
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
	query.Expand = c.Query("expand")
	products, err := h.svc.FindProductPriceLanguages(query)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, products)
}

func (h *ProductPriceLanguageHandler) FindProductPriceLanguage(c *gin.Context) {
	id := c.Param("id")
	product, err := h.svc.FindProductPriceLanguage(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, product)
}

func (h *ProductPriceLanguageHandler) CreateProductPrice(c *gin.Context) {
	var req dto.ProductPriceLanguage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.svc.CreateProductPriceLanguage(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "success"})
}
