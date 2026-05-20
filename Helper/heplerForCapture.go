package helper

import (
	"sync"

	model "github.com/S-Unknown047/LoadBalancer/Model"
)

func Check(ip string, b *model.Backend) *model.Server {
	for _, sIp := range b.Servers {
		if sIp.Url == ip {
			return &sIp
		}
	}
	return nil
}

func UpdateServerCount(mu *sync.Mutex, b *model.Backend) uint64 {
	mu.Lock()
	serverCount = (serverCount + 1) % b.TotalServer
	mu.Unlock()
	return serverCount
}
