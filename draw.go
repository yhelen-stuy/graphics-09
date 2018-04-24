package main

import (
	"errors"
	"fmt"
	"math"
)

const (
	sphereStepSize float64 = 1.0 / 15
	torusStepSize  float64 = 1.0 / 10
)

func (image Image) DrawPolygons(p *Matrix, c Color) {
	m := p.mat
	cnew := c
	for i := 0; i < p.cols-2; i += 3 {
		p0 := []float64{
			m[0][i],
			m[1][i],
			m[2][i],
		}
		p1 := []float64{
			m[0][i+1],
			m[1][i+1],
			m[2][i+1],
		}
		p2 := []float64{
			m[0][i+2],
			m[1][i+2],
			m[2][i+2],
		}
		v1 := MakeVector(p0[0], p0[1], p0[2], p1[0], p1[1], p1[2])
		v2 := MakeVector(p0[0], p0[1], p0[2], p2[0], p2[1], p2[2])
		cross, err := CrossProduct(v1, v2)
		if err != nil {
			fmt.Println(err)
			continue
		}
		if cross[2] > 0 {
			cnew = Color{
				r: cnew.r + 37%255,
				g: cnew.g + 41%255,
				b: cnew.b + 59%255,
			}
			image.scanline(p0, p1, p2, cnew)
			image.DrawLine(cnew, int(p0[0]), int(p0[1]), int(p1[0]), int(p1[1]))
			image.DrawLine(cnew, int(p1[0]), int(p1[1]), int(p2[0]), int(p2[1]))
			image.DrawLine(cnew, int(p2[0]), int(p2[1]), int(p0[0]), int(p0[1]))
		}
	}
}

func (image Image) scanline(p0, p1, p2 []float64, c Color) {
	// Sort points
	if p1[1] < p0[1] {
		p0, p1 = p1, p0
	}

	if p2[1] < p1[1] {
		p1, p2 = p2, p1
	}

	if p1[1] < p0[1] {
		p0, p1 = p1, p0
	}

	// BM
	x0 := p0[0]
	x1 := x0
	d0 := (p2[0] - p0[0]) / (p2[1] - p0[1])
	d1 := (p1[0] - p0[0]) / (p1[1] - p0[1])
	fmt.Printf("d0: %.2f, d1: %.2f\n", d0, d1)
	for y := int(p0[1]); y < int(p1[1]); y++ {
		image.DrawLine(c, int(x0), y, int(x1), y)
		x0 += d0
		x1 += d1
	}

	// MT
	x1 = p1[0]
	d0 = (p2[0] - p0[0]) / (p2[1] - p0[1])
	d1 = (p2[0] - p1[0]) / (p2[1] - p1[1])
	fmt.Printf("d0: %.2f, d1: %.2f\n", d0, d1)
	for y := int(p1[1]); y < int(p2[1]); y++ {
		image.DrawLine(c, int(x0), y, int(x1), y)
		x0 += d0
		x1 += d1
	}
}

func (image Image) DrawLines(edges *Matrix, c Color) {
	m := edges.mat
	for i := 0; i < edges.cols-1; i += 2 {
		image.DrawLine(c, int(m[0][i]), int(m[1][i]), int(m[0][i+1]), int(m[1][i+1]))
	}
}

// TODO: Calculate & Plot z
func (image Image) DrawLine(c Color, x0, y0 int, z0 float64, x1, y1 int, z1 float64) error {
	if x0 < 0 || y0 < 0 || x1 > image.width || y1 > image.height {
		return errors.New("Error: Coordinates out of bounds")
	}
	if x0 > x1 {
		x1, x0 = x0, x1
		y1, y0 = y0, y1
	}

	deltaX := x1 - x0
	deltaY := y1 - y0
	lA := deltaY
	lB := deltaX * -1
	if deltaY >= 0 {
		if math.Abs(float64(deltaY)) <= math.Abs(float64(deltaX)) {
			// Octant 1 and 5
			lD := 2*lA + lB
			for x := x0; x <= x1; x++ {
				err := image.plot(c, x, y)
				if err != nil {
					return err
				}
				if lD > 0 {
					y++
					lD += 2 * lB
				}
				lD += 2 * lA
			}
		} else {
			// Octant 2 and 6
			x := x0
			lD := lA + 2*lB
			for y := y0; y <= y1; y++ {
				err := image.plot(c, x, y)
				if err != nil {
					return err
				}
				if lD < 0 {
					x++
					lD += 2 * lA
				}
				lD += 2 * lB
			}
		}
	} else {
		if math.Abs(float64(deltaY)) > math.Abs(float64(deltaX)) {
			// Octant 7 and 3
			x := x0
			lD := lA - 2*lB
			for y := y0; y >= y1; y-- {
				err := image.plot(c, x, y)
				if err != nil {
					return err
				}
				if lD > 0 {
					x++
					lD += 2 * lA
				}
				lD -= 2 * lB
			}
		} else {
			// Octant 8 and 4
			y := y0
			lD := 2*lA - lB
			for x := x0; x <= x1; x++ {
				err := image.plot(c, x, y)
				if err != nil {
					return err
				}
				if lD < 0 {
					y--
					lD -= 2 * lB
				}
				lD += 2 * lA
			}
		}
	}
	return nil
}

