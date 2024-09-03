package main

import (
	"fmt"
	"image/color"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const N = 50
const PARTICULE_SIZE = 5

var screenHeigh = 768
var screenWidth = 1024

type Particule struct {
	Pos rl.Vector2
	Vel rl.Vector2
}

func (p *Particule) Draw() {
	direction := rl.Vector2{X: p.Pos.X + p.Vel.X*2, Y: p.Pos.Y + p.Vel.Y*2}
	rl.DrawCircleV(p.Pos, PARTICULE_SIZE, rl.Blue)
	rl.DrawLineV(p.Pos, direction, rl.Blue)
}

func (p *Particule) Move() {
	p.Pos.X += p.Vel.X
	p.Pos.Y += p.Vel.Y
}

func NewParticule() *Particule {
	x := rand.Float32()
	y := rand.Float32()
	return &Particule{
		Pos: rl.Vector2{
			X: x*float32(screenWidth)*0.20 + 0.40*float32(screenWidth),
			Y: y*float32(screenHeigh)*0.20 + 0.40*float32(screenHeigh),
		},
		Vel: rl.Vector2{
			X: rand.Float32() - 0.5,
			Y: rand.Float32() - 0.5,
		},
	}
}

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(int32(screenWidth), int32(screenHeigh), "Boids")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var particules [N]Particule
	for i := 0; i < N; i++ {
		particules[i] = *NewParticule()
	}

	for !rl.WindowShouldClose() {
		if rl.IsWindowResized() {
			screenWidth = rl.GetRenderWidth()
			screenHeigh = rl.GetRenderHeight()
		}

		rl.BeginDrawing()

		rl.ClearBackground(color.RGBA{22, 23, 31, 255})
		for i := 0; i < len(particules); i++ {

			particules[i].Move()
			particules[i].Draw()
		}

		rl.DrawText(fmt.Sprint("FPS: ", rl.GetFPS()), 20, 20, 20, rl.LightGray)

		rl.EndDrawing()
	}
}
