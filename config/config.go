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

package config

import (
	"encoding/hex"
	"fmt"
	"image/color"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/meisterluk/carpet/rules"
)

// Config encompasses all parameters of this application.
// Thus for one specific Config, the output will be deterministic.
type Config struct {
	Iterations   int
	Rules        []rules.Rule
	InitialColor color.Color
	OutputPath   string
}

// Read takes CLI `args` and potentially other parameters
// to initialize its Config struct.
func (c *Config) Read(args []string) error {
	if len(args) != 4 && len(args) != 5 {
		return fmt.Errorf(`expected 3 or 4 CLI arguments; got %d; see usage`, len(args)-1)
	}

	// prepare rules directory
	rulesDir := args[1]

	// prepare number of iterations
	iter, err := strconv.Atoi(args[2])
	if err != nil {
		return err
	}
	if iter < 0 || 1000 < iter {
		return fmt.Errorf(`expected 0 < iter < 1000; got iter == %d`, iter)
	}
	c.Iterations = iter

	// prepare output file
	c.OutputPath = args[3]

	// prepare rules
	refGrowth := -1
	c.Rules = make([]rules.Rule, 0, 8)
	fileInfos, err := ioutil.ReadDir(rulesDir)
	if err != nil {
		return err
	}

	// iterate over all PNGs
	for _, info := range fileInfos {
		if !strings.HasSuffix(info.Name(), `.png`) {
			continue
		}

		// reading rule from file
		fullPath := filepath.Join(rulesDir, info.Name())
		r, err := rules.NewRuleFromPNGFile(fullPath)
		if err != nil && strings.Contains(err.Error(), `RGBA`) {
			continue
		} else if err != nil {
			return err
		}

		// check dimensions
		if refGrowth < 0 {
			refGrowth = r.Size()
		} else if r.Size() != refGrowth {
			return fmt.Errorf(`expected dimensions %d×%d like the first rule parsed; got %d×%d in '%s'`, refGrowth, refGrowth, r.Size(), r.Size(), fullPath)
		}

		c.Rules = append(c.Rules, r)
	}

	if len(c.Rules) == 0 {
		return fmt.Errorf(`I need at least one rule file to run (directory '%s' does not contain any)`, rulesDir)
	}

	// prepare initial color
	if len(args) == 5 {
		if len(args[4]) != 8 {
			return fmt.Errorf(`expected 8 hexadecimal characters for initial color; got %s`, args[4])
		}
		h, err := hex.DecodeString(args[4])
		if err != nil {
			return err
		}
		c.InitialColor = color.RGBA{R: uint8(h[0]), G: uint8(h[1]), B: uint8(h[2]), A: uint8(h[3])}
	} else {
		c.InitialColor = color.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	}

	return nil
}