func (m *Matrix) AddCircle(cx, cy, cz, r float64) {
	var oldX, oldY float64 = -1, -1
	// TODO: No magic numbers wow i have so much to fix
	for i := 0; i <= 100; i++ {
		var t float64 = float64(i) / float64(100)
		x := r*math.Cos(2*math.Pi*t) + cx
		y := r*math.Sin(2*math.Pi*t) + cy
		if oldX < 0 || oldY < 0 {
			oldX = x
			oldY = y
			continue
		}
		m.AddEdge(oldX, oldY, cz, x, y, cz)
		oldX = x
		oldY = y
	}
}

func makeHermiteCoefs(p0, p1, rp0, rp1 float64) (*Matrix, error) {
	h := MakeMatrix(4, 0)
	h.AddCol([]float64{2, -3, 0, 1})
	h.AddCol([]float64{-2, 3, 0, 0})
	h.AddCol([]float64{1, -2, 1, 0})
	h.AddCol([]float64{1, -1, 0, 0})

	mat := MakeMatrix(4, 0)
	mat.AddCol([]float64{p0, p1, rp0, rp1})

	return mat.Mult(h)
}

func (m *Matrix) AddHermite(x0, y0, x1, y1, rx0, ry0, rx1, ry1, stepSize float64) error {
	xC, err := makeHermiteCoefs(x0, x1, rx0, rx1)
	if err != nil {
		return err
	}
	yC, err := makeHermiteCoefs(y0, y1, ry0, ry1)
	if err != nil {
		return err
	}
	// TODO: Figure out a better way to do this
	var oldX, oldY float64 = -1, -1
	var steps int = int(1 / stepSize)
	for i := 0; i <= steps; i++ {
		var t float64 = float64(i) / float64(steps)
		x := xC.mat[0][0]*math.Pow(t, 3.0) + xC.mat[1][0]*math.Pow(t, 2.0) + xC.mat[2][0]*t + xC.mat[3][0]
		y := yC.mat[0][0]*math.Pow(t, 3.0) + yC.mat[1][0]*math.Pow(t, 2.0) + yC.mat[2][0]*t + yC.mat[3][0]
		if oldX < 0 || oldY < 0 {
			oldX = x
			oldY = y
			continue
		}
		m.AddEdge(oldX, oldY, 0.0, x, y, 0.0)
		oldX = x
		oldY = y
	}
	return nil
}

func makeBezierCoefs(p0, p1, p2, p3 float64) (*Matrix, error) {
	h := MakeMatrix(4, 0)
	h.AddCol([]float64{-1, 3, -3, 1})
	h.AddCol([]float64{3, -6, 3, 0})
	h.AddCol([]float64{-3, 3, 0, 0})
	h.AddCol([]float64{1, 0, 0, 0})

	mat := MakeMatrix(4, 0)
	mat.AddCol([]float64{p0, p1, p2, p3})

	return mat.Mult(h)
}

// TODO: maybe combine with hermite bc a lot of duplicate code?
// Or make a separate parametric fxn
func (m *Matrix) AddBezier(x0, y0, x1, y1, x2, y2, x3, y3, stepSize float64) error {
	xC, err := makeBezierCoefs(x0, x1, x2, x3)
	if err != nil {
		return err
	}
	yC, err := makeBezierCoefs(y0, y1, y2, y3)
	if err != nil {
		return err
	}
	// TODO: Figure out a better way to do this
	var oldX, oldY float64 = -1, -1
	var steps int = int(1 / stepSize)
	for i := 0; i <= steps; i++ {
		var t float64 = float64(i) / float64(steps)
		x := xC.mat[0][0]*math.Pow(t, 3.0) + xC.mat[1][0]*math.Pow(t, 2.0) + xC.mat[2][0]*t + xC.mat[3][0]
		y := yC.mat[0][0]*math.Pow(t, 3.0) + yC.mat[1][0]*math.Pow(t, 2.0) + yC.mat[2][0]*t + yC.mat[3][0]
		if oldX < 0 || oldY < 0 {
			oldX = x
			oldY = y
			continue
		}
		m.AddEdge(oldX, oldY, 0.0, x, y, 0.0)
		oldX = x
		oldY = y
	}
	return nil
}

