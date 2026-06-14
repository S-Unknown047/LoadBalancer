package natmode

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"

	"syscall"

	helper "github.com/S-Unknown047/LoadBalancer/Helper"
	model "github.com/S-Unknown047/LoadBalancer/Model"
	routingalgo "github.com/S-Unknown047/LoadBalancer/Routing_Algo"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

var b *model.Backend = &helper.BackendData
var mu sync.Mutex
var rawSocketFd int

type conn struct {
	senderIP   string
	senderPort uint16
	dstIP      string
	dstPort    uint16
}

func initRawSocket() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		log.Fatalf("Failed to create raw socket: %v", err)
	}
	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
	if err != nil {
		log.Fatalf("Failed to set IP_HDRINCL: %v", err)
	}
	rawSocketFd = fd
}

func Test() {
	initRawSocket()

	// Specify your network interface (e.g., "eth0", "wlan0", or "lo")
	// busy waiting
	for !helper.Flag {
	}

	device := "lo"

	vipIP := "192.168.1.100"

	translationIP := "127.0.0.2"

	// Berkeley Packet Filter (BPF) is a filter check all the packets at the data link layer and filter out the packets
	// so that only required packet reached the kernel or application
	filter := fmt.Sprintf("(tcp and dst host %s and dst port 80) or (tcp and dst host %s and dst portrange 49152-65535)", vipIP, translationIP)

	// fmt.Println(filter)

	// handle is pcap handel through which we will write back the packet to the network
	// packetChan is channel where we are getting the packet from the StartCapture
	handle, packetChan, err := StartCapture(device, filter)

	defer handle.Close()

	fmt.Println("Get some packet")
	if err != nil {
		log.Fatalf("Failed to start capture: %v", err)
	}

	for packet := range packetChan {

		// here we are having the serializable layer eth, loop
		var serializableLinkLayer gopacket.SerializableLayer
		if eth := packet.Layer(layers.LayerTypeEthernet); eth != nil {
			serializableLinkLayer = eth.(*layers.Ethernet)
		} else if loop := packet.Layer(layers.LayerTypeLoopback); loop != nil {
			serializableLinkLayer = loop.(*layers.Loopback)
		}

		ipv4Layer := packet.Layer(layers.LayerTypeIPv4)
		if ipv4Layer == nil {
			continue
		}

		//  Cast the layer to access IP headers (Source IP, Dest IP)
		ipv4 := ipv4Layer.(*layers.IPv4)

		// 3. Check for the TCP layer
		tcpLayer := packet.Layer(layers.LayerTypeTCP)

		if tcpLayer == nil {
			continue // Not TCP, ignore
		}

		// tye assertion () Cast the layer to access TCP headers (Ports, Flags)
		tcp := tcpLayer.(*layers.TCP)

		fmt.Printf("Captured TCP packet: %s:%d -> %s:%d (Flags: SYN=%t ACK=%t, LinkLayer=%s, Serializable=%t)\n", ipv4.SrcIP, tcp.SrcPort, ipv4.DstIP, tcp.DstPort, tcp.SYN, tcp.ACK, packet.LinkLayer().LayerType().String(), serializableLinkLayer != nil)

		if ipv4.DstIP.String() == vipIP && (tcp.DstPort == 80) {
			fmt.Println("inside client -> server ")

			clientIP := ipv4.SrcIP.String()
			clientPort := uint16(tcp.SrcPort)

			// it is calling the round robin algrithom and getting the backend ip and port
			// check where ip+port is already present or not
			// if present than return same server else return differnt server with round robin

			state := GetOrAssignConnection(clientIP, clientPort, func() (string, uint16, uint16) {
				chosenIP, chosenPortStr := routingalgo.GetServerIpAndPort(b)
				chosenPort, err := strconv.Atoi(chosenPortStr)
				if err != nil {
					log.Printf("Invalid port returned by round robin: %v", err)
				}
				port_ := model.GetPort()
				model.UpdatePort(&mu)
				return chosenIP, uint16(chosenPort), port_
			})

			connection := &conn{
				senderIP:   translationIP,
				senderPort: state.MappedPort,
				dstIP:      state.BackendIP,
				dstPort:    state.BackendPort,
			}

			// fmt.Printf("inside client -> server: forwarding to %s:%d using mapped port %d\n", state.BackendIP, state.BackendPort, state.MappedPort)
			go requestSend(connection, serializableLinkLayer, ipv4, tcp, handle)

		} else if ipv4.DstIP.String() == translationIP && 49152 <= tcp.DstPort {

			clientIP, clientPort, exists := GetClientByMappedPort(uint16(tcp.DstPort))
			portMap := uint16(tcp.DstPort)
			if !exists {
				continue
			}

			fmt.Println("inside server -> client")

			connection := &conn{
				senderIP:   vipIP,
				senderPort: 80,
				dstIP:      clientIP,
				dstPort:    clientPort,
			}

			go requestSend(connection, serializableLinkLayer, ipv4, tcp, handle)

			if tcp.FIN {
				RemoveConnection(clientIP, clientPort, portMap, b)
			}
		}
	}

}

