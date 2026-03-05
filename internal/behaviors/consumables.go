package behaviors

import (
	"math/rand"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

// Consumable items: Road Spikes and Pineapple Bombs
// These are single-use items placed directly on the track (no block placement).
// Spikes persist in placement mode for multiple placements.
// Pineapples are single placement (exit placement after one).

const (
	consumableSpikeSel     = 1001.0 // towerselect value for spike placement
	consumablePineappleSel = 1002.0 // towerselect value for pineapple placement
	consumableCost         = 30.0   // $30 per item
)

// isConsumableSelected returns true if the current towerselect is a consumable item.
// Used by BlockBehavior to hide land tiles during consumable placement.
func isConsumableSelected(g *engine.Game) bool {
	sel := getGlobal(g, "towerselect")
	return sel >= 1000
}

// ─────────────────────────────────────────────────────────────────────────────
// Buy panel helpers
// ─────────────────────────────────────────────────────────────────────────────

// drawConsumablePanel draws a buy panel sprite directly to screen (HUD area,
// no viewport transform).
func drawConsumablePanel(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := spr.Frames[0]
	engine.DrawSpriteExt(screen, frame, spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, 0, inst.ImageAlpha)
}

// enterConsumablePlacement sets up placement mode for a consumable item.
func enterConsumablePlacement(inst *engine.Instance, g *engine.Game, towerSel float64, placeObj string) {
	if !inst.Visible {
		return
	}
	money := getGlobal(g, "money")
	if money < consumableCost {
		return
	}
	// cancel any existing placement or tower selection
	if getGlobal(g, "towerplace") == 1 {
		cancelTowerUI(g)
	}
	deselectAllTowers(g)

	g.GlobalVars["towerselect"] = towerSel
	g.GlobalVars["towerplace"] = 1.0

	// create placement ghost at mouse position
	mx, my := g.GetMouseRoomPos()
	place := g.InstanceMgr.Create(placeObj, mx, my)
	place.Depth = -100
}

// ─────────────────────────────────────────────────────────────────────────────
// Spike_Pile_buy — HUD buy panel for Road Spikes ($30)
// Placed in room at (880, 384), sprite "Road_Spikes_Panel" (64×64, origin 0,0)
// ─────────────────────────────────────────────────────────────────────────────

type SpikePileBuyBehavior struct {
	engine.DefaultBehavior
}

func (b *SpikePileBuyBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Road_Spikes_Panel"
	inst.Depth = -20
	inst.Visible = true
}

func (b *SpikePileBuyBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	drawConsumablePanel(inst, screen, g)
}

func (b *SpikePileBuyBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	enterConsumablePlacement(inst, g, consumableSpikeSel, "Pile_Place")
}

// ─────────────────────────────────────────────────────────────────────────────
// Pineapples — HUD buy panel for Pineapple Bombs ($30)
// Placed in room at (944, 384), sprite "Pineapple_Panel" (64×64, origin 0,0)
// ─────────────────────────────────────────────────────────────────────────────

type PineappleBuyBehavior struct {
	engine.DefaultBehavior
}

func (b *PineappleBuyBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Pineapple_Panel"
	inst.Depth = -20
	inst.Visible = true
}

func (b *PineappleBuyBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	drawConsumablePanel(inst, screen, g)
}

func (b *PineappleBuyBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	enterConsumablePlacement(inst, g, consumablePineappleSel, "Pineapple_Place")
}

// ─────────────────────────────────────────────────────────────────────────────
// Pile_Place — placement ghost for Road Spikes
// Follows mouse cursor, semi-transparent preview.
// Placement happens in Step (which runs after mouse events), so the cancel bar
// can intercept clicks without accidentally placing a spike.
// Stays alive for multiple placements until cancelled or out of money.
// ─────────────────────────────────────────────────────────────────────────────

type PilePlaceBehavior struct {
	engine.DefaultBehavior
}

func (b *PilePlaceBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Big_Spikes"
	inst.ImageSpeed = 0
	inst.ImageAlpha = 0.5
	inst.Depth = -100
}

func (b *PilePlaceBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// self-destruct if placement mode was cancelled
	if getGlobal(g, "towerselect") != consumableSpikeSel || getGlobal(g, "towerplace") != 1 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}

	// follow mouse
	mx, my := g.GetMouseRoomPos()
	inst.X = mx
	inst.Y = my

	// check for left click to place (Step runs after mouse events,
	// so cancelTowerUI from the X-bar has already been applied if clicked)
	if g.InputMgr.MouseLeftPressed() {
		// only in playfield area (left of HUD)
		if mx >= 0 && my >= 0 && mx < 864 && my < 480 {
			money := getGlobal(g, "money")
			if money >= consumableCost {
				g.InstanceMgr.Create("Spike_Pile", mx, my)
				g.GlobalVars["money"] = money - consumableCost
				g.AudioMgr.Play("Tower_Place")

				// exit placement if can't afford another
				if money-consumableCost < consumableCost {
					cancelTowerUI(g)
				}
			}
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Pineapple_Place — placement ghost for Pineapple Bombs
// Single-use: places one pineapple then exits placement mode.
// ─────────────────────────────────────────────────────────────────────────────

type PineapplePlaceBehavior struct {
	engine.DefaultBehavior
}

func (b *PineapplePlaceBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Pineapple_Spr"
	inst.ImageAlpha = 0.5
	inst.Depth = -100
}

func (b *PineapplePlaceBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "towerselect") != consumablePineappleSel || getGlobal(g, "towerplace") != 1 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}

	mx, my := g.GetMouseRoomPos()
	inst.X = mx
	inst.Y = my

	if g.InputMgr.MouseLeftPressed() {
		if mx >= 0 && my >= 0 && mx < 864 && my < 480 {
			money := getGlobal(g, "money")
			if money >= consumableCost {
				g.InstanceMgr.Create("Grilled_Pineapple", mx, my)
				g.GlobalVars["money"] = money - consumableCost
				g.AudioMgr.Play("Tower_Place")

				// single-use placement — exit immediately
				cancelTowerUI(g)
			}
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Spike_Pile — actual road spike placed on the track
// PP=12, LP=1, camopop=1, leadpop=0
// Shrinks visually as pierce points deplete.
// Self-destructs when PP reaches 0 or after a very long timeout.
// ─────────────────────────────────────────────────────────────────────────────

type SpikePileBehavior struct {
	engine.DefaultBehavior
}

func (b *SpikePileBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Big_Spikes" // 48x48, 6 frames for visual degradation
	inst.ImageSpeed = 0            // manual frame control
	inst.ImageIndex = 0            // full spikes
	inst.Depth = 3                 // behind bloons(0), blocks(-1); in front of water(4)
	inst.Vars["LP"] = 1.0
	inst.Vars["PP"] = 12.0
	inst.Vars["leadpop"] = 0.0
	inst.Vars["camopop"] = 1.0
	inst.ImageAngle = rand.Float64() * 360
	inst.Alarms[0] = 90000 // ~25 min lifetime at 60fps
}

func (b *SpikePileBehavior) Step(inst *engine.Instance, g *engine.Game) {
	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}

	// check collision with bloons in a small radius
	projectileHitBloons(inst, g, 20)

	// visual degradation: 6-frame Big_Spikes sprite (GML original)
	// frame 0 = full, frame 5 = nearly gone
	pp = getVar(inst, "PP") // re-read after potential hits
	if pp > 10 {
		inst.ImageIndex = 0
	} else if pp > 8 {
		inst.ImageIndex = 1
	} else if pp > 6 {
		inst.ImageIndex = 2
	} else if pp > 4 {
		inst.ImageIndex = 3
	} else if pp > 2 {
		inst.ImageIndex = 4
	} else {
		inst.ImageIndex = 5
	}
}

func (b *SpikePileBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Grilled_Pineapple — pineapple bomb placed on the track
// Sits for 30 frames (0.5s) then explodes into a Medium Explosion.
// ─────────────────────────────────────────────────────────────────────────────

type GrilledPineappleBehavior struct {
	engine.DefaultBehavior
}

func (b *GrilledPineappleBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Pineapple_Spr"
	inst.Depth = 3      // same layer as spikes, behind bloons
	inst.Alarms[0] = 30 // explode after 0.5 seconds
	inst.ImageAngle = rand.Float64() * 360
}

func (b *GrilledPineappleBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		// create medium explosion (reuses SmallExplosionBehavior with overridden params)
		explosion := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
		explosion.Vars["LP"] = 1.0
		explosion.Vars["PP"] = 40.0
		explosion.Vars["leadpop"] = 1.0
		explosion.Vars["camopop"] = 1.0
		explosion.Vars["explode_radius"] = 60.0
		explosion.SpriteName = "Medium_Explosion_Spr"
		explosion.ImageXScale = 1.1
		explosion.ImageYScale = 1.1
		explosion.Alarms[1] = 8 // destroy after 8 frames
		g.AudioMgr.Play("Small_Boom")

		g.InstanceMgr.Destroy(inst.ID)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Registration
// ─────────────────────────────────────────────────────────────────────────────

func RegisterConsumableBehaviors(im *engine.InstanceManager) {
	// buy panels (HUD)
	im.RegisterBehavior("Spike_Pile_buy", func() engine.InstanceBehavior { return &SpikePileBuyBehavior{} })
	im.RegisterBehavior("Pineapples", func() engine.InstanceBehavior { return &PineappleBuyBehavior{} })

	// placement ghosts
	im.RegisterBehavior("Pile_Place", func() engine.InstanceBehavior { return &PilePlaceBehavior{} })
	im.RegisterBehavior("Pineapple_Place", func() engine.InstanceBehavior { return &PineapplePlaceBehavior{} })

	// actual consumable items
	im.RegisterBehavior("Spike_Pile", func() engine.InstanceBehavior { return &SpikePileBehavior{} })
	im.RegisterBehavior("Grilled_Pineapple", func() engine.InstanceBehavior { return &GrilledPineappleBehavior{} })
}
