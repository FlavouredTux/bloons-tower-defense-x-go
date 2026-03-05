package behaviors

import (
	"math"

	"btdx/internal/engine"
)

// boomerang Thrower — fires curving boomerangs that return to sender.

type BoomerangThrowerBehavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
	camoDetect float64
	leadDetect float64
}

type boomerangTowerForm struct {
	Range        float64
	AlarmBase    float64
	Projectile   string
	Speed        float64
	LP           float64
	PP           float64
	Lead         float64
	Camo         float64
	Life         float64
	SpinRate     float64
	HitRadius    float64
	Offsets      []float64
	CurveStart   float64
	CurveAccel   float64
	Juggle       float64 // bounces on lifetime expire (Plasmarang/Masterang)
	ShieldPop    float64 // -1 strips shields (Ricochet/King/Lord/Masterang)
	AuraInterval float64 // extra_pop spawn interval (Lord_Glaive = 5)
	AbilityMax   float64 // ability charge threshold (0 = no ability)
}

var boomerangTowerForms = map[string]boomerangTowerForm{
	// path: base → tier 1 → tier 2 (linear chain)
	"Boomerang_Thrower": {Range: 110, AlarmBase: 37, Projectile: "Boomerang", Speed: 20, LP: 1, PP: 3, Lead: 0, Camo: 0, Life: 19, SpinRate: 0, HitRadius: 20, Offsets: []float64{-15}, CurveStart: 32, CurveAccel: 2},
	"Multi_Pop_Thrower": {Range: 110, AlarmBase: 37, Projectile: "Boomerang", Speed: 20, LP: 1, PP: 7, Lead: 0, Camo: 0, Life: 19, SpinRate: 0, HitRadius: 20, Offsets: []float64{-15}, CurveStart: 32, CurveAccel: 2},
	"Glaive_Thrower":    {Range: 115, AlarmBase: 30, Projectile: "Glaive", Speed: 23, LP: 1, PP: 12, Lead: 1, Camo: 0, Life: 19, SpinRate: 15, HitRadius: 22, Offsets: []float64{-15}, CurveStart: 32, CurveAccel: 2},
	// path 1 (left): Plasmarangs — straight flight, bounces on expire
	"Plasmarangs":   {Range: 115, AlarmBase: 22, Projectile: "Plasmarang", Speed: 24, LP: 1, PP: 25, Lead: 1, Camo: 0, Life: 20, SpinRate: 0, HitRadius: 22, Offsets: []float64{0}, Juggle: 1},
	"Masterangs":    {Range: 125, AlarmBase: 17, Projectile: "Masterang", Speed: 24, LP: 1, PP: 25, Lead: 1, Camo: 1, Life: 20, SpinRate: 0, HitRadius: 22, Offsets: []float64{0}, Juggle: 7, ShieldPop: -1},
	"Megarang_Toss": {Range: 125, AlarmBase: 15, Projectile: "Masterang", Speed: 24, LP: 1, PP: 25, Lead: 1, Camo: 1, Life: 20, SpinRate: 0, HitRadius: 22, Offsets: []float64{0}, Juggle: 7, ShieldPop: -1, AbilityMax: 29},
	// path 2 (middle): Glaive Ricochet — straight flight, spins, wall bounce
	"Glaive_Ricochet": {Range: 115, AlarmBase: 30, Projectile: "Ricochet_Glaive", Speed: 24, LP: 1, PP: 40, Lead: 1, Camo: 1, Life: 120, SpinRate: 15, HitRadius: 24, Offsets: []float64{0}, ShieldPop: -1},
	"Glaive_King":     {Range: 120, AlarmBase: 30, Projectile: "King_Glaive", Speed: 33, LP: 1, PP: 100, Lead: 1, Camo: 1, Life: 180, SpinRate: 15, HitRadius: 25, Offsets: []float64{0}, ShieldPop: -1},
	"Glaive_Lord":     {Range: 120, AlarmBase: 30, Projectile: "Lord_Glaive", Speed: 18, LP: 1, PP: 1000, Lead: 1, Camo: 1, Life: 210, SpinRate: 15, HitRadius: 26, Offsets: []float64{0}, ShieldPop: -1, AuraInterval: 5},
	// path 3 (right): Bionic Boomer — fast attack, curving glaives
	"Bionic_Boomer": {Range: 115, AlarmBase: 7, Projectile: "Glaive", Speed: 23, LP: 1, PP: 12, Lead: 1, Camo: 0, Life: 19, SpinRate: 15, HitRadius: 22, Offsets: []float64{-15}, CurveStart: 32, CurveAccel: 2},
	"Doublerangs":   {Range: 115, AlarmBase: 7, Projectile: "Glaive", Speed: 23, LP: 1, PP: 12, Lead: 1, Camo: 0, Life: 19, SpinRate: 15, HitRadius: 22, Offsets: []float64{-6, -24}, CurveStart: 32, CurveAccel: 2},
	"Turbo_Charge":  {Range: 120, AlarmBase: 7, Projectile: "Turbo_Glaive", Speed: 23, LP: 1, PP: 12, Lead: 1, Camo: 0, Life: 19, SpinRate: 15, HitRadius: 22, Offsets: []float64{-6, -24}, CurveStart: 32, CurveAccel: 2, AbilityMax: 39},
}

