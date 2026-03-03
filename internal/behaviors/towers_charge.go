package behaviors

import (
	"math"

	"btdx/internal/engine"
)

// ---------- Charge Tower ----------
//
// Mechanic: two concurrent alarm loops
//   alarm[11] fires every AlarmRate frames — increments charge, decrements stun
//   alarm[0]  fires every FireDelay frames  — fires one projectile if charge > 0 and not stunned
//
// When charge accumulates to ChargeThresh the next alarm[0] tick fires a burst.
// Each shot decrements charge by 1, so towers with large ChargeThresh fire rapid bursts.

type chargeTowerForm struct {
	Range        float64
	ChargeThresh float64
	AlarmRate    float64 // alarm[11] interval
	FireDelay    float64 // alarm[0] interval
	StunDec      float64 // stun reduction per alarm[11] tick
	Lead         float64
	Camo         float64
	Projectile   string
	ProjPP       float64
	ProjLP       float64
	ProjSpeed    float64
	AbilityMax   float64
	RotSpeed     float64 // degrees per step for tower sprite rotation
}

var chargeTowerForms = map[string]chargeTowerForm{
	// ── base chain ──────────────────────────────────────────────────
	"Charge_Tower": {
		Range: 110, ChargeThresh: 8, AlarmRate: 23, FireDelay: 3,
		StunDec: 2, Lead: 0, Camo: 0,
		Projectile: "Charge_Proj", ProjPP: 2, ProjLP: 1, ProjSpeed: 20,
		RotSpeed: 3,
	},
	"Charge_Storage": {
		Range: 110, ChargeThresh: 16, AlarmRate: 15, FireDelay: 2,
		StunDec: 2, Lead: 0, Camo: 0,
		Projectile: "Charge_Proj", ProjPP: 2, ProjLP: 1, ProjSpeed: 20,
		RotSpeed: 3,
	},
	"Powerful_Charges": {
		Range: 112, ChargeThresh: 16, AlarmRate: 15, FireDelay: 2,
		StunDec: 4, Lead: 1, Camo: 0,
		Projectile: "Powerful_Charge", ProjPP: 4, ProjLP: 1, ProjSpeed: 24,
		RotSpeed: 3,
	},
	// ── path 1 (left): Battery ──────────────────────────────────────
	"Charge_Battery": {
		Range: 112, ChargeThresh: 40, AlarmRate: 7, FireDelay: 1,
		StunDec: 4, Lead: 1, Camo: 0,
		Projectile: "Powerful_Charge", ProjPP: 4, ProjLP: 1, ProjSpeed: 24,
	},
	"Charge_Burst": {
		Range: 112, ChargeThresh: 40, AlarmRate: 7, FireDelay: 1,
		StunDec: 10, Lead: 1, Camo: 0,
		Projectile: "Burst_Charge", ProjPP: 4, ProjLP: 1, ProjSpeed: 26,
	},
	"Charge_Overload": {
		Range: 115, ChargeThresh: 40, AlarmRate: 6, FireDelay: 1,
		StunDec: 10, Lead: 1, Camo: 0,
		Projectile: "Burst_Charge", ProjPP: 4, ProjLP: 1, ProjSpeed: 26,
		AbilityMax: 33,
	},
	// ── path 2 (middle): Orbital ────────────────────────────────────
	"Orbital_Discharge": {
		Range: 112, ChargeThresh: 24, AlarmRate: 14, FireDelay: 2,
		StunDec: 10, Lead: 1, Camo: 1,
		Projectile: "Orbital_Charge", ProjPP: 10, ProjLP: 1, ProjSpeed: 29,
		RotSpeed: 3,
	},
	"Magnetic_Charge_Tower": {
		Range: 115, ChargeThresh: 30, AlarmRate: 13, FireDelay: 2,
		StunDec: 26, Lead: 1, Camo: 1,
		Projectile: "Magnetic_Charge", ProjPP: 13, ProjLP: 2, ProjSpeed: 29,
		RotSpeed: 1,
	},
	// ── path 3 (right): Tesla ───────────────────────────────────────
	"Tesla_Coil": {
		Range: 112, ChargeThresh: 24, AlarmRate: 12, FireDelay: 2,
		StunDec: 30, Lead: 1, Camo: 0,
		Projectile: "Small_Energy", ProjPP: 30, ProjLP: 1, ProjSpeed: 32,
	},
	"Giga_Pops": {
		Range: 112, ChargeThresh: 60, AlarmRate: 7, FireDelay: 1,
		StunDec: 70, Lead: 1, Camo: 1,
		Projectile: "Big_Energy", ProjPP: 70, ProjLP: 1, ProjSpeed: 36,
	},
	"Lightning_Bomb": {
		Range: 112, ChargeThresh: 60, AlarmRate: 7, FireDelay: 1,
		StunDec: 70, Lead: 1, Camo: 1,
		Projectile: "Big_Energy", ProjPP: 70, ProjLP: 1, ProjSpeed: 36,
		AbilityMax: 35,
	},
	// ── special path (0): Super / Mega ──────────────────────────────
	"Super_Charge_Tower": {
		Range: 120, ChargeThresh: 8, AlarmRate: 22, FireDelay: 2,
		StunDec: 30, Lead: 1, Camo: 0,
		Projectile: "Super_Charge", ProjPP: 10, ProjLP: 3, ProjSpeed: 25,
		RotSpeed: 2,
	},
	"Mega_Charger": {
		Range: 120, ChargeThresh: 8, AlarmRate: 30, FireDelay: 2,
		StunDec: 320, Lead: 1, Camo: 0,
		Projectile: "Mega_Charge", ProjPP: 20, ProjLP: 4, ProjSpeed: 11,
		RotSpeed: 2,
	},
	"Mega_Mega_Charger": {
		Range: 120, ChargeThresh: 8, AlarmRate: 30, FireDelay: 2,
		StunDec: 320, Lead: 1, Camo: 0,
		Projectile: "Mega_Mega_Charge", ProjPP: 100, ProjLP: 10, ProjSpeed: 11,
		RotSpeed: 2,
	},
}

