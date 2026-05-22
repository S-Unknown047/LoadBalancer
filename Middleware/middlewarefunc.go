package middleware

import (
	"encoding/json"
	"log"
	"net/http"

	helper "github.com/S-Unknown047/LoadBalancer/Helper"
	model "github.com/S-Unknown047/LoadBalancer/Model"
)

func GetServerPath(w http.ResponseWriter, r *http.Request) {
	var output []model.ReqServer

	err := json.NewDecoder(r.Body).Decode(&output)

	if err != nil {
		log.Fatal(err)
	}

	helper.HandelServer(&output)
}

func GetServerSetup(w http.ResponseWriter, r *http.Request) {
	var setup model.ReqSetup

	err := json.NewDecoder(r.Body).Decode(&setup)

	if err != nil {
		log.Fatal(err)
	}

	helper.HandelSetup(&setup)
}
