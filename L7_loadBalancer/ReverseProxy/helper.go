package reverseproxy

import (
	"fmt"
	"strconv"
	"sync"

	helper "github.com/S-Unknown047/LoadBalancer/Helper"
	model "github.com/S-Unknown047/LoadBalancer/Model"
	routingalgo "github.com/S-Unknown047/LoadBalancer/Routing_Algo"
)

var b *model.Backend = &helper.BackendData

type connection struct {
	BackendIp   string
	BackendPort uint16
}

var BackendStore = make(map[string]connection)
var mu sync.Mutex

func GetOrAssignConnection(host string, port string) (string, uint16) {
	mu.Lock()
	defer mu.Unlock()
	key := fmt.Sprintf("%s:%s", host, port)

	state, exists := BackendStore[key]

	if exists {
		return state.BackendIp, state.BackendPort
	}

	Ip, Port := routingalgo.GetServerIpAndPort(b)

	portU, err := (strconv.Atoi(Port))

	if err != nil {
		fmt.Println("error while Parsing")
	}

	states := connection{
		BackendIp:   Ip,
		BackendPort: uint16(portU),
	}

	BackendStore[key] = states

	return states.BackendIp, states.BackendPort
}
