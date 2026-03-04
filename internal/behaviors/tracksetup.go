package behaviors

import (
	"btdx/internal/engine"
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

// pendingModTooltip holds the modifier hover tooltip to draw this frame.
// It is populated by modifierToggle.Draw and consumed by scoreAndTimeDraw.Draw.
var pendingModTooltip struct {
	active      bool
	mx, my      float64
	description string // may contain '#' as a line-break
	pointsText  string
}

// track_Setup_II room behaviors
// mode selection, wave counts, modifiers, and Play button

// normal_Mode_enable — selects normal mode
type NormalModeEnable struct {
	engine.DefaultBehavior
}

func (b *NormalModeEnable) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["normalmodeselect"] = 0.0
}

func (b *NormalModeEnable) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["normalmodeselect"] = 1.0
	g.GlobalVars["impoppablemodeselect"] = 0.0
	g.GlobalVars["nightmaremodeselect"] = 0.0
	g.GlobalVars["timemodeselect"] = 0.0
}

func (b *NormalModeEnable) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := 0
	if getGlobal(g, "normalmodeselect") == 1 {
		frame = 1
	}
	if frame >= len(spr.Frames) {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// impoppable_Mode_enable — selects impoppable mode (rank >= 30)
type ImpoppableModeEnable struct {
	engine.DefaultBehavior
}

func (b *ImpoppableModeEnable) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "rank") >= 30 {
		g.GlobalVars["normalmodeselect"] = 0.0
		g.GlobalVars["impoppablemodeselect"] = 1.0
		g.GlobalVars["nightmaremodeselect"] = 0.0
		g.GlobalVars["timemodeselect"] = 0.0
	}
}

func (b *ImpoppableModeEnable) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := 0
	if getGlobal(g, "impoppablemodeselect") == 1 {
		frame = 1
	}
	if getGlobal(g, "rank") < 30 {
		frame = 2
	}
	if frame >= len(spr.Frames) {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// nightmare_Mode_enable — selects nightmare mode (rank >= 40)
type NightmareModeEnable struct {
	engine.DefaultBehavior
}

func (b *NightmareModeEnable) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "rank") >= 40 {
		g.GlobalVars["normalmodeselect"] = 0.0
		g.GlobalVars["impoppablemodeselect"] = 0.0
		g.GlobalVars["nightmaremodeselect"] = 1.0
		g.GlobalVars["timemodeselect"] = 0.0
	}
}

func (b *NightmareModeEnable) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := 0
	if getGlobal(g, "nightmaremodeselect") == 1 {
		frame = 1
	}
	if getGlobal(g, "rank") < 40 {
		frame = 2
	}
	if frame >= len(spr.Frames) {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// time_Trial_enable — selects time trial mode
type TimeTrialEnable struct {
	engine.DefaultBehavior
}

func (b *TimeTrialEnable) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0
}

func (b *TimeTrialEnable) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "rank") >= 20 {
		g.GlobalVars["normalmodeselect"] = 0.0
		g.GlobalVars["impoppablemodeselect"] = 0.0
		g.GlobalVars["nightmaremodeselect"] = 0.0
		g.GlobalVars["timemodeselect"] = 1.0
	}
}

