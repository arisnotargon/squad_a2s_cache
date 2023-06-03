package funcs

import (
	"fmt"
	"log"

	"github.com/arisnotargon/squad_a2s_cache/config"
	"github.com/davecgh/go-spew/spew"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	REQUEST_DATA *[]byte
	device       *pcap.Interface
	snaplen      = int32(1600)
	promisc      = false
	timeout      = pcap.BlockForever
	filter       = "udp and port 10002"
)

func init() {
	if REQUEST_DATA == nil {
		// b'\xFF\xFF\xFF\xFF\x54Source Engine Query\x00'
		REQUEST_DATA = new([]byte)
		*REQUEST_DATA = make([]byte, 0)
		*REQUEST_DATA = append(*REQUEST_DATA, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x54}...)
		*REQUEST_DATA = append(*REQUEST_DATA, []byte("Source Engine Query")...)
		*REQUEST_DATA = append(*REQUEST_DATA, 0x00)
	}

	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	if device == nil {
		for _, d := range devices {
			if d.Name == config.Conf.DevName {
				d := d
				device = &d
			}
		}
		if device == nil {
			panic(fmt.Errorf("device not found, device name: %s", config.Conf.DevName))
		}
	}
}

func Cap_a2s() {
	// spew.Dump("in Cap_a2s", REQUEST_DATA)
	// spew.Dump(device)
	handle, err := pcap.OpenLive(device.Name, snaplen, promisc, timeout)
	if err != nil {
		spew.Dump("OpenLive err====>", err)
		return
	}
	defer handle.Close()
	if err := handle.SetBPFFilter(filter); err != nil {
		spew.Dump("SetBPFFilter err====>", err)
		return
	}

	source := gopacket.NewPacketSource(handle, handle.LinkType())

	// fmt.Printf("source===> %+v\n", source)
	for packet := range source.Packets() {
		fmt.Println(packet)
		l := packet.TransportLayer()
		// fmt.Printf("packet===> %+v\n", l)
		fmt.Printf("packet payload===> \n")
		spew.Dump(l.LayerPayload())
		// fmt.Printf("packet content===> %+v\n", l.LayerContents())
		// spew.Dump(l.LayerContents())
	}
}
