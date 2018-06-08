package main

import (
	"log"
	"os"
	"time"

	"github.com/CN-TU/go-ipfix"
)

func main() {
	f, err := os.Create("test.ipfix")
	if err != nil {
		log.Panic(err)
	}
	ipfix.LoadIANASpec()
	msgStream := ipfix.MakeMessageStream(f, 0, 0)
	ie := ipfix.NewBasicList("testlist", ipfix.GetInformationElement("octetDeltaCount"), 0)
	id, err := msgStream.AddTemplate(time.Now(),
		ie,
	)
	if err != nil {
		log.Panic(err)
	}
	if err := msgStream.SendData(time.Now(), id, []uint64{1, 2, 3}); err != nil {
		log.Panic(err)
	}
	if err := msgStream.SendData(time.Now(), id, []uint64{4, 5}); err != nil {
		log.Panic(err)
	}
	if err := msgStream.SendData(time.Now(), id, []uint64{10, 20, 33, 100}); err != nil {
		log.Panic(err)
	}
	if err := msgStream.Finalize(time.Now()); err != nil {
		log.Panic(err)
	}
	f.Close()
}
