package product

type ProductUseCase struct {
	repository ProductRepositoryInterface
}

func NewProductUseCase(repository ProductRepositoryInterface) *ProductUseCase {
	return &ProductUseCase{
		repository: repository,
	}
}

func (uc *ProductUseCase) GetProduct(id int) (Product, error) {
	return uc.repository.GetProduct(id)
}
