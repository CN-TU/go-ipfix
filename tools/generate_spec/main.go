package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"

	"pm.cn.tuwien.ac.at/ipfix/go-ipfix"
)

type specType int

const (
	ieSpec specType = iota
	xmlSpec
)

func main() {
	var input, output *os.File
	var err error
	var spec specType
	if len(os.Args) != 3 {
		log.Panicf("Usage: %s Funcname file\n", os.Args[0])
	}
	funcName := os.Args[1]
	inputName := os.Args[2]
	if r, _ := utf8.DecodeRuneInString(funcName); r == utf8.RuneError || !unicode.IsUpper(r) {
		log.Panicln("Funcname must start with a legal UTF-8 uppercase letter")
	}
	ext := filepath.Ext(inputName)
	switch ext {
	case ".xml":
		spec = xmlSpec
	case ".iespec":
		spec = ieSpec
	default:
		log.Panicf("Unknown file extension '%s'; I only know xml and iespec!\n", ext)
	}
	outputName := inputName[:len(inputName)-len(ext)] + ".go"

	if output, err = os.Create(outputName); err != nil {
		log.Panicln("Couldn't open output file", outputName, err)
	}
	if input, err = os.Open(inputName); err != nil {
		log.Panicln("Couldn't open input file", inputName, err)
	}

	wr := bufio.NewWriter(output)
	fmt.Fprintf(wr, `package ipfix

// GENERATED BY COMMAND ABOVE; DO NOT CHANGE!

func %s() {
`, funcName)
	cb := func(ie ipfix.InformationElement) {
		iev := reflect.ValueOf(ie)
		iet := reflect.TypeOf(ie)
		wr.WriteString("	RegisterInformationElement(InformationElement{")
		for i := 0; i < iev.NumField(); i++ {
			if i > 0 {
				wr.WriteString(", ")
			}
			fmt.Fprintf(wr, "%s: %#v", iet.Field(i).Name, iev.Field(i))
		}
		wr.WriteString("})\n")
	}
	switch spec {
	case xmlSpec:
		xmlspec(input, cb)
	case ieSpec:
		iespec(input, cb)
	}
	wr.WriteString(`}
`)
	wr.Flush()
	input.Close()
	output.Close()
}

func iespec(spec *os.File, cb func(ipfix.InformationElement)) {
	rd := bufio.NewScanner(spec)
	for rd.Scan() {
		cb(ipfix.MakeIEFromSpec(rd.Bytes()))
	}
}

func xmlspec(spec *os.File, cb func(ipfix.InformationElement)) {
	type Record struct {
		XMLName   xml.Name `xml:"record"`
		Name      string   `xml:"name"`
		DataType  []byte   `xml:"dataType"`
		ElementID int      `xml:"elementId"`
	}
	dec := xml.NewDecoder(spec)
SEARCH_IES:
	for {
		if tok, err := dec.Token(); err != nil {
			log.Panic(err)
		} else {
			if tok, ok := tok.(xml.StartElement); ok {
				if tok.Name.Local == "registry" && len(tok.Attr) == 1 && tok.Attr[0].Name.Local == "id" && tok.Attr[0].Value == "ipfix-information-elements" {
					break SEARCH_IES
				}
			}
		}
	}
	var start xml.StartElement
SEARCH_FIRST_RECORD:
	for {
		if tok, err := dec.Token(); err != nil {
			log.Panic(err)
		} else {
			if tok, ok := tok.(xml.StartElement); ok {
				if tok.Name.Local == "record" {
					start = tok
					break SEARCH_FIRST_RECORD
				}
			}
		}
	}
NEXT_RECORD:
	for {
		var rec Record
		dec.DecodeElement(&rec, &start)
		rec.Name = strings.TrimSpace(rec.Name)
		rec.DataType = bytes.TrimSpace(rec.DataType)
		if rec.Name != "" && rec.ElementID != 0 && len(rec.DataType) != 0 &&
			string(rec.DataType) != "basicList" && string(rec.DataType) != "subTemplateList" && string(rec.DataType) != "subTemplateMultiList" {
			//SMELL: hardcoded iana stuff
			cb(ipfix.NewInformationElement(rec.Name, 0, uint16(rec.ElementID), ipfix.NameToType(rec.DataType), 0))
		}
		for {
			if tok, err := dec.Token(); err != nil {
				log.Panic(err)
			} else {
				switch t := tok.(type) {
				case xml.StartElement:
					if t.Name.Local == "record" {
						start = t
						continue NEXT_RECORD
					}
					if t.Name.Local == "registry" {
						goto FINISHED
					}
				case xml.EndElement:
					if t.Name.Local == "registry" {
						goto FINISHED
					}
				}
			}
		}
	}
FINISHED:
}
