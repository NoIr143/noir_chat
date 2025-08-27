package auth

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	common "github.com/noir143/noir_chat/src/common/dtos"
	authDto "github.com/noir143/noir_chat/src/modules/features/auth/dtos"
	"github.com/noir143/noir_chat/src/shared/exceptions"
	"github.com/noir143/noir_chat/src/shared/utils"
)

type AuthController struct {
	authService *AuthService
}

func AuthControllerConstructor(authService *AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (authController *AuthController) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /auth/register", authController.register)
	mux.HandleFunc("POST /auth/login", authController.login)

	log.Printf("Loaded: POST /auth/register")
	log.Printf("Loaded: POST /auth/login")
}

func (authController *AuthController) register(w http.ResponseWriter, r *http.Request) {
	var registerDto authDto.RegisterDTO
	if err := utils.ParseJSON(r, &registerDto); err != nil {
		utils.WriteError(w, exceptions.InternalException{Error: err})
		return
	}

	if err := utils.Validate.Struct(registerDto); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, exceptions.InvalidParameterException{ValidationErrors: errors})
		return
	}

	user, err := authController.authService.Register(registerDto)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, common.UpdateResponseDTO{ID: user.ID})
}

func (authController *AuthController) login(w http.ResponseWriter, r *http.Request) {
	var loginDto authDto.LoginDTO
	if err := utils.ParseJSON(r, &loginDto); err != nil {
		utils.WriteError(w, exceptions.InternalException{Error: err})
		return
	}

	if err := utils.Validate.Struct(loginDto); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, exceptions.InvalidParameterException{ValidationErrors: errors})
		return
	}

	loginResponseDTO, err := authController.authService.Login(loginDto)
	if err != nil {
		utils.WriteError(w, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, loginResponseDTO)
}
