package auth

import (
	"fmt"

	"github.com/noir143/noir_chat/src/configs"
	"github.com/noir143/noir_chat/src/database/entities"
	"github.com/noir143/noir_chat/src/database/repositories"
	authDto "github.com/noir143/noir_chat/src/modules/features/auth/dtos"
	"github.com/noir143/noir_chat/src/shared/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func AuthServiceConstructor(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo,
	}
}

func (authService *AuthService) Register(request authDto.RegisterDTO) (entities.User, error) {
	exists, err := authService.userRepo.UserExists(request.Email)
	if err != nil {
		return entities.User{}, fmt.Errorf("internal sever")
	}

	if exists {
		return entities.User{}, fmt.Errorf("user with email %s already exists", request.Email)
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return entities.User{}, fmt.Errorf("internal server")
	}

	// Create new user
	user := entities.User{
		FirstName:      request.FirstName,
		LastName:       request.LastName,
		Email:          request.Email,
		HashedPassword: hashedPassword,
	}

	createdUser, err := authService.userRepo.Create(user)
	if err != nil {
		return entities.User{}, err
	}

	return createdUser, nil
}

func (authService *AuthService) Login(request authDto.LoginDTO) (authDto.LoginResponseDto, error) {
	user, err := authService.userRepo.GetByEmail(request.Email)
	if err != nil {
		return authDto.LoginResponseDto{}, fmt.Errorf("invalid credentials")
	}

	if !utils.ComparePasswords(user.HashedPassword, []byte(request.Password)) {
		return authDto.LoginResponseDto{}, fmt.Errorf("invalid email or password")
	}

	secret := []byte(configs.EnvConfigs.JWT_SECRET)
	token, err := utils.CreateJWT(secret, user.ID)
	if err != nil {
		return authDto.LoginResponseDto{}, fmt.Errorf("invalid credentiaks")
	}

	return authDto.LoginResponseDto{AccessToken: token}, nil

}
