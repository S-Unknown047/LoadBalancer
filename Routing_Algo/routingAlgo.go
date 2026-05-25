package routingalgo

import (
	"sync"
	"sync/atomic"

	helper "github.com/S-Unknown047/LoadBalancer/Helper"
	model "github.com/S-Unknown047/LoadBalancer/Model"
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
