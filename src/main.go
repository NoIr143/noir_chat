package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/noir143/noir_chat/src/database"
	"github.com/noir143/noir_chat/src/modules/features/auth"
	"github.com/noir143/noir_chat/src/modules/features/users"
)

func main() {
	mainDB := initMainDB()
	defer mainDB.Close()
	mux := http.NewServeMux()

	appModule(mainDB, mux)

	fmt.Println("Server running at :8080")
	http.ListenAndServe(":8080", mux)
}

func initMainDB() *sql.DB {
	db, err := database.NewPostgresStorage()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected!")
	return db
}

func appModule(mainDB *sql.DB, mux *http.ServeMux) {
	users.UserModule(mainDB, mux)
	auth.AuthModule(mainDB, mux)
}
