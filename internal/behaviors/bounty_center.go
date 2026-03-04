package behaviors

import (
	"fmt"
	"image/color"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	etext "github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

// ---------------------------------------------------------------------------
// Bounty Center room (Bloons_Bounty_Center)
// Hub where players select boss bloon fights across 4 difficulty tiers.
// ---------------------------------------------------------------------------

// allCageEnableKeys lists every cage enable global in the order they appear.
var allCageEnableKeys = []string{
	"bullyenable", "mmoabenable", "ufoenable", "horrorenable",
	"superenable", "motherenable", "lolenable", "clownenable",
	"flowerenable", "crawlerenable", "destroyerenable",
}

// cageObjectNames maps enable-key → object name so we can cross-reference.
var cageObjectNames = []string{
	"Bully_Cage", "MMoab_Cage", "UFo_Cage", "Horror_Cage",
	"Super_Bloon_Cage", "Mother_Cage", "LUL_Cage_", "Clown_Cage",
	"Blooming_Cage", "Crawler_Cage", "Destroyer_Cage",
}

// BountyCageBehavior is the clickable boss icon in the Bounty Center.
// Left-click toggles selection (deselects other cages).
// When already selected, left-click starts the boss fight.
type BountyCageBehavior struct {
	engine.DefaultBehavior
	EnableKey      string    // e.g. "bullyenable"
	ProgressKey    string    // e.g. "b1" — highest unlocked difficulty (0-3)
	GoRooms        [4]string // rooms for difficulties 1-4
	BossName       string    // display name for the tooltip
	CageSprite     string    // e.g. "Big_Bully_Bloon_Cage_spr"
	InfoPanelOffset int      // offset into Bounty_Info_Panel_spr (0,4,8,...)
}

func (b *BountyCageBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["bossname"] = b.BossName
}

// MouseLeftPressed — toggle selection or launch fight.
// First click: select this cage (deselect others, set to saved progress+1).
// Second click (already selected): go to the fight room at current difficulty.
func (b *BountyCageBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	cur := getGlobal(g, b.EnableKey)
	if cur == 0 {
		// not selected → select at saved progress level
		newVal := getGlobal(g, b.ProgressKey) + 1
		if newVal > 4 {
			newVal = 4
		}
		if newVal < 1 {
			newVal = 1
		}
		// deselect all other cages
		for _, ek := range allCageEnableKeys {
			g.GlobalVars[ek] = 0.0
		}
		g.GlobalVars[b.EnableKey] = newVal
	} else {
		// already selected → show confirmation popup
		diff := int(cur)
		if diff < 1 || diff > 4 {
			return
		}
		roomName := b.GoRooms[diff-1]
		if roomName == "" {
			fmt.Printf("WARNING: no room for %s difficulty %d\n", b.EnableKey, diff)
			return
		}
		// store pending fight info and spawn confirmation dialog
		g.GlobalVars["bounty_pending_room"] = roomName
		g.GlobalVars["bounty_pending_boss"] = b.BossName
		g.GlobalVars["bounty_pending_diff"] = float64(diff)
		if g.InstanceMgr.InstanceCount("Bounty_Confirm") == 0 {
			g.InstanceMgr.Create("Bounty_Confirm", 512, 288)
		}
	}
}

// Draw — replicate the original GMX Draw event:
//  1. draw cage sprite at instance position, subimage = enable value
//  2. draw Bounty_Info_Panel_spr at absolute (32, 384), subimage = enable + offset
//  3. draw Stars at instance position, subimage = progress value
func (b *BountyCageBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	enable := int(getGlobal(g, b.EnableKey))
	progress := int(getGlobal(g, b.ProgressKey))

	// 1. cage icon — subimage selects the visual state (0=unselected, 1-4=selected difficulty)
	if cageSpr := g.AssetManager.GetSprite(b.CageSprite); cageSpr != nil && len(cageSpr.Frames) > 0 {
		frame := enable % len(cageSpr.Frames)
		engine.DrawSpriteExt(screen, cageSpr.Frames[frame], cageSpr.XOrigin, cageSpr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
	}

	// 2. info panel at bottom — only when this cage is selected
	if enable != 0 {
		if panel := g.AssetManager.GetSprite("Bounty_Info_Panel_spr"); panel != nil && len(panel.Frames) > 0 {
			panelFrame := (enable + b.InfoPanelOffset) % len(panel.Frames)
			engine.DrawSpriteExt(screen, panel.Frames[panelFrame], panel.XOrigin, panel.YOrigin,
				32, 384, 1, 1, 0, 1)
		}
	}

	// 3. star progress overlay — subimage = highest completed difficulty (skip when 0)
	if progress > 0 {
		if starsSpr := g.AssetManager.GetSprite("Stars"); starsSpr != nil && len(starsSpr.Frames) > 0 {
			starFrame := progress % len(starsSpr.Frames)
			engine.DrawSpriteExt(screen, starsSpr.Frames[starFrame], starsSpr.XOrigin, starsSpr.YOrigin,
				inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
		}
	}
}

// ---------------------------------------------------------------------------
// Bounty_Down — scroll button (shifts all cages up by 320px)
// ---------------------------------------------------------------------------
type BountyDownBehavior struct {
	engine.DefaultBehavior
}

func (b *BountyDownBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["trackpanel"] = 0.0
}

func (b *BountyDownBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	tp := getGlobal(g, "trackpanel")
	if tp < 2 {
		g.GlobalVars["trackpanel"] = tp + 1
		// shift all cage instances up by 320
		for _, name := range cageObjectNames {
			for _, cage := range g.InstanceMgr.FindByObject(name) {
				cage.Y -= 320
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Bounty_Up — scroll button (shifts all cages down by 320px)
// Currently not in the room but may be created dynamically.
// ---------------------------------------------------------------------------
type BountyUpBehavior struct {
	engine.DefaultBehavior
}

func (b *BountyUpBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	tp := getGlobal(g, "trackpanel")
	if tp > 0 {
		g.GlobalVars["trackpanel"] = tp - 1
		for _, name := range cageObjectNames {
			for _, cage := range g.InstanceMgr.FindByObject(name) {
				cage.Y += 320
			}
		}
	}
}

// ---------------------------------------------------------------------------
// Bloon_Bounty_Prompt — info/tooltip panel.
// In the original GMX this uses Bounty_Prompt_spr with default drawing (animated).
// No custom Draw event — the engine's default draw handles it.
// ---------------------------------------------------------------------------
type BloonBountyPromptBehavior struct {
	engine.DefaultBehavior
}

func (b *BloonBountyPromptBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0
	inst.Vars["prompt"] = 1.0
}

// ---------------------------------------------------------------------------
// Start_Over_butt — resets all bounty selections and scroll position
// ---------------------------------------------------------------------------
type StartOverButtBehavior struct {
	engine.DefaultBehavior
}

func (b *StartOverButtBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	// deselect all cages
	for _, ek := range allCageEnableKeys {
		g.GlobalVars[ek] = 0.0
	}
	// reset scroll position
	g.GlobalVars["trackpanel"] = 0.0
}

// ---------------------------------------------------------------------------
// RetryBountyBehavior — cycles difficulty of the currently selected cage.
// In the original GM, this is the info panel that shows the boss details.
// Left-click increments difficulty (wraps around).
// ---------------------------------------------------------------------------
type RetryBountyBehavior struct {
	engine.DefaultBehavior
}

// Create — Retry is invisible in the original GMX (visible=0); it acts as a
// click-target overlapping the info panel drawn by the selected cage.
func (b *RetryBountyBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Visible = false
}

func (b *RetryBountyBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	// for each cage that's currently selected, cycle its difficulty
	for idx, ek := range allCageEnableKeys {
		cur := getGlobal(g, ek)
		if cur == 0 {
			continue
		}
		pk := bountyCageConfigs[idx].ProgressKey
		maxUnlocked := getGlobal(g, pk) + 1
		if maxUnlocked > 4 {
			maxUnlocked = 4
		}
		next := cur + 1
		if next > maxUnlocked || next > 4 {
			next = 1 // wrap around
		}
		g.GlobalVars[ek] = next
	}
}

// ---------------------------------------------------------------------------
// Cage configurations table
// ---------------------------------------------------------------------------

type cageConfig struct {
	ObjectName      string
	EnableKey       string
	ProgressKey     string
	GoRooms         [4]string
	BossName        string
	CageSprite      string // sprite drawn in the Draw event
	InfoPanelOffset int    // offset for Bounty_Info_Panel_spr subimage
}

// bountyCageConfigs defines all 11 boss cage configurations.
// GoRooms are the challenge rooms containing the appropriate _Go controller.
// Rooms are mapped from GMX data — empty strings mean that difficulty is not yet available.
var bountyCageConfigs = []cageConfig{
	{
		ObjectName:      "Bully_Cage",
		EnableKey:       "bullyenable",
		ProgressKey:     "b1",
		GoRooms:         [4]string{"Village_Challenge", "Apprentice_Challenge", "Spactory_Challenge", ""},
		BossName:        "Big Bully Bloon",
		CageSprite:      "Big_Bully_Bloon_Cage_spr",
		InfoPanelOffset: 0,
	},
	{
		ObjectName:      "MMoab_Cage",
		EnableKey:       "mmoabenable",
		ProgressKey:     "b2",
		GoRooms:         [4]string{"Village_Challenge", "Apprentice_Challenge", "", ""},
		BossName:        "Mighty MOAB",
		CageSprite:      "Mighty_Moab_Cage_spr",
		InfoPanelOffset: 4,
	},
	{
		ObjectName:      "UFo_Cage",
		EnableKey:       "ufoenable",
		ProgressKey:     "b4",
		GoRooms:         [4]string{"Village_Challenge", "Apprentice_Challenge", "", ""},
		BossName:        "UFO Bloon",
		CageSprite:      "UFO_Bloon_Cage_spr",
		InfoPanelOffset: 12,
	},
	{
		ObjectName:      "Horror_Cage",
		EnableKey:       "horrorenable",
		ProgressKey:     "b3",
		GoRooms:         [4]string{"Village_Challenge", "Apprentice_Challenge", "", ""},
		BossName:        "Horror Bloon",
		CageSprite:      "Horror_Bloon_Cage_spr",
		InfoPanelOffset: 8,
	},
	{
		ObjectName:      "Super_Bloon_Cage",
		EnableKey:       "superenable",
		ProgressKey:     "b5",
		GoRooms:         [4]string{"Prison_Break_I", "Prison_Break_II", "Prison_Break_III", "Prison_Break_IV"},
		BossName:        "Super Bloon",
		CageSprite:      "Super_Bloon_Cage_spr",
		InfoPanelOffset: 16,
	},
	{
		ObjectName:      "Mother_Cage",
		EnableKey:       "motherenable",
		ProgressKey:     "b6",
		GoRooms:         [4]string{"Mortar_Challenge", "Ninja_Challenge", "", ""},
		BossName:        "The Mother",
		CageSprite:      "The_Mother_Cage_spr",
		InfoPanelOffset: 20,
	},
	{
		ObjectName:      "LUL_Cage_",
		EnableKey:       "lolenable",
		ProgressKey:     "b7",
		GoRooms:         [4]string{"Village_Challenge", "Charge_Challenge", "", "Super_Challenge"},
		BossName:        "LOL Bloon",
		CageSprite:      "Placeholder_cage_spr",
		InfoPanelOffset: 24,
	},
	{
		ObjectName:      "Clown_Cage",
		EnableKey:       "clownenable",
		ProgressKey:     "b8",
		GoRooms:         [4]string{"Sub_Challenge", "Bucc_Challenge", "", ""},
		BossName:        "Clown Bloon",
		CageSprite:      "Clown_Bloon_Cage_spr",
		InfoPanelOffset: 28,
	},
	{
		ObjectName:      "Blooming_Cage",
		EnableKey:       "flowerenable",
		ProgressKey:     "b9",
		GoRooms:         [4]string{"Glue_Challenge", "Ice_Challenge", "", ""},
		BossName:        "Blooming Bloon",
		CageSprite:      "Blooming_Cage_spr",
		InfoPanelOffset: 32,
	},
	{
		ObjectName:      "Crawler_Cage",
		EnableKey:       "crawlerenable",
		ProgressKey:     "b10",
		GoRooms:         [4]string{"Chipper_Challenge", "Plasma_Challenge", "", ""},
		BossName:        "Track Crawler",
		CageSprite:      "Track_Crawler_Cage_spr",
		InfoPanelOffset: 36,
	},
	{
		ObjectName:      "Destroyer_Cage",
		EnableKey:       "destroyerenable",
		ProgressKey:     "b11",
		GoRooms:         [4]string{"Village_Challenge", "Apprentice_Challenge", "", ""},
		BossName:        "The Destroyer",
		CageSprite:      "The_Destroyer_of_Monkeys_spr",
		InfoPanelOffset: 40,
	},
}

// ---------------------------------------------------------------------------
// Bounty_Confirm — "Enter Boss Fight?" popup before entering a boss fight.
// Uses Go_to_Main sprite as the dialog frame (has Yes/No buttons), but
// overrides Draw to paint "Enter Boss Fight?" over the baked-in text.
// Left half = Yes, right half / outside = No.
// ---------------------------------------------------------------------------
type BountyConfirmBehavior struct {
	engine.DefaultBehavior
}

func (b *BountyConfirmBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Go_to_Main"
	inst.Depth = -100
	inst.ImageSpeed = 0
}

func (b *BountyConfirmBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	// draw the base dialog sprite (has Yes/No buttons)
	engine.DrawSpriteExt(screen, spr.Frames[0], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)

	// paint over the baked-in "Go to Main Menu?" text area with the brown bg
	// The text sits roughly at y_offset -45 to -15 from center (sprite origin=120)
	bgColor := color.RGBA{126, 59, 0, 255}
	coverX := float32(inst.X - 170)
	coverY := float32(inst.Y - 45)
	vector.DrawFilledRect(screen, coverX, coverY, 340, 32, bgColor, false)

	// draw our replacement text centered
	label := "Enter Boss Fight?"
	textW := len(label) * 7 // basicfont is ~7px per char
	tx := int(inst.X) - textW/2
	ty := int(inst.Y) - 24
	etext.Draw(screen, label, basicfont.Face7x13, tx, ty, color.White)
}

func (b *BountyConfirmBehavior) MouseGlobalLeftPressed(inst *engine.Instance, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	halfW := 180.0
	halfH := 120.0
	if spr != nil && spr.Width > 0 {
		halfW = float64(spr.XOrigin) * inst.ImageXScale
		halfH = float64(spr.YOrigin) * inst.ImageYScale
		if halfW <= 0 {
			halfW = float64(spr.Width) * inst.ImageXScale / 2.0
		}
		if halfH <= 0 {
			halfH = float64(spr.Height) * inst.ImageYScale / 2.0
		}
	}

	roomX, roomY := g.GetMouseRoomPos()
	if roomX < inst.X-halfW || roomX > inst.X+halfW ||
		roomY < inst.Y-halfH || roomY > inst.Y+halfH {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	if roomX >= inst.X {
		// right half → No
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	// left half → Yes → launch the boss fight
	pendingRoom, ok := g.GlobalVars["bounty_pending_room"]
	if !ok {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	roomName, _ := pendingRoom.(string)
	if roomName == "" {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	diff := getGlobal(g, "bounty_pending_diff")
	bossName, _ := g.GlobalVars["bounty_pending_boss"].(string)

	g.GlobalVars["challenge"] = 0.0
	g.GlobalVars["normalmodeselect"] = 1.0
	g.GlobalVars["impoppablemodeselect"] = 0.0
	g.GlobalVars["nightmaremodeselect"] = 0.0
	fmt.Printf("Bounty: starting %s difficulty %.0f → room %s\n", bossName, diff, roomName)
	g.InstanceMgr.Destroy(inst.ID)
	g.RequestRoomGoto(roomName)
}

func (b *BountyConfirmBehavior) MouseRightPressed(inst *engine.Instance, g *engine.Game) {
	g.InstanceMgr.Destroy(inst.ID)
}

func (b *BountyConfirmBehavior) KeyPress(inst *engine.Instance, g *engine.Game) {
	if g.InputMgr.KeyPressed(ebiten.KeyEscape) || g.InputMgr.KeyPressed(ebiten.KeyX) {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// ---------------------------------------------------------------------------
// Registration
// ---------------------------------------------------------------------------

func RegisterBountyCenterBehaviors(im *engine.InstanceManager) {
	// register all 11 cage behaviors
	for _, cfg := range bountyCageConfigs {
		c := cfg // capture for closure
		im.RegisterBehavior(c.ObjectName, func() engine.InstanceBehavior {
			return &BountyCageBehavior{
				EnableKey:       c.EnableKey,
				ProgressKey:     c.ProgressKey,
				GoRooms:         c.GoRooms,
				BossName:        c.BossName,
				CageSprite:      c.CageSprite,
				InfoPanelOffset: c.InfoPanelOffset,
			}
		})
	}

	// scroll button
	im.RegisterBehavior("Bounty_Down", func() engine.InstanceBehavior { return &BountyDownBehavior{} })
	im.RegisterBehavior("Bounty_Up", func() engine.InstanceBehavior { return &BountyUpBehavior{} })

	// prompt tooltip
	im.RegisterBehavior("Bloon_Bounty_Prompt", func() engine.InstanceBehavior { return &BloonBountyPromptBehavior{} })

	// start over / retry
	im.RegisterBehavior("Start_Over_butt", func() engine.InstanceBehavior { return &StartOverButtBehavior{} })
	im.RegisterBehavior("Retry", func() engine.InstanceBehavior { return &RetryBountyBehavior{} })

	// confirmation dialog (uses Are_you_sure_tho sprite)
	im.RegisterBehavior("Bounty_Confirm", func() engine.InstanceBehavior { return &BountyConfirmBehavior{} })
}