func (b *TimeTrialEnable) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := 0
	if getGlobal(g, "timemodeselect") == 1 {
		frame = 1
	}
	if getGlobal(g, "rank") < 20 {
		frame = 2
	}
	if frame >= len(spr.Frames) {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// play_bar — starts the game, routes to the correct room
type PlayBar struct {
	engine.DefaultBehavior
}

func (b *PlayBar) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["challenge"] = 0.0
	g.GlobalVars["pointmultiplier"] = 1.0
	g.GlobalVars["towerlimit"] = 1000000.0

	// reset modifier flags
	g.GlobalVars["sixtowers"] = 0.0
	g.GlobalVars["randomtowers"] = 0.0
	g.GlobalVars["wavesqueeze"] = 0.0
	g.GlobalVars["waveskip"] = 0.0
	g.GlobalVars["strongerbloons"] = 0.0
	g.GlobalVars["fasterbloons"] = 0.0
	g.GlobalVars["noliveslost"] = 0.0

	// bloon info for new players
	if getGlobal(g, "rank") < 20 {
		g.GlobalVars["blooninfo"] = 1.0
	} else {
		g.GlobalVars["blooninfo"] = 0.0
	}

	// reset all tower locks
	towerLockNames := []string{
		"DMlock", "TSlock", "BMlock", "SnMlock", "NMlock", "BClock",
		"MSlock", "CTlock", "MBlock", "MElock", "GGlock", "IMlock",
		"MAlock", "BChlock", "MAplock", "MAllock", "MVlock", "BTlock",
		"DGlock", "MLlock", "SFlock", "HPlock", "PMlock", "SuMlock", "Derlock",
	}
	for _, ln := range towerLockNames {
		g.GlobalVars[ln] = 0.0
	}
}

func (b *PlayBar) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	ts := getGlobal(g, "trackselect")

	// normal mode room routing (default)
	trackRooms := map[float64]string{
		1: "Monkey_Meadows_Norm", 2: "Bloon_Oasis_Norm",
		3: "Swamp_Spirals_Norm", 4: "Monkey_Fort_Norm",
		5: "Monkey_Town_Docks_Norm", 6: "Conveyor_Belts_Norm",
		7: "The_Depths_Norm", 8: "Sun_Stone_Norm",
		9: "Shade_Woods_Norm", 10: "Minecarts_Norm",
		11: "Crimson_Creek_Norm", 12: "Xtreme_Park_Norm",
		13: "Portal_Lab_Norm", 14: "Omega_River_Norm",
		15: "Space_Portals_Norm", 17: "Bloonlight_Throwback_Norm",
		18: "Bloon_Circles_X_Norm", 19: "Autumn_Acres_Norm",
		20: "Graveyard_Norm", 21: "Village_Defense_Norm",
		22: "Circuit_Norm", 23: "Grand_Canyon_Norm",
		24: "Bloonside_River_Norm", 25: "Cotton_Fields_Norm",
		27: "Rubber_Rug_Norm", 28: "Frozen_Lake_Norm",
		29: "Sky_Battles_Norm", 30: "Lava_Stream_Norm",
		31: "Ravine_River_Norm", 32: "Peaceful_Bridge_Norm",
	}

	// time trial rooms
	trackTimeRooms := map[float64]string{
		1: "Monkey_Meadows_Time", 2: "Bloon_Oasis_Time",
		3: "Swamp_Spirals_Time", 4: "Monkey_Fort_Time",
		5: "Monkey_Town_Docks_Time", 6: "Conveyor_Belts_Time",
		7: "The_Depths_Time", 8: "Sun_Stone_Time",
		9: "Shade_Woods_Time", 10: "Minecarts_Time",
		11: "Crimson_Creek_Time", 12: "Xtreme_Park_Time",
		13: "Portal_Lab_Time", 14: "Omega_River_Time",
		15: "Space_Portals_Time", 17: "Bloonlight_Throwback_Time",
		18: "Bloon_Circles_X_Time", 19: "Autumn_Acres_Time",
		20: "Graveyard_Time", 21: "Village_Defense_Time",
		22: "Circuit_Time", 23: "Grand_Canyon_Time",
		24: "Bloonside_River_Time", 27: "Rubber_Rug_Time",
		28: "Frozen_Lake_Time", 29: "Sky_Battles_Time",
		30: "Lava_Stream_Time", 31: "Ravine_River_Time",
		32: "Peaceful_Bridge_Time",
	}

	// point multipliers for mode
	if getGlobal(g, "impoppablemodeselect") == 1 {
		g.GlobalVars["pointmultiplier"] = 1.5
	}
	if getGlobal(g, "nightmaremodeselect") == 1 {
		g.GlobalVars["pointmultiplier"] = 2.0
	}

	// modifier multipliers
	pm := getGlobal(g, "pointmultiplier")
	if getGlobal(g, "sixtowers") == 1 {
		pm *= 1.6
		g.GlobalVars["towerlimit"] = 6.0
	}
	if getGlobal(g, "wavesqueeze") == 1 {
		pm *= 1.6
	}
	if getGlobal(g, "waveskip") == 1 {
		pm *= 2.0
	}
	if getGlobal(g, "strongerbloons") == 1 {
		pm *= 1.35
	}
	if getGlobal(g, "fasterbloons") == 1 {
		pm *= 1.45
	}
	if getGlobal(g, "noliveslost") == 1 {
		pm *= 1.25
	}
	g.GlobalVars["pointmultiplier"] = pm

	// unlock towers based on rank
	b.unlockTowers(g)

	// route to room based on mode
	if getGlobal(g, "timemodeselect") == 1 {
		if room, ok := trackTimeRooms[ts]; ok {
			g.RequestRoomGoto(room)
		}
	} else {
		if room, ok := trackRooms[ts]; ok {
			g.RequestRoomGoto(room)
		}
	}
}

