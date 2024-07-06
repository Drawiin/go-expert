package database

import (
	"database/sql"

	"github.com/google/uuid"
)

type CourseDB struct {
	db          *sql.DB
	ID          string
	Name        string
	Description string
	CategoryID  string
}

func NewCourse(db *sql.DB) *CourseDB {
	return &CourseDB{db: db}
}

func (c *CourseDB) Create(name, description, categoryID string) (CourseDB, error) {
	_, err := c.db.Exec("CREATE TABLE IF NOT EXISTS courses (id TEXT PRIMARY KEY, name TEXT, description TEXT, category_id TEXT)")
	if err != nil {
		return CourseDB{}, err
	}
	id := uuid.New().String()
	_, err = c.db.Exec("INSERT INTO courses (id, name, description, category_id) VALUES ($1, $2, $3, $4)", id, name, description, categoryID)
	if err != nil {
		return CourseDB{}, err
	}
	return CourseDB{ID: id, Name: name, Description: description, CategoryID: categoryID}, nil
}

func (c *CourseDB) GetById(id string) (CourseDB, error) {
	row := c.db.QueryRow("SELECT id, name, description, category_id FROM courses WHERE id = $1", id)
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.CategoryID)
	if err != nil {
		return CourseDB{}, err
	}
	return *c, nil
}

func (c *CourseDB) GetAll() ([]CourseDB, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []CourseDB
	for rows.Next() {
		var course CourseDB
		err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.CategoryID)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}

func (c *CourseDB) GetByCategoryId(categoryID string) ([]CourseDB, error) {
	rows, err := c.db.Query("SELECT id, name, description, category_id FROM courses WHERE category_id = $1", categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var courses []CourseDB
	for rows.Next() {
		var course CourseDB
		err := rows.Scan(&course.ID, &course.Name, &course.Description, &course.CategoryID)
		if err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}
	return courses, nil
}