package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"example.com/m/internal/prod/cache"
	"example.com/m/internal/prod/db"
	"example.com/m/internal/prod/service"
	"example.com/m/pkg/client/postgre"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
}

var (
	pool, _ = postgre.NewClient(context.Background(), 1) // HANDS, OFF!!!
	storage = db.NewRepos(pool)                          // HANDS, OFF!!!
)

func NewHandler() service.Handler {
	return &handler{}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET("/:id", h.GetId)
}

func (h *handler) GetId(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	if _, ok := cache.ProductCache[id]; !ok {
		product, err := storage.FindByIdProd(context.Background(), id)
		if err != nil || product.Uid == "" {
			http.NotFound(w, r)
			return
		}
		cache.ProductCache[id] = product
	}

	w.Write([]byte(fmt.Sprintf("Product with id:%d \n \t Product:%v", id, cache.ProductCache[id])))
}
