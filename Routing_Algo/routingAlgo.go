package routingalgo

import (
	"sync"
	"sync/atomic"

	model "github.com/S-Unknown047/LoadBalancer/Model"
)

var serverCount uint64 = 0

var Mu sync.Mutex

func RoundRobin(b *model.Backend) string {

	url := (*b.Servers)[serverCount].IP + ":" + (*b.Servers)[serverCount].Port
	atomic.AddUint64(&((*b.Servers)[serverCount].Connection), 1)
	atomic.AddUint64(&(b.TotalServerConnection), 1)
	// UpdateServerCount(&Mu, b)
	return url
}
