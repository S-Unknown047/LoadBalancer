package natmode

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

type ConnectionState struct {
	MappedPort  uint16
	BackendIP   string
	BackendPort uint16
}

var clientToBackendMap = make(map[string]ConnectionState)
var portToClientMap = make(map[uint16]string)
var ipAndPortMapMu sync.Mutex

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
