package behaviors

import (
	"math"
	"math/rand"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

// ============================================================================
// Special bloon ability tick — called once per Step for special bloons
// Replaces GML alarm[0]/[1] cycles with frame counters
// ============================================================================

func specialBloonAbilityTick(inst *engine.Instance, g *engine.Game, stype int) {
	switch stype {
	case 2: // Ninja — shuriken attack
		timer := getVar(inst, "ability_timer")
		timer--
		if timer <= 0 {
			layer := getVar(inst, "bloonlayer")
			attackRange := 200 + layer*10
			target := findNearestTower(inst, g, attackRange)
			if target != nil {
				shuriken := g.InstanceMgr.Create("Ninja_Bloon_Shuriken", inst.X, inst.Y)
				if shuriken != nil {
					shuriken.Vars["potency"] = 10 + layer*6
					shuriken.Vars["target_id"] = float64(target.ID)
					dx := target.X - inst.X
					dy := target.Y - inst.Y
					dist := math.Sqrt(dx*dx + dy*dy)
					if dist > 0 {
						shuriken.HSpeed = (dx / dist) * 6
						shuriken.VSpeed = (dy / dist) * 6
					}
				}
			}
			inst.Vars["ability_timer"] = 90.0 + rand.Float64()*45
		} else {
			inst.Vars["ability_timer"] = timer
		}

	case 3: // Robo — laser attack
		timer := getVar(inst, "ability_timer")
		timer--
		if timer <= 0 {
			layer := getVar(inst, "bloonlayer")
			attackRange := 200 + layer*10
			target := findNearestTower(inst, g, attackRange)
			if target != nil {
				laser := g.InstanceMgr.Create("Robo_Bloon_Laser", inst.X, inst.Y)
				if laser != nil {
					laser.Vars["potency"] = 2 + layer
					dx := target.X - inst.X
					dy := target.Y - inst.Y
					dist := math.Sqrt(dx*dx + dy*dy)
					if dist > 0 {
						laser.HSpeed = (dx / dist) * 8
						laser.VSpeed = (dy / dist) * 8
						laser.ImageAngle = math.Atan2(-dy, dx) * 180 / math.Pi
					}
				}
			}
			inst.Vars["ability_timer"] = 400.0 + rand.Float64()*200
		} else {
			inst.Vars["ability_timer"] = timer
		}

	case 4: // Patrol — shooting at towers
		timer := getVar(inst, "ability_timer")
		timer--
		if timer <= 0 {
			attackRange := 200.0
			target := findNearestTower(inst, g, attackRange)
			if target != nil {
				shot := g.InstanceMgr.Create("Patrol_Bloon_Shot", inst.X, inst.Y)
				if shot != nil {
					shot.Vars["potency"] = 20.0
					dx := target.X - inst.X
					dy := target.Y - inst.Y
					dist := math.Sqrt(dx*dx + dy*dy)
					if dist > 0 {
						shot.HSpeed = (dx / dist) * 12
						shot.VSpeed = (dy / dist) * 12
						shot.ImageAngle = math.Atan2(-dy, dx) * 180 / math.Pi
					}
				}
			}
			inst.Vars["ability_timer"] = 30.0
		} else {
			inst.Vars["ability_timer"] = timer
		}

	case 5: // Barrier — spawns Block_Branch obstacles toward random tower
		timer := getVar(inst, "ability_timer")
		timer--
		if timer <= 0 {
			layer := getVar(inst, "bloonlayer")
			target := findRandomTower(g)
			if target != nil {
				for i := 0; i < 3; i++ {
					block := g.InstanceMgr.Create("Block_Branch", inst.X, inst.Y)
					if block != nil {
						shieldVal := layer * 3
						block.Vars["bloonlayer"] = shieldVal
						block.Vars["shield"] = shieldVal
						block.Vars["maxshield"] = shieldVal
						spd := 0.2 + layer/6
						dx := target.X - inst.X
						dy := target.Y - inst.Y
						dist := math.Sqrt(dx*dx + dy*dy)
						if dist > 0 {
							block.HSpeed = (dx/dist)*spd + (rand.Float64()*2 - 1)
							block.VSpeed = (dy/dist)*spd + (rand.Float64()*2 - 1)
						}
					}
				}
			}
			inst.Vars["ability_timer"] = 200.0 + rand.Float64()*30
		} else {
			inst.Vars["ability_timer"] = timer
		}

	case 6: // Planetarium — spawns orbiting Satellites
		timer := getVar(inst, "ability_timer")
		timer--
		if timer <= 0 {
			progress := getVar(inst, "path_progress")
			if progress > 0.125 && progress < 0.9 {
				layer := getVar(inst, "bloonlayer")
				for i := 0; i < 4; i++ {
					sat := g.InstanceMgr.Create("Satellite_Bloon_Branch", inst.X, inst.Y)
					if sat != nil {
						shieldVal := layer * 3
						sat.Vars["bloonlayer"] = shieldVal
						sat.Vars["shield"] = shieldVal
						sat.Vars["maxshield"] = shieldVal
						angle := float64(i) * 90 * math.Pi / 180
						sat.HSpeed = math.Cos(angle)
						sat.VSpeed = -math.Sin(angle)
					}
				}
			}
			inst.Vars["ability_timer"] = 130.0
		} else {
			inst.Vars["ability_timer"] = timer
		}
	}
}

// findNearestTower finds the nearest tower instance within the given range
func findNearestTower(bloon *engine.Instance, g *engine.Game, attackRange float64) *engine.Instance {
	var nearest *engine.Instance
	bestDist := attackRange

	for _, name := range allSelectableTowers {
		for _, tower := range g.InstanceMgr.FindByObject(name) {
			if tower.Destroyed {
				continue
			}
			dx := tower.X - bloon.X
			dy := tower.Y - bloon.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < bestDist {
				bestDist = dist
				nearest = tower
			}
		}
	}
	return nearest
}

// findRandomTower returns a random tower instance, or nil if none exist
func findRandomTower(g *engine.Game) *engine.Instance {
	var towers []*engine.Instance
	for _, name := range allSelectableTowers {
		for _, tower := range g.InstanceMgr.FindByObject(name) {
			if !tower.Destroyed {
				towers = append(towers, tower)
			}
		}
	}
	if len(towers) == 0 {
		return nil
	}
	return towers[rand.Intn(len(towers))]
}

// ============================================================================
// Floaty_Branch — free-floating bloon spawned by Stuffed on death
// Uses normal bloon sprites based on layer, self-destructs after timer
// ============================================================================

type FloatyBranch struct {
	engine.DefaultBehavior
}

func (b *FloatyBranch) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["bloonlayer"] = 1.0
	inst.Vars["bloonmaxlayer"] = 1.0
	inst.Vars["fast"] = 1.5
	inst.Vars["glue"] = 0.0
	inst.Vars["freeze"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["distraction"] = 0.0
	inst.Vars["tattered"] = 0.0
	inst.Vars["camo"] = 0.0
	inst.Vars["lead"] = 0.0
	inst.Vars["regrow"] = 0.0
	inst.Vars["shielded"] = 0.0
	inst.Friction = 0.02
	// self-destruct timer: 5-7 seconds
	inst.Alarms[0] = 300 + rand.Intn(120)
}

func (b *FloatyBranch) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0

	// friction slowdown, minimum speed
	if inst.Speed <= 0.4 {
		inst.Friction = 0
		inst.Speed = 0.4
	}

	layer := getVar(inst, "bloonlayer")
	// clamp layers
	if layer > 5 && layer < 6 {
		layer = 5
		inst.Vars["bloonlayer"] = layer
	}
	if layer > 9 {
		layer = 8
		inst.Vars["bloonlayer"] = layer
	}

	// update sprite from layer (same as normal bloons)
	updateBloonSprite(inst)

	// variant frame
	if getVar(inst, "tattered") == 1 {
		inst.ImageIndex = 1
	} else {
		inst.ImageIndex = 0
	}

	// face movement direction
	inst.ImageAngle = inst.Direction - 90

	// destroy if layer depleted
	if layer < 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *FloatyBranch) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		// self-destruct
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *FloatyBranch) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frameIdx := int(inst.ImageIndex) % len(spr.Frames)
	if frameIdx < 0 {
		frameIdx = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frameIdx], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale,
		inst.ImageAngle, inst.ImageAlpha)
}

