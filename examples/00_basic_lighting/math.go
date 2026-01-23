package main

import "math"

type Point2 struct {
	x float64
	y float64
}

type Point2Int struct {
	x int
	y int
}

type Vector2Int struct {
	x int
	y int
}

type Vector2 struct {
	x float64
	y float64
}

func (v1 Vector2) cross(v2 Vector2) float64 {
	return v1.x*v2.y - v1.y*v2.x
}

type Float3 struct {
	x float64
	y float64
	z float64
}

type Point3 struct {
	x float64
	y float64
	z float64
}

type Vector3 struct {
	x float64
	y float64
	z float64
}

type Color3 struct {
	r float64
	g float64
	b float64
}

func (v1 Vector3) cross(v2 Vector3) Vector3 {
	return Vector3{
		v1.y*v2.z - v1.z*v2.y,
		v1.z*v2.x - v1.x*v2.z,
		v1.x*v2.y - v1.y*v2.x,
	}
}

func (v Vector3) magnitude() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v Vector3) normalize() Vector3 {
	return Vector3{
		v.x / v.magnitude(),
		v.y / v.magnitude(),
		v.z / v.magnitude(),
	}
}

func (v1 Vector3) dot(v2 Vector3) float64 {
	return v1.x*v2.x + v1.y*v2.y + v1.z*v2.z
}

func (v1 Vector2Int) dot(v2 Vector2Int) int {
	return v1.x*v2.x + v1.y*v2.y
}

func (v1 Point2) round() Point2Int {
	return Point2Int{
		int(math.Round(v1.x)),
		int(math.Round(v1.y)),
	}
}

func (v1 Point3) subtract(v2 Point3) Vector3 {
	return Vector3{
		v1.x - v2.x,
		v1.y - v2.y,
		v1.z - v2.z,
	}
}

func (v1 Vector3) subtract(v2 Vector3) Vector3 {
	return Vector3{
		v1.x - v2.x,
		v1.y - v2.y,
		v1.z - v2.z,
	}
}

func (v1 Point2) subtract(v2 Point2) Vector2 {
	return Vector2{
		v1.x - v2.x,
		v1.y - v2.y,
	}
}

func (v1 Point2Int) subtract(v2 Point2Int) Vector2Int {
	return Vector2Int{
		v1.x - v2.x,
		v1.y - v2.y,
	}
}

func (v Vector3) multiply(s float64) Vector3 {
	return Vector3{
		v.x * s,
		v.y * s,
		v.z * s,
	}
}

func (c Color3) multiply(s float64) Color3 {
	return Color3{
		c.r * s,
		c.g * s,
		c.b * s,
	}
}

func (c Color3) add(c2 Color3) Color3 {
	return Color3{
		c.r + c2.r,
		c.g + c2.g,
		c.b + c2.b,
	}
}

func (p Point3) add(v Vector3) Point3 {
	return Point3{
		p.x + v.x,
		p.y + v.y,
		p.z + v.z,
	}
}

func (v Vector3) add(v2 Vector3) Vector3 {
	return Vector3{
		v.x + v2.x,
		v.y + v2.y,
		v.z + v2.z,
	}
}

func (p Point3) ToVector3() Vector3 {
	return Vector3{p.x, p.y, p.z}
}
