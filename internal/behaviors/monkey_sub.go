package behaviors

import (
	"math"
	"math/rand"

	"btdx/internal/engine"
)

// ──────────────────────────────────────────────────────────────────────────────
// Monkey Sub — water-only tower that fires homing darts.
// Upgrade chain: Barbed_Darts_Sub → Twin_Guns → (path 1) Airburst_Sub
//                                               (path 2) Support_Sub
//                                               (path 3) Torpedo_Sub
// ──────────────────────────────────────────────────────────────────────────────

type monkeySubForm struct {
	Range               float64
	Reload              float64 // alarm[0] reset after firing (frames)
	LP                  float64
	PP                  float64
	Camo                float64 // 1 = detects camo bloons
	Lead                float64 // 1 = pops lead bloons
	TorpoCycle          int     // >0: every TorpoCycle shots also fires a Torpedo/Missile
	TorpoProjectile     string  // projectile name for torpedo cycle (default "Torpedo")
	TorpoLP             float64 // LP for torpedo/missile projectile (default 1)
	TorpoPP             float64 // PP for torpedo/missile projectile (default 20)
	TorpoSpeed          float64 // speed for torpedo/missile (default 21→27 for Torpedo, 30 for Ballistic)
	Projectile          string  // object name of the regular dart projectile to spawn
	BurstReload         int     // >0: interval (frames) for Sub_Energy decamo burst (alarm[1])
	HasBloontoniumPulse bool    // true: also fire Bloontonium_Energy on burst (reactor forms)
	NoFire              bool    // true: skip regular dart fire (submerge-style forms)
	AbilityMax          float64 // >0: cooldown charge cap for first-strike style abilities
	StunDecrement       float64 // per-alarm stun reduction (default 11)
}

var monkeySubForms = map[string]monkeySubForm{
	// linear chain ──────────────────────────────────────────────────────────────
	"Monkey_Sub":       {Range: 115, Reload: 14, LP: 1, PP: 1, Projectile: "Monkey_Sub_Dart"},
	"Barbed_Darts_Sub": {Range: 115, Reload: 14, LP: 1, PP: 3, Projectile: "Barbed_Dart"},
	"Twin_Guns":        {Range: 115, Reload: 7, LP: 1, PP: 3, Projectile: "Barbed_Dart"},
	// path 1 (left) ─────────────────────────────────────────────────────────────
	"Airburst_Sub": {Range: 115, Reload: 7, LP: 1, PP: 3, Projectile: "Airburst_Dart"},
	// path 2 (middle) ────────────────────────────────────────────────────────────
	"Support_Sub": {Range: 115, Reload: 7, LP: 1, PP: 3, Camo: 1, NoFire: true, BurstReload: 50},
	// support_Sub tier-3 upgrades
	"Bloontonium_Reactor": {Range: 115, Reload: 7, LP: 1, PP: 3, Camo: 1, Lead: 1, Projectile: "Barbed_Dart", BurstReload: 43, HasBloontoniumPulse: true},
	"Anti_Matter_Reactor": {Range: 115, Reload: 7, LP: 1, PP: 3, Camo: 1, Lead: 1, Projectile: "Barbed_Dart", BurstReload: 43, HasBloontoniumPulse: true},
	// path 3 (right) ─────────────────────────────────────────────────────────────
	"Torpedo_Sub":           {Range: 125, Reload: 7, LP: 1, PP: 3, TorpoCycle: 3, TorpoProjectile: "Torpedo", TorpoLP: 1, TorpoPP: 20, TorpoSpeed: 27, Projectile: "Barbed_Dart", StunDecrement: 11},
	"Ballistic_Missile_Sub": {Range: 135, Reload: 7, LP: 1, PP: 3, TorpoCycle: 3, TorpoProjectile: "Ballistic_Missile", TorpoLP: 2, TorpoPP: 60, TorpoSpeed: 30, Projectile: "Barbed_Dart", StunDecrement: 35},
	"First_Strike_Sub":      {Range: 135, Reload: 7, LP: 1, PP: 3, TorpoCycle: 2, TorpoProjectile: "Ballistic_Missile", TorpoLP: 2, TorpoPP: 60, TorpoSpeed: 30, Projectile: "Barbed_Dart", AbilityMax: 42, StunDecrement: 88},
	// path 1 continued (left) — Airburst path tier 2 & 3
	"Assault_Wave_Sub": {Range: 125, Reload: 6, LP: 1, PP: 4, Projectile: "Airwave_Dart", StunDecrement: 28},
	"Blockade_Sub":     {Range: 135, Reload: 5, LP: 1, PP: 6, Projectile: "Airwave_Dart", AbilityMax: 42},
	// special path
	"Smart_Sub": {Range: 170, Reload: 7, LP: 1, PP: 3, Camo: 1, Lead: 1, Projectile: "Barbed_Dart"},
}

