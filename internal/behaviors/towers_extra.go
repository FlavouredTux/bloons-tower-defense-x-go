package behaviors

import (
	"math"
	"math/rand"

	"btdx/internal/engine"
)

// additional tower and projectile behaviors
// tack Shooter, Boomerang, Sniper, Ninja, Bomb, Ice, Glue

func towerTier(inst *engine.Instance) float64 {
	return getVar(inst, "tier")
}

func towerSelectValue(base float64, inst *engine.Instance) float64 {
	if code := getVar(inst, "tower_code"); code > 0 {
		return code
	}
	return base + towerTier(inst)*0.1
}

// effectiveTowerCode returns the tower code for an object, correcting
// branch-point objects whose graph edges all carry a .22x code.  The
// upgrade-panel system expects .20 at the branch point so all three
// path panels are displayed.
func effectiveTowerCode(objectName string) (float64, bool) {
	edges := legacyUpgradeEdgesFor(objectName)
	if len(edges) == 0 {
		return 0, false
	}
	code := edges[0].TowerCode
	if len(edges) > 1 {
		frac := towerCodeFraction(code)
		if frac < 200 || (frac >= 220 && frac < 230) {
			code = math.Floor(code) + 0.20
		}
	}
	return code, true
}

// fixed tier upgrade costs (tier 0→1, 1→2) for towers that stay as a single
// behavior instance. Sourced from the upgrade graph with hardcoded fallbacks.
var tierUpgradeCosts = func() map[string][]float64 {
	fallback := map[string][]float64{
		"Tack_Shooter":      {170, 230},
		"Boomerang_Thrower": {230, 330},
		"Sniper_Monkey":     {500, 1600},
		"Ninja_Monkey":      {370, 480},
		"Bomb_Cannon":       {500, 450},
		"Ice_Monkey":        {350, 390},
		"Glue_Gunner_L1":    {300, 360},
	}

	out := make(map[string][]float64, len(fallback))
	for obj, fb := range fallback {
		if costs := legacyLinearCostsFor(obj, 2); len(costs) == 2 {
			out[obj] = costs
		} else {
			out[obj] = fb
		}
	}
	return out
}()

// applyTierUpgrade applies a linear in-place upgrade chain.
// costs[0] is tier0->tier1, costs[1] is tier1->tier2, etc.
// returns the new tier (>=1) if upgraded this frame, or 0 if no upgrade happened.
func applyTierUpgrade(inst *engine.Instance, g *engine.Game, costs []float64) int {
	if getVar(inst, "select") != 1 || getGlobal(g, "up") != 1 {
		return 0
	}
	tier := int(math.Round(towerTier(inst)))
	if tier < 0 {
		tier = 0
	}
	if tier >= len(costs) {
		g.GlobalVars["up"] = 0.0
		return 0
	}
	cost := costs[tier]
	if getGlobal(g, "money") < cost {
		g.GlobalVars["up"] = 0.0
		return 0
	}

	g.GlobalVars["money"] = getGlobal(g, "money") - cost
	inst.Vars["invested"] = getVar(inst, "invested") + cost
	newTier := tier + 1
	inst.Vars["tier"] = float64(newTier)

	// close upgrade UI.
	g.GlobalVars["upgradeselect"] = 0.0
	g.GlobalVars["up"] = 0.0
	g.GlobalVars["tower"] = 0.0
	for _, sign := range g.InstanceMgr.FindByObject("Upgrade_Sign") {
		g.InstanceMgr.Destroy(sign.ID)
	}
	return newTier
}

func applyTowerUpgrade(inst *engine.Instance, g *engine.Game) int {
	costs, ok := tierUpgradeCosts[inst.ObjectName]
	if !ok {
		return 0
	}
	return applyTierUpgrade(inst, g, costs)
}

func towerIDForObjectName(objectName string) int {
	for id, obj := range towerObjects {
		if obj == objectName {
			return int(id)
		}
	}
	return 0
}

func pathGroupForFrac1000(frac int) int {
	digit := (frac / 10) % 10
	switch digit {
	case 1, 4:
		return 1
	case 2, 5:
		return 2
	case 3, 6:
		return 3
	default:
		return 0
	}
}

