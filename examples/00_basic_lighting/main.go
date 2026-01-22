package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 640
	perspective  = 0.002 // 1/500
	delta        = 0.01  // Rotation speed
)

var LightSource = Point3{200, 200, 350}
var EyePosition = Point3{0, 0, 600}
var OutlineColor = Color3{1.0, 0.2, 0.5} // Red-ish
var FillColor = Color3{1.0, 1.0, 1.0}

type Game struct {
	canvas           *ebiten.Image
	triangles        []*Triangle // Original geometry
	rotatedTriangles []*Triangle
	currentColor     color.RGBA
	theta            float64
	cullBackFaces    bool
	drawOutline      bool
}

func NewGame() *Game {
	// tris := makeSampleTriangle(100)
	tris := makeSphere(150, 10)

	rotatedTriangles := make([]*Triangle, len(tris))

	for i := range rotatedTriangles {
		rotatedTriangles[i] = &Triangle{}
	}

	return &Game{
		triangles:        tris,
		rotatedTriangles: rotatedTriangles,
		canvas:           ebiten.NewImage(screenWidth, screenHeight),
		currentColor:     color.RGBA{255, 0, 0, 255},
		theta:            0,
		cullBackFaces:    true,
		drawOutline:      true,
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.cullBackFaces = !g.cullBackFaces
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		g.drawOutline = !g.drawOutline
	}

	g.theta += delta
	for g.theta > math.Pi*2 {
		g.theta -= math.Pi * 2
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.canvas.Clear()

	g.RotateTriangles()
	g.DrawTriangles()
	screen.DrawImage(g.canvas, nil)
}

func (g *Game) SetColor(pixel_color Color3) {
	g.currentColor = color.RGBA{
		uint8(math.Min(float64(255), math.Max(float64(0), pixel_color.r*255))),
		uint8(math.Min(float64(255), math.Max(float64(0), pixel_color.g*255))),
		uint8(math.Min(float64(255), math.Max(float64(0), pixel_color.b*255))),
		255,
	}
}

func (g *Game) DrawTriangles() {
	for _, t := range g.rotatedTriangles {
		t.project()
		g.DrawTriangle(t)
	}
}

func (g *Game) DrawTriangle(t *Triangle) {
	vecA := t.pp3.subtract(t.pp1)
	vecB := t.pp2.subtract(t.pp1)
	cross := vecA.cross(vecB)
	if cross < 0 && g.cullBackFaces { // Backface culling
		return
	}

	g.FillTriangle(t)

	if g.drawOutline {
		g.DrawOutline(t)
	}
}

func (g *Game) FillTriangle(t *Triangle) {
	minx := t.min_px()
	maxx := t.max_px()
	miny := t.min_py()
	maxy := t.max_py()

	for y := maxy; y >= miny; y-- {
		for x := minx; x <= maxx; x++ {
			in_triangle, _ := t.baryCentricCoordinates(Point2Int{x, y})

			if !in_triangle {
				continue
			}

			g.SetColor(FillColor)
			g.DrawPixel(x, y)
		}
	}
}

func (g *Game) DrawOutline(t *Triangle) {
	g.SetColor(OutlineColor)
	g.DrawLine(t.pp1, t.pp2)
	g.DrawLine(t.pp2, t.pp3)
	g.DrawLine(t.pp3, t.pp1)
}

// set the pixel using current color. 0,0 is the middle, x axis right, y going up
func (g *Game) DrawPixel(x, y int) {
	x += screenWidth / 2                    // offset by half screen
	y = screenHeight - (y + screenHeight/2) // offset by half screen and reverse Y direction

	if x < 0 || x >= screenWidth || y < 0 || y >= screenHeight {
		return
	}

	g.canvas.Set(x, y, g.currentColor)
}

func (g *Game) DrawLine(start, end Point2) {
	dx := end.x - start.x
	dy := end.y - start.y
	absdx := math.Abs(dx)
	absdy := math.Abs(dy)

	stepx := 1
	if start.x > end.x {
		stepx = -1
	}

	stepy := 1
	if start.y > end.y {
		stepy = -1
	}

	errY := 0.0
	errX := 0.0
	x := int(start.x)
	y := int(start.y)

	// Line has more rise than run
	if max(-dy, dy) > max(-dx, dx) {
		slope := math.Abs(float64(dx) / float64(dy))

		for ystep := 0; ystep < int(absdy); ystep++ {
			g.DrawPixel(x, y)
			y += stepy
			errX += slope
			if errX > 0.5 {
				errX -= 1
				x += stepx
			}
		}
	} else {
		slope := math.Abs(float64(dy) / float64(dx))
		for xstep := 0; xstep < int(absdx); xstep++ {
			g.DrawPixel(x, y)
			x += stepx
			errY += slope
			if errY > 0.5 {
				errY -= 1
				y += stepy
			}
		}
	}
}

func (g *Game) RotateTriangles() {
	for i, original_tri := range g.triangles {
		rotated_tri := g.rotatedTriangles[i]
		rotated_tri.p1 = original_tri.p1 // Copy by value
		rotated_tri.p2 = original_tri.p2
		rotated_tri.p3 = original_tri.p3

		Rotate(&rotated_tri.p1, g.theta) // Rotate in place
		Rotate(&rotated_tri.p2, g.theta)
		Rotate(&rotated_tri.p3, g.theta)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Basic Lighting")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
