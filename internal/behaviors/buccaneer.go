package behaviors

import (
	"math"
	"math/rand"

	"btdx/internal/engine"
)

// Monkey Buccaneer tower family — water pirate tower (ID 11).
// Uses a form-based system like Monkey Sub.

// ─────────────────────────────────────────────────────────────────────────────
// Form data
// ─────────────────────────────────────────────────────────────────────────────

type buccaneerForm struct {
	Range       float64
	Reload      float64 // alarm[0] interval for main fire
	LP          float64
	PP          float64
	CamoDet     float64 // 1 = detects camo
	LeadDet     float64 // 1 = pops lead
	GrapeCycle  int     // >0: every N shots, fire grape volley
	GrapeCount  int     // how many grapes per volley (default 4)
	GrapeLP     float64
	GrapePP     float64
	BombCycle   int // >0: every N shots, also fire Ship_Bomb
	BombLP      float64
	BombPP      float64
	Projectile  string  // main dart projectile name
	DartSpeed   float64 // main dart speed (default 25)
	AgentObj    string  // spawned agent object name (Swashbuckler, etc.)
	AgentRate   int     // alarm[1] interval for spawning agents
	AbilityMax  float64 // >0: ability charge cap
	BurstCount  int     // >0: every BurstCycle shots, fire this many extra darts
	BurstCycle  int     // how often to do a burst (e.g., every 3rd or 8th shot)
	HarpoonRate int     // >0: alarm[3] interval for harpoon attack
	HarpoonLP   float64
	HarpoonPP   float64
	// Dreadnaut path: accelerating projectile
	Accelerate   bool // projectile accelerates (speed += 0.3 per step)
	CursedRate   int  // >0: alarm[2] interval for Cursed_Dart spawn (Ghost_Ship path)
	CursedDartLP float64
	CursedDartPP float64
}

var buccaneerForms = map[string]buccaneerForm{
	// linear chain ──────────────────────────────────────────────────────────────
	"Monkey_Buccaneer": {
		Range: 130, Reload: 25, LP: 1, PP: 5, Projectile: "Buccaneer_Dart", DartSpeed: 25,
	},
	"Grape_Shot": {
		Range: 130, Reload: 25, LP: 1, PP: 5, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
	},
	"Crows_Nest": {
		Range: 140, Reload: 18, LP: 1, PP: 5, CamoDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
	},
	// path 1 (left) — Swashbucklers ──────────────────────────────────────────
	"Swashbucklers": {
		Range: 140, Reload: 18, LP: 1, PP: 5, CamoDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		AgentObj: "Swashbuckler", AgentRate: 210,
	},
	"Monkey_Pirates": {
		Range: 140, Reload: 18, LP: 1, PP: 5, CamoDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		AgentObj: "Monkey_Pirate", AgentRate: 200,
	},
	"Pirate_Captain_Ship": {
		Range: 144, Reload: 18, LP: 1, PP: 5, CamoDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		AgentObj: "Monkey_Pirate", AgentRate: 160, AbilityMax: 41,
	},
	// path 2 (middle) — Destroyer ────────────────────────────────────────────
	"Destroyer": {
		Range: 140, Reload: 4, LP: 1, PP: 5, CamoDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 3, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
	},
	"Supreme_Battleship": {
		Range: 144, Reload: 2, LP: 1, PP: 5, CamoDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 3, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		BurstCount: 10, BurstCycle: 8,
	},
	"Aircraft_Carrier": {
		Range: 160, Reload: 2, LP: 1, PP: 5, CamoDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 3, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		BurstCount: 10, BurstCycle: 6,
		AgentObj: "Boat_Plane", AgentRate: 230,
	},
	// path 3 (right) — Cannon_Ship ───────────────────────────────────────────
	"Cannon_Ship": {
		Range: 145, Reload: 17, LP: 1, PP: 5, CamoDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		BombCycle: 3, BombLP: 1, BombPP: 50,
	},
	"Harpoon_Ship": {
		Range: 145, Reload: 17, LP: 1, PP: 5, CamoDet: 1, LeadDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		BombCycle: 3, BombLP: 1, BombPP: 50,
		HarpoonRate: 54, HarpoonLP: 6, HarpoonPP: 12,
	},
	"MOAB_Takedown": {
		Range: 148, Reload: 17, LP: 1, PP: 5, CamoDet: 1, LeadDet: 1, Projectile: "Buccaneer_Dart", DartSpeed: 25,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		BombCycle: 3, BombLP: 1, BombPP: 50,
		HarpoonRate: 47, HarpoonLP: 6, HarpoonPP: 12, AbilityMax: 57,
	},
	// secret path — Dreadnaut ────────────────────────────────────────────────
	"Dreadnaut_Ship": {
		Range: 140, Reload: 12, LP: 1, PP: 10, CamoDet: 1, LeadDet: 1,
		Projectile: "Dreadnaut_Dart", DartSpeed: 3, Accelerate: true,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		AgentObj: "Dread_Monkey", AgentRate: 210,
	},
	"Cursed_Pirate_Ship": {
		Range: 160, Reload: 12, LP: 1, PP: 10, CamoDet: 1, LeadDet: 1,
		Projectile: "Dreadnaut_Dart", DartSpeed: 3, Accelerate: true,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		AgentObj: "Cursed_Monkey", AgentRate: 210,
		CursedRate: 3, CursedDartLP: 1, CursedDartPP: 20,
	},
	"Ghost_Ship": {
		Range: 160, Reload: 8, LP: 3, PP: 25, CamoDet: 1, LeadDet: 1,
		Projectile: "Dreadnaut_Dart", DartSpeed: 3, Accelerate: true,
		GrapeCycle: 1, GrapeCount: 4, GrapeLP: 1, GrapePP: 2,
		AgentObj: "Haunted_Monkey", AgentRate: 150,
		CursedRate: 3, CursedDartLP: 2, CursedDartPP: 30,
	},
}