var legacyPathOverrides = map[string]map[string]int{
	// dart Monkey's first split has multiple branches sharing the same code fraction.
	"Dart_Monkey_3": {
		"Bloontonium_Darts":  1,
		"Spike_o_Pult":       2,
		"Spike_o_Pult_Plus":  2,
		"Triple_Dart_Monkey": 3,
	},
	"Even_More_Tacks": {
		"Blade_Shooter": 1,
		"Red_Hot_Tacks": 2,
		"Bomb_Sprayer":  2,
		"Tack_Sprayer":  3,
	},
	// ninja training shares a code fraction across several branches.
	// force explicit path mapping so the upgrade panels route correctly.
	"Ninja_Training": {
		"Distraction":      1,
		"Flash_Bombs":      1,
		"Mass_Distraction": 1,
		"Double_Shot":      2,
		"Bloonjitzu":       2,
		"Ninja_God":        2,
		"Sai_Ninja":        3,
		"Katana_Ninja":     3,
		"Cursed_Katana_Ninja": 3,
		"Cursed_Blade_Ninja":  3,
		"Hidden_Monkey":       3,
	},
	// glaive_Thrower branches into 3 paths (all share code 3.221).
	"Glaive_Thrower": {
		"Plasmarangs":          1,
		"Glaive_Ricochet":      2,
		"Bionic_Boomer":        3,
		"Plasmasaber_Thrower":  0, // special path, excluded
	},
	// heat_Sniper branches into 3 paths (all share code 4.221).
	"Heat_Sniper": {
		"Tactical_Shotgun":      1,
		"Deadly_Precision":      2,
		"Semi_Automatic_Rifle":  3,
		"Shotgun_Plus":          0, // special path, excluded
	},
	// missile_Launcher branches into 3 paths (+ special, all share code 6.223).
	"Missile_Launcher": {
		"Bloon_Buster_Cannon": 1,
		"Cluster_Bombs":       2,
		"Pineapple_Launcher":  3,
		"Pop_Cannon":          0, // special path, excluded
	},
}

func legacyUpgradePath(sourceObj, nextObj string) int {
	if byNext, ok := legacyPathOverrides[sourceObj]; ok {
		if p, ok := byNext[nextObj]; ok {
			return p
		}
	}
	nextCode, ok := legacyTowerCodeForObject(nextObj)
	if !ok {
		return 0
	}
	return pathGroupForFrac1000(towerCodeFraction(nextCode))
}

func cheapestUpgradeEdge(candidates []legacyUpgradeEdge) legacyUpgradeEdge {
	if len(candidates) == 0 {
		return legacyUpgradeEdge{}
	}
	best := candidates[0]
	for _, c := range candidates[1:] {
		if c.Cost < best.Cost {
			best = c
		}
	}
	return best
}

func resolvedUpgradeName(inst *engine.Instance) string {
	if s, ok := inst.Vars["legacy_object"].(string); ok && s != "" {
		return s
	}

	obj := inst.ObjectName
	steps := int(math.Round(towerTier(inst)))
	if steps < 0 {
		steps = 0
	}
	for i := 0; i < steps; i++ {
		rule, ok := legacyLinearUpgradeRuleFor(obj)
		if !ok {
			break
		}
		obj = rule.NextObject
	}
	return obj
}

// upgradeName returns the current upgrade object name for an instance,
// falling back to the given default when no upgrade has been applied.
func upgradeName(inst *engine.Instance, fallback string) string {
	obj := resolvedUpgradeName(inst)
	if obj == "" {
		return fallback
	}
	return obj
}

func inferNextCodeFromSource(sourceCode float64, pathup int) float64 {
	whole := math.Floor(sourceCode)
	srcFrac := towerCodeFraction(sourceCode)
	switch {
	case srcFrac == 0:
		return whole + 0.10
	case srcFrac == 100:
		return whole + 0.20
	case srcFrac >= 300 && srcFrac < 500:
		return whole + float64(srcFrac+100)/1000.0
	case srcFrac >= 500:
		return whole + float64(srcFrac)/1000.0
	}
	switch pathup {
	case 1:
		return whole + 0.31
	case 2:
		return whole + 0.32
	case 3:
		return whole + 0.33
	default:
		return whole
	}
}

