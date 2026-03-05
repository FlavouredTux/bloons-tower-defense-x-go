package behaviors

import (
	"math"
	"math/rand"

	"btdx/internal/engine"
)

// Bomb Cannon — fires bombs that explode into AoE damage.

type bombTowerForm struct {
	Range      float64
	AlarmBase  float64
	Projectile string
	Speed      float64
	LP         float64
	PP         float64
	ExplodeR   float64
	Lead       float64
	Camo       float64
	AbilityMax float64
}

var bombTowerForms = map[string]bombTowerForm{
	// base
	"Bomb_Cannon": {Range: 112, AlarmBase: 48, Projectile: "Bomb", Speed: 20, LP: 1, PP: 20, ExplodeR: 40, Lead: 1, Camo: 0},
	// tier 1
	"Big_Bombs": {Range: 115, AlarmBase: 48, Projectile: "Bomb", Speed: 20, LP: 1, PP: 40, ExplodeR: 40, Lead: 1, Camo: 0},
	// tier 2
	"Missile_Launcher": {Range: 135, AlarmBase: 36, Projectile: "Bomb", Speed: 30, LP: 1, PP: 40, ExplodeR: 40, Lead: 1, Camo: 0},
	// path 1 (left): Bloon_Buster_Cannon → Moab_Mauler → Moab_Assassin_Cannon
	"Bloon_Buster_Cannon":  {Range: 135, AlarmBase: 30, Projectile: "Bomb", Speed: 30, LP: 2, PP: 40, ExplodeR: 40, Lead: 1, Camo: 0},
	"Moab_Mauler":          {Range: 135, AlarmBase: 30, Projectile: "Bomb", Speed: 30, LP: 3, PP: 40, ExplodeR: 40, Lead: 1, Camo: 0},
	"Moab_Assassin_Cannon": {Range: 150, AlarmBase: 30, Projectile: "Bomb", Speed: 45, LP: 3, PP: 40, ExplodeR: 40, Lead: 1, Camo: 1, AbilityMax: 25},
	// path 2 (middle): Cluster_Bombs → Bloon_Impactor → Explosion_King
	"Cluster_Bombs":  {Range: 135, AlarmBase: 35, Projectile: "Cluster_Bomb", Speed: 24, LP: 1, PP: 40, ExplodeR: 40, Lead: 1, Camo: 0},
	"Bloon_Impactor": {Range: 135, AlarmBase: 33, Projectile: "Impact_Bomb", Speed: 24, LP: 1, PP: 50, ExplodeR: 40, Lead: 1, Camo: 0},
	"Explosion_King": {Range: 135, AlarmBase: 30, Projectile: "King_Bomb", Speed: 24, LP: 2, PP: 60, ExplodeR: 50, Lead: 1, Camo: 0},
	// path 3 (right): Pineapple_Launcher → Mega_Fruit_Cannon
	"Pineapple_Launcher": {Range: 112, AlarmBase: 16, Projectile: "Bomb", Speed: 24, LP: 1, PP: 40, ExplodeR: 90, Lead: 1, Camo: 1},
	"Mega_Fruit_Cannon":  {Range: 112, AlarmBase: 14, Projectile: "Bomb", Speed: 24, LP: 2, PP: 50, ExplodeR: 90, Lead: 1, Camo: 1},
	// special path: Pop_Cannon → Explosion_Machine
	"Pop_Cannon":        {Range: 135, AlarmBase: 30, Projectile: "Bomb", Speed: 6, LP: 1, PP: 100, ExplodeR: 90, Lead: 1, Camo: 0},
	"Explosion_Machine": {Range: 135, AlarmBase: 28, Projectile: "Bomb", Speed: 6, LP: 2, PP: 120, ExplodeR: 90, Lead: 1, Camo: 0},
}

func bombUpgradeName(inst *engine.Instance) string {
	return upgradeName(inst, "Bomb_Cannon")
}

func bombFormFor(inst *engine.Instance) bombTowerForm {
	if form, ok := bombTowerForms[bombUpgradeName(inst)]; ok {
		return form
	}
	return bombTowerForms["Bomb_Cannon"]
}

type BombCannonBehavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
	camoDetect float64
	leadDetect float64
}

func (b *BombCannonBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["ability"] = 0.0
	inst.Vars["ability_max"] = 0.0
	form := b.refreshForm(inst, g)
	inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
}