func chargeUpgradeName(inst *engine.Instance) string {
	return upgradeName(inst, "Charge_Tower")
}

func chargeFormFor(inst *engine.Instance) chargeTowerForm {
	if form, ok := chargeTowerForms[chargeUpgradeName(inst)]; ok {
		return form
	}
	return chargeTowerForms["Charge_Tower"]
}

// ChargeTowerBehavior handles all Charge Tower upgrade forms.
type ChargeTowerBehavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
	camoDetect float64
	leadDetect float64
}

// chargeImageIndex returns the sprite frame index for the current charge level.
// Tesla/Giga forms use charge/2 (large thresh, fewer frames); all others use charge directly.
func chargeImageIndex(form chargeTowerForm, charge float64) float64 {
	switch form.Projectile {
	case "Small_Energy", "Big_Energy":
		return math.Floor(charge / 2)
	}
	return charge
}

func (b *ChargeTowerBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["charge"] = 0.0
	inst.Vars["ability"] = 0.0
	inst.Vars["ability_max"] = 0.0
	inst.Vars["cycle"] = 1.0 // for orbital forms
	inst.ImageSpeed = 0      // freeze animation; frame = charge level
	inst.ImageIndex = 0
	form := b.refreshForm(inst, g)
	// Original sets alarm[0] = 16/attackrate on create (delayed first fire check
	// so initial charge can accumulate before any shot is attempted).
	// alarm[0] reschedules itself to FireDelay after the first tick.
	inst.Alarms[11] = int(math.Round(form.AlarmRate / b.attackRate))
	inst.Alarms[0] = int(math.Round(16 / b.attackRate))
}

func (b *ChargeTowerBehavior) refreshForm(inst *engine.Instance, g *engine.Game) chargeTowerForm {
	form := chargeFormFor(inst)
	b.rng = form.Range
	b.camoDetect = form.Camo
	b.leadDetect = form.Lead
	inst.Vars["range"] = form.Range
	inst.ImageSpeed = 0 // always frozen; frame = charge level
	if code, ok := effectiveTowerCode(chargeUpgradeName(inst)); ok {
		inst.Vars["tower_code"] = code
	}
	if spr := g.InstanceMgr.ObjectSpriteName(chargeUpgradeName(inst)); spr != "" && g.AssetManager.GetSprite(spr) != nil {
		inst.SpriteName = spr
	}
	if form.AbilityMax > 0 {
		inst.Vars["ability_max"] = form.AbilityMax
	} else {
		inst.Vars["ability_max"] = 0.0
	}
	return form
}

