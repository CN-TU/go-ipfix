package main

import (
	"log"
	"net"
	"os"
	"time"

	"pm.cn.tuwien.ac.at/ipfix/ipfix"
)

func main() {
	f, err := os.Create("test.ipfix")
	if err != nil {
		log.Panic(err)
	}
	msgStream := ipfix.MakeMessageStream(f, 0, 0)
	id, err := msgStream.AddTemplate(time.Now(),
		ipfix.InformationElement{
			"octetDeltaCount",
			0,
			1,
			ipfix.Unsigned64,
			8,
		},
		ipfix.InformationElement{
			"sourceIPv4Address",
			0,
			8,
			ipfix.Ipv4Address,
			4,
		},
		ipfix.InformationElement{
			"flowEndNanoseconds",
			0,
			157,
			ipfix.DateTimeNanoseconds,
			8,
		},
	)
	if err != nil {
		log.Panic(err)
	}
	if err := msgStream.SendData(time.Now(), id, uint64(5), net.IP{192, 168, 0, 1}, time.Now()); err != nil {
		log.Panic(err)
	}
	if err := msgStream.SendData(time.Now(), id, uint64(10), net.IP{192, 168, 0, 2}, time.Now()); err != nil {
		log.Panic(err)
	}
	if err := msgStream.SendData(time.Now(), id, uint64(2), net.IP{192, 168, 0, 3}, time.Now()); err != nil {
		log.Panic(err)
	}
	if err := msgStream.Finalize(time.Now()); err != nil {
		log.Panic(err)
	}
	f.Close()
}
