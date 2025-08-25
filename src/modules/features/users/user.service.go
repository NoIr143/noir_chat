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

func (userService *UserService) CreateUser(request dtos.UserCreateDTO) (entities.User, error) {
	return entities.User{}, nil
}

func (userService *UserService) GetUsers(page, pageSize int) ([]entities.User, error) {
	return []entities.User{}, nil
}

func (userService *UserService) GetUserByID(id int) (entities.User, error) {
	return userService.userRepo.GetByID(id)
}

func (userService *UserService) GetUserByEmail(email string) (entities.User, error) {
	return userService.userRepo.GetByEmail(email)
}

func (userService *UserService) UpdateUser(id int, request dtos.UserUpdateDTO) (entities.User, error) {
	return entities.User{}, nil
}

func (userService *UserService) DeleteUser(id int) error {
	return userService.userRepo.Delete(id)
}
