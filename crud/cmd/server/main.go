package main

import (
	"log"
	"net/http"

	"github.com/drawiin/go-expert/crud/config"
	_ "github.com/drawiin/go-expert/crud/docs"
	"github.com/drawiin/go-expert/crud/internal/entity"
	"github.com/drawiin/go-expert/crud/internal/infra/database"
	"github.com/drawiin/go-expert/crud/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// @title Go Expert CRUD API
// @version 1.0
// @description Product API with authentication
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://swagger.io/support
// @contact.email email@email.com

// @license.name MIT
// @license.url http://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apiKey jwt
// @in header
// @name Authorization
func main() {
	// LoadConfig function is called to load the configuration values from the .env file
	config, err := config.LoadConfig(".")
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

	userDb := database.NewUserDB(db)
	userHandler := handlers.NewUserHandler(userDb)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.WithValue("jwt", config.TokenAuth))
	r.Use(middleware.WithValue("jwtExpiresIn", config.JWTExpiresIn))

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier((config.TokenAuth)))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Get("/", productHandler.GetAllProducts)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Post("/seed", productHandler.Seed)
	})

	r.Post("/users/signup", userHandler.CreateUser)
	r.Post("/users/seed", userHandler.Seed)
	r.Post("/users/login", userHandler.Login)

	// Aparently the doc.json in the end is not optional
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/docs/doc.json")))

	http.ListenAndServe(":8080", r)
}
