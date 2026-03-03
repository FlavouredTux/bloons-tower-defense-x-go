package behaviors

import (
	"fmt"
	"math"
	"strings"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

// tower system: Buy panels, Block placement, Dart Monkey base tower

// tower costs for placement (towerselect → cost)
var towerCosts = map[float64]float64{
	1: 200, 2: 370, 3: 380, 4: 420, 5: 590, 6: 650,
	8: 450, 9: 450, 10: 510, 12: 470, 13: 810, 14: 600,
	15: 720, 16: 450, 17: 900, 18: 1200, 19: 780, 20: 950,
	21: 750, 22: 1250, 23: 1110, 24: 3000,
}

// tower select → object name mapping
var towerObjects = map[float64]string{
	1: "Dart_Monkey", 2: "Tack_Shooter", 3: "Boomerang_Thrower",
	4: "Sniper_Monkey", 5: "Ninja_Monkey", 6: "Bomb_Cannon",
	8: "Charge_Tower", 9: "Glue_Gunner_L1", 10: "Ice_Monkey",
	12: "Monkey_Engineer", 13: "Hanger_0X", 14: "Bloonchipper",
	15: "Monkey_Alchemist", 16: "Monkey_Apprentice", 17: "Banana_Tree",
	18: "Monkey_Village", 19: "Mortar_Launcher", 20: "Dartling_Gunner",
	21: "Spike_Factory", 22: "AHanger_0X", 23: "Plasma_Monkey_",
	24: "Super_Monkey",
}

// towerPanelBuy — generic buy panel for towers
// shows in sidebar, click to select for placement
// hides if the tower's lock global is 0
type towerPanelBuy struct {
	engine.DefaultBehavior
	towerSelect float64
	lockKey     string
	cost        float64
}

func (b *towerPanelBuy) Create(inst *engine.Instance, g *engine.Game) {
	// panel paging uses panelsee slots where only 1..4 are visible.
	// panel positions in rooms are laid out in 64px rows starting at y=96.
	slot := math.Round((inst.Y-96.0)/64.0) + 1.0
	inst.Vars["panelslot"] = slot
	inst.Vars["panelsee"] = slot
	g.GlobalVars["towerplace"] = 0.0
}

func (b *towerPanelBuy) Step(inst *engine.Instance, g *engine.Game) {
	// hide if tower is locked (not unlocked for this game)
	if getGlobal(g, b.lockKey) == 0 {
		inst.Visible = false
		inst.Depth = 20
		return
	}

	panelsee := getVar(inst, "panelsee")
	if panelsee < 1 || panelsee > 4 {
		inst.Visible = false
		inst.Depth = 20
		return
	}

	inst.Visible = true
	inst.Depth = -20
}

func (b *towerPanelBuy) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if !inst.Visible || getGlobal(g, b.lockKey) == 0 {
		return
	}
	money := getGlobal(g, "money")
	if money >= b.cost {
		// normal buy behavior.
		deselectAllTowers(g)

		g.GlobalVars["towerselect"] = b.towerSelect
		g.GlobalVars["towerplace"] = 1.0

		// show placement blocks
		for _, block := range g.InstanceMgr.FindByObject("Block") {
			block.Visible = true
			block.SpriteName = "Block_Spr"
		}
	}
}

