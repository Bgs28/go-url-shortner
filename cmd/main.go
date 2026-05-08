package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"url-shortener/internal/handler"
	"url-shortener/internal/repository"
	"url-shortener/internal/service"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	// load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// connect to database
	db := initDB()
	defer db.Close()

	// wire up layers
	urlRepo := repository.NewURLRepository(db)
	urlService := service.NewURLService(urlRepo, os.Getenv("BASE_URL"))
	URLHandler := handler.NewURLHandler(urlService)

	// routes
	http.HandleFunc("/shorten", URLHandler.Shorten)
	http.HandleFunc("/stats/", URLHandler.GetStats)
	http.HandleFunc("/", URLHandler.Redirect)

	// start server
	port := os.Getenv("APP_PORT")
	fmt.Println("Server running on port", port)
	log.Fatal(http.ListenAndServe(":" + port, nil))
}

func initDB() *sql.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_HOST"),
	os.Getenv("DB_PORT"),
	os.Getenv("DB_NAME"),
)

db, err := sql.Open("mysql", dsn)
if err != nil {
	log.Fatal("Failed to connect to database: ", err)
}

if err := db.Ping(); err != nil{
	log.Fatal("Database unreachable: ", err)
}

fmt.Println("Database Connected!")
return db
}