package main

import (
	"fmt"
	"gounter/api/handler"
	"gounter/api/route"
	"gounter/internal/repository"
	"gounter/internal/service"
	"gounter/util"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	token, err := util.GenerateValidJWT()
	if err != nil {
		log.Fatalf("Error generating JWT token: %v", err)
	}
	fmt.Println("Generated JWT token (valid for 5 minutes): ", token)

	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable", config.User, config.Password, config.DBName, config.Host, config.Port))
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}

	defer db.Close()

	counterRepo := repository.New(db)
	service := service.NewCounterService(counterRepo)
	counterHandler := handler.NewHandler(service)

	routes := route.InitRoutes(counterHandler)

	// Start the HTTP server on port 8081
	log.Println("Starting server on :8081...")
	if err := http.ListenAndServe(":8081", routes); err != nil {
		log.Fatal("Error starting server:", err)
	}
}