func isMouseOverPanelDraw(inst *engine.Instance, g *engine.Game) bool {
	if inst == nil || g == nil || inst.SpriteName == "" {
		return false
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil {
		return false
	}
	mx := float64(g.InputMgr.MouseX)
	my := float64(g.InputMgr.MouseY)
	left := inst.X - float64(spr.XOrigin)*inst.ImageXScale + float64(spr.BBox.Left)*inst.ImageXScale
	top := inst.Y - float64(spr.YOrigin)*inst.ImageYScale + float64(spr.BBox.Top)*inst.ImageYScale
	right := inst.X - float64(spr.XOrigin)*inst.ImageXScale + float64(spr.BBox.Right+1)*inst.ImageXScale
	bottom := inst.Y - float64(spr.YOrigin)*inst.ImageYScale + float64(spr.BBox.Bottom+1)*inst.ImageYScale
	return mx >= left && mx <= right && my >= top && my <= bottom
}

func towerDisplayName(sel float64) string {
	switch sel {
	case 1:
		return "Dart Monkey"
	case 2:
		return "Tack Shooter"
	case 3:
		return "Boomerang Thrower"
	case 4:
		return "Sniper Monkey"
	case 5:
		return "Ninja Monkey"
	case 6:
		return "Bomb Cannon"
	case 7:
		return "Monkey Sub"
	case 8:
		return "Charge Tower"
	case 9:
		return "Glue Gunner"
	case 10:
		return "Ice Monkey"
	case 11:
		return "Monkey Buccaneer"
	case 12:
		return "Monkey Engineer"
	case 13:
		return "Monkey Ace"
	case 14:
		return "Bloonchipper"
	case 15:
		return "Monkey Alchemist"
	case 16:
		return "Monkey Apprentice"
	case 17:
		return "Banana Tree"
	case 18:
		return "Monkey Village"
	case 19:
		return "Mortar Launcher"
	case 20:
		return "Dartling Gunner"
	case 21:
		return "Spike Factory"
	case 22:
		return "Heli Pilot"
	case 23:
		return "Plasma Monkey"
	case 24:
		return "Super Monkey"
	default:
		if obj, ok := towerObjects[sel]; ok {
			return strings.ReplaceAll(obj, "_", " ")
		}
		return ""
	}
}

func (b *towerPanelBuy) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr != nil && len(spr.Frames) > 0 {
		frame := int(inst.ImageIndex) % len(spr.Frames)
		if frame < 0 {
			frame = 0
		}
		engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
	}
	// tooltip is drawn in a single top-layer pass from DrawHUD to avoid overlap flicker.
}

func drawHoveredTowerBuyTooltip(screen *ebiten.Image, g *engine.Game) {
	if g == nil || getGlobal(g, "towerplace") == 1 {
		return
	}

	mx := float64(g.InputMgr.MouseX)
	my := float64(g.InputMgr.MouseY)
	bestDist := math.MaxFloat64
	bestName := ""
	bestPrice := 0

	for _, inst := range g.InstanceMgr.GetAll() {
		cfg, ok := panelConfigs[inst.ObjectName]
		if !ok || !inst.Visible {
			continue
		}
		if cfg.lockKey != "" && getGlobal(g, cfg.lockKey) == 0 {
			continue
		}
		if !isMouseOverPanelDraw(inst, g) {
			continue
		}
		name := towerDisplayName(cfg.towerSelect)
		if name == "" {
			continue
		}
		d := math.Abs(inst.X-mx) + math.Abs(inst.Y-my)
		if d < bestDist {
			bestDist = d
			bestName = name
			bestPrice = int(math.Round(cfg.cost))
		}
	}

	if bestName == "" {
		return
	}

	if tip := g.AssetManager.GetSprite("Tower_Info_Panel_Spr"); tip != nil && len(tip.Frames) > 0 {
		engine.DrawSpriteExt(screen, tip.Frames[0], tip.XOrigin, tip.YOrigin, mx, my, 1, 1, 0, 1)
	}
	drawHUDTextSmall(screen, g, bestName, mx-119, my-11, hudColorBlack)
	drawHUDTextSmall(screen, g, fmt.Sprintf("$%d", bestPrice), mx-60, my+16, hudColorBlack)
}

// block — tower placement spot. Click to place selected tower.
type BlockBehavior struct {
	engine.DefaultBehavior
}

func (b *BlockBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Visible = false // hidden until tower placement mode
}

func (b *BlockBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// show/hide based on towerplace mode
	if getGlobal(g, "towerplace") == 1 {
		if inst.SpriteName == "" {
			inst.SpriteName = "Block_Spr"
		}
		inst.Visible = true
	} else {
		inst.Visible = false
	}
}

