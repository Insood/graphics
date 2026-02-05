package main

import (
	"image/color"
	"math"

	matrix "github.com/go-gl/mathgl/mgl64"
)

type Gear struct {
	teeth         int
	x             float64
	y             float64
	radius        float64
	rotationSpeed float64
	rotation      float64
	color         color.RGBA
}

type RingGear struct {
	teeth     int
	x         float64
	y         float64
	radius    float64
	thickness float64
	rotation  float64
	color     color.RGBA
}

type Scene struct {
	active                     bool
	sunGear                    Gear
	ringGear                   RingGear
	planetaryGears             []Gear
	planetaryRadius            float64
	planetaryGearRotationSpeed float64
	planetaryGearRotation      float64
}

func makeScene() *Scene {
	scene := Scene{active: true, planetaryGearRotationSpeed: 0.008}

	scene.sunGear = Gear{teeth: 20, x: 0, y: 0, radius: 0.1, rotationSpeed: 0.042, rotation: 0.1, color: color.RGBA{R: 255, G: 255, B: 0, A: 255}}
	scene.ringGear = RingGear{teeth: 100, x: 0, y: 0, thickness: 0.9, radius: 0.64, rotation: 0.02, color: color.RGBA{R: 255, G: 0, B: 0, A: 255}}

	for i := range 3 {
		y := math.Cos(float64(i)*(2*math.Pi)/3) * 0.34
		x := math.Sin(float64(i)*(2*math.Pi)/3) * 0.34
		scene.planetaryGears = append(scene.planetaryGears, Gear{teeth: 40, x: x, y: y, radius: 0.2, rotationSpeed: -0.02, color: color.RGBA{R: 0, G: 255, B: 0, A: 255}})
	}

	return &scene
}

func (s *Scene) update() {
	if !s.active {
		return
	}

	s.sunGear.update()

	s.planetaryGearRotation += s.planetaryGearRotationSpeed

	if s.planetaryGearRotation > (2 * math.Pi) {
		s.planetaryGearRotation -= (2 * math.Pi)
	} else if s.planetaryGearRotation < 0 {
		s.planetaryGearRotation += (2 * math.Pi)
	}

	for i := range s.planetaryGears {
		s.planetaryGears[i].update()
	}
}

func (g *Gear) update() {
	g.rotation += g.rotationSpeed
	if g.rotation > (2 * math.Pi) {
		g.rotation -= (2 * math.Pi)
	} else if g.rotation < 0 {
		g.rotation += (2 * math.Pi)
	}
}

func (g *Game) DrawScene() {
	g.SetColor(color.RGBA{R: 255, G: 255, B: 255, A: 255})
	g.scene.update()

	g.DrawGear(&g.scene.sunGear)

	g.PushMatrix()
	g.RotateModel(g.scene.planetaryGearRotation)
	for i := range g.scene.planetaryGears {
		g.DrawGear(&g.scene.planetaryGears[i])
	}
	g.PopMatrix()

	g.DrawRingGear(&g.scene.ringGear)
}

func (g *Game) DrawGear(gear *Gear) {
	g.SetColor(gear.color)
	g.PushMatrix()
	g.TranslateModel(gear.x, gear.y, 0)
	g.RotateModel(gear.rotation)
	g.ScaleModel(gear.radius)
	g.DrawGearSegments(gear.teeth)
	g.PopMatrix()
}

func (g *Game) DrawGearSegments(teeth int) {
	arc := (2 * math.Pi) / float64(teeth)

	for i := range teeth {
		g.PushMatrix()
		g.RotateModel(arc * float64(i))
		g.DrawHubPiece(arc)
		g.TranslateModel(0, 1, 0)
		g.drawGearTooth(arc)
		g.PopMatrix()
	}
}

func (g *Game) drawGearTooth(arc float64) {
	g.DrawTriangle(
		matrix.Vec2{-math.Sin(arc) / 2, 0},
		matrix.Vec2{0, math.Sin(arc)},
		matrix.Vec2{math.Sin(arc) / 2, 0},
	)
}

func (g *Game) DrawHubPiece(arc float64) {
	g.DrawTriangle(
		matrix.Vec2{0, 0},
		matrix.Vec2{-math.Sin(arc) / 2, math.Cos(arc)},
		matrix.Vec2{math.Sin(arc) / 2, math.Cos(arc)},
	)
}

func (g *Game) DrawRingGear(gear *RingGear) {
	arc := (2 * math.Pi) / float64(gear.teeth)
	g.SetColor(gear.color)
	g.PushMatrix()
	g.ScaleModel(gear.radius)
	g.RotateModel(gear.rotation)

	for i := range gear.teeth {
		g.PushMatrix()
		g.RotateModel(arc * float64(i))
		g.DrawRingGearSegment(arc, gear.thickness)
		g.PushMatrix()
		g.TranslateModel(0, gear.thickness, 0)
		g.RotateModel(math.Pi)
		g.ScaleModel(0.75)
		g.drawGearTooth(arc)
		g.PopMatrix()
		g.PopMatrix()
	}

	g.PopMatrix()
}

func (g *Game) DrawRingGearSegment(arc, raceThickness float64) {
	a := matrix.Vec2{math.Cos(arc), math.Sin(arc)}
	b := matrix.Vec2{math.Cos(arc) * raceThickness, math.Sin(arc) * raceThickness}
	c := matrix.Vec2{1, 0}
	d := matrix.Vec2{raceThickness, 0}

	g.DrawTriangle(a, b, c) //
	g.DrawTriangle(b, c, d)
}
