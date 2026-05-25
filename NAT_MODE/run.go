package natmode

import (
	"log"
	"net"
	"strconv"

	model "github.com/S-Unknown047/LoadBalancer/Model"
	routingalgo "github.com/S-Unknown047/LoadBalancer/Routing_Algo"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var b *model.Backend

func Test() {

	// Specify your network interface (e.g., "eth0", "wlan0", or "lo")
	device := "lo"

	// BPF Filter: Only capture TCP traffic destined for our Load Balancer's VIP on port 80
	// This drops all unrelated traffic at the kernel level, saving massive CPU cycles.
	vipIP := "192.168.1.100"
	chosenIP, chosenPortStr := routingalgo.RoundRobin(b)
	chosenPort, err := strconv.Atoi(chosenPortStr)

	if err != nil {
		log.Printf("Invalid port returned by round robin: %v", err)
	}

	filter := "tcp and (" + vipIP + "or" + chosenIP + ")"

	handle, packetChan, err := StartCapture(device, filter)

	if err != nil {
		log.Fatalf("Failed to start capture: %v", err)
	}

	defer handle.Close()

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
        
		if ipv4.DstIP.String() == vipIP && tcp.DstPort.String() == "80" {
		log.Printf("Captured Request: %s:%d -> %s:%d (SYN: %t)",
			ipv4.SrcIP, tcp.SrcPort,
			ipv4.DstIP, tcp.DstPort,
			tcp.SYN)

		// Choose IP and Port using Round Robin algorithm

		// Set the new source IP and port (Load Balancer IP and Port)
		ipv4.SrcIP = net.ParseIP(vipIP)
		tcp.SrcPort = layers.TCPPort(model.GetPort()) // Assuming load balancer port is 80

		// Set the new destination IP and port (Chosen Server IP and Port)
		ipv4.DstIP = net.ParseIP(chosenIP)
		tcp.DstPort = layers.TCPPort(chosenPort)

		// Set the network layer for TCP checksum calculation
		err = tcp.SetNetworkLayerForChecksum(ipv4)
		if err != nil {
			log.Printf("Failed to set network layer for checksum: %v", err)
			continue
		}

		// Serialize the modified packet
		options := gopacket.SerializeOptions{
			ComputeChecksums: true,
			FixLengths:       true,
		}
		buffer := gopacket.NewSerializeBuffer()
		err = gopacket.SerializeLayers(buffer, options,
			ipv4,
			tcp,
			gopacket.Payload(tcp.Payload),
		)
		if err != nil {
			log.Printf("Failed to serialize packet: %v", err)
			continue
		}

		// Forward the request to the server
		err = handle.WritePacketData(buffer.Bytes())
		if err != nil {
			log.Printf("Failed to write packet data: %v", err)
		} else {
			log.Printf("Forwarded Request: %s:%d -> %s:%d",
				ipv4.SrcIP, tcp.SrcPort,
				ipv4.DstIP, tcp.DstPort)
		}
	} else if ipv4.DstIP.String() ==  vipIP && <=tcp.DstPort.String() {

	}
	}

}
