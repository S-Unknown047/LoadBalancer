package router

import (
	"net/http"

	reverseproxy "github.com/S-Unknown047/LoadBalancer/internal/proxy/l7"
)

func Routing() *http.ServeMux {

	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
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

	router.HandleFunc("/l7/", reverseproxy.CaptureProxy)

	return router
}
