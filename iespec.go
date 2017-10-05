package ipfix

import (
	"log"
	"regexp"
	"strconv"
)

// iespec according to draft-trammell-ipfix-text-iespec-01
var iespecRegex = regexp.MustCompile(`(?m)^(.+)\(((\d+)/)?(\d+)\)<(.+)>(\[(.+)\])?$`)

func MakeIEFromSpec(spec []byte) InformationElement {
	var err error
	x := iespecRegex.FindSubmatch(spec)
	if x == nil {
		log.Panicf("Could not parse iespec '%s'!\n", spec)
	}
	name := string(x[1])
	pen := 0
	if x[3] != nil {
		if pen, err = strconv.Atoi(string(x[3])); err != nil {
			log.Panic(err)
		}
	}
	var id int
	if id, err = strconv.Atoi(string(x[4])); err != nil {
		log.Panic(err)
	}
	t := NameToType(x[5])
	length := 0
	if x[7] != nil && x[7][0] != 'v' {
		if length, err = strconv.Atoi(string(x[7])); err != nil {
			log.Panic(err)
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

func RegisterInformationElement(x InformationElement) {
	if _, ok := informationElementRegistry[x.Name]; ok {
		log.Panicf("Information element with name %s already registered\n", x.Name)
	}
	informationElementRegistry[x.Name] = x
}

func GetInformationElement(name string) (ret InformationElement) {
	var ok bool
	if ret, ok = informationElementRegistry[name]; !ok {
		log.Panicf("No information element with name %s registered\n", name)
	}
	return ret
}