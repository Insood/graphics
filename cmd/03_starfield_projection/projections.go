package main

import (
	"math"

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

func viewFrustum(left, right, bottom, top, near, far float64) matrix.Mat4 {
	return matrix.Mat4FromRows(
		matrix.Vec4{2.0 * near / (right - left), 0, (left + right) / (left - right), 0},
		matrix.Vec4{0, 2.0 * near / (top - bottom), (bottom + top) / (bottom - top), 0},
		matrix.Vec4{0, 0, (far + near) / (near - far), (2 * far * near) / (far - near)},
		matrix.Vec4{0, 0, 1, 0},
	)
}

func fovToWidth(fov, near float64) float64 {
	return 2 * near * math.Tan(fov/2)
}
