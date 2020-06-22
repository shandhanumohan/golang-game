package main

import (
	"fmt"
	"image/color"
	"log"

	"./model"
	"./resolv"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

type Game struct {
	player *model.Player
	alien  *model.Alien
	hits   int
}

func NewGame() *Game {
	return &Game{
		hits: 0,
		player: &model.Player{
			XPos:      screenWidth / 2.0,
			YPos:      screenHeight / 2.0,
			Speed:     4,
			Shoot:     false,
			FireCount: 0,
		},
		alien: &model.Alien{
			XPos: alienDefaultX,
			YPos: alienDefaultY,
		},
	}
}

var (
	background      *ebiten.Image
	spaceShip       *ebiten.Image
	shootFire       *ebiten.Image
	explode         *ebiten.Image
	err             error
	moon            *ebiten.Image
	mplusNormalFont font.Face
)

const (
	screenWidth, screenHeight = 640, 480
	dpi                       = 72
	x                         = 20
	missileSpeed              = 20
	alienDefaultX             = 0
	alienDefaultY             = 10
	alienSpeed                = 10
	alienWidth                = 100
	missileWidth              = 100
)

//
func (g *Game) Update(screen *ebiten.Image) error {
	// Write your game's logical update.

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.player.YPos -= g.player.Speed

	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.player.YPos += g.player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.XPos -= g.player.Speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.XPos += g.player.Speed
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {

		g.player.FireCount = g.player.FireCount + 1

		m := model.Missile{XPos: g.player.XPos - 75, YPos: g.player.YPos - 150, Visible: true}
		g.player.MissilesFired = append(g.player.MissilesFired, m)

	}

	return nil
}

func (g *Game) reverseFire() {
	g.player.Shoot = false
}

func (g *Game) moveMissiles() {
	for i, _ := range g.player.MissilesFired {
		g.player.MissilesFired[i].YPos -= missileSpeed

		if g.player.MissilesFired[i].YPos < 0 {
			g.player.MissilesFired[i].Visible = false

		}
	}
}

func (g *Game) moveAlien() {
	g.alien.XPos += alienSpeed
	if g.alien.XPos > screenWidth {
		g.alien.XPos = alienDefaultX
	}
}

func init() {
	background, _, err = ebitenutil.NewImageFromFile("assets/space.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	spaceShip, _, err = ebitenutil.NewImageFromFile("assets/spaceship.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	shootFire, _, err = ebitenutil.NewImageFromFile("assets/fire.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	moon, _, err = ebitenutil.NewImageFromFile("assets/enemy.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	explode, _, err = ebitenutil.NewImageFromFile("assets/explode.png", ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}

	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	mplusNormalFont = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

//
func (g *Game) Draw(screen *ebiten.Image) {
	// Write your game's rendering.

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, 0)
	screen.DrawImage(background, op)

	moonOp := &ebiten.DrawImageOptions{}
	moonOp.GeoM.Translate(g.alien.XPos, g.alien.YPos)
	screen.DrawImage(moon, moonOp)

	alienShape := resolv.NewRectangle(int32(g.alien.XPos), int32(g.alien.YPos), alienWidth, alienWidth)

	g.moveAlien()

	//fmt.Println(len(g.player.MissilesFired))
	if len(g.player.MissilesFired) > 0 {

		for _, missile := range g.player.MissilesFired {
			if missile.Visible {
				shootOp := &ebiten.DrawImageOptions{}
				shootOp.GeoM.Translate(missile.XPos, missile.YPos)
				screen.DrawImage(shootFire, shootOp)

				missileShape := resolv.NewRectangle(int32(missile.XPos), int32(missile.YPos), missileWidth, missileWidth)

				colliding := alienShape.IsColliding(missileShape)
				if colliding {
					explodeOp := &ebiten.DrawImageOptions{}
					explodeOp.GeoM.Translate(g.alien.XPos, g.alien.YPos)
					screen.DrawImage(explode, explodeOp)

					g.hits++
					g.alien.XPos = alienDefaultX
				}

			}

		}
		g.moveMissiles()
		//g.reverseFire()

	}
	playerOp := &ebiten.DrawImageOptions{}

	playerOp.GeoM.Translate(g.player.XPos, g.player.YPos)
	screen.DrawImage(spaceShip, playerOp)

	//msg := fmt.Sprintf("FireCount: %d", g.player.FireCount)
	score := fmt.Sprintf("Score: %d", g.hits)
	//text.Draw(screen, msg, mplusNormalFont, x, 40, color.White)
	text.Draw(screen, score, mplusNormalFont, x, 40, color.White)
}

//
func (g *Game) Layout(outsideWidth, outsideHeight int) (width, height int) {
	return outsideWidth, outsideHeight
}

func main() {

	// Sepcify the window size as you like. Here, a doulbed size is specified.
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Your game's title")
	// Call ebiten.RunGame to start your game loop.

	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}
