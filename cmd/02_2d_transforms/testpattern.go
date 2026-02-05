package main

import "image/color"

func (g *Game) DrawTestPattern(scale float64) {
	left := -scale / 2
	right := scale / 2
	top := -scale / 2
	bottom := scale / 2

	g.SetProjection()

	red := 1.0
	green := 0.0
	blue := 0.0

	for i := range 10 {
		for j := range 10 {
			x := left + scale/testGridSize*float64(i)
			y := top + scale/testGridSize*float64(j)

			if i > testGridSize/2 {
				green = 1.0
			} else {
				green = 0.0
			}

			if j > testGridSize/2 {
				blue = 1.0
			} else {
				blue = 0.0
			}
			g.SetColor(color.RGBA{R: uint8(red * 255), G: uint8(green * 255), B: uint8(blue * 255), A: 255})

			v1 := g.Project(left, y)
			v2 := g.Project(right, y)
			v3 := g.Project(x, bottom)
			v4 := g.Project(x, top)

			g.DrawLine(v1, v2)
			g.DrawLine(v3, v4)
		}
	}
}
