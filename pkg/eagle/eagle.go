// Copyright 2019 John Slee <jslee@jslee.io>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package eagle

import (
	"encoding/xml"
	"io/ioutil"
	"log"
	"os"
)

// LoadEagleFile attempts to read and unmarshal an Eagle XML file.
func LoadEagleFile(filename string) (*Eagle, error) {
	xmlText, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var eagle Eagle
	if err := xml.Unmarshal(xmlText, &eagle); err != nil {
		return nil, err
	}
	return &eagle, nil
}

// WriteFile attempts to generate a valid Eagle XML board file from an
// Eagle data structure.
func (e *Eagle) WriteFile(filename string) error {
	wrapped := &struct {
		Eagle
		XMLName struct{} `xml:"eagle"`
	}{Eagle: *e}
	xml, err := xml.MarshalIndent(wrapped, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0644)
	defer file.Close()
	header := "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n"
	header += "<!DOCTYPE eagle SYSTEM \"eagle.dtd\">\n"
	if _, err := file.WriteString(header); err != nil {
		return err
	}
	if _, err := file.Write(xml); err != nil {
		return err
	}
	return nil
}

// CloneEmpty creates an empty shell based on an existing Eagle object, copying
// layer and grid definitions, and the Eagle version, but not copying any board
// elements, board outline, DRC, signals or libraries.
func (e *Eagle) CloneEmpty() *Eagle {
	clone := &Eagle{
		Version: e.Version,
		Grid:    e.Grid,
		Layers:  []Layer{},
		Board:   NewBoard(),
	}
	for _, layer := range e.Layers {
		clone.Layers = append(clone.Layers, layer)
	}
	return clone
}

// LayerByName attempts to find the layer number for a named layer. Eagle does
// appear to standardise these but it's easy to do a lookup, so let's be
// tolerant of future surprises. Aborts if the desired layer is not present,
// as it's extremely unlikely that any other action is desired in that case.
func (e *Eagle) LayerByName(name string) int {
	for _, layer := range e.Layers {
		if layer.Name == name {
			return layer.Number
		}
	}
	log.Fatalf("Requested layer %q not found!", name)
	return 0
}
