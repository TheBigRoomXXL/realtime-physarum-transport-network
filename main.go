package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const N = 50
const SCREEN_WIDTH = 1024
const SCREEN_HEIGHT = 768
const PARTICULE_SIZE = 5

type Particule struct {
	Pos rl.Vector2
	Vel rl.Vector2
}

func (p *Particule) Draw() {
	direction := rl.Vector2{
		X: p.Pos.X + p.Vel.X*PARTICULE_SIZE*2,
		Y: p.Pos.Y + p.Vel.Y*PARTICULE_SIZE*2,
	}
	rl.DrawCircleV(p.Pos, PARTICULE_SIZE, rl.Blue)
	rl.DrawLineV(p.Pos, direction, rl.Blue)
}

func (p *Particule) Move(grid *Grid) {
	p.Pos.X += p.Vel.X
	p.Pos.Y += p.Vel.Y
	grid.Add(p.Pos, 0.5)
}

// Adapted form https://github.com/fogleman/physarum/blob/main/pkg/physarum/grid.go
type Grid struct {
	W, H    int
	Image   *rl.Image
	Texture rl.Texture2D
}

func NewGrid(w, h int) *Grid {
	imgImg := image.NewGray(
		image.Rectangle{
			image.Point{0, 0},
			image.Point{SCREEN_WIDTH, SCREEN_HEIGHT},
		},
	)
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			imgImg.SetGray(i, j, color.Gray{0})

		}
	}
	img := rl.NewImageFromImage(imgImg)
	texture := rl.LoadTextureFromImage(img)
	return &Grid{w, h, img, texture}
}

func (g *Grid) Add(pos rl.Vector2, a float32) {
	rl.ImageDrawCircleV(
		g.Image, pos, PARTICULE_SIZE, color.RGBA{255, 255, 255, 25},
	)
}

func (g *Grid) Draw() {
	dataSlice := (*[1 << 24]rl.Color)(g.Image.Data)[:g.Image.Width*g.Image.Height]
	rl.UpdateTexture(g.Texture, dataSlice)
	rl.DrawTexture(g.Texture, 0, 0, color.RGBA{255, 255, 255, 255})
}

func NewParticule() *Particule {
	x := rand.Float32()
	y := rand.Float32()
	return &Particule{
		Pos: rl.Vector2{
			X: x*float32(SCREEN_WIDTH)*0.20 + 0.40*float32(SCREEN_WIDTH),
			Y: y*float32(SCREEN_HEIGHT)*0.20 + 0.40*float32(SCREEN_HEIGHT),
		},
		Vel: rl.Vector2{
			X: rand.Float32() - 0.5,
			Y: rand.Float32() - 0.5,
		},
	}
}

func main() {
	// Init rendering window
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), "Boids")
	defer rl.CloseWindow()

	rl.SetTargetFPS(30)

	// Init Particles
	var particules [N]Particule
	for i := 0; i < N; i++ {
		particules[i] = *NewParticule()
	}

	// Init Grid
	grid := NewGrid(SCREEN_WIDTH, SCREEN_HEIGHT)

	// Rendering loop
	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(color.RGBA{22, 23, 31, 100})

		for i := 0; i < len(particules); i++ {
			particules[i].Move(grid)
			particules[i].Draw()
		}
		grid.Draw()

		rl.DrawText(fmt.Sprint("FPS: ", rl.GetFPS()), 20, 20, 20, rl.LightGray)
		rl.EndDrawing()
	}
}
