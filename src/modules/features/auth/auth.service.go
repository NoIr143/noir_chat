package auth

import (
	"database/sql"

	"github.com/noir143/noir_chat/src/configs"
	"github.com/noir143/noir_chat/src/database/entities"
	"github.com/noir143/noir_chat/src/database/repositories"
	authDto "github.com/noir143/noir_chat/src/modules/features/auth/dtos"
	"github.com/noir143/noir_chat/src/shared/exceptions"
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

func (authService *AuthService) Register(request authDto.RegisterDTO) (entities.User, any) {
	exists, err := authService.userRepo.UserExists(request.Email)
	if err != nil {
		return entities.User{}, exceptions.InternalException{Error: err}
	}

	if exists {
		return entities.User{}, exceptions.BadRequestException{ErrorId: "NC_0404", Message: "user_email_already_exist"}
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return entities.User{}, exceptions.InternalException{Error: err}
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
		return entities.User{}, exceptions.InternalException{Error: err}
	}

	return createdUser, nil
}

func (authService *AuthService) Login(request authDto.LoginDTO) (authDto.LoginResponseDto, any) {
	user, err := authService.userRepo.GetByEmail(request.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return authDto.LoginResponseDto{}, exceptions.BadRequestException{ErrorId: "NC_600", Message: "invalid_credentials", Error: err}
		}
		return authDto.LoginResponseDto{}, exceptions.InternalException{Error: err}
	}

	if !utils.ComparePasswords(user.HashedPassword, []byte(request.Password)) {
		return authDto.LoginResponseDto{}, exceptions.BadRequestException{ErrorId: "NC_600", Message: "invalid_credentials", Error: err}
	}

	secret := []byte(configs.EnvConfigs.JWT_SECRET)
	token, err := utils.CreateJWT(secret, user.ID)
	if err != nil {
		return authDto.LoginResponseDto{}, exceptions.InternalException{Error: err}
	}

	return authDto.LoginResponseDto{AccessToken: token}, nil

}
