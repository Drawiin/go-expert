package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/drawiin/go-expert/crud/internal/dto"
	"github.com/drawiin/go-expert/crud/internal/entity"
	"github.com/drawiin/go-expert/crud/internal/infra/database"
	entityPkd "github.com/drawiin/go-expert/crud/pkg/entity"
	"github.com/go-chi/chi/v5"
)

type ProductHandler struct {
	ProductDB database.ProductDBInterface
}

func NewProductHandler(db database.ProductDBInterface) *ProductHandler {
	return &ProductHandler{
		ProductDB: db,
	}
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Adds a new product to the database
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductInput true "Create Product Input"
// @Success 201 {object} string "Product created successfully"
// @Failure 400 {object} string "Invalid request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /products [post]
// @Security jwt
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product dto.CreateProductInput
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	p, err := entity.NewProduct(product.Name, product.Price)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.ProductDB.Create(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	p, err := h.ProductDB.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(p)
}

// GetAllProducts godoc
// @Summary Get all products
// @Description Retrieves a list of products with pagination and sorting
// @Tags products
// @Accept json
// @Produce json
// @Param page query int false "Page number for pagination"
// @Param limit query int false "Limit of products per page for pagination"
// @Param sort query string false "Sort order"
// @Success 200 {array} []entity.Product "List of products"
// @Failure 500 {object} string "Internal Server Error"
// @Router /products [get]
// @Security jwt
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	sort := r.URL.Query().Get("sort")

	products, err := h.ProductDB.FindAll(page, limit, sort)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Updates a product by ID
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body entity.Product true "Product object to update"
// @Success 200 {string} string "OK"
// @Failure 400 {object} string "Invalid request"
// @Failure 404 {object} string "Product not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /products/{id} [put]
// @Security jwt
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}
	var updateProduct entity.Product
	err := json.NewDecoder(r.Body).Decode(&updateProduct)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	updateProduct.ID, err = entityPkd.ParseID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.ProductDB.FindByID(updateProduct.ID.String())
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = h.ProductDB.Update(&updateProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Deletes a product by ID
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} string "Product deleted successfully"
// @Failure 400 {object} string "Bad Request"
// @Failure 500 {object} string "Internal Server Error"
// @Router /products/{id} [delete]
// @Security jwt
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
	}

	err := h.ProductDB.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}


// @Summary Seed the database with products
// @Description Seeds the database with a specified number of products
// @Tags products
// @Accept json
// @Produce json
// @Param size query int false "Number of products to create (default: 25)"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /products/seed [post]
// @Security jwt
func (h *ProductHandler) Seed(w http.ResponseWriter, r *http.Request) {
	size, err := strconv.Atoi(r.URL.Query().Get("size"))
	if err != nil {
		size = 25
	}
	for i := 0; i < size; i++ {
		go func() {
			product, err := entity.NewProduct(fmt.Sprintf("product%d", i), rand.Float64()*100)
			if err != nil {
				panic(err)
			}
			err = h.ProductDB.Create(product)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()
	}
	w.WriteHeader(http.StatusOK)
}
