// package behaviors contains object behaviors for btdx.
package behaviors

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"btdx/internal/engine"
	"btdx/internal/savedata"
)

// main_Menu_Control — controls the Start_Screen title sequence.
// on any key/click: starts scrolling down, after 180 frames goes to Main_Menu.
// also spawns Menu_Bloon (alarm 0) and Menu_Blimp (alarm 5) decorations.
type MainMenuControl struct {
	engine.DefaultBehavior
}

func (b *MainMenuControl) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["start"] = 0
	inst.Vars["ticks"] = 0
	inst.Alarms[5] = 5
	inst.Alarms[0] = 1

	// Hide the static ANY_KEY tile so we can draw an animated pulse version!
	rm := g.RoomManager.GetCurrent()
	if rm != nil {
		for i := 0; i < len(rm.Tiles); i++ {
			if rm.Tiles[i].BGName == "ANY_KEY" {
				// Remove it from the renderer
				rm.Tiles = append(rm.Tiles[:i], rm.Tiles[i+1:]...)
				break // Only one ANY_KEY exists
			}
		}
	}
}

// Draw animates the pulsing ANY_KEY text (missing original animation request)
func (b *MainMenuControl) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	rm := g.RoomManager.GetCurrent()
	var viewX, viewY float64
	if rm != nil && len(rm.Views) > 0 {
		viewX = float64(rm.Views[0].XView)
		viewY = float64(rm.Views[0].YView)
	}

	targetRoomX := 256.0
	targetRoomY := 1248.0
	screenX := targetRoomX - viewX
	screenY := targetRoomY - viewY

	bg := g.AssetManager.GetBackground("ANY_KEY")
	if bg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(screenX, screenY)

		t, _ := inst.Vars["ticks"].(int)
		// Pulse the alpha smoothly back and forth!
		alpha := 0.65 + 0.35*math.Sin(float64(t)*0.08)
		op.ColorScale.ScaleAlpha(float32(alpha))

		screen.DrawImage(bg, op)
	}
}

func (b *MainMenuControl) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	switch idx {
	case 0:
		// spawn a random bloon at bottom of room floating up
		g.InstanceMgr.Create("Menu_Bloon", -100+rand.Float64()*1224, 1640)
		inst.Alarms[0] = 7 + int(rand.Float64()*15)
	case 1:
		// transition to Main_Menu
		g.RequestRoomGoto("Main_Menu")
		g.InstanceMgr.Create("Main_Menu_Music", inst.X, inst.Y)
	case 5:
		// spawn a blimp from the left
		g.InstanceMgr.Create("Menu_Blimp", -100, rand.Float64()*1400)
		inst.Alarms[5] = 120 + int(rand.Float64()*480)
	}
}

func (b *MainMenuControl) Step(inst *engine.Instance, g *engine.Game) {
	t, _ := inst.Vars["ticks"].(int)
	inst.Vars["ticks"] = t + 1

	if inst.Speed > 0 {
		inst.Speed += 0.04
		// recompute hspeed/vspeed from speed/direction
		inst.MotionSet(inst.Direction, inst.Speed)

		// To fix the "deadzone" jump where the view normally doesn't scroll
		// for 138 frames before sharply jumping down, we override the camera continuously!
		rm := g.RoomManager.GetCurrent()
		if rm != nil && len(rm.Views) > 0 {
			// Lock exactly 16 pixels above the target's visual coordinate
			// (Preserves exact starting parity where YView=976, Y=992)
			rm.Views[0].YView = int(inst.Y) - 16
		}
	}
}

func (b *MainMenuControl) startTransition(inst *engine.Instance, g *engine.Game) {
	start, _ := inst.Vars["start"].(int)
	if start == 0 {
		inst.Vars["start"] = 1
		// move upward (direction 90 = up); original GML uses action_move grid 000000010 = position 8 = UP
		inst.MotionSet(90, 0.1)
		inst.Alarms[1] = 180

		rm := g.RoomManager.GetCurrent()
		if rm != nil && len(rm.Views) > 0 {
			rm.Views[0].ObjName = "" // stop standard view bounding box logic that created the deadzone
		}
	}
}

func (b *MainMenuControl) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	b.startTransition(inst, g)
}

func (b *MainMenuControl) KeyPress(inst *engine.Instance, g *engine.Game) {
	b.startTransition(inst, g)
}

// menu_Bloon — decorative bloons floating up on the start screen.
// random color, speed, direction. Clickable (pops).
type MenuBloon struct {
	engine.DefaultBehavior
}

