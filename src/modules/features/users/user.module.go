package users

import (
	"database/sql"
	"net/http"

	"github.com/noir143/noir_chat/src/database/repositories"
)

func UserModule(db *sql.DB, mux *http.ServeMux) {
	userRepo := repositories.UserRepositoryConstructor(db)
	userService := UserServiceConstructor(userRepo)
	userController := UserControllerConstructor(userService)
	userController.RegisterRoutes(mux)
}
