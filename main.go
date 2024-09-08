package main

import (
	"fmt"
	"image/color"
	_ "image/png"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const N = 5000
const SCREEN_WIDTH = 1024
const SCREEN_HEIGHT = 1024
const PARTICULE_SIZE = 2
const PARTICULE_SPEED = 2
const SENSOR_DISTANCE = PARTICULE_SPEED
const SENSOR_ANGLE = 20
const FADING_FACTOR = 2

type Particule struct {
	Pos rl.Vector2
	Vel rl.Vector2
}

func (p *Particule) Move(image *rl.Image) {
	// Avoid edges
	if p.Pos.X+p.Vel.X < 40 ||
		p.Pos.X+p.Vel.X > SCREEN_WIDTH-40 ||
		p.Pos.Y+p.Vel.Y < 40 ||
		p.Pos.Y+p.Vel.Y > SCREEN_HEIGHT-40 {
		p.Vel = rl.Vector2Rotate(p.Vel, 120)
	}

	// Determine next direction
	sampleCenter := rl.GetImageColor(
		*image,
		int32(p.Pos.X+p.Vel.X),
		int32(p.Pos.Y+p.Vel.Y),
	).A
	velLeft := rl.Vector2Rotate(p.Vel, -SENSOR_ANGLE)
	sampleLeft := rl.GetImageColor(
		*image,
		int32(p.Pos.X+velLeft.X),
		int32(p.Pos.Y+velLeft.Y),
	).A
	velRight := rl.Vector2Rotate(p.Vel, SENSOR_ANGLE)
	sampleRight := rl.GetImageColor(
		*image,
		int32(p.Pos.X+velRight.X),
		int32(p.Pos.Y+velRight.Y),
	).A
	if sampleLeft > sampleCenter && sampleLeft > sampleRight {
		p.Vel.X = velLeft.X
		p.Vel.Y = velLeft.Y
	} else if sampleRight > sampleCenter && sampleRight > sampleLeft {
		p.Vel.X = velRight.X
		p.Vel.Y = velRight.Y
	}

	// Updated position
	p.Pos.X += p.Vel.X
	p.Pos.Y += p.Vel.Y
}

func NewParticule() *Particule {
	dist := rand.Float32() * float32(SCREEN_WIDTH) * 0.25
	angle := rand.Float32() * 360
	pos := rl.Vector2Rotate(rl.Vector2{X: dist, Y: 0}, angle)
	return &Particule{
		Pos: rl.Vector2{
			X: pos.X + float32(SCREEN_WIDTH)*0.5,
			Y: pos.Y + float32(SCREEN_HEIGHT)*0.5,
		},
		Vel: rl.Vector2Rotate(rl.Vector2{X: PARTICULE_SPEED, Y: 0}, angle),
	}
}

func main() {
	// Init rendering window
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), "Physarium Transport Network")
	defer rl.CloseWindow()

	rl.SetTargetFPS(30)

	// Init texture
	target := rl.LoadRenderTexture(SCREEN_WIDTH, SCREEN_HEIGHT)

	// Init Particles
	var particules [N]Particule
	for i := 0; i < N; i++ {
		particules[i] = *NewParticule()
	}

	// Load shader
	shader := rl.LoadShader("", "diffuse.fs")
	fmt.Println(shader)
	rectange := rl.NewRectangle(
		0,
		0,
		float32(SCREEN_WIDTH),
		float32(SCREEN_HEIGHT),
	)

	// Rendering loop
	for !rl.WindowShouldClose() {
		// Get latest texture data from the GPU
		image := rl.LoadImageFromTexture(target.Texture)

		// Simulate the movement of every particule
		for i := 0; i < len(particules); i++ {
			particules[i].Move(image)
		}

		rl.BeginDrawing()

		rl.BeginTextureMode(target)

		// Fade the existing trails
		rl.DrawRectangleRec(rectange, color.RGBA{0, 0, 0, FADING_FACTOR})

		// Add trail for the latest move
		for _, p := range particules {
			rl.DrawPixelV(p.Pos, rl.White)
		}

		// Diffuse the trails with a gaussian blur
		rl.BeginShaderMode(shader)
		rl.DrawTexture(target.Texture, 0, 0, rl.White)
		rl.EndShaderMode()

		rl.EndTextureMode()

		// image = rl.LoadImageFromTexture(target.Texture)
		// rl.ImageBlurGaussian(image, 10)
		// colors := rl.LoadImageColors(image)
		// rl.UpdateTexture(target.Texture, colors)

		rl.ClearBackground(rl.Black)
		rl.DrawTexture(target.Texture, 0, 0, rl.White)

		rl.DrawFPS(10, 10)
		rl.EndDrawing()

	}
}
