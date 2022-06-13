package service

import (
	"context"
	"encoding/json"
	"log"

	"example.com/m/internal/prod/models"
)

func InputDataConvert(ctx context.Context, data []byte) models.Product {
	var product models.Product
	if err := json.Unmarshal(data, &product); err != nil {
		log.Println(err)
	}
	return product
}