func buccaneerUpgradeName(inst *engine.Instance) string {
	return upgradeName(inst, "Monkey_Buccaneer")
}

func buccaneerFormFor(inst *engine.Instance) buccaneerForm {
	if f, ok := buccaneerForms[buccaneerUpgradeName(inst)]; ok {
		return f
	}
	return buccaneerForms["Monkey_Buccaneer"]
}

// distToNearestTrackPath returns the minimum distance from (x,y) to any bloon
// path segment on the current track (unconstrained).
func distToNearestTrackPath(x, y float64, g *engine.Game) float64 {
	d, _, _ := nearestTrackPoint(x, y, g, 0, 0, 0)
	return d
}

// nearestTrackPoint returns the distance and coordinates of the closest point
// on any bloon path segment for the current track.
// If cx,cy,maxR are nonzero, only considers path segments within maxR of (cx,cy).
func nearestTrackPoint(x, y float64, g *engine.Game, cx, cy, maxR float64) (dist, nx, ny float64) {
	track := int(getGlobal(g, "track"))
	paths, ok := trackPaths[track]
	if !ok || g.PathMgr == nil {
		return 999, x, y
	}
	best := 999.0
	bx, by := x, y
	for _, pName := range paths {
		pa := g.PathMgr.Get(pName)
		if pa == nil || len(pa.CurvePoints) < 2 {
			continue
		}
		for i := 1; i < len(pa.CurvePoints); i++ {
			ax := pa.CurvePoints[i-1].X
			ay := pa.CurvePoints[i-1].Y
			ex := pa.CurvePoints[i].X
			ey := pa.CurvePoints[i].Y
			d, px, py := pointSegDistXY(x, y, ax, ay, ex, ey)
			// if radius constraint, skip segments whose closest point is far from tower
			if maxR > 0 {
				tdx := px - cx
				tdy := py - cy
				if math.Sqrt(tdx*tdx+tdy*tdy) > maxR {
					continue
				}
			}
			if d < best {
				best = d
				bx = px
				by = py
			}
		}
	}
	return best, bx, by
}

// pointSegDistXY returns the distance from (px,py) to segment (ax,ay)-(bx,by)
// and the closest point on the segment.
func pointSegDistXY(px, py, ax, ay, bx, by float64) (float64, float64, float64) {
	dx := bx - ax
	dy := by - ay
	len2 := dx*dx + dy*dy
	if len2 == 0 {
		dx2 := px - ax
		dy2 := py - ay
		return math.Sqrt(dx2*dx2 + dy2*dy2), ax, ay
	}
	t := ((px-ax)*dx + (py-ay)*dy) / len2
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	cx := ax + t*dx
	cy := ay + t*dy
	dx2 := px - cx
	dy2 := py - cy
	return math.Sqrt(dx2*dx2 + dy2*dy2), cx, cy
}

// ─────────────────────────────────────────────────────────────────────────────
// MonkeyBuccaneerBehavior
// ─────────────────────────────────────────────────────────────────────────────

type MonkeyBuccaneerBehavior struct {
	engine.DefaultBehavior
	rng        float64
	camoDetect float64
	leadDetect float64
	shotCount  int // counts main shots for grape/bomb/burst cycles
}

func (b *MonkeyBuccaneerBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["select"] = 0.0
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["ability"] = 0.0
	inst.Vars["ability_max"] = 0.0
	inst.Vars["attackrate"] = 1.0
	if _, ok := inst.Vars["invested"]; !ok {
		inst.Vars["invested"] = 0.0
	}

	// Derive tier and tower_code from object name.
	// Linear chain: Monkey_Buccaneer(0) → Grape_Shot(1) → Crows_Nest(2, branch point)
	switch inst.ObjectName {
	case "Monkey_Buccaneer":
		inst.Vars["tier"] = 0.0
		inst.Vars["tower_code"] = 11.00
	case "Grape_Shot":
		inst.Vars["tier"] = 1.0
		inst.Vars["tower_code"] = 11.00
	default:
		// Crows_Nest and all branch descendants: tier >= 2
		if getVar(inst, "tier") < 2 {
			inst.Vars["tier"] = 2.0
		}
		if _, ok := inst.Vars["legacy_object"]; !ok {
			inst.Vars["legacy_object"] = inst.ObjectName
		}
		if getVar(inst, "tower_code") <= 0 {
			if code, ok := effectiveTowerCode(inst.ObjectName); ok {
				inst.Vars["tower_code"] = code
			} else {
				inst.Vars["tower_code"] = 11.20
			}
		}
	}

	form := b.refreshForm(inst, g)
	inst.Alarms[0] = int(math.Round(form.Reload))
	if form.AgentRate > 0 {
		inst.Alarms[1] = form.AgentRate
	}
	if form.AbilityMax > 0 {
		inst.Vars["ability_max"] = form.AbilityMax
		inst.Alarms[2] = 30
	}
	if form.HarpoonRate > 0 {
		inst.Alarms[3] = form.HarpoonRate
	}
}

