package behaviors

import (
	"math"
	"math/rand"

	"btdx/internal/engine"
)

// tack Shooter — radial burst tower with multiple upgrade paths.

type tackShooterForm struct {
	Sprite    string
	Range     float64
	AlarmBase float64
	Camo      float64
	Lead      float64
}

var tackShooterForms = map[string]tackShooterForm{
	"Tack_Shooter":      {Sprite: "Tack_Shooter_Sprite", Range: 80, AlarmBase: 51, Camo: 0, Lead: 0},
	"Faster_Shooting":   {Sprite: "Even_Mo_Tacks_Sprite", Range: 80, AlarmBase: 35, Camo: 0, Lead: 0},
	"Even_More_Tacks":   {Sprite: "Even_Mo_Tacks_Sprite", Range: 80, AlarmBase: 24, Camo: 0, Lead: 0},
	"Blade_Shooter":     {Sprite: "Blade_Shooter_Sprite", Range: 80, AlarmBase: 17, Camo: 0, Lead: 0},
	"Torque_Blades":     {Sprite: "Razor_Blade_Shooter_Sprite", Range: 82, AlarmBase: 15, Camo: 0, Lead: 0},
	"Blade_Maelstrom":   {Sprite: "Blade_Maelstrom_Sprite", Range: 85, AlarmBase: 15, Camo: 0, Lead: 0},
	"Red_Hot_Tacks":     {Sprite: "Red_Hot_Tacks_Sprite", Range: 80, AlarmBase: 21, Camo: 0, Lead: 1},
	"Flame_Jets":        {Sprite: "Flame_Jets_Sprite", Range: 86, AlarmBase: 18, Camo: 0, Lead: 1},
	"Ring_of_Fire":      {Sprite: "Ring_of_Fire_Sprite", Range: 90, AlarmBase: 15, Camo: 0, Lead: 1},
	"Tack_Sprayer":      {Sprite: "Tack_Sprayer_Sprite", Range: 82, AlarmBase: 24, Camo: 0, Lead: 0},
	"Tack_Storm":        {Sprite: "Tack_Storm_Sprite", Range: 92, AlarmBase: 2, Camo: 0, Lead: 0},
	"Tack_Typhoon":      {Sprite: "Tack_Hurricane_Sprite", Range: 96, AlarmBase: 2, Camo: 0, Lead: 0},
	"Bomb_Sprayer":      {Sprite: "Bomb_Sprayer_spr", Range: 84, AlarmBase: 21, Camo: 0, Lead: 1},
	"Fire_Crackers":     {Sprite: "Firecrackers_Spr", Range: 88, AlarmBase: 21, Camo: 0, Lead: 1},
	"Fireworks_Shooter": {Sprite: "Fire_Works_Shooter_Spr", Range: 94, AlarmBase: 15, Camo: 0, Lead: 1},
}

func tackUpgradeName(inst *engine.Instance) string {
	return upgradeName(inst, "Tack_Shooter")
}

func tackShooterFormFor(inst *engine.Instance) tackShooterForm {
	if form, ok := tackShooterForms[tackUpgradeName(inst)]; ok {
		return form
	}
	return tackShooterForms["Tack_Shooter"]
}

// TackShooterBehavior — base tower + all path upgrades
type TackShooterBehavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
	camoDetect float64
	leadDetect float64
}

func activateTowerAbility(inst *engine.Instance, g *engine.Game) bool {
	max := getVar(inst, "ability_max")
	if max <= 0 || getVar(inst, "ability") < max {
		return false
	}
	ppbuff := getVar(inst, "ppbuff")
	lead := getVar(inst, "lead_detect")
	camo := getVar(inst, "camo_detect")
	switch int(math.Round(getVar(inst, "ability_type"))) {
	case 1:
		mael := getVar(inst, "maelstrom")
		for i := 0; i < 4; i++ {
			dir := 90*float64(i) + mael
			fireInDirection(g, inst, "Torque_Blade", dir, 21, 1, 500+ppbuff, lead, camo, 30)
		}
		inst.Vars["maelstrom"] = mael + 1
	case 2:
		for i := 0; i < 18; i++ {
			dir := rand.Float64() * 360
			fireInDirection(g, inst, "Storm_Tack", dir, 21, 1, 2+ppbuff, lead, camo, 6)
		}
	case 3:
		multi := getVar(inst, "multi")
		fireworks := []struct {
			name string
			dir  float64
		}{
			{name: "Firework_I", dir: 12*multi + 90},
			{name: "Firework_II", dir: 12*multi + 270},
			{name: "Firework_III", dir: -12*multi + 180},
			{name: "Firework_IV", dir: -12*multi + 360},
		}
		for _, fw := range fireworks {
			speed := 10 + rand.Float64()*30
			life := 8 + rand.Float64()*24
			fireInDirection(g, inst, fw.name, fw.dir, speed, 8, 40+ppbuff, lead, camo, life)
		}
		inst.Vars["multi"] = multi + 1
	case 4:
		inst.Vars["hidden_timer"] = 540.0
	default:
		return false
	}
	inst.Vars["ability"] = 0.0
	g.AudioMgr.Play("Upgrade")
	return true
}

