package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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

var LightSource = Vector3{200, 200, 350}
var EyePosition = Vector3{0, 0, 600}
var OutlineColor = Color3{1.0, 0.2, 0.5} // Red-ish
var FillColor = Color3{1.0, 1.0, 1.0}
var NormalColor = Color3{0.0, 1.0, 0.0}

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
	currentColor     Color3
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
		currentColor:     Color3{},
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

func (g *Game) SetColor(pixel_color Color3) {
	g.currentColor = pixel_color
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
	v1Color := g.PhongLighting(t.p1.ToVector3().normalize())
	v2Color := g.PhongLighting(t.p2.ToVector3().normalize())
	v3Color := g.PhongLighting(t.p3.ToVector3().normalize())

	for y := maxy; y >= miny; y-- {
		for x := minx; x <= maxx; x++ {
			in_triangle, uv := t.baryCentricCoordinates(Point2Int{x, y})

			if !in_triangle {
				continue
			}

			switch g.drawMode {
			case Flat:
				g.SetColor(FillColor)
			case Barycentric:
				g.SetColor(Color3{uv.x, uv.y, 1 - uv.x - 1.*uv.y})
			case PhongFace:
				g.SetColor(faceColor)
			case PhongVertex:
				g.SetColor(averageVertexColor)
			case PhongGourand:
				a := v3Color.multiply(uv.x)
				b := v2Color.multiply(uv.y)
				c := v1Color.multiply(1 - uv.x - uv.y)

				g.SetColor(a.add(b).add(c))
			case PhongShading:
				a := t.p1.ToVector3().multiply(1 - uv.x - uv.y)
				b := t.p2.ToVector3().multiply(uv.y)
				c := t.p3.ToVector3().multiply(uv.x)
				normal := a.add(b).add(c).normalize()
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
	start := Point3{
		(t.p1.x + t.p2.x + t.p3.x) / 3,
		(t.p1.y + t.p2.y + t.p3.y) / 3,
		(t.p1.z + t.p2.z + t.p3.z) / 3,
	}

	end := start.add(t.normal().multiply(20))

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
		uint8(math.Min(float64(255), math.Max(float64(0), g.currentColor.r*255))),
		uint8(math.Min(float64(255), math.Max(float64(0), g.currentColor.g*255))),
		uint8(math.Min(float64(255), math.Max(float64(0), g.currentColor.b*255))),
		255,
	}

	g.canvas.Set(x, y, pixelColor)
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

func (g *Game) PhongLighting(normal Vector3) Color3 {
	face_color := Color3{0.0, 0.0, 0.0}

	ambient := FillColor.multiply(ambientMaterial)
	face_color = face_color.add(ambient)

	light_normal := LightSource.normalize()
	diffuse_component := normal.dot(light_normal)

	diffuse := FillColor.multiply(diffuse_component * diffuseMaterial)
	face_color = face_color.add(diffuse)

	reflection := normal.multiply(2).multiply(LightSource.dot(normal)).subtract(LightSource)
	specular_component := reflection.normalize().dot(EyePosition.normalize())

	specular := FillColor.multiply(specularMaterial * math.Pow(specular_component, shininess))
	face_color = face_color.add(specular)

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
