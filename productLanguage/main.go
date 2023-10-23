package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sing3demons/productLanguage/db"
	"github.com/sing3demons/productLanguage/language/handler"
	"github.com/sing3demons/productLanguage/language/repository"
	"github.com/sing3demons/productLanguage/language/service"
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

	repo := repository.NewProductLanguage(col)
	service := service.NewProductLanguageService(repo)
	handler := handler.NewProductLanguage(service)

	r := gin.Default()
	r.GET("/productLanguage", handler.FindAllCategory)
	r.GET("/productLanguage/:id", handler.FindCategory)
	r.POST("/productLanguage", handler.CreateCategory)

	r.Run(":9090")

}
