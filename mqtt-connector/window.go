package main

import (
	"math"
)

type DataWindow struct {
	Size int16
	Data []Payload
}

func newDataWindow() *DataWindow {
	d := DataWindow{0, []Payload{}}
	//d.Size = 0
	//d.Data = make([]Payload, 0)
	return &d
}

func (window *DataWindow) Add(payload Payload) {
	window.Data = append(window.Data, payload)
	window.Size = window.Size + 1
}

func (window *DataWindow) Reset() {
	window.Size = 0
	window.Data = window.Data[:0]
}

func (window *DataWindow) getRootMeanSquareOfVelocity() float64 {
	var square float64

	var n int = 1
	// Calculate square.
	for _, payload := range window.Data {
		for _, val := range payload.Velocity {
			square += math.Pow(float64(val), 2)
			n++
		}
	}
	// Calculate Mean.
	mean := square / float64(n)

	// Calculate Root.
	root := math.Sqrt(mean)
	return root
}

func (window *DataWindow) getRootMeanSquareOfDistance() float64 {
	var square float64

	var n int = 1
	// Calculate square.
	for _, payload := range window.Data {
		for _, val := range payload.Distance {
			square += math.Pow(float64(val), 2)
			n++
		}
	}
	// Calculate Mean.
	mean := square / float64(n)

	// Calculate Root.
	root := math.Sqrt(mean)
	return root
}

func (window *DataWindow) getMeanOfVelocity() float64 {
	var n int = 1
	var sum float64 = 0

	// Calculate square.
	for _, payload := range window.Data {
		for _, val := range payload.Velocity {
			sum += float64(val)
			n++
		}
	}
	// Calculate Mean.
	mean := sum / float64(n)
	return mean
}

func (window *DataWindow) getMeanOfDistance() float64 {
	var n int = 1
	var sum float64 = 0

	// Calculate square.
	for _, payload := range window.Data {
		for _, val := range payload.Distance {
			sum += float64(val)
			n++
		}
	}
	// Calculate Mean.
	mean := sum / float64(n)
	return mean
}

func (window *DataWindow) getMinOfDistance() float64 {
	var min float64 = 100000000.0

	// min
	for _, payload := range window.Data {
		for _, val := range payload.Distance {
			min = math.Min(min, float64(val))
		}
	}

	return min
}

func (window *DataWindow) getMinOfVelocity() float64 {
	var min float64 = 100000000.0

	// min
	for _, payload := range window.Data {
		for _, val := range payload.Velocity {
			min = math.Min(min, float64(val))
		}
	}

	return min
}

func (window *DataWindow) getMaxOfDistance() float64 {
	var max float64 = -100000000.0

	// max
	for _, payload := range window.Data {
		for _, val := range payload.Distance {
			max = math.Max(max, float64(val))
		}
	}

	return max
}

func (window *DataWindow) getMaxOfVelocity() float64 {
	var max float64 = -100000000.0

	// max
	for _, payload := range window.Data {
		for _, val := range payload.Velocity {
			max = math.Max(max, float64(val))
		}
	}

	return max
}
