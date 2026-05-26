package natmode

import (
	"sync"

	model "github.com/S-Unknown047/LoadBalancer/Model"
)

var ipAndPortMap map[uint16]model.IpAndPort
var ipAndPortMapMu sync.Mutex

func StoreIpAndPort(ip string, port uint16) {
	ipAndPortMapMu.Lock()
	ipAndPortMap[port] = model.IpAndPort{IP: ip, Port: port}
	ipAndPortMapMu.Unlock()
}

func GetIpAndPort(port uint16) *model.IpAndPort {
	ipAndPortMapMu.Lock()
	defer ipAndPortMapMu.Unlock()
	val, ok := ipAndPortMap[port]

	if ok != true {
		return nil
	}
	return &val
}
