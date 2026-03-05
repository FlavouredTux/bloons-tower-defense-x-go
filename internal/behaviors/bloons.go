package behaviors

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
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
	// Ceramic & Brick
	18: "Ceramic_Bloon_Spr",
	48: "Brick_Bloon_Spr",
	// MOAB class
	93:   "Mini_Moab_New_Spr",
	348:  "New_Moab_Spr",
	1248: "New_BFB_Spr",
	5248: "New_ZOMG_Spr",
	68.5: "New_HTA_Spr",
	593:  "New_BRC_Spr",
	351:  "New_DDT_Spr",
	318:  "LPZ_Spr",
	// Nightmare MOAB class
	248:     "Rocket_Blimp_Spr",
	918:     "Storm_LPZ_Spr",
	2593:    "Mega_BRC_Spr",
	3351:    "Deadly_DDT_Spr",
	10068.5: "Prismatic_HTA_Spr",
	17248:   "Party_Blimp_Spr",
}

// bloon layer → speed multiplier
var bloonSpeeds = map[float64]float64{
	1: 1.6, 1.5: 2.2, 2: 2.0, 2.5: 2.6, 3: 2.4, 3.5: 3.0,
	4: 4.4, 4.5: 5.0, 5: 4.8, 5.5: 5.4, 6: 2.4, 6.1: 2.8,
	7: 2.8, 8: 3.8, 8.5: 4.4,
	// Ceramic & Brick
	18: 3.0, 48: 2.6,
	// MOAB class (from Branch objects × bspeed)
	93: 2.1, 348: 1.5, 1248: 1.15, 5248: 0.95,
	68.5: 1.8, 593: 1.0, 351: 4.0, 318: 1.5,
	// Nightmare MOAB class
	248: 7.0, 918: 1.5, 2593: 0.6, 3351: 1.0,
	10068.5: 1.8, 17248: 0.6,
}