func applyPathUpgrade(inst *engine.Instance, g *engine.Game) bool {
	if getVar(inst, "select") != 1 || getGlobal(g, "up") != 1 {
		return false
	}
	pathup := int(getGlobal(g, "pathup"))
	if pathup < 1 || pathup > 3 {
		return false
	}
	if towerTier(inst) < 2 {
		return false
	}

	sourceObj := resolvedUpgradeName(inst)
	edges := legacyUpgradeEdgesFor(sourceObj)
	if len(edges) == 0 {
		g.GlobalVars["up"] = 0.0
		return false
	}

	towerID := towerIDForObjectName(inst.ObjectName)
	val := towerCodeFraction(getGlobal(g, "tower"))
	if towerID > 0 && isPathLocked(g, towerID, pathup, val) {
		g.GlobalVars["up"] = 0.0
		return false
	}

	// collect candidates matching the requested path.
	var candidates []legacyUpgradeEdge
	for _, e := range edges {
		if legacyUpgradePath(sourceObj, e.NextObject) == pathup {
			candidates = append(candidates, e)
		}
	}

	choice := legacyUpgradeEdge{}
	matchedPath := false
	switch {
	case len(candidates) > 0:
		matchedPath = true
		choice = cheapestUpgradeEdge(candidates)
	case len(edges) == 1:
		matchedPath = true
		choice = edges[0]
	case len(edges) >= pathup:
		choice = edges[pathup-1]
	default:
		g.GlobalVars["up"] = 0.0
		return false
	}

	if getGlobal(g, "money") < choice.Cost {
		g.GlobalVars["up"] = 0.0
		return false
	}

	g.GlobalVars["money"] = getGlobal(g, "money") - choice.Cost
	inst.Vars["invested"] = getVar(inst, "invested") + choice.Cost
	inst.Vars["legacy_object"] = choice.NextObject
	inst.Vars["tier"] = towerTier(inst) + 1

	nextCode, ok := legacyTowerCodeForObject(choice.NextObject)
	if !ok {
		srcCode, srcOK := legacyTowerCodeForObject(sourceObj)
		if srcOK {
			nextCode = inferNextCodeFromSource(srcCode, pathup)
		}
	}
	if nextCode > 0 && !matchedPath && pathGroupForFrac1000(towerCodeFraction(nextCode)) != pathup {
		srcCode, srcOK := legacyTowerCodeForObject(sourceObj)
		if srcOK {
			nextCode = inferNextCodeFromSource(srcCode, pathup)
		}
	}
	if nextCode > 0 {
		inst.Vars["tower_code"] = nextCode
	}

	if spr := g.InstanceMgr.ObjectSpriteName(choice.NextObject); spr != "" {
		inst.SpriteName = spr
	}

	g.GlobalVars["upgradeselect"] = 0.0
	g.GlobalVars["up"] = 0.0
	g.GlobalVars["tower"] = 0.0
	g.GlobalVars["pathup"] = 0.0
	for _, sign := range g.InstanceMgr.FindByObject("Upgrade_Sign") {
		g.InstanceMgr.Destroy(sign.ID)
	}
	return true
}

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

// tack Shooter — base tower + all path upgrades
type TackShooterBehavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
	camoDetect float64
	leadDetect float64
}

func abilitySpecForObject(obj string) (abilityType int, abilityMax float64, ok bool) {
	switch obj {
	case "Blade_Maelstrom":
		return 1, 28, true
	case "Tack_Typhoon":
		return 2, 39, true
	case "Fireworks_Shooter":
		return 3, 47, true
	case "Hidden_Monkey":
		return 4, 37, true
	default:
		return 0, 0, false
	}
}

func setTowerAbilityCharge(inst *engine.Instance, delta float64) {
	max := getVar(inst, "ability_max")
	if max <= 0 || delta <= 0 {
		return
	}
	cur := getVar(inst, "ability") + delta
	if cur > max {
		cur = max
	}
	inst.Vars["ability"] = cur
}