// ============================================================================
// Ninja_Bloon_Shuriken — homing projectile that stuns towers
// ============================================================================

type NinjaBloonShuriken struct {
	engine.DefaultBehavior
}

func (b *NinjaBloonShuriken) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Ninja_Bloon_Shuriken_spr"
	inst.ImageSpeed = 0
	inst.Depth = -5
}

func (b *NinjaBloonShuriken) Step(inst *engine.Instance, g *engine.Game) {
	// home toward target
	targetID := int(getVar(inst, "target_id"))
	target := g.InstanceMgr.GetByID(targetID)
	if target != nil && !target.Destroyed {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0 {
			inst.HSpeed = (dx / dist) * 4
			inst.VSpeed = (dy / dist) * 4
		}

		// check collision with target
		if dist < 16 {
			// stun the tower
			potency := getVar(inst, "potency")
			stunVal, _ := target.Vars["bloon_stun"].(float64)
			target.Vars["bloon_stun"] = stunVal + potency
			g.InstanceMgr.Destroy(inst.ID)
			return
		}
	}

	// apply motion
	inst.X += inst.HSpeed
	inst.Y += inst.VSpeed

	// rotate visually
	inst.ImageAngle += 15

	// destroy if out of room
	if inst.X < -50 || inst.X > 920 || inst.Y < -50 || inst.Y > 530 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *NinjaBloonShuriken) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frameIdx := int(inst.ImageIndex) % len(spr.Frames)
	if frameIdx < 0 {
		frameIdx = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frameIdx], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale,
		inst.ImageAngle, inst.ImageAlpha)
}

