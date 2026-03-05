package behaviors

import (
	"math"

	"btdx/internal/engine"
)

// Glue Gunner — glues bloons to slow them.

type GlueGunnerBehavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
	alarmBase  float64
}

func (b *GlueGunnerBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	b.rng = 111.0
	b.alarmBase = 33.0
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["range"] = b.rng
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["gluepotency"] = 1.0
	inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
}

func (b *GlueGunnerBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		if getVar(inst, "stun") > 0 {
			inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
			return
		}
		// smart targeting: prefer unglued bloons
		target := findUngluedBloon(inst, g, b.rng, false)
		if target != nil {
			glob := g.InstanceMgr.Create("Glue_Glob", inst.X, inst.Y)
			if glob != nil {
				dx := target.X - inst.X
				dy := target.Y - inst.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				speed := 21.0
				if dist > 0 {
					glob.HSpeed = (dx / dist) * speed
					glob.VSpeed = (dy / dist) * speed
					glob.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
				}
				glob.ImageAngle = glob.Direction
				glob.Vars["LP"] = 0.0
				glob.Vars["PP"] = 1.0 + getVar(inst, "ppbuff")
				glob.Vars["leadpop"] = 1.0
				glob.Vars["camopop"] = 0.0
				glob.Vars["gluepotency"] = getVar(inst, "gluepotency")
				glob.Alarms[0] = 11

				inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
			}
		}
		inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
	}
}

func (b *GlueGunnerBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.Vars["range"] = b.rng
	if applyPathUpgrade(inst, g) {
		return
	}
	switch applyTowerUpgrade(inst, g) {
	case 1:
		// piercing Glue
		inst.Vars["ppbuff"] = getVar(inst, "ppbuff") + 1
		inst.Vars["gluepotency"] = getVar(inst, "gluepotency") + 0.3
	case 2:
		// corrosive Glue (partial until DoT behavior is added)
		inst.Vars["ppbuff"] = getVar(inst, "ppbuff") + 1
		inst.Vars["gluepotency"] = getVar(inst, "gluepotency") + 0.3
	}

	target := findUngluedBloon(inst, g, b.rng, false)
	if target != nil {
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}
}

func (b *GlueGunnerBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	towerClickSelect(inst, g, towerSelectValue(9.0, inst))
}

// GlueGlobBehavior — applies glue slow on contact
type GlueGlobBehavior struct {
	engine.DefaultBehavior
}

func (b *GlueGlobBehavior) Create(inst *engine.Instance, g *engine.Game) {
	setProjDefaults(inst, 0, 1, 1, 0)
	inst.Vars["gluepotency"] = 1.0
}

func (b *GlueGlobBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	pp := getVar(inst, "PP")
	if pp <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	gluePotency := getVar(inst, "gluepotency")
	// check for bloon contact
	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < 18 {
			// apply glue
			bloon.Vars["glue"] = gluePotency
			inst.Vars["PP"] = pp - 1
			if pp-1 <= 0 {
				g.InstanceMgr.Destroy(inst.ID)
				return
			}
		}
	}
}

func (b *GlueGlobBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}
