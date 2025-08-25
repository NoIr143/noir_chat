package auth

import (
	"database/sql"
	"net/http"

	"github.com/noir143/noir_chat/src/database/repositories"
)

func AuthModule(db *sql.DB, mux *http.ServeMux) {
	userRepo := repositories.UserRepositoryConstructor(db)
	authService := AuthServiceConstructor(userRepo)
	authController := AuthControllerConstructor(authService)
	authController.RegisterRoutes(mux)
}
