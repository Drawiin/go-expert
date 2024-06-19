package database

import (
	"github.com/drawiin/go-expert/crud/internal/entity"
	"gorm.io/gorm"
)

type ProductDBInterface interface {
	Create(product *entity.Product) error
	FindAll(page, limit int, sort string) ([]entity.Product, error)
	FindByID(id string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id string) error
}

// Implements ProductDBInterface
type ProductDB struct {
	DB *gorm.DB
}

func NewProductDB(db *gorm.DB) *ProductDB {
	return &ProductDB{
		DB: db,
	}
}

func (db *ProductDB) Create(product *entity.Product) error {
	return db.DB.Create(product).Error
}

func (db *ProductDB) FindAll(page, limit int, sort string) ([]entity.Product, error) {
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if page <= 0 {
		page = 1
	}

	if limit <= 0 {
		limit = 10
	}

	var products []entity.Product
	err := db.DB.Limit(limit).Offset((page - 1) * limit).Order("created_at " + sort).Find(&products).Error
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (db *ProductDB) FindByID(id string) (*entity.Product, error) {
	var product entity.Product
	if err := db.DB.First(&product, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (db *ProductDB) Update(product *entity.Product) error {
	if _, err := db.FindByID(product.ID.String()); err != nil {
		return err
	}
	return db.DB.Save(product).Error
}

func (db *ProductDB) Delete(id string) error {
	product, err := db.FindByID(id)
	if err != nil {
		return err
	}
	return db.DB.Delete(product).Error
}
