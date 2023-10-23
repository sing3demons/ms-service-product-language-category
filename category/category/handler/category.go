package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sing3demons/category/category/model"
	"github.com/sing3demons/category/category/service"
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

	// category, err := h.service.CreateCategory(req)
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

	query.LifecycleStatus = c.Query("lifecycleStatus")
	category, err := h.service.FindAllCategory(query)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}
