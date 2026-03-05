package behaviors

import (
	"math"
	"math/rand"

	"btdx/internal/engine"
)

// sniper Monkey — global range, furthest-bloon targeting, homing darts.

type sniperTowerForm struct {
	Range      float64
	AlarmBase  float64 // alarm[0] fire delay
	Reload     float64 // alarm[11] re-target delay
	LP         float64
	PP         float64
	Lead       float64
	Camo       float64
	Speed      float64 // projectile speed
	Stun       float64 // stun decrement per shot
	Projectile string  // projectile object name
	AbilityMax float64 // 0 = no ability
}

var sniperTowerForms = map[string]sniperTowerForm{
	// base chain (tier 0→1→2)
	"Sniper_Monkey":     {Range: 2000, AlarmBase: 72, Reload: 62, LP: 3, PP: 1, Lead: 0, Camo: 0, Speed: 30, Stun: 3, Projectile: "Sniper_Dart"},
	"Full_Metal_Jacket": {Range: 2000, AlarmBase: 72, Reload: 61, LP: 5, PP: 1, Lead: 1, Camo: 0, Speed: 30, Stun: 5, Projectile: "Sniper_Dart"},
	"Heat_Sniper":       {Range: 2000, AlarmBase: 60, Reload: 47, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 7, Projectile: "Sniper_Dart"},
	// path 1 (left): Tactical Shotgun → Bloonzooka → RPG Strike
	"Tactical_Shotgun": {Range: 2000, AlarmBase: 61, Reload: 43, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 39, Projectile: "Shotgun_Slug"},
	"Bloonzooka":       {Range: 2000, AlarmBase: 61, Reload: 43, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 247, Projectile: "Bloonzooka_Shot"},
	"RPG_Strike":       {Range: 2000, AlarmBase: 61, Reload: 40, LP: 10, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 300, Projectile: "Bloonzooka_Shot", AbilityMax: 43},
	// path 2 (middle): Deadly Precision → Brick Layer → Moab Crippler
	"Deadly_Precision": {Range: 2000, AlarmBase: 57, Reload: 45, LP: 18, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 18, Projectile: "Sniper_Dart"},
	"Brick_Layer":      {Range: 2000, AlarmBase: 54, Reload: 42, LP: 48, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 48, Projectile: "Sniper_Dart"},
	"Moab_Crippler":    {Range: 2000, AlarmBase: 54, Reload: 39, LP: 75, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 75, Projectile: "Sniper_Dart"},
	// path 3 (right): Semi Auto → Machine Gun → Supply Drones
	"Semi_Automatic_Rifle": {Range: 2000, AlarmBase: 18, Reload: 14, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 54, Stun: 7, Projectile: "Sniper_Dart"},
	"Machine_Gun":          {Range: 2000, AlarmBase: 8, Reload: 7, LP: 8, PP: 1, Lead: 1, Camo: 1, Speed: 54, Stun: 8, Projectile: "Sniper_Dart"},
	"Supply_Drones":        {Range: 2000, AlarmBase: 8, Reload: 7, LP: 8, PP: 1, Lead: 1, Camo: 1, Speed: 54, Stun: 8, Projectile: "Sniper_Dart", AbilityMax: 48},
	// special path
	"Shotgun_Plus":    {Range: 2000, AlarmBase: 61, Reload: 43, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 39, Projectile: "Shotgun_Slug"},
	"Bloonzooka_Plus": {Range: 2000, AlarmBase: 61, Reload: 43, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 247, Projectile: "Bloonzooka_Shot"},
	"Railgun_Tank":    {Range: 2000, AlarmBase: 54, Reload: 36, LP: 100, PP: 1, Lead: 1, Camo: 1, Speed: 54, Stun: 500, Projectile: "Sniper_Dart"},
}

func sniperUpgradeName(inst *engine.Instance) string {
	return upgradeName(inst, "Sniper_Monkey")
}

func sniperFormFor(inst *engine.Instance) sniperTowerForm {
	if form, ok := sniperTowerForms[sniperUpgradeName(inst)]; ok {
		return form
	}
	return sniperTowerForms["Sniper_Monkey"]
}

type SniperMonkeyBehavior struct {
	engine.DefaultBehavior
	rng        float64
	camoDetect float64
	leadDetect float64
}

func (b *SniperMonkeyBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.rng = 2000.0
	b.camoDetect = 0.0
	b.leadDetect = 0.0
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["range"] = b.rng
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["ability"] = 0.0
	inst.Vars["ability_max"] = 0.0
	form := b.refreshForm(inst, g)
	inst.Alarms[0] = int(math.Round(form.AlarmBase))
}

