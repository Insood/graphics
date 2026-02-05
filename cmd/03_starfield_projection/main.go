package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	matrix "github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	near         = -10
	far          = -500
	starCount    = 500
	fov          = math.Pi / 2
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

	scene *Scene
}

func NewGame() *Game {
	game := Game{
		debugMode: false,

		canvas:       ebiten.NewImage(screenWidth, screenHeight),
		currentColor: color.RGBA{},

		projectionMatrix: matrix.Ident4(),
		viewMatrix:       matrix.Ident4(),
		viewportMatrix:   getViewport(screenWidth, screenHeight),
		modelMatrixStack: []matrix.Mat4{matrix.Ident4()},
	}

	game.scene = makeScene(&game)

	return &game
}

func (g *Game) Update() error {
	g.scene.Update()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debugMode = !g.debugMode
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.scene.active = !g.scene.active
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	g.canvas.Clear()
	g.viewMatrix = matrix.Ident4() // At 0,0, looking in

	right := fovToWidth(fov, math.Abs(near))
	top := right * float64(screenHeight) / float64(screenWidth)

	g.projectionMatrix = viewFrustum(-right, right, -top, top, near, far)
	g.scene.Draw()
	screen.DrawImage(g.canvas, nil)
}

func (g *Game) SetColor(color color.RGBA) {
	g.currentColor = color
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
		matrix.Vec4{0, 0, s, 0},
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

	ndc_corrected := ndc.Mul(1.0 / ndc.W())

	screen := g.viewportMatrix.Mul4x1(ndc_corrected)

	if g.debugMode {
		fmt.Println("world    : ", world)
		fmt.Println("model    : ", model)
		fmt.Println("camera   : ", camera)
		fmt.Println("ndc      : ", ndc)
		fmt.Println("ndc W-Cor: ", ndc_corrected)
		fmt.Println("screen   : ", screen)
	}

	return screen.Vec2()
}

func (g *Game) DrawPixel(x, y int) {
	g.canvas.Set(x, y, g.currentColor)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("3D Starfield")

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