// ============================================================================
// Robo_Bloon_Laser — piercing laser that stuns towers (passes through)
// ============================================================================

type RoboBloonLaser struct {
	engine.DefaultBehavior
}

func (b *RoboBloonLaser) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Robo_Bloon_Laser_spr"
	inst.ImageSpeed = 0
	inst.Depth = -5
}

func (b *RoboBloonLaser) Step(inst *engine.Instance, g *engine.Game) {
	// apply motion
	inst.X += inst.HSpeed
	inst.Y += inst.VSpeed

	// face movement direction
	inst.ImageAngle = math.Atan2(-inst.VSpeed, inst.HSpeed) * 180 / math.Pi

	// check collision with towers (piercing — hits all in path, but each only once)
	potency := getVar(inst, "potency")
	for _, name := range allSelectableTowers {
		for _, tower := range g.InstanceMgr.FindByObject(name) {
			if tower.Destroyed {
				continue
			}
			dx := tower.X - inst.X
			dy := tower.Y - inst.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < 20 {
				// check if already hit this tower
				hitKey := hitKeyForID(tower.ID)
				if getVar(inst, hitKey) == 0 {
					inst.Vars[hitKey] = 1.0
					stunVal, _ := tower.Vars["bloon_stun"].(float64)
					tower.Vars["bloon_stun"] = stunVal + potency
				}
			}
		}
	}

	// destroy if out of room
	if inst.X < -50 || inst.X > 920 || inst.Y < -50 || inst.Y > 530 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *RoboBloonLaser) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frameIdx := int(inst.ImageIndex) % len(spr.Frames)
	if frameIdx < 0 {
		frameIdx = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frameIdx], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale,
		inst.ImageAngle, inst.ImageAlpha)
}

// hitKeyForID generates a per-instance hit tracking key
func hitKeyForID(id int) string {
	return "hit_" + string(rune('0'+id%10000))
}

// ============================================================================
// Patrol_Bloon_Shot — shot that stuns towers on contact
// ============================================================================

type PatrolBloonShot struct {
	engine.DefaultBehavior
}

func (b *PatrolBloonShot) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Patrol_Bloon_Shot_Spr"
	inst.ImageSpeed = 0
	inst.Depth = -5
}

func (b *PatrolBloonShot) Step(inst *engine.Instance, g *engine.Game) {
	// apply motion
	inst.X += inst.HSpeed
	inst.Y += inst.VSpeed
	inst.ImageAngle = math.Atan2(-inst.VSpeed, inst.HSpeed) * 180 / math.Pi

	// check collision with towers
	potency := getVar(inst, "potency")
	for _, name := range allSelectableTowers {
		for _, tower := range g.InstanceMgr.FindByObject(name) {
			if tower.Destroyed {
				continue
			}
			dx := tower.X - inst.X
			dy := tower.Y - inst.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist < 20 {
				stunVal, _ := tower.Vars["bloon_stun"].(float64)
				tower.Vars["bloon_stun"] = stunVal + potency
				g.InstanceMgr.Destroy(inst.ID)
				return
			}
		}
	}

	// destroy if out of room
	if inst.X < -50 || inst.X > 920 || inst.Y < -50 || inst.Y > 530 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *PatrolBloonShot) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	drawSimpleSprite(inst, screen, g)
}

// ============================================================================
// Block_Branch — shield-based obstacle spawned by Barrier bloon
// Self-destructs after timer. Takes damage from projectiles via shield.
// ============================================================================

type BlockBranch struct {
	engine.DefaultBehavior
}

func (b *BlockBranch) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Blocker_Bloon_Spr"
	inst.ImageSpeed = 0
	inst.ImageXScale = 1.15
	inst.ImageYScale = 1.15
	inst.Depth = -1
	inst.Friction = 0.02
	// self-destruct timer
	inst.Alarms[0] = 300 + rand.Intn(120)
}

func (b *BlockBranch) Step(inst *engine.Instance, g *engine.Game) {
	// minimum friction
	if inst.Speed <= 0.4 {
		inst.Friction = 0
		inst.Speed = 0.4
	}

	// frame based on bloonlayer
	layer := getVar(inst, "bloonlayer")
	inst.ImageIndex = math.Max(0, layer/3-1)

	// face movement direction
	inst.ImageAngle = inst.Direction - 90

	// destroy if shield depleted
	shield := getVar(inst, "shield")
	if shield < 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *BlockBranch) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *BlockBranch) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	drawSimpleSprite(inst, screen, g)
}

