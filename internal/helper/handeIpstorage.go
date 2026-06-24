package helper

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	model "github.com/S-Unknown047/LoadBalancer/internal/model"
)

type ConnectionState struct {
	MappedPort  uint16
	BackendIP   string
	BackendPort uint16
}

var clientToBackendMap = make(map[string]ConnectionState)
var portToClientMap = make(map[uint16]string)
var ipAndPortMapMu sync.Mutex

func RemoveConnection(clientIP string, clientPort uint16, portMap uint16, b *model.Backend) {
	ipAndPortMapMu.Lock()
	key := fmt.Sprintf("%s:%d", clientIP, clientPort)
	state, exists := clientToBackendMap[key]
	if !exists {
		ipAndPortMapMu.Unlock()
		return
	}
	delete(clientToBackendMap, key)
	delete(portToClientMap, portMap)
	ipAndPortMapMu.Unlock()

	for i := range *b.Servers {
		srv := &(*b.Servers)[i]
		temp, _ := strconv.Atoi(srv.Port)
		port := uint16(temp)

		if srv.IP == state.BackendIP && port == state.BackendPort {
			atomic.AddUint64(&srv.Connection, ^uint64(0)) // Atomic decrement by 1
			HeapFix(srv)
			atomic.AddUint64(&b.TotalServerConnection, ^uint64(0)) // Atomic decrement by 1
			break
		}
	}
}

func GetOrAssignConnection(clientIP string, clientPort uint16, assignBackend func() (string, uint16, uint16)) ConnectionState {
	key := fmt.Sprintf("%s:%d", clientIP, clientPort)
	ipAndPortMapMu.Lock()
	defer ipAndPortMapMu.Unlock()

	if state, exists := clientToBackendMap[key]; exists {
		return state
	}

	backendIP, backendPort, mappedPort := assignBackend()

	state := ConnectionState{
		MappedPort:  mappedPort,
		BackendIP:   backendIP,
		BackendPort: backendPort,
	}

	clientToBackendMap[key] = state
	portToClientMap[mappedPort] = key

	return state
}

func GetClientByMappedPort(mappedPort uint16) (string, uint16, bool) {
	ipAndPortMapMu.Lock()
	defer ipAndPortMapMu.Unlock()

	key, exists := portToClientMap[mappedPort]
	if !exists {
		return "", 0, false
	}

	var ip string
	var port uint16
	parts := strings.Split(key, ":")
	if len(parts) == 2 {
		ip = parts[0]
		portV, err := strconv.Atoi(parts[1])
		if err != nil {
			log.Printf("Invalid port returned by round robin: %v", err)
			return "", 0, false
		}
		port = uint16(portV)

	}
	return ip, port, true
}
