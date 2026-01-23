package main

import "math"

type Triangle struct {
	p1 Point3
	p2 Point3
	p3 Point3

	// projected data. On the screen raster
	pp1 Point2
	pp2 Point2
	pp3 Point2
}

func newTriangle(p1, p2, p3 Point3) *Triangle {
	return &Triangle{p1, p2, p3, Point2{}, Point2{}, Point2{}}
}

func (t *Triangle) project() {
	t.pp1, _ = Project(t.p1)
	t.pp2, _ = Project(t.p2)
	t.pp3, _ = Project(t.p3)
}

func (t *Triangle) min_px() int {
	return int(math.Min(t.pp1.x, math.Min(t.pp2.x, t.pp3.x)))
}

func (t *Triangle) max_px() int {
	return int(math.Max(t.pp1.x, math.Max(t.pp2.x, t.pp3.x)))
}

func (t *Triangle) min_py() int {
	return int(math.Min(t.pp1.y, math.Min(t.pp2.y, t.pp3.y)))
}

func (t *Triangle) max_py() int {
	return int(math.Max(t.pp1.y, math.Max(t.pp2.y, t.pp3.y)))
}

func (t *Triangle) normal() Vector3 {
	v1 := t.p2.subtract(t.p1)
	v2 := t.p3.subtract(t.p1)
	return v2.cross(v1).normalize()
}

// [!] This assumes that the triangle belongs to a sphere and that the points p1,p2,p3
//     are points on the sphere and are effectively the normal
func (t *Triangle) sphericalFaceNormal() Vector3 {
	average_face_normal := Vector3{
		(t.p1.x + t.p2.x + t.p3.x) / 3,
		(t.p1.y + t.p2.y + t.p3.y) / 3,
		(t.p1.z + t.p2.z + t.p3.z) / 3,
	}

	return average_face_normal.normalize()
}

func (t *Triangle) baryCentricCoordinates(p Point2Int) (bool, Point2) {
	v0 := t.pp3.round().subtract(t.pp1.round())
	v1 := t.pp2.round().subtract(t.pp1.round())
	v2 := p.subtract(t.pp1.round())

	dot00 := v0.dot(v0)
	dot01 := v0.dot(v1)
	dot02 := v0.dot(v2)
	dot11 := v1.dot(v1)
	dot12 := v1.dot(v2)

	invDenom := float64(1) / float64(dot00*dot11-dot01*dot01)
	u := float64(dot11*dot02-dot01*dot12) * invDenom
	v := float64(dot00*dot12-dot01*dot02) * invDenom

	return u >= 0 && v >= 0 && u+v <= 1, Point2{u, v}
}

func makeSampleTriangle(size int) []*Triangle {
	return []*Triangle{
		newTriangle(
			Point3{float64(0), float64(size), float64(size)},
			Point3{float64(size), float64(-size), float64(size)},
			Point3{float64(-size), float64(-size), float64(size)},
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

			pt1 := Point3{x1, y12, z1}
			pt2 := Point3{x2, y12, z2}
			pt3 := Point3{x3, y34, z3}
			pt4 := Point3{x4, y34, z4}

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
