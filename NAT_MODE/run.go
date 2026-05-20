package natmode

import (
	"log"

	helper "github.com/S-Unknown047/LoadBalancer/Helper"
	model "github.com/S-Unknown047/LoadBalancer/Model"
	"github.com/google/gopacket/layers"
)

var b *model.Backend

func Test() {

	// Specify your network interface (e.g., "eth0", "wlan0", or "lo")

	device := "lo"

	// BPF Filter: Only capture TCP traffic destined for our Load Balancer's VIP on port 80
	// This drops all unrelated traffic at the kernel level, saving massive CPU cycles.
	vipIP := "192.168.1.100"
	filter := "tcp and dst host " + vipIP + " and dst port 80"

	packetChan, err := StartCapture(device, filter)

	if err != nil {
		log.Fatalf("Failed to start capture: %v", err)
	}

	log.Println("Listening for raw packets...")

	// The Main Routing Loop
	for packet := range packetChan {
		// 1. Check if the packet has an IPv4 layer
		ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ipv4Layer == nil {
			continue // Not IPv4, ignore
		}

		// 2. Cast the layer to access IP headers (Source IP, Dest IP)
		ipv4 := ipv4Layer.(*layers.IPv4)

		// 3. Check for the TCP layer
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue // Not TCP, ignore
		}

		// 4. Cast the layer to access TCP headers (Ports, Flags)
		tcp := tcpLayer.(*layers.TCP)
		if string(ipv4.DstIP) == "192.168.1.102" {
			serv := helper.Check(string(ipv4.SrcIP), b)
			if serv != nil {

			}
		}
		log.Printf("Captured Request: %s:%d -> %s:%d (SYN: %t)",
			ipv4.SrcIP, tcp.SrcPort,
			ipv4.DstIP, tcp.DstPort,
			tcp.SYN)
	}
}