func (b *BlockBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "towerplace") != 1 {
		return
	}

	towerSel := getGlobal(g, "towerselect")
	cost, hasCost := towerCosts[towerSel]
	if !hasCost {
		return
	}

	money := getGlobal(g, "money")
	if money < cost {
		return
	}

	objName, hasObj := towerObjects[towerSel]
	if !hasObj {
		return
	}

	// check tower limit
	towerLimit := getGlobal(g, "towerlimit")
	// count placed towers (rough count of Monkey parent instances)
	towerCount := float64(countTowers(g))
	if towerCount >= towerLimit {
		return
	}

	// place tower at block center (block is 32x32 with 0,0 origin)
	tower := g.InstanceMgr.Create(objName, inst.X+16, inst.Y+16)
	if tower != nil {
		tower.Depth = -2
		// track total cash invested for sell refunds.
		tower.Vars["invested"] = cost
		fmt.Printf("Placed %s at (%.0f, %.0f) for $%.0f\n", objName, inst.X+16, inst.Y+16, cost)
	}

	// deduct cost
	g.GlobalVars["money"] = money - cost

	// play placement sound
	g.AudioMgr.Play("Tower_Place")

	// end placement mode
	g.GlobalVars["towerplace"] = 0.0
	g.GlobalVars["towerselect"] = 0.0

	// hide all blocks
	for _, block := range g.InstanceMgr.FindByObject("Block") {
		block.Visible = false
		block.SpriteName = ""
	}

	// destroy this block (tower occupies the spot)
	g.InstanceMgr.Destroy(inst.ID)
}

// countTowers counts placed tower instances
func countTowers(g *engine.Game) int {
	count := 0
	towerObjectNames := []string{
		"Dart_Monkey", "Tack_Shooter", "Boomerang_Thrower", "Sniper_Monkey",
		"Ninja_Monkey", "Bomb_Cannon", "Charge_Tower", "Glue_Gunner_L1",
		"Ice_Monkey", "Monkey_Engineer", "Hanger_0X", "Bloonchipper",
		"Monkey_Alchemist", "Monkey_Apprentice", "Banana_Tree", "Monkey_Village",
		"Mortar_Launcher", "Dartling_Gunner", "Spike_Factory", "AHanger_0X",
		"Plasma_Monkey_", "Super_Monkey",
	}
	for _, name := range towerObjectNames {
		count += g.InstanceMgr.InstanceCount(name)
	}
	return count
}

// dartMonkey — basic tower. Fires darts at nearest bloon.
type DartMonkeyBehavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
	alarmBase  float64
	camoDetect float64
	leadDetect float64
}

func (b *DartMonkeyBehavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	b.rng = 105.0
	b.alarmBase = 30.0
	b.camoDetect = 0.0
	b.leadDetect = 0.0
	inst.Vars["select"] = 0.0
	inst.Vars["range"] = b.rng
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
}

func (b *DartMonkeyBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		if getVar(inst, "stun") > 0 {
			inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
			return
		}

		// find target
		target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
		if target != nil {
			// fire dart
			dart := g.InstanceMgr.Create("Dart", inst.X, inst.Y)
			if dart != nil {
				// calculate direction to target
				dx := target.X - inst.X
				dy := target.Y - inst.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > 0 {
					speed := 15.0
					dart.HSpeed = (dx / dist) * speed
					dart.VSpeed = (dy / dist) * speed
					dart.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
					dart.ImageAngle = dart.Direction
				}
				dart.Vars["LP"] = 1.0
				dart.Vars["PP"] = 1.0 + getVar(inst, "ppbuff")
				dart.Vars["leadpop"] = b.leadDetect
				dart.Vars["camopop"] = b.camoDetect
				dart.Vars["range"] = 11.0
				dart.Alarms[0] = 11

				// face target
				inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
			}
		}

		inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
	}
}

func (b *DartMonkeyBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.Vars["range"] = b.rng
	// face nearest target
	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target != nil {
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}

	// check for linear upgrade.
	if rule, ok := linearUpgradeRules[inst.ObjectName]; ok {
		_ = applyLinearUpgrade(inst, g, rule)
	} else if rule, ok := legacyLinearUpgradeRuleFor(inst.ObjectName); ok {
		_ = applyLinearUpgrade(inst, g, rule)
	}
}

func (b *DartMonkeyBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	towerClickSelect(inst, g, 1.00)
}

// dart — projectile fired by Dart Monkey
type DartBehavior struct {
	engine.DefaultBehavior
}

func (b *DartBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	inst.Vars["LP"] = 1.0
	inst.Vars["PP"] = 1.0
	inst.Vars["leadpop"] = 0.0
	inst.Vars["camopop"] = 0.0
	inst.Vars["range"] = 11.0
}