func (m *Matrix) AddBox(x, y, z, width, height, depth float64) {
	x1 := x + width
	y1 := y - height
	z1 := z - depth

	// Front
	m.AddPolygon(x, y, z, x, y1, z, x1, y1, z)
	m.AddPolygon(x, y, z, x1, y1, z, x1, y, z)

	// Back
	m.AddPolygon(x, y, z1, x1, y, z1, x, y1, z1)
	m.AddPolygon(x1, y1, z1, x, y1, z1, x1, y, z1)

	// Top
	m.AddPolygon(x, y, z, x1, y, z, x1, y, z1)
	m.AddPolygon(x1, y, z1, x, y, z1, x, y, z)

	// Bottom
	m.AddPolygon(x1, y1, z, x, y1, z, x, y1, z1)
	m.AddPolygon(x1, y1, z1, x1, y1, z, x, y1, z1)

	// Left
	m.AddPolygon(x, y1, z, x, y, z, x, y, z1)
	m.AddPolygon(x, y, z1, x, y1, z1, x, y1, z)

	// Right
	m.AddPolygon(x1, y, z, x1, y1, z1, x1, y, z1)
	m.AddPolygon(x1, y, z, x1, y1, z, x1, y1, z1)
}

func (m *Matrix) AddSphere(cx, cy, cz, r float64) {
	points := generateSpherePoints(cx, cy, cz, r)
	p := points.mat
	steps := int(1 / sphereStepSize)
	latStart, lonStart := 0, 0
	latEnd, lonEnd := steps, steps
	steps++
	for lat := latStart; lat < latEnd; lat++ {
		lat1 := lat * steps
		lat2 := (lat1 + steps) % points.cols
		for lon := lonStart; lon < lonEnd; lon++ {
			index := lat1 + lon
			indexLat2 := lat2 + lon
			// Only draw one triangle at poles
			if lon > 0 {
				m.AddPolygon(p[0][index], p[1][index], p[2][index],
					p[0][index+1], p[1][index+1], p[2][index+1],
					p[0][indexLat2], p[1][indexLat2], p[2][indexLat2])
			}
			if lon != lonEnd-1 {
				m.AddPolygon(
					p[0][index+1], p[1][index+1], p[2][index+1],
					p[0][indexLat2+1], p[1][indexLat2+1], p[2][indexLat2+1],
					p[0][indexLat2], p[1][indexLat2], p[2][indexLat2])
			}
		}
	}
}

func generateSpherePoints(cx, cy, cz, r float64) *Matrix {
	m := MakeMatrix(4, 0)
	var end float64
	// Rotating
	if sphereStepSize > 15 {
		end = 1 + sphereStepSize
	} else {
		end = 1.0
	}
	for i := 0.0; i <= end; i += sphereStepSize {
		phi := 2.0 * math.Pi * i
		// Semicircle
		for j := 0.0; j <= end; j += sphereStepSize {
			theta := math.Pi * j
			x := r*math.Cos(theta) + cx
			y := r*math.Sin(theta)*math.Cos(phi) + cy
			z := r*math.Sin(theta)*math.Sin(phi) + cz
			m.AddPoint(x, y, z)
		}
	}
	return m
}

func (m *Matrix) AddTorus(cx, cy, cz, r1, r2 float64) {
	points := generateTorusPoints(cx, cy, cz, r1, r2)
	p := points.mat
	steps := int(1 / torusStepSize)
	latStart, lonStart := 0, 0
	latEnd, lonEnd := steps, steps
	steps++
	for lat := latStart; lat < latEnd; lat++ {
		lat1 := lat * steps
		lat2 := (lat1 + steps) % points.cols
		for lon := lonStart; lon < lonEnd; lon++ {
			index := lat1 + lon
			indexLat2 := lat2 + lon
			// fmt.Printf("*(%d, %d, %d)\n", index, index+1, indexLat2)
			m.AddPolygon(p[0][index], p[1][index], p[2][index],
				p[0][index+1], p[1][index+1], p[2][index+1],
				p[0][indexLat2], p[1][indexLat2], p[2][indexLat2])
			// fmt.Printf("+(%d, %d, %d)\n", index+1, indexLat2+1, indexLat2)
			m.AddPolygon(
				p[0][index+1], p[1][index+1], p[2][index+1],
				p[0][indexLat2+1], p[1][indexLat2+1], p[2][indexLat2+1],
				p[0][indexLat2], p[1][indexLat2], p[2][indexLat2])
		}
	}
}

// r1: Radius of circle
// r2: Radius of torus
func generateTorusPoints(cx, cy, cz, r1, r2 float64) *Matrix {
	m := MakeMatrix(4, 0)
	var end float64
	if torusStepSize > 10 {
		end = 1 + torusStepSize
	} else {
		end = 1.0
	}
	// Rotating
	for i := 0.0; i < end; i += torusStepSize {
		phi := 2.0 * math.Pi * i
		// Circle
		for j := 0.0; j < end; j += torusStepSize {
			theta := 2.0 * math.Pi * j
			x := math.Cos(phi)*(r1*math.Cos(theta)+r2) + cx
			y := r1*math.Sin(theta) + cy
			z := -math.Sin(phi)*(r1*math.Cos(theta)+r2) + cz
			m.AddPoint(x, y, z)
		}
	}
	return m
}
