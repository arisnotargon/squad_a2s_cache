package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/gopacket/pcap"
)

func main() {
	devList, err := pcap.FindAllDevs()
	if err != nil {
		spew.Dump("FindAllDevs err==>", err)
	}
	for idx, device := range devList {
		fmt.Printf("device %d ===> [%+v]\n", idx, device)
	}
}