func (b *MonkeyBuccaneerBehavior) refreshForm(inst *engine.Instance, g *engine.Game) buccaneerForm {
	form := buccaneerFormFor(inst)
	b.rng = form.Range
	b.camoDetect = form.CamoDet
	b.leadDetect = form.LeadDet
	inst.Vars["range"] = form.Range
	if form.AbilityMax > 0 {
		inst.Vars["ability_max"] = form.AbilityMax
	}
	if code, ok := effectiveTowerCode(buccaneerUpgradeName(inst)); ok {
		inst.Vars["tower_code"] = code
	}
	return form
}

func (b *MonkeyBuccaneerBehavior) Step(inst *engine.Instance, g *engine.Game) {
	form := b.refreshForm(inst, g)

	// linear upgrades: Monkey_Buccaneer → Grape_Shot → Crows_Nest
	if rule, ok := legacyLinearUpgradeRuleFor(inst.ObjectName); ok {
		if applyLinearUpgrade(inst, g, rule) {
			return // destroyed & recreated — bail out
		}
	}

	if applyPathUpgrade(inst, g) {
		b.shotCount = 0
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.Reload))
		if form.AgentRate > 0 && inst.Alarms[1] <= 0 {
			inst.Alarms[1] = form.AgentRate
		}
		if form.AbilityMax > 0 && inst.Alarms[2] <= 0 {
			inst.Vars["ability_max"] = form.AbilityMax
			inst.Alarms[2] = 30
		}
		if form.HarpoonRate > 0 && inst.Alarms[3] <= 0 {
			inst.Alarms[3] = form.HarpoonRate
		}
		if form.CursedRate > 0 && inst.Alarms[4] <= 0 {
			inst.Alarms[4] = form.CursedRate
		}
		return
	}

	// face nearest bloon
	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target != nil {
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}
}

func (b *MonkeyBuccaneerBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	form := b.refreshForm(inst, g)

	switch idx {
	case 0:
		// main fire
		b.fireBuccaneerDart(inst, g, form)
		inst.Alarms[0] = int(math.Round(form.Reload))
	case 1:
		// spawn agent
		if form.AgentObj != "" {
			b.spawnAgent(inst, g, form)
			inst.Alarms[1] = form.AgentRate
		}
	case 2:
		// ability charge
		if getGlobal(g, "wavenow") == 1 {
			inst.Vars["ability"] = getVar(inst, "ability") + 1
		}
		inst.Alarms[2] = 30
	case 3:
		// harpoon fire
		if form.HarpoonRate > 0 {
			b.fireHarpoon(inst, g, form)
			inst.Alarms[3] = form.HarpoonRate
		}
	case 4:
		// cursed dart spawn (Ghost_Ship / Cursed_Pirate_Ship)
		if form.CursedRate > 0 {
			b.fireCursedDart(inst, g, form)
			inst.Alarms[4] = form.CursedRate
		}
	}
}

func (b *MonkeyBuccaneerBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	towerClickSelect(inst, g, towerSelectValue(11.0, inst))
}

// ─────────────────────────────────────────────────────────────────────────────
// Fire main dart
// ─────────────────────────────────────────────────────────────────────────────

func (b *MonkeyBuccaneerBehavior) fireBuccaneerDart(inst *engine.Instance, g *engine.Game, form buccaneerForm) {
	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target == nil {
		return
	}
	ppbuff := getVar(inst, "ppbuff")
	b.shotCount++

	// main dart
	dartSpeed := form.DartSpeed
	if dartSpeed <= 0 {
		dartSpeed = 25
	}
	proj := form.Projectile
	if proj == "" {
		proj = "Buccaneer_Dart"
	}
	dart := g.InstanceMgr.Create(proj, inst.X, inst.Y)
	if dart != nil {
		aimProjectile(dart, target, dartSpeed)
		dart.Vars["LP"] = form.LP
		dart.Vars["PP"] = form.PP + ppbuff
		dart.Vars["leadpop"] = b.leadDetect
		dart.Vars["camopop"] = b.camoDetect
		dart.Alarms[0] = 21
		if form.Accelerate {
			dart.Vars["accelerate"] = 1.0
		}
	}

	// burst fire (Supreme_Battleship, Aircraft_Carrier)
	if form.BurstCount > 0 && form.BurstCycle > 0 && b.shotCount%form.BurstCycle == 0 {
		for i := 0; i < form.BurstCount; i++ {
			burstDart := g.InstanceMgr.Create(proj, inst.X, inst.Y)
			if burstDart == nil {
				continue
			}
			// spread the burst darts in a fan pattern
			angle := math.Atan2(-(target.Y - inst.Y), target.X-inst.X)
			spread := (float64(i) - float64(form.BurstCount)/2.0) * 0.15
			bAngle := angle + spread
			burstDart.HSpeed = math.Cos(bAngle) * dartSpeed
			burstDart.VSpeed = -math.Sin(bAngle) * dartSpeed
			burstDart.Direction = bAngle * 180 / math.Pi
			burstDart.ImageAngle = burstDart.Direction
			burstDart.Vars["LP"] = form.LP
			burstDart.Vars["PP"] = form.PP + ppbuff
			burstDart.Vars["leadpop"] = b.leadDetect
			burstDart.Vars["camopop"] = b.camoDetect
			burstDart.Alarms[0] = 21
		}
	}

	// grape volley
	if form.GrapeCycle > 0 && b.shotCount%form.GrapeCycle == 0 {
		b.fireGrapeVolley(inst, g, target, form, ppbuff)
	}

	// ship bomb (Cannon_Ship path)
	if form.BombCycle > 0 && b.shotCount%form.BombCycle == 0 {
		b.fireShipBomb(inst, g, target, form, ppbuff)
	}

	inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
}

