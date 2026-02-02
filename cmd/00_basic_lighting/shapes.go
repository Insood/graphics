package main

import (
	"math"

	mymath "github.com/insood/graphics/internal/math"
)

type Triangle struct {
	p1 mymath.Vector3
	p2 mymath.Vector3
	p3 mymath.Vector3

	// projected data. On the screen raster
	pp1 mymath.Vector2
	pp2 mymath.Vector2
	pp3 mymath.Vector2
}

func newTriangle(p1, p2, p3 mymath.Vector3) *Triangle {
	return &Triangle{p1, p2, p3, mymath.Vector2{}, mymath.Vector2{}, mymath.Vector2{}}
}

func (t *Triangle) project() {
	t.pp1, _ = Project(t.p1)
	t.pp2, _ = Project(t.p2)
	t.pp3, _ = Project(t.p3)
}

func (t *Triangle) min_px() int {
	return int(math.Min(t.pp1.X, math.Min(t.pp2.X, t.pp3.X)))
}

func (t *Triangle) max_px() int {
	return int(math.Max(t.pp1.X, math.Max(t.pp2.X, t.pp3.X)))
}

func (t *Triangle) min_py() int {
	return int(math.Min(t.pp1.Y, math.Min(t.pp2.Y, t.pp3.Y)))
}

func (t *Triangle) max_py() int {
	return int(math.Max(t.pp1.Y, math.Max(t.pp2.Y, t.pp3.Y)))
}

func (t *Triangle) normal() mymath.Vector3 {
	v1 := t.p2.Subtract(t.p1)
	v2 := t.p3.Subtract(t.p1)
	return v2.Cross(v1).Normalize()
}

// [!] This assumes that the triangle belongs to a sphere and that the points p1,p2,p3
//
//	are points on the sphere and are effectively the normal
func (t *Triangle) sphericalFaceNormal() mymath.Vector3 {
	average_face_normal := mymath.Vector3{
		X: (t.p1.X + t.p2.X + t.p3.X) / 3,
		Y: (t.p1.Y + t.p2.Y + t.p3.Y) / 3,
		Z: (t.p1.Z + t.p2.Z + t.p3.Z) / 3,
	}

	return average_face_normal.Normalize()
}

func (t *Triangle) baryCentricCoordinates(p mymath.Vector2Int) (bool, mymath.Vector2) {
	v0 := t.pp3.Round().Subtract(t.pp1.Round())
	v1 := t.pp2.Round().Subtract(t.pp1.Round())
	v2 := p.Subtract(t.pp1.Round())

	dot00 := v0.Dot(v0)
	dot01 := v0.Dot(v1)
	dot02 := v0.Dot(v2)
	dot11 := v1.Dot(v1)
	dot12 := v1.Dot(v2)

	invDenom := float64(1) / float64(dot00*dot11-dot01*dot01)
	u := float64(dot11*dot02-dot01*dot12) * invDenom
	v := float64(dot00*dot12-dot01*dot02) * invDenom

	return u >= 0 && v >= 0 && u+v <= 1, mymath.Vector2{X: u, Y: v}
}

func makeSampleTriangle(size int) []*Triangle {
	return []*Triangle{
		newTriangle(
			mymath.Vector3{X: float64(0), Y: float64(size), Z: float64(size)},
			mymath.Vector3{X: float64(size), Y: float64(-size), Z: float64(size)},
			mymath.Vector3{X: float64(-size), Y: float64(-size), Z: float64(size)},
		),
	}
}

func makeSphere(radius int, divisions int) []*Triangle {
	tris := []*Triangle{}

	for phi_step := 0; phi_step < divisions; phi_step++ {
		for theta_step := 0; theta_step < (divisions * 2); theta_step++ {
			phi1 := math.Pi * float64(phi_step) / float64(divisions)
			phi2 := math.Pi * float64(phi_step+1) / float64(divisions)
			theta1 := 2 * math.Pi * float64(theta_step) / float64(divisions*2)
			theta2 := 2 * math.Pi * float64(theta_step+1) / float64(divisions*2)

			y12 := float64(radius) * math.Cos(phi1)
			y34 := float64(radius) * math.Cos(phi2)

			x1 := float64(radius) * math.Sin(phi1) * math.Sin(theta1)
			x2 := float64(radius) * math.Sin(phi1) * math.Sin(theta2)
			x3 := float64(radius) * math.Sin(phi2) * math.Sin(theta2)
			x4 := float64(radius) * math.Sin(phi2) * math.Sin(theta1)

			z1 := float64(radius) * math.Sin(phi1) * math.Cos(theta1)
			z2 := float64(radius) * math.Sin(phi1) * math.Cos(theta2)
			z3 := float64(radius) * math.Sin(phi2) * math.Cos(theta2)
			z4 := float64(radius) * math.Sin(phi2) * math.Cos(theta1)

			pt1 := mymath.Vector3{X: x1, Y: y12, Z: z1}
			pt2 := mymath.Vector3{X: x2, Y: y12, Z: z2}
			pt3 := mymath.Vector3{X: x3, Y: y34, Z: z3}
			pt4 := mymath.Vector3{X: x4, Y: y34, Z: z4}

			switch phi_step {
			case 0: // Top
				tris = append(tris, newTriangle(pt1, pt3, pt4))
			case divisions - 1: // Bottom
				tris = append(tris, newTriangle(pt1, pt2, pt4))
			default:
				tris = append(tris, newTriangle(pt1, pt2, pt3))
				tris = append(tris, newTriangle(pt1, pt3, pt4))
			}
		}
	}

	return tris
}
