// Package parser reads a JSON config file and scans values into a struct defined externally,
// all fields in the external struct must be exported (capitalised)
package parser

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

// Parser contains the data from the config file
type Parser struct {
	data []byte
}

// NewParser reads the file passed in to create a parser object
func NewParser(file string) (*Parser, error) {
	p := new(Parser)

	var err error
	p.data, err = ioutil.ReadFile(file)
	if err != nil {
		return p, err
	}

	return p, nil
}

// Scan decodes the file data into a struct
func (p *Parser) Scan(o interface{}) error {
	return json.NewDecoder(bytes.NewReader(p.data)).Decode(&o)
}
