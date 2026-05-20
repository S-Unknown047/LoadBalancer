package natmode

import (
	"log"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// StartCapture opens a network interface and returns a channel of raw packets.
func StartCapture(device string, bpfFilter string) (chan gopacket.Packet, error) {
	// 1. Configuration variables

	var snapshotLen int32 = 65535                // Capture the entire packet
	var promiscuous bool = true                  // See all traffic on the wire
	var timeout time.Duration = -1 * time.Second // Non-blocking/immediate delivery

	log.Printf("Opening capture on interface %s...", device)

	// 2. this is to capture all the request on the device (such type of network), max length of capture packet in 65535, promisucous trafice on wire,
	// after which interval i caturing the request
	handle, err := pcap.OpenLive(device, snapshotLen, promiscuous, timeout)
	if err != nil {
		return nil, err
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

	// Create the packet source
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Return the channel so the main router can loop over incoming packets
	return packetSource.Packets(), nil
}
