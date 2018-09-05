package ipfix_test

import (
	"fmt"

	ipfix "github.com/CN-TU/go-ipfix"
)

func ExampleInformationElement_Reverse() {
	ipfix.LoadIANASpec()

	ie, err := ipfix.GetInformationElement("octetDeltaCount")
	if err != nil {
		fmt.Println("GetInformationElement failed:", err)
		return
	}
	fmt.Println(ie)

	revie := ie.Reverse()
	fmt.Println(revie)

	revrevie := revie.Reverse()
	fmt.Println(revrevie)
	// Output:
	// octetDeltaCount
	// reverseOctetDeltaCount(29305/1)<unsigned64>
	// octetDeltaCount
}
