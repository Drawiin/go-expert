package database

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/drawiin/go-expert/crud/internal/entity"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreate(t *testing.T) {
	db := CreateProductTestDB(t)
	product, _ := entity.NewProduct("product1", 10)
	productDB := NewProductDB(db)

	err := productDB.Create(product)
	assert.Nil(t, err)

	var productFound entity.Product
	err = db.First(&productFound, "id = ?", product.ID).Error

	assert.Nil(t, err)
	assert.NotNil(t, productFound)
	assert.Equal(t, product.Name, productFound.Name)
	assert.Equal(t, product.Price, productFound.Price)
	assert.Equal(t, product.ID, productFound.ID)
}

func TestFindAll(t *testing.T) {
	db := CreateProductTestDB(t)
	productDB := NewProductDB(db)
	fakeProducts := CreateFakeProductList(25)
	for _, p := range fakeProducts {
		db.Create(&p)
	}

	result, err := productDB.FindAll(1, 10, "asc")

	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, 10, len(result))
	assert.Equal(t, fakeProducts[0].Name, result[0].Name)
	assert.Equal(t, fakeProducts[9].Name, result[9].Name)
}

func TestFindAllDesc(t *testing.T) {
	db := CreateProductTestDB(t)
	productDB := NewProductDB(db)
	fakeProducts := CreateFakeProductList(25)
	for _, p := range fakeProducts {
		db.Create(&p)
	}

	result, err := productDB.FindAll(1, 10, "desc")

	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, 10, len(result))
	assert.Equal(t, fakeProducts[len(fakeProducts) -1].Name, result[0].Name)
}

func TestFindAllNextFullPage(t *testing.T) {
	db := CreateProductTestDB(t)
	productDB := NewProductDB(db)
	fakeProducts := CreateFakeProductList(25)
	for _, p := range fakeProducts {
		db.Create(&p)
	}

	result, err := productDB.FindAll(2, 10, "asc")

	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, 10, len(result))
	assert.Equal(t, fakeProducts[10].Name, result[0].Name)
	assert.Equal(t, fakeProducts[19].Name, result[9].Name)
}

func TestFindAllNextPartialPage(t *testing.T) {
	db := CreateProductTestDB(t)
	productDB := NewProductDB(db)
	fakeProducts := CreateFakeProductList(25)
	for _, p := range fakeProducts {
		db.Create(&p)
	}

	result, err := productDB.FindAll(3, 10, "asc")

	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, 5, len(result))
	assert.Equal(t, fakeProducts[20].Name, result[0].Name)
	assert.Equal(t, fakeProducts[24].Name, result[4].Name)
}

func TestFindAllDefault(t *testing.T) {
	db := CreateProductTestDB(t)
	productDB := NewProductDB(db)
	fakeProducts := CreateFakeProductList(25)
	for _, p := range fakeProducts {
		db.Create(&p)
	}

	result, err := productDB.FindAll(0, -1, "")

	assert.Nil(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, 10, len(result))
	assert.Equal(t, fakeProducts[0].Name, result[0].Name)
}

func TestFindAllEmpty(t *testing.T) {
	db := CreateProductTestDB(t)
	productDB := NewProductDB(db)

	result, err := productDB.FindAll(1, 10, "asc")

	assert.Nil(t, err)
	assert.Empty(t, result)
}

func CreateProductTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to the database: %v", err)
	}
	db.AutoMigrate(&entity.Product{})
	return db
}

func CreateFakeProductList(len int) []entity.Product {
	var products []entity.Product
	for i := 0; i < len; i++ {
		product, err := entity.NewProduct(fmt.Sprintf("product%d", i), rand.Float64()*100)
		if err != nil {
			panic(err)
		}
		products = append(products, *product)
	}
	return products
}
