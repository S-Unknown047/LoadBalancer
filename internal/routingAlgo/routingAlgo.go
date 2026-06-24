package routingalgo

import (
	"log"
	"sync"
	"sync/atomic"

	helper "github.com/S-Unknown047/LoadBalancer/internal/helper"
	model "github.com/S-Unknown047/LoadBalancer/internal/model"
)

var Mu sync.Mutex

func RoundRobin(b *model.Backend) (ip string, port string) {
	var serverCount uint64 = helper.GetserverCount()

	ip = (*b.Servers)[serverCount].IP
	port = (*b.Servers)[serverCount].Port
	helper.UpdateServerCount(&Mu, b.TotalServer)
	atomic.AddUint64(&(b.TotalServerConnection), 1)
	return ip, port
}

func LeastConnectionServer(b *model.Backend) (string, string) {
	server := helper.GetServer()
	if server == nil {
		return "", ""
	}
	atomic.AddUint64(&(b.TotalServerConnection), 1)
	return server.IP, server.Port
}

func GetServerIpAndPort(b *model.Backend) (string, string) {

	algo := b.Algo
	var ip string
	var port string
	switch algo {
	case "RoundRobin", "rr":
		ip, port = RoundRobin(b)
		break
	case "LeastConnection", "lc":
		ip, port = LeastConnectionServer(b)
		break
	default:
		log.Fatal("Not a right server")
		return "", ""
	}
	return ip, port

}