func configureTowerAbility(inst *engine.Instance, obj string) {
	if abilityType, abilityMax, ok := abilitySpecForObject(obj); ok {
		inst.Vars["ability_type"] = float64(abilityType)
		inst.Vars["ability_max"] = abilityMax
		cur := getVar(inst, "ability")
		if cur < 0 {
			cur = 0
		}
		if cur > abilityMax {
			cur = abilityMax
		}
		inst.Vars["ability"] = cur
		return
	}
	inst.Vars["ability_type"] = 0.0
	inst.Vars["ability_max"] = 0.0
	inst.Vars["ability"] = 0.0
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
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
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
	inst.Vars["LP"] = 1.0
	inst.Vars["PP"] = 1.0
	inst.Vars["leadpop"] = 0.0
	inst.Vars["camopop"] = 0.0
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
	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = 1.0
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = 7.0
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = 0.0
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = 0.0
	}
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

// boomerang Thrower — fires curving boomerangs
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

// boomerangBehavior — unified projectile behavior for all boomerang variants.
// handles curving return (Boomerang/Glaive/Turbo_Glaive), straight flight
// (Plasmarang/Masterang/Ricochet/King/Lord glaives), juggle bouncing,
// and Lord_Glaive's Extra_pop aura spawning.
type BoomerangBehavior struct {
	engine.DefaultBehavior
}

func (b *BoomerangBehavior) Create(inst *engine.Instance, g *engine.Game) {
	// defaults — usually overridden by fireBoomerang.
	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = 1.0
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = 1.0
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = 0.0
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = 0.0
	}
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

// sniper Monkey — global range, furthest-bloon targeting, homing darts.
// range=2000 (whole map), fires Sniper_Dart that homes at speed 54.
// supply_Drones has ability: spawns 10 Supply_Drone objects for bonus cash.

type sniperTowerForm struct {
	Range     float64
	AlarmBase float64 // alarm[0] fire delay
	Reload    float64 // alarm[11] re-target delay
	LP        float64
	PP        float64
	Lead      float64
	Camo      float64
	Speed     float64 // projectile speed
	Stun      float64 // stun decrement per shot
	Projectile string // projectile object name
	AbilityMax float64 // 0 = no ability
}

