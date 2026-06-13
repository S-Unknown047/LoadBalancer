package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	natmode "github.com/S-Unknown047/LoadBalancer/NAT_MODE"
	router "github.com/S-Unknown047/LoadBalancer/Router"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server := router.Routing()
	port := os.Getenv("APP_PORT")
	fmt.Println("running on server ", port)
	go natmode.Test()
	log.Fatal(http.ListenAndServe(port, server))
}
