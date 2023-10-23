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
	"github.com/sing3demons/product/db"
	"github.com/sing3demons/product/product/handler"
	"github.com/sing3demons/product/product/repository"
	"github.com/sing3demons/product/product/service"

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

	repo := repository.NewProduct(col)
	service := service.NewProductService(repo)
	handler := handler.NewProduct(service)

	r := gin.Default()
	r.GET("/products", handler.FindAllProduct)
	r.GET("/products/:id", handler.FindProduct)
	r.POST("/products", handler.CreateProduct)

	serveHttp(":2566", "product-service", r)

}

func serveHttp(addr, serviceName string, router http.Handler) {
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
