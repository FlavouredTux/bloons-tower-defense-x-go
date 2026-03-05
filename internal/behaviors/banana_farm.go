package behaviors

import (
	"math"
	"math/rand"
	"time"

	"btdx/internal/engine"
)

// == Banana Tree / Farm income tower ==
//
// GML alarm base: 16 steps at ~30fps ≈ 0.53s per cycle.
// Ebiten runs at 60fps → 32 frames per cycle.
const bananaFarmAlarmBase = 32

// bananaFarmMaxBananas returns the number of bananas (×$20 each) produced per
// cycle for a given object name.  These match the `maxbananas` variable set in
// the original GML Create events.
func bananaFarmMaxBananas(objName string) float64 {
	switch objName {
	case "Banana_Tree":
		return 3 // $60 / cycle
	case "Banana_Farm":
		return 5 // $100 / cycle
	case "Banana_Plantation":
		return 8 // $160 / cycle
	case "Banana_Republic":
		return 20 // $400 / cycle
	case "Healthy_Bananas":
		return 9 // $180 / cycle
	case "Passive_Income":
		return 11 // $220 / cycle
	case "Rubberlust_Farm":
		return 7 // $140 / cycle
	case "Banana_Factory":
		return 25 // $500 / cycle
	case "Banana_Replicator":
		return 50 // $1000 / cycle
	default:
		return 5
	}
}

// bananaFarmSprite returns the sprite name for a given banana tower object.
func bananaFarmSprite(objName string) string {
	switch objName {
	case "Banana_Farm":
		return "Banana_Farm_Spr"
	case "Banana_Plantation":
		return "Banana_Plantation_Spr"
	case "Banana_Republic":
		return "Banana_Republic_Spr"
	case "Banana_Factory":
		return "Banana_Factory_Spr"
	case "Banana_Replicator":
		return "Banana_Replicator_Spr"
	default:
		return "Banana_Farm_Spr"
	}
}

// bananaFarmTowerCode returns the global.tower value to set when a banana farm
// upgrade object is clicked for selection.  The integer part is always 17 (farm
// tower ID) and the fractional encodes the upgrade position.
var bananaFarmTowerCodes = map[string]float64{
	"Banana_Farm":       17.10,
	"Banana_Plantation": 17.20,
	"Banana_Republic":   17.32,
	"Healthy_Bananas":   17.32,
	"Passive_Income":    17.32,
	"Rubberlust_Farm":   17.32,
	"Banana_Factory":    17.42,
	"Banana_Replicator": 17.52,
}

// ─── Banana Tree (base tower, select ID 17) ───────────────────────────────

// BananaTreeBehavior is the base Banana Tree tower (17, code 17.00).
// It spawns banana pickups at the end of each wave.
// Clicking it selects the tower for the upgrade panel.
// The upgrade panel can upgrade it to Banana_Farm for $600.
type BananaTreeBehavior struct {
	engine.DefaultBehavior
	prevWavenow float64
}

func (b *BananaTreeBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Banana_Tree_Spr"
	inst.Depth = 1
	inst.Vars["select"] = 0.0
	inst.Vars["tier"] = 0.0
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["maxbananas"] = bananaFarmMaxBananas("Banana_Tree")
	inst.Vars["attackrate"] = 1.0
	inst.Vars["range"] = 110.0
	b.prevWavenow = getGlobal(g, "wavenow")
}

func (b *BananaTreeBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// Spawn bananas when the wave ends (wavenow transitions 1 → 0).
	wavenow := getGlobal(g, "wavenow")
	if b.prevWavenow == 1 && wavenow == 0 {
		spawnBananas(inst, g, int(getVar(inst, "maxbananas")))
	}
	b.prevWavenow = wavenow

	// Upgrade: Banana_Tree → Banana_Farm ($600) via legacy linear upgrade graph.
	if rule, ok := linearUpgradeRules[inst.ObjectName]; ok {
		_ = applyLinearUpgrade(inst, g, rule)
	} else if rule, ok := legacyLinearUpgradeRuleFor(inst.ObjectName); ok {
		_ = applyLinearUpgrade(inst, g, rule)
	}
}

func (b *BananaTreeBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	towerClickSelect(inst, g, towerSelectValue(17.0, inst))
}

// ─── Banana Farm (and all further banana upgrade objects) ─────────────────

// BananaFarmBehavior is shared by Banana_Farm and all its further upgrades
// (Banana_Plantation, Banana_Republic, etc.).  Each upgrade object gets the
// correct sprite and income rate via bananaFarmSprite / bananaFarmMaxBananas.
type BananaFarmBehavior struct {
	engine.DefaultBehavior
	prevWavenow float64
}

func (b *BananaFarmBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = bananaFarmSprite(inst.ObjectName)
	inst.Depth = 1
	inst.Vars["select"] = 0.0
	// Banana_Plantation is the branching point in the upgrade tree.
	// applyPathUpgrade requires tier >= 2, so we start it at 2 so that the
	// three-way path panel fires immediately without extra tier increments.
	if inst.ObjectName == "Banana_Plantation" {
		inst.Vars["tier"] = 2.0
	} else {
		inst.Vars["tier"] = 0.0
	}
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Vars["maxbananas"] = bananaFarmMaxBananas(inst.ObjectName)
	inst.Vars["attackrate"] = 1.0
	inst.Vars["range"] = 110.0
	b.prevWavenow = getGlobal(g, "wavenow")
}

