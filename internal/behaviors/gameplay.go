package behaviors

import (
	"fmt"
	"image/color"
	"math"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
	etext "github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// gameplay room behaviors
// go (wave controller), Control (HUD init), Draw (HUD render),
// home_Button, BloonSpawn, End, Wave_Panel, track controllers,
// auto_start_button, upgrade panels

// go — the wave start/speed button. Clicking starts next wave.
// sprite "sprite278" when idle, "Going" when wave is active.
// draw event renders speed indicator bars via scr_Fast_Forward logic.
type GoBehavior struct {
	engine.DefaultBehavior
	shiftpress        int
	afterwave         int
	autoButtonChecked bool
}

func (b *GoBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["autostart"] = 0.0
	g.GlobalVars["cashinflate"] = 0.0
	g.GlobalVars["freeplay"] = 0.0
	g.GlobalVars["cashflow"] = 1.0
	g.GlobalVars["cashwavereward"] = 0.0
	g.GlobalVars["bpower"] = 1.0 + (getGlobal(g, "strongerbloons") / 4.0)
	g.GlobalVars["bspeed"] = 1.0 + (getGlobal(g, "fasterbloons") / 4.0)
	g.GlobalVars["wave"] = 1.0
	g.GlobalVars["wavenow"] = 0.0
	g.GlobalVars["cycle"] = 0.0
	g.GlobalVars["money"] = 750.0
	g.GlobalVars["endsequence"] = 0.0
	g.GlobalVars["life"] = 200.0 - (199.0 * getGlobal(g, "noliveslost"))
	g.GlobalVars["points"] = 0.0
	g.GlobalVars["gamespeed"] = 30.0
	g.SetGameSpeed(30) // ensure actual game tick rate matches the variable
	b.shiftpress = 0
	b.afterwave = -1
	b.autoButtonChecked = false
}

func (b *GoBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// iMPORTANT: Do this in Step (not Create) to avoid lock re-entry deadlock
	// when room instances are being created.
	if !b.autoButtonChecked {
		b.autoButtonChecked = true
		if g.InstanceMgr.InstanceCount("auto_start_button") == 0 {
			g.InstanceMgr.Create("auto_start_button", 944, 544)
		}
	}

	wavenow := getGlobal(g, "wavenow")
	wave := getGlobal(g, "wave")
	bloonCount := float64(countBloons(g))

	// when timeline finishes spawning all bloons, mark wavenow=0
	if wavenow == 1 && g.ActiveTimeline != nil && !g.ActiveTimeline.Running {
		g.GlobalVars["wavenow"] = 0.0
		wavenow = 0
	}

	// switch sprite based on wave state
	if wavenow == 1 || bloonCount > 0 {
		inst.SpriteName = "Going"
	} else {
		inst.SpriteName = "sprite278"
	}

	// wave-end logic: all bloons popped, wave was active
	if wavenow == 0 && bloonCount == 0 && b.afterwave == 0 {
		// reset sprite to idle
		inst.SpriteName = "sprite278"

		// cash reward
		cashReward := getGlobal(g, "cashwavereward")
		cashFlow := getGlobal(g, "cashflow")
		cashInflate := getGlobal(g, "cashinflate")
		money := getGlobal(g, "money")
		wealthiness := getGlobal(g, "wealthiness")

		reward := math.Round(cashReward * cashFlow * (1.0 + cashInflate*0.1))
		money += reward
		money += wealthiness*20.0 - 1.0 + wave

		// points
		pm := getGlobal(g, "pointmultiplier")
		pts := getGlobal(g, "points")
		pts += (100.0 + wave*3.0) * pm
		pts += math.Sqrt(cashReward) * pm
		if pts > 0 {
			pts += math.Sqrt(money) * pm
		}

		// wave 91 bonus
		if wave == 91 {
			pts += 5000 * pm
			g.GlobalVars["XP"] = getGlobal(g, "XP") + pts/2.0
		}

		// cash inflation bonus
		if cashInflate > 0 {
			money = math.Round(money * (1.0 + 0.03*cashInflate))
		}

		// no lives lost decay
		if getGlobal(g, "noliveslost") == 1 {
			life := getGlobal(g, "life")
			if life > 1 {
				g.GlobalVars["life"] = math.Floor(life * 0.92)
			}
		}

		g.GlobalVars["money"] = money
		g.GlobalVars["points"] = pts
		g.GlobalVars["cashwavereward"] = 0.0
		b.afterwave = 1
	}

	// auto-start: when autostart=1 and between waves, auto-trigger next wave
	if getGlobal(g, "autostart") == 1 && bloonCount == 0 && wavenow == 0 && b.afterwave == 1 {
		b.startNextWave(inst, g)
	}
}

// mouseLeftPressed — start next wave or toggle speed
func (b *GoBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	b.handleWaveStart(inst, g)
}

// keyPress — space bar also starts waves / toggles speed
func (b *GoBehavior) KeyPress(inst *engine.Instance, g *engine.Game) {
	if g.InputMgr.KeyPressed(ebiten.KeySpace) || g.InputMgr.KeyPressed(ebiten.KeyEnter) {
		b.handleWaveStart(inst, g)
	}
}

func (b *GoBehavior) handleWaveStart(inst *engine.Instance, g *engine.Game) {
	wavenow := getGlobal(g, "wavenow")
	bloonCount := float64(countBloons(g))

	// if wave is active and bloons exist — toggle speed
	if wavenow == 1 || bloonCount > 0 {
		// cycle: normal(30) → fast(60) → superfast(90) → normal(30)
		switch b.shiftpress {
		case 0:
			// currently at normal → go to fast
			g.SetGameSpeed(60)
			g.GlobalVars["gamespeed"] = 60.0
			b.shiftpress = 1
		case 1:
			// currently at fast → go to superfast
			g.SetGameSpeed(90)
			g.GlobalVars["gamespeed"] = 90.0
			b.shiftpress = 2
		case 2:
			// currently at superfast → go to normal
			g.SetGameSpeed(30)
			g.GlobalVars["gamespeed"] = 30.0
			b.shiftpress = 0
		}
		return
	}

	// if no bloons and wave not active — start next wave
	if bloonCount == 0 && wavenow == 0 {
		b.startNextWave(inst, g)
	}
}

func (b *GoBehavior) startNextWave(inst *engine.Instance, g *engine.Game) {
	wave := int(getGlobal(g, "wave"))
	if wave > 90 {
		return // past last wave
	}

	// switch to active sprite
	inst.SpriteName = "Going"

	// dispatch timeline for this wave
	tlName := fmt.Sprintf("N%d", wave)
	tl := g.TimelineMgr.Get(tlName)
	if tl != nil {
		wsqueeze := getGlobal(g, "wavesqueeze")
		runner := engine.NewTimelineRunner(tl)
		runner.Speed = 1.0 + wsqueeze
		runner.OnAction = func(action engine.TimelineAction) {
			executeBloonSpawn(g, action)
		}
		g.ActiveTimeline = runner
		fmt.Printf("Wave %d started (timeline %s, %d steps)\n", wave, tlName, tl.MaxStep)
	} else {
		fmt.Printf("WARNING: Timeline %s not found for wave %d\n", tlName, wave)
		inst.SpriteName = "sprite278"
		return
	}

	g.GlobalVars["wavenow"] = 1.0
	g.GlobalVars["cycle"] = math.Mod(getGlobal(g, "cycle")+1, 4)
	if getGlobal(g, "cycle") == 0 {
		g.GlobalVars["cycle"] = 1.0
	}
	b.afterwave = 0
	g.GlobalVars["wave"] = float64(wave + 1)
}

// draw renders the Go button sprite plus speed indicator bars
func (b *GoBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr != nil && len(spr.Frames) > 0 {
		engine.DrawSpriteExt(screen, spr.Frames[0], spr.XOrigin, spr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
	}

	// speed indicator bars (left and right vertical bars on the button)
	// fill percentage based on room_speed / 0.9 (capped at 100)
	gs := getGlobal(g, "gamespeed")
	fill := gs / 90.0 // 0..1: 30→0.33, 60→0.67, 90→1.0
	if fill > 1 {
		fill = 1
	}

	barH := 44.0 // total bar height
	filledH := barH * fill
	barTop := inst.Y + 10
	barBottom := inst.Y + 54

	// left bar at x+8
	drawRect(screen, inst.X+8, barBottom-filledH, 2, filledH, speedBarColor(fill))
	drawRect(screen, inst.X+8, barTop, 2, barH-filledH, color0x000000)

	// right bar at x+53
	drawRect(screen, inst.X+53, barBottom-filledH, 2, filledH, speedBarColor(fill))
	drawRect(screen, inst.X+53, barTop, 2, barH-filledH, color0x000000)
}

var color0x000000 = [3]uint8{0, 0, 0}

func speedBarColor(fill float64) [3]uint8 {
	// interpolate yellow to red based on fill
	if fill < 0.5 {
		return [3]uint8{255, 255, 0} // yellow
	}
	return [3]uint8{255, uint8(255 * (1 - fill) * 2), 0} // yellow→red
}

func drawRect(screen *ebiten.Image, x, y, w, h float64, col [3]uint8) {
	if w <= 0 || h <= 0 {
		return
	}
	vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), color.RGBA{col[0], col[1], col[2], 255}, false)
}

