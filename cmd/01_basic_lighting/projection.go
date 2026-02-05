package main

import (
	"errors"
	"math"

	mymath "github.com/insood/graphics/internal/math"
)

func Project(p mymath.Vector3) (mymath.Vector2, error) {
	adjZ := p.Z - EyePosition.Z // RHS, z negative into screen; relative to eye.Z
	if adjZ > 0 {
		return mymath.Vector2{}, errors.New("point is behind the camera")
	}
	adjZ *= -1 // Absolute z value for division

	return mymath.Vector2{X: p.X / (adjZ * perspective), Y: p.Y / (adjZ * perspective)}, nil
}

func Rotate(v *mymath.Vector3, theta float64) {
	rx := v.X*math.Cos(theta) - v.Z*math.Sin(theta)
	rz := v.X*math.Sin(theta) + v.Z*math.Cos(theta)
	v.X = rx
	v.Z = rz

	rx = v.X*math.Cos(theta) - v.Y*math.Sin(theta)
	ry := v.X*math.Sin(theta) + v.Y*math.Cos(theta)
	v.X = rx
	v.Y = ry

	ry = v.Y*math.Cos(theta) - v.Z*math.Sin(theta)
	rz = v.Y*math.Sin(theta) + v.Z*math.Cos(theta)
	v.Y = ry
	v.Z = rz
}