func (b *SniperMonkeyBehavior) refreshForm(inst *engine.Instance, g *engine.Game) sniperTowerForm {
	form := sniperFormFor(inst)
	b.rng = form.Range
	b.camoDetect = form.Camo
	b.leadDetect = form.Lead
	inst.Vars["range"] = form.Range
	if code, ok := effectiveTowerCode(sniperUpgradeName(inst)); ok {
		inst.Vars["tower_code"] = code
	}
	if spr := g.InstanceMgr.ObjectSpriteName(sniperUpgradeName(inst)); spr != "" && g.AssetManager.GetSprite(spr) != nil {
		inst.SpriteName = spr
	}
	if form.AbilityMax > 0 {
		inst.Vars["ability_max"] = form.AbilityMax
	} else {
		inst.Vars["ability_max"] = 0.0
	}
	return form
}

func (b *SniperMonkeyBehavior) fireSniperDart(inst *engine.Instance, g *engine.Game, form sniperTowerForm) {
	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target == nil {
		return
	}
	dart := g.InstanceMgr.Create(form.Projectile, inst.X, inst.Y)
	if dart == nil {
		return
	}
	dx := target.X - inst.X
	dy := target.Y - inst.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist > 0 {
		dart.HSpeed = (dx / dist) * form.Speed
		dart.VSpeed = (dy / dist) * form.Speed
		dart.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
	}
	dart.ImageAngle = dart.Direction
	dart.Vars["LP"] = form.LP
	dart.Vars["PP"] = form.PP + getVar(inst, "ppbuff")
	dart.Vars["leadpop"] = form.Lead
	dart.Vars["camopop"] = form.Camo
	dart.Vars["targetID"] = float64(target.ID)
	dart.Alarms[1] = 40

	inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
}

func (b *SniperMonkeyBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	form := b.refreshForm(inst, g)
	switch idx {
	case 0:
		// fire cycle: shoot then wait for alarm[11] retarget loop.
		stunVal := getVar(inst, "stun")
		if stunVal > 0 {
			inst.Vars["stun"] = stunVal - form.Stun
			if getVar(inst, "stun") < 0 {
				inst.Vars["stun"] = 0.0
			}
			inst.Alarms[0] = int(math.Round(form.AlarmBase))
			return
		}
		b.fireSniperDart(inst, g, form)
		inst.Alarms[11] = int(math.Round(form.Reload)) - 1
	case 2:
		// ability charge loop (Supply_Drones): +1 every 30 frames during a wave.
		inst.Alarms[2] = 30
		if getVar(inst, "stun") > 0 {
			return
		}
		wavenow := getGlobal(g, "wavenow")
		hasBloons := len(g.InstanceMgr.FindByObject("Normal_Bloon_Branch")) > 0
		if wavenow == 1 || hasBloons {
			ability := getVar(inst, "ability") + 1
			max := getVar(inst, "ability_max")
			if max > 0 && ability > max {
				ability = max
			}
			inst.Vars["ability"] = ability
		}
	case 11:
		// retarget idle loop: find a target and schedule immediate fire.
		target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
		if target != nil {
			inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
			inst.Alarms[0] = 1
		} else {
			inst.Alarms[11] = 1
		}
	}
}

func (b *SniperMonkeyBehavior) activateSniperAbility(inst *engine.Instance, g *engine.Game) bool {
	max := getVar(inst, "ability_max")
	if max <= 0 || getVar(inst, "ability") < max {
		return false
	}
	formName := sniperUpgradeName(inst)
	if formName == "RPG_Strike" {
		proj := g.InstanceMgr.Create("RPG_Projectile", inst.X, inst.Y)
		if proj != nil {
			proj.Vars["LP"] = 400.0
			proj.Vars["PP"] = 120.0
			proj.Vars["leadpop"] = 1.0
			proj.Vars["camopop"] = 1.0
			proj.ImageXScale = 2.0
			proj.ImageYScale = 2.0
			proj.Alarms[0] = 30
		}
		inst.Vars["ability"] = 0.0
		g.AudioMgr.Play("Upgrade")
		return true
	}
	// supply Drop: spawn 10 Supply_Drone instances across the screen.
	positions := [][2]float64{
		{0, 48}, {0, 144}, {0, 240}, {0, 336}, {0, 432},
		{-48, 96}, {-48, 192}, {-48, 288}, {-48, 384}, {-48, 480},
	}
	for _, pos := range positions {
		g.InstanceMgr.Create("Supply_Drone", pos[0], pos[1])
	}
	inst.Vars["ability"] = 0.0
	g.AudioMgr.Play("Upgrade")
	return true
}