func (b *PlayBar) unlockTowers(g *engine.Game) {
	unlockTowersForRank(g)
}

// unlockTowersForRank sets tower lock globals based on rank.
// rank 0: Dart + Tack always available; higher ranks unlock more towers in pairs.
// Exported as a package-level function so both PlayBar and BountyGoBehavior can call it.
func unlockTowersForRank(g *engine.Game) {
	rank := getGlobal(g, "rank")

	// base towers always unlocked
	g.GlobalVars["DMlock"] = 1.0 // dart Monkey
	g.GlobalVars["TSlock"] = 1.0 // tack Shooter

	// rank-based unlocks (pairs)
	if rank >= 1 {
		g.GlobalVars["BMlock"] = 1.0  // boomerang
		g.GlobalVars["SnMlock"] = 1.0 // sniper
	}
	if rank >= 2 {
		g.GlobalVars["NMlock"] = 1.0 // ninja
		g.GlobalVars["BClock"] = 1.0 // bomb Cannon
	}
	if rank >= 3 {
		g.GlobalVars["MSlock"] = 1.0 // monkey Sub
		g.GlobalVars["CTlock"] = 1.0 // charge Tower
	}
	if rank >= 4 {
		g.GlobalVars["GGlock"] = 1.0 // glue Gunner
		g.GlobalVars["IMlock"] = 1.0 // ice Monkey
	}
	if rank >= 5 {
		g.GlobalVars["MElock"] = 1.0 // monkey Engineer
		g.GlobalVars["MBlock"] = 1.0 // monkey Buccaneer
	}
	if rank >= 6 {
		g.GlobalVars["MAlock"] = 1.0  // monkey Ace
		g.GlobalVars["BChlock"] = 1.0 // bloonchipper
	}
	if rank >= 7 {
		g.GlobalVars["MAllock"] = 1.0 // monkey Alchemist
		g.GlobalVars["MAplock"] = 1.0 // monkey Apprentice
	}
	if rank >= 8 {
		g.GlobalVars["BTlock"] = 1.0 // banana Tree
		g.GlobalVars["MVlock"] = 1.0 // monkey Village
	}
	if rank >= 9 {
		g.GlobalVars["MLlock"] = 1.0 // mortar Launcher
		g.GlobalVars["DGlock"] = 1.0 // dartling Gunner
	}
	if rank >= 10 {
		g.GlobalVars["SFlock"] = 1.0 // spike Factory
		g.GlobalVars["HPlock"] = 1.0 // heli Pilot
	}
	if rank >= 11 {
		g.GlobalVars["PMlock"] = 1.0  // plasma Monkey
		g.GlobalVars["SuMlock"] = 1.0 // super Monkey
	}
}

