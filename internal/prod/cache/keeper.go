package cache

import (
	"context"
	"log"

	"example.com/m/internal/prod/models"
	"example.com/m/internal/prod/service"
)

var ProductCache = make(map[int]models.Product)

func CheckOut(ctx context.Context, storage service.Storage) {
	allData, err := service.Storage.GetAllDataDB(storage, context.Background())
	if err != nil {
		log.Println(err)
	}
	for i, value := range allData {
		ProductCache[i] = value //i+1
	}
}
