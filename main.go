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
	"log"
	"os"

	"github.com/meisterluk/carpet/config"
	d "github.com/meisterluk/carpet/draw"
)

const HEADER = `  𝛑𝛑𝛑 Carpet 𝛑𝛑𝛑

  This binary implements carpets as presented by Jack Hodkinson
  in a Youtube video "Carpets, Genetics, and the Pi Fractal" [0].
  Implementation by meisterluk under MIT license.

  Version 1.0.0

  usage: %s %s %s %s
            [%s]

  [0] https://friendlyfieldsandopenmaps.com/2017/09/18/the-pi-fractal/

`

// main routine
func main() {
	fmt.Printf(HEADER, os.Args[0], osArgOr(1, `<rule-files-dir:string>`), osArgOr(2, `<#iterations:int>`), osArgOr(3, `<output-file.png:string>`), osArgOr(4, `<initial-color:8hex>`))
	if len(os.Args) == 2 && (os.Args[1] == `-h` || os.Args[1] == `--help`) {
		return
	}
	if err := do(); err != nil {
		panic(err)
	}
}

// do is the core of the main routine
func do() error {
	var conf config.Config
	if err := conf.Read(os.Args); err != nil {
		return err
	}
	log.Printf("%d rules found\n", len(conf.Rules))

	growth := conf.Rules[0].Size()

	// check color responsibilities in rules
	ruleColors := make(map[uint32]color.Color)
	var oneColor color.Color
	for _, rule := range conf.Rules {
		for y := 0; y < growth; y++ {
			for x := 0; x < growth; x++ {
				col := rule.Colorize(x, y)
				ruleColors[col2u32(col)] = col
				if len(ruleColors) == 1 || col2u32(col) > col2u32(oneColor) {
					oneColor = col
				}
			}
		}
	}

	out := ""
	i := 0
	for _, col := range ruleColors {
		r, g, b, a := col.RGBA()
		if i%9 == 0 {
			out += "\n  "
		}
		out += fmt.Sprintf("%02X%02X%02X%02X ", r>>8, g>>8, b>>8, a>>8)
		i++
	}
	log.Printf("colors found in rules:%s\n", out)

	for _, col := range ruleColors {
		anyResponsible := false
		for _, rule := range conf.Rules {
			if rule.Matches(col) {
				anyResponsible = true
			}
		}
		if !anyResponsible {
			r, g, b, a := col.RGBA()
			return fmt.Errorf("color %02X%02X%02X%02X occurs in rules, but no rule matches this color", r>>8, g>>8, b>>8, a>>8)
		}
	}
	anyResponsible := false
	for _, rule := range conf.Rules {
		if rule.Matches(conf.InitialColor) {
			anyResponsible = true
		}
	}
	if !anyResponsible {
		return fmt.Errorf("initialization color is %s, please provide a rule for %s", colorRepr(conf.InitialColor), colorRepr(conf.InitialColor))
	}

	return d.Draw(&conf)
}
