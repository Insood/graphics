package graphicsmath

import "math"

type Vector2Int struct {
	X int
	Y int
}

type Vector2 struct {
	X float64
	Y float64
}

func (v1 Vector2) Cross(v2 Vector2) float64 {
	return v1.X*v2.Y - v1.Y*v2.X
}

type Vector3 struct {
	X float64
	Y float64
	Z float64
}

type Color3 struct {
	R float64
	G float64
	B float64
}

func (v1 Vector3) Cross(v2 Vector3) Vector3 {
	return Vector3{
		v1.Y*v2.Z - v1.Z*v2.Y,
		v1.Z*v2.X - v1.X*v2.Z,
		v1.X*v2.Y - v1.Y*v2.X,
	}
}

func (v Vector3) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector3) Normalize() Vector3 {
	return Vector3{
		v.X / v.Magnitude(),
		v.Y / v.Magnitude(),
		v.Z / v.Magnitude(),
	}
}

func (v1 Vector3) Dot(v2 Vector3) float64 {
	return v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z
}

func (v1 Vector2Int) Dot(v2 Vector2Int) int {
	return v1.X*v2.X + v1.Y*v2.Y
}

func (v1 Vector2) Round() Vector2Int {
	return Vector2Int{
		int(math.Round(v1.X)),
		int(math.Round(v1.Y)),
	}
}

func (v1 Vector3) Subtract(v2 Vector3) Vector3 {
	return Vector3{
		v1.X - v2.X,
		v1.Y - v2.Y,
		v1.Z - v2.Z,
	}
}

func (v1 Vector2) Subtract(v2 Vector2) Vector2 {
	return Vector2{
		v1.X - v2.X,
		v1.Y - v2.Y,
	}
}

func (v1 Vector2Int) Subtract(v2 Vector2Int) Vector2Int {
	return Vector2Int{
		v1.X - v2.X,
		v1.Y - v2.Y,
	}
}

func (v Vector3) Multiply(s float64) Vector3 {
	return Vector3{
		v.X * s,
		v.Y * s,
		v.Z * s,
	}
}

func (c Color3) Multiply(s float64) Color3 {
	return Color3{
		c.R * s,
		c.G * s,
		c.B * s,
	}
}

func (c Color3) Add(c2 Color3) Color3 {
	return Color3{
		c.R + c2.R,
		c.G + c2.G,
		c.B + c2.B,
	}
}

func (p Vector3) Add(v Vector3) Vector3 {
	return Vector3{
		p.X + v.X,
		p.Y + v.Y,
		p.Z + v.Z,
	}
}