func (b *TackShooterBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	b.rng = 80.0
	b.camoDetect = 0.0
	b.leadDetect = 0.0
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["range"] = b.rng
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["ability"] = 0.0
	inst.Vars["ability_max"] = 0.0
	inst.Vars["ability_type"] = 0.0
	form := tackShooterFormFor(inst)
	inst.SpriteName = form.Sprite
	inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
}

func (b *TackShooterBehavior) refreshForm(inst *engine.Instance, g *engine.Game) tackShooterForm {
	form := tackShooterFormFor(inst)
	b.rng = form.Range
	b.camoDetect = form.Camo
	b.leadDetect = form.Lead
	inst.Vars["range"] = b.rng
	if form.Sprite != "" && g.AssetManager.GetSprite(form.Sprite) != nil {
		inst.SpriteName = form.Sprite
	}
	inst.Vars["lead_detect"] = b.leadDetect
	inst.Vars["camo_detect"] = b.camoDetect
	configureTowerAbility(inst, tackUpgradeName(inst))
	return form
}

func fireInDirection(g *engine.Game, inst *engine.Instance, projectile string, directionDeg, speed, lp, pp, lead, camo, life float64) {
	if inst == nil || life <= 0 {
		return
	}
	proj := g.InstanceMgr.Create(projectile, inst.X, inst.Y)
	if proj == nil {
		return
	}
	rad := directionDeg * math.Pi / 180.0
	proj.HSpeed = math.Cos(rad) * speed
	proj.VSpeed = -math.Sin(rad) * speed
	proj.Direction = directionDeg
	proj.ImageAngle = directionDeg
	proj.Vars["LP"] = lp
	proj.Vars["PP"] = pp
	proj.Vars["leadpop"] = lead
	proj.Vars["camopop"] = camo
	proj.Vars["range"] = life
	proj.Alarms[0] = int(math.Round(life))
}

func fireRadialBurst(g *engine.Game, inst *engine.Instance, projectile string, count int, stepDeg, speed, lp, pp, lead, camo, life float64) {
	if inst == nil || count <= 0 {
		return
	}
	for i := 0; i < count; i++ {
		fireInDirection(g, inst, projectile, float64(i)*stepDeg, speed, lp, pp, lead, camo, life)
	}
}

func (b *TackShooterBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx != 0 {
		return
	}
	form := b.refreshForm(inst, g)
	if getVar(inst, "stun") > 0 {
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
		return
	}
	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target != nil {
		ppbuff := getVar(inst, "ppbuff")
		lead := b.leadDetect
		camo := b.camoDetect
		obj := tackUpgradeName(inst)
		switch obj {
		case "Tack_Shooter", "Faster_Shooting", "Even_More_Tacks":
			fireRadialBurst(g, inst, "Tack", 8, 45, 19, 1, 1+ppbuff, lead, camo, 4)
		case "Blade_Shooter":
			fireRadialBurst(g, inst, "Blade", 8, 45, 21, 1, 3+ppbuff, lead, camo, 5)
		case "Torque_Blades":
			fireRadialBurst(g, inst, "Torque_Blade", 8, 45, 21, 1, 7+ppbuff, lead, camo, 30)
		case "Blade_Maelstrom":
			fireRadialBurst(g, inst, "Torque_Blade", 8, 45, 21, 1, 7+ppbuff, lead, camo, 30)
			setTowerAbilityCharge(inst, 1)
		case "Red_Hot_Tacks":
			fireRadialBurst(g, inst, "Red_Hot_Tack", 8, 45, 19, 1, 1+ppbuff, lead, camo, 4)
		case "Flame_Jets":
			fireRadialBurst(g, inst, "Flame_Jet", 4, 90, 2, 2, 12+ppbuff, lead, camo, 10)
		case "Ring_of_Fire":
			rof := g.InstanceMgr.Create("RoF", inst.X, inst.Y)
			if rof != nil {
				// ring of Fire damage area should track tower range, not a fixed small radius.
				hitRadius := form.Range
				if hitRadius < 30 {
					hitRadius = 30
				}
				rof.Vars["LP"] = 3.0
				rof.Vars["PP"] = 60.0 + ppbuff
				rof.Vars["leadpop"] = lead
				rof.Vars["camopop"] = camo
				rof.Vars["hit_radius"] = hitRadius
				rof.Vars["range"] = 9.0
				rof.Alarms[0] = 9
			}
		case "Tack_Sprayer":
			fireRadialBurst(g, inst, "Tack", 16, 22.5, 19, 1, 1+ppbuff, lead, camo, 4)
		case "Tack_Storm":
			waterCycle := getVar(inst, "watercycle")
			for i := 0; i < 6; i++ {
				dir := (60 * float64(i)) + (-4 * waterCycle)
				fireInDirection(g, inst, "Water_Tack", dir, 15, 1, 1+ppbuff, lead, camo, 7)
				waterCycle += 1
			}
			inst.Vars["watercycle"] = waterCycle
		case "Tack_Typhoon":
			for i := 0; i < 2; i++ {
				dir := rand.Float64() * 360
				fireInDirection(g, inst, "Storm_Tack", dir, 21, 1, 2+ppbuff, lead, camo, 6)
			}
			waterCycle := getVar(inst, "watercycle")
			for i := 0; i < 6; i++ {
				dir := (60 * float64(i)) + (-4 * waterCycle)
				fireInDirection(g, inst, "Water_Tack", dir, 16, 1, 1+ppbuff, lead, camo, 7)
				waterCycle += 1
			}
			inst.Vars["watercycle"] = waterCycle
			setTowerAbilityCharge(inst, 1)
		case "Bomb_Sprayer":
			fireRadialBurst(g, inst, "Bomb", 8, 45, 19, 1, 20+ppbuff, lead, camo, 5)
		case "Fire_Crackers":
			fireRadialBurst(g, inst, "Firecracker", 8, 45, 23, 2, 40+ppbuff, lead, camo, 18)
		case "Fireworks_Shooter":
			fireRadialBurst(g, inst, "Firecracker", 8, 45, 23, 2, 40+ppbuff, lead, camo, 18)
			setTowerAbilityCharge(inst, 1)
		default:
			fireRadialBurst(g, inst, "Tack", 8, 45, 19, 1, 1+ppbuff, lead, camo, 4)
		}
		inst.ImageXScale = 0.85
		inst.ImageYScale = 0.85
		// Tack Shooter fires radially in all directions — do not rotate toward target.
	}
	inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
}

