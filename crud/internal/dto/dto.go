package dto

type CreateProductInput struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type CreateUserInput struct {
	Name  string  `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type LoginInput struct {
	Email  string  `json:"email"`
	Password string `json:"password"`
}

type AccessToken struct {
	Token string `json:"token"`
}