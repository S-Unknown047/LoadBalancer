package router

import (
	"net/http"

	middlewar "github.com/S-Unknown047/LoadBalancer/Middleware"
)

func Routing() *http.ServeMux {

	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("<h1> This is home Page</h1>"))
	})

	router.HandleFunc("POST /server", middlewar.GetServerPath)

	router.HandleFunc("POST /server/setup", middlewar.GetServerSetup)

	return router
}
