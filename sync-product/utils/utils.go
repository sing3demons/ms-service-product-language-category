package utils

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetHost() string {
	host := os.Getenv("HOST")
	if host == "" {
		host = "http://localhost:2566"
	}
	return host
}

func Href(key, value string) string {
	return fmt.Sprintf("%s/%s/%s", GetHost(), strings.ToLower(key), value)
}

func ConvertTimeBangkok(dataTime string) string {
	loc, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		fmt.Println("error load location")
		fmt.Println(err)
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", dataTime)
	if err != nil {
		fmt.Println("error parse time")
		fmt.Println(err)
	}
	return t.In(loc).Format("2006-01-02T15:04:05Z07:00")
}

func GetTransactionID() string {
	return time.Now().Format("20060102150405")
}

func GetMultiple[T any](db *mongo.Collection, filter primitive.D) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cur, err := db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var results []T
	for cur.Next(ctx) {
		var result T
		err := cur.Decode(&result)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

func GetOne[T any](db *mongo.Collection, filter primitive.D) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var result T
	err := db.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
