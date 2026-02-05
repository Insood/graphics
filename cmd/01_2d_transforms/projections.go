package main

import (
	matrix "github.com/go-gl/mathgl/mgl64"
)

func getViewport(width, height float64) matrix.Mat4 {
	viewportMatrix := matrix.Ident4()

	translate := matrix.Mat4FromRows(
		matrix.Vec4{1, 0, 0, width / 2},
		matrix.Vec4{0, 1, 0, height / 2},
		matrix.Vec4{0, 0, 1, 0},
		matrix.Vec4{0, 0, 0, 1},
	)

	scale := matrix.Mat4FromRows(
		matrix.Vec4{width / 2, 0, 0, 0},
		matrix.Vec4{0, -height / 2, 0, 0},
		matrix.Vec4{0, 0, 1, 0},
		matrix.Vec4{0, 0, 0, 1},
	)

	viewportMatrix = viewportMatrix.Mul4(translate)
	viewportMatrix = viewportMatrix.Mul4(scale)
	return viewportMatrix
}

func getCamera(up matrix.Vec2, center matrix.Vec2, zoom float64) matrix.Mat4 {
	translate := matrix.Mat4FromRows(
		matrix.Vec4{1, 0, 0, -center.X()},
		matrix.Vec4{0, 1, 0, -center.Y()},
		matrix.Vec4{0, 0, 1, 0},
		matrix.Vec4{0, 0, 0, 1},
	)

	scale := matrix.Mat4FromRows(
		matrix.Vec4{zoom, 0, 0, 0},
		matrix.Vec4{0, zoom, 0, 0},
		matrix.Vec4{0, 0, 1, 0},
		matrix.Vec4{0, 0, 0, 1},
	)

	orient := matrix.Mat4FromRows(
		matrix.Vec4{up.Y(), -up.X(), 0, 0},
		matrix.Vec4{up.X(), up.Y(), 0, 0},
		matrix.Vec4{0, 0, 1, 0},
		matrix.Vec4{0, 0, 0, 1},
	)

	viewMatrix := scale.Mul4(translate)
	viewMatrix = orient.Mul4(viewMatrix)
	return viewMatrix
}

func getOrtho(left, right, bottom, top float64) matrix.Mat4 {
	far := -1.0
	near := 1.0

	translate := matrix.Mat4FromRows(
		matrix.Vec4{1, 0, 0, -(left + right) / 2},
		matrix.Vec4{0, 1, 0, -(top + bottom) / 2},
		matrix.Vec4{0, 0, -1, -(far + near) / 2}, // includes flip here
		matrix.Vec4{0, 0, 0, 1},
	)

	scale := matrix.Mat4FromRows(
		matrix.Vec4{2 / (right - left), 0, 0, 0},
		matrix.Vec4{0, 2 / (top - bottom), 0, 0},
		matrix.Vec4{0, 0, 2 / (far - 1), 0},
		matrix.Vec4{0, 0, 0, 1},
	)

	return scale.Mul4(translate)
}
