package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	mymath "github.com/insood/graphics/internal/math"
)

const (
	screenWidth      = 640
	screenHeight     = 640
	perspective      = 0.002 // 1/500
	delta            = 0.01  // Rotation speed
	ambientMaterial  = 0.35
	diffuseMaterial  = 0.45
	specularMaterial = 0.3
	shininess        = 30
)

var LightSource = mymath.Vector3{X: 200, Y: 200, Z: 350}
var EyePosition = mymath.Vector3{X: 0, Y: 0, Z: 600}
var OutlineColor = mymath.Color3{R: 1.0, G: 0.2, B: 0.5} // Red-ish
var FillColor = mymath.Color3{R: 1.0, G: 1.0, B: 1.0}
var NormalColor = mymath.Color3{R: 0.0, G: 1.0, B: 0.0}

const (
	None = iota
	Flat
	Barycentric
	PhongFace
	PhongVertex
	PhongGourand
	PhongShading
)

type Game struct {
	canvas           *ebiten.Image
	triangles        []*Triangle // Original geometry
	rotatedTriangles []*Triangle
	currentColor     mymath.Color3
	theta            float64
	rotate           bool
	cullBackFaces    bool
	drawOutline      bool
	drawNormals      bool
	drawMode         int
}

func NewGame() *Game {
	// tris := makeSampleTriangle(100)
	tris := makeSphere(250, 20)

	rotatedTriangles := make([]*Triangle, len(tris))

	for i := range rotatedTriangles {
		rotatedTriangles[i] = &Triangle{}
	}

	return &Game{
		triangles:        tris,
		rotatedTriangles: rotatedTriangles,
		canvas:           ebiten.NewImage(screenWidth, screenHeight),
		currentColor:     mymath.Color3{},
		theta:            0,
		rotate:           false,
		cullBackFaces:    true,
		drawOutline:      true,
		drawNormals:      false,
		drawMode:         None,
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.rotate = !g.rotate
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		g.cullBackFaces = !g.cullBackFaces
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		g.drawOutline = !g.drawOutline
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		g.drawNormals = !g.drawNormals
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		g.drawMode++
		if g.drawMode > PhongShading {
			g.drawMode = None
		}
	}

	if g.rotate {
		g.theta += delta
		for g.theta > math.Pi*2 {
			g.theta -= math.Pi * 2
		}
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

func (g *Game) SetColor(pixel_color mymath.Color3) {
	g.currentColor = pixel_color
}

func (g *Game) DrawTriangles() {
	for _, t := range g.rotatedTriangles {
		t.project()
		g.DrawTriangle(t)
	}
}

func (g *Game) DrawTriangle(t *Triangle) {
	vecA := t.pp3.Subtract(t.pp1)
	vecB := t.pp2.Subtract(t.pp1)
	cross := vecA.Cross(vecB)
	if cross < 0 && g.cullBackFaces { // Backface culling
		return
	}

	if g.drawMode != None {
		g.FillTriangle(t)
	}

	if g.drawOutline {
		g.DrawOutline(t)
	}

	if g.drawNormals {
		g.DrawNormal(t)
	}
}

func (g *Game) FillTriangle(t *Triangle) {
	minx := t.min_px()
	maxx := t.max_px()
	miny := t.min_py()
	maxy := t.max_py()

	faceColor := g.PhongLighting(t.normal())
	averageVertexColor := g.PhongLighting(t.sphericalFaceNormal())
	v1Color := g.PhongLighting(t.p1.Normalize())
	v2Color := g.PhongLighting(t.p2.Normalize())
	v3Color := g.PhongLighting(t.p3.Normalize())

	for y := maxy; y >= miny; y-- {
		for x := minx; x <= maxx; x++ {
			in_triangle, uv := t.baryCentricCoordinates(mymath.Vector2Int{X: x, Y: y})

			if !in_triangle {
				continue
			}

			switch g.drawMode {
			case Flat:
				g.SetColor(FillColor)
			case Barycentric:
				g.SetColor(mymath.Color3{R: uv.X, G: uv.Y, B: 1 - uv.X - 1.*uv.Y})
			case PhongFace:
				g.SetColor(faceColor)
			case PhongVertex:
				g.SetColor(averageVertexColor)
			case PhongGourand:
				a := v3Color.Multiply(uv.X)
				b := v2Color.Multiply(uv.Y)
				c := v1Color.Multiply(1 - uv.X - uv.Y)

				g.SetColor(a.Add(b).Add(c))
			case PhongShading:
				a := t.p1.Multiply(1 - uv.X - uv.Y)
				b := t.p2.Multiply(uv.Y)
				c := t.p3.Multiply(uv.X)
				normal := a.Add(b).Add(c).Normalize()
				g.SetColor(g.PhongLighting(normal))
			}

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

func (g *Game) DrawNormal(t *Triangle) {
	start := mymath.Vector3{
		X: (t.p1.X + t.p2.X + t.p3.X) / 3,
		Y: (t.p1.Y + t.p2.Y + t.p3.Y) / 3,
		Z: (t.p1.Z + t.p2.Z + t.p3.Z) / 3,
	}

	end := start.Add(t.normal().Multiply(20))

	screenStart, _ := Project(start)
	screenEnd, _ := Project(end)

	g.SetColor(NormalColor)
	g.DrawLine(screenStart, screenEnd)
}

// set the pixel using current color. 0,0 is the middle, x axis right, y going up
func (g *Game) DrawPixel(x, y int) {
	x += screenWidth / 2                    // offset by half screen
	y = screenHeight - (y + screenHeight/2) // offset by half screen and reverse Y direction

	if x < 0 || x >= screenWidth || y < 0 || y >= screenHeight {
		return
	}

	pixelColor := color.RGBA{
		uint8(math.Min(float64(255), math.Max(float64(0), g.currentColor.R*255))),
		uint8(math.Min(float64(255), math.Max(float64(0), g.currentColor.G*255))),
		uint8(math.Min(float64(255), math.Max(float64(0), g.currentColor.B*255))),
		255,
	}

	g.canvas.Set(x, y, pixelColor)
}

func (g *Game) DrawLine(start, end mymath.Vector2) {
	dx := end.X - start.X
	dy := end.Y - start.Y
	absdx := math.Abs(dx)
	absdy := math.Abs(dy)

	stepx := 1
	if start.X > end.X {
		stepx = -1
	}

	stepy := 1
	if start.Y > end.Y {
		stepy = -1
	}

	errY := 0.0
	errX := 0.0
	x := int(start.X)
	y := int(start.Y)

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

func (g *Game) PhongLighting(normal mymath.Vector3) mymath.Color3 {
	face_color := mymath.Color3{R: 0.0, G: 0.0, B: 0.0}

	ambient := FillColor.Multiply(ambientMaterial)
	face_color = face_color.Add(ambient)

	light_normal := LightSource.Normalize()
	diffuse_component := normal.Dot(light_normal)

	diffuse := FillColor.Multiply(diffuse_component * diffuseMaterial)
	face_color = face_color.Add(diffuse)

	reflection := normal.Multiply(2).Multiply(LightSource.Dot(normal)).Subtract(LightSource)
	specular_component := reflection.Normalize().Dot(EyePosition.Normalize())

	specular := FillColor.Multiply(specularMaterial * math.Pow(specular_component, shininess))
	face_color = face_color.Add(specular)

	return face_color
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
