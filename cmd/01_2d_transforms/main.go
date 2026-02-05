package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	matrix "github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 640
	testGridSize = 10
)

// Draw Mode (Test Pattern or Scene)
const (
	TestPattern = iota
	SceneLayout
)

// Projection & View Matrix Mode
const (
	Identity      = iota // Use identity matrix
	Center640            // 0,0 at center
	BottomLeft640        // 0,0 at bottom left
	FlipX                // same as Center640, but X is flipped
	Aspect               // uneven aspect ratio: x is from -320 to 320, y is from -100 to 100
)

type Game struct {
	drawMode  int
	debugMode bool

	canvas       *ebiten.Image
	currentColor color.RGBA

	projectionMode   int
	projectionMatrix matrix.Mat4
	viewMatrix       matrix.Mat4
	viewportMatrix   matrix.Mat4
	modelMatrixStack []matrix.Mat4

	cameraTarget matrix.Vec2
	cameraZoom   float64
	cameraRotate float64

	mouseDragging bool
	mouseLastX    int
	mouseLastY    int

	scene *Scene
}

func NewGame() *Game {
	return &Game{
		drawMode:  SceneLayout,
		debugMode: false,

		canvas:       ebiten.NewImage(screenWidth, screenHeight),
		currentColor: color.RGBA{},

		projectionMode:   Identity,
		projectionMatrix: matrix.Ident4(),
		viewMatrix:       matrix.Ident4(),
		viewportMatrix:   getViewport(screenWidth, screenHeight),
		modelMatrixStack: []matrix.Mat4{matrix.Ident4()},

		cameraTarget: matrix.Vec2{},
		cameraZoom:   1.0,
		cameraRotate: 0,

		mouseDragging: false,
		mouseLastX:    0,
		mouseLastY:    0,

		scene: makeScene(),
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		g.projectionMode++
		if g.projectionMode > Aspect {
			g.projectionMode = Identity
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debugMode = !g.debugMode
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		g.drawMode++
		if g.drawMode > SceneLayout {
			g.drawMode = TestPattern
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) { // Actually +
		g.cameraZoom *= 1.1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) { // -
		g.cameraZoom *= 0.9
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBracketRight) {
		g.cameraRotate += math.Pi / (2 * 10)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyBracketLeft) {
		g.cameraRotate -= math.Pi / (2 * 10)
	}

	mouseX, mouseY := ebiten.CursorPosition()

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.mouseDragging = true
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.mouseDragging = false
	}

	if g.mouseDragging {
		g.mouseDragged(mouseX-g.mouseLastX, mouseY-g.mouseLastY)
	}

	_, wheelY := ebiten.Wheel()

	if wheelY != 0 {
		g.wheelScrolled(wheelY)
	}

	g.mouseLastX = mouseX
	g.mouseLastY = mouseY

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.scene.active = !g.scene.active
	}

	return nil
}

func (g *Game) mouseDragged(dx, dy int) {
	xmove := float64(dx) / float64(screenWidth)
	ymove := float64(dy) / float64(screenHeight)

	// xMove/Ymove give a percentage of how much along the entire actual screen
	// the mouse has moved. By multiplying by 2, we can convert that to a delta NDC
	// PVector deltaNDC = new PVector(-xMove*2, yMove*2);

	deltaNDC := matrix.Vec4{-xmove * 2, ymove * 2, 0, 0}

	inverseProjection := matrix.Mat4(g.projectionMatrix)
	inverseProjection.Set(0, 3, 0)
	inverseProjection.Set(1, 3, 0)
	inverseProjection.Set(2, 3, 0)
	inverseProjection = inverseProjection.Inv()

	inverseViewMatrix := matrix.Mat4(g.viewMatrix)
	inverseViewMatrix.Set(0, 3, 0)
	inverseViewMatrix.Set(1, 3, 0)
	inverseViewMatrix.Set(2, 3, 0)
	inverseViewMatrix = inverseViewMatrix.Inv()

	deltaView := inverseProjection.Mul4x1(deltaNDC)
	deltaWorld := inverseViewMatrix.Mul4x1(deltaView)

	g.cameraTarget[0] += deltaWorld.X()
	g.cameraTarget[1] += deltaWorld.Y()
}

