package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"example.com/m/internal/prod/cache"
	"example.com/m/internal/prod/dataworker"
	"example.com/m/internal/prod/db"
	"example.com/m/internal/prod/handlers"
	"example.com/m/internal/prod/models"
	"example.com/m/internal/prod/nats"
	"example.com/m/pkg/client/postgre"
	"github.com/julienschmidt/httprouter"
)

func main() {
	postgeSqlClient, err := postgre.NewClient(context.Background(), 5)
	if err != nil {
		log.Fatal(err)
	}
	storage := db.NewRepos(postgeSqlClient)

	product := models.Product{}
	q := make(chan []byte)
	go func() {
		router := httprouter.New()
		handler := handlers.NewHandler()
		handler.Register(router)

		startServ(router)
	}()

	// cache.CheckOut(context.Background(), storage)

	cacheKey := len(cache.ProductCache) + 1

	// fmt.Printf("Len of Cash from start:%d \n", len(cache.ProductCache))

	sc := nats.ConnectStan("client-124")

	go nats.GetData("foo", "test", "c", sc, q)

	for value := range q {
		product = dataworker.ConvertToProd(value)
		if product.Uid != "" {
			cache.ProductCache[cacheKey] = product
			// fmt.Println(len(cache.ProductCache))
			storage.InsertProd(context.Background(), &product)
			cacheKey++
		}
	}
}

func startServ(router *httprouter.Router) {
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		panic(err)
	}
	server := &http.Server{
		Handler:      router,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}
	log.Println("Server is listening")
	log.Fatal(server.Serve(listener))

}