func (b *MenuBloon) Create(inst *engine.Instance, g *engine.Game) {
	inst.Speed = 0.5 + rand.Float64()*2
	inst.Direction = 72 + rand.Float64()*36
	inst.Depth += int(rand.Float64() * 4.5)
	inst.ImageSpeed = 0

	bloonSprites := []string{
		"Red_Bloon_Spr", "Blue_Bloon_Spr", "Green_Bloon_Spr",
		"Yellow_Bloon_Spr", "Pink_Bloon_Spr",
	}
	pick := rand.Intn(len(bloonSprites))
	inst.SpriteName = bloonSprites[pick]
	inst.ImageIndex = 0

	inst.MotionSet(inst.Direction, inst.Speed)
	inst.Alarms[0] = 1500
}

func (b *MenuBloon) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		inst.Destroyed = true
	}
}

func (b *MenuBloon) MouseGlobalLeftPressed(inst *engine.Instance, g *engine.Game) {
	// pop when clicked anywhere (check if mouse is over this bloon)
	if g.IsMouseOverInstance(inst) {
		inst.Destroyed = true
		g.InstanceMgr.Create("Harmless_Pop", inst.X, inst.Y)
	}
}

// menu_Blimp — decorative blimps floating across the start screen.
type MenuBlimp struct {
	engine.DefaultBehavior
}

func (b *MenuBlimp) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Mini_Moab_New_Spr"
	inst.ImageSpeed = 0.5
	inst.Speed = 0.2 + rand.Float64()*1.0
	inst.Direction = 0
	inst.Depth += int(rand.Float64() * 4.5)
	inst.MotionSet(inst.Direction, inst.Speed)
	inst.Alarms[0] = 5000
}

func (b *MenuBlimp) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		inst.Destroyed = true
	}
}

func (b *MenuBlimp) Step(inst *engine.Instance, g *engine.Game) {
	if inst.ImageIndex > 8 {
		inst.ImageIndex = 0
	}
}

func (b *MenuBlimp) MouseGlobalLeftPressed(inst *engine.Instance, g *engine.Game) {
	if g.IsMouseOverInstance(inst) {
		inst.Destroyed = true
		g.InstanceMgr.Create("Harmless_Pop", inst.X, inst.Y)
	}
}

// harmless_Pop — simple pop effect, destroys itself after animation
type HarmlessPop struct {
	engine.DefaultBehavior
}

func (b *HarmlessPop) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Pop_Spr"
	inst.ImageSpeed = 0.5
	inst.Alarms[0] = 15
}

func (b *HarmlessPop) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		inst.Destroyed = true
	}
}

// real_MM_Music — plays main menu music on Create
type RealMMMusic struct {
	engine.DefaultBehavior
}

func (b *RealMMMusic) Create(inst *engine.Instance, g *engine.Game) {
	g.AudioMgr.PlayMusic("Main_Menu0")
}

// main_Menu_Music — plays music when entering Main_Menu room
type MainMenuMusic struct {
	engine.DefaultBehavior
}

func (b *MainMenuMusic) Create(inst *engine.Instance, g *engine.Game) {
	mute, _ := g.GlobalVars["mute"].(float64)
	if mute == 0 {
		g.AudioMgr.PlayMusic("Main_Menu0")
	}
}

// career_Control — persistent save data manager
// initializes all career variables with defaults.
type CareerControl struct {
	engine.DefaultBehavior
}

