package service

import (
	"context"

	"example.com/m/internal/prod/models"
)

type Storage interface {
	InsertProd(ctx context.Context, product *models.Product) error
	FindByIdProd(ctx context.Context, product_id int) (models.Product, error)
	GetAllDataDB(ctx context.Context) ([]models.Product, error)
}
