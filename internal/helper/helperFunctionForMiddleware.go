package helper

import (
	"log"
	"os"
	"time"

	model "github.com/S-Unknown047/LoadBalancer/internal/model"
)

//	func Check(ip string, b *model.Backend) *model.Server {
//		for _, sIp := range b.Servers {
//			if sIp.Url == ip {
//				return &sIp
//			}
//		}
//		return nil
//	}

var BackendData model.Backend
var ServerData []model.Server

func HandelServer(obj *[]model.ReqServer) {
	newServer := storingServer(obj)
	ServerData = append(ServerData, newServer...)

}

func FlagBackend() chan bool {
	var FlagChan = make(chan bool)
	FlagChan <- true
	return FlagChan
}

func FlagServer() chan bool {
	var FlagChan = make(chan bool)
	FlagChan <- true
	return FlagChan
}

func storingServer(obj *[]model.ReqServer) []model.Server {
	defer FlagServer()

	output := make([]model.Server, 0, len(*obj))
	for _, val := range *obj {
		ip := val.IP
		port := val.Port
		// there can be some code for getting response time so that we can use the algo
		// which choose server based on response time
		temp := model.Server{
			IP:         ip,
			Port:       port,
			Connection: 0,
			Weight:     val.Weight,
		}
		output = append(output, temp)
	}

	return output
}

func HandelSetup(obj *model.ReqSetup) {
	BackendData = *backendSetup(obj, &ServerData)
}

func backendSetup(obj *model.ReqSetup, server *[]model.Server) *model.Backend {
	defer FlagBackend()

	if len(*server) == 0 {
		log.Fatal("Error no server Present use addserver cmd to add servers")
		os.Exit(1)
	}
	var tempBackend model.Backend
	tempBackend.Servers = server

	algo := obj.Algo

	if algo == "rr" || algo == "RoundRobin" {
		algo = "RoundRobin"
	} else if algo == "lc" || algo == "LeastConnection" {
		algo = "LeastConnection"
	}

	tempBackend.Algo = algo

	// tempBackend.Mode = obj.ModeBalance
	tempBackend.TotalServer = (uint64)(len((*server)))
	tempBackend.TotalServerConnection = 0

	Checker()

	if algo == "LeastConnection" {
		EnterData(&ServerData)
	}

	return &tempBackend
}

func Checker() {
	ServerupAliveServer(&ServerData)
	ServerData = ServerTest()
	ticker := time.NewTicker(20 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				ServerData = ServerTest()
				if BackendData.Algo == "LeastConnection" {
					EnterData(&ServerData)
				}
			}
		}
	}()
}