func boomerangUpgradeName(inst *engine.Instance) string {
	return upgradeName(inst, "Boomerang_Thrower")
}

func boomerangFormFor(inst *engine.Instance) boomerangTowerForm {
	if form, ok := boomerangTowerForms[boomerangUpgradeName(inst)]; ok {
		return form
	}
	return boomerangTowerForms["Boomerang_Thrower"]
}

func fireBoomerang(g *engine.Game, tower, target *engine.Instance, form boomerangTowerForm, angleOffset, ppbuff, lead, camo float64) {
	if g == nil || tower == nil || target == nil {
		return
	}
	dx := target.X - tower.X
	dy := target.Y - tower.Y
	dir := 0.0
	if dx != 0 || dy != 0 {
		dir = math.Atan2(-dy, dx) * 180 / math.Pi
	}
	dir += angleOffset
	rad := dir * math.Pi / 180.0

	boom := g.InstanceMgr.Create(form.Projectile, tower.X, tower.Y)
	if boom == nil {
		return
	}

	boom.HSpeed = math.Cos(rad) * form.Speed
	boom.VSpeed = -math.Sin(rad) * form.Speed
	boom.Direction = dir
	boom.ImageAngle = dir
	boom.Vars["LP"] = form.LP
	boom.Vars["PP"] = form.PP + ppbuff
	boom.Vars["leadpop"] = lead
	boom.Vars["camopop"] = camo
	boom.Vars["speed"] = form.Speed
	boom.Vars["spin_rate"] = form.SpinRate
	boom.Vars["hit_radius"] = form.HitRadius
	boom.Vars["curve_start"] = form.CurveStart
	boom.Vars["curve_accel"] = form.CurveAccel
	boom.Vars["juggle"] = form.Juggle
	boom.Vars["shieldpop"] = form.ShieldPop
	boom.Vars["up"] = 1.0
	if form.AuraInterval > 0 {
		boom.Alarms[0] = int(form.AuraInterval)
	}
	boom.Alarms[1] = int(math.Round(form.Life))
}

func (b *BoomerangThrowerBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	b.rng = 110.0
	b.camoDetect = 0.0
	b.leadDetect = 0.0
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["range"] = b.rng
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["ability"] = 0.0
	inst.Vars["ability_max"] = 0.0
	form := boomerangFormFor(inst)
	inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
}

func (b *BoomerangThrowerBehavior) refreshForm(inst *engine.Instance, g *engine.Game) boomerangTowerForm {
	form := boomerangFormFor(inst)
	b.rng = form.Range
	b.camoDetect = form.Camo
	b.leadDetect = form.Lead
	inst.Vars["range"] = form.Range
	if code, ok := effectiveTowerCode(boomerangUpgradeName(inst)); ok {
		inst.Vars["tower_code"] = code
	}
	if spr := g.InstanceMgr.ObjectSpriteName(boomerangUpgradeName(inst)); spr != "" && g.AssetManager.GetSprite(spr) != nil {
		inst.SpriteName = spr
	}
	// configure ability threshold from form data.
	if form.AbilityMax > 0 {
		inst.Vars["ability_max"] = form.AbilityMax
	} else {
		inst.Vars["ability_max"] = 0.0
	}
	return form
}