// ─────────────────────────────────────────────────────────────────────────────
// Grape volley (4 grapes spread in a fan)
// ─────────────────────────────────────────────────────────────────────────────

func (b *MonkeyBuccaneerBehavior) fireGrapeVolley(
	inst *engine.Instance, g *engine.Game,
	target *engine.Instance, form buccaneerForm, ppbuff float64,
) {
	count := form.GrapeCount
	if count <= 0 {
		count = 4
	}
	baseAngle := math.Atan2(-(target.Y - inst.Y), target.X-inst.X)
	// offsets: +30°, +10°, -10°, -30° (from GML cycle pattern)
	offsets := []float64{30, 10, -10, -30}
	if count > len(offsets) {
		count = len(offsets)
	}
	for i := 0; i < count; i++ {
		grape := g.InstanceMgr.Create("Grape", inst.X, inst.Y)
		if grape == nil {
			continue
		}
		angle := baseAngle + offsets[i]*math.Pi/180
		const grapeSpeed = 23.0
		grape.HSpeed = math.Cos(angle) * grapeSpeed
		grape.VSpeed = -math.Sin(angle) * grapeSpeed
		grape.Direction = angle * 180 / math.Pi
		grape.ImageAngle = grape.Direction
		grape.Vars["LP"] = form.GrapeLP
		grape.Vars["PP"] = form.GrapePP + ppbuff
		grape.Vars["leadpop"] = b.leadDetect
		grape.Vars["camopop"] = b.camoDetect
		grape.Alarms[0] = 20
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Ship bomb (Cannon_Ship path — creates explosion on impact)
// ─────────────────────────────────────────────────────────────────────────────

func (b *MonkeyBuccaneerBehavior) fireShipBomb(
	inst *engine.Instance, g *engine.Game,
	target *engine.Instance, form buccaneerForm, ppbuff float64,
) {
	bomb := g.InstanceMgr.Create("Ship_Bomb", inst.X, inst.Y)
	if bomb == nil {
		return
	}
	aimProjectile(bomb, target, 28)
	bomb.Vars["LP"] = form.BombLP
	bomb.Vars["PP"] = form.BombPP + ppbuff
	bomb.Vars["leadpop"] = 1.0
	bomb.Vars["camopop"] = b.camoDetect
	bomb.Alarms[0] = 30
}

// ─────────────────────────────────────────────────────────────────────────────
// Harpoon attack (Harpoon_Ship / MOAB_Takedown)
// ─────────────────────────────────────────────────────────────────────────────

func (b *MonkeyBuccaneerBehavior) fireHarpoon(inst *engine.Instance, g *engine.Game, form buccaneerForm) {
	target := findStrongestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target == nil {
		return
	}
	ppbuff := getVar(inst, "ppbuff")
	harpoon := g.InstanceMgr.Create("Harpoon", inst.X, inst.Y)
	if harpoon == nil {
		return
	}
	aimProjectile(harpoon, target, 1.0) // slow-moving harpoon
	harpoon.Vars["LP"] = form.HarpoonLP
	harpoon.Vars["PP"] = form.HarpoonPP + ppbuff
	harpoon.Vars["leadpop"] = 1.0
	harpoon.Vars["camopop"] = b.camoDetect
	harpoon.ImageXScale = 1.2
	harpoon.ImageYScale = 1.2
	harpoon.Alarms[0] = 12
}

// ─────────────────────────────────────────────────────────────────────────────
// Cursed dart spawn (Cursed_Pirate_Ship / Ghost_Ship)
// Spawns at random screen position and moves toward tower
// ─────────────────────────────────────────────────────────────────────────────

func (b *MonkeyBuccaneerBehavior) fireCursedDart(inst *engine.Instance, g *engine.Game, form buccaneerForm) {
	// random spawn at screen edges
	spawnX := rand.Float64() * 1024
	spawnY := rand.Float64() * 576
	ppbuff := getVar(inst, "ppbuff")

	projName := "Cursed_Dart"
	if buccaneerUpgradeName(inst) == "Ghost_Ship" {
		projName = "Ghost_Dart"
	}
	dart := g.InstanceMgr.Create(projName, spawnX, spawnY)
	if dart == nil {
		return
	}
	// aim toward tower position
	dx := inst.X - spawnX
	dy := inst.Y - spawnY
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist > 0 {
		const speed = 3.0
		dart.HSpeed = (dx / dist) * speed
		dart.VSpeed = (dy / dist) * speed
		dart.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
	}
	dart.ImageAngle = dart.Direction
	dart.Vars["LP"] = form.CursedDartLP
	dart.Vars["PP"] = form.CursedDartPP + ppbuff
	dart.Vars["leadpop"] = 1.0
	dart.Vars["camopop"] = 1.0
	dart.Vars["accelerate"] = 1.0
	dart.Vars["ownerID"] = float64(inst.ID)
	dart.Alarms[0] = 200 // long lifetime
}

// ─────────────────────────────────────────────────────────────────────────────
// Spawn agents (Swashbuckler, Monkey_Pirate, etc.)
// ─────────────────────────────────────────────────────────────────────────────

func (b *MonkeyBuccaneerBehavior) spawnAgent(inst *engine.Instance, g *engine.Game, form buccaneerForm) {
	if form.AgentObj == "" {
		return
	}
	agent := g.InstanceMgr.Create(form.AgentObj, inst.X, inst.Y)
	if agent == nil {
		return
	}
	// launch in random direction from tower — use HSpeed/VSpeed (engine moves via these)
	dir := rand.Float64() * 360
	rad := dir * math.Pi / 180
	const launchSpeed = 8.0 // moderate launch — don't fly too far from the ship
	agent.HSpeed = math.Cos(rad) * launchSpeed
	agent.VSpeed = -math.Sin(rad) * launchSpeed
	agent.Speed = launchSpeed
	agent.Direction = dir
	agent.Friction = 0.2
	agent.Depth = -1
	// store parent tower position + range so agent stays nearby
	agent.Vars["tower_x"] = inst.X
	agent.Vars["tower_y"] = inst.Y
	agent.Vars["tower_range"] = form.Range
}

// ─────────────────────────────────────────────────────────────────────────────
// Helper: aim projectile toward target
// ─────────────────────────────────────────────────────────────────────────────

func aimProjectile(proj, target *engine.Instance, speed float64) {
	dx := target.X - proj.X
	dy := target.Y - proj.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist > 0 {
		proj.HSpeed = (dx / dist) * speed
		proj.VSpeed = (dy / dist) * speed
		proj.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
	}
	proj.ImageAngle = proj.Direction
}

// ─────────────────────────────────────────────────────────────────────────────
// Projectile behaviors
// ─────────────────────────────────────────────────────────────────────────────

// BuccaneerDartBehavior — standard dart for Monkey Buccaneer
type BuccaneerDartBehavior struct {
	engine.DefaultBehavior
}

func (b *BuccaneerDartBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 1, 0, 0)
}

func (b *BuccaneerDartBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	if getVar(inst, "accelerate") == 1 {
		inst.Speed += 0.3
	}
	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 14)
}

