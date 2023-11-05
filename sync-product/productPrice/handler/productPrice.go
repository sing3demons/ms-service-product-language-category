package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/product.product.sync/common/constants"
	"github.com/sing3demons/product.product.sync/common/dto"
	"github.com/sing3demons/product.product.sync/product/product/model"
	"github.com/sing3demons/product.product.sync/productPrice/service"
)

type ProductPriceHandler struct {
	svc *service.ProductPriceService
}

func NewProductPriceHandler(svc *service.ProductPriceService) *ProductPriceHandler {
	return &ProductPriceHandler{svc}
}

func (h *ProductPriceHandler) FindProductPrices(c *gin.Context) {
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

	result, err := h.svc.FindProductPrices(query)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, result)
}

func (h *ProductPriceHandler) FindProductPrice(c *gin.Context) {
	id := c.Param("id")
	product, err := h.svc.FindProductPrice(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, product)
}

func (h *ProductPriceHandler) CreateProductPrice(c *gin.Context) {
	var req dto.ProductPrice
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if req.LifecycleStatus != constants.Active && req.LifecycleStatus != constants.Inactive {
		c.JSON(400, gin.H{"error": "lifecycleStatus is invalid"})
		return
	}

	if err := h.svc.CreateProductPrice(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"message": "success"})
}
