package router

import (
	"net/http"

	reverseproxy "github.com/S-Unknown047/LoadBalancer/L7_loadBalancer/ReverseProxy"
	middlewar "github.com/S-Unknown047/LoadBalancer/Middleware"
)

func Routing() *http.ServeMux {

	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Write([]byte("<h1> This is home Page</h1>"))
	})

	router.HandleFunc("POST /server", middlewar.GetServerPath)

	router.HandleFunc("POST /server/setup", middlewar.GetServerSetup)

	router.HandleFunc("/proxyLoadBalancer", reverseproxy.CaptureProxy)

	return router
}