func (g *Game) wheelScrolled(wheelY float64) {
	mouseX, mouseY := ebiten.CursorPosition()

	mouseNDC := matrix.Vec4{
		float64(mouseX)*2.0/float64(screenWidth) - 1.0,
		-(float64(mouseY)*2.0/float64(screenHeight) - 1.0),
		0,
		0,
	}

	invertedProjectedMatrix := matrix.Mat4(g.projectionMatrix).Inv()
	invertedViewMatrix := matrix.Mat4(g.viewMatrix).Inv()

	mouseView := invertedProjectedMatrix.Mul4x1(mouseNDC)
	mouseWorld := invertedViewMatrix.Mul4x1(mouseView)

	g.cameraTarget = mouseWorld.Vec2()

	if wheelY > 0 {
		g.cameraZoom *= 0.9
	} else if wheelY < 0 {
		g.cameraZoom *= 1.1
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.canvas.Clear()
	g.SetProjection()

	switch g.drawMode {
	case TestPattern:
		g.DrawTestPattern(1)
		g.DrawTestPattern(100)
		g.DrawTestPattern(1000)

	case SceneLayout:
		g.DrawScene()
	}

	screen.DrawImage(g.canvas, nil)
}

func (g *Game) SetColor(color color.RGBA) {
	g.currentColor = color
}

func (g *Game) SetProjection() {
	upVector := matrix.Vec2{0, 1}
	rotate := matrix.Mat2FromRows(
		matrix.Vec2{math.Cos(g.cameraRotate), math.Sin(g.cameraRotate)},
		matrix.Vec2{-math.Sin(g.cameraRotate), math.Cos(g.cameraRotate)},
	)

	upVector = rotate.Mul2x1(upVector)

	switch g.projectionMode {
	case Identity:
		g.projectionMatrix = matrix.Ident4()
	case Center640:
		g.projectionMatrix = getOrtho(-320, 320, -320, 320)
	case BottomLeft640:
		g.projectionMatrix = getOrtho(0, 640, 0, 640)
	case FlipX:
		g.projectionMatrix = getOrtho(320, -320, -320, 320)
	case Aspect:
		g.projectionMatrix = getOrtho(-320, 320, -100, 100)
	}

	g.viewMatrix = getCamera(upVector, g.cameraTarget, g.cameraZoom)
}

func (g *Game) PushMatrix() {
	g.modelMatrixStack = append(g.modelMatrixStack, matrix.Mat4(g.modelMatrixStack[len(g.modelMatrixStack)-1]))
}

func (g *Game) PopMatrix() {
	g.modelMatrixStack = g.modelMatrixStack[:len(g.modelMatrixStack)-1]
}

func (g *Game) ScaleModel(s float64) {
	scale := matrix.Mat4FromRows(
		matrix.Vec4{s, 0, 0, 0},
		matrix.Vec4{0, s, 0, 0},
		matrix.Vec4{0, 0, 1, 0},
		matrix.Vec4{0, 0, 0, 1},
	)

	g.modelMatrixStack[len(g.modelMatrixStack)-1] = g.modelMatrixStack[len(g.modelMatrixStack)-1].Mul4(scale)
}

func (g *Game) TranslateModel(x, y, z float64) {
	translate := matrix.Mat4FromRows(
		matrix.Vec4{1, 0, 0, x},
		matrix.Vec4{0, 1, 0, y},
		matrix.Vec4{0, 0, 1, z},
		matrix.Vec4{0, 0, 0, 1},
	)

	g.modelMatrixStack[len(g.modelMatrixStack)-1] = g.modelMatrixStack[len(g.modelMatrixStack)-1].Mul4(translate)
}

func (g *Game) RotateModel(angle float64) {
	rotate := matrix.Mat4FromRows(
		matrix.Vec4{math.Cos(angle), -math.Sin(angle), 0, 0},
		matrix.Vec4{math.Sin(angle), math.Cos(angle), 0, 0},
		matrix.Vec4{0, 0, 1, 0},
		matrix.Vec4{0, 0, 0, 1},
	)

	g.modelMatrixStack[len(g.modelMatrixStack)-1] = g.modelMatrixStack[len(g.modelMatrixStack)-1].Mul4(rotate)
}

func (g *Game) Project(worldx, worldy float64) matrix.Vec2 {
	world := matrix.Vec4{worldx, worldy, 0, 1}

	model := g.modelMatrixStack[len(g.modelMatrixStack)-1].Mul4x1(world)
	camera := g.viewMatrix.Mul4x1(model)
	ndc := g.projectionMatrix.Mul4x1(camera)

	screen := g.viewportMatrix.Mul4x1(ndc)

	if g.debugMode {
		fmt.Println("camera: ", camera)
		fmt.Println("ndc: ", ndc)
		fmt.Println("screen: ", screen)
	}

	return screen.Vec2()
}

func (g *Game) DrawLine(start, end matrix.Vec2) {
	vector.StrokeLine(g.canvas, float32(start[0]), float32(start[1]), float32(end[0]), float32(end[1]), 1, g.currentColor, false)
}

func (g *Game) DrawTriangle(modelA, modelB, modelC matrix.Vec2) {
	screenA := g.Project(modelA[0], modelA[1])
	screenB := g.Project(modelB[0], modelB[1])
	screenC := g.Project(modelC[0], modelC[1])

	g.DrawLine(screenA, screenB)
	g.DrawLine(screenB, screenC)
	g.DrawLine(screenC, screenA)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2D Transforms")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
