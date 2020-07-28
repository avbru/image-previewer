//The MIT License (MIT)
//
//Copyright (c) 2020 Aleksey Bakin
//
//Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

//nolint:testpackage
package router

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShiftPath(t *testing.T) {
	tests := []struct {
		path string
		head string
		tail string
	}{
		{"", "", "/"},
		{"/", "", "/"},
		{"/path", "path", "/"},
		{"/path/", "path", "/"},
		{"/path/path2", "path", "/path2"},
		{"/path/path2/", "path", "/path2"},
	}

	//nolint:scopelint
	for _, tc := range tests {
		t.Run(tc.path, func(t *testing.T) {
			head, tail := shiftPath(tc.path)
			assert.Equal(t, tc.head, head)
			assert.Equal(t, tc.tail, tail)
		})
	}
}