var sniperTowerForms = map[string]sniperTowerForm{
	// base chain (tier 0→1→2)
	"Sniper_Monkey":    {Range: 2000, AlarmBase: 72, Reload: 62, LP: 3, PP: 1, Lead: 0, Camo: 0, Speed: 30, Stun: 3, Projectile: "Sniper_Dart"},
	"Full_Metal_Jacket": {Range: 2000, AlarmBase: 72, Reload: 61, LP: 5, PP: 1, Lead: 1, Camo: 0, Speed: 30, Stun: 5, Projectile: "Sniper_Dart"},
	"Heat_Sniper":       {Range: 2000, AlarmBase: 60, Reload: 47, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 7, Projectile: "Sniper_Dart"},
	// path 1 (left): Tactical Shotgun → Bloonzooka → RPG Strike
	"Tactical_Shotgun":  {Range: 2000, AlarmBase: 61, Reload: 43, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 39, Projectile: "Sniper_Dart"},
	"Bloonzooka":        {Range: 2000, AlarmBase: 61, Reload: 43, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 247, Projectile: "Sniper_Dart"},
	"RPG_Strike":        {Range: 2000, AlarmBase: 61, Reload: 40, LP: 10, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 300, Projectile: "Sniper_Dart"},
	// path 2 (middle): Deadly Precision → Brick Layer → Moab Crippler
	"Deadly_Precision":  {Range: 2000, AlarmBase: 57, Reload: 45, LP: 18, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 18, Projectile: "Sniper_Dart"},
	"Brick_Layer":       {Range: 2000, AlarmBase: 54, Reload: 42, LP: 48, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 48, Projectile: "Sniper_Dart"},
	"Moab_Crippler":     {Range: 2000, AlarmBase: 54, Reload: 39, LP: 75, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 75, Projectile: "Sniper_Dart"},
	// path 3 (right): Semi Auto → Machine Gun → Supply Drones
	"Semi_Automatic_Rifle": {Range: 2000, AlarmBase: 18, Reload: 14, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 54, Stun: 7, Projectile: "Sniper_Dart"},
	"Machine_Gun":          {Range: 2000, AlarmBase: 8, Reload: 7, LP: 8, PP: 1, Lead: 1, Camo: 1, Speed: 54, Stun: 8, Projectile: "Sniper_Dart"},
	"Supply_Drones":        {Range: 2000, AlarmBase: 8, Reload: 7, LP: 8, PP: 1, Lead: 1, Camo: 1, Speed: 54, Stun: 8, Projectile: "Sniper_Dart", AbilityMax: 48},
	// special path
	"Shotgun_Plus":     {Range: 2000, AlarmBase: 61, Reload: 43, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 39, Projectile: "Sniper_Dart"},
	"Bloonzooka_Plus":  {Range: 2000, AlarmBase: 61, Reload: 43, LP: 7, PP: 1, Lead: 1, Camo: 1, Speed: 30, Stun: 247, Projectile: "Sniper_Dart"},
	"Railgun_Tank":     {Range: 2000, AlarmBase: 54, Reload: 36, LP: 100, PP: 1, Lead: 1, Camo: 1, Speed: 54, Stun: 500, Projectile: "Sniper_Dart"},
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

// sniperDartBehavior — homing projectile (speed 54, retargets if target dies)
type SniperDartBehavior struct {
	engine.DefaultBehavior
}

func (b *SniperDartBehavior) Create(inst *engine.Instance, g *engine.Game) {
	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = 3.0
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = 1.0
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = 0.0
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = 0.0
	}
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

// ninja Monkey — fast attack, camo detection, semi-homing shurikens
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
	"Ninja_Monkey":         {Range: 120, AlarmBase: 16, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 2, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Sharp_Shurikens":      {Range: 120, AlarmBase: 16, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 4, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Ninja_Training":       {Range: 130, AlarmBase: 12, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Distraction":          {Range: 125, AlarmBase: 12, BaseProjectile: "Distraction_Shot", BaseSpeed: 24, BaseLP: 1, BasePP: 7, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 100, TrackSpeed: 21, SpinRate: 20},
	"Double_Shot":          {Range: 135, AlarmBase: 12, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Flash_Bombs":          {Range: 125, AlarmBase: 11, BaseProjectile: "Distraction_Shot", BaseSpeed: 24, BaseLP: 1, BasePP: 7, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 100, TrackSpeed: 21, SpinRate: 20},
	"Mass_Distraction":     {Range: 125, AlarmBase: 11, BaseProjectile: "Distraction_Shot", BaseSpeed: 24, BaseLP: 1, BasePP: 9, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 100, TrackSpeed: 21, SpinRate: 20},
	"Sai_Ninja":            {Range: 125, AlarmBase: 11, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Katana_Ninja":         {Range: 125, AlarmBase: 8, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Hidden_Monkey":        {Range: 125, AlarmBase: 8, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Cursed_Katana_Ninja":  {Range: 133, AlarmBase: 6, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 9, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Cursed_Blade_Ninja":   {Range: 143, AlarmBase: 6, BaseProjectile: "Shuriken", BaseSpeed: 24, BaseLP: 1, BasePP: 9, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Bloonjitzu":           {Range: 139, AlarmBase: 10, BaseProjectile: "Shuriken", BaseSpeed: 21, BaseLP: 1, BasePP: 5, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 90, TrackSpeed: 21, SpinRate: 20},
	"Ninja_God":            {Range: 144, AlarmBase: 6, BaseProjectile: "Golden_Ninja_Star", BaseSpeed: 23, BaseLP: 2, BasePP: 6, BaseLead: 0, BaseCamo: 1, BaseLife: 15, HitRadius: 18, TrackRadius: 100, TrackSpeed: 23, SpinRate: 20},
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

		// upgrades add extra projectiles or melee weapons on top.
		// sai/Katana melee weapons: slow speed (4), short life (5 steps), direction offset +96,
		// spin -36 deg/step via turn_per_step. No homing. They sweep a ~180 degree arc
		// around the tower, hitting bloons with collision each frame.
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
			// sai, speed=4, dir+=96, LP=1, PP=30, life=5, turn=-36/step, no homing
			cycle++
			if cycle >= 2 {
				fireNinjaProjectile(g, inst, target, "Sai", 4, 96, 1, 30+ppbuff, 0, 1, 5, 24, 0, 0, 0, -36)
				cycle = 0
			}
		case "Katana_Ninja", "Hidden_Monkey":
			// katana, speed=4, dir+=96, LP=2, PP=45, life=5, turn=-36/step, no homing
			cycle++
			if cycle >= 2 {
				fireNinjaProjectile(g, inst, target, "Katana", 4, 96, 2, 45+ppbuff, 0, 1, 5, 24, 0, 0, 0, -36)
				cycle = 0
			}
			if obj == "Hidden_Monkey" {
				setTowerAbilityCharge(inst, 1)
			}
		case "Cursed_Katana_Ninja":
			// cursed_katana, speed=4, dir+=96, LP=2, PP=66, life=5, turn=-36/step, no homing
			cycle++
			if cycle >= 2 {
				fireNinjaProjectile(g, inst, target, "Cursed_Katana", 4, 96, 2, 66+ppbuff, 0, 1, 5, 24, 0, 0, 0, -36)
				cycle = 0
			}
		case "Cursed_Blade_Ninja":
			// fires all three melee weapons every attack, no cycle
			// cursed_Blade: speed=4, dir+=96, LP=7, PP=99, leadpop=1 (always), life=5
			fireNinjaProjectile(g, inst, target, "Cursed_Blade", 4, 96, 7, 99+ppbuff, 1, 1, 5, 24, 0, 0, 0, -36)
			// sai: speed=3, dir+=114, LP=1, PP=33, life=6
			fireNinjaProjectile(g, inst, target, "Sai", 3, 114, 1, 33+ppbuff, 0, 1, 6, 24, 0, 0, 0, -36)
			// alt_Sai: speed=3, dir-=114, LP=1, PP=33, life=6
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

// shurikenBehavior — semi-homing projectile
type ShurikenBehavior struct {
	engine.DefaultBehavior
}

func (b *ShurikenBehavior) Create(inst *engine.Instance, g *engine.Game) {
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
		inst.Vars["camopop"] = 1.0
	}
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
	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = 1.0
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = 60.0
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = 0.0
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = 1.0
	}
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
	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = 1.0
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = 60.0
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = 0.0
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = 1.0
	}
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

// ---------- Bomb Cannon form-based tower ----------

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
	"Bloon_Buster_Cannon": {Range: 135, AlarmBase: 30, Projectile: "Bomb", Speed: 30, LP: 2, PP: 40, ExplodeR: 40, Lead: 1, Camo: 0},
	"Moab_Mauler":         {Range: 135, AlarmBase: 30, Projectile: "Bomb", Speed: 30, LP: 3, PP: 40, ExplodeR: 40, Lead: 1, Camo: 0},
	"Moab_Assassin_Cannon": {Range: 150, AlarmBase: 30, Projectile: "Bomb", Speed: 45, LP: 3, PP: 40, ExplodeR: 40, Lead: 1, Camo: 1, AbilityMax: 25},
	// path 2 (middle): Cluster_Bombs → Bloon_Impactor → Explosion_King
	"Cluster_Bombs":  {Range: 135, AlarmBase: 35, Projectile: "Bomb", Speed: 24, LP: 1, PP: 40, ExplodeR: 40, Lead: 1, Camo: 0},
	"Bloon_Impactor":  {Range: 135, AlarmBase: 33, Projectile: "Bomb", Speed: 24, LP: 1, PP: 50, ExplodeR: 40, Lead: 1, Camo: 0},
	"Explosion_King":  {Range: 135, AlarmBase: 30, Projectile: "Bomb", Speed: 24, LP: 2, PP: 60, ExplodeR: 50, Lead: 1, Camo: 0},
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

// bomb Cannon — fires bombs that explode into AoE damage
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
			bomb := g.InstanceMgr.Create("Bomb", inst.X, inst.Y)
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

// bombBehavior — explodes on contact with any bloon, or on timeout
type BombBehavior struct {
	engine.DefaultBehavior
}

func (b *BombBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["LP"] = 1.0
	inst.Vars["PP"] = 20.0
	inst.Vars["leadpop"] = 1.0
	inst.Vars["camopop"] = 0.0
}

func (b *BombBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
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
	// create explosion at bomb position
	explosion := g.InstanceMgr.Create("Small_Explosion", inst.X, inst.Y)
	if explosion != nil {
		explosion.Vars["LP"] = getVar(inst, "LP")
		explosion.Vars["PP"] = getVar(inst, "PP")
		explosion.Vars["leadpop"] = getVar(inst, "leadpop")
		explosion.Vars["camopop"] = getVar(inst, "camopop")
		if r := getVar(inst, "explode_radius"); r > 0 {
			explosion.Vars["explode_radius"] = r
		}
		explosion.ImageXScale = 1.1
		explosion.ImageYScale = 1.1
		explosion.Alarms[0] = 8
	}
	g.AudioMgr.Play("Small_Boom")
	g.InstanceMgr.Destroy(inst.ID)
}

// smallExplosionBehavior — AoE damage with high pierce
type SmallExplosionBehavior struct {
	engine.DefaultBehavior
}

func (b *SmallExplosionBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["LP"] = 1.0
	inst.Vars["PP"] = 20.0
	inst.Vars["leadpop"] = 1.0
	inst.Vars["camopop"] = 0.0
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
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// extraPopBehavior — Lord_Glaive's periodic AoE damage pulse.
// spawned by Lord_Glaive every 5 frames; hits bloons in a small radius then dies.
type ExtraPopBehavior struct {
	engine.DefaultBehavior
}

func (b *ExtraPopBehavior) Create(inst *engine.Instance, g *engine.Game) {
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

// ice Monkey — freezes nearby bloons (does NOT pop them)
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

// iceAuraBehavior — stationary freeze AoE
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

// glue Gunner — glues bloons to slow them
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

// glueGlobBehavior — applies glue slow on contact
type GlueGlobBehavior struct {
	engine.DefaultBehavior
}

func (b *GlueGlobBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["LP"] = 0.0
	inst.Vars["PP"] = 1.0
	inst.Vars["leadpop"] = 1.0
	inst.Vars["camopop"] = 0.0
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

// shared projectile collision helper

// projectileHitBloons checks for bloon collisions and applies pop damage
func projectileHitBloons(inst *engine.Instance, g *engine.Game, radius float64) {
	lp := getVar(inst, "LP")
	leadpop := getVar(inst, "leadpop")
	camopop := getVar(inst, "camopop")

	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		pp := getVar(inst, "PP")
		if pp <= 0 {
			return
		}
		// lead check
		if getVar(bloon, "lead") == 1 && leadpop == 0 {
			continue
		}
		// camo check
		if getVar(bloon, "camo") == 1 && camopop == 0 {
			continue
		}

		dx := inst.X - bloon.X
		dy := inst.Y - bloon.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < radius {
			popBloon(bloon, lp, g)
			inst.Vars["PP"] = pp - 1
			if pp-1 <= 0 {
				g.InstanceMgr.Destroy(inst.ID)
				return
			}
		}
	}
}

// findStrongestBloon targets the strongest bloon (highest layer) in range
func findStrongestBloon(inst *engine.Instance, g *engine.Game, rng float64, detectCamo bool) *engine.Instance {
	var best *engine.Instance
	bestLayer := -1.0

	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		if !detectCamo && getVar(bloon, "camo") == 1 {
			continue
		}
		dx := bloon.X - inst.X
		dy := bloon.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > rng {
			continue
		}
		layer := getVar(bloon, "bloonlayer")
		if layer > bestLayer {
			bestLayer = layer
			best = bloon
		}
	}
	return best
}

// findUngluedBloon prefers bloons that aren't already glued
func findUngluedBloon(inst *engine.Instance, g *engine.Game, rng float64, detectCamo bool) *engine.Instance {
	var bestUnglued *engine.Instance
	var bestAny *engine.Instance
	bestUngluedProgress := -1.0
	bestAnyProgress := -1.0

	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}
		if !detectCamo && getVar(bloon, "camo") == 1 {
			continue
		}
		dx := bloon.X - inst.X
		dy := bloon.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > rng {
			continue
		}
		progress := getVar(bloon, "path_progress")
		if getVar(bloon, "glue") == 0 && progress > bestUngluedProgress {
			bestUngluedProgress = progress
			bestUnglued = bloon
		}
		if progress > bestAnyProgress {
			bestAnyProgress = progress
			bestAny = bloon
		}
	}
	if bestUnglued != nil {
		return bestUnglued
	}
	return bestAny
}

// registerExtraTowerBehaviors registers all additional towers
func RegisterExtraTowerBehaviors(im *engine.InstanceManager) {
	// towers
	im.RegisterBehavior("Tack_Shooter", func() engine.InstanceBehavior { return &TackShooterBehavior{} })
	im.RegisterBehavior("Boomerang_Thrower", func() engine.InstanceBehavior { return &BoomerangThrowerBehavior{} })
	im.RegisterBehavior("Sniper_Monkey", func() engine.InstanceBehavior { return &SniperMonkeyBehavior{} })
	im.RegisterBehavior("Ninja_Monkey", func() engine.InstanceBehavior { return &NinjaMonkeyBehavior{} })
	im.RegisterBehavior("Bomb_Cannon", func() engine.InstanceBehavior { return &BombCannonBehavior{} })
	im.RegisterBehavior("Ice_Monkey", func() engine.InstanceBehavior { return &IceMonkeyBehavior{} })
	im.RegisterBehavior("Glue_Gunner_L1", func() engine.InstanceBehavior { return &GlueGunnerBehavior{} })

	// projectiles
	im.RegisterBehavior("Tack", func() engine.InstanceBehavior { return &TackBehavior{} })
	im.RegisterBehavior("Blade", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 18} })
	im.RegisterBehavior("Red_Hot_Tack", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 16} })
	im.RegisterBehavior("Water_Tack", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 20} })
	im.RegisterBehavior("Storm_Tack", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 22} })
	im.RegisterBehavior("Firecracker", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 22} })
	im.RegisterBehavior("Flame_Jet", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 26} })
	im.RegisterBehavior("Firework_I", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 24} })
	im.RegisterBehavior("Firework_II", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 24} })
	im.RegisterBehavior("Firework_III", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 24} })
	im.RegisterBehavior("Firework_IV", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 24} })
	im.RegisterBehavior("Torque_Blade", func() engine.InstanceBehavior { return &SpinningBladeBehavior{} })
	im.RegisterBehavior("RoF", func() engine.InstanceBehavior { return &RingOfFireBehavior{} })
	im.RegisterBehavior("Boomerang", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Plasmarang", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Masterang", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Ricochet_Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("King_Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Lord_Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Turbo_Glaive", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Megarang", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Extra_pop", func() engine.InstanceBehavior { return &ExtraPopBehavior{} })
	im.RegisterBehavior("Super_Glaive_Proj", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Turbo_Glaive_Proj", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Glaive_Lord_Proj", func() engine.InstanceBehavior { return &BoomerangBehavior{} })
	im.RegisterBehavior("Sniper_Dart", func() engine.InstanceBehavior { return &SniperDartBehavior{} })
	im.RegisterBehavior("Shuriken", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Distraction_Shot", func() engine.InstanceBehavior { return &DistractionShotBehavior{} })
	im.RegisterBehavior("Flash_Bomb_Proj", func() engine.InstanceBehavior { return &FlashBombProjBehavior{} })
	im.RegisterBehavior("Flash", func() engine.InstanceBehavior { return &FlashBehavior{} })
	im.RegisterBehavior("Sai", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Alt_Sai", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Katana", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Cursed_Katana", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Cursed_Blade", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Crouching_Blade", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Golden_Ninja_Star", func() engine.InstanceBehavior { return &ShurikenBehavior{} })
	im.RegisterBehavior("Bomb", func() engine.InstanceBehavior { return &BombBehavior{} })
	im.RegisterBehavior("Small_Explosion", func() engine.InstanceBehavior { return &SmallExplosionBehavior{} })
	im.RegisterBehavior("Ice_Aura", func() engine.InstanceBehavior { return &IceAuraBehavior{} })
	im.RegisterBehavior("Glue_Glob", func() engine.InstanceBehavior { return &GlueGlobBehavior{} })
}
