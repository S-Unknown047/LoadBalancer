when you restart your system yout virtual ip will be gone 
sudo ip addr add 192.168.1.100/24 dev wlp0s20f3
execute this command to add virtual ip to your system
<!-- 
sudo iptables -A INPUT -p tcp -d 192.168.1.100 --dport 80 -j DROP
# This stops the kernel from sending RST to your backend servers' responses
sudo iptables -A INPUT -p tcp -d 192.168.1.100 --dport 49152:65535 -j DROP -->


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

what are the x-forward are and why they are used in a proxy server request

The SetXForwarded function fixes this by having the proxy inject the original client's information into the headers before forwarding the request. Here is exactly what each header does:
1. X-Forwarded-For (XFF)

What it does: Identifies the originating IP address of the client.
Why it matters: Without this, your backend server would log the load balancer's IP address for every single request. You need X-Forwarded-For for rate limiting, geo-location, analytics, and security auditing (like banning malicious IP addresses).

    Example: X-Forwarded-For: 203.0.113.195 (If it passes through multiple proxies, it becomes a comma-separated list: client-ip, proxy1-ip, proxy2-ip).

2. X-Forwarded-Host (XFH)

What it does: Identifies the original domain name the client requested (the original Host header).
Why it matters: Proxies often rewrite the Host header to route traffic internally (e.g., changing www.myapp.com to internal-server-01.local). If your backend application needs to generate absolute URLs (like sending a password reset email with a link), it needs to know the original domain the user was looking at.

    Example: X-Forwarded-Host: www.myapp.com

3. X-Forwarded-Proto (XFP)

What it does: Identifies the original protocol the client used to connect (http or https).
Why it matters: As mentioned, load balancers usually handle the heavy lifting of SSL/TLS decryption and forward the traffic to your backend over plain HTTP. If your backend app doesn't know the user originally connected securely, it might accidentally trigger an infinite redirect loop (constantly trying to redirect the user to HTTPS, not realizing they are already there).

    Example: X-Forwarded-Proto: https

X-Forwarded-For: 198.51.100.1
X-Forwarded-Host: example.com
X-Forwarded-Proto: https