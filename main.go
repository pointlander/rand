// Copyright 2021 The Rand Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func rng(x float64) float64 {
	return math.Sin(7.0 * math.Pi * x)
}

func main() {
	x := 7.0
	random := make(plotter.Values, 0, 1024)
	for i := 0; i < 1024*1024; i++ {
		x = rng(x)
		random = append(random, x)
	}

	p := plot.New()
	p.Title.Text = "random numbers"
	histogram, err := plotter.NewHist(random, 256)
	if err != nil {
		panic(err)
	}
	p.Add(histogram)
	err = p.Save(8*vg.Inch, 8*vg.Inch, "histogram.png")
	if err != nil {
		panic(err)
	}
}
