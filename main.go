package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	server := http.NewServeMux()
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal(err)
	}

	server.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1> This is Server Page </h1>"))
	})
	port := os.Getenv("APP_PORT")
	serverReq := http.Server{
		Addr:    port,
		Handler: server,
	}

	log.Fatal(serverReq.ListenAndServe())
}