func (b *BuccaneerDartBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// GrapeBehavior — grape shot projectile
type GrapeBehavior struct {
	engine.DefaultBehavior
}

func (b *GrapeBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 1, 0, 0)
}

func (b *GrapeBehavior) Step(inst *engine.Instance, g *engine.Game) {
	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 12)
}

func (b *GrapeBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// ShipBombBehavior — cannonball that explodes on bloon impact
type ShipBombBehavior struct {
	engine.DefaultBehavior
}

func (b *ShipBombBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 50, 1, 0)
}

func (b *ShipBombBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	// check for bloon contact — explode on hit (with swept collision for high speed)
	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	prevX := inst.XPrevious
	prevY := inst.YPrevious
	curX := inst.X
	curY := inst.Y
	const contactRadius = 24.0
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		// current position check
		dx := curX - bloon.X
		dy := curY - bloon.Y
		d := math.Sqrt(dx*dx + dy*dy)
		if d < contactRadius {
			b.explode(inst, g)
			return
		}
		// swept collision: closest point on movement segment
		segDx := curX - prevX
		segDy := curY - prevY
		segLen2 := segDx*segDx + segDy*segDy
		if segLen2 > contactRadius*contactRadius {
			t := ((bloon.X-prevX)*segDx + (bloon.Y-prevY)*segDy) / segLen2
			if t < 0 {
				t = 0
			} else if t > 1 {
				t = 1
			}
			cx := prevX + t*segDx
			cy := prevY + t*segDy
			cdx := cx - bloon.X
			cdy := cy - bloon.Y
			if math.Sqrt(cdx*cdx+cdy*cdy) < contactRadius {
				b.explode(inst, g)
				return
			}
		}
	}
}

func (b *ShipBombBehavior) explode(inst *engine.Instance, g *engine.Game) {
	explosion := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
	if explosion != nil {
		explosion.Vars["LP"] = getVar(inst, "LP")
		explosion.Vars["PP"] = getVar(inst, "PP")
		explosion.Vars["leadpop"] = 1.0
		explosion.Vars["camopop"] = getVar(inst, "camopop")
		explosion.Vars["explode_radius"] = 48.0
		explosion.SpriteName = "Medium_Explosion_Spr"
		explosion.ImageXScale = 1.4
		explosion.ImageYScale = 1.4
		explosion.Alarms[1] = 8
	}
	g.AudioMgr.Play("Small_Boom")
	g.InstanceMgr.Destroy(inst.ID)
}

func (b *ShipBombBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		// timeout — explode at current position as fallback
		b.explode(inst, g)
	}
}

// HarpoonBehavior — slow-moving piercing harpoon
type HarpoonBehavior struct {
	engine.DefaultBehavior
}