func (b *BombCannonBehavior) refreshForm(inst *engine.Instance, g *engine.Game) bombTowerForm {
	form := bombFormFor(inst)
	b.rng = form.Range
	b.camoDetect = form.Camo
	b.leadDetect = form.Lead
	inst.Vars["range"] = form.Range
	if code, ok := effectiveTowerCode(bombUpgradeName(inst)); ok {
		inst.Vars["tower_code"] = code
	}
	if spr := g.InstanceMgr.ObjectSpriteName(bombUpgradeName(inst)); spr != "" && g.AssetManager.GetSprite(spr) != nil {
		inst.SpriteName = spr
	}
	if form.AbilityMax > 0 {
		inst.Vars["ability_max"] = form.AbilityMax
	} else {
		inst.Vars["ability_max"] = 0.0
	}
	return form
}

func (b *BombCannonBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		if getVar(inst, "stun") > 0 {
			form := bombFormFor(inst)
			inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
			return
		}
		form := bombFormFor(inst)
		target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
		if target != nil {
			projectileName := form.Projectile
			if projectileName == "" {
				projectileName = "Bomb"
			}
			bomb := g.InstanceMgr.Create(projectileName, inst.X, inst.Y)
			if bomb != nil {
				dx := target.X - inst.X
				dy := target.Y - inst.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > 0 {
					bomb.HSpeed = (dx / dist) * form.Speed
					bomb.VSpeed = (dy / dist) * form.Speed
					bomb.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
				}
				bomb.ImageAngle = bomb.Direction
				bomb.Vars["LP"] = form.LP
				bomb.Vars["PP"] = form.PP + getVar(inst, "ppbuff")
				bomb.Vars["leadpop"] = form.Lead
				bomb.Vars["camopop"] = form.Camo
				bomb.Vars["explode_radius"] = form.ExplodeR
				bomb.Alarms[1] = 18

				inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
			}
		}
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
	}
	if idx == 2 {
		// ability charge (Moab_Assassin_Cannon)
		if getGlobal(g, "play") == 1 {
			inst.Vars["ability"] = getVar(inst, "ability") + 1
		}
		inst.Alarms[2] = 30
	}
}

func (b *BombCannonBehavior) Step(inst *engine.Instance, g *engine.Game) {
	form := b.refreshForm(inst, g)
	if applyPathUpgrade(inst, g) {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
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

func (b *BombCannonBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if activateBombAbility(inst, g) {
		return
	}
	towerClickSelect(inst, g, towerSelectValue(6.0, inst))
}

// activateBombAbility — Moab Assassin Cannon's MOAB missile ability
func activateBombAbility(inst *engine.Instance, g *engine.Game) bool {
	abilityMax := getVar(inst, "ability_max")
	if abilityMax <= 0 {
		return false
	}
	if getVar(inst, "ability") < abilityMax {
		return false
	}
	inst.Vars["ability"] = 0.0
	// fire a massive missile at the strongest bloon
	target := findNearestBloon(inst, g, 2000, true)
	if target == nil {
		return true
	}
	missile := g.InstanceMgr.Create("Bomb", inst.X, inst.Y)
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
		missile.ImageXScale = 2.0
		missile.ImageYScale = 2.0
		missile.Vars["LP"] = 1000.0
		missile.Vars["PP"] = 1.0 + getVar(inst, "ppbuff")
		missile.Vars["leadpop"] = 1.0
		missile.Vars["camopop"] = 1.0
		missile.Vars["explode_radius"] = 40.0
		missile.Alarms[1] = 30
	}
	g.AudioMgr.Play("Upgrade")
	return true
}

// BombBehavior — explodes on contact with any bloon, or on timeout
type BombBehavior struct {
	engine.DefaultBehavior
}

func (b *BombBehavior) Create(inst *engine.Instance, g *engine.Game) {
	setProjDefaults(inst, 1, 20, 1, 0)
}

func (b *BombBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction

	// Cluster fragments don't explode on collision, only on timeout
	if inst.ObjectName == "Cluster" || inst.ObjectName == "Impactor" {
		return
	}

	// check for contact with any bloon to explode
	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < 24 {
			b.explode(inst, g)
			return
		}
	}
}

func (b *BombBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 1 {
		// timeout — explode anyway
		b.explode(inst, g)
	}
}