func monkeySubUpgradeName(inst *engine.Instance) string {
	return upgradeName(inst, "Monkey_Sub")
}

func monkeySubFormFor(inst *engine.Instance) monkeySubForm {
	if f, ok := monkeySubForms[monkeySubUpgradeName(inst)]; ok {
		return f
	}
	return monkeySubForms["Monkey_Sub"]
}

type MonkeySubBehavior struct {
	engine.DefaultBehavior
	rng                 float64
	camoDetect          float64
	leadDetect          float64
	cycle               int  // torpedo cycle counter
	burstReload         int  // alarm[1] period for Sub_Energy decamo burst (0=none)
	hasBloontoniumPulse bool // also fire Bloontonium_Energy on burst
}

func (b *MonkeySubBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["ability"] = 0.0
	inst.Vars["ability_max"] = 0.0
	if _, ok := inst.Vars["invested"]; !ok {
		inst.Vars["invested"] = 0.0
	}
	form := b.refreshForm(inst, g)
	inst.Alarms[0] = 16
	if b.burstReload > 0 {
		inst.Alarms[1] = b.burstReload
	}
	if form.AbilityMax > 0 {
		inst.Vars["ability_max"] = form.AbilityMax
		inst.Alarms[2] = 30
	}
}

func (b *MonkeySubBehavior) refreshForm(inst *engine.Instance, g *engine.Game) monkeySubForm {
	form := monkeySubFormFor(inst)
	b.rng = form.Range
	b.camoDetect = form.Camo
	b.leadDetect = form.Lead
	b.burstReload = form.BurstReload
	b.hasBloontoniumPulse = form.HasBloontoniumPulse
	// stop dart fire alarm when form switches to a NoFire form (submerge path)
	if form.NoFire && inst.Alarms[0] > 0 {
		inst.Alarms[0] = 0
	}
	inst.Vars["range"] = form.Range
	if form.AbilityMax > 0 {
		inst.Vars["ability_max"] = form.AbilityMax
	}
	if code, ok := effectiveTowerCode(monkeySubUpgradeName(inst)); ok {
		inst.Vars["tower_code"] = code
	}
	if spr := g.InstanceMgr.ObjectSpriteName(monkeySubUpgradeName(inst)); spr != "" && g.AssetManager.GetSprite(spr) != nil {
		inst.SpriteName = spr
	}
	return form
}

func (b *MonkeySubBehavior) fireTorpedo(inst *engine.Instance, g *engine.Game, target *engine.Instance, ppbuff float64, form monkeySubForm) {
	projName := form.TorpoProjectile
	if projName == "" {
		projName = "Torpedo"
	}
	torp := g.InstanceMgr.Create(projName, inst.X, inst.Y)
	if torp == nil {
		return
	}
	torpSpeed := form.TorpoSpeed
	if torpSpeed <= 0 {
		torpSpeed = 21.0
	}
	dx := target.X - inst.X
	dy := target.Y - inst.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist > 0 {
		torp.HSpeed = (dx / dist) * torpSpeed
		torp.VSpeed = (dy / dist) * torpSpeed
		torp.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
	}
	torp.ImageAngle = torp.Direction
	torpLP := form.TorpoLP
	if torpLP <= 0 {
		torpLP = 1.0
	}
	torpPP := form.TorpoPP
	if torpPP <= 0 {
		torpPP = 20.0
	}
	torp.Vars["LP"] = torpLP
	torp.Vars["PP"] = torpPP + ppbuff
	torp.Vars["leadpop"] = 1.0
	torp.Vars["camopop"] = b.camoDetect
	torp.Vars["targetID"] = float64(target.ID)
	torp.Alarms[1] = 30
}