// auto_start_button — toggles autostart on/off
// frame 0 = off, frame 1 = on
type AutoStartButton struct {
	engine.DefaultBehavior
}

func (b *AutoStartButton) Create(inst *engine.Instance, g *engine.Game) {
	inst.Depth = -1000
}

func (b *AutoStartButton) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "autostart") == 0 {
		g.GlobalVars["autostart"] = 1.0
	} else {
		g.GlobalVars["autostart"] = 0.0
	}
}

func (b *AutoStartButton) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite("Auto_Start_butt_spr")
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := int(getGlobal(g, "autostart"))
	if frame >= len(spr.Frames) {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// settings — creates wave panel cards used in the side HUD.
// spawns 90 Wave_Panel instances.
type SettingsBehavior struct {
	engine.DefaultBehavior
}

func (b *SettingsBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["sandbox"] = 0.0
	if g.InstanceMgr.InstanceCount("Wave_Panel") > 0 {
		return
	}
	for i := 1; i <= 90; i++ {
		g.InstanceMgr.Create("Wave_Panel", inst.X, inst.Y+float64(i*64))
	}
}

// control — HUD controller, initializes gameplay globals
type ControlBehavior struct {
	engine.DefaultBehavior
	cheatProgress int
}

func (b *ControlBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["stuncycle"] = 0.0
	g.GlobalVars["Wavepanel"] = 0.0
	g.GlobalVars["amplification"] = 0.0
	g.GlobalVars["upgradelock"] = 0.0
	g.GlobalVars["upgradeselect"] = 0.0
	g.GlobalVars["spin"] = 0.0
	g.GlobalVars["select"] = 0.0
	g.GlobalVars["monkeypop"] = 0.0
	g.GlobalVars["idset"] = 0.0
	g.GlobalVars["biglayer"] = 0.0
	if _, ok := g.GlobalVars["life"]; !ok {
		g.GlobalVars["life"] = 200.0
	}
	g.GlobalVars["dmgreduction"] = 1.0
	g.GlobalVars["wiredfunds"] = 0.0
	g.GlobalVars["tower"] = 0.0
	g.GlobalVars["towerselect"] = 0.0
	g.GlobalVars["towerplace"] = 0.0
	g.GlobalVars["up"] = 0.0
	g.GlobalVars["pathup"] = 0.0
	g.GlobalVars["maxlayer"] = 50.0
	g.GlobalVars["panelsee"] = 1.0
	b.cheatProgress = 0

	// create the sell panel.
	if g.InstanceMgr.InstanceCount("Sell") == 0 {
		g.InstanceMgr.Create("Sell", 64, 416)
	}
}

var debugKonamiSequence = []ebiten.Key{
	ebiten.KeyArrowUp,
	ebiten.KeyArrowUp,
	ebiten.KeyArrowDown,
	ebiten.KeyArrowDown,
	ebiten.KeyArrowLeft,
	ebiten.KeyArrowRight,
	ebiten.KeyArrowLeft,
	ebiten.KeyArrowRight,
}

func pressedArrowKey(g *engine.Game) (ebiten.Key, bool) {
	keys := []ebiten.Key{
		ebiten.KeyArrowUp,
		ebiten.KeyArrowDown,
		ebiten.KeyArrowLeft,
		ebiten.KeyArrowRight,
	}
	for _, k := range keys {
		if g.InputMgr.KeyPressed(k) {
			return k, true
		}
	}
	return 0, false
}

func activateDebugCheat(g *engine.Game) {
	// money cap in Control.Step is 2,000,000,000.
	g.GlobalVars["money"] = 2000000000.0

	// unlock all tower/unit buy panels for this run.
	for _, cfg := range panelConfigs {
		if cfg.lockKey == "" {
			continue
		}
		g.GlobalVars[cfg.lockKey] = 1.0
	}
	// extra legacy lock used by some old menus/challenges.
	g.GlobalVars["Derlock"] = 1.0

	// unlock all tower tier-4/tier-5 paths for debug.
	for _, triplet := range towerPathProgressVars {
		for _, key := range triplet {
			if key == "" {
				continue
			}
			g.GlobalVars[key] = 99.0
		}
	}
	fmt.Println("DEBUG CHEAT ACTIVATED: max money + all units + all tier4/tier5 paths unlocked")
}

func (b *ControlBehavior) handleDebugCheat(g *engine.Game) {
	key, ok := pressedArrowKey(g)
	if !ok {
		return
	}

	expected := debugKonamiSequence[b.cheatProgress]
	if key == expected {
		b.cheatProgress++
		if b.cheatProgress >= len(debugKonamiSequence) {
			b.cheatProgress = 0
			activateDebugCheat(g)
		}
		return
	}

	if key == debugKonamiSequence[0] {
		b.cheatProgress = 1
		return
	}
	b.cheatProgress = 0
}

func (b *ControlBehavior) Step(inst *engine.Instance, g *engine.Game) {
	b.handleDebugCheat(g)

	// clamp money
	money := getGlobal(g, "money")
	if money < 0 {
		g.GlobalVars["money"] = 0.0
	}
	if money > 2000000000 {
		g.GlobalVars["money"] = 2000000000.0
	}

	// right-click cancels tower placement mode
	if g.InputMgr.MouseRightPressed() {
		if getGlobal(g, "towerplace") == 1 {
			g.GlobalVars["towerplace"] = 0.0
			g.GlobalVars["towerselect"] = 0.0
			for _, block := range g.InstanceMgr.FindByObject("Block") {
				block.Visible = false
				block.SpriteName = ""
			}
		}
	}
}

// drawHUD — the HUD overlay (Draw object, depth -22)
// shows money, lives, rank, points
// black for money/rank/points, red for life
type DrawHUD struct {
	engine.DefaultBehavior
	font *engine.BMFont
}

// colors used in the HUD draw
var hudColorBlack = [3]uint8{0, 0, 0}
var hudColorRed = [3]uint8{255, 0, 0}

func (b *DrawHUD) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	// lazy-load the bitmap font once
	if b.font == nil {
		b.font = g.BMFont
	}

	// draw HUD icons
	drawHUDSprite(screen, g, "Pop_Icon", 0, 878, 59)
	drawHUDSprite(screen, g, "Banana_Icon", 0, 878, 2)
	drawHUDSprite(screen, g, "Heart_Icon", 0, 880, 30)

	// draw Visual_Dart decorations
	drawHUDSprite(screen, g, "Visual_Dart", 0, 895, 462)
	drawHUDSprite(screen, g, "Visual_Dart", 1, 992, 462)

	// money — black text
	money := int(getGlobal(g, "money"))
	drawHUDTextColored(screen, g, b.font, fmt.Sprintf("%d", money), 911, 4, hudColorBlack)

	// life — red text
	life := int(getGlobal(g, "life"))
	drawHUDTextColored(screen, g, b.font, fmt.Sprintf("%d", life), 911, 30, hudColorRed)

	// rank (next to Pop_Icon) — black text
	rank := int(getGlobal(g, "monkeypop"))
	drawHUDTextColored(screen, g, b.font, fmt.Sprintf("%d", rank), 954, 59, hudColorBlack)

	// points — black text
	pts := int(math.Floor(getGlobal(g, "points")))
	drawHUDTextColored(screen, g, b.font, fmt.Sprintf("%d", pts), 942, 451, hudColorBlack)

	// xP health bar
	xp := getGlobal(g, "XP")
	criteria := getGlobal(g, "criteria")
	if criteria > 0 {
		pct := (xp / criteria)
		if pct > 1 {
			pct = 1
		}
		barW := 85.0 // 999-914
		filled := barW * pct
		// background
		drawRect(screen, 914, 82, barW, 4, [3]uint8{0, 0, 0})
		// fill
		if filled > 0 {
			drawRect(screen, 914, 82, filled, 4, [3]uint8{255, 255, 255})
		}
	}

	// draw tower buy tooltip last so it always stays above the panel icons.
	drawHoveredTowerBuyTooltip(screen, g)
}