func (b *BoomerangThrowerBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	form := b.refreshForm(inst, g)
	switch idx {
	case 0:
		// fire projectile
		if getVar(inst, "stun") > 0 {
			inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
			return
		}
		target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
		if target != nil {
			offsets := form.Offsets
			if len(offsets) == 0 {
				offsets = []float64{0}
			}
			for _, off := range offsets {
				fireBoomerang(g, inst, target, form, off, getVar(inst, "ppbuff"), b.leadDetect, b.camoDetect)
			}
			inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
		}
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))

	case 2:
		// ability charge loop (every 30 frames, +1 if not stunned and bloons present)
		inst.Alarms[2] = 30
		if getVar(inst, "ability_max") <= 0 || getVar(inst, "stun") > 0 {
			return
		}
		bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
		if len(bloons) > 0 || getGlobal(g, "wavenow") == 1 {
			cur := getVar(inst, "ability") + 1
			max := getVar(inst, "ability_max")
			if cur > max {
				cur = max
			}
			inst.Vars["ability"] = cur
		}

	case 3:
		// turbo Charge ability activation (triggered 1 frame after click)
		if getVar(inst, "ability") >= getVar(inst, "ability_max") && getVar(inst, "ability_max") > 0 {
			b.attackRate = 4.0
			inst.Vars["ability"] = 0.0
			inst.Alarms[4] = 300 // duration: 300 frames
			inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
		}

	case 4:
		// turbo Charge ability ends
		b.attackRate = 1.0
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
	}
}

func (b *BoomerangThrowerBehavior) Step(inst *engine.Instance, g *engine.Game) {
	form := b.refreshForm(inst, g)
	if applyPathUpgrade(inst, g) {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
		// start ability charge loop if the upgrade grants an ability.
		if form.AbilityMax > 0 && inst.Alarms[2] <= 0 {
			inst.Alarms[2] = 30
		}
		return
	}
	if applyTowerUpgrade(inst, g) > 0 {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
	}

	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target != nil {
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}
}

func (b *BoomerangThrowerBehavior) activateBoomerangAbility(inst *engine.Instance, g *engine.Game) bool {
	max := getVar(inst, "ability_max")
	if max <= 0 || getVar(inst, "ability") < max {
		return false
	}
	obj := boomerangUpgradeName(inst)
	switch obj {
	case "Turbo_Charge":
		// speed boost — trigger via alarm so it activates next frame.
		inst.Alarms[3] = 1
		return true
	case "Megarang_Toss":
		// fire a super Megarang at the strongest bloon in extended range.
		target := findNearestBloon(inst, g, b.rng+1000, b.camoDetect == 1)
		if target != nil {
			abilityForm := boomerangTowerForm{
				Projectile: "Megarang", Speed: 45, LP: 15, PP: 1000,
				Life: 20, HitRadius: 22,
				Juggle: 7, ShieldPop: -1,
			}
			fireBoomerang(g, inst, target, abilityForm, 0, getVar(inst, "ppbuff"), b.leadDetect, b.camoDetect)
		}
		inst.Vars["ability"] = 0.0
		return true
	}
	return false
}

func (b *BoomerangThrowerBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if b.activateBoomerangAbility(inst, g) {
		return
	}
	val := getVar(inst, "tower_code")
	if val <= 0 {
		val = towerSelectValue(3.0, inst)
	}
	towerClickSelect(inst, g, val)
}

// BoomerangBehavior — unified projectile behavior for all boomerang variants.
// handles curving return (Boomerang/Glaive/Turbo_Glaive), straight flight
// (Plasmarang/Masterang/Ricochet/King/Lord glaives), juggle bouncing,
// and Lord_Glaive's Extra_pop aura spawning.
type BoomerangBehavior struct {
	engine.DefaultBehavior
}

func (b *BoomerangBehavior) Create(inst *engine.Instance, g *engine.Game) {
	// defaults — usually overridden by fireBoomerang.
	initProjDefaults(inst, 1, 1, 0, 0)
	if _, ok := inst.Vars["speed"]; !ok {
		spd := math.Sqrt(inst.HSpeed*inst.HSpeed + inst.VSpeed*inst.VSpeed)
		if spd <= 0 {
			spd = 20
		}
		inst.Vars["speed"] = spd
	}
	if _, ok := inst.Vars["hit_radius"]; !ok {
		inst.Vars["hit_radius"] = 20.0
	}
	if _, ok := inst.Vars["spin_rate"]; !ok {
		inst.Vars["spin_rate"] = 0.0
	}
	if _, ok := inst.Vars["up"]; !ok {
		inst.Vars["up"] = 1.0
	}
	if _, ok := inst.Vars["curve_start"]; !ok {
		inst.Vars["curve_start"] = 0.0
	}
	if _, ok := inst.Vars["curve_accel"]; !ok {
		inst.Vars["curve_accel"] = 0.0
	}
	if _, ok := inst.Vars["juggle"]; !ok {
		inst.Vars["juggle"] = 0.0
	}
	if _, ok := inst.Vars["shieldpop"]; !ok {
		inst.Vars["shieldpop"] = 0.0
	}
}