func requestSend(connection *conn, linkLayer gopacket.SerializableLayer, ipv4 *layers.IPv4, tcp *layers.TCP, handle *pcap.Handle) {
	srcIP := connection.senderIP
	srcPort := connection.senderPort
	dstIP := connection.dstIP
	dstPort := connection.dstPort

	ipv4.SrcIP = net.ParseIP(srcIP).To4()
	tcp.SrcPort = layers.TCPPort(srcPort)

	ipv4.DstIP = net.ParseIP(dstIP).To4()
	tcp.DstPort = layers.TCPPort(dstPort)

	err := tcp.SetNetworkLayerForChecksum(ipv4)
	if err != nil {
		log.Printf("Failed to set network layer for checksum: %v", err)
		return
	}

	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	var layersToSerialize []gopacket.SerializableLayer
	// if using the mac address approch pass mac address too in this part
	// first macaddress, ipv4, tcp, payload
	layersToSerialize = append(layersToSerialize, ipv4, tcp, gopacket.Payload(tcp.Payload))

	// this is serialization which we need to do
	// buffer is of gopacket type serializeBuffer used for storing the serialized packet
	buffer := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buffer, options, layersToSerialize...)
	if err != nil {
		log.Printf("Failed to serialize packet: %v", err)
		return
	}

	// // Verify the serialized packet (starts with IPv4)
	// decodedPacket := gopacket.NewPacket(buffer.Bytes(), layers.LayerTypeIPv4, gopacket.Default)
	// decIPv4 := decodedPacket.Layer(layers.LayerTypeIPv4)
	// decTCP := decodedPacket.Layer(layers.LayerTypeTCP)
	// if decIPv4 != nil && decTCP != nil {
	// 	ipL := decIPv4.(*layers.IPv4)
	// 	tcpL := decTCP.(*layers.TCP)
	// 	fmt.Printf("Serialized packet check: %s:%d -> %s:%d (TCP Checksum: 0x%x, IP Checksum: 0x%x)\n", ipL.SrcIP, tcpL.SrcPort, ipL.DstIP, tcpL.DstPort, tcpL.Checksum, ipL.Checksum)
	// } else {
	// 	fmt.Printf("Warning: Failed to decode serialized packet layers! (decIPv4=%t, decTCP=%t)\n", decIPv4 != nil, decTCP != nil)
	// }

	// Send via raw socket
	var addr syscall.SockaddrInet4
	addr.Port = int(dstPort)
	copy(addr.Addr[:], ipv4.DstIP.To4())

	err = syscall.Sendto(rawSocketFd, buffer.Bytes(), 0, &addr)
	if err != nil {
		log.Printf("Failed to send packet via raw socket: %v", err)
	} else {
		fmt.Println("Sucessfull")
	}
}