func (b *ChargeTowerBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	switch idx {
	case 11:
		// heartbeat: accumulate charge, handle stun
		form := chargeFormFor(inst)
		inst.Alarms[11] = int(math.Round(form.AlarmRate / b.attackRate))

		stun := getVar(inst, "stun")
		if stun > 0 {
			stun -= form.StunDec
			if stun < 0 {
				stun = 0
			}
			inst.Vars["stun"] = stun
			return
		}
		charge := getVar(inst, "charge")
		if charge < form.ChargeThresh {
			charge++
		}
		if charge > form.ChargeThresh {
			charge = form.ChargeThresh
		}
		inst.Vars["charge"] = charge
		// Update sprite frame to reflect charge level (visual dots building up)
		inst.ImageIndex = chargeImageIndex(form, charge)

	case 0:
		// fire: if charged, launch one projectile
		form := chargeFormFor(inst)
		inst.Alarms[0] = int(math.Round(form.FireDelay))

		if getVar(inst, "stun") > 0 {
			return
		}
		charge := getVar(inst, "charge")
		if charge <= 0 {
			return
		}

		target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
		if target == nil {
			return
		}

		b.fireProjectile(inst, g, form, target)
			newCharge := charge - 1
			inst.Vars["charge"] = newCharge
			// Update sprite frame to show one fewer charge dot
			inst.ImageIndex = chargeImageIndex(form, newCharge)

	case 2:
		// ability charge tick (Charge_Overload, Lightning_Bomb)
		if getGlobal(g, "play") == 1 {
			abilityMax := getVar(inst, "ability_max")
			if abilityMax > 0 {
				ab := getVar(inst, "ability") + 1
				if ab > abilityMax {
					ab = abilityMax
				}
				inst.Vars["ability"] = ab
			}
		}
		inst.Alarms[2] = 30
	}
}

// fireProjectile creates the appropriate projectile for the current form.
func (b *ChargeTowerBehavior) fireProjectile(inst *engine.Instance, g *engine.Game, form chargeTowerForm, target *engine.Instance) {
	dx := target.X - inst.X
	dy := target.Y - inst.Y

	// base angle toward target (used for tower sprite rotation)
	baseAngle := math.Atan2(-dy, dx)
	inst.ImageAngle = baseAngle * 180 / math.Pi

	// Orbital/Magnetic projectiles follow spiral GMX paths relative to the tower.
	// They do NOT aim at the target — they expand outward in a square spiral.
	switch form.Projectile {
	case "Orbital_Charge", "Magnetic_Charge":
		cycle := int(getVar(inst, "cycle"))
		next := cycle + 1
		if next > 4 {
			next = 1
		}
		inst.Vars["cycle"] = float64(next)

		proj := g.InstanceMgr.Create(form.Projectile, inst.X, inst.Y)
		if proj == nil {
			return
		}
		// Path-following: no velocity; Step will set X/Y from path each frame.
		proj.HSpeed = 0
		proj.VSpeed = 0
		proj.Vars["pathCycle"] = float64(cycle)
		proj.Vars["pathDist"] = 0.0
		proj.Vars["anchorX"] = inst.X
		proj.Vars["anchorY"] = inst.Y
		proj.Vars["pathSpeed"] = form.ProjSpeed // 29
		proj.Vars["LP"] = form.ProjLP
		proj.Vars["PP"] = form.ProjPP + getVar(inst, "ppbuff")
		proj.Vars["leadpop"] = form.Lead
		proj.Vars["camopop"] = form.Camo
		// Lifetime matches original: Orbital=100 frames, Magnetic=132 frames.
		if form.Projectile == "Orbital_Charge" {
			proj.Alarms[0] = 100
		} else {
			proj.Alarms[0] = 132
		}
		return
	}

	// All other projectiles: aim at target.
	proj := g.InstanceMgr.Create(form.Projectile, inst.X, inst.Y)
	if proj == nil {
		return
	}
	proj.HSpeed = math.Cos(baseAngle) * form.ProjSpeed
	proj.VSpeed = -math.Sin(baseAngle) * form.ProjSpeed
	proj.Direction = baseAngle * 180 / math.Pi
	proj.ImageAngle = proj.Direction
	proj.Vars["LP"] = form.ProjLP
	proj.Vars["PP"] = form.ProjPP + getVar(inst, "ppbuff")
	proj.Vars["leadpop"] = form.Lead
	proj.Vars["camopop"] = form.Camo

	// AoE projectiles (Super_Charge, Mega_Charge, Mega_Mega_Charge) use explode_radius
	switch form.Projectile {
	case "Super_Charge":
		proj.Vars["explode_radius"] = 20.0
		proj.Alarms[1] = 30 // lifetime before exploding
	case "Mega_Charge":
		proj.Vars["explode_radius"] = 30.0
		proj.Alarms[1] = 40
	case "Mega_Mega_Charge":
		proj.Vars["explode_radius"] = 40.0
		proj.ImageXScale = 1.5
		proj.ImageYScale = 1.5
		proj.Alarms[1] = 50
	}
}