func (b *MonkeySubBehavior) fireSubDart(inst *engine.Instance, g *engine.Game, form monkeySubForm) {
	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target == nil {
		return
	}

	ppbuff := getVar(inst, "ppbuff")

	// Torpedo_Sub / Ballistic_Missile_Sub / First_Strike_Sub:
	// every TorpoCycle shots, also fire a Torpedo/Ballistic_Missile
	if form.TorpoCycle > 0 {
		b.cycle++
		if b.cycle >= form.TorpoCycle {
			b.cycle = 0
			b.fireTorpedo(inst, g, target, ppbuff, form)
		}
	}

	proj := form.Projectile
	if proj == "" {
		proj = "Monkey_Sub_Dart"
	}
	dart := g.InstanceMgr.Create(proj, inst.X, inst.Y)
	if dart == nil {
		return
	}
	dx := target.X - inst.X
	dy := target.Y - inst.Y
	dist := math.Sqrt(dx*dx + dy*dy)
	if dist > 0 {
		const speed = 21.0
		dart.HSpeed = (dx / dist) * speed
		dart.VSpeed = (dy / dist) * speed
		dart.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
	}
	dart.ImageAngle = dart.Direction
	dart.Vars["LP"] = form.LP
	dart.Vars["PP"] = form.PP + ppbuff
	dart.Vars["leadpop"] = b.leadDetect
	dart.Vars["camopop"] = b.camoDetect
	dart.Vars["targetID"] = float64(target.ID)
	dart.Alarms[0] = 20

	inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
}

func (b *MonkeySubBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		// Sub_Energy decamo burst (Support_Sub / Bloontonium_Reactor)
		b.fireBurst(inst, g)
		if b.burstReload > 0 {
			inst.Alarms[1] = b.burstReload
		}
		return
	}
	if idx == 2 {
		// ability charge (First_Strike_Sub / Blockade_Sub)
		if getGlobal(g, "wavenow") == 1 {
			inst.Vars["ability"] = getVar(inst, "ability") + 1
			// charge faster when bloons are present
			bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
			if len(bloons) > 0 {
				inst.Vars["ability"] = getVar(inst, "ability") + 1
			}
		}
		inst.Alarms[2] = 30
		return
	}
	if idx != 0 {
		return
	}
	form := b.refreshForm(inst, g)
	// submerge-style forms: no dart fire
	if form.NoFire {
		return
	}
	if stunVal := getVar(inst, "stun"); stunVal > 0 {
		stunDec := form.StunDecrement
		if stunDec <= 0 {
			stunDec = 11
		}
		inst.Vars["stun"] = stunVal - stunDec
		if getVar(inst, "stun") < 0 {
			inst.Vars["stun"] = 0.0
		}
		inst.Alarms[0] = int(math.Round(form.Reload))
		return
	}
	b.fireSubDart(inst, g, form)
	inst.Alarms[0] = int(math.Round(form.Reload))
}

// fireBurst fires the Sub_Energy (and optionally Bloontonium_Energy) decamo pulse.
func (b *MonkeySubBehavior) fireBurst(inst *engine.Instance, g *engine.Game) {
	// sub_Energy: LP=0, PP=200 → decamos camo bloons (no lead pop)
	se := g.InstanceMgr.Create("Sub_Energy", inst.X, inst.Y)
	if se != nil {
		se.Vars["LP"] = 0.0
		se.Vars["PP"] = 200.0
		se.Vars["leadpop"] = b.leadDetect
		se.Vars["camopop"] = b.camoDetect
	}
	if b.hasBloontoniumPulse {
		// bloontonium_Energy: LP=1, PP=200 → additional energy burst
		be := g.InstanceMgr.Create("Bloontonium_Energy", inst.X, inst.Y)
		if be != nil {
			be.Vars["LP"] = 1.0
			be.Vars["PP"] = 200.0
			be.Vars["leadpop"] = b.leadDetect
			be.Vars["camopop"] = b.camoDetect
		}
	}
}

