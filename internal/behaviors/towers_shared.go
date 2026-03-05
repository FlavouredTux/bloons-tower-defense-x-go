package behaviors

import (
	"math"

	"btdx/internal/engine"
)

// shared tower utility functions, upgrade helpers, and projectile collision logic.

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
		"Charge_Tower":      {220, 330},
		"Monkey_Sub":        {330, 530},
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

	// Sync global.tower so the upgrade panel reflects the new tier immediately.
	base := math.Floor(getGlobal(g, "tower"))
	g.GlobalVars["tower"] = base + float64(newTier)*0.1

	// keep upgrade UI open so the player can continue upgrading without reselecting.
	g.GlobalVars["up"] = 0.0
	for _, sign := range g.InstanceMgr.FindByObject("Upgrade_Sign") {
		g.InstanceMgr.Destroy(sign.ID)
	}
	g.InstanceMgr.Create("Upgrade_Sign", inst.X-16, inst.Y-16)
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
		"Distraction":         1,
		"Flash_Bombs":         1,
		"Mass_Distraction":    1,
		"Double_Shot":         2,
		"Bloonjitzu":          2,
		"Ninja_God":           2,
		"Sai_Ninja":           3,
		"Katana_Ninja":        3,
		"Cursed_Katana_Ninja": 3,
		"Cursed_Blade_Ninja":  3,
		"Hidden_Monkey":       3,
	},
	// twin_Guns (Monkey Sub tier-2) branches into 3 active paths + 1 special.
	"Twin_Guns": {
		"Airburst_Sub": 1, // left path
		"Support_Sub":  2, // middle path
		"Torpedo_Sub":  3, // right path
		"Smart_Sub":    0, // special path - excluded from normal panel
	},
	// glaive_Thrower branches into 3 paths (all share code 3.221).
	"Glaive_Thrower": {
		"Plasmarangs":         1,
		"Glaive_Ricochet":     2,
		"Bionic_Boomer":       3,
		"Plasmasaber_Thrower": 0, // special path, excluded
	},
	// heat_Sniper branches into 3 paths (all share code 4.221).
	"Heat_Sniper": {
		"Tactical_Shotgun":     1,
		"Deadly_Precision":     2,
		"Semi_Automatic_Rifle": 3,
		"Shotgun_Plus":         0, // special path, excluded
	},
	// missile_Launcher branches into 3 paths (+ special, all share code 6.223).
	"Missile_Launcher": {
		"Bloon_Buster_Cannon": 1,
		"Cluster_Bombs":       2,
		"Pineapple_Launcher":  3,
		"Pop_Cannon":          0, // special path, excluded
	},
	// powerful_Charges branches into 3 paths (+ special, all share code 8.221).
	"Powerful_Charges": {
		"Charge_Battery":     1,
		"Orbital_Discharge":  2,
		"Tesla_Coil":         3,
		"Super_Charge_Tower": 0, // special path, excluded
	},
	// banana_Plantation branches into 3 income paths + 1 secret variant.
	// All share code 17.222 so explicit overrides are required.
	"Banana_Plantation": {
		"Healthy_Bananas": 1, // left path  — cheapest income upgrade
		"Banana_Republic": 2, // middle path — standard republic upgrade
		"Passive_Income":  3, // right path  — passive money stream
		"Rubberlust_Farm": 0, // secret path — excluded from normal panel
	},
	// crows_Nest branches into 3 pirate paths + 1 secret dreadnaut path.
	"Crows_Nest": {
		"Swashbucklers":  1, // left path  — melee agents
		"Destroyer":      2, // middle path — rapid fire
		"Cannon_Ship":    3, // right path  — explosive cannons
		"Dreadnaut_Ship": 0, // secret path — accelerating fire
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

	// unlock next tier in this path so the next upgrade panel is accessible.
	// The original game uses persistent CTL/CTM/CTR counters grinded across sessions;
	// here we increment by 2 per purchase which reliably unlocks each successive tier
	// (tier-4 needs >=2, tier-5 needs >=3 so +2 per step keeps pace).
	if towerID > 0 {
		if vars, ok := towerPathProgressVars[towerID]; ok && pathup >= 1 && pathup <= 3 {
			key := vars[pathup-1]
			if key != "" {
				g.GlobalVars[key] = getGlobal(g, key) + 2
			}
		}
	}

	// Sync global.tower so the upgrade panel reflects the chosen path immediately.
	// mirrors towerSelectValue: prefer tower_code, fall back to base+tier*0.1.
	if code := getVar(inst, "tower_code"); code > 0 {
		g.GlobalVars["tower"] = code
	} else {
		base := math.Floor(getGlobal(g, "tower"))
		g.GlobalVars["tower"] = base + towerTier(inst)*0.1
	}

	// keep upgrade UI open so the player can continue upgrading without reselecting.
	g.GlobalVars["up"] = 0.0
	g.GlobalVars["pathup"] = 0.0
	for _, sign := range g.InstanceMgr.FindByObject("Upgrade_Sign") {
		g.InstanceMgr.Destroy(sign.ID)
	}
	g.InstanceMgr.Create("Upgrade_Sign", inst.X-16, inst.Y-16)
	return true
}

// ── ability helpers (used by Tack Shooter, Ninja, and other towers) ──

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

// ── shared projectile collision helpers ──

// projectileHitBloons checks for bloon collisions and applies pop damage.
// uses swept-circle check along the projectile's movement vector so fast
// projectiles can't tunnel through bloons between frames.
func projectileHitBloons(inst *engine.Instance, g *engine.Game, radius float64) {
	lp := getVar(inst, "LP")
	leadpop := getVar(inst, "leadpop")
	camopop := getVar(inst, "camopop")

	// previous position for swept collision
	prevX := inst.XPrevious
	prevY := inst.YPrevious
	curX := inst.X
	curY := inst.Y

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

		// check current position first (cheap)
		dx := curX - bloon.X
		dy := curY - bloon.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < radius {
			popBloon(bloon, lp, g)
			inst.Vars["PP"] = pp - 1
			if pp-1 <= 0 {
				g.InstanceMgr.Destroy(inst.ID)
				return
			}
			continue
		}

		// swept check: find closest point on the movement segment to the bloon
		// this catches fast projectiles that skip over bloons
		segDx := curX - prevX
		segDy := curY - prevY
		segLen2 := segDx*segDx + segDy*segDy
		if segLen2 > radius*radius { // only bother if we moved more than the hit radius
			t := ((bloon.X-prevX)*segDx + (bloon.Y-prevY)*segDy) / segLen2
			if t < 0 {
				t = 0
			} else if t > 1 {
				t = 1
			}
			closestX := prevX + t*segDx
			closestY := prevY + t*segDy
			cdx := closestX - bloon.X
			cdy := closestY - bloon.Y
			cdist := math.Sqrt(cdx*cdx + cdy*cdy)
			if cdist < radius {
				popBloon(bloon, lp, g)
				inst.Vars["PP"] = pp - 1
				if pp-1 <= 0 {
					g.InstanceMgr.Destroy(inst.ID)
					return
				}
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
