package behaviors

import (
	"math"

	"btdx/internal/engine"
)

// Ice Monkey — freezes nearby bloons (does NOT pop them).

type IceMonkeyBehavior struct {
	engine.DefaultBehavior
	attackRate    float64
	rng           float64
	alarmBase     float64
	alarmCooldown float64
}

func (b *IceMonkeyBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	b.rng = 85.0
	b.alarmBase = 33.0
	b.alarmCooldown = 100.0 // 33 warmup + 67 cooldown
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["range"] = b.rng
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["icepotency"] = 1.0
	inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
}

func (b *IceMonkeyBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		if getVar(inst, "stun") > 0 {
			inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
			return
		}
		// check for bloons in range
		target := findNearestBloon(inst, g, b.rng, false)
		if target != nil {
			// create ice aura
			aura := g.InstanceMgr.Create("Ice_Aura", inst.X, inst.Y)
			if aura != nil {
				aura.Vars["PP"] = 40.0 + getVar(inst, "ppbuff")
				aura.Vars["icepotency"] = getVar(inst, "icepotency")
				aura.Alarms[0] = 9
			}
			// visual feedback
			inst.ImageXScale = 0.9
			inst.ImageYScale = 0.9
		}
		inst.Alarms[0] = int(math.Round(b.alarmCooldown / b.attackRate))
	}
}

func (b *IceMonkeyBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.Vars["range"] = b.rng
	if applyPathUpgrade(inst, g) {
		return
	}
	switch applyTowerUpgrade(inst, g) {
	case 1:
		// enhanced Freeze
		inst.Vars["icepotency"] = getVar(inst, "icepotency") + 0.5
		b.rng += 8
	case 2:
		// perma Frost
		inst.Vars["icepotency"] = getVar(inst, "icepotency") + 0.5
		b.rng += 8
		inst.Vars["ppbuff"] = getVar(inst, "ppbuff") + 10
	}

	if inst.ImageXScale < 1.0 {
		inst.ImageXScale += 0.02
		inst.ImageYScale += 0.02
		if inst.ImageXScale > 1.0 {
			inst.ImageXScale = 1.0
			inst.ImageYScale = 1.0
		}
	}
}

func (b *IceMonkeyBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	towerClickSelect(inst, g, towerSelectValue(10.0, inst))
}

// IceAuraBehavior — stationary freeze AoE
type IceAuraBehavior struct {
	engine.DefaultBehavior
}

func (b *IceAuraBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["PP"] = 40.0
	inst.Vars["icepotency"] = 1.0
	inst.HSpeed = 0
	inst.VSpeed = 0
}

func (b *IceAuraBehavior) Step(inst *engine.Instance, g *engine.Game) {
	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	// freeze bloons in radius (does NOT pop them)
	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	for _, bloon := range bloons {
		if bloon.Destroyed || pp <= 0 {
			continue
		}
		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < 85 {
			// apply freeze
			bloon.Vars["freeze"] = 1.0
			bloon.Alarms[6] = 30 // frozen for 30 frames
			pp--
			inst.Vars["PP"] = pp
		}
	}
}

func (b *IceAuraBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}
