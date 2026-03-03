package behaviors

import (
	"fmt"
	"image/color"
	"math"
	"strings"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
	etext "github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

// normal_Bloon_Branch — the core bloon object
// follows a path based on global.track, sprite changes by bloonlayer

// bloon layer → sprite name mapping
var bloonSprites = map[float64]string{
	1:   "Red_Bloon_Spr",
	1.5: "Orange_Bloon_Spr",
	2:   "Blue_Bloon_Spr",
	2.5: "Cyan_Bloon_Spr",
	3:   "Green_Bloon_Spr",
	3.5: "Lime_Bloon_Spr",
	4:   "Yellow_Bloon_Spr",
	4.5: "Amber_Bloon_Spr",
	5:   "Pink_Bloon_Spr",
	5.5: "Purple_Bloon_Spr",
	6:   "Black_Bloon_Spr",
	6.1: "White_Bloon_Spr",
	7:   "Zebra_Bloon_Spr",
	8:   "Rainbow_Bloon_Spr",
	8.5: "Prismatic_Bloon_Spr",
}

// bloon layer → speed multiplier
var bloonSpeeds = map[float64]float64{
	1: 1.6, 1.5: 2.2, 2: 2.0, 2.5: 2.6, 3: 2.4, 3.5: 3.0,
	4: 4.4, 4.5: 5.0, 5: 4.8, 5.5: 5.4, 6: 2.4, 6.1: 2.8,
	7: 2.8, 8: 3.8, 8.5: 4.4,
}

// bloon layer → image scale
var bloonScales = map[float64]float64{
	1: 0.9, 1.5: 1.0, 2: 0.95, 2.5: 1.05, 3: 1.0, 3.5: 1.1,
	4: 1.05, 4.5: 1.15, 5: 1.1, 5.5: 1.2, 6: 0.85, 6.1: 0.85,
	7: 1.15, 8: 1.2, 8.5: 1.3,
}

// track number → path name(s) mapping
// for single-path tracks, only index 0 is used
// for multi-path tracks, paths cycle based on BloonSpawn.path counter
var trackPaths = map[int][]string{
	1:  {"Monkey_Meadows"},
	2:  {"Bloon_Oasis_Path"},
	3:  {"Spiral_Swamp_Path_A", "Spiral_Swamp_Path_B"},
	4:  {"Monkey_Fort_Path"},
	5:  {"Monkey_Town_Docks_Path_A", "Monkey_Town_Docks_Path_B"},
	6:  {"Conveyor_Belt_Path"},
	7:  {"The_Depths_Path_A", "The_Depths_Path_B"},
	8:  {"Sun_Dial_Path_A", "Sun_Dial_Path_B", "Sun_Dial_Path_C", "Sun_Dial_Path_D"},
	9:  {"Shade_Woods_Path_A", "Shade_Woods_Path_B"},
	10: {"Minecarts_Path_A", "Minecarts_Path_B"},
	11: {"Crimson_Creek_Path_A", "Crimson_Creek_Path_B"},
	12: {"Xtreme_Park_Path_A", "Xtreme_Park_Path_B", "Xtreme_Park_Path_C", "Xtreme_Park_Path_D"},
	13: {"Portal_Lab_Path"},
	14: {"Omega_River_Path"},
	15: {"Space_Portals_Path_A", "Space_Portals_Path_B", "Space_Portals_Path_C"},
	17: {"Bloonlight_Throwback_Path"},
	18: {"Bloon_Circles_X_Path_A", "Bloon_Circles_X_Path_B"},
	19: {"Autumn_Acres_Path_A", "Autumn_Acres_Path_B"},
	20: {"Graveyard_Path_A", "Graveyard_Path_B", "Graveyard_Path_C"},
	21: {"Village_A", "Village_B", "Village_C", "Village_D", "Village_E", "Village_F", "Village_G",
		"Village_H", "Village_I", "Village_J", "Village_K", "Village_L", "Village_M", "Village_N", "Village_O"},
	22: {"Circuit_Path_A", "Circuit_Path_B", "Circuit_Path_C", "Circuit_Path_D"},
	23: {"Grand_Canyon_Path_A", "Grand_Canyon_Path_B"},
	24: {"Bloonside_River_Path_A", "Bloonside_River_Path_B"},
	25: {"Cotton_Fields_Path_A", "Cotton_Fields_Path_B"},
	27: {"Rubber_Rug_Path_A", "Rubber_Rug_Path_B"},
	28: {"Frozen_Lake_Path_A", "Frozen_Lake_Path_B", "Frozen_Lake_Path_C", "Frozen_Lake_Path_D"},
	29: {"Sky_Battles_Path_A", "Sky_Battles_Path_B"},
	30: {"Lava_Stream_Path_A", "Lava_Stream_Path_B", "Lava_Stream_Path_C"},
	31: {"Ravine_River_Path_A", "Ravine_River_Path_B"},
	32: {"Peaceful_Bridge_Path"},
}

type NormalBloonBranch struct {
	engine.DefaultBehavior
}

func (b *NormalBloonBranch) Create(inst *engine.Instance, g *engine.Game) {
	// initialize bloon vars
	inst.Vars["glue"] = 0.0
	inst.Vars["freeze"] = 0.0
	inst.Vars["shielded"] = 0.0
	inst.Vars["bigbloon"] = 0.0
	inst.Vars["electric"] = 0.0
	inst.Vars["normal"] = 1.0
	inst.Vars["camo"] = 0.0
	inst.Vars["lead"] = 0.0
	inst.Vars["regrow"] = 0.0
	inst.Vars["tattered"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["distraction"] = 0.0
	inst.Vars["corrosion"] = 0.0
	inst.Vars["radiation"] = 0.0
	inst.Vars["fast"] = 1.6

	// default bloonlayer (set by timeline spawn code)
	if _, ok := inst.Vars["bloonlayer"]; !ok {
		inst.Vars["bloonlayer"] = 1.0
	}
	if _, ok := inst.Vars["bloonmaxlayer"]; !ok {
		inst.Vars["bloonmaxlayer"] = 1.0
	}

	// assign path based on global.track
	inst.Vars["path_progress"] = 0.0
	inst.Vars["path_name"] = ""
	assignBloonPath(inst, g)

	// set initial sprite based on layer
	updateBloonSprite(inst)
	updateBloonDepth(inst)
}

func (b *NormalBloonBranch) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0

	layer := getVar(inst, "bloonlayer")

	// clamp bloonlayer
	if layer < 5.5 && layer > 5 {
		inst.Vars["bloonlayer"] = 5.0
		layer = 5
	}
	if layer > 9 {
		inst.Vars["bloonlayer"] = 8.0
		layer = 8
	}

	// update sprite based on layer
	updateBloonSprite(inst)
	updateBloonDepth(inst)

	// calculate speed
	bspeed := getGlobal(g, "bspeed")
	if bspeed == 0 {
		bspeed = 1
	}

	speedMul := 1.6 // default red
	if s, ok := bloonSpeeds[layer]; ok {
		speedMul = s
	}

	tattered := getVar(inst, "tattered")
	fast := bspeed * speedMul * (tattered + 1)

	// glue effect
	glue := getVar(inst, "glue")
	if glue > 0 {
		fast = fast * (0.6 - glue*0.1)
	}

	// freeze effect
	freeze := getVar(inst, "freeze")
	if freeze >= 1 {
		fast = 0
	} else if freeze > 0 {
		fast = fast / (1 + freeze*2)
	}

	// stun
	if getVar(inst, "stun") == 1 {
		fast = 0
	}

	// distraction (moves backward)
	distraction := getVar(inst, "distraction")
	if distraction > 0 {
		fast = bspeed * -3 * distraction
	}

	inst.Vars["fast"] = fast

	// advance along path
	pathName, _ := inst.Vars["path_name"].(string)
	if pathName != "" && g.PathMgr != nil {
		pa := g.PathMgr.Get(pathName)
		if pa != nil && pa.TotalLength > 0 {
			progress := getVar(inst, "path_progress")
			// convert speed to progress increment
			// speed is in pixels/frame, total length is in pixels
			progressInc := fast / pa.TotalLength
			progress += progressInc
			inst.Vars["path_progress"] = progress

			if progress >= 1.0 {
				// reached end — will be caught by End's collision check
				inst.Vars["path_progress"] = 1.0
			}

			// update position from path
			x, y := g.PathMgr.GetPositionAtProgress(pathName, progress)
			inst.X = x
			inst.Y = y

			// update direction for rendering
			dir := g.PathMgr.GetDirectionAtProgress(pathName, progress)
			_ = dir // could use for image_angle if needed
		}
	}

	// tattered visual: use frame 1 if tattered
	if tattered == 1 {
		inst.ImageIndex = 1
	} else {
		inst.ImageIndex = 0
	}
}

// updateBloonDepth enforces deterministic overlap ordering:
// stronger bloons render above weaker ones when they intersect.
func updateBloonDepth(inst *engine.Instance) {
	layer := getVar(inst, "bloonlayer")
	progress := getVar(inst, "path_progress")

	// keep bloons behind towers/UI while ordering among themselves.
	// lower depth draws later (in front) in this engine.
	strengthPart := int(math.Round(layer * 100)) // stronger layer => larger value
	progressPart := int(math.Round(progress * 10))
	tiePart := inst.ID % 10 // stable tie-breaker for equal layer/progress
	inst.Depth = 5000 - strengthPart - progressPart - tiePart
}

// alarm handles freeze/glue timers
func (b *NormalBloonBranch) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 6 {
		// freeze expired
		inst.Vars["freeze"] = 0.0
		return
	}
	if idx == 8 {
		// flash bomb stun expired
		inst.Vars["stun"] = 0.0
		return
	}
	if idx == 9 {
		// distraction shot pullback expired
		inst.Vars["distraction"] = 0.0
	}
}