func (b *BananaFarmBehavior) Step(inst *engine.Instance, g *engine.Game) {
	curObj := resolvedUpgradeName(inst)

	// Passive_Income (and its successors) generate $5 every ~1.2 s during
	// active waves or when no bloons remain on screen (matching the original
	// GML Alarm[0] = 36 at 30fps → 72 frames at 60fps behaviour).
	if passiveIncomeFamily(curObj) {
		alarm, _ := inst.Vars["passive_alarm"].(float64)
		if alarm <= 0 {
			alarm = 72 // 36 GML steps × 2 for 60fps
		}
		alarm--
		inst.Vars["passive_alarm"] = alarm
		if alarm <= 0 {
			inst.Vars["passive_alarm"] = 72
			// Only drip cash during an active wave — not in the limbo between rounds.
			if getGlobal(g, "wavenow") == 1 {
				g.GlobalVars["money"] = getGlobal(g, "money") + 5.0
			}
		}
	} else {
		// All other banana farm tiers spawn pickups at the end of each wave.
		wavenow := getGlobal(g, "wavenow")
		if b.prevWavenow == 1 && wavenow == 0 {
			spawnBananas(inst, g, int(bananaFarmMaxBananas(curObj)))
		}
		b.prevWavenow = wavenow
	}

	// Try a simple linear upgrade first (e.g. Banana_Farm → Banana_Plantation).
	if rule, ok := linearUpgradeRules[inst.ObjectName]; ok {
		_ = applyLinearUpgrade(inst, g, rule)
		return
	}
	if rule, ok := legacyLinearUpgradeRuleFor(inst.ObjectName); ok {
		_ = applyLinearUpgrade(inst, g, rule)
		return
	}
	// Branching tiers (Banana_Plantation → multiple paths).
	_ = applyPathUpgrade(inst, g)
}

// passiveIncomeFamily returns true for Passive_Income and all of its upgrade
// successors.  These towers drip cash continuously rather than spawning banana
// pickups at wave-end.
func passiveIncomeFamily(objName string) bool {
	switch objName {
	case "Passive_Income", "Life_Insurance", "Long_term_Investment", "Super_Insurer":
		return true
	}
	return false
}

func (b *BananaFarmBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	code, ok := bananaFarmTowerCodes[inst.ObjectName]
	if !ok {
		code = 17.10
	}
	// Use towerSelectValue so that after a path upgrade the stored tower_code
	// var (set by applyPathUpgrade) is used instead of the base object code.
	// This prevents re-selecting an already-upgraded tower from showing the
	// Plantation branch-point panel again.
	towerClickSelect(inst, g, towerSelectValue(code, inst))
}

// ─── Spawn helper ─────────────────────────────────────────────────────────

// bananaRng is a package-level RNG for banana spawn directions.
var bananaRng = rand.New(rand.NewSource(time.Now().UnixNano()))

// spawnBananas creates count Banana pickup instances near src.
// If there are already 80+ Banana objects in the room the new bananas
// auto-collect immediately (matching the original GML "if number > 80" logic).
func spawnBananas(src *engine.Instance, g *engine.Game, count int) {
	existing := len(g.InstanceMgr.FindByObject("Banana"))
	for i := 0; i < count; i++ {
		if existing >= 80 {
			// Too many bananas on screen — auto-collect this one immediately.
			g.GlobalVars["money"] = getGlobal(g, "money") + 20.0
			continue
		}
		banana := g.InstanceMgr.Create("Banana", src.X, src.Y)
		if banana != nil {
			existing++
		}
	}
}

// ─── Banana pickup object ─────────────────────────────────────────────────

// BananaBehavior is the hoverable $20 banana pickup spawned by the farm.
// Physics is handled by the engine's built-in HSpeed/VSpeed/Friction system.
// Hover over it to collect $20.  Auto-collects after 8 s or when it leaves the
// playfield area (matching the GML "outside room" auto-collect event).
type BananaBehavior struct {
	engine.DefaultBehavior
}

func (b *BananaBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Banana_Spr"
	inst.Depth = -1 // draw in front of towers
	inst.Visible = true
	inst.ImageSpeed = 0
	inst.ImageIndex = 0

	// Random launch velocity — scaled for 60 fps (GML: speed 12–36 at 30fps).
	speed := 3.0 + bananaRng.Float64()*9.0
	angle := bananaRng.Float64() * 2 * math.Pi
	inst.HSpeed = speed * math.Cos(angle)
	inst.VSpeed = speed * math.Sin(angle)
	inst.Friction = 0.4 // engine reduces total speed by 0.4 per frame

	// Auto-collect after 8 seconds (480 frames at 60 fps).
	inst.Alarms[0] = 480
}

func (b *BananaBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// Collect on hover.
	if g.IsMouseOverInstance(inst) {
		collectBanana(inst, g)
		return
	}
	// Auto-collect when the banana drifts outside the playfield.
	if inst.X < 0 || inst.X > 863 || inst.Y < 0 || inst.Y > 575 {
		collectBanana(inst, g)
	}
}

func (b *BananaBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		collectBanana(inst, g) // timeout auto-collect
	}
}

// collectBanana gives the player $20 and removes the banana instance.
func collectBanana(inst *engine.Instance, g *engine.Game) {
	if inst.Destroyed {
		return
	}
	g.GlobalVars["money"] = getGlobal(g, "money") + 20.0
	g.InstanceMgr.Destroy(inst.ID)
}
