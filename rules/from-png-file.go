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

package rules

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
)

// RuleFromPNGFile reads rules from SVG files
type RuleFromPNGFile struct {
	size   int
	match  color.Color
	colors []color.Color
}

// NewRuleFromPNGFile creates a new Rule
func NewRuleFromPNGFile(path string) (*RuleFromPNGFile, error) {
	rgbQuadruple := regexp.MustCompile(`([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})([0-9a-fA-F]{2})`)

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	// parse match color from filename
	filename := filepath.Base(path)
	matches := rgbQuadruple.FindStringSubmatch(filename)
	if len(matches) == 0 {
		return nil, fmt.Errorf(`expected RGBA quadruple in hex like '204060FF' in rule PNG filename; got '%s'`, filename)
	}
	var cols [4]uint8
	for i := 1; i <= 4; i++ {
		b, err := hex.DecodeString(matches[i])
		if err != nil {
			return nil, err
		}
		cols[i-1] = uint8(b[0])
	}

	// decode image
	img, err := png.Decode(fd)
	if err != nil {
		return nil, err
	}
	b := img.Bounds()

	// check dimensions
	height := b.Max.Y - b.Min.Y
	width := b.Max.X - b.Min.X
	if height != width {
		return nil, fmt.Errorf(`expected rule image to be square; got dimensions %d×%d with '%s'`, width, height, filename)
	}

	// create Rule instance
	rule := new(RuleFromPNGFile)
	rule.size = width
	rule.match = color.RGBA{cols[0], cols[1], cols[2], cols[3]}
	rule.colors = make([]color.Color, 0, 16)

	// read colors
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			rule.colors = append(rule.colors, img.At(x, y))
		}
	}

	return rule, nil
}

// Size returns the width (= height) of the replaced colors
func (r *RuleFromPNGFile) Size() int {
	return r.size
}

// Matches returns whether this rule is responsible for the given color (Chain-of-responsibility pattern)
func (r *RuleFromPNGFile) Matches(c color.Color) bool {
	r1, g1, b1, a1 := r.match.RGBA()
	r2, g2, b2, a2 := c.RGBA()
	return r1 == r2 && g1 == g2 && b1 == b2 && a1 == a2
}

// Colorize returns the replacement color at a given position
// where x and y are between 0 and Size() inclusively
func (r *RuleFromPNGFile) Colorize(x, y int) color.Color {
	if 0 <= x && x < r.size && 0 <= y && y < r.size {
		return r.colors[y*r.size+x]
	} else {
		return nil
	}
}

// String gives a human-readable representation of this Rule
func (r *RuleFromPNGFile) String() string {
	runes := `abcdefghijklmnopqrstuvwxyz`
	assoc := make(map[uint32]rune)
	i := 0

	out := "{"
	for y := 0; y < r.size; y++ {
		for x := 0; x < r.size; x++ {
			col := col2u32(r.colors[y*r.size+x])
			r, ok := assoc[col]
			if ok {
				out += string(r)
			} else {
				assoc[col] = rune(runes[i])
				out += string(assoc[col])
				i++
			}
		}
		if y != r.size-1 {
			out += " "
		}
	}

	return out + "}"
}

// col2u32 represents a color as uint32
func col2u32(c color.Color) uint32 {
	r, g, b, a := c.RGBA()
	return ((r >> 8) << 24) | ((g >> 8) << 16) | (b & 0xff00) | (a >> 8)
}