func (b *NormalBloonBranch) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}

	layer := getVar(inst, "bloonlayer")
	scale := 1.0
	if s, ok := bloonScales[layer]; ok {
		scale = s
	}

	frameIdx := int(inst.ImageIndex) % len(spr.Frames)
	if frameIdx < 0 {
		frameIdx = 0
	}

	engine.DrawSpriteExt(screen, spr.Frames[frameIdx], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, scale*inst.ImageXScale, scale*inst.ImageYScale,
		inst.ImageAngle, inst.ImageAlpha)

	// show bloon info tooltip when the setting is on and mouse hovers the bloon
	if getGlobal(g, "blooninfo") == 1 && g.IsMouseOverInstance(inst) {
		name := bloonLayerName(layer)
		info := fmt.Sprintf("%s (%.1f)", name, layer)
		drawBloonInfoTag(screen, g, info, inst.X, inst.Y-20)
	}
}

// bloonLayerName returns a readable name for the bloon layer
func bloonLayerName(layer float64) string {
	switch {
	case layer <= 1:
		return "Red"
	case layer <= 1.5:
		return "Orange"
	case layer <= 2:
		return "Blue"
	case layer <= 2.5:
		return "Cyan"
	case layer <= 3:
		return "Green"
	case layer <= 3.5:
		return "Lime"
	case layer <= 4:
		return "Yellow"
	case layer <= 4.5:
		return "Amber"
	case layer <= 5:
		return "Pink"
	case layer <= 5.5:
		return "Purple"
	case layer <= 6:
		return "Black"
	case layer <= 6.1:
		return "White"
	case layer <= 7:
		return "Zebra"
	case layer <= 8:
		return "Rainbow"
	case layer <= 8.5:
		return "Prismatic"
	default:
		return fmt.Sprintf("Layer %.1f", layer)
	}
}