func (b *MonkeySubBehavior) Step(inst *engine.Instance, g *engine.Game) {
	form := b.refreshForm(inst, g)
	if applyPathUpgrade(inst, g) {
		b.cycle = 0
		form = b.refreshForm(inst, g)
		if !form.NoFire {
			inst.Alarms[0] = int(math.Round(form.Reload))
		} else {
			inst.Alarms[0] = 0 // stop dart fire
		}
		// restart burst timer if new form has one
		if b.burstReload > 0 && inst.Alarms[1] <= 0 {
			inst.Alarms[1] = b.burstReload
		}
		// start ability charge timer if new form has an ability
		if form.AbilityMax > 0 && inst.Alarms[2] <= 0 {
			inst.Vars["ability_max"] = form.AbilityMax
			inst.Alarms[2] = 30
		}
		return
	}
	if applyTowerUpgrade(inst, g) > 0 {
		b.cycle = 0
		form = b.refreshForm(inst, g)
		if !form.NoFire {
			inst.Alarms[0] = int(math.Round(form.Reload))
		} else {
			inst.Alarms[0] = 0
		}
		if b.burstReload > 0 && inst.Alarms[1] <= 0 {
			inst.Alarms[1] = b.burstReload
		}
		if form.AbilityMax > 0 && inst.Alarms[2] <= 0 {
			inst.Vars["ability_max"] = form.AbilityMax
			inst.Alarms[2] = 30
		}
		return
	}
	// submerged forms: don't rotate the sprite toward targets
	if form.NoFire {
		return
	}
	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target != nil {
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}
}

func (b *MonkeySubBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	// First Strike / Blockade ability activation
	if activateSubAbility(inst, g, b.camoDetect, b.leadDetect) {
		return
	}
	towerClickSelect(inst, g, towerSelectValue(7.0, inst))
}

// activateSubAbility — First_Strike_Sub's massive missile ability.
func activateSubAbility(inst *engine.Instance, g *engine.Game, camoDetect, leadDetect float64) bool {
	abilityMax := getVar(inst, "ability_max")
	if abilityMax <= 0 {
		return false
	}
	if getVar(inst, "ability") < abilityMax {
		return false
	}
	inst.Vars["ability"] = 0.0
	// fire a First Strike Missile at the strongest/furthest bloon
	target := findNearestBloon(inst, g, 99999, camoDetect == 1)
	if target == nil {
		return true
	}
	missile := g.InstanceMgr.Create("First_Strike_Missile", inst.X, inst.Y)
	if missile != nil {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > 0 {
			missile.HSpeed = (dx / dist) * 55
			missile.VSpeed = (dy / dist) * 55
			missile.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
		missile.ImageAngle = missile.Direction
		missile.Vars["LP"] = 5000.0
		missile.Vars["PP"] = 1.0 + getVar(inst, "ppbuff")
		missile.Vars["leadpop"] = leadDetect
		missile.Vars["camopop"] = camoDetect
		missile.Vars["targetID"] = float64(target.ID)
		missile.Alarms[1] = 40
	}
	g.AudioMgr.Play("Upgrade")
	return true
}

// ──────────────────────────────────────────────────────────────────────────────
// Sub projectile behaviors
// ──────────────────────────────────────────────────────────────────────────────

// SubHomingDartBehavior — homing dart for Monkey Sub (Monkey_Sub_Dart, Barbed_Dart).
// Homes to targetID at speed 21; self-destructs on alarm[0] or when PP <= 0.
type SubHomingDartBehavior struct {
	engine.DefaultBehavior
	hitRadius float64
}

func (b *SubHomingDartBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 1, 0, 0)
	if _, ok := inst.Vars["targetID"]; !ok {
		inst.Vars["targetID"] = 0.0
	}
}

