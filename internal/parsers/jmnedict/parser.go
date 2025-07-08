package jmnedict

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
)

type JMnedict struct {
	XMLName xml.Name `xml:"JMnedict"`
	Entries []Entry  `xml:"entry"`
}

type Entry struct {
	Kanji       string `xml:"k_ele>keb"`
	Kana        string `xml:"r_ele>reb"`
	Translation string `xml:"trans>trans_det"`
	NameType    string `xml:"trans>name_type"`
}

func Parse(filepath string) (*JMnedict, error) {
	dict := new(JMnedict)

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	d := xml.NewDecoder(bytes.NewReader(data))
	d.Entity = map[string]string{
		"given":        "(given)",
		"fem":          "(fem)",
		"surname":      "(surname)",
		"company":      "(company)",
		"place":        "(place)",
		"organization": "(organization)",
		"serv":         "(service)",
		"station":      "(station)",
		"work":         "(work)",
		"product":      "(product)",
		"masc":         "(male name)",
		"group":        "(group)",
		"person":       "(person)",
		"unclass":      "(unclass)",
		"char":         "(char)",
		"obj":          "(obj)",
		"dei":          "(dei)",
		"fict":         "(fict)",
		"creat":        "(creat)",
		"myth":         "(myth)",
		"ship":         "(ship)",
		"ev":           "(ev)",
		"leg":          "(leg)",
		"doc":          "(doc)",
	}

	err = d.Decode(dict)
	if err != nil {
		return nil, fmt.Errorf("failed to parse xml: %v", err)
	}

	return dict, nil
}
