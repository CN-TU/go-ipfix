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
	ie := ipfix.NewBasicList("testlist", ipfix.GetInformationElement("applicationName"), 0)
	id, err := msgStream.AddTemplate(time.Now(),
		ie,
	)
	if err != nil {
		log.Panic(err)
	}
	if err := msgStream.SendData(time.Now(), id, []string{"testA", "2", "testB"}); err != nil {
		log.Panic(err)
	}
	if err := msgStream.SendData(time.Now(), id, []string{"something longer"}); err != nil {
		log.Panic(err)
	}
	if err := msgStream.SendData(time.Now(), id, []string{"short", "test", "some", "more", "tests"}); err != nil {
		log.Panic(err)
	}
	if err := msgStream.Finalize(time.Now()); err != nil {
		log.Panic(err)
	}
	f.Close()
}