// bloon layer → image scale
var bloonScales = map[float64]float64{
	1: 0.9, 1.5: 1.0, 2: 0.95, 2.5: 1.05, 3: 1.0, 3.5: 1.1,
	4: 1.05, 4.5: 1.15, 5: 1.1, 5.5: 1.2, 6: 0.85, 6.1: 0.85,
	7: 1.15, 8: 1.2, 8.5: 1.3,
	// Ceramic & Brick
	18: 1.0, 48: 1.0,
	// MOAB class (from Branch objects)
	93: 1.3, 348: 1.3, 1248: 1.4, 5248: 1.5,
	68.5: 1.25, 593: 1.35, 351: 1.25, 318: 1.0,
	// Nightmare MOAB class (bigbloon-derived scales)
	248: 1.3, 918: 1.0, 2593: 1.35, 3351: 1.25,
	10068.5: 1.25, 17248: 1.4,
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

	// clamp ambiguous fractional layers (5.1-5.4 → 5.0)
	if layer > 5 && layer < 5.5 {
		inst.Vars["bloonlayer"] = 5.0
		layer = 5
	}

	// Auto-initialize MOAB shield on first step (for timeline-spawned MOABs
	// that don't get shield set by the spawning code)
	if getVar(inst, "shield_init") == 0 {
		inst.Vars["shield_init"] = 1.0
		if layer >= 68.5 && getVar(inst, "shielded") == 0 {
			moabBaseShields := map[float64]float64{
				68.5: 60, 93: 75, 318: 300, 348: 300,
				351: 303, 593: 500, 1248: 900, 5248: 4000,
				// Nightmare MOAB class
				248: 200, 918: 3500, 2593: 2000, 3351: 75,
				10068.5: 10000, 17248: 2000,
			}
			if hp, ok := moabBaseShields[layer]; ok {
				bpower := getGlobal(g, "bpower")
				if bpower == 0 {
					bpower = 1
				}
				shieldVal := math.Round(hp * bpower)
				inst.Vars["shielded"] = 1.0
				inst.Vars["shield_hp"] = shieldVal
				inst.Vars["shield_max"] = shieldVal
			}
		}
		// DDT and Deadly DDT inherently have camo + lead
		if layer == 351 || layer == 3351 {
			inst.Vars["camo"] = 1.0
			inst.Vars["lead"] = 1.0
		}
	}

	// if layer is not in any known sprite map, find the closest valid one
	if _, ok := bloonSprites[layer]; !ok {
		// check if it's close to a known layer (within 0.05)
		found := false
		for l := range bloonSprites {
			if math.Abs(layer-l) < 0.05 {
				found = true
				break
			}
		}
		if !found && layer > 8.5 && layer < 18 {
			// unknown layer between prismatic and ceramic — clamp to rainbow
			inst.Vars["bloonlayer"] = 8.0
			layer = 8
		}
	}

	// update sprite based on layer
	updateBloonSprite(inst)
	updateBloonDepth(inst)

	// --- Special bloon: per-type speed, scale and armour overrides ---
	stype := int(getVar(inst, "special_type"))
	bspeed := getGlobal(g, "bspeed")
	if bspeed == 0 {
		bspeed = 1
	}
	bpower := getGlobal(g, "bpower")
	if bpower == 0 {
		bpower = 1
	}

	var fast float64
	if stype > 0 {
		switch stype {
		case 1: // Stuffed
			inst.Vars["maxarmour"] = math.Ceil(layer * bpower)
			inst.ImageXScale = 0.95 + layer/20
			inst.ImageYScale = 0.95 + layer/20
			fast = (1.5 + layer/4) * bspeed
		case 2: // Ninja
			inst.Vars["maxarmour"] = math.Ceil(layer * 9 * bpower)
			inst.ImageXScale = 1.15 + layer/40
			inst.ImageYScale = 1.15 + layer/40
			fast = (3.5 + layer/8) * bspeed
		case 3: // Robo
			inst.Vars["maxarmour"] = math.Ceil(layer * 12 * bpower)
			inst.ImageXScale = 1.15 + layer/40
			inst.ImageYScale = 1.15 + layer/40
			fast = (2 + layer/16) * bspeed
		case 4: // Patrol
			inst.Vars["maxarmour"] = math.Ceil(60 * bpower)
			inst.ImageXScale = 1.4
			inst.ImageYScale = 1.4
			fast = 3 * bspeed
			// Patrol doesn't follow path — uses engine motion
			inst.Vars["path_name"] = ""
		case 5: // Barrier
			inst.Vars["maxarmour"] = math.Ceil(layer * 10 * bpower)
			inst.ImageXScale = 1.1 + layer/20
			inst.ImageYScale = 1.1 + layer/20
			fast = (1.2 + layer/10) * bspeed
		case 6: // Planetarium
			if layer >= 10 {
				inst.Vars["maxarmour"] = math.Ceil(300 * bpower)
			} else {
				inst.Vars["maxarmour"] = math.Ceil(layer * 30 * bpower)
			}
			inst.ImageXScale = 1 + layer/20
			inst.ImageYScale = 1 + layer/20
			fast = 0.85 * bspeed
		case 7: // Spectrum
			inst.Vars["maxarmour"] = math.Ceil(200 * bpower)
			inst.ImageXScale = 1.2
			inst.ImageYScale = 1.2
			fast = 2.4 * bspeed
		}

		// Status effects for special bloons
		if getVar(inst, "stun") == 1 {
			if stype == 2 || stype == 3 {
				fast = 0.3 * bspeed // Ninja/Robo: slow when stunned, not stopped
			} else {
				fast = 0
			}
		}
		if getVar(inst, "distraction") > 0 {
			distraction := getVar(inst, "distraction")
			if stype == 2 || stype == 3 {
				fast = -0.5 * distraction * bspeed
			} else {
				fast = bspeed * -3 * distraction
			}
		}
		if getVar(inst, "freeze") >= 1 {
			fast = 0
		}
		glue := getVar(inst, "glue")
		if glue > 0 {
			fast = fast * (0.6 - glue*0.1)
		}

		inst.Vars["fast"] = fast
	} else {
		// Normal bloon speed calculation
		speedMul := 1.6 // default red
		if s, ok := bloonSpeeds[layer]; ok {
			speedMul = s
		}

		tattered := getVar(inst, "tattered")
		fast = bspeed * speedMul * (tattered + 1)

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
	}

	// regrow mechanic: gradually regrow layers back to bloonmaxlayer
	if getVar(inst, "regrow") == 1 {
		if stype == 1 {
			// Stuffed: regrow armour, not layers
			maxArmour := getVar(inst, "maxarmour")
			currentArmour := getVar(inst, "armour")
			if currentArmour < maxArmour {
				regrowTimer := getVar(inst, "regrow_timer")
				regrowTimer++
				if regrowTimer >= 60 { // regrow 1 armour every second
					currentArmour++
					if currentArmour > maxArmour {
						currentArmour = maxArmour
					}
					inst.Vars["armour"] = currentArmour
					inst.Vars["regrow_timer"] = 0.0
				} else {
					inst.Vars["regrow_timer"] = regrowTimer
				}
			}
		} else {
			// Normal bloons: regrow layers
			maxLayer := getVar(inst, "bloonmaxlayer")
			if layer < maxLayer && layer < 8.5 { // only non-MOAB layers regrow
				regrowTimer := getVar(inst, "regrow_timer")
				regrowTimer++
				if regrowTimer >= 180 { // regrow every 3 seconds (180 frames)
					// regrow by 1 layer
					newLayer := layer + 1
					if newLayer > maxLayer {
						newLayer = maxLayer
					}
					inst.Vars["bloonlayer"] = newLayer
					inst.Vars["regrow_timer"] = 0.0
				} else {
					inst.Vars["regrow_timer"] = regrowTimer
				}
			}
		}
	}

	// Special bloon ability timers (replacement for GML alarm[0]/[1]/[2])
	if stype > 0 {
		specialBloonAbilityTick(inst, g, stype)
	}

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

			// update direction for rendering — only MOABs face their movement direction
			dir := g.PathMgr.GetDirectionAtProgress(pathName, progress)
			if layer >= 68.5 {
				inst.ImageAngle = dir
			}
		}
	}

	// Patrol bloon: custom free-flying movement with bouncing (NO path)
	if stype == 4 {
		// Initialize patrol movement on first step
		if getVar(inst, "patrol_init") == 0 {
			inst.Vars["patrol_init"] = 1.0
			// Remove path — Patrol flies freely, doesn't follow a track
			inst.Vars["path_name"] = ""
			// Start off the left edge, random Y (GML: x=-45, y=32+random(448))
			inst.X = -45
			inst.Y = 32 + rand.Float64()*448
			// Initial speed: 4 px/frame moving right (GML: speed=4, direction=0)
			spd := 4.0
			inst.HSpeed = spd
			inst.VSpeed = 0
			inst.Vars["patrol_life"] = 2400.0            // self-destruct timer (40s)
			inst.Vars["patrol_bounce"] = 1200.0 / spd    // first bounce timer (GML: alarm[0]=1200/speed)
		}

		// apply motion
		inst.X += inst.HSpeed
		inst.Y += inst.VSpeed

		// patrol bounce timer — reverse direction, new random Y, new speed
		bounceTimer := getVar(inst, "patrol_bounce") - 1
		if bounceTimer <= 0 {
			// GML alarm[0]: direction+=180, y=32+random(448), speed=2*bspeed+random(2*bspeed)
			inst.HSpeed = -inst.HSpeed
			spd := 2*bspeed + rand.Float64()*2*bspeed
			if inst.HSpeed > 0 {
				inst.HSpeed = spd
			} else {
				inst.HSpeed = -spd
			}
			inst.VSpeed = 0
			inst.Y = 32 + rand.Float64()*448
			inst.Vars["patrol_bounce"] = 1200.0 / spd
		} else {
			inst.Vars["patrol_bounce"] = bounceTimer
		}

		// lifetime timer — self-destruct after 2400 frames (40s)
		life := getVar(inst, "patrol_life") - 1
		if life <= 0 {
			g.InstanceMgr.Destroy(inst.ID)
			return
		}
		inst.Vars["patrol_life"] = life

		// fake path_position so End collision doesn't trigger
		inst.Vars["path_progress"] = 0.7
	}

	// visual variant frame selection (each bloon sprite has 12 frames):
	// 0=normal, 1=tattered, 2=shielded, 3=regrow, 4=camo, 5=lead,
	// 7=camo+lead, 8=regrow+tattered
	camo := getVar(inst, "camo")
	lead := getVar(inst, "lead")
	regrow := getVar(inst, "regrow")
	shielded := getVar(inst, "shielded")
	tattered := getVar(inst, "tattered")

	// MOAB-class bloons (layer >= 18): per-type animation + damage stages
	if layer >= 18 {
		shieldHP := getVar(inst, "shield_hp")
		maxShield := getVar(inst, "shield_max")
		if maxShield == 0 {
			// initialize max shield on first check
			if shieldHP > 0 {
				inst.Vars["shield_max"] = shieldHP
				maxShield = shieldHP
			}
		}

		ratio := 1.0
		if maxShield > 0 && shielded == 1 {
			ratio = shieldHP / maxShield
		}

		switch layer {
		case 351, 3351: // DDT, Deadly DDT — 3/1 frames, reversed logic (0=healthy, 1=medium, 2=most damaged)
			if ratio >= 2.0/3.0 {
				inst.ImageIndex = 0
			} else if ratio >= 1.0/3.0 {
				inst.ImageIndex = 1
			} else {
				inst.ImageIndex = 2
			}

		case 318: // LPZ — sprite swap between Shielded_LPZ_Spr and LPZ_Spr
			// continuous animation at 0.3 speed (no damage frames)
			if getVar(inst, "shield_hp") > 0 {
				inst.SpriteName = "Shielded_LPZ_Spr"
			} else {
				inst.SpriteName = "LPZ_Spr"
			}
			inst.ImageIndex = math.Mod(inst.ImageIndex+0.3, 9)

		case 918: // Storm LPZ — 2 frames, continuous animation
			inst.ImageIndex = math.Mod(inst.ImageIndex+0.3, 2)

		case 68.5: // HTA — 26 frames: 0-23 rotation animation, 24-25 damage
			if ratio <= 1.0/3.0 {
				inst.ImageIndex = 25
			} else if ratio <= 2.0/3.0 {
				inst.ImageIndex = 24
			} else {
				// animate frames 0-23
				idx := inst.ImageIndex + 0.5
				if idx > 23 {
					idx = 0
				}
				inst.ImageIndex = idx
			}

		case 5248: // ZOMG — 12 frames: 0-9 animation, 10-11 damage
			if ratio < 1.0/3.0 {
				inst.ImageIndex = 11
			} else if ratio < 2.0/3.0 {
				inst.ImageIndex = 10
			} else {
				// animate frames 0-9
				inst.ImageIndex = math.Mod(inst.ImageIndex+0.5, 10)
			}

		case 248, 2593: // Rocket Blimp, Mega BRC — 1 frame (static)
			inst.ImageIndex = 0

		case 10068.5: // Prismatic HTA — single-frame sprite, no animation
			inst.ImageIndex = 0
		case 17248: // Party Blimp — single-frame sprite, no animation
			inst.ImageIndex = 0
		default: // Mini-MOAB, MOAB, BFB, BRC, Ceramic, Brick — 11 frames: 0-8 anim, 9-10 damage
			if maxShield > 0 && shielded == 1 {
				if ratio < 1.0/3.0 {
					inst.ImageIndex = 10
				} else if ratio < 2.0/3.0 {
					inst.ImageIndex = 9
				} else {
					// animate frames 0-8
					inst.ImageIndex = math.Mod(inst.ImageIndex+0.5, 9)
				}
			} else {
				// no shield (ceramic/brick or shield broken) — just animate
				inst.ImageIndex = math.Mod(inst.ImageIndex+0.5, 9)
			}
		}
	} else if getVar(inst, "special_type") > 0 {
		// Special bloons use their own sprite with custom frame logic
		stype := int(getVar(inst, "special_type"))
		switch stype {
		case 1: // Stuffed — 5 frames, frame = layer-1 (0-4)
			f := layer - 1
			if f < 0 {
				f = 0
			}
			if f > 4 {
				f = 4
			}
			inst.ImageIndex = f
		case 2: // Ninja — 12 frames, frame = layer-1 (0-11)
			f := layer - 1
			if f < 0 {
				f = 0
			}
			if f > 11 {
				f = 11
			}
			inst.ImageIndex = f
		case 3: // Robo — 11 frames, frame = layer-1 (0-10)
			f := layer - 1
			if f < 0 {
				f = 0
			}
			if f > 10 {
				f = 10
			}
			inst.ImageIndex = f
		case 4: // Patrol — 1 frame
			inst.ImageIndex = 0
		case 5: // Barrier — 6 frames, frame = layer-1 (0-5)
			f := layer - 1
			if f < 0 {
				f = 0
			}
			if f > 5 {
				f = 5
			}
			inst.ImageIndex = f
		case 6: // Planetarium — 9 frames, animated
			inst.ImageIndex = math.Mod(inst.ImageIndex+0.3, 9)
		case 7: // Spectrum — 36 frames, animated
			inst.ImageIndex = math.Mod(inst.ImageIndex+0.5, 36)
		}
	} else if camo == 1 && lead == 1 {
		if tattered == 1 {
			inst.ImageIndex = 1
		} else {
			inst.ImageIndex = 7
		}
	} else if lead == 1 {
		if tattered == 1 {
			inst.ImageIndex = 1
		} else {
			inst.ImageIndex = 5
		}
	} else if camo == 1 {
		if tattered == 1 {
			inst.ImageIndex = 1
		} else {
			inst.ImageIndex = 4
		}
	} else if regrow == 1 {
		if tattered == 1 {
			inst.ImageIndex = 8
		} else {
			inst.ImageIndex = 3
		}
	} else if shielded == 1 {
		if tattered == 1 {
			inst.ImageIndex = 1
		} else {
			inst.ImageIndex = 2
		}
	} else if tattered == 1 {
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

// Destroy handles special bloon death effects (extra children spawning)
func (b *NormalBloonBranch) Destroy(inst *engine.Instance, g *engine.Game) {
	stype := int(getVar(inst, "special_type"))
	if stype == 0 {
		return // not a special bloon
	}

	layer := getVar(inst, "bloonlayer")

	switch stype {
	case 1: // Stuffed — spawns 15 Floaty_Branch at bloonlayer + 3 at bloonlayer+0.5
		// find nearest tower for direction
		target := findNearestTower(inst, g, 9999)
		for i := 0; i < 15; i++ {
			child := g.InstanceMgr.Create("Floaty_Branch", inst.X, inst.Y)
			if child == nil {
				continue
			}
			child.Vars["bloonlayer"] = layer
			child.Vars["bloonmaxlayer"] = layer
			spd := 1.2 + layer/2
			if target != nil {
				dx := target.X - inst.X
				dy := target.Y - inst.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > 0 {
					child.HSpeed = (dx/dist)*spd + (rand.Float64()*3 - 1.5)
					child.VSpeed = (dy/dist)*spd + (rand.Float64()*3 - 1.5)
				}
			} else {
				child.HSpeed = (rand.Float64() - 0.5) * spd * 2
				child.VSpeed = (rand.Float64() - 0.5) * spd * 2
			}
		}
		for i := 0; i < 3; i++ {
			child := g.InstanceMgr.Create("Floaty_Branch", inst.X, inst.Y)
			if child == nil {
				continue
			}
			child.Vars["bloonlayer"] = layer + 0.5
			child.Vars["bloonmaxlayer"] = layer + 0.5
			spd := 1.2 + layer/2
			if target != nil {
				dx := target.X - inst.X
				dy := target.Y - inst.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > 0 {
					child.HSpeed = (dx/dist)*spd + (rand.Float64()*3 - 1.5)
					child.VSpeed = (dy/dist)*spd + (rand.Float64()*3 - 1.5)
				}
			} else {
				child.HSpeed = (rand.Float64() - 0.5) * spd * 2
				child.VSpeed = (rand.Float64() - 0.5) * spd * 2
			}
		}

	case 2: // Ninja — spawns camo bloon at layer+0.5
		childLayer := layer + 0.5
		if layer >= 10 {
			childLayer = 8
		}
		if childLayer >= 1 {
			child := spawnBloonChild(inst, childLayer, g)
			if child != nil {
				child.Vars["camo"] = 1.0
			}
		}

	case 3: // Robo — spawns lead bloon at layer+0.5
		childLayer := layer + 0.5
		if layer >= 10 {
			childLayer = 8
		}
		if childLayer >= 1 {
			child := spawnBloonChild(inst, childLayer, g)
			if child != nil {
				child.Vars["lead"] = 1.0
			}
		}

	case 4: // Patrol — just pops, no children
		// nothing

	case 5: // Barrier — spawns 6 Block_Branch on death
		for i := 0; i < 6; i++ {
			block := g.InstanceMgr.Create("Block_Branch", inst.X, inst.Y)
			if block != nil {
				shieldVal := layer * 3
				block.Vars["bloonlayer"] = shieldVal
				block.Vars["shield"] = shieldVal
				block.Vars["maxshield"] = shieldVal
				spd := 0.2 + layer/6
				angle := rand.Float64() * 360
				angleRad := angle * math.Pi / 180
				block.HSpeed = math.Cos(angleRad) * spd
				block.VSpeed = -math.Sin(angleRad) * spd
			}
		}

	case 6: // Planetarium — spawns 12 Satellites + 1 Comet on death
		for i := 0; i < 12; i++ {
			sat := g.InstanceMgr.Create("Satellite_Bloon_Branch", inst.X, inst.Y)
			if sat != nil {
				shieldVal := layer * 3
				sat.Vars["bloonlayer"] = shieldVal
				sat.Vars["shield"] = shieldVal
				sat.Vars["maxshield"] = shieldVal
				angle := float64(i) * 30 * math.Pi / 180
				sat.HSpeed = math.Cos(angle)
				sat.VSpeed = -math.Sin(angle)
				sat.Speed = 1
			}
		}
		comet := g.InstanceMgr.Create("Comet_Bloon_Branch", inst.X, inst.Y)
		if comet != nil {
			comet.Vars["bloonlayer"] = layer * 3
			comet.Vars["shield"] = layer * 10
			comet.Vars["maxshield"] = layer * 10
		}

	case 7: // Spectrum — spawns 1 bloon at layer-0.5
		childLayer := layer - 0.5
		if childLayer >= 1 {
			spawnBloonChild(inst, childLayer, g)
		}
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

	// Bloons should stay inside the gameplay area and never overlap the HUD.
	playfield := screen.SubImage(image.Rect(0, 0, 864, 480)).(*ebiten.Image)
	engine.DrawSpriteExt(playfield, spr.Frames[frameIdx], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, scale*inst.ImageXScale, scale*inst.ImageYScale,
		inst.ImageAngle, inst.ImageAlpha)

	// show bloon info tooltip when the setting is on and mouse hovers the bloon
	if getGlobal(g, "blooninfo") == 1 && g.IsMouseOverInstance(inst) {
		stype := int(getVar(inst, "special_type"))
		var name string
		if stype > 0 {
			specialNames := map[int]string{
				1: "Stuffed", 2: "Ninja", 3: "Robo",
				4: "Patrol", 5: "Barrier", 6: "Planetarium", 7: "Spectrum",
			}
			name = specialNames[stype]
		} else {
			name = bloonLayerName(layer)
		}
		info := fmt.Sprintf("%s (%.1f)", name, layer)
		drawBloonInfoTag(screen, g, info, inst.X, inst.Y-20)
	}
}

// bloonLayerName returns a readable name for the bloon layer
func bloonLayerName(layer float64) string {
	// exact matches for nightmare/special MOABs first
	exactNames := map[float64]string{
		248:     "Rocket Blimp",
		918:     "Storm LPZ",
		2593:    "Mega BRC",
		3351:    "Deadly DDT",
		10068.5: "Prismatic HTA",
		17248:   "Party Blimp",
	}
	if name, ok := exactNames[layer]; ok {
		return name
	}

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
	case layer <= 18:
		return "Ceramic"
	case layer <= 48:
		return "Brick"
	case layer <= 68.5:
		return "HTA"
	case layer <= 93:
		return "Mini-MOAB"
	case layer <= 318:
		return "LPZ"
	case layer <= 348:
		return "MOAB"
	case layer <= 351:
		return "DDT"
	case layer <= 593:
		return "BRC"
	case layer <= 1248:
		return "BFB"
	case layer <= 5248:
		return "ZOMG"
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
	// Special bloons use their own sprite, not the layer's
	if specialSprite, ok := inst.Vars["special_sprite"].(string); ok && specialSprite != "" {
		inst.SpriteName = specialSprite
		return
	}

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

	// --- shield handling (MOAB-class + fortified) ---
	if getVar(bloon, "shielded") == 1 {
		shieldHP := getVar(bloon, "shield_hp")
		if shieldHP > 0 {
			shieldHP -= lp
			if shieldHP <= 0 {
				bloon.Vars["shield_hp"] = 0.0
				bloon.Vars["shielded"] = 0.0
				// shield broken — damage passes through
				remaining := -shieldHP
				if remaining > 0 {
					popBloon(bloon, remaining, g)
				}
			} else {
				bloon.Vars["shield_hp"] = shieldHP
			}
			return
		}
		// shield_hp already 0 — remove shield flag and continue to normal pop
		bloon.Vars["shielded"] = 0.0
	}

	// --- Special bloon: Spectrum damage cap (max 100/hit) ---
	if getVar(bloon, "damage_cap") > 0 && lp > getVar(bloon, "damage_cap") {
		lp = getVar(bloon, "damage_cap")
	}

	// --- Special bloon: Armour acts as extra HP pool (like ceramic HP) ---
	armour := getVar(bloon, "armour")
	if armour > 0 {
		armour -= lp
		if armour <= 0 {
			bloon.Vars["armour"] = 0.0
			// Special bloons: armour IS their HP — destroy when empty
			if getVar(bloon, "special_type") > 0 {
				g.GlobalVars["cashwavereward"] = cashReward + layer
				g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + layer
				g.InstanceMgr.Destroy(bloon.ID)
				playBloonPopSound(g)
				return
			}
			// Normal bloons: armour broken, pass remaining damage through
			lp = -armour
			if lp <= 0 {
				playBloonPopSound(g)
				return
			}
		} else {
			bloon.Vars["armour"] = armour
			playBloonPopSound(g)
			return // damage fully absorbed by armour
		}
	}

	// --- MOAB-class bloons (layer >= 68.5): reduce to children on death ---
	if layer >= 68.5 {
		moabChildren := map[float64]struct {
			childLayer float64
			count      int
		}{
			5248: {1248, 4},   // ZOMG → 4 BFB
			1248: {348, 4},    // BFB → 4 MOAB
			348:  {18, 4},     // MOAB → 4 Ceramic
			93:   {8, 4},      // Mini-MOAB → 4 Rainbow
			593:  {348, 2},    // BRC → 2 MOAB
			351:  {18, 6},     // DDT → 6 Ceramic
			318:  {93, 4},     // LPZ → 4 Mini-MOAB
			68.5: {8.5, 6},    // HTA → 6 Prismatic
			// Nightmare MOAB class
			248:     {8, 4},      // Rocket Blimp → 4 Rainbow
			918:     {93, 4},     // Storm LPZ → 4 Mini-MOAB
			2593:    {348, 4},    // Mega BRC → 4 MOAB
			3351:    {18, 6},     // Deadly DDT → 6 Ceramic
			10068.5: {8.5, 8},    // Prismatic HTA → 8 Prismatic
			17248:   {1248, 4},   // Party Blimp → 4 BFB
		}

		if children, ok := moabChildren[layer]; ok {
			g.GlobalVars["cashwavereward"] = cashReward + layer
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + layer
			for i := 0; i < children.count; i++ {
				spawnBloonChild(bloon, children.childLayer, g)
			}
			g.InstanceMgr.Destroy(bloon.ID)
			playBloonPopSound(g)
			return
		}

		// unknown MOAB-class: just destroy
		g.GlobalVars["cashwavereward"] = cashReward + layer
		g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + layer
		g.InstanceMgr.Destroy(bloon.ID)
		playBloonPopSound(g)
		return
	}

	// --- Ceramic (18) and Brick (48): multi-hit then split ---
	if layer == 18 || layer == 48 {
		hp := getVar(bloon, "ceramic_hp")
		if hp == 0 {
			// first hit: initialize HP (Ceramic=10, Brick=30)
			if layer == 18 {
				hp = 10
			} else {
				hp = 30
			}
			bloon.Vars["ceramic_hp"] = hp
		}
		hp -= lp
		if hp <= 0 {
			// break into children
			g.GlobalVars["cashwavereward"] = cashReward + layer
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + layer
			childLayer := 8.0  // Ceramic → Rainbow
			childCount := 2
			if layer == 48 {
				childLayer = 18.0 // Brick → Ceramic
				childCount = 3
			}
			for i := 0; i < childCount; i++ {
				spawnBloonChild(bloon, childLayer, g)
			}
			g.InstanceMgr.Destroy(bloon.ID)
		} else {
			bloon.Vars["ceramic_hp"] = hp
			g.GlobalVars["cashwavereward"] = cashReward + lp
			g.GlobalVars["monkeypop"] = getGlobal(g, "monkeypop") + lp
		}
		playBloonPopSound(g)
		return
	}

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
func spawnBloonChild(parent *engine.Instance, layer float64, g *engine.Game) *engine.Instance {
	child := g.InstanceMgr.Create("Normal_Bloon_Branch", parent.X, parent.Y)
	if child == nil {
		return nil
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
	return child
}

// registerBloonBehaviors registers bloon-related behaviors
func RegisterBloonBehaviors(im *engine.InstanceManager) {
	im.RegisterBehavior("Normal_Bloon_Branch", func() engine.InstanceBehavior { return &NormalBloonBranch{} })
}
