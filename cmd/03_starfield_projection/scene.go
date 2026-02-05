package main

import (
	"image/color"
	"math/rand"
)

type Star struct {
	x float64
	y float64
	z float64
}

type Scene struct {
	game               *Game
	stars              []Star
	active             bool
	speed              float64
	starAppearDistance float64
}

func makeScene(game *Game) *Scene {
	scene := Scene{game: game, active: true, speed: 1, starAppearDistance: -500}

	for range starCount {
		x := float64(rand.Intn(screenWidth) - screenWidth/2)
		y := float64(rand.Intn(screenHeight) - screenHeight/2)
		z := rand.Float64() * scene.starAppearDistance
		scene.stars = append(scene.stars, Star{x, y, z})
	}

	return &scene
}

func (s *Scene) Update() {
	if !s.active {
		return
	}

	s.speed += 0.01

	for i := range s.stars {
		s.UpdateStar(&s.stars[i])
	}
}

func (s *Scene) Draw() {
	for i := range s.stars {
		s.game.PushMatrix()
		s.DrawStar(&s.stars[i])
		s.game.PopMatrix()
	}
}

func (s *Scene) UpdateStar(star *Star) {
	star.z += s.speed

	if star.z > 0 {
		star.z = s.starAppearDistance
	}
}

func (s *Scene) DrawStar(star *Star) {
	s.game.TranslateModel(star.x, star.y, star.z)
	xy := s.game.Project(star.x, star.y)

	// Ebiten images are stored as premultiplied alpha
	alpha := uint8(255 * (1 - (star.z / s.starAppearDistance)))
	s.game.SetColor(color.RGBA{R: alpha, G: alpha, B: alpha, A: 0})
	s.game.DrawPixel(int(xy[0]), int(xy[1]))
}