func (b *CareerControl) Create(inst *engine.Instance, g *engine.Game) {
	inst.Persistent = true
	inst.Depth = -21

	g.GlobalVars["soundandmusic"] = 1.0
	g.GlobalVars["mute"] = 0.0
	g.GlobalVars["soundmute"] = 0.0
	g.GlobalVars["trackselect"] = 0.0
	g.GlobalVars["sandbox"] = 0.0

	// initialize career defaults
	g.GlobalVars["BP"] = 5.0
	g.GlobalVars["monkeymoney"] = 150.0
	g.GlobalVars["bsouls"] = 0.0
	g.GlobalVars["trophies"] = 0.0
	g.GlobalVars["XP"] = 0.0
	g.GlobalVars["rank"] = 0.0
	g.GlobalVars["criteria"] = 50.0
	g.GlobalVars["careershow"] = 0.0

	// agents
	g.GlobalVars["angrysquirrel"] = 0.0
	g.GlobalVars["bloonburybush"] = 0.0
	g.GlobalVars["sprinkler"] = 0.0
	g.GlobalVars["monkeynurse"] = 0.0
	g.GlobalVars["bananamobile"] = 0.0

	// track milestones, hardstones, best waves, scores, trials, etc.
	for i := 1; i <= 32; i++ {
		g.GlobalVars[fmt.Sprintf("track%dmilestone", i)] = 0.0
		g.GlobalVars[fmt.Sprintf("track%dhardstone", i)] = 0.0
		g.GlobalVars[fmt.Sprintf("track%dbestwave", i)] = 0.0
		g.GlobalVars[fmt.Sprintf("track%dbesthardwave", i)] = 0.0
		g.GlobalVars[fmt.Sprintf("track%dnightstone", i)] = 0.0
		g.GlobalVars[fmt.Sprintf("x%d", i)] = 0.0
		g.GlobalVars[fmt.Sprintf("xx%d", i)] = 0.0
		g.GlobalVars[fmt.Sprintf("t%d", i)] = 9999.0
		g.GlobalVars[fmt.Sprintf("n%d", i)] = 0.0
	}

	// special missions
	for i := 1; i <= 16; i++ {
		g.GlobalVars[fmt.Sprintf("specialmission%d", i)] = 0.0
	}

	// bounties
	for i := 1; i <= 12; i++ {
		g.GlobalVars[fmt.Sprintf("b%d", i)] = 0.0
	}

	// challenges
	for i := 1; i <= 6; i++ {
		g.GlobalVars[fmt.Sprintf("c%d", i)] = 0.0
		g.GlobalVars[fmt.Sprintf("c%d", i+100)] = 0.0
	}

	// tower upgrade paths (all start at 1 = locked, 0 = none)
	towerPaths := []string{
		"DML", "DMM", "DMR", "TSL", "TSM", "TSR", "BML", "BMM", "BMR",
		"SnML", "SnMM", "SnMR", "NML", "NMM", "NMR", "BCL", "BCM", "BCR",
		"MSL", "MSM", "MSR", "CTL", "CTM", "CTR",
		"GGL", "GGM", "GGR", "IML", "IMM", "IMR", "MBL", "MBM", "MBR",
		"MEL", "MEM", "MER", "MAL", "MAM", "MAR", "BChL", "BChM", "BChR",
		"MApL", "MApM", "MApR", "MAlL", "MAlM", "MAlR",
		"MVL", "MVM", "MVR", "BTL", "BTM", "BTR", "DGL", "DGM", "DGR",
		"MLL", "MLM", "MLR", "HPL", "HPM", "HPR", "SFL", "SFM", "SFR",
		"PML", "PMM", "PMR", "SuML", "SuMM", "SuMR",
	}
	for _, p := range towerPaths {
		g.GlobalVars[p] = 1.0
	}

	// bloon toggles
	for _, b := range []string{
		"bullyenable", "mmoabenable", "horrorenable", "ufoenable",
		"superenable", "motherenable", "lolenable", "clownenable",
		"flowerenable", "crawlerenable", "destroyerenable",
	} {
		g.GlobalVars[b] = 0.0
	}

	g.GlobalVars["autostart"] = 0.0
	g.GlobalVars["challenge"] = 0.0
	g.GlobalVars["normalmodeselect"] = 0.0
	g.GlobalVars["impoppablemodeselect"] = 0.0
	g.GlobalVars["nightmaremodeselect"] = 0.0
	g.GlobalVars["totalachievements"] = 0.0
	g.GlobalVars["PopUp"] = 0.0

	// load saved career data (overwrites defaults for any key present in the file)
	if err := savedata.Load(g); err != nil {
		fmt.Printf("WARNING: could not load save data: %v\n", err)
	}

	// move to origin
	inst.X = 0
	inst.Y = 0
}

// view follower helper — the Start_Screen view follows Main_Menu_Control
// we handle this in the engine by checking if the view's ObjName is set.

// registerStartScreenBehaviors registers all Start_Screen room behaviors
func RegisterStartScreenBehaviors(im *engine.InstanceManager) {
	im.RegisterBehavior("Main_Menu_Control", func() engine.InstanceBehavior { return &MainMenuControl{} })
	im.RegisterBehavior("Career_Control", func() engine.InstanceBehavior { return &CareerControl{} })
	im.RegisterBehavior("Real_MM_Music", func() engine.InstanceBehavior { return &RealMMMusic{} })
	im.RegisterBehavior("Menu_Bloon", func() engine.InstanceBehavior { return &MenuBloon{} })
	im.RegisterBehavior("Menu_Blimp", func() engine.InstanceBehavior { return &MenuBlimp{} })
	im.RegisterBehavior("Harmless_Pop", func() engine.InstanceBehavior { return &HarmlessPop{} })
	im.RegisterBehavior("Main_Menu_Music", func() engine.InstanceBehavior { return &MainMenuMusic{} })
}
