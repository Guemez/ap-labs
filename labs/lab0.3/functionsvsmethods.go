package main

import (
	"math"
	"fmt"
)

type Point struct{
	X float64 
	Y float64
}

// traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

type Path []Point

func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}

func printPoint(p Point) {
    fmt.Printf("Points: %f, %f", p.X, p.Y)
}

func main() {
    p := Point{2,4}
    printPoint(p)
	
}
