package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/drawiin/go-expert/crud/internal/dto"
	"github.com/drawiin/go-expert/crud/internal/entity"
	"github.com/drawiin/go-expert/crud/internal/infra/database"
	"github.com/go-chi/jwtauth"
)

type UserHandler struct {
	UserDB database.UserDBInterface
}

func NewUserHandler(userDB database.UserDBInterface) *UserHandler {
	return &UserHandler{UserDB: userDB}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUserInput true "User object that needs to be created"
// @Success 201 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dto.CreateUserInput
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	u, err := entity.NewUser(user.Name, user.Email, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.UserDB.Create(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Seed(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < 10; i++ {
		u, err := entity.NewUser(fmt.Sprintf("user%v", i), fmt.Sprintf("user%v@email.com", i), "123456")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = h.UserDB.Create(u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	jwt := r.Context().Value("jwt").(*jwtauth.JWTAuth)
	jwtExpiresIn := r.Context().Value("jwtExpiresIn").(int)

	var login dto.LoginInput
	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	u, err := h.UserDB.FindByEmail(login.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	valid := u.ValidatePassword(login.Password)
	if !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	_, token, _ := jwt.Encode(map[string]interface{}{
		"sub": u.ID.String(),
		"exp": time.Now().Add(time.Duration(jwtExpiresIn) * time.Second).Unix(),
	})

	loginResponse := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: token,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponse)
}