func drawHUDSprite(screen *ebiten.Image, g *engine.Game, spriteName string, frame int, x, y float64) {
	spr := g.AssetManager.GetSprite(spriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	f := frame % len(spr.Frames)
	engine.DrawSpriteExt(screen, spr.Frames[f], spr.XOrigin, spr.YOrigin,
		x, y, 1, 1, 0, 1)
}

func drawHUDTextColored(screen *ebiten.Image, g *engine.Game, font *engine.BMFont, value string, x, y float64, clr [3]uint8) {
	// use bitmap font if loaded
	if font != nil && len(font.Glyphs) > 0 {
		font.DrawText(screen, value, x, y, clr)
		return
	}
	// fallback: use Go's text rendering with matching color
	c := color.RGBA{clr[0], clr[1], clr[2], 255}
	etext.Draw(screen, value, basicfont.Face7x13, int(x), int(y)+10, c)
}

const hudSmallTextScale = 0.65

// drawHUDTextSmall draws compact crisp UI text for small panels (wave cards, sell value).
func drawHUDTextSmall(screen *ebiten.Image, g *engine.Game, value string, x, y float64, clr [3]uint8) {
	if g != nil && g.BMFont != nil && len(g.BMFont.Glyphs) > 0 {
		g.BMFont.DrawTextScaled(screen, value, math.Round(x), math.Round(y), hudSmallTextScale, clr)
		return
	}
	c := color.RGBA{clr[0], clr[1], clr[2], 255}
	etext.Draw(screen, value, basicfont.Face7x13, int(math.Round(x)), int(math.Round(y))+10, c)
}

// home_Button — exit button near GO.
// shows Exit icon only when no tower is selected/placing.
// click opens Wanna_go_to_main_ confirmation popup.
type HomeButtonBehavior struct {
	engine.DefaultBehavior
	exist int
}

func (b *HomeButtonBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.exist = 0
	inst.Depth = -20
	inst.ImageSpeed = 0
	inst.SpriteName = "sprite277"
	inst.Alarms[0] = 30
}

func (b *HomeButtonBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		b.exist = 1
	}
}