// modifier toggle behaviors (Six_Towers, Wave_Squeeze, etc.)
// each toggles a global between 0 and 1 on click
type modifierToggle struct {
	engine.DefaultBehavior
	GlobalKey   string
	Description string // may contain '#' as a line-break (GameMaker convention)
	PointsText  string // e.g. "x1.6 Points"
}

func (b *modifierToggle) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	cur := getGlobal(g, b.GlobalKey)
	if cur == 0 {
		g.GlobalVars[b.GlobalKey] = 1.0
	} else {
		g.GlobalVars[b.GlobalKey] = 0.0
	}
}

func (b *modifierToggle) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := 0
	active := getGlobal(g, b.GlobalKey) == 1
	hover := g.IsMouseOverInstance(inst)

	if active {
		frame = 1
		// if sprite has dedicated hover-selected frame, use it
		if hover && len(spr.Frames) >= 4 {
			frame = 3
		}
	} else if hover {
		// if sprite has dedicated hover frame, use it
		if len(spr.Frames) >= 3 {
			frame = 2
		} else if len(spr.Frames) >= 2 {
			// fallback: highlight with selected frame when only 2 frames exist
			frame = 1
		}
	}
	if frame >= len(spr.Frames) {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)

	if hover {
		// Store tooltip data to be drawn this frame by the top-layer draw object.
		pendingModTooltip.active = true
		pendingModTooltip.mx = float64(g.InputMgr.MouseX)
		pendingModTooltip.my = float64(g.InputMgr.MouseY)
		pendingModTooltip.description = b.Description
		pendingModTooltip.pointsText = b.PointsText
	}
}

// scoreAndTimeDraw — Score_and_Time_Draw object in Track_Setup_II.
// Draws nothing itself but forces depth = -1001 so its Draw runs last in the
// room, letting it render the modifier hover tooltip on top of everything else.
type scoreAndTimeDraw struct {
	engine.DefaultBehavior
}

func (b *scoreAndTimeDraw) Create(inst *engine.Instance, g *engine.Game) {
	inst.Depth = -1001
}

func (b *scoreAndTimeDraw) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if !pendingModTooltip.active {
		return
	}
	// Consume the pending tooltip so it clears after one frame.
	p := pendingModTooltip
	pendingModTooltip.active = false

	drawModifierTooltip(screen, g, p.mx, p.my, p.description, p.pointsText)
}

// drawModifierTooltip renders the hover tooltip for a modifier button using
// the same Tower_Info_Panel_Spr sprite used by the in-game tower buy tooltip.
// Description may contain '#' as a line-break (GameMaker convention).
// When rank < 20 a "Rank 20 Needed" message is shown instead.
func drawModifierTooltip(screen *ebiten.Image, g *engine.Game, mx, my float64, description, pointsText string) {
	var line1, line2, bottom string
	if getGlobal(g, "rank") < 20 {
		line1 = "Rank 20 Needed"
	} else {
		parts := strings.SplitN(description, "#", 2)
		line1 = parts[0]
		if len(parts) > 1 {
			line2 = parts[1]
		}
		bottom = pointsText
	}

	const panelScale = 1.3
	if tip := g.AssetManager.GetSprite("Tower_Info_Panel_Spr"); tip != nil && len(tip.Frames) > 0 {
		engine.DrawSpriteExt(screen, tip.Frames[0], tip.XOrigin, tip.YOrigin, mx, my, panelScale, panelScale, 0, 1)
	}

	// Text positions scaled proportionally from the sprite origin (125, 16).
	if line2 != "" {
		drawHUDTextSmall(screen, g, line1, mx-155, my-14, hudColorBlack)
		drawHUDTextSmall(screen, g, line2, mx-155, my-2, hudColorBlack)
	} else {
		drawHUDTextSmall(screen, g, line1, mx-155, my-8, hudColorBlack)
	}
	if bottom != "" {
		drawHUDTextSmall(screen, g, bottom, mx-78, my+21, hudColorBlack)
	}
}

