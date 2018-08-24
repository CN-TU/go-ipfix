package ipfix

import (
	"fmt"
	"regexp"
	"strconv"
)

// iespec according to draft-trammell-ipfix-text-iespec-01
var iespecRegex = regexp.MustCompile(`(?m)^(.+)\(((\d+)/)?(\d+)\)<(.+)>(\[(.+)\])?$`)

// MakeIEFromSpec returns an InformationElement as specified by the provided specification.
// The specification format must follow draft-trammell-ipfix-text-iespec-01
func MakeIEFromSpec(spec []byte) InformationElement {
	var err error
	x := iespecRegex.FindSubmatch(spec)
	if x == nil {
		panic(fmt.Sprintf("Could not parse iespec '%s'!\n", spec))
	}
	name := string(x[1])
	pen := 0
	if x[3] != nil {
		if pen, err = strconv.Atoi(string(x[3])); err != nil {
			panic(err)
		}
	}
	var id int
	if id, err = strconv.Atoi(string(x[4])); err != nil {
		panic(err)
	}
	t := NameToType(x[5])
	length := 0
	if x[7] != nil && x[7][0] != 'v' {
		if length, err = strconv.Atoi(string(x[7])); err != nil {
			panic(err)
		}
	}
	if length == 0 {
		length = int(DefaultSize[t])
	}
	return NewInformationElement(name, uint32(pen), uint16(id), t, uint16(length))
}

var informationElementRegistry map[string]InformationElement

func init() {
	informationElementRegistry = make(map[string]InformationElement)
}

// RegisterInformationElement registers the given InformationElement. This can later be queried by name with GetInformationElement.
func RegisterInformationElement(x InformationElement) {
	if _, ok := informationElementRegistry[x.Name]; ok {
		panic(fmt.Sprintf("Information element with name %s already registered\n", x.Name))
	}
	informationElementRegistry[x.Name] = x
}

// GetInformationElement retrieves an InformationElement by name.
func GetInformationElement(name string) (ret InformationElement) {
	var ok bool
	if ret, ok = informationElementRegistry[name]; !ok {
		panic(fmt.Sprintf("No information element with name %s registered\n", name))
	}
	return ret
}
