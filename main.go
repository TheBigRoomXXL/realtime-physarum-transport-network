package main

import (
	"image"
	"image/color"
	_ "image/png"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const N = 10_000
const SCREEN_WIDTH = 1024
const SCREEN_HEIGHT = 1024
const PARTICULE_SIZE = 1
const PARTICULE_SPEED = 4
const SENSOR_DISTANCE = PARTICULE_SPEED
const SENSOR_ANGLE = 45
const DIFFUSION_SIZE = 2
const TRAIL_STRENGTH = 255

type Particule struct {
	Pos rl.Vector2
	Vel rl.Vector2
}

func (p *Particule) Move(image *rl.Image) {
	// Avoid edges
	if p.Pos.X+p.Vel.X < 10 ||
		p.Pos.X+p.Vel.X > SCREEN_WIDTH-10 ||
		p.Pos.Y+p.Vel.Y < 10 ||
		p.Pos.Y+p.Vel.Y > SCREEN_HEIGHT-10 {
		p.Vel = rl.Vector2Rotate(p.Vel, 90)
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
		Vel: rl.Vector2Rotate(rl.Vector2{X: -PARTICULE_SPEED, Y: 0}, angle),
	}
}

func main() {
	// Init rendering window
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.SetConfigFlags(rl.FlagWindowResizable)

	rl.InitWindow(int32(SCREEN_WIDTH), int32(SCREEN_HEIGHT), "Physarium Transport Network")
	defer rl.CloseWindow()

	rl.SetTargetFPS(30)

	//  Init texture in CPU and GPU
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
	image := rl.NewImageFromImage(imgImg)
	target := rl.LoadRenderTexture(SCREEN_WIDTH, SCREEN_HEIGHT)

	// Init Particles
	var particules [N]Particule
	for i := 0; i < N; i++ {
		particules[i] = *NewParticule()
	}

	// Load shader
	shader := rl.LoadShader("", "diffuse.fs")
	if !rl.IsShaderReady(shader) {
		panic(shader)
	}
	rectange := rl.NewRectangle(
		0,
		0,
		float32(SCREEN_WIDTH),
		float32(SCREEN_HEIGHT),
	)

	// Rendering loop
	// i := 0
	for !rl.WindowShouldClose() {
		// Debug
		// i++
		// if i%30 == 0 {
		// 	rl.ExportImage(*image, fmt.Sprintf("frames/%d.png", i/30))
		// }

		// Update the image data with the particule simulation
		for i := 0; i < len(particules); i++ {
			particules[i].Move(image)
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		rl.BeginTextureMode(target)
		rl.DrawRectangleRec(rectange, color.RGBA{0, 0, 0, 10})
		for _, p := range particules {
			rl.DrawCircleV(p.Pos, 2, rl.White)
		}
		rl.EndTextureMode()

		rl.BeginShaderMode(shader)
		rl.DrawTextureRec(target.Texture, rectange, rl.Vector2Zero(), rl.White)
		rl.EndShaderMode()

		rl.DrawFPS(10, 10)
		rl.EndDrawing()

		// Get texture data after shader processing
		image = rl.LoadImageFromTexture(target.Texture)
	}
}
