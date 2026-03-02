package main

import (
    "fmt"
    "github.com/hajimehoshi/ebiten/v2"
)

type Game struct{}
func (g *Game) Update() error { return nil }
func (g *Game) Draw(screen *ebiten.Image) {}
func (g *Game) Layout(w, h int) (int, int) { return w, h }
func (g *Game) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
    op := &ebiten.DrawImageOptions{}
    op.Filter = ebiten.FilterLinear
    op.GeoM = geoM
    screen.DrawImage(offscreen, op)
}
func main() {
    g := &Game{}
    _ = ebiten.RunGame(g)
}
