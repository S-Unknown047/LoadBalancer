package proxy

import (
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// StartCapture opens a network interface and returns a handle and a channel of raw packets.
func StartCapture(device string, bpfFilter string) (*pcap.Handle, chan gopacket.Packet, error) {
	// 1. Configuration variables

	var snapshotLen int32 = 65535                // Capture the entire packet
	var promiscuous bool = true                  // See all traffic on the wire
	var timeout time.Duration = -1 * time.Second // Non-blocking/immediate delivery

	log.Printf("Opening capture on interface %s...", device)

	// 2. this is to capture all the request on the device (such type of network), max length of capture packet in 65535, promisucous trafice on wire,
	// after which interval i caturing the request
	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		return nil, nil, err
	}

	// 3. Apply the BPF (Berkeley Packet Filter)\
	// this is BPF which is filter := "tcp port 80"
	// This is CRITICAL for performance. We only want load balancer traffic, not SSH or background noise.
	if bpfFilter != "" {
		log.Printf("Applying BPF Filter: %s", bpfFilter)
		err = handle.SetBPFFilter(bpfFilter)
		if err != nil {
			log.Fatalf("Error setting BPF filter: %v", err)
		}
	}

	// this is packet source create it take handle which is pcap.handle and provide packet, and handle.LinkType() tell about link layer type
	// it return a sturct *gopacket.PacketSource through which we can read packets and feed into
	// correct decoder based of linkType

	// This is documentation of packetSource  we will use the packetSource.Packets() method which returs a channel of packets
	// PacketSource reads in packets from a PacketDataSource, decodes them, and
	// returns them.
	//
	// There are currently two different methods for reading packets in through
	// a PacketSource:
	//
	// Reading With Packets Function
	//
	// This method is the most convenient and easiest to code, but lacks
	// flexibility.  Packets returns a 'chan Packet', then asynchronously writes
	// packets into that channel.  Packets uses a blocking channel, and closes
	// it if an io.EOF is returned by the underlying PacketDataSource.  All other
	// PacketDataSource errors are ignored and discarded.
	//  for packet := range packetSource.Packets() {
	//    ...
	//  }
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Return the handle and the channel so the main router can loop over incoming packets
	return handle, packetSource.Packets(), nil
}
