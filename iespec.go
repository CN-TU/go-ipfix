package ipfix

import (
	"fmt"
	"regexp"
	"strconv"
)

// iespec according to RFC7013 section 10.1
var iespecRegex = regexp.MustCompile(`(?m)^(.+)\(((\d+)/)?(\d+)\)<(.+)>(\[(.+)\])?$`)

// MakeIEFromSpec returns an InformationElement as specified by the provided specification.
// The specification format must follow RFC7013 section 10.1
func MakeIEFromSpec(spec []byte) (InformationElement, error) {
	var err error
	x := iespecRegex.FindSubmatch(spec)
	if x == nil {
		return InformationElement{}, fmt.Errorf("ipfix: Could not parse iespec '%s'", spec)
	}
	name := string(x[1])
	pen := 0
	if x[3] != nil {
		if pen, err = strconv.Atoi(string(x[3])); err != nil {
			return InformationElement{}, fmt.Errorf("ipfix: Could not parse pen '%s'. Must be valid number", x[3])
		}
	}
	var id int
	if id, err = strconv.Atoi(string(x[4])); err != nil {
		return InformationElement{}, fmt.Errorf("ipfix: Could not parse id '%s'. Must be valid number", x[4])
	}
	t := NameToType(x[5])
	if t == IllegalType {
		return InformationElement{}, fmt.Errorf("ipfix: Could not parse id '%s'. Must be valid ipfix type", x[4])
	}
	length := 0
	if x[7] != nil && x[7][0] != 'v' {
		if length, err = strconv.Atoi(string(x[7])); err != nil {
			return InformationElement{}, fmt.Errorf("ipfix: Could not length '%s'. Must be valid number", x[4])
		}
	}
	if length == 0 {
		length = int(DefaultSize[t])
	}
	return NewInformationElement(name, uint32(pen), uint16(id), t, uint16(length)), nil
}

var informationElementRegistry map[string]InformationElement

func init() {
	informationElementRegistry = make(map[string]InformationElement)
}

// RegisterInformationElement registers the given InformationElement. This can later be queried by name with GetInformationElement.
func RegisterInformationElement(x InformationElement) error {
	if _, ok := informationElementRegistry[x.Name]; ok {
		return fmt.Errorf("ipfix: Information element with name '%s' already registered", x.Name)
	}
	informationElementRegistry[x.Name] = x
	return nil
}

// GetInformationElement retrieves an InformationElement by name.
func GetInformationElement(name string) (ret InformationElement, err error) {
	var ok bool
	if ret, ok = informationElementRegistry[name]; !ok {
		err = fmt.Errorf("ipfix: No information element with name '%s' registered", name)
	}
	return
}
