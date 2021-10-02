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

// Size is the size of the network
const Size = 7

// Weights are the weights of the network
var Weights = [Size]float64{3, 5, 7, 9, 11, 13, 15}

func rng(x [Size]float64) [Size]float64 {
	var y [Size]float64
	for key, value := range Weights {
		y[key] = math.Sin(value * math.Pi * x[key])
	}
	return y
}

func main() {
	x := Weights
	random := make(plotter.Values, 0, 1024)
	for i := 0; i < 1024*1024; i++ {
		x = rng(x)
		sum := 0.0
		for _, value := range x {
			sum += value
		}
		random = append(random, sum)
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