func (b *TackShooterBehavior) Step(inst *engine.Instance, g *engine.Game) {
	form := b.refreshForm(inst, g)
	if applyPathUpgrade(inst, g) {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
		return
	}
	if upgradedTier := applyTowerUpgrade(inst, g); upgradedTier > 0 {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
	}

	// tack shooter fires in all directions, it should never rotate to face a target
	inst.ImageAngle = 0

	// restore scale
	if inst.ImageXScale < 1.0 {
		inst.ImageXScale += 0.03
		inst.ImageYScale += 0.03
		if inst.ImageXScale > 1.0 {
			inst.ImageXScale = 1.0
			inst.ImageYScale = 1.0
		}
	}
}

func (b *TackShooterBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	_ = b.refreshForm(inst, g)
	if activateTowerAbility(inst, g) {
		return
	}
	towerClickSelect(inst, g, towerSelectValue(2.0, inst))
}

// tackBehavior — short-lived projectile from Tack Shooter
type TackBehavior struct {
	engine.DefaultBehavior
}

func (b *TackBehavior) Create(inst *engine.Instance, g *engine.Game) {
	setProjDefaults(inst, 1, 1, 0, 0)
}

func (b *TackBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 16)
}

func (b *TackBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

type SpinningBladeBehavior struct {
	engine.DefaultBehavior
}

func (b *SpinningBladeBehavior) Create(inst *engine.Instance, g *engine.Game) {
	initProjDefaults(inst, 1, 7, 0, 0)
}

func (b *SpinningBladeBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle += 45
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	projectileHitBloons(inst, g, 20)
}

func (b *SpinningBladeBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

type RingOfFireBehavior struct {
	engine.DefaultBehavior
}

func (b *RingOfFireBehavior) Create(inst *engine.Instance, g *engine.Game) {
	// keep RoF cycling so it expands to full ring.
	inst.ImageSpeed = 1.0

	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = 3.0
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = 60.0
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = 1.0
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = 0.0
	}
	if _, ok := inst.Vars["hit_radius"]; !ok {
		// roF sprite is ~204px wide (about 102px radius).
		inst.Vars["hit_radius"] = 102.0
	}

	// keep RoF visual and hit area aligned to avoid "fires but no damage" desync.
	if spr := g.AssetManager.GetSprite(inst.SpriteName); spr != nil && spr.Width > 0 {
		hitRadius := getVar(inst, "hit_radius")
		if hitRadius <= 0 {
			hitRadius = 102.0
		}
		scale := (2.0 * hitRadius) / float64(spr.Width)
		if scale <= 0 {
			scale = 1.0
		}
		inst.ImageXScale = scale
		inst.ImageYScale = scale
	}
}

func (b *RingOfFireBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle += 18
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	hitRadius := getVar(inst, "hit_radius")
	if hitRadius <= 0 {
		hitRadius = 102.0
	}
	projectileHitBloons(inst, g, hitRadius)
}

func (b *RingOfFireBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}