func (b *SubHomingDartBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	const speed = 21.0
	targetID := int(getVar(inst, "targetID"))
	target := g.InstanceMgr.GetByID(targetID)
	if target == nil || target.Destroyed {
		target = findNearestBloon(inst, g, 99999, getVar(inst, "camopop") == 1)
		if target != nil {
			inst.Vars["targetID"] = float64(target.ID)
		}
	}
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
	hr := b.hitRadius
	if hr <= 0 {
		hr = 12
	}
	projectileHitBloons(inst, g, hr)
}

func (b *SubHomingDartBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// AirburstDartBehavior — Airburst_Sub's dart: homes like a SubHomingDart but
// detonates into a small AoE explosion when PP runs out.
type AirburstDartBehavior struct {
	SubHomingDartBehavior
}

func (b *AirburstDartBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.SubHomingDartBehavior.Create(inst, g)
}

func (b *AirburstDartBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		expl := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
		if expl != nil {
			expl.Vars["LP"] = 1.0
			expl.Vars["PP"] = 3.0
			expl.Vars["leadpop"] = getVar(inst, "leadpop")
			expl.Vars["camopop"] = getVar(inst, "camopop")
			expl.Vars["explode_radius"] = 30.0
			expl.Alarms[1] = 8
		}
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	b.SubHomingDartBehavior.Step(inst, g)
}

func (b *AirburstDartBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	b.SubHomingDartBehavior.Alarm(inst, idx, g)
}

// TorpedoProjectileBehavior — homing torpedo fired by Torpedo_Sub every 3rd shot.
// Explodes (Small_Explosion) on contact with a bloon or after alarm[1]=30 frames.
type TorpedoProjectileBehavior struct {
	engine.DefaultBehavior
}

func (b *TorpedoProjectileBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 20, 1, 0)
	if _, ok := inst.Vars["targetID"]; !ok {
		inst.Vars["targetID"] = 0.0
	}
}

func (b *TorpedoProjectileBehavior) torpedoExplode(inst *engine.Instance, g *engine.Game) {
	expl := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
	if expl != nil {
		expl.Vars["LP"] = getVar(inst, "LP")
		expl.Vars["PP"] = getVar(inst, "PP")
		expl.Vars["leadpop"] = getVar(inst, "leadpop")
		expl.Vars["camopop"] = getVar(inst, "camopop")
		expl.Vars["explode_radius"] = 40.0
		expl.Alarms[1] = 8
	}
	g.InstanceMgr.Destroy(inst.ID)
}