// wave count arrow displays (Norm_40, Impo_35, Night_30, etc.)
// each shows locked/unlocked frame based on track milestone progress
type waveCountBehavior struct {
	engine.DefaultBehavior
	stoneKey      string // "milestone", "hardstone", "nightstone"
	threshold     int    // 1-4 (which milestone tier)
	lockedFrame   int    // 0 for normal, 2 for impoppable, 4 for nightmare
	unlockedFrame int    // 1 for normal, 3 for impoppable, 5 for nightmare
}

func (b *waveCountBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0
}

func (b *waveCountBehavior) Step(inst *engine.Instance, g *engine.Game) {
	ts := int(getGlobal(g, "trackselect"))
	key := fmt.Sprintf("track%d%s", ts, b.stoneKey)
	if getGlobal(g, key) >= float64(b.threshold) {
		inst.ImageIndex = float64(b.unlockedFrame)
	} else {
		inst.ImageIndex = float64(b.lockedFrame)
	}
}

// freeplay displays — show locked/unlocked/gold based on milestones
// freeplay_Spr has 9 frames: normal 0-2, impoppable 3-5, nightmare 6-8
type freeplayBehavior struct {
	engine.DefaultBehavior
	stoneKey  string // "milestone" or "hardstone"
	waveKey   string // "bestwave" or "besthardwave"
	baseFrame int    // 0, 3, or 6
}

func (b *freeplayBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0
}

