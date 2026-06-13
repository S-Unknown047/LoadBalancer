package serverTesting

import (
	"fmt"
	"net"
	"time"

	model "github.com/S-Unknown047/LoadBalancer/Model"
)

var allServer *([]model.Server)

func SerupAliveServer(ser *[]model.Server) {
	allServer = ser
}

func testConection(port string, ip string) bool {
	address := net.JoinHostPort(ip, port)
	timeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", address, timeout)
	defer conn.Close()
	if err != nil {
		fmt.Println("not reachable")
		return false
	}
	return true
}

func ServerTest() *([]model.Server) {
	var ser []model.Server
	for _, server := range *allServer {

		ip := server.IP
		port := server.Port
		if testConection(port, ip) {
			ser = append(ser, server)
		}
	}

	return &ser
}
