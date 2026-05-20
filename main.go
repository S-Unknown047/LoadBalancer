package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	model "github.com/S-Unknown047/LoadBalancer/Model"
	nat "github.com/S-Unknown047/LoadBalancer/NAT_MODE"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	server := http.NewServeMux()

	ServerCount := os.Getenv("SERVER_TURN_ROUND_ROBIN")

	server.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("this is Main Gett")
		w.Write([]byte("<h1> this is header </h1>"))
	})

	server.HandleFunc("POST /server", func(w http.ResponseWriter, r *http.Request) {
		var server model.ReqServer
		err := json.NewDecoder(r.Body).Decode(server)

		if err != nil {
			fmt.Println("Error in json Decoder")
		}

		fmt.Println("value ", server)
	})

	go nat.Test()

	port := os.Getenv("APP_PORT")
	fmt.Println("running on server ", port)
	log.Fatal(http.ListenAndServe(port, server))

}
