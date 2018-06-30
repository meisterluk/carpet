/*
 * Copyright (c) 2018 meisterluk
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"fmt"
	"image/color"
	"os"
	"strings"
)

// osArgOr returns CLI argument at index `i` or string `alt`
func osArgOr(i int, alt string) string {
	if len(os.Args) <= i || os.Args[i] == "" {
		return alt
	}
	if len(os.Args[i]) > 0 && strings.HasPrefix(os.Args[i], `-`) {
		return alt
	}
	return os.Args[i]
}

// col2u32 represents a color as uint32
func col2u32(c color.Color) uint32 {
	r, g, b, a := c.RGBA()
	return ((r >> 8) << 24) | ((g >> 8) << 16) | (b & 0xff00) | (a >> 8)
}

// colorRepr returns a human-readable color representation
func colorRepr(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("%02X%02X%02X%02X", r>>8, g>>8, b>>8, a>>8)
}
