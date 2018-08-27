package ipfix_test

import (
	"bytes"
	"fmt"
	"net"
	"time"

	ipfix "github.com/CN-TU/go-ipfix"
)

func ExampleMakeMessageStream() {
	// output of this example will be in buf
	buf := new(bytes.Buffer)

	// load the iana information elements
	ipfix.LoadIANASpec()

	now := time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC) // simulated fixed time

	// First create a message stream; mtu=0 chooses the default size
	msgStream := ipfix.MakeMessageStream(buf, 0, 0)

	// Add a new template with three information elements
	id, err := msgStream.AddTemplate(now,
		ipfix.GetInformationElement("octetDeltaCount"),
		ipfix.GetInformationElement("sourceIPv4Address"),
		ipfix.GetInformationElement("flowEndNanoseconds"),
	)
	if err != nil {
		fmt.Println("MessageStream.AddTemplate failed:", err)
		return
	}

	// Export data for this information Element
	if err := msgStream.SendData(now, id, uint64(5), net.IP{192, 168, 0, 1}, now); err != nil {
		fmt.Println("MessageStream.SendData failed:", err)
		return
	}

	now = now.Add(1 * time.Second)
	if err := msgStream.SendData(now, id, uint64(10), net.IP{192, 168, 0, 2}, now); err != nil {
		fmt.Println("MessageStream.SendData failed:", err)
		return
	}

	now = now.Add(1 * time.Minute)
	if err := msgStream.SendData(now.Add(10*time.Second), id, uint64(2), net.IP{192, 168, 0, 3}, now); err != nil {
		fmt.Println("MessageStream.SendData failed:", err)
		return
	}

	// Call finalize
	if err := msgStream.Finalize(now); err != nil {
		fmt.Println("MessageStream.Finalize failed:", err)
		return
	}

	// buf holds now the complete ipfix data of this example
	fmt.Printf("% x", buf.Bytes())
	// Output: 00 0a 00 64 5a 49 7a 3d 00 00 00 00 00 00 00 00 00 02 00 14 01 00 00 03 00 01 00 08 00 08 00 04 00 9d 00 08 01 00 00 40 00 00 00 00 00 00 00 05 c0 a8 00 01 dd f3 f8 80 00 00 00 00 00 00 00 00 00 00 00 0a c0 a8 00 02 dd f3 f8 81 00 00 00 00 00 00 00 00 00 00 00 02 c0 a8 00 03 dd f3 f8 bd 00 00 00 00
}

func ExampleMakeMessageStream_basicList() {
	// output of this example will be in buf
	buf := new(bytes.Buffer)

	// load the iana information elements
	ipfix.LoadIANASpec()

	now := time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC) // simulated fixed time

	// First create a message stream; mtu=0 chooses the default size
	msgStream := ipfix.MakeMessageStream(buf, 0, 0)

	// Create an information elemenet for the basiclist, holding a variable number of octetDeltaCount
	ie := ipfix.NewBasicList("testlist", ipfix.GetInformationElement("octetDeltaCount"), 0)

	// Add a new template with this information element
	id, err := msgStream.AddTemplate(now,
		ie,
	)
	if err != nil {
		fmt.Println("MessageStream.AddTemplate failed:", err)
		return
	}

	// write out some data
	// Note that despite octetDeltaCount being defined as uint64, you can use differnt types here
	// The data will be converted to the right type; There is no automation to export data with shorter
	// values
	if err := msgStream.SendData(now, id, []uint64{1, 2, 3}); err != nil {
		fmt.Println("MessageStream.SendData failed:", err)
	}
	now = now.Add(1 * time.Second)
	if err := msgStream.SendData(now, id, []uint8{4, 5}); err != nil {
		fmt.Println("MessageStream.SendData failed:", err)
	}
	now = now.Add(1 * time.Second)
	if err := msgStream.SendData(now, id, []int32{10, 20, 33, 100}); err != nil {
		fmt.Println("MessageStream.SendData failed:", err)
	}
	now = now.Add(1 * time.Second)
	if err := msgStream.Finalize(now); err != nil {
		fmt.Println("MessageStream.Finalize failed:", err)
	}

	// buf holds now the complete ipfix data of this example
	fmt.Printf("% x", buf.Bytes())
	// Output: 00 0a 00 80 5a 49 7a 03 00 00 00 00 00 00 00 00 00 02 00 0c 01 00 00 01 01 23 ff ff 01 00 00 64 ff 00 1d ff 00 01 00 08 00 00 00 00 00 00 00 01 00 00 00 00 00 00 00 02 00 00 00 00 00 00 00 03 ff 00 15 ff 00 01 00 08 00 00 00 00 00 00 00 04 00 00 00 00 00 00 00 05 ff 00 25 ff 00 01 00 08 00 00 00 00 00 00 00 0a 00 00 00 00 00 00 00 14 00 00 00 00 00 00 00 21 00 00 00 00 00 00 00 64
}

func ExampleMakeMessageStream_basicListVariable() {
	// output of this example will be in buf
	buf := new(bytes.Buffer)

	// load the iana information elements
	ipfix.LoadIANASpec()

	now := time.Date(2018, 01, 01, 0, 0, 0, 0, time.UTC) // simulated fixed time

	// First create a message stream; mtu=0 chooses the default size
	msgStream := ipfix.MakeMessageStream(buf, 0, 0)

	// Create an information elemenet for the basiclist, holding a variable number of applicationName, which in turn are also variable length
	ie := ipfix.NewBasicList("testlist", ipfix.GetInformationElement("applicationName"), 0)

	// Add a new template with this information element
	id, err := msgStream.AddTemplate(now,
		ie,
	)
	if err != nil {
		fmt.Println("MessageStream.AddTemplate failed:", err)
		return
	}

	// write out some data
	if err := msgStream.SendData(now, id, []string{"testA", "2", "testB"}); err != nil {
		fmt.Println("MessageStream.SendData failed:", err)
	}
	now = now.Add(1 * time.Second)
	if err := msgStream.SendData(now, id, []string{"something longer"}); err != nil {
		fmt.Println("MessageStream.SendData failed:", err)
	}
	now = now.Add(1 * time.Second)
	if err := msgStream.SendData(now, id, []string{"short", "test", "some", "more", "tests"}); err != nil {
		fmt.Println("MessageStream.SendData failed:", err)
	}
	now = now.Add(1 * time.Second)
	if err := msgStream.Finalize(now); err != nil {
		fmt.Println("MessageStream.Finalize failed:", err)
	}

	// buf holds now the complete ipfix data of this example
	fmt.Printf("% x", buf.Bytes())
	// Output: 00 0a 00 72 5a 49 7a 03 00 00 00 00 00 00 00 00 00 02 00 0c 01 00 00 01 01 23 ff ff 01 00 00 56 ff 00 13 ff 00 60 ff ff 05 74 65 73 74 41 01 32 05 74 65 73 74 42 ff 00 16 ff 00 60 ff ff 10 73 6f 6d 65 74 68 69 6e 67 20 6c 6f 6e 67 65 72 ff 00 20 ff 00 60 ff ff 05 73 68 6f 72 74 04 74 65 73 74 04 73 6f 6d 65 04 6d 6f 72 65 05 74 65 73 74 73
}