func (b *ChargeTowerBehavior) Step(inst *engine.Instance, g *engine.Game) {
	form := b.refreshForm(inst, g)

	// handle sprite rotation
	if form.RotSpeed > 0 {
		inst.ImageAngle += form.RotSpeed
	}

	// path upgrade (branch point → Charge_Battery / Orbital_Discharge / Tesla_Coil / Super_Charge_Tower)
	if applyPathUpgrade(inst, g) {
		form = b.refreshForm(inst, g)
		// reset charge on upgrade (mirrors original instance_change Create event)
		inst.Vars["charge"] = 0.0
		inst.ImageIndex = 0
		inst.Alarms[11] = int(math.Round(form.AlarmRate / b.attackRate))
		inst.Alarms[0] = int(math.Round(16 / b.attackRate))
		if form.AbilityMax > 0 && inst.Alarms[2] <= 0 {
			inst.Alarms[2] = 30
		}
		return
	}

	// linear tier upgrades (Charge_Tower → Charge_Storage → Powerful_Charges)
	if applyTowerUpgrade(inst, g) > 0 {
		// reset charge on upgrade
		inst.Vars["charge"] = 0.0
		inst.ImageIndex = 0
		form = b.refreshForm(inst, g)
		inst.Alarms[11] = int(math.Round(form.AlarmRate / b.attackRate))
		inst.Alarms[0] = int(math.Round(16 / b.attackRate))
		return
	}
}

func (b *ChargeTowerBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if activateChargeAbility(inst, g) {
		return
	}
	towerClickSelect(inst, g, towerSelectValue(8.0, inst))
}

// activateChargeAbility handles Charge_Overload and Lightning_Bomb abilities.
func activateChargeAbility(inst *engine.Instance, g *engine.Game) bool {
	abilityMax := getVar(inst, "ability_max")
	if abilityMax <= 0 {
		return false
	}
	if getVar(inst, "ability") < abilityMax {
		return false
	}
	inst.Vars["ability"] = 0.0
	g.AudioMgr.Play("Upgrade")

	name := chargeUpgradeName(inst)
	switch name {
	case "Charge_Overload":
		// supercharge burst: fire 12 Burst_Charge projectiles in all directions
		for i := 0; i < 12; i++ {
			angle := float64(i) * (2 * math.Pi / 12)
			proj := g.InstanceMgr.Create("Burst_Charge", inst.X, inst.Y)
			if proj != nil {
				speed := 28.0
				proj.HSpeed = math.Cos(angle) * speed
				proj.VSpeed = -math.Sin(angle) * speed
				proj.Direction = angle * 180 / math.Pi
				proj.ImageAngle = proj.Direction
				proj.Vars["LP"] = 2.0
				proj.Vars["PP"] = 20.0 + getVar(inst, "ppbuff")
				proj.Vars["leadpop"] = 1.0
				proj.Vars["camopop"] = 1.0
			}
		}
	case "Lightning_Bomb":
		// giant energy strike at the farthest bloon in range
		target := findNearestBloon(inst, g, 2000, true)
		tx, ty := inst.X, inst.Y
		if target != nil {
			tx, ty = target.X, target.Y
		}
		bomb := g.InstanceMgr.Create("Energy_Bomb", tx, ty)
		if bomb != nil {
			bomb.Vars["LP"] = 30.0
			bomb.Vars["PP"] = 2500.0
			bomb.Vars["leadpop"] = 1.0
			bomb.Vars["camopop"] = 1.0
			bomb.Vars["explode_radius"] = 150.0
			bomb.Alarms[1] = 40
		}
	}
	return true
}