func (b *HomeButtonBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "tower") == 0 && getGlobal(g, "towerplace") == 0 {
		inst.SpriteName = "Exit"
	} else {
		inst.SpriteName = "sprite277"
	}
}

func (b *HomeButtonBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if b.exist != 1 {
		return
	}
	if getGlobal(g, "tower") != 0 || getGlobal(g, "towerplace") != 0 {
		return
	}
	if g.InstanceMgr.InstanceCount("Wanna_go_to_main_") == 0 {
		g.InstanceMgr.Create("Wanna_go_to_main_", 480, 256)
	}
}

// wanna_go_to_main_ confirmation popup
// left click confirms and returns to Main_Menu, right click cancels popup.
type WannaGoToMainBehavior struct {
	engine.DefaultBehavior
}

func (b *WannaGoToMainBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Depth = -100
	inst.ImageSpeed = 0
}

func (b *WannaGoToMainBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.ActiveTimeline = nil
	g.SetGameSpeed(60)
	g.AudioMgr.PlayMusic("Main_Menu0")
	g.RequestRoomGoto("Main_Menu")
}

func (b *WannaGoToMainBehavior) MouseRightPressed(inst *engine.Instance, g *engine.Game) {
	g.InstanceMgr.Destroy(inst.ID)
}

func (b *WannaGoToMainBehavior) KeyPress(inst *engine.Instance, g *engine.Game) {
	if g.InputMgr.KeyPressed(ebiten.KeyEscape) || g.InputMgr.KeyPressed(ebiten.KeyX) {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// bloonSpawn — invisible spawn point, tracks path cycling
type BloonSpawnBehavior struct {
	engine.DefaultBehavior
}

func (b *BloonSpawnBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["path"] = 0.0
	inst.Vars["jump"] = 0.0
}

// end — bloon leak point, subtracts lives on collision
type EndBehavior struct {
	engine.DefaultBehavior
}

func (b *EndBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["time"] = 0.0
}

func (b *EndBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// check for bloons reaching the end
	life := getGlobal(g, "life")
	if life <= 0 {
		// game over - could create Failure popup
		return
	}

	// check collision with bloons at the end point
	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		pp := 0.0
		if v, ok := bloon.Vars["path_progress"]; ok {
			pp = v.(float64)
		}
		// bloon reached end of path
		if pp >= 0.99 {
			layer := 1.0
			if v, ok := bloon.Vars["bloonlayer"]; ok {
				layer = v.(float64)
			}
			ouch := bloonLayerDamage(layer)
			dmgReduction := getGlobal(g, "dmgreduction")
			if dmgReduction <= 0 {
				dmgReduction = 1
			}
			ouch = ouch / dmgReduction

			g.GlobalVars["life"] = getGlobal(g, "life") - ouch
			pts := getGlobal(g, "points")
			pts -= (pts * ouch) / 400.0
			g.GlobalVars["points"] = pts

			g.InstanceMgr.Destroy(bloon.ID)

			if getGlobal(g, "life") <= 0 {
				// game over
				fmt.Println("GAME OVER! Lives reached 0")
			}
		}
	}
}

// bloonLayerDamage maps bloon layer to leak damage
func bloonLayerDamage(layer float64) float64 {
	switch {
	case layer <= 1:
		return 1
	case layer <= 1.5:
		return 4
	case layer <= 2:
		return 2
	case layer <= 2.5:
		return 7
	case layer <= 3:
		return 3
	case layer <= 3.5:
		return 10
	case layer <= 4:
		return 4
	case layer <= 4.5:
		return 13
	case layer <= 5:
		return 5
	case layer <= 5.5:
		return 16
	case layer <= 6:
		return 11
	case layer <= 6.1:
		return 11
	case layer <= 7:
		return 23
	case layer <= 8:
		return 47
	case layer <= 8.5:
		return 142
	default:
		return math.Floor(layer * 6)
	}
}

// wave_Panel — shows current wave number in the HUD sidebar
type WavePanelBehavior struct {
	engine.DefaultBehavior
	waveup float64
}

func (b *WavePanelBehavior) Create(inst *engine.Instance, g *engine.Game) {
	wp := getGlobal(g, "Wavepanel") + 1
	g.GlobalVars["Wavepanel"] = wp
	b.waveup = wp
	inst.Vars["preview_sprite"] = wavePreviewSprite(g, int(wp))
}

func (b *WavePanelBehavior) Step(inst *engine.Instance, g *engine.Game) {
	wave := int(getGlobal(g, "wave"))
	waveup := int(b.waveup)
	displayBase := wave

	// during an active wave, global.wave already points to next wave.
	// show the currently active wave at the first slot.
	if getGlobal(g, "wavenow") == 1 || countBloons(g) > 0 {
		displayBase = wave - 1
	}
	if displayBase < 1 {
		displayBase = 1
	}

	// keep a readable stacked sidebar: base/current + upcoming waves.
	offset := waveup - displayBase
	if offset >= 0 && offset < 8 {
		inst.Y = 64 + float64(offset*64)
		inst.Depth = -19
		return
	}

	// move non-visible cards off-screen and behind gameplay.
	inst.Y = 1200
	inst.Depth = 200
}

func (b *WavePanelBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	modeFrame := 0
	if getGlobal(g, "impoppablemodeselect") == 1 {
		modeFrame = 1
	}
	if getGlobal(g, "nightmaremodeselect") == 1 {
		modeFrame = 2
	}

	panelSpr := g.AssetManager.GetSprite("Wave_Paper_spr")
	if panelSpr != nil && len(panelSpr.Frames) > 0 {
		frame := modeFrame
		if frame >= len(panelSpr.Frames) {
			frame = 0
		}
		engine.DrawSpriteExt(screen, panelSpr.Frames[frame], panelSpr.XOrigin, panelSpr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
	}

	drawHUDTextSmall(screen, g, fmt.Sprintf("%d", int(b.waveup)), inst.X+7, inst.Y+3, hudColorBlack)

	preview, _ := inst.Vars["preview_sprite"].(string)
	if preview == "" {
		preview = "Red_Bloon_Spr"
	}
	spr := g.AssetManager.GetSprite(preview)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	engine.DrawSpriteExt(screen, spr.Frames[0], spr.XOrigin, spr.YOrigin,
		inst.X+31, inst.Y+32, 1, 1, 0, 1)
}

func wavePreviewSprite(g *engine.Game, wave int) string {
	tl := g.TimelineMgr.Get(fmt.Sprintf("N%d", wave))
	if tl == nil {
		return "Red_Bloon_Spr"
	}

	maxLayer := 1.0
	found := false
	for _, step := range tl.Steps {
		for _, act := range step.Actions {
			layer, maxLayerInAct := parseBloonSpawnCode(act.Code)
			if layer <= 0 {
				layer = 1
			}
			if maxLayerInAct > layer {
				layer = maxLayerInAct
			}
			if !found || layer > maxLayer {
				maxLayer = layer
				found = true
			}
		}
	}
	if !found {
		return "Red_Bloon_Spr"
	}
	return previewSpriteForLayer(maxLayer)
}

func previewSpriteForLayer(layer float64) string {
	switch {
	case layer >= 5000:
		return "New_DDT_Spr"
	case layer >= 1200:
		return "New_ZOMG_Spr"
	case layer >= 600:
		return "New_BFB_Spr"
	case layer >= 90:
		return "Moab_Spr"
	case layer >= 18:
		return "Ceramic_Bloon_Spr"
	case layer <= 1:
		return "Red_Bloon_Spr"
	case layer <= 1.5:
		return "Orange_Bloon_Spr"
	case layer <= 2:
		return "Blue_Bloon_Spr"
	case layer <= 2.5:
		return "Cyan_Bloon_Spr"
	case layer <= 3:
		return "Green_Bloon_Spr"
	case layer <= 3.5:
		return "Lime_Bloon_Spr"
	case layer <= 4:
		return "Yellow_Bloon_Spr"
	case layer <= 4.5:
		return "Amber_Bloon_Spr"
	case layer <= 5:
		return "Pink_Bloon_Spr"
	case layer <= 5.5:
		return "Purple_Bloon_Spr"
	case layer <= 6:
		return "Black_Bloon_Spr"
	case layer <= 6.1:
		return "White_Bloon_Spr"
	case layer <= 7:
		return "Zebra_Bloon_Spr"
	case layer <= 8:
		return "Rainbow_Bloon_Spr"
	case layer <= 8.5:
		return "Prismatic_Bloon_Spr"
	default:
		return "Ceramic_Bloon_Spr"
	}
}

// track-specific controllers
// each sets global.track and track-specific vars

// makeTrackController creates a controller behavior for a specific track
func makeTrackController(trackNum float64) engine.InstanceBehavior {
	return &trackControllerBehavior{trackNum: trackNum}
}

type trackControllerBehavior struct {
	engine.DefaultBehavior
	trackNum float64
}

func (b *trackControllerBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["track"] = b.trackNum
	g.GlobalVars["showhints"] = 0.0

	// wealthiness based on mode
	if getGlobal(g, "impoppablemodeselect") == 1 {
		g.GlobalVars["wealthiness"] = 6.0
	} else if getGlobal(g, "nightmaremodeselect") == 1 {
		g.GlobalVars["wealthiness"] = 8.0
	} else {
		g.GlobalVars["wealthiness"] = 5.0
	}
	g.GlobalVars["healthiness"] = 0.0
}

func (b *trackControllerBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// track high score
	key := fmt.Sprintf("x%d", int(b.trackNum))
	pts := getGlobal(g, "points")
	if pts > getGlobal(g, key) {
		g.GlobalVars[key] = math.Floor(pts)
	}
}

// countBloons counts active bloon instances
func countBloons(g *engine.Game) int {
	count := 0
	for _, inst := range g.InstanceMgr.FindByObject("Normal_Bloon_Branch") {
		if !inst.Destroyed {
			count++
		}
	}
	return count
}

// executeBloonSpawn handles a timeline action to spawn a bloon
// parses the spawn code to extract bloonlayer/bloonmaxlayer
func executeBloonSpawn(g *engine.Game, action engine.TimelineAction) {
	// timeline actions spawn bloons at BloonSpawn position
	spawns := g.InstanceMgr.FindByObject("BloonSpawn")
	if len(spawns) == 0 {
		return
	}
	spawn := spawns[0]

	// parse bloonlayer and bloonmaxlayer from the spawn code
	bloonlayer, bloonmaxlayer := parseBloonSpawnCode(action.Code)

	// create bloon at spawn point
	bloon := g.InstanceMgr.Create("Normal_Bloon_Branch", spawn.X, spawn.Y)
	if bloon != nil {
		bloon.Vars["bloonlayer"] = bloonlayer
		bloon.Vars["bloonmaxlayer"] = bloonmaxlayer
	}
}

// parseBloonSpawnCode extracts bloonlayer and bloonmaxlayer from spawn code
func parseBloonSpawnCode(code string) (bloonlayer, bloonmaxlayer float64) {
	bloonlayer = 1
	bloonmaxlayer = 1

	// simple parser for patterns like "bloonlayer = 1;" and "bloonmaxlayer = 1;"
	bloonlayer = extractCodeVar(code, "bloonlayer")
	bloonmaxlayer = extractCodeVar(code, "bloonmaxlayer")

	if bloonlayer == 0 {
		bloonlayer = 1
	}
	if bloonmaxlayer == 0 {
		bloonmaxlayer = bloonlayer
	}
	return
}

// extractCodeVar extracts a numeric value from code like "varname = 1.5;"
func extractCodeVar(code, varName string) float64 {
	// find "varname = " followed by a number
	searchFor := varName + " = "
	idx := 0
	for {
		pos := indexOf(code[idx:], searchFor)
		if pos < 0 {
			// try without spaces around =
			searchFor2 := varName + "="
			pos = indexOf(code[idx:], searchFor2)
			if pos < 0 {
				return 0
			}
			idx += pos + len(searchFor2)
		} else {
			idx += pos + len(searchFor)
		}

		// parse number
		numStr := ""
		for i := idx; i < len(code); i++ {
			ch := code[i]
			if (ch >= '0' && ch <= '9') || ch == '.' || ch == '-' {
				numStr += string(ch)
			} else {
				break
			}
		}
		if numStr != "" {
			val := 0.0
			fmt.Sscanf(numStr, "%f", &val)
			return val
		}
		break
	}
	return 0
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// registerGameplayBehaviors registers all gameplay behaviors
func RegisterGameplayBehaviors(im *engine.InstanceManager) {
	// core controllers
	im.RegisterBehavior("Go", func() engine.InstanceBehavior { return &GoBehavior{} })
	im.RegisterBehavior("Control", func() engine.InstanceBehavior { return &ControlBehavior{} })
	im.RegisterBehavior("Draw", func() engine.InstanceBehavior { return &DrawHUD{} })
	im.RegisterBehavior("Home_Button", func() engine.InstanceBehavior { return &HomeButtonBehavior{} })
	im.RegisterBehavior("Settings", func() engine.InstanceBehavior { return &SettingsBehavior{} })
	im.RegisterBehavior("BloonSpawn", func() engine.InstanceBehavior { return &BloonSpawnBehavior{} })
	im.RegisterBehavior("End", func() engine.InstanceBehavior { return &EndBehavior{} })
	im.RegisterBehavior("Wave_Panel", func() engine.InstanceBehavior { return &WavePanelBehavior{} })
	im.RegisterBehavior("Time_Wave_Panel", func() engine.InstanceBehavior { return &WavePanelBehavior{} })
	im.RegisterBehavior("auto_start_button", func() engine.InstanceBehavior { return &AutoStartButton{} })
	im.RegisterBehavior("Upgrade_Sign", func() engine.InstanceBehavior { return &engine.DefaultBehavior{} })
	im.RegisterBehavior("Wanna_go_to_main_", func() engine.InstanceBehavior { return &WannaGoToMainBehavior{} })

	// track controllers — each sets global.track
	trackControllers := map[string]float64{
		"Monkey_Meadows_Controler":       1,
		"Bloon_Oasis_Controler":          2,
		"Swamp_Spirals_Controler":        3,
		"Monkey_Fort_Controler":          4,
		"Monkey_Town_Docks_Controler":    5,
		"Conveyor_Belts_Controler":       6,
		"The_Depths_Controler":           7,
		"Sun_Stone_Controler":            8,
		"Shade_Woods_Controler":          9,
		"Minecarts_Controler":            10,
		"Crimson_Creek_Controler":        11,
		"Xtreme_Park_Controler":          12,
		"Portal_Lab_Controler":           13,
		"Omega_River_Controler":          14,
		"Space_Portals_Controler":        15,
		"Bloonlight_Throwback_Controler": 17,
		"Bloon_Circles_X_Controler":      18,
		"Autumn_Acres_Controler":         19,
		"Graveyard_Controler":            20,
		"Village_Defense_Controler":      21,
		"Circuit_Controler":              22,
		"Grand_Canyon_Controler":         23,
		"Bloonside_River_Controler":      24,
		"Cotton_Fields_Controler":        25,
		"Rubber_Rug_Controler":           27,
		"Frozen_Lake_Controler":          28,
		"Sky_Battles_Controler":          29,
		"Lava_Stream_Controler":          30,
		"Ravine_River_Controler":         31,
		"Peaceful_Bridge_Controler":      32,

		// aliases used by extracted rooms (Controller/Control naming variants).
		"Hard_Monkey_Meadows_Controller":   1,
		"Hard_Bloon_Oasis_Controller":      2,
		"Spiral_Swamp_Controler":           3,
		"Monkey_Fort_Controller":           4,
		"Monkey_Docks_Controller":          5,
		"Conveyor_Belt_Controller":         6,
		"The_Depths_Controller":            7,
		"Sun_Dial_Controller":              8,
		"Shade_Woods_Controller":           9,
		"Hard_Shade_Woods_Controller":      9,
		"Minecarts_Controller":             10,
		"Minecarts_Hard_Controller":        10,
		"Crimson_Creek_Controller":         11,
		"Crimson_Creek_Controller_Hard":    11,
		"Xtreme_Park_Controller":           12,
		"Portal_Lab_Controller":            13,
		"Omega_River_Controller":           14,
		"Space_Portals_Controller":         15,
		"Bloon_Light_Throwback_Controller": 17,
		"Bloon_Circles_X_Controller":       18,
		"Autum_Control":                    19,
		"Graveyard_Control":                20,
		"Village_Control":                  21,
		"Circuit_Control":                  22,
		"Canyon_Control":                   23,
		"Bloonside_Control":                24,
		"Cotton_Controller":                25,
		"Rubber_Rug_Control":               27,
		"Frozen_Lake_Controller":           28,
		"Sky_Battles_Controller":           29,
		"Lava_Stream_Controller":           30,
		"Ravine_River_Controller":          31,
		"Peaceful_Bridge_Controller":       32,
	}
	for name, trackNum := range trackControllers {
		tn := trackNum
		im.RegisterBehavior(name, func() engine.InstanceBehavior {
			return makeTrackController(tn)
		})
	}
}