func (b *SniperMonkeyBehavior) Step(inst *engine.Instance, g *engine.Game) {
	form := b.refreshForm(inst, g)
	if applyPathUpgrade(inst, g) {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase))
		// start ability charge loop if gained an ability.
		if form.AbilityMax > 0 && inst.Alarms[2] <= 0 {
			inst.Alarms[2] = 30
		}
		return
	}
	if applyTowerUpgrade(inst, g) > 0 {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase))
	}

	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target != nil {
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}
}

func (b *SniperMonkeyBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if b.activateSniperAbility(inst, g) {
		return
	}
	towerClickSelect(inst, g, towerSelectValue(4.0, inst))
}

// RPGProjectileBehavior — the RPG Strike ability projectile.
// Moves at speed 33 toward the nearest bloon; on bloon hit or after alarm[0]=30 steps,
// detonates a massive Medium_Explosion (LP=400, PP=120, explode_radius=64).
type RPGProjectileBehavior struct {
	engine.DefaultBehavior
}

func (b *RPGProjectileBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 400, 120, 1, 1)
	if inst.Alarms[0] <= 0 {
		inst.Alarms[0] = 30
	}
}

func (b *RPGProjectileBehavior) rpgExplode(inst *engine.Instance, g *engine.Game) {
	expl := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
	if expl != nil {
		expl.Vars["LP"] = getVar(inst, "LP")
		expl.Vars["PP"] = getVar(inst, "PP")
		expl.Vars["leadpop"] = 1.0
		expl.Vars["camopop"] = 1.0
		expl.Vars["explode_radius"] = 64.0
		expl.ImageXScale = 2.5
		expl.ImageYScale = 2.5
		expl.Alarms[0] = 8
	}
	g.InstanceMgr.Destroy(inst.ID)
}

func (b *RPGProjectileBehavior) Step(inst *engine.Instance, g *engine.Game) {
	const speed = 33.0
	camopop := getVar(inst, "camopop") == 1
	target := findNearestBloon(inst, g, 99999, camopop)
	if target != nil {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= speed {
			inst.X = target.X
			inst.Y = target.Y
			b.rpgExplode(inst, g)
			return
		} else if dist > 0 {
			inst.HSpeed = (dx / dist) * speed
			inst.VSpeed = (dy / dist) * speed
			inst.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
	}
	inst.ImageAngle = inst.Direction

	// Collision check: detonate on bloon contact
	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		if math.Sqrt(dx*dx+dy*dy) < 20 {
			b.rpgExplode(inst, g)
			return
		}
	}
}

func (b *RPGProjectileBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		b.rpgExplode(inst, g)
	}
}

// ShotgunSlugBehavior — homing sniper slug that spawns 5 Frag shrapnel on impact.
type ShotgunSlugBehavior struct {
	engine.DefaultBehavior
}

func (b *ShotgunSlugBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 7, 1, 1, 1)
	if _, ok := inst.Vars["targetID"]; !ok {
		inst.Vars["targetID"] = 0.0
	}
}

