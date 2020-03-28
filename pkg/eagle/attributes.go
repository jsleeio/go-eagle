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
	"fmt"
	"strconv"
	"strings"
)

// AttributeCarrier allows applying attribute operations to any objects that
// contain attributes
type AttributeCarrier interface {
	GetAttributes() []Attribute
}

// AttributeString returns an attribute's string value, or if it isn't found,
// a provided default value. Whitespace is not trimmed as that would hinder
// supplying an attribute value of just whitespace (unlikely?)
func AttributeString(c AttributeCarrier, name string, def string) string {
	for _, attribute := range c.GetAttributes() {
		if attribute.Name == name {
			return attribute.Value
		}
	}
	return def
}

// AttributeFloat returns an attribute's numeric value as a float64, or if
// it isn't found, a provided default value.
func AttributeFloat(c AttributeCarrier, name string, def float64) (float64, error) {
	s := AttributeString(c, name, fmt.Sprint(def))
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, fmt.Errorf("unparseable numeric floating-point attribute value for %s: %v", name, err)
	}
	return f, nil
}

// AttributeBool returns an attribute's value as a boolean value, or if
// it isn't found, a provided default value. Valid values are "yes"/"true"
// or "no"/"false". Case insensitive.
func AttributeBool(c AttributeCarrier, name string, def bool) (bool, error) {
	s := strings.ToLower(strings.TrimSpace(AttributeString(c, name, fmt.Sprint(def))))
	switch {
	case s == "yes" || s == "true":
		return true, nil
	case s == "no" || s == "false":
		return false, nil
	default:
		return false, fmt.Errorf("unparseable boolean attribute value for %s: %q", name, s)
	}
}

// AttributeInt returns an attribute's numeric value as an int, or if
// it isn't found, a provided default value.
func AttributeInt(c AttributeCarrier, name string, def int) (int, error) {
	s := AttributeString(c, name, fmt.Sprint(def))
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("unparseable numeric integer attribute value for %s: %v", name, err)
	}
	return n, nil
}
