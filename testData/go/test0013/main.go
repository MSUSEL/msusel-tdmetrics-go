//go:build test

package main

// Test nesting structs and interfaces

type XCoord struct{ x int }

func (x XCoord) GetX() int { return x.x }

type YCoord struct{ y int }

func (y YCoord) GetY() int { return y.y }

type Point struct {
	XCoord
	YCoord
}

func (p Point) Sum() int { return p.x + p.y }

type (
	IXCoord interface{ GetX() int }
	IYCoord interface{ GetY() int }
)

type IPoint interface {
	IXCoord
	IYCoord
	interface {
		Sum() int
	}
}

func PrintPoint(p IPoint) {
	println(p.GetX(), `+`, p.GetY(), `=`, p.Sum()) // 12 + 34 = 46
}

func main() {
	p := Point{XCoord{x: 12}, YCoord{y: 34}}
	PrintPoint(p)
}
