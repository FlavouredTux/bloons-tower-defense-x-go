package behaviors

import (
	"fmt"
	"math"
	"math/rand"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

// ===========================================================================
// Sandbox mode — faithful port of the original GML sandbox GUI objects:
//   Sanbox_Bar, Sandbox_Settings, Sandbox_Go, Bloon_Button,
//   Special_Bloon_Button, Bloon_Up, Bloon_Down, sandbox_sender
// ===========================================================================

// ---------------------------------------------------------------------------
// Sanbox_Bar — sandbox play button on Track_Setup_II  (original GML typo kept)
// ---------------------------------------------------------------------------
// Sprite: Sandbox_spr. On click: unlocks all towers, routes to _Sand room.
type SanboxBarBehavior struct {
	engine.DefaultBehavior
}

func (b *SanboxBarBehavior) Create(inst *engine.Instance, g *engine.Game) {
	// original GML: Sanbox_Bar.Create resets all modifiers to 0,
	// sets towerlimit=1000000, pointmultiplier=1, and unlocks all tower locks.
	g.GlobalVars["challenge"] = 0.0
	g.GlobalVars["pointmultiplier"] = 1.0
	g.GlobalVars["towerlimit"] = 1000000.0

	g.GlobalVars["sixtowers"] = 0.0
	g.GlobalVars["randomtowers"] = 0.0
	g.GlobalVars["wavesqueeze"] = 0.0
	g.GlobalVars["waveskip"] = 0.0
	g.GlobalVars["strongerbloons"] = 0.0
	g.GlobalVars["fasterbloons"] = 0.0
	g.GlobalVars["noliveslost"] = 0.0

	if getGlobal(g, "rank") < 20 {
		g.GlobalVars["blooninfo"] = 1.0
	} else {
		g.GlobalVars["blooninfo"] = 0.0
	}

	// zero all tower path-upgrade locks (they'll be set to 1 on room goto)
	unlockAllTowers(g)
}

func (b *SanboxBarBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	ts := getGlobal(g, "trackselect")

	sandRooms := map[float64]string{
		1: "Monkey_Meadows_Sand", 2: "Bloon_Oasis_Sand",
		3: "Swamp_Spirals_Sand", 4: "Monkey_Fort_Sand",
		5: "Monkey_Town_Docks_Sand", 6: "Conveyor_Belts_Sand",
		7: "The_Depths_Sand", 8: "Sun_Stone_Sand",
		9: "Shade_Woods_Sand", 10: "Minecarts_Sand",
		11: "Crimson_Creek_Sand", 12: "Xtreme_Park_Sand",
		13: "Portal_Lab_Sand", 14: "Omega_River_Sand",
		15: "Space_Portals_Sand", 17: "Bloonlight_Throwback_Sand",
		18: "Bloon_Circles_X_Sand", 19: "Autumn_Acres_Sand",
		20: "Graveyard_Sand", 21: "Village_Defense_Sand",
		22: "Circuit_Sand", 23: "Grand_Canyon_Sand",
		24: "Bloonside_River_Sand",
		27: "Rubber_Rug_Sand", 28: "Frozen_Lake_Sand",
		29: "Sky_Battles_Sand", 30: "Lava_Stream_Sand",
		31: "Ravine_River_Sand", 32: "Peaceful_Bridge_Sand",
	}

	room, ok := sandRooms[ts]
	if ok {
		g.RequestRoomGoto(room)
	}
}

// ---------------------------------------------------------------------------
// Sandbox_Settings — placed at (0,0) in each _Sand room
// ---------------------------------------------------------------------------
// Sets sandbox=1, disables agent toggles, spawns the entire scrollable
// bloon picker panel (Bloon_Button + Special_Bloon_Button instances),
// transforms Go → Sandbox_Go, and creates Bloon_Up/Bloon_Down scroll buttons.
type SandboxSettingsBehavior struct {
	engine.DefaultBehavior
}

func (b *SandboxSettingsBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["sandbox"] = 1.0

	// disable boss/agent bloon enables
	g.GlobalVars["bullyenable"] = 0.0
	g.GlobalVars["mmoabenable"] = 0.0
	g.GlobalVars["ufoenable"] = 0.0
	g.GlobalVars["horrorenable"] = 0.0
	g.GlobalVars["superenable"] = 0.0
	g.GlobalVars["motherenable"] = 0.0
	g.GlobalVars["lolenable"] = 0.0

	// --- Spawn bloon picker buttons (original GML Sandbox_Settings.Create) ---
	// addon tracks the Y position slot; each button is 64px tall.
	addon := 1
	ox := inst.X // 0
	oy := inst.Y // 0

	// ========== Normal Bloons ==========
	// Red (layer 1), 7 variants (uptype 0-6)
	for i := 0; i < 7; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 1.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	// Blue (layer 2)
	for i := 0; i < 7; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 2.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	// Green (layer 3)
	for i := 0; i < 7; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 3.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	// Yellow (layer 4)
	for i := 0; i < 7; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 4.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	// Pink (layer 5)
	for i := 0; i < 7; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 5.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	// Black (layer 6), 6 variants
	for i := 0; i < 6; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 6.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	addon++ // gap
	// White (layer 6.1), 6 variants
	for i := 0; i < 6; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 6.1
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	addon++ // gap
	// Zebra (layer 7), 6 variants
	for i := 0; i < 6; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 7.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	addon++ // gap
	// Rainbow (layer 8), 7 variants
	for i := 0; i < 7; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 8.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	// Ceramic (layer 18), 7 variants
	for i := 0; i < 7; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 18.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	// Brick (layer 48), 6 variants
	for i := 0; i < 6; i++ {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = 48.0
			bb.Vars["uptype"] = float64(i)
		}
		addon++
	}
	addon++ // gap

	// ========== MOAB Class Bloons ==========
	moabEntries := []struct {
		bloonup float64
		uptype  float64
	}{
		{93, 0}, {348, 0}, {1248, 0}, {5248, 0},
		{68.5, 0}, {593, 0}, {351, 0}, {318, 0},
	}
	for _, e := range moabEntries {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = e.bloonup
			bb.Vars["uptype"] = e.uptype
		}
		addon++
	}
	// Fortified MOAB variants (uptype=2)
	fortifiedMoabs := []float64{93, 348, 1248, 5248}
	for _, bl := range fortifiedMoabs {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 2.0
		}
		addon++
	}

	addon += 2 // gap

	// ========== Splitting Bloons ==========
	splittingLayers := []float64{1.5, 2.5, 3.5, 4.5, 5.5, 8.5}
	for _, layer := range splittingLayers {
		for i := 0; i < 6; i++ {
			bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
			if bb != nil {
				bb.Vars["Bloonup"] = layer
				bb.Vars["uptype"] = float64(i)
			}
			addon++
		}
		addon++ // gap between each splitting type
	}

	// ========== Nightmare Bloons (Special_Bloon_Button) ==========
	// Stuffed (uptype=1), layers 1-5
	for i := 1; i <= 5; i++ {
		sb := g.InstanceMgr.Create("Special_Bloon_Button", ox, oy+float64(addon)*64)
		if sb != nil {
			sb.Vars["Bloonup"] = float64(i)
			sb.Vars["uptype"] = 1.0
		}
		addon++
	}
	addon++ // gap
	// Ninja (uptype=2), layers 1-5
	for i := 1; i <= 5; i++ {
		sb := g.InstanceMgr.Create("Special_Bloon_Button", ox, oy+float64(addon)*64)
		if sb != nil {
			sb.Vars["Bloonup"] = float64(i)
			sb.Vars["uptype"] = 2.0
		}
		addon++
	}
	addon++ // gap
	// Robo (uptype=3), layers 1-5
	for i := 1; i <= 5; i++ {
		sb := g.InstanceMgr.Create("Special_Bloon_Button", ox, oy+float64(addon)*64)
		if sb != nil {
			sb.Vars["Bloonup"] = float64(i)
			sb.Vars["uptype"] = 3.0
		}
		addon++
	}
	addon++ // gap
	// Patrol (uptype=4)
	{
		sb := g.InstanceMgr.Create("Special_Bloon_Button", ox, oy+float64(addon)*64)
		if sb != nil {
			sb.Vars["Bloonup"] = 1.0
			sb.Vars["uptype"] = 4.0
		}
		addon++
	}
	// Barrier (uptype=5), layers 3 and 6
	for _, layer := range []float64{3, 6} {
		sb := g.InstanceMgr.Create("Special_Bloon_Button", ox, oy+float64(addon)*64)
		if sb != nil {
			sb.Vars["Bloonup"] = layer
			sb.Vars["uptype"] = 5.0
		}
		addon++
	}
	// Planetarium (uptype=6), layers 10 and 20
	for _, layer := range []float64{10, 20} {
		sb := g.InstanceMgr.Create("Special_Bloon_Button", ox, oy+float64(addon)*64)
		if sb != nil {
			sb.Vars["Bloonup"] = layer
			sb.Vars["uptype"] = 6.0
		}
		addon++
	}
	// Spectrum (uptype=7)
	{
		sb := g.InstanceMgr.Create("Special_Bloon_Button", ox, oy+float64(addon)*64)
		if sb != nil {
			sb.Vars["Bloonup"] = 9.0
			sb.Vars["uptype"] = 7.0
		}
		addon++
	}
	addon++ // gap

	// ========== Nightmare MOAB-class ==========
	nightmareMoabs := []float64{10068.5, 2593, 248, 918, 3351, 17248}
	for _, bl := range nightmareMoabs {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}
	addon++ // gap

	// ========== Boss Bloons ==========
	bossEntries := []float64{
		10014, 10015, 10016, 10017, // Bully colored
		10011, 10012, 10013,        // Bully, Rage, Bully_Rage
	}
	for _, bl := range bossEntries {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}
	addon++ // gap

	// Mighty MOABs
	for _, bl := range []float64{10021, 10022, 10023, 10024} {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}
	addon++ // gap

	// UFOs
	for _, bl := range []float64{10041, 10042, 10043} {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}
	addon++ // gap

	// Supers
	for _, bl := range []float64{10051, 10052, 10053, 10054} {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}
	addon++ // gap

	// Mothers
	for _, bl := range []float64{10061, 10062, 10063, 10064} {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}
	addon++ // gap

	// Terribles
	for _, bl := range []float64{10071, 10072, 10073, 10074} {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}
	addon++ // gap

	// Clowns
	for _, bl := range []float64{10081, 10082, 10083, 10084} {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}
	addon++ // gap

	// Bloomers
	for _, bl := range []float64{10091, 10092, 10093, 10094} {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}
	addon++ // gap

	// Destroyers
	for _, bl := range []float64{10111, 10112} {
		bb := g.InstanceMgr.Create("Bloon_Button", ox, oy+float64(addon)*64)
		if bb != nil {
			bb.Vars["Bloonup"] = bl
			bb.Vars["uptype"] = 0.0
		}
		addon++
	}

	// --- Create scroll buttons ---
	g.InstanceMgr.Create("Bloon_Up", 0, 512)
	g.InstanceMgr.Create("Bloon_Down", 0, 544)
}

