package main

import (
	"errors"
	"math"
)

// Given 2 points P0 and P1, return the vector P0P1
func MakeVector(x0, y0, z0, x1, y1, z1 float64) []float64 {
	p := make([]float64, 3)
	p[0] = x1 - x0
	p[1] = y1 - y0
	p[2] = z1 - z0
	return p
}

func CrossProduct(a, b []float64) ([]float64, error) {
	if len(a) != 3 || len(b) != 3 {
		return nil, errors.New("Invalid vectors")
	}
	p := make([]float64, 3)
	p[0] = a[1]*b[2] - a[2]*b[1]
	p[1] = a[2]*b[0] - a[0]*b[2]
	p[2] = a[0]*b[1] - a[1]*b[0]
	return p, nil
}

func DotProduct(a, b []float64) (float64, error) {
	if len(a) != len(b) {
		// Is this really the best course of action...?
		return math.NaN(), errors.New("Cannot dot vectors of different length")
	}
	p := 0.0
	for i := range a {
		p += a[i] * b[i]
	}
	return p, nil
}

func Magnitude(v []float64) float64 {
	sum := 0.0
	for i := range v {
		sum += v[i] * v[i]
	}
	return math.Sqrt(sum)
}
