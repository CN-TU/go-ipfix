package main

import (
	"log"
	"net"
	"os"
	"time"

	"pm.cn.tuwien.ac.at/ipfix/ipfix"
	_ "pm.cn.tuwien.ac.at/ipfix/ipfix/specs/iana"
)

func main() {
	f, err := os.Create("test.ipfix")
	if err != nil {
		log.Panic(err)
	}
	msgStream := ipfix.MakeMessageStream(f, 0, 0)
	id, err := msgStream.AddTemplate(time.Now(),
		ipfix.GetInformationElement("octetDeltaCount"),
		ipfix.GetInformationElement("sourceIPv4Address"),
		ipfix.GetInformationElement("flowEndNanoseconds"),
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