func (b *HarpoonBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 1, 1, 0)
	inst.Vars["armourpop"] = 5.0
	inst.Vars["shellpop"] = 5.0
	inst.Vars["shieldpop"] = 5.0
}

func (b *HarpoonBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 18)
}

func (b *HarpoonBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 || idx == 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// DreadnautDartBehavior — accelerating fire dart
type DreadnautDartBehavior struct {
	engine.DefaultBehavior
}

func (b *DreadnautDartBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 1, 1, 1)
	inst.Vars["shieldpop"] = -0.25
}

func (b *DreadnautDartBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	inst.Speed += 0.3 // accelerating
	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 16)
}

func (b *DreadnautDartBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// CursedDartBehavior — ghost/cursed projectile that spawns at random position
type CursedDartBehavior struct {
	engine.DefaultBehavior
}

func (b *CursedDartBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 1, 1, 1)
	inst.Vars["shieldpop"] = -1.0
}

func (b *CursedDartBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	inst.Speed += 0.3 // accelerating
	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 16)
}

func (b *CursedDartBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Agent behaviors (spawned units with independent targeting)
// ─────────────────────────────────────────────────────────────────────────────

// BuccaneerAgentBehavior — generic agent with independent targeting and limited lifespan.
// Used for Swashbuckler, Monkey_Pirate, Dread_Monkey, Cursed_Monkey, Haunted_Monkey.
type BuccaneerAgentBehavior struct {
	engine.DefaultBehavior
	rng        float64
	camoDetect bool
	leadDetect bool
	projName   string // main projectile to fire
	projLP     float64
	projPP     float64
	projSpeed  float64
	fireRate   int     // alarm[0] reload
	lifespan   int     // alarm[1] self-destruct
	projScaleX float64 // optional projectile scale
	projScaleY float64
	spinning   bool // projectile spins (Pirate_Sword)
	// Extra projectile (e.g. Monkey_Pirate fires a sword AND 2 grapes)
	extraProjName  string
	extraProjLP    float64
	extraProjPP    float64
	extraProjSpeed float64
	extraProjCount int     // how many extra projectiles per attack
	extraSpread    float64 // angular spread between extras (degrees)
}

func (b *BuccaneerAgentBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["select"] = 0.0
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["launch_timer"] = 30.0 // fly outward for 30 frames before checking path proximity
	inst.Alarms[0] = b.fireRate
	inst.Alarms[1] = b.lifespan
}

func (b *BuccaneerAgentBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// Clamp to playfield bounds (don't leave the map)
	if inst.X < 8 {
		inst.X = 8
		inst.HSpeed = 0
	}
	if inst.X > 856 {
		inst.X = 856
		inst.HSpeed = 0
	}
	if inst.Y < 8 {
		inst.Y = 8
		inst.VSpeed = 0
	}
	if inst.Y > 472 {
		inst.Y = 472
		inst.VSpeed = 0
	}

	// Launch phase: let the agent fly outward from the tower before seeking the path.
	lt := getVar(inst, "launch_timer")
	if lt > 0 {
		inst.Vars["launch_timer"] = lt - 1
	} else {
		// Path-seeking: find nearest point on the bloon track within tower range.
		tx := getVar(inst, "tower_x")
		ty := getVar(inst, "tower_y")
		tr := getVar(inst, "tower_range")
		if tr <= 0 {
			tr = 140
		}
		maxR := tr + 50 // search radius slightly beyond tower range
		d, nx, ny := nearestTrackPoint(inst.X, inst.Y, g, tx, ty, maxR)

		if d < 15 {
			// Close enough — park on the track.
			inst.Speed = 0
			inst.HSpeed = 0
			inst.VSpeed = 0
			inst.Friction = 0
		} else if d < 60 {
			// Approaching — steer toward nearest track point.
			dx := nx - inst.X
			dy := ny - inst.Y
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist > 0 {
				spd := inst.Speed
				if spd < 1.5 {
					spd = 1.5
				}
				if spd > 3 {
					spd = 3
				}
				inst.HSpeed = (dx / dist) * spd
				inst.VSpeed = (dy / dist) * spd
				inst.Speed = spd
				inst.Friction = 0.3
			}
		}
		// Beyond 60px from any nearby path: keep drifting with natural friction
	}

	// face nearest bloon
	target := findNearestBloon(inst, g, b.rng, b.camoDetect)
	if target != nil {
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}
}

func (b *BuccaneerAgentBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	if idx != 0 {
		return
	}
	// fire projectile
	target := findNearestBloon(inst, g, b.rng, b.camoDetect)
	if target == nil {
		inst.Alarms[0] = b.fireRate
		return
	}
	ppbuff := getVar(inst, "ppbuff")
	leadVal := 0.0
	if b.leadDetect {
		leadVal = 1.0
	}
	camoVal := 0.0
	if b.camoDetect {
		camoVal = 1.0
	}

	proj := g.InstanceMgr.Create(b.projName, inst.X, inst.Y)
	if proj != nil {
		aimProjectile(proj, target, b.projSpeed)
		proj.Vars["LP"] = b.projLP
		proj.Vars["PP"] = b.projPP + ppbuff
		proj.Vars["leadpop"] = leadVal
		proj.Vars["camopop"] = camoVal
		if b.projScaleX != 0 {
			proj.ImageXScale = b.projScaleX
		}
		if b.projScaleY != 0 {
			proj.ImageYScale = b.projScaleY
		}
		if b.spinning {
			// Sword/melee: offset direction +96° for arc sweep, short lifespan.
			// Must update HSpeed/VSpeed too (engine moves via those).
			newDir := proj.Direction + 96
			rad := newDir * math.Pi / 180
			proj.HSpeed = math.Cos(rad) * b.projSpeed
			proj.VSpeed = -math.Sin(rad) * b.projSpeed
			proj.Direction = newDir
			proj.ImageAngle = newDir
			proj.Alarms[0] = 5
		} else {
			proj.Alarms[0] = 20
		}
	}

	// Extra projectiles (e.g. Monkey_Pirate fires 2 Grapes in addition to Cutlass)
	if b.extraProjName != "" && b.extraProjCount > 0 {
		baseAngle := math.Atan2(-(target.Y - inst.Y), target.X-inst.X)
		spread := b.extraSpread * math.Pi / 180
		for i := 0; i < b.extraProjCount; i++ {
			ep := g.InstanceMgr.Create(b.extraProjName, inst.X, inst.Y)
			if ep == nil {
				continue
			}
			// spread symmetrically: -spread/2, +spread/2 for 2 projectiles
			offset := (float64(i) - float64(b.extraProjCount-1)/2.0) * spread
			angle := baseAngle + offset
			ep.HSpeed = math.Cos(angle) * b.extraProjSpeed
			ep.VSpeed = -math.Sin(angle) * b.extraProjSpeed
			ep.Direction = angle * 180 / math.Pi
			ep.ImageAngle = ep.Direction
			ep.Vars["LP"] = b.extraProjLP
			ep.Vars["PP"] = b.extraProjPP + ppbuff
			ep.Vars["leadpop"] = leadVal
			ep.Vars["camopop"] = camoVal
			ep.Alarms[0] = 21
		}
	}

	inst.Alarms[0] = b.fireRate
	inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
}

// PirateSwordBehavior — short-range melee strike from Swashbuckler / Monkey_Pirate.
// In the original GML the sword gets direction += 96 (offset arc) and lives ~5 frames.
type PirateSwordBehavior struct {
	engine.DefaultBehavior
}

func (b *PirateSwordBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 20, 0, 0)
	inst.Vars["armourpop"] = 2.0
	inst.Vars["shellpop"] = 2.0
	inst.Vars["shieldpop"] = 2.0
}

