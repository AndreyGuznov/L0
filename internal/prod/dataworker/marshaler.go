package dataworker

import (
	"encoding/json"
	"log"

	"example.com/m/internal/prod/models"
)

func ConvertToProd(data []byte) models.Product {
	var product models.Product
	err := json.Unmarshal(data, &product)
	if err != nil {
		log.Println("Undirect message in chanel")
		return models.Product{}
	}
	return product
}
