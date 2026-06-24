package proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	helper "github.com/S-Unknown047/LoadBalancer/internal/helper"
	model "github.com/S-Unknown047/LoadBalancer/internal/model"
	routingalgo "github.com/S-Unknown047/LoadBalancer/internal/routingAlgo"
)

var b *model.Backend = &helper.BackendData

func CaptureProxy(w http.ResponseWriter, r *http.Request) {

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			IdleConnTimeout:     90 * time.Second,
			TLSHandshakeTimeout: 10 * time.Second,
		},
		Timeout: 30 * time.Second,
	}
	url := r.URL
	// Clientport := url.Port()
	receiverUrlPath := url.Path
	Serverhost := url.Host

	dstIp, dstPort := routingalgo.GetServerIpAndPort(b)

	url_ := fmt.Sprintf("http://%s:%s%s", dstIp, dstPort, filterPath(receiverUrlPath))

	req, err := http.NewRequest(r.Method, url_, r.Body)

	if err != nil {
		log.Fatal(err)
	}

	for k, v := range r.Header {
		if k != "X-Forwarded-For" && k != "X-Forwarded-Host" && k != "X-Forwarded-Proto" {
			req.Header[k] = v
		}
	}

	ipAddr, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("error while getting address")
	}

	setHeader := func(out *http.Request) {
		out.Header["X-Forwarded-For"] = []string{ipAddr}
		out.Header["X-Forward-Host"] = []string{Serverhost}
		out.Header["X-Forward-Proto"] = []string{r.Proto}
	}

	setHeader(req)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	for k, val := range res.Header {
		if k != "X-Forwarded-For" && k != "X-Forwarded-Host" && k != "X-Forwarded-Proto" {
			w.Header()[k] = val
		}
	}

	w.WriteHeader(res.StatusCode)

	io.Copy(w, res.Body)
}

func filterPath(path string) string {
	var Path_ string
	Path_ = ""
	for i := 18; i < len(path); i++ {
		Path_ += string(path[i])
	}
	return Path_
}
