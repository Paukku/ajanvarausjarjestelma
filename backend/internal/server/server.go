package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/handler"
	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/repository"
	"github.com/Paukku/ajanvarausjarjestelma/backend/internal/user/service"
	pbHTTP "github.com/Paukku/ajanvarausjarjestelma/backend/pb/http"
)

var jwtSecret []byte

func Run() {
	_ = godotenv.Load()
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewPostgresUserRepository(db)
	userService := service.NewUserServiceServer(userRepo)
	userHandler := handler.NewUserHandler(userService)

	converter := pbHTTP.NewBusinessCustomerAPIHTTPConverter(userHandler)
	mux := http.NewServeMux()
	RegisterRoutes(mux, converter)

	fmt.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