// ---------- Charge Tower projectile behaviors ----------

// ChargeProjBehavior handles Charge_Proj, Powerful_Charge, and Burst_Charge.
// Simple linear projectile that pops bloons on contact and expires after a
// short lived alarm.
type ChargeProjBehavior struct {
	engine.DefaultBehavior
	hitRadius float64
}

func (b *ChargeProjBehavior) Create(inst *engine.Instance, g *engine.Game) {
	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = 1.0
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = 2.0
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = 0.0
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = 0.0
	}
	// short lifetime: expires after ~11 frames (enough at speed 20–26 to clear tower range)
	inst.Alarms[0] = 14
}

func (b *ChargeProjBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	r := b.hitRadius
	if r <= 0 {
		r = 10
	}
	projectileHitBloons(inst, g, r)
}

func (b *ChargeProjBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// OrbitalChargeBehavior handles Orbital_Charge and Magnetic_Charge.
// Each frame the projectile advances along a spiral GMX path (C_Path_Right/Up/Left/Down
// or S1/S2/S3/S4) at 29 px/frame relative to the tower that spawned it.
// Lifetime: alarm[0] set by the tower (100 or 132 frames). Hits bloons along the way.
type OrbitalChargeBehavior struct {
	engine.DefaultBehavior
}

func (b *OrbitalChargeBehavior) Create(inst *engine.Instance, g *engine.Game) {
	// Defaults — tower will override these before the first Step.
	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = 1.0
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = 10.0
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = 1.0
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = 1.0
	}
	if _, ok := inst.Vars["pathCycle"]; !ok {
		inst.Vars["pathCycle"] = 1.0
	}
	if _, ok := inst.Vars["pathDist"]; !ok {
		inst.Vars["pathDist"] = 0.0
	}
	if _, ok := inst.Vars["pathSpeed"]; !ok {
		inst.Vars["pathSpeed"] = 29.0
	}
	// anchor defaults to creation position; tower sets anchorX/anchorY explicitly.
	if _, ok := inst.Vars["anchorX"]; !ok {
		inst.Vars["anchorX"] = inst.X
	}
	if _, ok := inst.Vars["anchorY"]; !ok {
		inst.Vars["anchorY"] = inst.Y
	}
}