func (b *PirateSwordBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// The sword doesn't spin continuously — it keeps its arc direction and
	// ImageAngle follows Direction (set once at creation with +96° offset).
	inst.ImageAngle = inst.Direction
	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 16)
}

func (b *PirateSwordBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 || idx == 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// BoatPlaneBehavior — Aircraft_Carrier's flying plane unit.
// Orbits in a circle firing radial bursts and directional shots.
type BoatPlaneBehavior struct {
	engine.DefaultBehavior
}

func (b *BoatPlaneBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["select"] = 0.0
	inst.Vars["ppbuff"] = 0.0
	// Start with a random direction, constant speed 2 (circular flight).
	dir := rand.Float64() * 360
	inst.Direction = dir
	inst.Speed = 2
	inst.Friction = 0 // planes don't decelerate
	// Set HSpeed/VSpeed from direction (engine needs these for actual movement)
	rad := dir * math.Pi / 180
	inst.HSpeed = math.Cos(rad) * 2
	inst.VSpeed = -math.Sin(rad) * 2
	inst.Alarms[0] = 8   // radial burst fire
	inst.Alarms[2] = 2   // directional fire
	inst.Alarms[3] = 790 // lifespan
}

func (b *BoatPlaneBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// Override friction/speed in case spawnAgent set them
	inst.Friction = 0
	// circular motion: rotate direction slightly each frame
	inst.Direction += 0.9
	inst.Speed = 2
	rad := inst.Direction * math.Pi / 180
	inst.HSpeed = math.Cos(rad) * 2
	inst.VSpeed = -math.Sin(rad) * 2
	inst.ImageAngle = inst.Direction
}

func (b *BoatPlaneBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	ppbuff := getVar(inst, "ppbuff")
	switch idx {
	case 0:
		// radial burst: 8 darts at 45° intervals
		for i := 0; i < 8; i++ {
			angle := float64(i) * 45 * math.Pi / 180
			dart := g.InstanceMgr.Create("Buccaneer_Dart", inst.X, inst.Y)
			if dart == nil {
				continue
			}
			const speed = 27.0
			dart.HSpeed = math.Cos(angle) * speed
			dart.VSpeed = -math.Sin(angle) * speed
			dart.Direction = angle * 180 / math.Pi
			dart.ImageAngle = dart.Direction
			dart.Vars["LP"] = 1.0
			dart.Vars["PP"] = 8.0 + ppbuff
			dart.Vars["leadpop"] = 0.0
			dart.Vars["camopop"] = 1.0
			dart.Alarms[0] = 40
		}
		inst.Alarms[0] = 8
	case 2:
		// directional fire: 2 darts from heading
		for _, off := range []float64{-4, 4} {
			dart := g.InstanceMgr.Create("Buccaneer_Dart", inst.X, inst.Y)
			if dart == nil {
				continue
			}
			angle := (inst.Direction + off) * math.Pi / 180
			const speed = 30.0
			dart.HSpeed = math.Cos(angle) * speed
			dart.VSpeed = -math.Sin(angle) * speed
			dart.Direction = (inst.Direction + off)
			dart.ImageAngle = dart.Direction
			dart.Vars["LP"] = 1.0
			dart.Vars["PP"] = 6.0 + ppbuff
			dart.Vars["leadpop"] = 0.0
			dart.Vars["camopop"] = 1.0
			dart.Alarms[0] = 40
		}
		inst.Alarms[2] = 2
	case 3:
		// self-destruct
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Registration
// ─────────────────────────────────────────────────────────────────────────────

func RegisterBuccaneerBehaviors(im *engine.InstanceManager) {
	// all tower forms use the same behavior
	buccaneerObjects := []string{
		"Monkey_Buccaneer", "Grape_Shot", "Crows_Nest",
		"Swashbucklers", "Monkey_Pirates", "Pirate_Captain_Ship",
		"Destroyer", "Supreme_Battleship", "Aircraft_Carrier",
		"Cannon_Ship", "Harpoon_Ship", "MOAB_Takedown",
		"Dreadnaut_Ship", "Cursed_Pirate_Ship", "Ghost_Ship",
	}
	for _, n := range buccaneerObjects {
		im.RegisterBehavior(n, func() engine.InstanceBehavior { return &MonkeyBuccaneerBehavior{} })
	}

	// projectiles
	im.RegisterBehavior("Buccaneer_Dart", func() engine.InstanceBehavior { return &BuccaneerDartBehavior{} })
	im.RegisterBehavior("Grape", func() engine.InstanceBehavior { return &GrapeBehavior{} })
	im.RegisterBehavior("Ship_Bomb", func() engine.InstanceBehavior { return &ShipBombBehavior{} })
	im.RegisterBehavior("Harpoon", func() engine.InstanceBehavior { return &HarpoonBehavior{} })
	im.RegisterBehavior("Dreadnaut_Dart", func() engine.InstanceBehavior { return &DreadnautDartBehavior{} })
	im.RegisterBehavior("Cursed_Dart", func() engine.InstanceBehavior { return &CursedDartBehavior{} })
	im.RegisterBehavior("Ghost_Dart", func() engine.InstanceBehavior { return &CursedDartBehavior{} })
	im.RegisterBehavior("Pirate_Sword", func() engine.InstanceBehavior { return &PirateSwordBehavior{} })

	// agents
	im.RegisterBehavior("Swashbuckler", func() engine.InstanceBehavior {
		return &BuccaneerAgentBehavior{
			rng: 100, camoDetect: true, projName: "Pirate_Sword",
			projLP: 1, projPP: 20, projSpeed: 4, fireRate: 40, lifespan: 800,
			spinning: true, projScaleX: 1.5, projScaleY: 1.5,
		}
	})
	im.RegisterBehavior("Monkey_Pirate", func() engine.InstanceBehavior {
		return &BuccaneerAgentBehavior{
			rng: 120, camoDetect: true, projName: "Pirate_Sword",
			projLP: 1, projPP: 30, projSpeed: 4, fireRate: 27, lifespan: 800,
			spinning: true, projScaleX: 1.5, projScaleY: 1.5,
			extraProjName: "Grape", extraProjLP: 1, extraProjPP: 2,
			extraProjSpeed: 25, extraProjCount: 2, extraSpread: 14,
		}
	})
	im.RegisterBehavior("Monkey_Pirate_Captain", func() engine.InstanceBehavior {
		return &BuccaneerAgentBehavior{
			rng: 135, camoDetect: true, projName: "Pirate_Sword",
			projLP: 1, projPP: 30, projSpeed: 4, fireRate: 7, lifespan: 800,
			spinning: true, projScaleX: 1.5, projScaleY: 1.5,
			extraProjName: "Grape", extraProjLP: 1, extraProjPP: 2,
			extraProjSpeed: 25, extraProjCount: 2, extraSpread: 14,
		}
	})
	im.RegisterBehavior("Dread_Monkey", func() engine.InstanceBehavior {
		return &BuccaneerAgentBehavior{
			rng: 100, camoDetect: true, leadDetect: true, projName: "Dreadnaut_Dart",
			projLP: 1, projPP: 10, projSpeed: 3, fireRate: 20, lifespan: 800,
			projScaleX: 0.65, projScaleY: 0.7,
		}
	})
	im.RegisterBehavior("Cursed_Monkey", func() engine.InstanceBehavior {
		return &BuccaneerAgentBehavior{
			rng: 120, camoDetect: true, leadDetect: true, projName: "Dreadnaut_Dart",
			projLP: 1, projPP: 10, projSpeed: 3, fireRate: 8, lifespan: 800,
			projScaleX: 0.7, projScaleY: 0.75,
		}
	})
	im.RegisterBehavior("Haunted_Monkey", func() engine.InstanceBehavior {
		return &BuccaneerAgentBehavior{
			rng: 120, camoDetect: true, leadDetect: true, projName: "Dreadnaut_Dart",
			projLP: 2, projPP: 10, projSpeed: 2, fireRate: 5, lifespan: 800,
			projScaleX: 0.8, projScaleY: 0.85,
		}
	})

	// boat plane (Aircraft_Carrier)
	im.RegisterBehavior("Boat_Plane", func() engine.InstanceBehavior { return &BoatPlaneBehavior{} })
}