// drawBloonInfoTag draws a small tooltip above a bloon
func drawBloonInfoTag(screen *ebiten.Image, g *engine.Game, text string, x, y float64) {
	tw := float64(len(text)*7 + 6)
	th := 16.0
	bx := x - tw/2
	by := y - th

	vector.DrawFilledRect(screen, float32(bx), float32(by), float32(tw), float32(th),
		color.RGBA{0, 0, 0, 180}, false)
	etext.Draw(screen, text, basicfont.Face7x13, int(bx)+3, int(by)+12,
		color.RGBA{255, 255, 255, 255})
}

// assignBloonPath assigns a path to the bloon based on global.track
func assignBloonPath(inst *engine.Instance, g *engine.Game) {
	track := int(getGlobal(g, "track"))
	if track == 0 {
		track = inferTrackFromRoomName(g.CurrentRoom)
		if track > 0 {
			g.GlobalVars["track"] = float64(track)
		}
	}
	paths, ok := trackPaths[track]
	if !ok || len(paths) == 0 {
		fmt.Printf("WARNING: No path for track %d\n", track)
		return
	}

	// get path index from BloonSpawn's cycling counter
	pathIdx := 0
	spawns := g.InstanceMgr.FindByObject("BloonSpawn")
	if len(spawns) > 0 {
		spawn := spawns[0]
		p := getVar(spawn, "path")
		pathIdx = int(p) % len(paths)

		// cycle to next path for the next bloon
		nextPath := p + 1
		if int(nextPath) >= len(paths) {
			nextPath = 0
		}
		spawn.Vars["path"] = nextPath
	}

	pathName := paths[pathIdx]

	// verify path exists
	if g.PathMgr.Get(pathName) == nil {
		fmt.Printf("WARNING: Path '%s' not found for track %d\n", pathName, track)
		// try without the last character variations
		return
	}

	inst.Vars["path_name"] = pathName
	inst.Vars["path_progress"] = 0.0

	// set initial position from path start
	x, y := g.PathMgr.GetPositionAtProgress(pathName, 0)
	inst.X = x
	inst.Y = y
}

