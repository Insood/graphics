package main

import (
	"errors"
	"math"
)

func Project(p Point3) (Point2, error) {
	adjZ := p.z - EyePosition.z // RHS, z negative into screen; relative to eye.z
	if adjZ > 0 {
		return Point2{}, errors.New("point is behind the camera")
	}
	adjZ *= -1 // Absolute z value for division

	return Point2{p.x / (adjZ * perspective), p.y / (adjZ * perspective)}, nil
}

func Rotate(v *Point3, theta float64) {
	rx := v.x*math.Cos(theta) - v.z*math.Sin(theta)
	rz := v.x*math.Sin(theta) + v.z*math.Cos(theta)
	v.x = rx
	v.z = rz

	rx = v.x*math.Cos(theta) - v.y*math.Sin(theta)
	ry := v.x*math.Sin(theta) + v.y*math.Cos(theta)
	v.x = rx
	v.y = ry

	ry = v.y*math.Cos(theta) - v.z*math.Sin(theta)
	rz = v.y*math.Sin(theta) + v.z*math.Cos(theta)
	v.y = ry
	v.z = rz
}