func (b *BombBehavior) explode(inst *engine.Instance, g *engine.Game) {
	radius := getVar(inst, "explode_radius")
	if radius <= 0 {
		radius = 40
	}

	// create explosion at bomb position
	explosion := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
	if explosion != nil {
		explosion.Vars["LP"] = getVar(inst, "LP")
		explosion.Vars["PP"] = getVar(inst, "PP")
		explosion.Vars["leadpop"] = getVar(inst, "leadpop")
		explosion.Vars["camopop"] = getVar(inst, "camopop")
		explosion.Vars["explode_radius"] = radius

		// Port old GameMaker VFX logic by updating sprite and scale based on blast radius
		if radius >= 90 {
			explosion.SpriteName = "Large_Explosion"
			explosion.ImageXScale = 1.5
			explosion.ImageYScale = 1.5
		} else if radius > 45 {
			explosion.SpriteName = "Medium_Explosion_Spr"
			explosion.ImageXScale = 1.1
			explosion.ImageYScale = 1.1
		} else {
			explosion.SpriteName = "S_Explosion_Spr"
			explosion.ImageXScale = 1.1
			explosion.ImageYScale = 1.1
		}

		explosion.Alarms[0] = 8
	}

	if radius >= 90 {
		g.AudioMgr.Play("Large_Boom")
	} else {
		g.AudioMgr.Play("Small_Boom")
	}

	// Port old GameMaker sub-munitions for Cluster and Impact variants
	switch inst.ObjectName {
	case "Cluster_Bomb", "Cluster_BombS_Proj":
		mul := 1.0
		for i := 0; i < 8; i++ {
			c := g.InstanceMgr.Create("Cluster", inst.X, inst.Y)
			if c != nil {
				c.Vars["LP"] = 1.0
				c.Vars["PP"] = 25.0
				c.Vars["leadpop"] = getVar(inst, "leadpop")
				c.Vars["camopop"] = getVar(inst, "camopop")
				c.Vars["explode_radius"] = 40.0
				c.MotionSet(45.0*mul, 27)
				c.Alarms[1] = 4
			}
			mul += 1.0
		}
	case "Impact_Bomb":
		mul := 1.0
		for i := 0; i < 8; i++ {
			c := g.InstanceMgr.Create("Impactor", inst.X, inst.Y)
			if c != nil {
				c.Vars["LP"] = 1.0
				c.Vars["PP"] = 20.0
				c.Vars["leadpop"] = getVar(inst, "leadpop")
				c.Vars["camopop"] = getVar(inst, "camopop")
				c.Vars["explode_radius"] = 40.0
				c.MotionSet(45.0*mul, 27)
				c.Alarms[1] = 4
			}
			mul += 1.0
		}
	case "King_Bomb":
		mul := 1.0
		for i := 0; i < 20; i++ {
			c := g.InstanceMgr.Create("Impactor", inst.X, inst.Y)
			if c != nil {
				c.Vars["LP"] = 1.0
				c.Vars["PP"] = 20.0
				c.Vars["leadpop"] = getVar(inst, "leadpop")
				c.Vars["camopop"] = getVar(inst, "camopop")
				c.Vars["explode_radius"] = 40.0
				c.MotionSet(36.0*mul-15.0+float64(rand.Intn(30)), 20)
				c.Alarms[1] = 1 + rand.Intn(7)
			}
			mul += 1.0
		}
	}

	g.InstanceMgr.Destroy(inst.ID)
}

// SmallExplosionBehavior — AoE damage with high pierce
type SmallExplosionBehavior struct {
	engine.DefaultBehavior
}

func (b *SmallExplosionBehavior) Create(inst *engine.Instance, g *engine.Game) {
	setProjDefaults(inst, 1, 20, 1, 0)

	// Set the image speed so the animation plays naturally (GM default is 1 frame per step)
	inst.ImageSpeed = 1.0
	// By default scale out a bit and rotate randomly for extra "pow" effect even if GM didn't manually enforce angle here
	inst.ImageAngle = rand.Float64() * 360
}

func (b *SmallExplosionBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	radius := getVar(inst, "explode_radius")
	if radius <= 0 {
		radius = 40 // default AoE radius
	}
	projectileHitBloons(inst, g, radius)
}

func (b *SmallExplosionBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	// Both idx=0 and idx=1 are used as destroy signals by different callers.
	if idx == 0 || idx == 1 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}
