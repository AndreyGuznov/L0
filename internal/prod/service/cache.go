package service

import (
	"context"

	"example.com/m/internal/prod/models"
)

type Cach interface {
	CheckOut(context.Context, Storage) map[int]models.Product
}
