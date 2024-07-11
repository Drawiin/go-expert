package service

import (
	"context"

	"github.com/drawiin/go-expert/grpc/internal/database"
	"github.com/drawiin/go-expert/grpc/internal/pb"
)

type CategoryService struct {
	pb.CategoryServiceServer
	CategoryDB *database.CategoryDB
}

func NewCategoryService(db *database.CategoryDB) *CategoryService {
	return &CategoryService{CategoryDB: db}
}

func (s *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := s.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}
