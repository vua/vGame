package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/vua/vGame/component"
	"image/color"
	"log"
	"math/rand"
	"time"
)

// Game implements ebiten.Game interface.
type Game struct {
	bgc    color.NRGBA
	recv   *component.Square
	ball   *component.Square
	awards []*component.Square
	cls    chan int
	w      int
	h      int
}

// Update proceeds the game state.
// Update is called every tick (1/60 [s] by default).
var a byte = 0xff

func (g *Game) Update() error {
	// Write your game's logical update.

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.ball.IsRun = true
	}
	if g.ball.IsRun {
		g.ball.CollisionDetection(g.recv, float64(g.w), float64(g.h))
		g.ball.HitDetection(&g.awards)
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.recv.Move(g.w, -5, g.ball)

	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.recv.Move(g.w, 5, g.ball)
	}
	/*if !g.ball.IsAlive() {
		close(g.cls)
		return errors.New("Game Over")
	}*/
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
func (g *Game) Draw(screen *ebiten.Image) {

	screen.DrawImage(g.recv.Image, g.recv.Opts)
	screen.DrawImage(g.ball.Image, g.ball.Opts)
	for _, i := range g.awards {
		screen.DrawImage(i.Image, i.Opts)
	}
	if !g.ball.IsAlive() {
		ebitenutil.DebugPrint(screen, "Game Over")

	} else {
		ebitenutil.DebugPrint(screen, "Score:"+g.ball.GetScore())
	}
	// Write your game's rendering.
}

// Layout takes the outside size (e.g., the window size) and returns the (logical) screen size.
// If you don't have to adjust the screen size with the outside size, just return a fixed size.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.w, g.h
}

func (g *Game) awardGenerator() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			if !g.ball.IsAlive() {
				return
			}
			if len(g.awards) == 5 {
				g.awards = g.awards[1:5]
			}
			g.awards = append(g.awards, component.NewSquare(Yellow, 5, 5, float64(rand.Intn(300)+10), float64(rand.Intn(200)+10), 0))
		case <-g.cls:
			return
		}
	}
}

var (
	Red    = color.RGBA{0xff, 0x00, 0x00, 0xff}
	Yellow = color.RGBA{0xff, 0xff, 0x00, 0xff}
	White  = color.RGBA{0xff, 0xff, 0xff, 0xff}
)

func main() {
	game := &Game{
		recv:   component.NewSquare(White, 5, 40, 140, 230, 0),
		ball:   component.NewSquare(Red, 5, 5, 157.5, 225, 3),
		awards: make([]*component.Square, 0),
		cls:    make(chan int, 0),
		w:      320,
		h:      240,
	}
	go game.awardGenerator()
	// Specify the window size as you like. Here, a doubled size is specified.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("vGame")
	// Call ebiten.RunGame to start your game loop.
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