func (b *OrbitalChargeBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}

	// Determine which spiral path to follow.
	cycle := int(getVar(inst, "pathCycle"))
	var pathName string
	if inst.ObjectName == "Orbital_Charge" {
		switch cycle {
		case 1:
			pathName = "C_Path_Right"
		case 2:
			pathName = "C_Path_Up"
		case 3:
			pathName = "C_Path_Left"
		case 4:
			pathName = "C_Path_Down"
		default:
			pathName = "C_Path_Right"
		}
	} else { // Magnetic_Charge
		switch cycle {
		case 1:
			pathName = "S1"
		case 2:
			pathName = "S2"
		case 3:
			pathName = "S3"
		case 4:
			pathName = "S4"
		default:
			pathName = "S1"
		}
	}

	pa := g.PathMgr.Get(pathName)
	if pa == nil || len(pa.Points) == 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}

	// Advance distance along path.
	pathSpeed := getVar(inst, "pathSpeed")
	if pathSpeed <= 0 {
		pathSpeed = 29
	}
	pathDist := getVar(inst, "pathDist") + pathSpeed
	inst.Vars["pathDist"] = pathDist

	// Progress [0..1] along the path.
	progress := pathDist / pa.TotalLength
	if progress > 1.0 {
		progress = 1.0
	}

	// Current position on path (absolute path coordinates).
	rx, ry := g.PathMgr.GetPositionAtProgress(pathName, progress)
	originX := pa.Points[0].X
	originY := pa.Points[0].Y

	// Apply path relative to tower anchor (mirrors GM path_start with absolute=0).
	anchorX := getVar(inst, "anchorX")
	anchorY := getVar(inst, "anchorY")
	newX := anchorX + (rx - originX)
	newY := anchorY + (ry - originY)

	// Update image rotation from direction of travel.
	prevDist := pathDist - pathSpeed
	if prevDist >= 0 {
		prevProg := prevDist / pa.TotalLength
		if prevProg > 1.0 {
			prevProg = 1.0
		}
		px, py := g.PathMgr.GetPositionAtProgress(pathName, prevProg)
		ddx := (rx - originX) - (px - originX)
		ddy := (ry - originY) - (py - originY)
		if math.Abs(ddx) > 0.001 || math.Abs(ddy) > 0.001 {
			inst.ImageAngle = math.Atan2(-ddy, ddx) * 180 / math.Pi
		}
	}

	// Move to new position directly; set HSpeed/VSpeed to 0 so applyMotion won't drift.
	inst.X = newX
	inst.Y = newY
	inst.HSpeed = 0
	inst.VSpeed = 0

	projectileHitBloons(inst, g, 14)
}

func (b *OrbitalChargeBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// EnergyProjBehavior handles Small_Energy and Big_Energy (tesla lightning bolts).
// These travel fast and have a wide hit radius, approximating instant-strike chain lightning.
type EnergyProjBehavior struct {
	engine.DefaultBehavior
	hitRadius float64
}

func (b *EnergyProjBehavior) Create(inst *engine.Instance, g *engine.Game) {
	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = 1.0
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = 30.0
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = 1.0
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = 0.0
	}
	inst.Alarms[0] = 8 // short lifetime
}

func (b *EnergyProjBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	r := b.hitRadius
	if r <= 0 {
		r = 22
	}
	projectileHitBloons(inst, g, r)
}

func (b *EnergyProjBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// EnergyBombBehavior handles the Lightning_Bomb ability strike (Energy_Bomb).
// Grows into a huge AoE explosion dealing massive damage.
type EnergyBombBehavior struct {
	engine.DefaultBehavior
}

func (b *EnergyBombBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["LP"] = 30.0
	inst.Vars["PP"] = 2500.0
	inst.Vars["leadpop"] = 1.0
	inst.Vars["camopop"] = 1.0
	inst.Vars["explode_radius"] = 150.0
	inst.ImageXScale = 0.01
	inst.ImageYScale = 0.01
	inst.Alarms[1] = 40
}

func (b *EnergyBombBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// grow scale
	inst.ImageXScale += 0.025
	inst.ImageYScale += 0.025
	if inst.ImageXScale >= 1.0 {
		inst.ImageXScale = 1.0
		inst.ImageYScale = 1.0
	}
}

func (b *EnergyBombBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		// explode: huge AoE
		radius := getVar(inst, "explode_radius")
		if radius <= 0 {
			radius = 150
		}
		lp := getVar(inst, "LP")
		pp := getVar(inst, "PP")
		leadpop := getVar(inst, "leadpop")
		camopop := getVar(inst, "camopop")

		for _, bloon := range g.InstanceMgr.FindByObject("Normal_Bloon_Branch") {
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
			if math.Sqrt(dx*dx+dy*dy) > radius {
				continue
			}
			popBloon(bloon, lp, g)
			pp--
			if pp <= 0 {
				break
			}
		}

		explosion := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
		if explosion != nil {
			explosion.SpriteName = "Large_Explosion"
			explosion.ImageXScale = 3.0
			explosion.ImageYScale = 3.0
			explosion.Vars["LP"] = 1.0
			explosion.Vars["PP"] = 0.0
			explosion.Vars["explode_radius"] = 0.0
			explosion.Alarms[0] = 12
		}
		g.AudioMgr.Play("Large_Boom")
		g.InstanceMgr.Destroy(inst.ID)
	}
}