func (b *DartBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction

	// check pierce exhaustion
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}

	// use the shared swept collision check so fast darts don't tunnel through
	projectileHitBloons(inst, g, 20)
}

func (b *DartBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		// out of range — self-destruct
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// findNearestBloon finds the nearest bloon within range
// prefers bloons furthest along the path (highest path_progress)
func findNearestBloon(inst *engine.Instance, g *engine.Game, rng float64, detectCamo bool) *engine.Instance {
	var best *engine.Instance
	bestProgress := -1.0

	bloons := g.InstanceMgr.FindByObject("Normal_Bloon_Branch")
	for _, bloon := range bloons {
		if bloon.Destroyed {
			continue
		}

		// check camo
		if !detectCamo && getVar(bloon, "camo") == 1 {
			continue
		}

		// check range
		dx := bloon.X - inst.X
		dy := bloon.Y - inst.Y
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist > rng {
			continue
		}

		// prefer bloon furthest along path
		progress := getVar(bloon, "path_progress")
		if progress > bestProgress {
			bestProgress = progress
			best = bloon
		}
	}

	return best
}

// panel buy configuration: object name → (towerselect, lock key, cost)
type panelConfig struct {
	towerSelect float64
	lockKey     string
	cost        float64
}

var panelConfigs = map[string]panelConfig{
	"Dart_Panel_buy":       {1, "DMlock", 200},
	"Tack_Panel_buy":       {2, "TSlock", 370},
	"Boomerang_Panel_buy":  {3, "BMlock", 380},
	"Sniper_Panel_buy":     {4, "SnMlock", 420},
	"Ninja_Panel_buy":      {5, "NMlock", 590},
	"Bomb_Panel_buy":       {6, "BClock", 650},
	"Sub_Panel_buy":        {7, "MSlock", 550},
	"Charge_Panel_buy":     {8, "CTlock", 450},
	"Glue_Panel_buy":       {9, "GGlock", 450},
	"Ice_Panel_buy":        {10, "IMlock", 510},
	"Buccaneer_Panel_buy":  {11, "MBlock", 700},
	"Engineer_Panel_buy":   {12, "MElock", 470},
	"Ace_Panel_buy":        {13, "MAlock", 810},
	"Chipper_Panel_buy":    {14, "BChlock", 600},
	"Alchemist_Panel_buy":  {15, "MAllock", 720},
	"Apprentice_Panel_buy": {16, "MAplock", 450},
	"Farm_Panel_buy":       {17, "BTlock", 900},
	"Village_Panel_buy":    {18, "MVlock", 1200},
	"Mortar_Panel_buy":     {19, "MLlock", 780},
	"Dartling_Panel_buy":   {20, "DGlock", 950},
	"Spike_Panel_buy":      {21, "SFlock", 750},
	"Heli_Pilot_buy":       {22, "HPlock", 1250},
	// alias from extracted data.
	"Static_Panel_buy": {16, "MAplock", 450},
	"Plasma_Panel_buy": {23, "PMlock", 1110},
	"Super_Panel_buy":  {24, "SuMlock", 3000},
}

type ScrollUpBehavior struct {
	engine.DefaultBehavior
}

type ScrollDownBehavior struct {
	engine.DefaultBehavior
}

func isTowerPanelObject(name string) bool {
	_, ok := panelConfigs[name]
	return ok
}

func towerPanelMaxScroll(g *engine.Game) int {
	maxPanelSlot := 4.0
	for _, inst := range g.InstanceMgr.GetAll() {
		if !isTowerPanelObject(inst.ObjectName) {
			continue
		}
		// respect progression lock state when deciding how far the list can scroll.
		if cfg, ok := panelConfigs[inst.ObjectName]; ok && cfg.lockKey != "" {
			if getGlobal(g, cfg.lockKey) == 0 {
				continue
			}
		}
		slot := getVar(inst, "panelslot")
		if slot == 0 {
			slot = getVar(inst, "panelsee")
		}
		if slot > maxPanelSlot {
			maxPanelSlot = slot
		}
	}
	maxScroll := int(math.Ceil(maxPanelSlot - 4.0))
	if maxScroll < 0 {
		return 0
	}
	return maxScroll
}

func towerPanelScrollOffset(g *engine.Game) int {
	minPanelSee := 1.0
	found := false
	for _, inst := range g.InstanceMgr.GetAll() {
		if !isTowerPanelObject(inst.ObjectName) {
			continue
		}
		ps := getVar(inst, "panelsee")
		if !found || ps < minPanelSee {
			minPanelSee = ps
			found = true
		}
	}
	if !found {
		return 0
	}
	offset := int(math.Round(1.0 - minPanelSee))
	if offset < 0 {
		return 0
	}
	return offset
}

func shiftTowerPanels(g *engine.Game, delta float64) {
	for _, inst := range g.InstanceMgr.GetAll() {
		if !isTowerPanelObject(inst.ObjectName) {
			continue
		}
		inst.Vars["panelsee"] = getVar(inst, "panelsee") + delta
		inst.Y += 64.0 * delta
	}
}

func mouseInSidebarScrollArea(g *engine.Game) bool {
	mx, my := g.GetMouseRoomPos()
	return mx >= 864 && mx <= 1024 && my >= 88 && my <= 392
}

func pointInArrowRect(inst *engine.Instance, g *engine.Game) bool {
	mx, my := g.GetMouseRoomPos()
	return mx >= inst.X && mx <= inst.X+64 && my >= inst.Y && my <= inst.Y+32
}

func (b *ScrollUpBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["totalpanels"] = 0.0
}

func (b *ScrollUpBehavior) scrollUp(g *engine.Game) {
	offset := towerPanelScrollOffset(g)
	if offset <= 0 {
		return
	}
	shiftTowerPanels(g, +1)
	g.GlobalVars["totalpanels"] = float64(offset - 1)
}

func (b *ScrollUpBehavior) MouseGlobalLeftPressed(inst *engine.Instance, g *engine.Game) {
	if !pointInArrowRect(inst, g) {
		return
	}
	b.scrollUp(g)
}

func (b *ScrollUpBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if !mouseInSidebarScrollArea(g) {
		return
	}
	_, wy := g.InputMgr.WheelDelta()
	if wy > 0 {
		b.scrollUp(g)
	}
}

func (b *ScrollDownBehavior) scrollDown(g *engine.Game) {
	offset := towerPanelScrollOffset(g)
	maxScroll := towerPanelMaxScroll(g)
	if offset >= maxScroll {
		return
	}
	shiftTowerPanels(g, -1)
	g.GlobalVars["totalpanels"] = float64(offset + 1)
}

func (b *ScrollDownBehavior) MouseGlobalLeftPressed(inst *engine.Instance, g *engine.Game) {
	if !pointInArrowRect(inst, g) {
		return
	}
	b.scrollDown(g)
}

func (b *ScrollDownBehavior) Step(inst *engine.Instance, g *engine.Game) {
	if !mouseInSidebarScrollArea(g) {
		return
	}
	_, wy := g.InputMgr.WheelDelta()
	if wy < 0 {
		b.scrollDown(g)
	}
}

// registerTowerBehaviors registers all tower and panel behaviors
func RegisterTowerBehaviors(im *engine.InstanceManager) {
	// tower buy panels
	for objName, cfg := range panelConfigs {
		c := cfg // capture
		im.RegisterBehavior(objName, func() engine.InstanceBehavior {
			return &towerPanelBuy{
				towerSelect: c.towerSelect,
				lockKey:     c.lockKey,
				cost:        c.cost,
			}
		})
	}

	// block placement
	im.RegisterBehavior("Block", func() engine.InstanceBehavior { return &BlockBehavior{} })
	im.RegisterBehavior("Scroll_Up", func() engine.InstanceBehavior { return &ScrollUpBehavior{} })
	im.RegisterBehavior("Scroll_Down", func() engine.InstanceBehavior { return &ScrollDownBehavior{} })

	// dart Monkey tower
	im.RegisterBehavior("Dart_Monkey", func() engine.InstanceBehavior { return &DartMonkeyBehavior{} })

	// dart projectile
	im.RegisterBehavior("Dart", func() engine.InstanceBehavior { return &DartBehavior{} })
}
