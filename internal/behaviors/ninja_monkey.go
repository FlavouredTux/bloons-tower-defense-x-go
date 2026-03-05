package behaviors

import (
	"math"
	"math/rand"

	"btdx/internal/engine"
)

// ninja Monkey — fast attack, camo detection, semi-homing shurikens.

type ninjaTowerForm struct {
	Range          float64
	AlarmBase      float64
	BaseProjectile string
	BaseSpeed      float64
	BaseLP         float64
	BasePP         float64
	BaseLead       float64
	BaseCamo       float64
	BaseLife       float64
	HitRadius      float64
	TrackRadius    float64
	TrackSpeed     float64
	SpinRate       float64
	TurnPerStep    float64
}

var ninjaTowerForms = map[string]ninjaTowerForm{
	"Ninja_Monkey":        {Range: 120, AlarmBase: 16, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 2, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Sharp_Shurikens":     {Range: 120, AlarmBase: 16, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 4, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Ninja_Training":      {Range: 130, AlarmBase: 12, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Distraction":         {Range: 125, AlarmBase: 12, BaseProjectile: "Distraction_Shot", BaseSpeed: 24, BaseLP: 1, BasePP: 7, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 100, TrackSpeed: 21, SpinRate: 20},
	"Double_Shot":         {Range: 135, AlarmBase: 12, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Flash_Bombs":         {Range: 125, AlarmBase: 11, BaseProjectile: "Distraction_Shot", BaseSpeed: 24, BaseLP: 1, BasePP: 7, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 100, TrackSpeed: 21, SpinRate: 20},
	"Mass_Distraction":    {Range: 125, AlarmBase: 11, BaseProjectile: "Distraction_Shot", BaseSpeed: 24, BaseLP: 1, BasePP: 9, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 100, TrackSpeed: 21, SpinRate: 20},
	"Sai_Ninja":           {Range: 125, AlarmBase: 11, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Katana_Ninja":        {Range: 125, AlarmBase: 8, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Hidden_Monkey":       {Range: 125, AlarmBase: 8, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Cursed_Katana_Ninja": {Range: 133, AlarmBase: 6, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 9, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Cursed_Blade_Ninja":  {Range: 143, AlarmBase: 6, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 9, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Bloonjitzu":          {Range: 139, AlarmBase: 10, BaseProjectile: "Shuriken", BaseSpeed: 21, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Ninja_God":           {Range: 144, AlarmBase: 6, BaseProjectile: "Golden_Ninja_Star", BaseSpeed: 23, BaseLP: 2, BasePP: 6, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 100, TrackSpeed: 23, SpinRate: 20},
}

func ninjaUpgradeName(inst *engine.Instance) string {
	return upgradeName(inst, "Ninja_Monkey")
}

func ninjaFormFor(inst *engine.Instance) ninjaTowerForm {
	if form, ok := ninjaTowerForms[ninjaUpgradeName(inst)]; ok {
		return form
	}
	return ninjaTowerForms["Ninja_Monkey"]
}

type NinjaMonkeyBehavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
}

func (b *NinjaMonkeyBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	b.rng = 120.0
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["range"] = b.rng
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["cycle"] = 0.0
	form := b.refreshForm(inst, g)
	inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
}

func (b *NinjaMonkeyBehavior) refreshForm(inst *engine.Instance, g *engine.Game) ninjaTowerForm {
	form := ninjaFormFor(inst)
	b.rng = form.Range
	inst.Vars["range"] = form.Range
	if code, ok := effectiveTowerCode(ninjaUpgradeName(inst)); ok {
		inst.Vars["tower_code"] = code
	}

	obj := ninjaUpgradeName(inst)

	if spr := g.InstanceMgr.ObjectSpriteName(obj); spr != "" && g.AssetManager.GetSprite(spr) != nil {
		inst.SpriteName = spr
	}

	configureTowerAbility(inst, obj)
	if obj == "Hidden_Monkey" {
		if getVar(inst, "hidden_revealed") > 0 {
			inst.SpriteName = "Hidden_Monkey_Spr"
		} else if getVar(inst, "hidden_timer") > 0 {
			inst.SpriteName = "UnChained"
		}
	}

	return form
}

func fireNinjaProjectile(g *engine.Game, tower, target *engine.Instance, projectile string, speed, angleOffset, lp, pp, lead, camo, life, hitRadius, trackRadius, trackSpeed, spinRate, turnPerStep float64) {
	if g == nil || tower == nil || target == nil || projectile == "" || life <= 0 {
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
	proj := g.InstanceMgr.Create(projectile, tower.X, tower.Y)
	if proj == nil {
		return
	}
	proj.HSpeed = math.Cos(rad) * speed
	proj.VSpeed = -math.Sin(rad) * speed
	proj.Direction = dir
	proj.ImageAngle = dir
	proj.Vars["LP"] = lp
	proj.Vars["PP"] = pp
	proj.Vars["leadpop"] = lead
	proj.Vars["camopop"] = camo
	proj.Vars["targetID"] = float64(target.ID)
	proj.Vars["track_radius"] = trackRadius
	proj.Vars["track_speed"] = trackSpeed
	proj.Vars["hit_radius"] = hitRadius
	proj.Vars["spin_rate"] = spinRate
	proj.Vars["turn_per_step"] = turnPerStep
	proj.Vars["speed"] = speed
	proj.Alarms[0] = int(math.Round(life))
}

func fireNinjaVolley(g *engine.Game, tower, target *engine.Instance, projectile string, speed, lp, pp, lead, camo, life, hitRadius, trackRadius, trackSpeed, spinRate float64) {
	offsets := []float64{8, 4, 0, -4, -8}
	for _, off := range offsets {
		fireNinjaProjectile(g, tower, target, projectile, speed, off, lp, pp, lead, camo, life, hitRadius, trackRadius, trackSpeed, spinRate, 0)
	}
}

func (b *NinjaMonkeyBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx != 0 {
		return
	}
	form := b.refreshForm(inst, g)
	if getVar(inst, "stun") > 0 {
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
		return
	}
	target := findNearestBloon(inst, g, b.rng, true)
	if target != nil {
		obj := ninjaUpgradeName(inst)
		ppbuff := getVar(inst, "ppbuff")
		cycle := int(math.Round(getVar(inst, "cycle")))

		// every ninja variant fires its base shuriken every attack.
		fireNinjaProjectile(g, inst, target, form.BaseProjectile, form.BaseSpeed, 0,
			form.BaseLP, form.BasePP+ppbuff, form.BaseLead, form.BaseCamo, form.BaseLife,
			form.HitRadius, form.TrackRadius, form.TrackSpeed, form.SpinRate, form.TurnPerStep)

		switch obj {
		case "Double_Shot":
			fireNinjaProjectile(g, inst, target, "Shuriken", form.BaseSpeed, 0, 1, 5+ppbuff, 0, 1, 15, 18, 90, 21, 20, 0)
		case "Flash_Bombs", "Mass_Distraction":
			cycle++
			if cycle >= 4 {
				fireNinjaProjectile(g, inst, target, "Flash_Bomb_Proj", 24, 0, 1, 60+ppbuff, 0, 1, 15, 18, 0, 0, 0, 0)
				cycle = 0
			}
		case "Sai_Ninja":
			cycle++
			if cycle >= 2 {
				fireNinjaProjectile(g, inst, target, "Sai", 4, 96, 1, 30+ppbuff, 0, 1, 5, 24, 0, 0, 0, -36)
				cycle = 0
			}
		case "Katana_Ninja", "Hidden_Monkey":
			cycle++
			if cycle >= 2 {
				fireNinjaProjectile(g, inst, target, "Katana", 4, 96, 2, 45+ppbuff, 0, 1, 5, 24, 0, 0, 0, -36)
				cycle = 0
			}
			if obj == "Hidden_Monkey" {
				setTowerAbilityCharge(inst, 1)
			}
		case "Cursed_Katana_Ninja":
			cycle++
			if cycle >= 2 {
				fireNinjaProjectile(g, inst, target, "Cursed_Katana", 4, 96, 2, 66+ppbuff, 0, 1, 5, 24, 0, 0, 0, -36)
				cycle = 0
			}
		case "Cursed_Blade_Ninja":
			fireNinjaProjectile(g, inst, target, "Cursed_Blade", 4, 96, 7, 99+ppbuff, 1, 1, 5, 24, 0, 0, 0, -36)
			fireNinjaProjectile(g, inst, target, "Sai", 3, 114, 1, 33+ppbuff, 0, 1, 6, 24, 0, 0, 0, -36)
			fireNinjaProjectile(g, inst, target, "Alt_Sai", 3, -114, 1, 33+ppbuff, 0, 1, 6, 24, 0, 0, 0, -36)
		case "Bloonjitzu":
			fireNinjaVolley(g, inst, target, "Shuriken", 21, 1, 5+ppbuff, 0, 1, 15, 18, 90, 21, 20)
		case "Ninja_God":
			fireNinjaVolley(g, inst, target, "Golden_Ninja_Star", 23, 2, 6+ppbuff, 0, 1, 15, 18, 100, 23, 20)
		}
		inst.Vars["cycle"] = float64(cycle)
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}
	inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
}

func (b *NinjaMonkeyBehavior) Step(inst *engine.Instance, g *engine.Game) {
	form := b.refreshForm(inst, g)
	if applyPathUpgrade(inst, g) {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
		return
	}
	if applyTowerUpgrade(inst, g) > 0 {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
	}

	obj := ninjaUpgradeName(inst)
	if obj == "Hidden_Monkey" {
		timer := getVar(inst, "hidden_timer")
		if timer > 0 {
			if timer == 540 {
				b.refreshForm(inst, g)
			}
			timer--
			inst.Vars["hidden_timer"] = timer

			if timer <= 0 {
				inst.Vars["hidden_revealed"] = 1.0
				b.refreshForm(inst, g)
			} else {
				if int(timer)%9 == 0 {
					target := findNearestBloon(inst, g, b.rng, true)
					if target != nil {
						// 50 degree melee sweep
						fireNinjaProjectile(g, inst, target, "Crouching_Blade", 0, 25, 15, 45+getVar(inst, "ppbuff"), 0, 1, 5, 24, 0, 0, 0, -10)
					}
				}
			}
		}
	}

	target := findNearestBloon(inst, g, b.rng, true)
	if target != nil {
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}
}

func (b *NinjaMonkeyBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	towerClickSelect(inst, g, towerSelectValue(5.0, inst))
}

// ShurikenBehavior — semi-homing projectile
type ShurikenBehavior struct {
	engine.DefaultBehavior
}

func (b *ShurikenBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 2, 0, 1)
	if _, ok := inst.Vars["targetID"]; !ok {
		inst.Vars["targetID"] = 0.0
	}
	if _, ok := inst.Vars["track_radius"]; !ok {
		inst.Vars["track_radius"] = 90.0
	}
	if _, ok := inst.Vars["track_speed"]; !ok {
		inst.Vars["track_speed"] = 21.0
	}
	if _, ok := inst.Vars["hit_radius"]; !ok {
		inst.Vars["hit_radius"] = 18.0
	}
	if _, ok := inst.Vars["spin_rate"]; !ok {
		inst.Vars["spin_rate"] = 20.0
	}
	if _, ok := inst.Vars["turn_per_step"]; !ok {
		inst.Vars["turn_per_step"] = 0.0
	}
	if _, ok := inst.Vars["speed"]; !ok {
		speed := math.Sqrt(inst.HSpeed*inst.HSpeed + inst.VSpeed*inst.VSpeed)
		if speed <= 0 {
			speed = 24.0
		}
		inst.Vars["speed"] = speed
	}
}

func (b *ShurikenBehavior) Step(inst *engine.Instance, g *engine.Game) {
	turnPerStep := getVar(inst, "turn_per_step")
	if turnPerStep != 0 {
		inst.Direction += turnPerStep
		speed := getVar(inst, "speed")
		if speed <= 0 {
			speed = 24
		}
		rad := inst.Direction * math.Pi / 180.0
		inst.HSpeed = math.Cos(rad) * speed
		inst.VSpeed = -math.Sin(rad) * speed
	}

	// semi-homing for star-style projectiles.
	targetID := int(getVar(inst, "targetID"))
	target := g.InstanceMgr.GetByID(targetID)
	trackRadius := getVar(inst, "track_radius")
	trackSpeed := getVar(inst, "track_speed")
	if target != nil && !target.Destroyed {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if trackRadius > 0 && dist < trackRadius && dist > 0 {
			inst.HSpeed = (dx / dist) * trackSpeed
			inst.VSpeed = (dy / dist) * trackSpeed
			inst.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
	}
	spinRate := getVar(inst, "spin_rate")
	if spinRate != 0 {
		inst.ImageAngle += spinRate
	} else {
		inst.ImageAngle = inst.Direction
	}

	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	hitRadius := getVar(inst, "hit_radius")
	if hitRadius <= 0 {
		hitRadius = 18
	}
	projectileHitBloons(inst, g, hitRadius)
}

func (b *ShurikenBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

type DistractionShotBehavior struct {
	engine.DefaultBehavior
}

func (b *DistractionShotBehavior) Create(inst *engine.Instance, g *engine.Game) {
	(&ShurikenBehavior{}).Create(inst, g)
	inst.Vars["track_radius"] = 100.0
}

func (b *DistractionShotBehavior) Step(inst *engine.Instance, g *engine.Game) {
	targetID := int(getVar(inst, "targetID"))
	target := g.InstanceMgr.GetByID(targetID)
	if target != nil && !target.Destroyed {
		dx := target.X - inst.X
		dy := target.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		trackRadius := getVar(inst, "track_radius")
		trackSpeed := getVar(inst, "track_speed")
		if trackRadius > 0 && dist < trackRadius && dist > 0 {
			inst.HSpeed = (dx / dist) * trackSpeed
			inst.VSpeed = (dy / dist) * trackSpeed
			inst.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
		}
	}
	inst.ImageAngle += getVar(inst, "spin_rate")

	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	lp := getVar(inst, "LP")
	leadpop := getVar(inst, "leadpop")
	camopop := getVar(inst, "camopop")
	hitRadius := getVar(inst, "hit_radius")
	if hitRadius <= 0 {
		hitRadius = 18
	}

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
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist >= hitRadius {
			continue
		}
		popBloon(bloon, lp, g)
		// distraction shot has a 1/3 chance to push bloons back.
		if rand.Intn(3) == 0 {
			bloon.Vars["distraction"] = 1.0
			if bloon.Alarms[9] < 30 {
				bloon.Alarms[9] = 30
			}
		}
		pp -= 1
		inst.Vars["PP"] = pp
		if pp <= 0 {
			g.InstanceMgr.Destroy(inst.ID)
			return
		}
	}
}

func (b *DistractionShotBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

type FlashBombProjBehavior struct {
	engine.DefaultBehavior
}

func (b *FlashBombProjBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	initProjDefaults(inst, 1, 60, 0, 1)
}

func (b *FlashBombProjBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
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
		if math.Sqrt(dx*dx+dy*dy) > 16 {
			continue
		}
		flash := g.InstanceMgr.Create("Flash", inst.X, inst.Y)
		if flash != nil {
			flash.ImageXScale = 1.4
			flash.ImageYScale = 1.4
			flash.Vars["LP"] = getVar(inst, "LP")
			flash.Vars["PP"] = getVar(inst, "PP")
			flash.Vars["leadpop"] = leadpop
			flash.Vars["camopop"] = camopop
			flash.Vars["impact"] = 1.0
			flash.Vars["range"] = 8.0
			flash.Alarms[1] = 8
		}
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
}

type FlashBehavior struct {
	engine.DefaultBehavior
}

func (b *FlashBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 60, 0, 1)
	if _, ok := inst.Vars["impact"]; !ok {
		inst.Vars["impact"] = 0.0
	}
}

func (b *FlashBehavior) Step(inst *engine.Instance, g *engine.Game) {
	lp := getVar(inst, "LP")
	pp := getVar(inst, "PP")
	leadpop := getVar(inst, "leadpop")
	camopop := getVar(inst, "camopop")
	impact := getVar(inst, "impact")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}

	// flash lives briefly and applies a short AoE hit+stun pulse.
	const radius = 44.0
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
		if impact == 1 {
			bloon.Vars["stun"] = 1.0
			if bloon.Alarms[8] < 20 {
				bloon.Alarms[8] = 20
			}
		}
		pp -= 1
		inst.Vars["PP"] = pp
		if pp <= 0 {
			break
		}
	}
}

func (b *FlashBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}