func (b *TorpedoProjectileBehavior) Step(inst *engine.Instance, g *engine.Game) {
	const speed = 21.0
	targetID := int(getVar(inst, "targetID"))
	target := g.InstanceMgr.GetByID(targetID)
	if target == nil || target.Destroyed {
		target = findNearestBloon(inst, g, 99999, getVar(inst, "camopop") == 1)
		if target != nil {
			inst.Vars["targetID"] = float64(target.ID)
		}
	}
	if target != nil && !target.Destroyed {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= speed {
			inst.X = target.X
			inst.Y = target.Y
			b.torpedoExplode(inst, g)
			return
		} else if dist > 0 {
			inst.HSpeed = (dx / dist) * speed
			inst.VSpeed = (dy / dist) * speed
			inst.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
	}
	inst.ImageAngle = inst.Direction
	// contact detonation at radius 14
	for _, bloon := range g.InstanceMgr.FindByObject("Normal_Bloon_Branch") {
		if bloon.Destroyed {
			continue
		}
		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		if math.Sqrt(dx*dx+dy*dy) < 14 {
			b.torpedoExplode(inst, g)
			return
		}
	}
}

func (b *TorpedoProjectileBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		b.torpedoExplode(inst, g)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Ballistic_Missile — homing missile fired by Ballistic_Missile_Sub / First_Strike_Sub.
// Similar to Torpedo but detonates into a Moab_Explosion (large AoE, armourpop/shellpop).
// ──────────────────────────────────────────────────────────────────────────────

type BallisticMissileBehavior struct {
	engine.DefaultBehavior
}

func (b *BallisticMissileBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 2, 60, 1, 0)
	if _, ok := inst.Vars["targetID"]; !ok {
		inst.Vars["targetID"] = 0.0
	}
}

func (b *BallisticMissileBehavior) ballisticExplode(inst *engine.Instance, g *engine.Game) {
	// Creates Moab_Explosion — large AoE with armourpop/shellpop
	expl := g.InstanceMgr.Create("Moab_Explosion", inst.X, inst.Y)
	if expl != nil {
		expl.Vars["LP"] = getVar(inst, "LP")
		expl.Vars["PP"] = getVar(inst, "PP")
		expl.Vars["leadpop"] = getVar(inst, "leadpop")
		expl.Vars["camopop"] = getVar(inst, "camopop")
		expl.Vars["explode_radius"] = 50.0
		expl.Vars["armourpop"] = 8.0
		expl.Vars["shellpop"] = 8.0
		expl.ImageXScale = 1.4
		expl.ImageYScale = 1.4
		expl.Alarms[1] = 8
	}
	g.AudioMgr.Play("Large_Boom")
	g.InstanceMgr.Destroy(inst.ID)
}

func (b *BallisticMissileBehavior) Step(inst *engine.Instance, g *engine.Game) {
	const speed = 21.0
	targetID := int(getVar(inst, "targetID"))
	target := g.InstanceMgr.GetByID(targetID)
	if target == nil || target.Destroyed {
		target = findNearestBloon(inst, g, 99999, getVar(inst, "camopop") == 1)
		if target != nil {
			inst.Vars["targetID"] = float64(target.ID)
		}
	}
	if target != nil && !target.Destroyed {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= speed {
			inst.X = target.X
			inst.Y = target.Y
			b.ballisticExplode(inst, g)
			return
		} else if dist > 0 {
			inst.HSpeed = (dx / dist) * speed
			inst.VSpeed = (dy / dist) * speed
			inst.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
	}
	inst.ImageAngle = inst.Direction
	// contact detonation at radius 14
	for _, bloon := range g.InstanceMgr.FindByObject("Normal_Bloon_Branch") {
		if bloon.Destroyed {
			continue
		}
		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		if math.Sqrt(dx*dx+dy*dy) < 14 {
			b.ballisticExplode(inst, g)
			return
		}
	}
}

func (b *BallisticMissileBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		b.ballisticExplode(inst, g)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// First_Strike_Missile — massive LP missile from First_Strike_Sub's ability.
// Speed 55, LP=5000 — devastating against MOABs.
// ──────────────────────────────────────────────────────────────────────────────

type FirstStrikeMissileBehavior struct {
	engine.DefaultBehavior
}

func (b *FirstStrikeMissileBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 5000, 1, 1, 0)
	if _, ok := inst.Vars["targetID"]; !ok {
		inst.Vars["targetID"] = 0.0
	}
	inst.ImageXScale = 2.0
	inst.ImageYScale = 2.0
}

func (b *FirstStrikeMissileBehavior) firstStrikeExplode(inst *engine.Instance, g *engine.Game) {
	// massive explosion
	expl := g.InstanceMgr.Create("Moab_Explosion", inst.X, inst.Y)
	if expl != nil {
		expl.Vars["LP"] = 1000.0
		expl.Vars["PP"] = 75.0
		expl.Vars["leadpop"] = getVar(inst, "leadpop")
		expl.Vars["camopop"] = getVar(inst, "camopop")
		expl.Vars["explode_radius"] = 80.0
		expl.Vars["armourpop"] = 20.0
		expl.Vars["shellpop"] = 20.0
		expl.ImageXScale = 2.5
		expl.ImageYScale = 2.5
		expl.Alarms[1] = 10
	}
	g.AudioMgr.Play("Large_Boom")
	g.InstanceMgr.Destroy(inst.ID)
}

func (b *FirstStrikeMissileBehavior) Step(inst *engine.Instance, g *engine.Game) {
	const speed = 55.0
	targetID := int(getVar(inst, "targetID"))
	target := g.InstanceMgr.GetByID(targetID)
	if target == nil || target.Destroyed {
		target = findNearestBloon(inst, g, 99999, getVar(inst, "camopop") == 1)
		if target != nil {
			inst.Vars["targetID"] = float64(target.ID)
		}
	}
	if target != nil && !target.Destroyed {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist <= speed {
			inst.X = target.X
			inst.Y = target.Y
			b.firstStrikeExplode(inst, g)
			return
		} else if dist > 0 {
			inst.HSpeed = (dx / dist) * speed
			inst.VSpeed = (dy / dist) * speed
			inst.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
	}
	inst.ImageAngle = inst.Direction
	// contact detonation
	for _, bloon := range g.InstanceMgr.FindByObject("Normal_Bloon_Branch") {
		if bloon.Destroyed {
			continue
		}
		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		if math.Sqrt(dx*dx+dy*dy) < 20 {
			b.firstStrikeExplode(inst, g)
			return
		}
	}
}

func (b *FirstStrikeMissileBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		b.firstStrikeExplode(inst, g)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Moab_Explosion — large AoE explosion from Ballistic_Missile and First_Strike_Missile.
// Same as SmallExplosionBehavior but supports armourpop/shellpop vars for MOAB stripping.
// ──────────────────────────────────────────────────────────────────────────────

type MoabExplosionBehavior struct {
	engine.DefaultBehavior
}

func (b *MoabExplosionBehavior) Create(inst *engine.Instance, g *engine.Game) {
	setProjDefaults(inst, 2, 60, 1, 0)
	inst.ImageSpeed = 1.0
	inst.ImageAngle = rand.Float64() * 360
}

func (b *MoabExplosionBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	radius := getVar(inst, "explode_radius")
	if radius <= 0 {
		radius = 50
	}
	projectileHitBloons(inst, g, radius)
}

func (b *MoabExplosionBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 || idx == 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Airwave_Dart — homing dart for Assault_Wave_Sub and Blockade_Sub.
// Similar to AirburstDart but with a slightly larger burst on depletion.
// ──────────────────────────────────────────────────────────────────────────────

type AirwaveDartBehavior struct {
	SubHomingDartBehavior
}

func (b *AirwaveDartBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.SubHomingDartBehavior.Create(inst, g)
}

func (b *AirwaveDartBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		// burst into a small AoE on depletion (like airburst but stronger)
		expl := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
		if expl != nil {
			expl.Vars["LP"] = 1.0
			expl.Vars["PP"] = 4.0
			expl.Vars["leadpop"] = getVar(inst, "leadpop")
			expl.Vars["camopop"] = getVar(inst, "camopop")
			expl.Vars["explode_radius"] = 35.0
			expl.Alarms[1] = 8
		}
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	b.SubHomingDartBehavior.Step(inst, g)
}

func (b *AirwaveDartBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	b.SubHomingDartBehavior.Alarm(inst, idx, g)
}

// SubEnergyBehavior — expanding energy ring fired by Support_Sub (and Bloontonium_Reactor).
// Decamos all camo bloons it touches during its 21-frame lifetime.
// Mirrors GML: alarm[0]=21; step: image_xscale/yscale += 0.04; collision → camo=0.
type SubEnergyBehavior struct {
	engine.DefaultBehavior
	radius float64 // current hit radius (grows each step)
}

func (b *SubEnergyBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.radius = 40 // starting radius (px); sprite origin ~center
	inst.Alarms[0] = 21
	inst.ImageXScale = 1.2
	inst.ImageYScale = 1.2
	// slight random rotation so overlapping rings look distinct
	inst.ImageAngle = math.Mod(inst.X+inst.Y, 360)
}

func (b *SubEnergyBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageXScale += 0.04
	inst.ImageYScale += 0.04
	b.radius += 2.0 // expand hit radius each step

	// decamo all camo bloons within the growing ring
	camopop := getVar(inst, "camopop")
	if camopop != 1 {
		return
	}
	for _, bloon := range g.InstanceMgr.FindByObject("Normal_Bloon_Branch") {
		if bloon.Destroyed {
			continue
		}
		if getVar(bloon, "camo") != 1 {
			continue
		}
		dx := bloon.X - inst.X
		dy := bloon.Y - inst.Y
		if math.Sqrt(dx*dx+dy*dy) <= b.radius {
			bloon.Vars["camo"] = 0.0
		}
	}
}

func (b *SubEnergyBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}
