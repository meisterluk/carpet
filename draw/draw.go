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

package draw

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/meisterluk/carpet/config"
	"github.com/meisterluk/carpet/rules"
)

// colorRepr returns a human-readable color representation
func colorRepr(c color.Color) string {
	r, g, b, a := c.RGBA()
	return fmt.Sprintf("%02X%02X%02X%02X", r>>8, g>>8, b>>8, a>>8)
}

// absolute returns the absolute coordinate for a given relative coordinate.
// i.e. it maps a coordinate `coord` at level `iter` to level `iterations`.
//      where one level is expected to grow by factor `growth` per iteration.
// example. if there is only 1 square (iter=0), but there are 3 iterations with growth 3,
//          the coord 0 is mapped to 26.
func absolute(coord, iter, iterations, growth int) int {
	f := int(math.Pow(float64(growth), float64(iterations-iter)))
	return (coord+1)*f - 1
}

// Draw takes a config and applies the carpet algorithm.
// It generates a PNG image based on the data in the Config.
func Draw(conf *config.Config) error {
	// create output file
	growth := conf.Rules[0].Size()
	totalSize := int(math.Pow(float64(growth), float64(conf.Iterations)))
	img := image.NewNRGBA(image.Rect(0, 0, totalSize, totalSize))
	b := img.Bounds()

	// initialize image with white
	log.Printf("initializing an image with color %s\n", colorRepr(conf.InitialColor))
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			img.Set(x, y, conf.InitialColor)
		}
	}

	// apply iterations
	for i := 0; i < conf.Iterations; i++ {
		log.Printf("Iteration %d", i)
		size := int(math.Pow(float64(growth), float64(i)))
		for prevY := 0; prevY < size; prevY++ {
			for prevX := 0; prevX < size; prevX++ {
				absPrevY := absolute(prevY, i, conf.Iterations, growth)
				absPrevX := absolute(prevX, i, conf.Iterations, growth)

				givenColor := img.At(absPrevX, absPrevY)
				//fmt.Printf("absolute coord of %d,%d is %d,%d - its color is %v\n", prevY, prevX, absPrevY, absPrevX, givenColor)

				var matchingRule rules.Rule
				for _, rule := range conf.Rules {
					if rule.Matches(givenColor) {
						matchingRule = rule
					}
				}
				if matchingRule == nil {
					return fmt.Errorf(`no rule found for color %s at %d,%d`, colorRepr(givenColor), prevX, prevY)
				}
				//fmt.Printf("rule %s is responsible for color %v at %d,%d\n", matchingRule.String(), givenColor, absPrevY, absPrevX)

				// for each square in the replacement pattern
				for y := 0; y < growth; y++ {
					for x := 0; x < growth; x++ {
						// we determine the absolute and set its color
						subY := absolute(prevY, i, i+1, growth)
						subX := absolute(prevX, i, i+1, growth)
						subY = subY - growth + y + 1
						subX = subX - growth + x + 1
						subY = absolute(subY, i+1, conf.Iterations, growth)
						subX = absolute(subX, i+1, conf.Iterations, growth)
						//fmt.Printf("%d,%d+%d,%d → %d,%d   %d\n", prevY, prevX, y, x, subY, subX, size)

						//fmt.Printf("set %d,%d with offset %d,%d (which is %d,%d) to color %v\n", prevY, prevX, y, x, subY, subX, matchingRule.Colorize(x, y))
						img.Set(subX, subY, matchingRule.Colorize(x, y))
					}
				}
			}
		}
	}

	// write file
	f, err := os.Create(conf.OutputPath)
	if err != nil {
		return err
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		return err
	}

	return f.Close()
}
