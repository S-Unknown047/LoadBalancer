package model

import "sync"

type Backend struct {
	Servers               *([]Server)
	TotalServer           uint64
	TotalServerConnection uint64
	Algo                  string
	Mode                  string
}

type Server struct {
	IP         string
	Port       string
	Connection uint64
	Weight     uint64
	mu         *sync.Mutex
	// Response_time time.Duration
}

type ReqServer struct {
	IP     string
	Port   string
	Weight uint64
}

type ReqSetup struct {
	Algo        string
	ModeBalance string // what is the balance to use  l4 or l7
}

func (b *Backend) Check(ip string) *Server {
	if b == nil || b.Servers == nil {
		return nil
	}
	for i := range *b.Servers {
		if (*b.Servers)[i].IP == ip {
			return &(*b.Servers)[i]
		}
	}
	return nil
}
