package helper

import (
	"sync"

	model "github.com/S-Unknown047/LoadBalancer/Model"
)

//	func Check(ip string, b *model.Backend) *model.Server {
//		for _, sIp := range b.Servers {
//			if sIp.Url == ip {
//				return &sIp
//			}
//		}
//		return nil
//	}
var serverCount uint64 = 0
var BackendData model.Backend
var ServerData []model.Server

func HandelServer(obj *[]model.ReqServer) {
	ServerData = *storingServer(obj)
}

func storingServer(obj *[]model.ReqServer) *[]model.Server {
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

	return &output
}

func HandelSetup(obj *model.ReqSetup) {
	BackendData = *backendSetup(obj, &ServerData)

}

func backendSetup(obj *model.ReqSetup, server *[]model.Server) *model.Backend {
	var tempBackend model.Backend
	tempBackend.Servers = server
	tempBackend.Algo = obj.Algo
	tempBackend.Mode = obj.ModeBalance
	tempBackend.TotalServer = (uint64)(len((*server)))
	tempBackend.TotalServerConnection = 0
	return &tempBackend
}

func GetserverCount() uint64 {
	return serverCount
}

func updateServerCount(mu *sync.Mutex, n uint64) {
	mu.Lock()
	serverCount = (serverCount + 1) % n
	mu.Unlock()
}