func (b *freeplayBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	ts := int(getGlobal(g, "trackselect"))
	stoneK := fmt.Sprintf("track%d%s", ts, b.stoneKey)
	waveK := fmt.Sprintf("track%d%s", ts, b.waveKey)
	frame := b.baseFrame // locked
	if getGlobal(g, stoneK) >= 5 {
		frame = b.baseFrame + 1 // unlocked
		if getGlobal(g, waveK) >= 100 {
			frame = b.baseFrame + 2 // gold
		}
	}
	if frame >= len(spr.Frames) {
		frame = b.baseFrame
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// time trial achievement displays (8min, 6.5min, 5min bars)
// frame 0 = not achieved, frame 1 = achieved
type timeBehavior struct {
	engine.DefaultBehavior
	threshold float64 // 480, 390, or 300 seconds
}

func (b *timeBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0
}

func (b *timeBehavior) Step(inst *engine.Instance, g *engine.Game) {
	ts := int(getGlobal(g, "trackselect"))
	key := fmt.Sprintf("t%d", ts)
	t := getGlobal(g, key)
	// t > 0 means a time was recorded; t < threshold means achievement earned
	if t > 0 && t < b.threshold {
		inst.ImageIndex = 1
	} else {
		inst.ImageIndex = 0
	}
}

// bloon_Info toggle — shows/hides bloon info panel
// frame 0 = off, frame 1 = on
type bloonInfoBehavior struct {
	engine.DefaultBehavior
}

func (b *bloonInfoBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0
}

func (b *bloonInfoBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageIndex = getGlobal(g, "blooninfo")
}

func (b *bloonInfoBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "blooninfo") == 0 {
		g.GlobalVars["blooninfo"] = 1.0
	} else {
		g.GlobalVars["blooninfo"] = 0.0
	}
}

// registerTrackSetupBehaviors registers all Track_Setup_II behaviors
func RegisterTrackSetupBehaviors(im *engine.InstanceManager) {
	// mode selectors
	im.RegisterBehavior("Normal_Mode_enable", func() engine.InstanceBehavior { return &NormalModeEnable{} })
	im.RegisterBehavior("Impoppable_Mode_enable", func() engine.InstanceBehavior { return &ImpoppableModeEnable{} })
	im.RegisterBehavior("Nightmare_Mode_enable", func() engine.InstanceBehavior { return &NightmareModeEnable{} })
	im.RegisterBehavior("Time_Trial_enable", func() engine.InstanceBehavior { return &TimeTrialEnable{} })

	// play button
	im.RegisterBehavior("Play_bar", func() engine.InstanceBehavior { return &PlayBar{} })

	// wave count arrow displays
	waveConfigs := []struct {
		name          string
		stoneKey      string
		threshold     int
		lockedFrame   int
		unlockedFrame int
	}{
		{"Norm_40", "milestone", 1, 0, 1},
		{"Norm_60", "milestone", 2, 0, 1},
		{"Norm_75", "milestone", 3, 0, 1},
		{"Norm_90", "milestone", 4, 0, 1},
		{"Impo_35", "hardstone", 1, 2, 3},
		{"Impo_55", "hardstone", 2, 2, 3},
		{"Impo_70", "hardstone", 3, 2, 3},
		{"Impo_85", "hardstone", 4, 2, 3},
		{"Night_30", "nightstone", 1, 4, 5},
		{"Night_40", "nightstone", 2, 4, 5},
		{"Night_50", "nightstone", 3, 4, 5},
		{"Night_60", "nightstone", 4, 4, 5},
	}
	for _, cfg := range waveConfigs {
		c := cfg
		im.RegisterBehavior(c.name, func() engine.InstanceBehavior {
			return &waveCountBehavior{
				stoneKey:      c.stoneKey,
				threshold:     c.threshold,
				lockedFrame:   c.lockedFrame,
				unlockedFrame: c.unlockedFrame,
			}
		})
	}

	// freeplay displays
	im.RegisterBehavior("Freeplay", func() engine.InstanceBehavior {
		return &freeplayBehavior{stoneKey: "milestone", waveKey: "bestwave", baseFrame: 0}
	})
	im.RegisterBehavior("Hard_Freeplay", func() engine.InstanceBehavior {
		return &freeplayBehavior{stoneKey: "hardstone", waveKey: "besthardwave", baseFrame: 3}
	})
	im.RegisterBehavior("Night_Freeplay", func() engine.InstanceBehavior {
		return &freeplayBehavior{stoneKey: "hardstone", waveKey: "besthardwave", baseFrame: 6}
	})

	// time trial displays
	im.RegisterBehavior("eightmins", func() engine.InstanceBehavior {
		return &timeBehavior{threshold: 480}
	})
	im.RegisterBehavior("sixandahalfmins", func() engine.InstanceBehavior {
		return &timeBehavior{threshold: 390}
	})
	im.RegisterBehavior("fivemins", func() engine.InstanceBehavior {
		return &timeBehavior{threshold: 300}
	})

	// bloon info toggle
	im.RegisterBehavior("Bloon_Info", func() engine.InstanceBehavior {
		return &bloonInfoBehavior{}
	})

	// top-layer draw object — renders modifier hover tooltip above all other UI
	im.RegisterBehavior("Score_and_Time_Draw", func() engine.InstanceBehavior {
		return &scoreAndTimeDraw{}
	})

	// modifier toggles (description uses '#' as a line-break, matching the original GML)
	type modCfg struct {
		globalKey   string
		description string
		pointsText  string
	}
	modifiers := map[string]modCfg{
		"Six_Towers":      {"sixtowers", "Six Towers Only", "x1.6 Points"},
		"Random_Towers":   {"randomtowers", "Random#Selection", "x2.2 Points"},
		"Wave_Squeeze":    {"wavesqueeze", "Shorter Waves", "x1.6 Points"},
		"Wave_Skip":       {"waveskip", "Skips Waves#Sometimes", "x2 Points"},
		"Faster_Bloons":   {"fasterbloons", "Bloons move#faster", "x1.45 Points"},
		"Stronger_Bloons": {"strongerbloons", "Big Bloons#have more hp", "x1.35 Points"},
		"No_Lives_Lost":   {"noliveslost", "1 Life Only", "x1.25 Points"},
	}
	for objName, cfg := range modifiers {
		c := cfg // capture
		im.RegisterBehavior(objName, func() engine.InstanceBehavior {
			return &modifierToggle{
				GlobalKey:   c.globalKey,
				Description: c.description,
				PointsText:  c.pointsText,
			}
		})
	}
}