func (b *BoomerangBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}

	speed := getVar(inst, "speed")
	if speed <= 0 {
		speed = 20
	}

	// curve logic: up += accel; if up < curveStart, direction += up
	curveAccel := getVar(inst, "curve_accel")
	curveStart := getVar(inst, "curve_start")
	if curveAccel > 0 && curveStart > 0 {
		up := getVar(inst, "up") + curveAccel
		inst.Vars["up"] = up
		if up < curveStart {
			inst.Direction += up
		}
	}

	// recalculate velocity from direction.
	rad := inst.Direction * math.Pi / 180.0
	inst.HSpeed = math.Cos(rad) * speed
	inst.VSpeed = -math.Sin(rad) * speed

	// wall bounce: reflect off room edges.
	if room := g.RoomManager.GetCurrent(); room != nil {
		rw := float64(room.Width)
		rh := float64(room.Height)
		bounced := false
		if inst.X <= 0 || inst.X >= rw {
			inst.HSpeed = -inst.HSpeed
			bounced = true
			if inst.X <= 0 {
				inst.X = 1
			} else {
				inst.X = rw - 1
			}
		}
		if inst.Y <= 0 || inst.Y >= rh {
			inst.VSpeed = -inst.VSpeed
			bounced = true
			if inst.Y <= 0 {
				inst.Y = 1
			} else {
				inst.Y = rh - 1
			}
		}
		if bounced {
			inst.Direction = math.Atan2(-inst.VSpeed, inst.HSpeed) * 180 / math.Pi
		}
	}

	// visual spin vs facing direction.
	spinRate := getVar(inst, "spin_rate")
	if spinRate != 0 {
		inst.ImageAngle += spinRate
	} else {
		inst.ImageAngle = inst.Direction
	}

	hitRadius := getVar(inst, "hit_radius")
	if hitRadius <= 0 {
		hitRadius = 20
	}
	projectileHitBloons(inst, g, hitRadius)
}

func (b *BoomerangBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	switch idx {
	case 0:
		// lord_Glaive aura: spawn Extra_pop every 5 frames while alive.
		exp := g.InstanceMgr.Create("Extra_pop", inst.X, inst.Y)
		if exp != nil {
			exp.Vars["LP"] = getVar(inst, "LP")
			exp.Vars["PP"] = getVar(inst, "PP")
			exp.Vars["leadpop"] = getVar(inst, "leadpop")
			exp.Vars["camopop"] = getVar(inst, "camopop")
			exp.Alarms[0] = 5
		}
		inst.Alarms[0] = 5

	case 1:
		// lifetime expired — juggle bounce or destroy.
		juggle := getVar(inst, "juggle")
		if juggle > 0 {
			speed := getVar(inst, "speed")
			// reverse direction and spawn a fresh copy with juggle-1.
			dir := inst.Direction + 180
			for dir >= 360 {
				dir -= 360
			}
			for dir < 0 {
				dir += 360
			}
			clone := g.InstanceMgr.Create(inst.ObjectName, inst.X, inst.Y)
			if clone != nil {
				rad := dir * math.Pi / 180.0
				clone.HSpeed = math.Cos(rad) * speed
				clone.VSpeed = -math.Sin(rad) * speed
				clone.Direction = dir
				clone.ImageAngle = dir
				clone.Vars["LP"] = getVar(inst, "LP")
				clone.Vars["PP"] = getVar(inst, "PP")
				clone.Vars["leadpop"] = getVar(inst, "leadpop")
				clone.Vars["camopop"] = getVar(inst, "camopop")
				clone.Vars["speed"] = speed
				clone.Vars["juggle"] = juggle - 1
				clone.Vars["spin_rate"] = getVar(inst, "spin_rate")
				clone.Vars["hit_radius"] = getVar(inst, "hit_radius")
				clone.Vars["shieldpop"] = getVar(inst, "shieldpop")
				clone.Vars["curve_start"] = 0.0
				clone.Vars["curve_accel"] = 0.0
				clone.Vars["up"] = 0.0
				clone.Alarms[1] = 20 // range=20 per bounce
			}
		}
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// ExtraPopBehavior — Lord_Glaive's periodic AoE damage pulse.
// spawned by Lord_Glaive every 5 frames; hits bloons in a small radius then dies.
type ExtraPopBehavior struct {
	engine.DefaultBehavior
}

func (b *ExtraPopBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 10, 1, 1)
}

func (b *ExtraPopBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 30)
}

func (b *ExtraPopBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}
