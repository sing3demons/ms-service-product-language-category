package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/product.product.sync/category/category/model"
	"github.com/sing3demons/product.product.sync/category/category/service"
	"github.com/sing3demons/product.product.sync/common/constants"
)

type Category struct {
	service *service.CategoryService
}

func NewCategory(service *service.CategoryService) *Category {
	return &Category{service: service}
}

func (h *Category) CreateCategory(c *gin.Context) {
	var req model.Category
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if req.LifecycleStatus != constants.Active && req.LifecycleStatus != constants.Inactive {
		c.JSON(400, gin.H{"error": "lifecycleStatus is invalid"})
		return
	}

	if err := h.service.CreateCategory(req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}

func (h *Category) FindCategory(c *gin.Context) {
	id := c.Param("id")
	category, err := h.service.FindCategory(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Category) FindAllCategory(c *gin.Context) {
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
	expand := c.Query("expand")
	if expand != "" {
		query.Expand = expand
	}

	lifecycleStatus := c.Query("lifecycleStatus")
	if lifecycleStatus != "" {
		query.LifecycleStatus = lifecycleStatus
	}
	category, err := h.service.FindAllCategory(query)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func (h *Category) UpdateCategory(c *gin.Context) {
	var req model.UpdateCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	id := c.Param("id")

	if req.LifecycleStatus != "" {
		if req.LifecycleStatus != constants.Active && req.LifecycleStatus != constants.Inactive {
			c.JSON(400, gin.H{"error": "lifecycleStatus is invalid"})
			return
		}
	}

	if err := h.service.UpdateCategory(id, req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "success"})
}
