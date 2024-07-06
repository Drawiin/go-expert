package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type CategoryDB struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
}

func NewCategory(db *sql.DB) *CategoryDB {
	return &CategoryDB{db: db}
}

func (c *CategoryDB) Create(name, description string) (CategoryDB, error) {
	id := uuid.New().String()
	_, err := c.db.Exec("CREATE TABLE IF NOT EXISTS categories (id TEXT PRIMARY KEY, name TEXT, description TEXT)")
	if err != nil {
		return CategoryDB{}, err
	}
	_, err = c.db.Exec("INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)", id, name, description)
	if err != nil {
		return CategoryDB{}, err
	}
	return CategoryDB{ID: id, Name: name, Description: description}, nil
}

func (c *CategoryDB) GetById(id string) (CategoryDB, error) {
	row := c.db.QueryRow("SELECT id, name, description FROM categories WHERE id = $1", id)
	err := row.Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return CategoryDB{}, err
	}
	return *c, nil
}

func (c *CategoryDB) GetAll() ([]CategoryDB, error) {
	rows, err := c.db.Query("SELECT id, name, description FROM categories")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var categories []CategoryDB
	for rows.Next() {
		var category CategoryDB
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}