func (b *SandboxSettingsBehavior) Step(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["points"] = 0.0
}

func (b *SandboxSettingsBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr != nil && len(spr.Frames) > 0 {
		engine.DrawSpriteExt(screen, spr.Frames[0], spr.XOrigin, spr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
	}
}

// ---------------------------------------------------------------------------
// Sandbox_Go — speed control only (no waves in sandbox)
// ---------------------------------------------------------------------------
// Sprite: sprite278 (idle) / Going (when bloons active). Placed in _Sand rooms.
// Left-click or Space/Enter cycles game speed: 30→60→90→30.
// Right-click toggles 10x speed (300).
type SandboxGoBehavior struct {
	engine.DefaultBehavior
	shiftpress int
}

func (b *SandboxGoBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["autostart"] = 0.0
	g.GlobalVars["freeplay"] = 1.0 // sandbox is always freeplay
	g.GlobalVars["wavenow"] = 0.0
	g.GlobalVars["wave"] = 1.0
	g.GlobalVars["endsequence"] = 0.0
	g.GlobalVars["bpower"] = 1.0
	g.GlobalVars["bspeed"] = 1.0

	// sandbox: infinite resources
	g.GlobalVars["money"] = 75000000.0
	g.GlobalVars["life"] = 99999.0
	g.GlobalVars["points"] = 0.0
	g.GlobalVars["gamespeed"] = 30.0
	g.SetGameSpeed(30)
	b.shiftpress = 0
}

func (b *SandboxGoBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// infinite resources every frame
	g.GlobalVars["life"] = 999999.0
	g.GlobalVars["money"] = 9999999.0

	// switch sprite based on whether bloons exist
	bloonCount := countBloons(g)
	if bloonCount > 0 {
		inst.SpriteName = "Going"
	} else {
		inst.SpriteName = "sprite278"
	}
}

func (b *SandboxGoBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	b.cycleSpeed(inst, g)
}

func (b *SandboxGoBehavior) KeyPress(inst *engine.Instance, g *engine.Game) {
	if g.InputMgr.KeyPressed(ebiten.KeySpace) || g.InputMgr.KeyPressed(ebiten.KeyEnter) {
		b.cycleSpeed(inst, g)
	}
}

func (b *SandboxGoBehavior) MouseRightPressed(inst *engine.Instance, g *engine.Game) {
	// 10x speed toggle
	if getGlobal(g, "gamespeed") == 300 {
		g.SetGameSpeed(30)
		g.GlobalVars["gamespeed"] = 30.0
		b.shiftpress = 0
	} else {
		g.SetGameSpeed(300)
		g.GlobalVars["gamespeed"] = 300.0
	}
}

func (b *SandboxGoBehavior) cycleSpeed(inst *engine.Instance, g *engine.Game) {
	switch b.shiftpress {
	case 0:
		g.SetGameSpeed(60)
		g.GlobalVars["gamespeed"] = 60.0
		b.shiftpress = 1
	case 1:
		g.SetGameSpeed(90)
		g.GlobalVars["gamespeed"] = 90.0
		b.shiftpress = 2
	case 2:
		g.SetGameSpeed(30)
		g.GlobalVars["gamespeed"] = 30.0
		b.shiftpress = 0
	}
}

func (b *SandboxGoBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr != nil && len(spr.Frames) > 0 {
		engine.DrawSpriteExt(screen, spr.Frames[0], spr.XOrigin, spr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
	}

	// speed indicator bars (same as GoBehavior)
	gs := getGlobal(g, "gamespeed")
	fill := gs / 90.0
	if fill > 1 {
		fill = 1
	}
	barH := 44.0
	filledH := barH * fill
	barTop := inst.Y + 10
	barBottom := inst.Y + 54

	drawRect(screen, inst.X+8, barBottom-filledH, 2, filledH, speedBarColor(fill))
	drawRect(screen, inst.X+8, barTop, 2, barH-filledH, [3]uint8{0, 0, 0})
	drawRect(screen, inst.X+53, barBottom-filledH, 2, filledH, speedBarColor(fill))
	drawRect(screen, inst.X+53, barTop, 2, barH-filledH, [3]uint8{0, 0, 0})
}

// ---------------------------------------------------------------------------
// Bloon_Button — scrollable bloon spawn button in the sandbox picker panel
// ---------------------------------------------------------------------------
// Sprite: Wave_Paper_spr (64x64). Vars: Bloonup (bloon code), uptype (variant).
// Left-click spawns 1 bloon, right-click spawns 20 (mass).
// Buttons hide (depth=200) when scrolled off-screen (Y>511 or Y<64).
type BloonButtonBehavior struct {
	engine.DefaultBehavior
}

func isSandboxBossCode(bloonup float64) bool {
	switch bloonup {
	case 10011, 10012, 10013, 10014, 10015, 10016, 10017,
		10021, 10022, 10023, 10024,
		10041, 10042, 10043,
		10051, 10052, 10053, 10054,
		10061, 10062, 10063, 10064,
		10071, 10072, 10073, 10074,
		10081, 10082, 10083, 10084,
		10091, 10092, 10093, 10094,
		10111, 10112:
		return true
	default:
		return false
	}
}

func (b *BloonButtonBehavior) Create(inst *engine.Instance, g *engine.Game) {
	if _, ok := inst.Vars["Bloonup"]; !ok {
		inst.Vars["Bloonup"] = 1.0
	}
	if _, ok := inst.Vars["uptype"]; !ok {
		inst.Vars["uptype"] = 0.0
	}
	inst.SpriteName = "Wave_Paper_spr"
	inst.Depth = -19
}

func (b *BloonButtonBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// depth-based visibility: hide when scrolled off-screen
	if inst.Y > 511 {
		inst.Depth = 200
	} else if inst.Y < 64 {
		inst.Depth = 200
	} else {
		inst.Depth = -19
	}
}

// MouseLeftPressed — spawn 1 bloon (primary click)
func (b *BloonButtonBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if inst.Depth > 0 {
		return
	}
	bloonup := getVar(inst, "Bloonup")
	uptype := getVar(inst, "uptype")
	fmt.Printf("[sandbox] click normal single: button=%d bloonup=%.3f uptype=%.0f\n", inst.ID, bloonup, uptype)

	// Boss bloons spawn directly. Nightmare MOAB-class entries like 10068.5
	// and 17248 are not bosses and must still go through the normal sender.
	if isSandboxBossCode(bloonup) {
		spawnBossBloon(inst, g, bloonup)
		return
	}

	// create sandbox_sender for normal/moab bloons
	sender := g.InstanceMgr.Create("sandbox_sender", inst.X, inst.Y)
	if sender == nil {
		return
	}
	sender.Vars["bloonsetlayer"] = bloonup
	sender.Vars["stack1type"] = uptype
	configureSenderStack1(sender, bloonup, uptype, 1, g) // single spawn
}

// MouseRightPressed — spawn 20 bloons (mass)
func (b *BloonButtonBehavior) MouseRightPressed(inst *engine.Instance, g *engine.Game) {
	if inst.Depth > 0 {
		return
	}
	bloonup := getVar(inst, "Bloonup")
	uptype := getVar(inst, "uptype")
	fmt.Printf("[sandbox] click normal mass: button=%d bloonup=%.3f uptype=%.0f\n", inst.ID, bloonup, uptype)

	if isSandboxBossCode(bloonup) {
		spawnBossBloon(inst, g, bloonup)
		return
	}

	sender := g.InstanceMgr.Create("sandbox_sender", inst.X, inst.Y)
	if sender == nil {
		return
	}
	sender.Vars["bloonsetlayer"] = bloonup
	sender.Vars["stack1type"] = uptype
	configureSenderStack1(sender, bloonup, uptype, 20, g) // mass spawn
}

func (b *BloonButtonBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.Depth > 0 {
		return
	}
	// draw button background
	spr := g.AssetManager.GetSprite("Wave_Paper_spr")
	if spr != nil && len(spr.Frames) > 0 {
		engine.DrawSpriteExt(screen, spr.Frames[0], spr.XOrigin, spr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, 0, 1)
	}

	// draw bloon icon
	bloonup := getVar(inst, "Bloonup")
	uptype := int(getVar(inst, "uptype"))
	drawBloonButtonIcon(screen, g, bloonup, uptype, inst.X+31, inst.Y+32)
}

// ---------------------------------------------------------------------------
// Special_Bloon_Button — nightmare bloon spawn button
// ---------------------------------------------------------------------------
// Parent is Bloon_Button. Uses stack10/alarm[10] for special bloon types.
// uptype: 1=Stuffed, 2=Ninja, 3=Robo, 4=Patrol, 5=Barrier, 6=Planetarium, 7=Spectrum
type SpecialBloonButtonBehavior struct {
	engine.DefaultBehavior
}

func (b *SpecialBloonButtonBehavior) Create(inst *engine.Instance, g *engine.Game) {
	if _, ok := inst.Vars["Bloonup"]; !ok {
		inst.Vars["Bloonup"] = 1.0
	}
	if _, ok := inst.Vars["uptype"]; !ok {
		inst.Vars["uptype"] = 0.0
	}
	inst.SpriteName = "Wave_Paper_spr"
	inst.Depth = -19
}

func (b *SpecialBloonButtonBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if inst.Y > 511 {
		inst.Depth = 200
	} else if inst.Y < 64 {
		inst.Depth = 200
	} else {
		inst.Depth = -19
	}
}

// MouseLeftPressed — spawn 1 special bloon
func (b *SpecialBloonButtonBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if inst.Depth > 0 {
		return
	}
	bloonup := getVar(inst, "Bloonup")
	uptype := getVar(inst, "uptype")
	fmt.Printf("[sandbox] click special single: button=%d bloonup=%.3f uptype=%.0f\n", inst.ID, bloonup, uptype)

	sender := g.InstanceMgr.Create("sandbox_sender", inst.X, inst.Y)
	if sender == nil {
		return
	}
	sender.Vars["bloonsetlayer"] = bloonup
	configureSenderStack10(sender, bloonup, uptype, 1, g)
}

// MouseRightPressed — mass spawn special bloons
func (b *SpecialBloonButtonBehavior) MouseRightPressed(inst *engine.Instance, g *engine.Game) {
	if inst.Depth > 0 {
		return
	}
	bloonup := getVar(inst, "Bloonup")
	uptype := getVar(inst, "uptype")
	fmt.Printf("[sandbox] click special mass: button=%d bloonup=%.3f uptype=%.0f\n", inst.ID, bloonup, uptype)

	sender := g.InstanceMgr.Create("sandbox_sender", inst.X, inst.Y)
	if sender == nil {
		return
	}
	sender.Vars["bloonsetlayer"] = bloonup
	configureSenderStack10(sender, bloonup, uptype, 20, g) // mass
}

func (b *SpecialBloonButtonBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.Depth > 0 {
		return
	}
	spr := g.AssetManager.GetSprite("Wave_Paper_spr")
	if spr != nil && len(spr.Frames) > 0 {
		engine.DrawSpriteExt(screen, spr.Frames[0], spr.XOrigin, spr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, 0, 1)
	}

	// draw special bloon icon based on uptype
	uptype := int(getVar(inst, "uptype"))
	bloonup := getVar(inst, "Bloonup")
	drawSpecialBloonIcon(screen, g, uptype, bloonup, inst.X+31, inst.Y+32)
}

// ---------------------------------------------------------------------------
// Bloon_Up / Bloon_Down — scroll the bloon picker panel
// ---------------------------------------------------------------------------

type BloonUpBehavior struct {
	engine.DefaultBehavior
}

func (b *BloonUpBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Scroll_Up_Spr"
	inst.Depth = -20
}

func (b *BloonUpBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	scrollBloonPanel(g, 1) // scroll up = buttons move down (Y increases)
}

func (b *BloonUpBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// mouse wheel up also scrolls
	_, wy := g.InputMgr.WheelDelta()
	if wy > 0 {
		// only scroll if mouse is on left side (bloon panel area)
		mx, _ := g.GetMouseRoomPos()
		if mx < 80 {
			scrollBloonPanel(g, 1)
		}
	}
}

type BloonDownBehavior struct {
	engine.DefaultBehavior
}

func (b *BloonDownBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Scroll_Down_Spr"
	inst.Depth = -20
	g.GlobalVars["bloonpanel"] = 0.0
}

func (b *BloonDownBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	scrollBloonPanel(g, -1) // scroll down = buttons move up (Y decreases)
}

func (b *BloonDownBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// mouse wheel down also scrolls
	_, wy := g.InputMgr.WheelDelta()
	if wy < 0 {
		mx, _ := g.GetMouseRoomPos()
		if mx < 80 {
			scrollBloonPanel(g, -1)
		}
	}
}

// scrollBloonPanel moves all Bloon_Button and Special_Bloon_Button instances.
// dir=1 scrolls up (page), dir=-1 scrolls down.
func scrollBloonPanel(g *engine.Game, dir int) {
	panel := getGlobal(g, "bloonpanel")
	if dir == 1 && panel <= 0 {
		return // already at top
	}
	if dir == -1 && panel >= 6300 {
		return // limit
	}

	// each scroll step moves buttons by 448px (7 button heights of 64px)
	delta := float64(dir) * 448.0
	g.GlobalVars["bloonpanel"] = panel - float64(dir)

	for _, bb := range g.InstanceMgr.FindByObject("Bloon_Button") {
		bb.Y += delta
	}
	for _, sb := range g.InstanceMgr.FindByObject("Special_Bloon_Button") {
		sb.Y += delta
	}
}

// ---------------------------------------------------------------------------
// sandbox_sender — invisible bloon spawning engine (alarm-driven)
// ---------------------------------------------------------------------------
// Created by Bloon_Button clicks. Has 11 stacks, each using a dedicated alarm.
// Alarm[0] = self-destruct. Alarms[1-9] = normal/MOAB stacks. Alarms[10-11] = special.
type SandboxSenderBehavior struct {
	engine.DefaultBehavior
}

func (b *SandboxSenderBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Visible = false
	inst.Depth = -20

	// initialize all stacks to 0
	for i := 1; i <= 11; i++ {
		inst.Vars[fmt.Sprintf("stack%dtype", i)] = 0.0
		inst.Vars[fmt.Sprintf("stack%dlayer", i)] = 0.0
		inst.Vars[fmt.Sprintf("stack%damount", i)] = 0.0
		inst.Vars[fmt.Sprintf("stack%ddelay", i)] = 0.0
		inst.Vars[fmt.Sprintf("stack%dshield", i)] = 0.0
	}
	inst.Vars["bloonsetlayer"] = 0.0
}

func (b *SandboxSenderBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// randomize cycle each frame (original GML: global.cycle = ceil(0.001 + random(3.99)))
	g.GlobalVars["cycle"] = math.Ceil(0.001 + rand.Float64()*3.99)
}

func (b *SandboxSenderBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		// self-destruct
		g.InstanceMgr.Destroy(inst.ID)
		return
	}

	// stacks 1-9 use alarms 1-9 (normal/MOAB bloons)
	if idx >= 1 && idx <= 9 {
		b.spawnFromStack(inst, g, idx)
		return
	}

	// stacks 10-11 use alarms 10-11 (special bloons)
	if idx >= 10 && idx <= 11 {
		b.spawnSpecialFromStack(inst, g, idx)
		return
	}
}

// spawnFromStack handles alarm-triggered spawning for stacks 1-9
func (b *SandboxSenderBehavior) spawnFromStack(inst *engine.Instance, g *engine.Game, stackIdx int) {
	prefix := fmt.Sprintf("stack%d", stackIdx)
	amount := getVar(inst, prefix+"amount")
	if amount <= 0 {
		return
	}

	layer := getVar(inst, prefix+"layer")
	stype := getVar(inst, prefix+"type")
	delay := getVar(inst, prefix+"delay")
	shield := getVar(inst, prefix+"shield")
	fmt.Printf("[sandbox] sender=%d stack=%d normal spawn request: layer=%.3f type=%.0f shield=%.1f remaining=%.0f\n",
		inst.ID, stackIdx, layer, stype, shield, amount)

	// spawn one bloon
	spawnSandboxBloon(g, layer, stype, shield)

	// decrement amount
	amount--
	inst.Vars[prefix+"amount"] = amount

	if amount > 0 {
		inst.Alarms[stackIdx] = int(math.Max(1, delay))
	}
}

// spawnSpecialFromStack handles alarm-triggered spawning for stacks 10-11
func (b *SandboxSenderBehavior) spawnSpecialFromStack(inst *engine.Instance, g *engine.Game, stackIdx int) {
	prefix := fmt.Sprintf("stack%d", stackIdx)
	amount := getVar(inst, prefix+"amount")
	if amount <= 0 {
		return
	}

	layer := getVar(inst, prefix+"layer")
	stype := getVar(inst, prefix+"type")
	delay := getVar(inst, prefix+"delay")
	shield := getVar(inst, prefix+"shield")
	fmt.Printf("[sandbox] sender=%d stack=%d special spawn request: layer=%.3f type=%.0f shield=%.1f remaining=%.0f\n",
		inst.ID, stackIdx, layer, stype, shield, amount)

	// spawn one special bloon
	spawnSpecialSandboxBloon(g, layer, stype, shield)

	amount--
	inst.Vars[prefix+"amount"] = amount

	if amount > 0 {
		inst.Alarms[stackIdx] = int(math.Max(1, delay))
	}
}

// ---------------------------------------------------------------------------
// Bloon spawning helpers
// ---------------------------------------------------------------------------

// spawnSandboxBloon spawns a bloon at BloonSpawn using Normal_Bloon_Branch.
// stype maps to GML uptype: 0=normal, 1=tattered, 2=shielded(fortified),
// 3=regrow, 4=camo, 5=lead, 6=static(camo_lead), 7=camo_lead, 8=regrow_tattered
func spawnSandboxBloon(g *engine.Game, layer, stype, shield float64) {
	spawns := g.InstanceMgr.FindByObject("BloonSpawn")
	if len(spawns) == 0 {
		return
	}
	sp := spawns[0]

	bloon := g.InstanceMgr.Create("Normal_Bloon_Branch", sp.X, sp.Y)
	if bloon == nil {
		return
	}
	bloon.Vars["bloonlayer"] = layer
	bloon.Vars["bloonmaxlayer"] = layer

	// Lock nightmare MOAB-class bloons to their exact sprite to avoid the
	// generic layer resolver snapping to a neighboring MOAB family sprite.
	if spriteName, ok := map[float64]string{
		248:     "Rocket_Blimp_Spr",
		918:     "Storm_LPZ_Spr",
		2593:    "Mega_BRC_Spr",
		3351:    "Deadly_DDT_Spr",
		10068.5: "Prismatic_HTA_Spr",
		17248:   "Party_Blimp_Spr",
	}[layer]; ok {
		bloon.Vars["special_sprite"] = spriteName
	}

	// apply modifiers based on uptype
	switch int(stype) {
	case 1: // tattered
		bloon.Vars["tattered"] = 1.0
	case 2: // shielded (fortified)
		bloon.Vars["shielded"] = 1.0
		if shield > 0 {
			bloon.Vars["shield_hp"] = shield
		}
	case 3: // regrow
		bloon.Vars["regrow"] = 1.0
	case 4: // camo
		bloon.Vars["camo"] = 1.0
	case 5: // lead
		bloon.Vars["lead"] = 1.0
	case 6: // camo + lead (static)
		bloon.Vars["camo"] = 1.0
		bloon.Vars["lead"] = 1.0
		if shield > 0 {
			bloon.Vars["shielded"] = 1.0
			bloon.Vars["shield_hp"] = shield
		}
	case 7: // camo + lead
		bloon.Vars["camo"] = 1.0
		bloon.Vars["lead"] = 1.0
	case 8: // regrow + tattered
		bloon.Vars["regrow"] = 1.0
		bloon.Vars["tattered"] = 1.0
	}

	// MOAB-class bloons always get their shield (it's their health system)
	if layer >= 68.5 && shield > 0 && getVar(bloon, "shielded") == 0 {
		bloon.Vars["shielded"] = 1.0
		bloon.Vars["shield_hp"] = shield
	}

	updateBloonSprite(bloon)
	updateBloonDepth(bloon)
	applySandboxSpawnStagger(bloon, g)
	fmt.Printf("[sandbox] spawned normal bloon: id=%d layer=%.3f type=%.0f shielded=%.0f shield_hp=%.1f sprite=%s\n",
		bloon.ID, getVar(bloon, "bloonlayer"), stype, getVar(bloon, "shielded"), getVar(bloon, "shield_hp"), bloon.SpriteName)
}

// spawnSpecialSandboxBloon spawns a special/nightmare bloon.
// stype: 1=Stuffed, 2=Ninja, 3=Robo, 4=Patrol, 5=Barrier, 6=Planetarium, 7=Spectrum
func spawnSpecialSandboxBloon(g *engine.Game, layer, stype, shield float64) {
	spawns := g.InstanceMgr.FindByObject("BloonSpawn")
	if len(spawns) == 0 {
		return
	}
	sp := spawns[0]

	// Most special bloons don't have registered behaviors yet;
	// spawn as Normal_Bloon_Branch with the properties set.
	bloon := g.InstanceMgr.Create("Normal_Bloon_Branch", sp.X, sp.Y)
	if bloon == nil {
		return
	}
	bloon.Vars["bloonlayer"] = layer
	bloon.Vars["bloonmaxlayer"] = layer
	bloon.Vars["special_type"] = stype

	// special bloons use their own sprite, not the layer's sprite
	specialSprites := map[float64]string{
		1: "Stuffed_Bloon_Spr",
		2: "Ninja_Bloon_Spr",
		3: "Robo_Bloon_Spr",
		4: "Patrol_Bloon_Spr",
		5: "Barrier_Bloon_Spr",
		6: "Planetarium_Bloon_Spr",
		7: "Spectrum_Bloon_Spr",
	}
	if spriteName, ok := specialSprites[stype]; ok {
		bloon.Vars["special_sprite"] = spriteName
	}

	bpower := getGlobal(g, "bpower")
	if bpower == 0 {
		bpower = 1
	}

	switch int(stype) {
	case 1: // Stuffed — regrow armour (GML: armour starts at 0, regrows to maxarmour)
		bloon.Vars["stuffed"] = 1.0
		maxArmour := math.Ceil(layer * bpower)
		bloon.Vars["armour"] = maxArmour
		bloon.Vars["maxarmour"] = maxArmour
		bloon.Vars["regrow"] = 1.0
	case 2: // Ninja — camo, armoured (GML: armour = 12 initially, maxarmour = layer×9×bpower)
		bloon.Vars["camo"] = 1.0
		maxArmour := math.Ceil(layer * 9 * bpower)
		bloon.Vars["armour"] = maxArmour
		bloon.Vars["maxarmour"] = maxArmour
		bloon.Vars["ability_timer"] = 90.0 + rand.Float64()*45
		if shield > 0 {
			bloon.Vars["shielded"] = 1.0
			bloon.Vars["shield_hp"] = shield
		}
	case 3: // Robo — lead, armoured (GML: armour = 20 initially, maxarmour = layer×12×bpower)
		bloon.Vars["lead"] = 1.0
		maxArmour := math.Ceil(layer * 12 * bpower)
		bloon.Vars["armour"] = maxArmour
		bloon.Vars["maxarmour"] = maxArmour
		bloon.Vars["ability_timer"] = 400.0 + rand.Float64()*200
		if shield > 0 {
			bloon.Vars["shielded"] = 1.0
			bloon.Vars["shield_hp"] = shield
		}
	case 4: // Patrol — no path, self-destruct (GML: armour = 20, maxarmour = 60×bpower)
		maxArmour := math.Ceil(60 * bpower)
		bloon.Vars["armour"] = maxArmour
		bloon.Vars["maxarmour"] = maxArmour
		bloon.Vars["ability_timer"] = 30.0 // shooting timer
	case 5: // Barrier — spawns blockers (GML: armour starts 0, maxarmour = layer×10×bpower)
		maxArmour := math.Ceil(layer * 10 * bpower)
		bloon.Vars["armour"] = maxArmour
		bloon.Vars["maxarmour"] = maxArmour
		bloon.Vars["ability_timer"] = 200.0 + rand.Float64()*30
		if shield > 0 {
			bloon.Vars["shielded"] = 1.0
			bloon.Vars["shield_hp"] = shield
		}
	case 6: // Planetarium — spawns satellites (GML: maxarmour depends on layer)
		var maxArmour float64
		if layer >= 10 {
			maxArmour = math.Ceil(300 * bpower)
		} else {
			maxArmour = math.Ceil(layer * 30 * bpower)
		}
		bloon.Vars["armour"] = maxArmour
		bloon.Vars["maxarmour"] = maxArmour
		bloon.Vars["ability_timer"] = 90.0
		if shield > 0 {
			bloon.Vars["shielded"] = 1.0
			bloon.Vars["shield_hp"] = shield
		}
	case 7: // Spectrum — damage capped at 100/hit (GML: armour = 200×bpower)
		maxArmour := math.Ceil(200 * bpower)
		bloon.Vars["armour"] = maxArmour
		bloon.Vars["maxarmour"] = maxArmour
		bloon.Vars["damage_cap"] = 100.0
	}

	updateBloonSprite(bloon)
	updateBloonDepth(bloon)
	applySandboxSpawnStagger(bloon, g)
	fmt.Printf("[sandbox] spawned special bloon: id=%d layer=%.3f type=%.0f sprite=%s\n",
		bloon.ID, getVar(bloon, "bloonlayer"), stype, bloon.SpriteName)
}

func applySandboxSpawnStagger(bloon *engine.Instance, g *engine.Game) {
	pathName, _ := bloon.Vars["path_name"].(string)
	if pathName == "" || g.PathMgr == nil {
		return
	}

	active := 0
	for _, other := range g.InstanceMgr.FindByObject("Normal_Bloon_Branch") {
		if other == nil || other.ID == bloon.ID || other.Destroyed {
			continue
		}
		active++
	}

	progress := 0.0
	if active > 0 {
		progress = getGlobal(g, "sandbox_spawn_progress") + 0.09
		if progress > 0.9 {
			progress = 0.09
		}
		g.GlobalVars["sandbox_spawn_progress"] = progress
	} else {
		g.GlobalVars["sandbox_spawn_progress"] = 0.0
	}
	if progress <= 0 {
		return
	}

	bloon.Vars["path_progress"] = progress
	x, y := g.PathMgr.GetPositionAtProgress(pathName, progress)
	bloon.X = x
	bloon.Y = y
	fmt.Printf("[sandbox] staggered bloon: id=%d path=%s progress=%.3f active=%d\n",
		bloon.ID, pathName, progress, active)
}

// spawnBossBloon spawns a boss-type bloon directly at BloonSpawn.
// Since boss behaviors aren't ported yet, these spawn as Normal_Bloon_Branch
// with high layer + special flags.
func spawnBossBloon(inst *engine.Instance, g *engine.Game, bloonup float64) {
	spawns := g.InstanceMgr.FindByObject("BloonSpawn")
	if len(spawns) == 0 {
		return
	}
	sp := spawns[0]

	bloon := g.InstanceMgr.Create("Normal_Bloon_Branch", sp.X, sp.Y)
	if bloon == nil {
		return
	}
	bloon.Vars["boss"] = bloonup

	// Sandbox boss buttons still route through Normal_Bloon_Branch. Give the
	// Mighty MOAB family a MOAB-class layer and explicit sprite so they don't
	// fall back to a red bloon placeholder.
	bossConfigs := map[float64]struct {
		layer      float64
		spriteName string
		shieldHP   float64
	}{
		10021: {layer: 348, spriteName: "Mighty_Moab_spr", shieldHP: 300},
		10022: {layer: 348, spriteName: "Mighty_2", shieldHP: 300},
		10023: {layer: 348, spriteName: "Mighty_Fist_spr", shieldHP: 300},
		10024: {layer: 1248, spriteName: "Mightiest_Moab_Spr", shieldHP: 900},
	}
	if cfg, ok := bossConfigs[bloonup]; ok {
		bpower := getGlobal(g, "bpower")
		if bpower == 0 {
			bpower = 1
		}
		shieldHP := math.Round(cfg.shieldHP * bpower)
		bloon.Vars["bloonlayer"] = cfg.layer
		bloon.Vars["bloonmaxlayer"] = cfg.layer
		bloon.Vars["special_sprite"] = cfg.spriteName
		bloon.Vars["shielded"] = 1.0
		bloon.Vars["shield_hp"] = shieldHP
		bloon.Vars["shield_max"] = shieldHP
		updateBloonSprite(bloon)
		updateBloonDepth(bloon)
		fmt.Printf("[sandbox] spawned boss bloon: id=%d boss=%.0f layer=%.3f sprite=%s shield_hp=%.1f\n",
			bloon.ID, bloonup, getVar(bloon, "bloonlayer"), bloon.SpriteName, getVar(bloon, "shield_hp"))
		return
	}

	bloon.Vars["bloonlayer"] = 1.0
	bloon.Vars["bloonmaxlayer"] = 1.0
	updateBloonSprite(bloon)
	updateBloonDepth(bloon)
	fmt.Printf("[sandbox] spawned fallback boss bloon: id=%d boss=%.0f layer=%.3f sprite=%s\n",
		bloon.ID, bloonup, getVar(bloon, "bloonlayer"), bloon.SpriteName)
}

// configureSenderStack1 sets up stack1 on a sandbox_sender for Bloon_Button spawns.
// multiplier: 1 for single spawn, 20 for mass spawn.
func configureSenderStack1(sender *engine.Instance, bloonup, uptype, multiplier float64, g *engine.Game) {
	bpower := getGlobal(g, "bpower")
	if bpower == 0 {
		bpower = 1
	}
	wsqueeze := getGlobal(g, "wavesqueeze")
	delayMul := 1.0 - wsqueeze*0.5

	sender.Vars["stack1type"] = uptype

	isMass := multiplier > 1

	switch {
	case bloonup < 18:
		// normal bloons
		sender.Vars["stack1layer"] = bloonup
		if isMass {
			sender.Vars["stack1amount"] = 20.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 30

		// shield for fortified variants
		if uptype == 2 {
			sender.Vars["stack1shield"] = math.Round(bloonup * 5 * bpower)
		}
		if uptype == 6 {
			sender.Vars["stack1shield"] = math.Round(bloonup * 10 * bpower)
		}

		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup < 49:
		// ceramic / brick
		sender.Vars["stack1layer"] = bloonup
		if isMass {
			sender.Vars["stack1amount"] = 20.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 30

		if uptype == 2 {
			sender.Vars["stack1shield"] = math.Round(bloonup * 5 * bpower)
		}
		if uptype == 6 {
			sender.Vars["stack1shield"] = math.Round(bloonup * 10 * bpower)
		}

		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 93:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(75 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 10.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 60
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 348:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(300 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 10.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 60
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 1248:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(900 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 5.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 120
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 5248:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(4000 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 5.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 120
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 68.5:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(60 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 10.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 60
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 593:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(500 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 5.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 120
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 351:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(303 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 6.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 100
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 318:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(300 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 5.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 120
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	// Nightmare MOAB-class
	case bloonup == 10068.5:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(10000 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 10.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 60
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 2593:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(1500 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 5.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 120
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 248:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(300 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 10.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 60
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 918:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(3500 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 5.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 120
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 3351:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(3000 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 5.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 120
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	case bloonup == 17248:
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1shield"] = math.Round(15000 * bpower)
		if isMass {
			sender.Vars["stack1amount"] = 5.0
		} else {
			sender.Vars["stack1amount"] = 1.0
		}
		sender.Vars["stack1delay"] = delayMul * 120
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}

	default:
		// fallback for unhandled bloon types
		sender.Vars["stack1layer"] = bloonup
		sender.Vars["stack1amount"] = multiplier
		sender.Vars["stack1delay"] = delayMul * 30
		sender.Alarms[1] = 1
		if isMass {
			sender.Alarms[0] = 700
		} else {
			sender.Alarms[0] = 70
		}
	}
}

// configureSenderStack10 sets up stack10 on a sandbox_sender for Special_Bloon_Button spawns.
func configureSenderStack10(sender *engine.Instance, bloonup, uptype, multiplier float64, g *engine.Game) {
	wsqueeze := getGlobal(g, "wavesqueeze")
	delayMul := 1.0 - wsqueeze*0.5

	sender.Vars["stack10type"] = uptype
	sender.Vars["stack10layer"] = bloonup

	isMass := multiplier > 1
	amount := 1.0
	shield := 0.0

	switch int(uptype) {
	case 1: // Stuffed
		if isMass {
			amount = 20
		}
	case 2: // Ninja
		shield = bloonup * 8
		if isMass {
			amount = 10
		}
	case 3: // Robo
		shield = bloonup * 12
		if isMass {
			amount = 10
		}
	case 4: // Patrol
		shield = 80
		if isMass {
			amount = 10
		}
	case 5: // Barrier
		shield = 10 * bloonup
		if isMass {
			amount = 10
		}
	case 6: // Planetarium
		shield = 30 * bloonup
		if isMass {
			amount = 5
		}
	case 7: // Spectrum
		shield = 250
		if isMass {
			amount = 20
		}
	}

	sender.Vars["stack10amount"] = amount
	sender.Vars["stack10shield"] = shield
	sender.Vars["stack10delay"] = delayMul * 30
	sender.Alarms[10] = 1

	if isMass {
		sender.Alarms[0] = 700
	} else {
		sender.Alarms[0] = 70
	}
}

// ---------------------------------------------------------------------------
// Drawing helpers for bloon icons in the picker panel
// ---------------------------------------------------------------------------

// drawBloonButtonIcon draws the correct bloon icon for Bloon_Button.
func drawBloonButtonIcon(screen *ebiten.Image, g *engine.Game, bloonup float64, uptype int, cx, cy float64) {
	// normal bloon sprites (layer → sprite name)
	normalSprites := map[float64]string{
		1: "Red_Bloon_Spr", 2: "Blue_Bloon_Spr", 3: "Green_Bloon_Spr",
		4: "Yellow_Bloon_Spr", 5: "Pink_Bloon_Spr",
		1.5: "Orange_Bloon_Spr", 2.5: "Cyan_Bloon_Spr", 3.5: "Lime_Bloon_Spr",
		4.5: "Amber_Bloon_Spr", 5.5: "Purple_Bloon_Spr",
		6: "Black_Bloon_Spr", 6.1: "White_Bloon_Spr",
		7: "Zebra_Bloon_Spr", 8: "Rainbow_Bloon_Spr", 8.5: "Prismatic_Bloon_Spr",
		18: "Ceramic_Bloon_Spr", 48: "Brick_Bloon_Spr",
	}

	// try normal bloon sprite
	if spriteName, ok := normalSprites[bloonup]; ok {
		spr := g.AssetManager.GetSprite(spriteName)
		if spr != nil && len(spr.Frames) > 0 {
			frame := uptype
			if frame >= len(spr.Frames) {
				frame = 0
			}
			engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
				cx, cy, 1, 1, 0, 1)
		}
		return
	}

	// MOAB panel sprites (full-size icons)
	moabPanelSprites := map[float64]string{
		93:   "Panel_Mini",
		348:  "Panel_Moab",
		1248: "Panel_BFB",
		5248: "Panel_ZOMG",
		68.5: "Panel_HTA",
		593:  "Panel_BRC",
		351:  "Panel_DDT",
		318:  "New_LPZ_Spr",
	}
	if spriteName, ok := moabPanelSprites[bloonup]; ok {
		spr := g.AssetManager.GetSprite(spriteName)
		if spr != nil && len(spr.Frames) > 0 {
			frame := uptype
			if frame >= len(spr.Frames) {
				frame = 0
			}
			engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
				cx, cy, 1, 1, 0, 1)
		}
		return
	}

	// Nightmare MOAB-class (scaled sprites)
	nightmareSprites := map[float64]struct {
		name  string
		frame int
		scale float64
	}{
		10068.5: {"Prismatic_HTA_Spr", 0, 0.5},
		2593:    {"Mega_BRC_Spr", 0, 0.3},
		3351:    {"Deadly_DDT_Spr", 0, 0.3},
		918:     {"Storm_LPZ_Spr", 1, 0.3},
		248:     {"Rocket_Blimp_Spr", 0, 0.5},
		17248:   {"Party_Blimp_Spr", 0, 0.25},
	}
	if info, ok := nightmareSprites[bloonup]; ok {
		spr := g.AssetManager.GetSprite(info.name)
		if spr != nil && len(spr.Frames) > 0 {
			frame := info.frame
			if frame >= len(spr.Frames) {
				frame = 0
			}
			engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
				cx, cy, info.scale, info.scale, 0, 1)
		}
		return
	}

	// Boss bloon sprites (all scaled to fit the 64x64 button)
	bossSprites := map[float64]struct {
		name  string
		frame int
		scale float64
	}{
		10011: {"Bully_Bloon", 0, 0.5},
		10012: {"Big_Rage_Bloon_spr", 0, 0.5},
		10013: {"Big_Bully_Rage_Spr", 0, 0.5},
		10014: {"Bully_Bloon_Layers", 0, 0.5},
		10015: {"Bully_Bloon_Layers", 1, 0.5},
		10016: {"Bully_Bloon_Layers", 2, 0.5},
		10017: {"Bully_Bloon_Layers", 3, 0.5},
		10021: {"Mighty_Moab_spr", 0, 0.3},
		10022: {"Mighty_2", 0, 0.3},
		10023: {"Mighty_Fist_spr", 0, 0.3},
		10024: {"Mightiest_Moab_Spr", 0, 0.3},
		10041: {"UFO_Bloon_spr", 0, 0.5},
		10042: {"UFO_Yellow", 0, 0.5},
		10043: {"Mothership_Spr", 0, 0.25},
		10051: {"The_Super_Bloon_spr", 0, 0.5},
		10052: {"Super_Bloon_B", 0, 0.5},
		10053: {"Super_Bloon_C", 0, 0.5},
		10054: {"Super_Bloon_2_spr", 0, 0.5},
		10061: {"Small_Mother_Bloon", 0, 0.3},
		10062: {"The_Mother_spr", 0, 0.2},
		10063: {"Monster_Mom_spr", 0, 0.3},
		10064: {"Monster_Queen_spr", 0, 0.2},
		10071: {"Big_Laughing_Bloon", 0, 0.5},
		10072: {"Emoji_Heart", 0, 0.5},
		10073: {"Emoji_DDT_spr", 0, 0.4},
		10074: {"Jumbo_Emoji", 0, 0.3},
		10081: {"Clown_Bloon_spr", 0, 0.5},
		10082: {"Joker_Bloon_spr", 0, 0.5},
		10083: {"Demon_Clown_spr", 0, 0.4},
		10084: {"Demon_T2_dpr", 0, 0.4},
		10091: {"Blooming_Bloon_spr", 119, 0.5},
		10092: {"Sprouter", 119, 0.5},
		10093: {"Monkey_Eater_spr", 100, 0.3},
		10094: {"Spring_Bloon_spr", 119, 0.3},
		10111: {"tDoM_spr", 0, 0.3},
		10112: {"tDoM_T3", 0, 0.3},
	}
	if info, ok := bossSprites[bloonup]; ok {
		spr := g.AssetManager.GetSprite(info.name)
		if spr != nil && len(spr.Frames) > 0 {
			frame := info.frame
			if frame >= len(spr.Frames) {
				frame = 0
			}
			engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
				cx, cy, info.scale, info.scale, 0, 1)
		}
		return
	}

	// fallback: draw a placeholder colored rect
	drawRect(screen, cx-10, cy-10, 20, 20, [3]uint8{200, 100, 100})
}

// drawSpecialBloonIcon draws nightmare bloon icon based on uptype.
func drawSpecialBloonIcon(screen *ebiten.Image, g *engine.Game, uptype int, bloonup, cx, cy float64) {
	spriteNames := map[int]string{
		1: "Stuffed_Bloon_Spr",
		2: "Ninja_Bloon_Spr",
		3: "Robo_Bloon_Spr",
		4: "Patrol_Bloon_Spr",
		5: "Barrier_Bloon_Spr",
		6: "Planetarium_Bloon_Spr",
		7: "Spectrum_Bloon_Spr",
	}

	spriteName, ok := spriteNames[uptype]
	if !ok {
		drawRect(screen, cx-10, cy-10, 20, 20, [3]uint8{150, 50, 150})
		return
	}

	spr := g.AssetManager.GetSprite(spriteName)
	if spr == nil || len(spr.Frames) == 0 {
		drawRect(screen, cx-10, cy-10, 20, 20, [3]uint8{150, 50, 150})
		return
	}

	// frame selection matches original GML
	frame := 0
	switch uptype {
	case 1, 2, 3, 5: // Bloonup-1 frame index
		frame = int(bloonup) - 1
	case 4: // always frame 0
		frame = 0
	case 6: // Bloonup/10
		frame = int(bloonup / 10)
	case 7: // always frame 0
		frame = 0
	}
	if frame < 0 {
		frame = 0
	}
	if frame >= len(spr.Frames) {
		frame = 0
	}

	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		cx, cy, 1, 1, 0, 1)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// unlockAllTowers sets every tower lock global to 1.0.
func unlockAllTowers(g *engine.Game) {
	locks := []string{
		"DMlock", "TSlock", "BMlock", "SnMlock", "NMlock", "BClock",
		"MSlock", "CTlock", "MBlock", "MElock", "GGlock", "IMlock",
		"MAlock", "BChlock", "MAplock", "MAllock", "MVlock", "BTlock",
		"DGlock", "MLlock", "SFlock", "HPlock", "PMlock", "SuMlock", "Derlock",
	}
	for _, l := range locks {
		g.GlobalVars[l] = 1.0
	}
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func RegisterSandboxBehaviors(im *engine.InstanceManager) {
	im.RegisterBehavior("Sanbox_Bar", func() engine.InstanceBehavior { return &SanboxBarBehavior{} })
	im.RegisterBehavior("Sandbox_Settings", func() engine.InstanceBehavior { return &SandboxSettingsBehavior{} })
	im.RegisterBehavior("Sandbox_Go", func() engine.InstanceBehavior { return &SandboxGoBehavior{} })
	im.RegisterBehavior("Bloon_Button", func() engine.InstanceBehavior { return &BloonButtonBehavior{} })
	im.RegisterBehavior("Special_Bloon_Button", func() engine.InstanceBehavior { return &SpecialBloonButtonBehavior{} })
	im.RegisterBehavior("Bloon_Up", func() engine.InstanceBehavior { return &BloonUpBehavior{} })
	im.RegisterBehavior("Bloon_Down", func() engine.InstanceBehavior { return &BloonDownBehavior{} })
	im.RegisterBehavior("sandbox_sender", func() engine.InstanceBehavior { return &SandboxSenderBehavior{} })
}
