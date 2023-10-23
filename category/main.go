package main

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sing3demons/category/category/handler"
	"github.com/sing3demons/category/category/repository"
	"github.com/sing3demons/category/category/service"
	"github.com/sing3demons/category/constants"
	"github.com/sing3demons/category/consume"
	"github.com/sing3demons/category/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	col, err := db.ConnectMonoDB()
	if err != nil {
		panic(err)
	}

	col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"id": 1,
		},
		Options: nil,
	})

	repo := repository.NewCategory(col)
	service := service.NewCategoryService(repo)
	handler := handler.NewCategory(service)

	r := gin.Default()
	r.GET("/category", handler.FindAllCategory)
	r.GET("/category/:id", handler.FindCategory)
	r.POST("/category", handler.CreateCategory)

	topics := []string{
		constants.CREATE_CATEGORY,
		constants.CREATE_CATEGORY_FAILED,
		constants.CREATE_CATEGORY_SUCCESS,
		constants.UPDATE_CATEGORY,
		constants.UPDATE_CATEGORY_FAILED,
		constants.UPDATE_CATEGORY_SUCCESS,
	}
	servers := "localhost:9092"
	groupID := "go.category.product"
	timeOut := time.Duration(-1)
	consume.Consumer(servers, groupID, topics, timeOut)

	r.Run(":8080")

}