func (b *ShotgunSlugBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	targetID := int(getVar(inst, "targetID"))
	target := g.InstanceMgr.GetByID(targetID)
	if target == nil || target.Destroyed {
		target = findNearestBloon(inst, g, 2000, getVar(inst, "camopop") == 1)
		if target != nil {
			inst.Vars["targetID"] = float64(target.ID)
		}
	}
	const speed = 30.0
	if target != nil && !target.Destroyed {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= speed {
			inst.X = target.X
			inst.Y = target.Y
			inst.HSpeed = 0
			inst.VSpeed = 0
		} else if dist > 0 {
			inst.HSpeed = (dx / dist) * speed
			inst.VSpeed = (dy / dist) * speed
			inst.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
	}
	inst.ImageAngle = inst.Direction

	// Collision: on hitting the target, pop it and spray shrapnel
	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	lp := getVar(inst, "LP")
	leadpop := getVar(inst, "leadpop")
	camopop := getVar(inst, "camopop")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		if getVar(bloon, "lead") == 1 && leadpop == 0 {
			continue
		}
		if getVar(bloon, "camo") == 1 && camopop == 0 {
			continue
		}
		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		if math.Sqrt(dx*dx+dy*dy) >= 20 {
			continue
		}
		// spawn 5 Frag shrapnel in a pentagonal spread
		baseAngle := rand.Float64() * 72.0
		for i := 0; i < 5; i++ {
			dirDeg := baseAngle + 72.0*float64(i)
			rad := dirDeg * math.Pi / 180.0
			frag := g.InstanceMgr.Create("Frag", inst.X, inst.Y)
			if frag != nil {
				frag.Vars["LP"] = 4.0
				frag.Vars["PP"] = 1.0
				frag.Vars["leadpop"] = 1.0
				frag.Vars["camopop"] = 1.0
				frag.HSpeed = math.Cos(rad) * 24.0
				frag.VSpeed = -math.Sin(rad) * 24.0
				frag.Direction = dirDeg
				frag.ImageAngle = dirDeg
				frag.Alarms[0] = 5
			}
		}
		popBloon(bloon, lp, g)
		inst.Vars["PP"] = getVar(inst, "PP") - 1
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
}

func (b *ShotgunSlugBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// BloonzookaShotBehavior — homing slug that spawns a Medium_Explosion AoE on impact.
type BloonzookaShotBehavior struct {
	engine.DefaultBehavior
}

func (b *BloonzookaShotBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 7, 1, 1, 1)
	if _, ok := inst.Vars["targetID"]; !ok {
		inst.Vars["targetID"] = 0.0
	}
}

func (b *BloonzookaShotBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	targetID := int(getVar(inst, "targetID"))
	target := g.InstanceMgr.GetByID(targetID)
	if target == nil || target.Destroyed {
		target = findNearestBloon(inst, g, 2000, getVar(inst, "camopop") == 1)
		if target != nil {
			inst.Vars["targetID"] = float64(target.ID)
		}
	}
	const speed = 30.0
	if target != nil && !target.Destroyed {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= speed {
			inst.X = target.X
			inst.Y = target.Y
			inst.HSpeed = 0
			inst.VSpeed = 0
		} else if dist > 0 {
			inst.HSpeed = (dx / dist) * speed
			inst.VSpeed = (dy / dist) * speed
			inst.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
	}
	inst.ImageAngle = inst.Direction

	// Collision: on hitting any bloon, pop it and spawn explosion AoE
	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	lp := getVar(inst, "LP")
	leadpop := getVar(inst, "leadpop")
	camopop := getVar(inst, "camopop")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		if getVar(bloon, "lead") == 1 && leadpop == 0 {
			continue
		}
		if getVar(bloon, "camo") == 1 && camopop == 0 {
			continue
		}
		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		if math.Sqrt(dx*dx+dy*dy) >= 20 {
			continue
		}
		popBloon(bloon, lp, g)
		// spawn medium explosion AoE at impact point
		expl := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
		if expl != nil {
			expl.Vars["LP"] = 4.0
			expl.Vars["PP"] = 60.0
			expl.Vars["leadpop"] = 1.0
			expl.Vars["camopop"] = 1.0
			expl.Vars["explode_radius"] = 36.0
			expl.ImageXScale = 1.4
			expl.ImageYScale = 1.4
			expl.Alarms[0] = 8
		}
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
}

func (b *BloonzookaShotBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// FragBehavior — shrapnel shard fired by Shotgun_Slug on impact.
type FragBehavior struct {
	engine.DefaultBehavior
}

func (b *FragBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 4, 1, 1, 1)
}

func (b *FragBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 8)
}

func (b *FragBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// SniperDartBehavior — homing projectile (speed 54, retargets if target dies)
type SniperDartBehavior struct {
	engine.DefaultBehavior
}

func (b *SniperDartBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 3, 1, 0, 0)
	if _, ok := inst.Vars["targetID"]; !ok {
		inst.Vars["targetID"] = 0.0
	}
}

func (b *SniperDartBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	// home toward target at speed 54.
	targetID := int(getVar(inst, "targetID"))
	target := g.InstanceMgr.GetByID(targetID)
	if target == nil || target.Destroyed {
		// retarget furthest bloon if original target died.
		target = findNearestBloon(inst, g, 2000, getVar(inst, "camopop") == 1)
		if target != nil {
			inst.Vars["targetID"] = float64(target.ID)
		}
	}
	if target != nil && !target.Destroyed {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		speed := math.Max(54.0, getVar(inst, "speed"))
		if dist <= speed {
			// snap right onto the target to prevent overshooting
			inst.X = target.X
			inst.Y = target.Y
			inst.HSpeed = 0
			inst.VSpeed = 0
		} else if dist > 0 {
			inst.HSpeed = (dx / dist) * speed
			inst.VSpeed = (dy / dist) * speed
			inst.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
	}
	inst.ImageAngle = inst.Direction
	projectileHitBloons(inst, g, 30)
}

func (b *SniperDartBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}
