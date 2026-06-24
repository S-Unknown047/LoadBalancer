package helper

import "sync"

// Round Robin
var serverCount uint64 = 0

func GetserverCount() uint64 {
	return serverCount
}

func UpdateServerCount(mu *sync.Mutex, n uint64) {
	mu.Lock()
	serverCount = (serverCount + 1) % n
	mu.Unlock()
}

// LeastConnection
