// Copyright 2021 The Rand Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
	"math/rand"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	"github.com/pointlander/gradient/sf64"
)

// Size is the size of the network
const Size = 7

// Weights are the weights of the network
var Weights = RNG([]float64{3, 5, 7, 9, 11, 13, 15})

// RNG is an RNG network
type RNG []float64

// RNG setps the network
func (r RNG) RNG(state []float64) {
	for key, value := range Weights {
		state[key] = math.Sin(value * math.Pi * state[key])
	}
}

func main() {
	s := make([]float64, len(Weights))
	copy(s, Weights)
	random := make(plotter.Values, 0, 1024)
	for i := 0; i < 1024*1024; i++ {
		Weights.RNG(s)
		sum := 0.0
		for _, value := range s {
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

	var weights = [Size]sf64.V{}
	for i := range weights {
		weights[i].X = 7 * rand.ExpFloat64()
	}
	var state = [Size]sf64.Meta{}
	for i := range state {
		state[i] = weights[i].Meta()
	}
	var outputs = [3]sf64.V{}
	for i := range outputs {
		outputs[i].X = 3 * rand.ExpFloat64()
	}
	var sum = [3]sf64.Meta{}
	for j := range sum {
		for i := range state {
			state[i] = sf64.Sin(sf64.Mul(weights[0].Meta(), state[i]))
			if sum[j] == nil {
				sum[j] = state[i]
			} else {
				sum[j] = sf64.Add(sum[j], state[i])
			}
		}
	}
	cost := sf64.Sub(outputs[0].Meta(), sum[0])
	cost = sf64.Mul(cost, cost)
	for i := 1; i < len(outputs); i++ {
		diff := sf64.Sub(outputs[i].Meta(), sum[i])
		diff = sf64.Mul(diff, diff)
		cost = sf64.Add(cost, diff)
	}

	eta, iterations := .001, 1024
	points := make(plotter.XYs, 0, iterations)
	i := 0
	for i < iterations {
		total := 0.0
		for i := range weights {
			weights[i].D = 0
		}

		total += sf64.Gradient(cost).X
		s := 0.0
		for i := range weights {
			d := weights[i].D
			s += d * d
		}
		norm := math.Sqrt(s)
		scaling := 1.0
		if norm > 1 {
			scaling = 1 / norm
		}

		for i := range weights {
			weights[i].X -= eta * weights[i].D * scaling
		}

		points = append(points, plotter.XY{X: float64(i), Y: total})

		fmt.Println(i, total)
		i++
		if total < 1e-4 {
			for j := range outputs {
				var y float64
				sum[j](func(a *sf64.V) bool {
					y = a.X
					return true
				})
				fmt.Println(outputs[j].X, y)
			}
			break
		}
	}

	p = plot.New()

	p.Title.Text = "epochs vs cost"
	p.X.Label.Text = "epochs"
	p.Y.Label.Text = "cost"

	scatter, err := plotter.NewScatter(points)
	if err != nil {
		panic(err)
	}
	scatter.GlyphStyle.Radius = vg.Length(1)
	scatter.GlyphStyle.Shape = draw.CircleGlyph{}
	p.Add(scatter)

	err = p.Save(8*vg.Inch, 8*vg.Inch, "cost.png")
	if err != nil {
		panic(err)
	}

	for i := range weights {
		Weights[i] = weights[i].X
	}
	s = make([]float64, len(Weights))
	copy(s, Weights)
	random = make(plotter.Values, 0, 1024)
	for i := 0; i < 1024*1024; i++ {
		Weights.RNG(s)
		sum := 0.0
		for _, value := range s {
			sum += value
		}
		random = append(random, sum)
	}

	p = plot.New()
	p.Title.Text = "random numbers"
	histogram, err = plotter.NewHist(random, 256)
	if err != nil {
		panic(err)
	}
	p.Add(histogram)
	err = p.Save(8*vg.Inch, 8*vg.Inch, "histogram_learned.png")
	if err != nil {
		panic(err)
	}
}
