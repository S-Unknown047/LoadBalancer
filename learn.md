when you restart your system yout virtual ip will be gone 
sudo ip addr add 192.168.1.100/24 dev wlp0s20f3
execute this command to add virtual ip to your system

sudo iptables -A INPUT -p tcp -d 192.168.1.100 --dport 80 -j DROP
# This stops the kernel from sending RST to your backend servers' responses
sudo iptables -A INPUT -p tcp -d 192.168.1.100 --dport 49152:65535 -j DROP

sudo iptables -A INPUT -p tcp -d 127.0.0.2 --dport 49152:65535 -j DROP

sudo iptables -A INPUT -p tcp -d 127.0.0.2 --dport 49152:65535 -j DROP


Ticker usage to make sure some code execute at every x seconds
https://leapcell.medium.com/working-with-scheduled-tasks-in-go-timer-and-ticker-5b6c4289a63c


# this code for when we have device on different device

func requestSend(
	connection *conn, 
	ipv4 *layers.IPv4, 
	tcp *layers.TCP, 
	handle *pcap.Handle,
	srcMAC net.HardwareAddr, // Load Balancer's MAC address
	dstMAC net.HardwareAddr, // Destination Server's (or Gateway's) MAC address
) {
	srcIP := connection.senderIP
	srcPort := connection.senderPort
	dstIP := connection.dstIP
	dstPort := connection.dstPort

	// 1. Rewrite IP and TCP Headers
	ipv4.SrcIP = net.ParseIP(srcIP).To4()
	tcp.SrcPort = layers.TCPPort(srcPort)

	ipv4.DstIP = net.ParseIP(dstIP).To4()
	tcp.DstPort = layers.TCPPort(dstPort)

	// 2. Set network layer for checksum calculations
	err := tcp.SetNetworkLayerForChecksum(ipv4)
	if err != nil {
		log.Printf("Failed to set network layer for checksum: %v", err)
		return
	}

	// 3. Construct the manual Ethernet header
	ethernetLayer := &layers.Ethernet{
		SrcMAC:       srcMAC,
		DstMAC:       dstMAC,
		EthernetType: layers.EthernetTypeIPv4,
	}

	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}

	// 4. Group layers (including the link layer) for serialization
	var layersToSerialize []gopacket.SerializableLayer
	layersToSerialize = append(layersToSerialize, ethernetLayer, ipv4, tcp, gopacket.Payload(tcp.Payload))

	buffer := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buffer, options, layersToSerialize...)
	if err != nil {
		log.Printf("Failed to serialize packet: %v", err)
		return
	}

	// 5. Write raw Ethernet frame directly to the network interface via pcap
	err = handle.WritePacketData(buffer.Bytes())
	if err != nil {
		log.Printf("Failed to write packet data: %v", err)
	} else {
		fmt.Println("Successfully sent Layer 2 frame")
	}
}
