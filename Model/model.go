package model

import "time"

type Backend struct {
	Servers               []Server
	TotalServer           uint64
	TotalServerConnection uint64
}

type Server struct {
	Url           string
	Connection    uint64
	Weight        uint64
	Response_time time.Duration
}

type ReqServer struct {
	IP     string
	Weight uint64
}
