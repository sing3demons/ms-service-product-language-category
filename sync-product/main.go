package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	category_db "github.com/sing3demons/product.product.sync/category/db"
	handler_category "github.com/sing3demons/product.product.sync/category/handler"
	repository_category "github.com/sing3demons/product.product.sync/category/repository"
	service_category "github.com/sing3demons/product.product.sync/category/service"
	product_db "github.com/sing3demons/product.product.sync/product/db"
	handler_product "github.com/sing3demons/product.product.sync/product/handler"
	repository_product "github.com/sing3demons/product.product.sync/product/repository"
	service_product "github.com/sing3demons/product.product.sync/product/service"
	productLanguage_db "github.com/sing3demons/product.product.sync/productLanguage/db"
	handler_productLanguage "github.com/sing3demons/product.product.sync/productLanguage/handler"
	repository_productLanguage "github.com/sing3demons/product.product.sync/productLanguage/repository"
	service_productLanguage "github.com/sing3demons/product.product.sync/productLanguage/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	if os.Getenv("ZONE") == "PROD" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
	}

	r := gin.Default()

	{
		col, err := product_db.ConnectMonoDB()
		if err != nil {
			panic(err)
		}
		col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.M{"id": 1},
		})
		repo := repository_product.NewProduct(col)
		service := service_product.NewProductService(repo)
		handler := handler_product.NewProduct(service)

		r.GET("/products", handler.FindAllProduct)
		r.GET("/products/:id", handler.FindProduct)
		r.POST("/products", handler.CreateProduct)
	}
	{
		col, err := productLanguage_db.ConnectMonoDB()

		if err != nil {
			panic(err)
		}

		col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.M{
				"id": 1,
			},
			Options: nil,
		})

		repo := repository_productLanguage.NewProductLanguage(col)
		service := service_productLanguage.NewProductLanguageService(repo)
		handler := handler_productLanguage.NewProductLanguage(service)

		r.GET("/productLanguage", handler.FindAllCategory)
		r.GET("/productLanguage/:id", handler.FindCategory)
		r.POST("/productLanguage", handler.CreateCategory)
		// 	topics := []string{
		// 		constants.CREATE_CATEGORY,
		// 		constants.CREATE_CATEGORY_FAILED,
		// 		constants.CREATE_CATEGORY_SUCCESS,
		// 		constants.UPDATE_CATEGORY,
		// 		constants.UPDATE_CATEGORY_FAILED,
		// 		constants.UPDATE_CATEGORY_SUCCESS,
		// 	}
		// 	servers := "localhost:9092"
		// 	groupID := "go.category.product"
		// 	timeOut := time.Duration(-1)
		// 	consume.Consumer(servers, groupID, topics, timeOut)
	}

	{
		col, err := category_db.ConnectMonoDB()
		if err != nil {
			panic(err)
		}

		col.Indexes().CreateOne(context.Background(), mongo.IndexModel{
			Keys: bson.M{
				"id": 1,
			},
			Options: nil,
		})

		repo := repository_category.NewCategory(col)
		service := service_category.NewCategoryService(repo)
		handler := handler_category.NewCategory(service)

		r.GET("/category", handler.FindAllCategory)
		r.GET("/category/:id", handler.FindCategory)
		r.POST("/category", handler.CreateCategory)
		r.PATCH("/category/:id", handler.UpdateCategory)
	}
	RunServer(":2566", "sync-product-service", r)
}

func RunServer(addr, serviceName string, router http.Handler) {
	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		fmt.Printf("[%s] http listen: %s\n", serviceName, srv.Addr)

		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server listen err: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown: ", err)
	}

	fmt.Println("server exited")
}
