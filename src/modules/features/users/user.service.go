package users

import (
	"github.com/noir143/noir_chat/src/database/entities"
	"github.com/noir143/noir_chat/src/database/repositories"
	"github.com/noir143/noir_chat/src/modules/features/users/dtos"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func UserServiceConstructor(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// CreateUser creates a new user
func (userService *UserService) CreateUser(request dtos.UserCreateDTO) (entities.User, error) {
	user := entities.User{
		Name:  request.Name,
		Email: request.Email,
	}

	return userService.userRepo.Create(user)
}

// GetUsers returns paginated users
func (userService *UserService) GetUsers(page, pageSize int) ([]entities.User, error) {
	// return userService.userRepo.GetPaginated(page, pageSize)
	users := []entities.User{
		{ID: 1, Name: "Alice", Email: "alice@gmail.com"},
		{ID: 2, Name: "Bob", Email: "bob@gmail.com"},
	}
	return users, nil
}

// GetUserByID returns a user by ID
func (userService *UserService) GetUserByID(id int) (entities.User, error) {
	return userService.userRepo.GetByID(id)
}

// GetUserByEmail returns a user by email
func (userService *UserService) GetUserByEmail(email string) (entities.User, error) {
	return userService.userRepo.GetByEmail(email)
}

// UpdateUser updates a user
func (userService *UserService) UpdateUser(id int, request dtos.UserUpdateDTO) (entities.User, error) {
	user := entities.User{
		Name:  request.Name,
		Email: request.Email,
	}
	return userService.userRepo.Update(id, user)
}

// DeleteUser deletes a user
func (userService *UserService) DeleteUser(id int) error {
	return userService.userRepo.Delete(id)
}