// ============================================================================
// Satellite_Bloon_Branch — orbiting shield-based bloon from Planetarium
// Gradually spirals outward, shield-based HP
// ============================================================================

type SatelliteBloonBranch struct {
	engine.DefaultBehavior
}

func (b *SatelliteBloonBranch) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Satellite_Bloon_Spr"
	inst.ImageSpeed = 0
	inst.Depth = -1
	// self-destruct after 50 seconds
	inst.Alarms[0] = 3000
}

func (b *SatelliteBloonBranch) Step(inst *engine.Instance, g *engine.Game) {
	// orbit: rotate direction and gradually accelerate
	inst.Direction += 1.5
	dirRad := inst.Direction * math.Pi / 180
	spd := inst.Speed + 0.005/(inst.Speed+0.01)
	inst.Speed = spd
	inst.HSpeed = math.Cos(dirRad) * spd
	inst.VSpeed = -math.Sin(dirRad) * spd

	inst.X += inst.HSpeed
	inst.Y += inst.VSpeed

	// visual
	layer := getVar(inst, "bloonlayer")
	inst.ImageIndex = layer / 15
	inst.ImageXScale = 0.9 + layer/150
	inst.ImageYScale = 0.9 + layer/150

	// destroy if shield depleted
	if getVar(inst, "shield") < 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}

	// destroy if out of bounds (generous margin)
	if inst.X < -100 || inst.X > 1000 || inst.Y < -100 || inst.Y > 600 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *SatelliteBloonBranch) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *SatelliteBloonBranch) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	drawSimpleSprite(inst, screen, g)
}

// ============================================================================
// Comet_Bloon_Branch — homes toward End point, spawned by Planetarium on death
// ============================================================================

type CometBloonBranch struct {
	engine.DefaultBehavior
}

func (b *CometBloonBranch) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Comet_Bloon_Spr"
	inst.ImageSpeed = 0
	inst.Depth = -1
	inst.Speed = 0.5
	inst.Alarms[0] = 3000
}

func (b *CometBloonBranch) Step(inst *engine.Instance, g *engine.Game) {
	// home toward End
	ends := g.InstanceMgr.FindByObject("End")
	if len(ends) > 0 {
		target := ends[0]
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		spd := inst.Speed + 0.005/(inst.Speed+0.01)
		inst.Speed = spd
		if dist > 0 {
			inst.HSpeed = (dx / dist) * spd
			inst.VSpeed = (dy / dist) * spd
		}
	}

	inst.X += inst.HSpeed
	inst.Y += inst.VSpeed
	inst.ImageAngle = math.Atan2(-inst.VSpeed, inst.HSpeed) * 180 / math.Pi

	layer := getVar(inst, "bloonlayer")
	inst.ImageIndex = layer / 30
	inst.ImageXScale = 1.3 + layer/150
	inst.ImageYScale = 1.3 + layer/150

	// destroy if shield depleted
	if getVar(inst, "shield") < 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *CometBloonBranch) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

func (b *CometBloonBranch) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	drawSimpleSprite(inst, screen, g)
}

// ============================================================================
// Shared drawing helper
// ============================================================================

func drawSimpleSprite(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frameIdx := int(inst.ImageIndex) % len(spr.Frames)
	if frameIdx < 0 {
		frameIdx = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frameIdx], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale,
		inst.ImageAngle, inst.ImageAlpha)
}

// ============================================================================
// Register all special bloon helper behaviors
// ============================================================================

func RegisterSpecialBloonBehaviors(im *engine.InstanceManager) {
	im.RegisterBehavior("Floaty_Branch", func() engine.InstanceBehavior { return &FloatyBranch{} })
	im.RegisterBehavior("Ninja_Bloon_Shuriken", func() engine.InstanceBehavior { return &NinjaBloonShuriken{} })
	im.RegisterBehavior("Robo_Bloon_Laser", func() engine.InstanceBehavior { return &RoboBloonLaser{} })
	im.RegisterBehavior("Patrol_Bloon_Shot", func() engine.InstanceBehavior { return &PatrolBloonShot{} })
	im.RegisterBehavior("Block_Branch", func() engine.InstanceBehavior { return &BlockBranch{} })
	im.RegisterBehavior("Satellite_Bloon_Branch", func() engine.InstanceBehavior { return &SatelliteBloonBranch{} })
	im.RegisterBehavior("Comet_Bloon_Branch", func() engine.InstanceBehavior { return &CometBloonBranch{} })
}
