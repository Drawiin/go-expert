package service

import (
	"context"
	"io"

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

func (s *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(categories)
			return nil
		}
		if err != nil {
			return err
		}
		createdCategory, err := s.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          createdCategory.ID,
			Name:        createdCategory.Name,
			Description: createdCategory.Description,
		})
	}
}

func (s *CategoryService) CreateCategoryStreamBidirectional(stream pb.CategoryService_CreateCategoryStreamBidirectionalServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		createdCategory, err := s.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		err = stream.Send(&pb.Category{
			Id:          createdCategory.ID,
			Name:        createdCategory.Name,
			Description: createdCategory.Description,
		})
		if err != nil {
			return err
		}
	}
}

func (s *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := s.CategoryDB.GetAll()
	if err != nil {
		return nil, err
	}
	var pbCategories []*pb.Category
	for _, category := range categories {
		pbCategories = append(pbCategories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}
	return &pb.CategoryList{Categories: pbCategories}, nil
}

func (s *CategoryService) GetCategoryById(ctx context.Context, in *pb.GetCategoryByIdRequest) (*pb.Category, error) {
	category, err := s.CategoryDB.GetById(in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}, nil
}
