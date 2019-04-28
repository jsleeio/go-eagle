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

// BoardOutlineWires generates a rectangular outline out of wires, suitable
// for placing in Eagle's Dimension layer to create a board outline
func BoardOutlineWires(width, height float64, layer int) []Wire {
	return []Wire{
		// origin at bottom left.
		{X1: 0, Y1: height, X2: width, Y2: height, Width: 0, Layer: layer}, // top
		{X1: 0, Y1: 0, X2: width, Y2: 0, Width: 0, Layer: layer},           // bottom
		{X1: 0, Y1: 0, X2: 0, Y2: height, Width: 0, Layer: layer},          // left
		{X1: width, Y1: 0, X2: width, Y2: height, Width: 0, Layer: layer},  // right
	}
}
