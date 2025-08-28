package users

import (
	"net/http"

	"github.com/noir143/noir_chat/src/middlewares"
	"github.com/noir143/noir_chat/src/modules/features/users/dtos"
	"github.com/noir143/noir_chat/src/shared/utils"
)

type UserController struct {
	userService *UserService
}

func UserControllerConstructor(userService *UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

func (userController *UserController) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /users", middlewares.AuthMiddleware(http.HandlerFunc(userController.GetUsers)))
	mux.Handle("POST /users", middlewares.AuthMiddleware(http.HandlerFunc(userController.CreateUser)))
	mux.Handle("GET /users/{id}", middlewares.AuthMiddleware(http.HandlerFunc(userController.GetUser)))
	mux.Handle("PUT /users/{id}", middlewares.AuthMiddleware(http.HandlerFunc(userController.UpdateUser)))
	mux.Handle("DELETE /users/{id}", middlewares.AuthMiddleware(http.HandlerFunc(userController.DeleteUser)))
}

// CreateUser handles POST /users
func (userController *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userCreateDTO dtos.UserCreateDTO
	if err := utils.ParseJSON(r, &userCreateDTO); err != nil {
		utils.WriteError(w, err)
		return
	}

	userController.userService.CreateUser(userCreateDTO)
}

// GetUsers handles GET /users
func (userController *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := userController.userService.GetUsers(1, 10)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, users)
}

// GetUser handles GET /users/{id}
func (userController *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	userController.userService.GetUserByID(1)
}

// UpdateUser handles PUT /users/{id}
func (userController *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var userUpdateDTO dtos.UserUpdateDTO
	if err := utils.ParseJSON(r, &userUpdateDTO); err != nil {
		utils.WriteError(w, err)
		return
	}
	userController.userService.UpdateUser(1, userUpdateDTO)
}

// DeleteUser handles DELETE /users/{id}
func (userController *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userController.userService.DeleteUser(1)
}
