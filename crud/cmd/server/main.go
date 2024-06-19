package main

import (
	"log"
	"net/http"

	"github.com/drawiin/go-expert/crud/config"
	"github.com/drawiin/go-expert/crud/internal/entity"
	"github.com/drawiin/go-expert/crud/internal/infra/database"
	"github.com/drawiin/go-expert/crud/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	// LoadConfig function is called to load the configuration values from the .env file
	_, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := gorm.Open((sqlite.Open("test.db")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})

	productDb := database.NewProductDB(db)
	productHandler := handlers.NewProductHandler(productDb)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Get("/products", productHandler.GetAllProducts)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)
	r.Post("/producs/seed", productHandler.Seed)

	http.ListenAndServe(":8080", r)
}