func inferTrackFromRoomName(room string) int {
	switch {
	case strings.Contains(room, "Monkey_Meadows"):
		return 1
	case strings.Contains(room, "Bloon_Oasis"):
		return 2
	case strings.Contains(room, "Swamp"):
		return 3
	case strings.Contains(room, "Monkey_Fort"):
		return 4
	case strings.Contains(room, "Docks"):
		return 5
	case strings.Contains(room, "Conveyor"):
		return 6
	case strings.Contains(room, "Depths"):
		return 7
	case strings.Contains(room, "Sun"):
		return 8
	case strings.Contains(room, "Shade_Woods"):
		return 9
	case strings.Contains(room, "Minecarts"):
		return 10
	case strings.Contains(room, "Crimson_Creek"):
		return 11
	case strings.Contains(room, "Xtreme_Park"):
		return 12
	case strings.Contains(room, "Portal_Lab"):
		return 13
	case strings.Contains(room, "Omega_River"):
		return 14
	case strings.Contains(room, "Space_Portals"):
		return 15
	case strings.Contains(room, "Throwback"):
		return 17
	case strings.Contains(room, "Bloon_Circles_X"):
		return 18
	case strings.Contains(room, "Autumn_Acres"):
		return 19
	case strings.Contains(room, "Graveyard"):
		return 20
	case strings.Contains(room, "Village"):
		return 21
	case strings.Contains(room, "Circuit"):
		return 22
	case strings.Contains(room, "Grand_Canyon"):
		return 23
	case strings.Contains(room, "Bloonside"):
		return 24
	case strings.Contains(room, "Cotton_Fields"):
		return 25
	case strings.Contains(room, "Rubber_Rug"):
		return 27
	case strings.Contains(room, "Frozen_Lake"):
		return 28
	case strings.Contains(room, "Sky_Battles"):
		return 29
	case strings.Contains(room, "Lava_Stream"):
		return 30
	case strings.Contains(room, "Ravine_River"):
		return 31
	case strings.Contains(room, "Peaceful_Bridge"):
		return 32
	default:
		return 0
	}
}

// updateBloonSprite sets the sprite based on bloonlayer
func updateBloonSprite(inst *engine.Instance) {
	layer := getVar(inst, "bloonlayer")

	// find exact match first
	if spriteName, ok := bloonSprites[layer]; ok {
		inst.SpriteName = spriteName
		return
	}

	// find closest match
	bestLayer := 1.0
	bestDist := math.Abs(layer - 1.0)
	for l := range bloonSprites {
		d := math.Abs(layer - l)
		if d < bestDist {
			bestDist = d
			bestLayer = l
		}
	}
	inst.SpriteName = bloonSprites[bestLayer]
}

// getVar gets a float64 from instance vars
func getVar(inst *engine.Instance, key string) float64 {
	v, _ := inst.Vars[key].(float64)
	return v
}

// popBloon handles bloon popping with proper splitting logic
//
// layer categories:
//
//	Integer layers (1-5): just reduce layer, no children
//	Half layers (.5): self + 2 children = 3 bloons
//	Layer 6.1 (White): self + 1 child = 2 bloons
//	Layers 6/7/8 (Black/Zebra/Rainbow): self + 1 child = 2 bloons
//	  bonus children at higher LP values
func popBloon(bloon *engine.Instance, lp float64, g *engine.Game) {
	layer := getVar(bloon, "bloonlayer")
	cashReward := getGlobal(g, "cashwavereward")

	// reset tattered on hit
	bloon.Vars["tattered"] = 0.0

	frac := layer - math.Floor(layer)

	switch {
	// category 1: Normal (integer layers 1-5)
	case frac < 0.05 && layer <= 5:
		newLayer := layer - lp
		if newLayer < 1 {
			g.GlobalVars["cashwavereward"] = cashReward + layer
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + layer
			g.InstanceMgr.Destroy(bloon.ID)
		} else {
			g.GlobalVars["cashwavereward"] = cashReward + lp
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + lp
			bloon.Vars["bloonlayer"] = newLayer
		}

	// category 2: Multi bloons (half layers: .5)
	case math.Abs(frac-0.5) < 0.05:
		newLayer := layer - (lp - 0.5)
		if newLayer < 1 {
			g.GlobalVars["cashwavereward"] = cashReward + layer
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + layer
			g.InstanceMgr.Destroy(bloon.ID)
		} else {
			newLayer = math.Floor(newLayer)
			g.GlobalVars["cashwavereward"] = cashReward + lp
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + lp
			bloon.Vars["bloonlayer"] = newLayer
			// spawn 2 children at the same new layer
			spawnBloonChild(bloon, newLayer, g)
			spawnBloonChild(bloon, newLayer, g)
		}

	// category 3: White (6.1)
	case math.Abs(frac-0.1) < 0.05:
		newLayer := math.Round(layer - lp)
		if newLayer < 1 {
			g.GlobalVars["cashwavereward"] = cashReward + layer
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + layer
			g.InstanceMgr.Destroy(bloon.ID)
		} else {
			g.GlobalVars["cashwavereward"] = cashReward + lp
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + lp
			bloon.Vars["bloonlayer"] = newLayer
			// 1 child
			spawnBloonChild(bloon, newLayer, g)
		}

	// category 4: Black(6)/Zebra(7)/Rainbow(8)
	default:
		newLayer := layer - lp
		if newLayer < 1 {
			g.GlobalVars["cashwavereward"] = cashReward + layer
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + layer
			g.InstanceMgr.Destroy(bloon.ID)
		} else {
			g.GlobalVars["cashwavereward"] = cashReward + lp
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + lp
			bloon.Vars["bloonlayer"] = newLayer
			// base: 1 child
			spawnBloonChild(bloon, newLayer, g)
			// bonus children at higher LP
			if lp > 2 {
				spawnBloonChild(bloon, newLayer, g)
				spawnBloonChild(bloon, newLayer, g)
			}
			if lp > 3 {
				for i := 0; i < 4; i++ {
					spawnBloonChild(bloon, newLayer, g)
				}
			}
		}
	}

	// play pop sound
	playBloonPopSound(g)
}

func playBloonPopSound(g *engine.Game) {
	candidates := []string{"Popp", "Pop", "Popping"}
	for _, s := range candidates {
		if g.AudioMgr.HasSound(s) {
			g.AudioMgr.Play(s)
			return
		}
	}
}

// spawnBloonChild creates a child bloon at the parent's position/path
func spawnBloonChild(parent *engine.Instance, layer float64, g *engine.Game) {
	child := g.InstanceMgr.Create("Normal_Bloon_Branch", parent.X, parent.Y)
	if child == nil {
		return
	}
	child.Vars["bloonlayer"] = layer
	child.Vars["bloonmaxlayer"] = getVar(parent, "bloonmaxlayer")
	// inherit path state so child continues from same position
	// (overrides Create's assignBloonPath which set progress=0)
	child.Vars["path_name"] = parent.Vars["path_name"]
	child.Vars["path_progress"] = parent.Vars["path_progress"]
	// fix position to parent's current location (Create moved it to path start)
	child.X = parent.X
	child.Y = parent.Y
	// inherit special flags
	child.Vars["camo"] = getVar(parent, "camo")
	child.Vars["regrow"] = getVar(parent, "regrow")
	child.Vars["lead"] = getVar(parent, "lead")
	child.Vars["shielded"] = getVar(parent, "shielded")
}

// registerBloonBehaviors registers bloon-related behaviors
func RegisterBloonBehaviors(im *engine.InstanceManager) {
	im.RegisterBehavior("Normal_Bloon_Branch", func() engine.InstanceBehavior { return &NormalBloonBranch{} })
}
