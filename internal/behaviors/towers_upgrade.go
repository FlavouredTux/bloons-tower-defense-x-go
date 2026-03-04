package behaviors

import (
	"fmt"
	"math"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

// tower upgrade system: panel behaviors, selection helpers,
// dart Monkey upgrade chain (tier 2, tier 3)

// tower ID → panel sprite name mapping
var towerPanelSprites = map[int]string{
	1:  "Dart_Monkey_Panel_Spr",
	2:  "Tack_Shooter_Panel_Spr",
	3:  "Boomerang_Thrower_Panel_Spr",
	4:  "Sniper_Monkey_Panel_Spr",
	5:  "Ninja_Monkey_Panel_Spr",
	6:  "Bomb_Cannon_Panel_Spr",
	7:  "Monkey_Submarine_Panel_Spr",
	8:  "Charge_Tower_Panel_Spr",
	9:  "Glue_Gunner_Panel_Spr",
	10: "Ice_Monkey_Panel_Spr",
	11: "Monkey_Buccaneer_Panel_Spr",
	12: "Monkey_Engineer_Panel_Spr",
	13: "Monkey_Ace_Panel_Spr",
	14: "Bloonchipper_Panel_Spr",
	15: "Monkey_Alchemist_Panel_Spr",
	16: "Monkey_Apprentice_Panel_Spr",
	17: "Banana_Tree_Panel_Spr",
	18: "Monkey_Village_Panel_Spr",
	19: "Mortar_Launcher_Panel_Spr",
	20: "Dartling_Gunner_Panel_Spr",
	21: "Spike_Factory_Panel_Spr",
	22: "Heli_Pilot_Panel_Spr",
	23: "Plasma_Monkey_Panel_Spr",
	24: "Super_Monkey_Panel_Spr",
}

type linearUpgradeRule struct {
	Cost       float64
	NextObject string
	TowerVal   float64 // global.tower value to assign after upgrade (0 = keep current)
}

var linearUpgradeRules = map[string]linearUpgradeRule{
	"Dart_Monkey":   {Cost: 160, NextObject: "Dart_Monkey_2", TowerVal: 1.10},
	"Dart_Monkey_2": {Cost: 180, NextObject: "Dart_Monkey_3", TowerVal: 1.20},
}

type pathChoiceRule struct {
	Cost       float64
	TowerVal   float64
	Sprite     string
	LeadDetect float64
}

// tier-3 branch rules by tower ID, then path (1=left,2=middle,3=right).
// start with Dart; add other towers here as they get implemented.
var tier3PathRules = map[int]map[int]pathChoiceRule{
	1: {
		1: {Cost: 210, TowerVal: 1.31, Sprite: "Dart_Monkey_L3_Sprite", LeadDetect: 1},
		2: {Cost: 710, TowerVal: 1.32, Sprite: "Spike_o_Pult_Sprite", LeadDetect: 0},
		3: {Cost: 440, TowerVal: 1.33, Sprite: "Triple_Darts_Sprite", LeadDetect: 0},
	},
}

// tower ID -> {left, middle, right} upgrade progression globals.
// these vars are also used by achievements and represent per-path unlock level.
var towerPathProgressVars = map[int][3]string{
	1:  {"DML", "DMM", "DMR"},
	2:  {"TSL", "TSM", "TSR"},
	3:  {"BML", "BMM", "BMR"},
	4:  {"SnML", "SnMM", "SnMR"},
	5:  {"NML", "NMM", "NMR"},
	6:  {"BCL", "BCM", "BCR"},
	7:  {"MSL", "MSM", "MSR"},
	8:  {"CTL", "CTM", "CTR"},
	9:  {"GGL", "GGM", "GGR"},
	10: {"IML", "IMM", "IMR"},
	11: {"MBL", "MBM", "MBR"},
	12: {"MEL", "MEM", "MER"},
	13: {"MAL", "MAM", "MAR"},
	14: {"BChL", "BChM", "BChR"},
	15: {"MAlL", "MAlM", "MAlR"},
	16: {"MApL", "MApM", "MApR"},
	17: {"BTL", "BTM", "BTR"},
	18: {"MVL", "MVM", "MVR"},
	19: {"MLL", "MLM", "MLR"},
	20: {"DGL", "DGM", "DGR"},
	21: {"SFL", "SFM", "SFR"},
	22: {"HPL", "HPM", "HPR"},
	23: {"PML", "PMM", "PMR"},
	24: {"SuML", "SuMM", "SuMR"},
}

// all tower object names that can be selected
var allSelectableTowers = []string{
	"Dart_Monkey", "Dart_Monkey_2", "Dart_Monkey_3",
	"Tack_Shooter", "Boomerang_Thrower", "Sniper_Monkey",
	"Ninja_Monkey", "Bomb_Cannon", "Charge_Tower",
	"Glue_Gunner_L1", "Ice_Monkey", "Monkey_Engineer",
	"Hanger_0X", "Bloonchipper", "Monkey_Alchemist",
	"Monkey_Apprentice", "Banana_Tree", "Monkey_Village",
	"Mortar_Launcher", "Dartling_Gunner", "Spike_Factory",
	"AHanger_0X", "Plasma_Monkey_", "Super_Monkey",
	"Monkey_Sub", "Barbed_Darts_Sub", "Twin_Guns",
	"Torpedo_Sub", "Airburst_Sub", "Support_Sub", "Smart_Sub",
	"Bloontonium_Reactor", "Anti_Matter_Reactor",
}

// towerCodeFraction converts frac(global.tower) to an integer ×1000 for clean comparison
func towerCodeFraction(tower float64) int {
	frac := tower - math.Floor(tower)
	return int(math.Round(frac * 1000))
}

// towerID returns the integer part of global.tower
func towerID(tower float64) int {
	return int(math.Floor(tower))
}

// selectedTower returns the currently selected tower instance, if any.
func selectedTower(g *engine.Game) *engine.Instance {
	for _, name := range allSelectableTowers {
		for _, inst := range g.InstanceMgr.FindByObject(name) {
			if getVar(inst, "select") == 1 {
				return inst
			}
		}
	}
	return nil
}

// towerBaseCost returns the best-known buy cost for an instance object.
func towerBaseCost(objectName string) float64 {
	for towerSel, objName := range towerObjects {
		if objName == objectName {
			if cost, ok := towerCosts[towerSel]; ok {
				return cost
			}
		}
	}

	// upgraded objects currently implemented.
	switch objectName {
	case "Dart_Monkey_2":
		return 200 + 160
	case "Dart_Monkey_3":
		return 200 + 160 + 180
	}
	return 0
}

// towerInvestment returns tracked invested cash for a tower, with fallback.
func towerInvestment(inst *engine.Instance) float64 {
	if inst == nil {
		return 0
	}
	if invested := getVar(inst, "invested"); invested > 0 {
		return invested
	}
	return towerBaseCost(inst.ObjectName)
}

func pathRequiredLevelForTowerVal(val int) float64 {
	// val=200 is the first 3-way split and must stay available.
	// val=300+ corresponds to later branch tiers that are progression-gated.
	switch {
	case val >= 400:
		return 3
	case val >= 300:
		return 2
	default:
		return 0
	}
}

// towerPathProgress returns progression level for tower/path.
// path: 1=left, 2=middle, 3=right.
func towerPathProgress(g *engine.Game, towerID, path int) float64 {
	vars, ok := towerPathProgressVars[towerID]
	if !ok {
		return 999 // unknown tower: never lock.
	}
	if path < 1 || path > 3 {
		return 999
	}
	key := vars[path-1]
	if key == "" {
		return 999
	}
	return getGlobal(g, key)
}

func isPathLocked(g *engine.Game, towerID, path, val int) bool {
	required := pathRequiredLevelForTowerVal(val)
	if required <= 0 {
		return false
	}
	return towerPathProgress(g, towerID, path) < required
}

func firstBranchChoicePath(val int) int {
	switch val {
	case 310:
		return 1
	case 320:
		return 2
	case 330:
		return 3
	default:
		return 0
	}
}

func firstBranchChoiceMismatch(val, panelPath int) bool {
	choice := firstBranchChoicePath(val)
	return choice != 0 && choice != panelPath
}

func getPathChoiceRule(towerID, path int) (pathChoiceRule, bool) {
	paths, ok := tier3PathRules[towerID]
	if !ok {
		return pathChoiceRule{}, false
	}
	rule, ok := paths[path]
	return rule, ok
}

func applyLinearUpgrade(inst *engine.Instance, g *engine.Game, rule linearUpgradeRule) bool {
	if getVar(inst, "select") != 1 || getGlobal(g, "up") != 1 {
		return false
	}
	if getGlobal(g, "money") < rule.Cost {
		g.GlobalVars["up"] = 0.0
		return true
	}

	g.GlobalVars["up"] = 0.0
	g.GlobalVars["money"] = getGlobal(g, "money") - rule.Cost

	// Remove the old selection indicator.
	for _, sign := range g.InstanceMgr.FindByObject("Upgrade_Sign") {
		g.InstanceMgr.Destroy(sign.ID)
	}

	newInst := g.InstanceMgr.Create(rule.NextObject, inst.X, inst.Y)
	if newInst != nil {
		newInst.Depth = inst.Depth
		newInst.Vars["ppbuff"] = getVar(inst, "ppbuff")
		newInst.Vars["invested"] = getVar(inst, "invested") + rule.Cost

		// Re-select the upgraded tower so the panel stays open.
		newInst.Vars["select"] = 1.0
		if rule.TowerVal != 0 {
			g.GlobalVars["tower"] = rule.TowerVal
		}
		g.GlobalVars["upgradeselect"] = 1.0
		g.InstanceMgr.Create("Upgrade_Sign", newInst.X-16, newInst.Y-16)
	}
	g.InstanceMgr.Destroy(inst.ID)
	return true
}

func showLockedPanel(inst *engine.Instance) {
	inst.SpriteName = "Locked"
	inst.ImageIndex = 0
	inst.ImageAlpha = 1
}

type SellBehavior struct {
	engine.DefaultBehavior
}

func (b *SellBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Depth = -101
	inst.SpriteName = "No_Panel_Spr"
	inst.ImageSpeed = 0
	inst.Visible = false
}

func (b *SellBehavior) Step(inst *engine.Instance, g *engine.Game) {
	sel := selectedTower(g)
	if sel == nil || getGlobal(g, "upgradeselect") != 1 {
		inst.Visible = false
		inst.SpriteName = "No_Panel_Spr"
		inst.Vars["sellval"] = 0.0
		return
	}

	inst.Visible = true
	inst.SpriteName = "Sell_Tower_butt_spr"
	inst.ImageIndex = 0

	invested := towerInvestment(sel)
	inst.Vars["sellval"] = invested
}

func (b *SellBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if !inst.Visible || getGlobal(g, "upgradeselect") != 1 {
		return
	}

	sel := selectedTower(g)
	if sel == nil {
		return
	}

	sellVal := towerInvestment(sel)
	refund := math.Round(sellVal * 0.8)
	if refund < 1 && sellVal > 0 {
		refund = 1
	}
	g.GlobalVars["money"] = getGlobal(g, "money") + refund

	g.InstanceMgr.Destroy(sel.ID)
	deselectAllTowers(g)
	g.GlobalVars["pathup"] = 0.0
}

func (b *SellBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if !inst.Visible {
		return
	}

	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr != nil && len(spr.Frames) > 0 {
		frame := int(inst.ImageIndex) % len(spr.Frames)
		if frame < 0 {
			frame = 0
		}
		engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
	}

	sellVal := getVar(inst, "sellval")
	refund := math.Round(sellVal * 0.8)
	if refund > 0 {
		drawHUDTextSmall(screen, g, fmt.Sprintf("%d", int(math.Round(refund))), inst.X+24, inst.Y+34, hudColorBlack)
	}
}

// deselectAllTowers — deselect all towers, reset upgrade globals
// deselects all towers
func deselectAllTowers(g *engine.Game) {
	for _, name := range allSelectableTowers {
		for _, inst := range g.InstanceMgr.FindByObject(name) {
			inst.Vars["select"] = 0.0
		}
	}
	g.GlobalVars["tower"] = 0.0
	g.GlobalVars["up"] = 0.0
	if getGlobal(g, "upgradeselect") == 1 {
		g.GlobalVars["upgradeselect"] = 0.0
		for _, sign := range g.InstanceMgr.FindByObject("Upgrade_Sign") {
			g.InstanceMgr.Destroy(sign.ID)
		}
	}
}

func cancelTowerUI(g *engine.Game) {
	deselectAllTowers(g)
	g.GlobalVars["towerselect"] = 0.0
	g.GlobalVars["towerplace"] = 0.0
	g.GlobalVars["pathup"] = 0.0
	for _, block := range g.InstanceMgr.FindByObject("Block") {
		block.Visible = false
		block.SpriteName = "sprite277"
	}
}

type CancelUpgradeBarBehavior struct {
	engine.DefaultBehavior
}

func (b *CancelUpgradeBarBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Depth = -21
	inst.ImageSpeed = 0
	inst.SpriteName = "sprite277"
}

func (b *CancelUpgradeBarBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// show while a tower is selected (upgrade UI) OR while placing a new tower.
	showUpgrade := getGlobal(g, "tower") > 0 && getGlobal(g, "towerplace") == 0 && getGlobal(g, "upgradeselect") == 1
	showPlace := getGlobal(g, "towerplace") == 1
	if showUpgrade || showPlace {
		inst.SpriteName = "Cancel_Upgrade_Bar_Spr"
		return
	}
	inst.SpriteName = "sprite277"
}

func (b *CancelUpgradeBarBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "towerplace") == 1 {
		// cancel placement mode
		cancelTowerUI(g)
		return
	}
	if getGlobal(g, "tower") <= 0 || getGlobal(g, "upgradeselect") != 1 {
		return
	}
	cancelTowerUI(g)
}

func (b *CancelUpgradeBarBehavior) KeyPress(inst *engine.Instance, g *engine.Game) {
	if !g.InputMgr.KeyPressed(ebiten.KeyX) {
		return
	}
	if getGlobal(g, "tower") <= 0 && getGlobal(g, "towerplace") == 0 && getGlobal(g, "upgradeselect") == 0 {
		return
	}
	cancelTowerUI(g)
}

func (b *CancelUpgradeBarBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := int(inst.ImageIndex) % len(spr.Frames)
	if frame < 0 {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// towerClickSelect — common click handler for all towers
// called from MouseLeftPressed. towerVal is the full global.tower value (e.g. 1.00, 1.10, 1.20)
func towerClickSelect(inst *engine.Instance, g *engine.Game, towerVal float64) {
	// during placement mode, tower clicks should not change upgrade selection.
	if getGlobal(g, "towerplace") == 1 {
		return
	}
	// ability activation: clicking a charged tower triggers ability.
	if activateTowerAbility(inst, g) {
		return
	}
	deselectAllTowers(g)
	// after deselectAllTowers, upgradeselect=0, so always enter this branch
	if getGlobal(g, "upgradeselect") == 0 {
		g.GlobalVars["tower"] = towerVal
		g.GlobalVars["upgradeselect"] = 1.0
		inst.Vars["select"] = 1.0
		// create selection indicator
		g.InstanceMgr.Create("Upgrade_Sign", inst.X-16, inst.Y-16)
	}
}

// getPanelSpriteName returns the panel sprite for current tower
func getPanelSpriteName(tower float64) string {
	id := towerID(tower)
	if spr, ok := towerPanelSprites[id]; ok {
		return spr
	}
	return "No_Panel_Spr"
}

// upgrade_Panel0 — center panel at (336, 480)
// shows tower info for base/tier1. When val >= 0.2, acts as
// panelMiddle (path 2 upgrade option).
// click triggers upgrade (global.up = 1).
type UpgradePanel0Behavior struct {
	engine.DefaultBehavior
}

func (b *UpgradePanel0Behavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "Blank_Panels"
	inst.Depth = -21
}

func (b *UpgradePanel0Behavior) StepBegin(inst *engine.Instance, g *engine.Game) {
	tower := getGlobal(g, "tower")
	if tower == 0 || getGlobal(g, "upgradeselect") == 0 {
		inst.SpriteName = "Blank_Panels"
		inst.ImageAlpha = 0
		return
	}

	val := towerCodeFraction(tower)
	panelSpr := getPanelSpriteName(tower)

	if val >= 200 {
		// act as PanelMiddle — show path 2 upgrade
		inst.SpriteName = panelSpr
		inst.ImageAlpha = 1
		// frame mapping for middle panel
		switch val {
		case 200:
			inst.ImageIndex = 5
		case 320:
			inst.ImageIndex = 6
		case 420:
			inst.ImageIndex = 7
		case 222:
			inst.ImageIndex = 13
		case 350:
			inst.ImageIndex = 14
		case 450:
			inst.ImageIndex = 15
		default:
			// check hidden states
			hiddenMiddle := []int{310, 330, 410, 430, 510, 520, 530, 360, 460, 340, 440, 221, 223, 540, 550, 560}
			for _, h := range hiddenMiddle {
				if val == h {
					inst.ImageAlpha = 0
					break
				}
			}
		}

		if firstBranchChoiceMismatch(val, 2) {
			inst.ImageAlpha = 0
			return
		}

		if isPathLocked(g, towerID(tower), 2, val) {
			showLockedPanel(inst)
		}
	} else {
		// base panel mode — show tower info
		inst.SpriteName = panelSpr
		inst.ImageAlpha = 1
		switch val {
		case 0:
			inst.ImageIndex = 0
		case 100:
			inst.ImageIndex = 1
		default:
			inst.ImageIndex = 0
		}
	}
}

func (b *UpgradePanel0Behavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	tower := getGlobal(g, "tower")
	if tower <= 0 {
		return
	}
	val := towerCodeFraction(tower)
	if val >= 200 {
		// acting as PanelMiddle — path 2
		// don't allow clicks on locked panels
		if isPathLocked(g, towerID(tower), 2, val) {
			return
		}
		g.GlobalVars["up"] = 1.0
		g.GlobalVars["pathup"] = 2.0
	} else {
		g.GlobalVars["up"] = 1.0
		// execute in-place linear upgrades immediately to avoid frame-order issues
		// where click is registered but tower Step consumes no upgrade.
		if sel := selectedTower(g); sel != nil {
			_ = applyTowerUpgrade(sel, g)
		}
	}
	g.AudioMgr.Play("Upgrade")
}

func (b *UpgradePanel0Behavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.ImageAlpha <= 0 {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := int(inst.ImageIndex) % len(spr.Frames)
	if frame < 0 {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// upgrade_PanelLeft — left panel at (64, 480)
// shows path 1 upgrade option when val >= 0.2
// click: global.up = 1, global.pathup = 1
type UpgradePanelLeftBehavior struct {
	engine.DefaultBehavior
}

func (b *UpgradePanelLeftBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "No_Panel_Spr"
	inst.Depth = -21
}

func (b *UpgradePanelLeftBehavior) StepBegin(inst *engine.Instance, g *engine.Game) {
	tower := getGlobal(g, "tower")
	inst.ImageAlpha = 1

	if tower == 0 || getGlobal(g, "upgradeselect") == 0 {
		inst.SpriteName = "No_Panel_Spr"
		inst.ImageAlpha = 0
		return
	}

	val := towerCodeFraction(tower)
	if val < 200 {
		inst.SpriteName = "No_Panel_Spr"
		inst.ImageAlpha = 0
		return
	}

	panelSpr := getPanelSpriteName(tower)
	inst.SpriteName = panelSpr

	// frame mapping for left panel (path 1)
	switch val {
	case 200:
		inst.ImageIndex = 2
	case 310:
		inst.ImageIndex = 3
	case 410:
		inst.ImageIndex = 4
	case 221:
		inst.ImageIndex = 13
	case 340:
		inst.ImageIndex = 14
	case 440:
		inst.ImageIndex = 15
	default:
		hiddenLeft := []int{320, 330, 420, 430, 510, 520, 530, 350, 450, 222, 223, 540, 550, 560, 360, 460}
		for _, h := range hiddenLeft {
			if val == h {
				inst.ImageAlpha = 0
				break
			}
		}
	}

	if firstBranchChoiceMismatch(val, 1) {
		inst.ImageAlpha = 0
		return
	}

	if isPathLocked(g, towerID(tower), 1, val) {
		showLockedPanel(inst)
	}
}

func (b *UpgradePanelLeftBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	tower := getGlobal(g, "tower")
	if tower <= 0 || getGlobal(g, "upgradeselect") != 1 {
		return
	}
	if inst.ImageAlpha <= 0 {
		return // hidden panel, can't click
	}
	// don't allow clicks on locked panels
	val := towerCodeFraction(tower)
	if isPathLocked(g, towerID(tower), 1, val) {
		return
	}
	g.GlobalVars["up"] = 1.0
	g.GlobalVars["pathup"] = 1.0
	g.AudioMgr.Play("Upgrade")
}

func (b *UpgradePanelLeftBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.ImageAlpha <= 0 {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := int(inst.ImageIndex) % len(spr.Frames)
	if frame < 0 {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// upgrade_PanelRight — right panel at (608, 480)
// shows path 3 upgrade option when val >= 0.2
// click: global.up = 1, global.pathup = 3
type UpgradePanelRightBehavior struct {
	engine.DefaultBehavior
}

func (b *UpgradePanelRightBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.SpriteName = "No_Panel_Spr"
	inst.Depth = -21
}

func (b *UpgradePanelRightBehavior) StepBegin(inst *engine.Instance, g *engine.Game) {
	tower := getGlobal(g, "tower")
	inst.ImageAlpha = 1

	if tower == 0 || getGlobal(g, "upgradeselect") == 0 {
		inst.SpriteName = "No_Panel_Spr"
		inst.ImageAlpha = 0
		return
	}

	val := towerCodeFraction(tower)
	if val < 200 {
		inst.SpriteName = "No_Panel_Spr"
		inst.ImageAlpha = 0
		return
	}

	panelSpr := getPanelSpriteName(tower)
	inst.SpriteName = panelSpr

	// frame mapping for right panel (path 3)
	switch val {
	case 200:
		inst.ImageIndex = 8
	case 330:
		inst.ImageIndex = 9
	case 430:
		inst.ImageIndex = 10
	case 223:
		inst.ImageIndex = 13
	case 360:
		inst.ImageIndex = 14
	case 460:
		inst.ImageIndex = 15
	default:
		hiddenRight := []int{310, 320, 350, 450, 340, 440, 410, 420, 510, 520, 530, 222, 221, 550, 540, 560}
		for _, h := range hiddenRight {
			if val == h {
				inst.ImageAlpha = 0
				break
			}
		}
	}

	if firstBranchChoiceMismatch(val, 3) {
		inst.ImageAlpha = 0
		return
	}

	if isPathLocked(g, towerID(tower), 3, val) {
		showLockedPanel(inst)
	}
}

func (b *UpgradePanelRightBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	tower := getGlobal(g, "tower")
	if tower <= 0 || getGlobal(g, "upgradeselect") != 1 {
		return
	}
	if inst.ImageAlpha <= 0 {
		return
	}
	// don't allow clicks on locked panels
	val := towerCodeFraction(tower)
	if isPathLocked(g, towerID(tower), 3, val) {
		return
	}
	g.GlobalVars["up"] = 1.0
	g.GlobalVars["pathup"] = 3.0
	g.AudioMgr.Play("Upgrade")
}

func (b *UpgradePanelRightBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.ImageAlpha <= 0 {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := int(inst.ImageIndex) % len(spr.Frames)
	if frame < 0 {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// dart_Monkey_2 — tier 2 Dart Monkey
type DartMonkey2Behavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
	alarmBase  float64
	camoDetect float64
	leadDetect float64
}

func (b *DartMonkey2Behavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	b.rng = 120.0
	b.alarmBase = 27.0
	b.camoDetect = 0.0
	b.leadDetect = 0.0
	inst.Vars["select"] = 0.0
	inst.Vars["range"] = b.rng
	inst.Vars["ppbuff"] = 0.0
	inst.Vars["stun"] = 0.0
	inst.SpriteName = "Dart_Monkey_L2_Sprite"
	inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
}

func (b *DartMonkey2Behavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		if getVar(inst, "stun") > 0 {
			inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
			return
		}
		target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
		if target != nil {
			dart := g.InstanceMgr.Create("Dart", inst.X, inst.Y)
			if dart != nil {
				dx := target.X - inst.X
				dy := target.Y - inst.Y
				dist := math.Sqrt(dx*dx + dy*dy)
				if dist > 0 {
					speed := 16.0
					dart.HSpeed = (dx / dist) * speed
					dart.VSpeed = (dy / dist) * speed
					dart.Direction = math.Atan2(-dy, dx) * 180 / math.Pi
					dart.ImageAngle = dart.Direction
				}
				dart.Vars["LP"] = 1.0
				dart.Vars["PP"] = 2.0 + getVar(inst, "ppbuff")
				dart.Vars["leadpop"] = b.leadDetect
				dart.Vars["camopop"] = b.camoDetect
				dart.Vars["range"] = 12.0
				dart.Alarms[0] = 12
				inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
			}
		}
		inst.Alarms[0] = int(math.Round(b.alarmBase / b.attackRate))
	}
}

func (b *DartMonkey2Behavior) Step(inst *engine.Instance, g *engine.Game) {
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

func (b *DartMonkey2Behavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	towerClickSelect(inst, g, 1.10)
}

func (b *DartMonkey2Behavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := int(inst.ImageIndex) % len(spr.Frames)
	if frame < 0 {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

type dartMonkeyForm struct {
	Sprite    string
	Range     float64
	AlarmBase float64
	Camo      float64
	Lead      float64
}

var dartMonkeyForms = map[string]dartMonkeyForm{
	"Dart_Monkey_3":       {Sprite: "Dart_Monkey_L3_Sprite", Range: 120, AlarmBase: 24, Camo: 1, Lead: 0},
	"Bloontonium_Darts":   {Sprite: "Dart_Monkey_L3_Sprite", Range: 120, AlarmBase: 24, Camo: 1, Lead: 1},
	"Dart_Monkey_Gunner":  {Sprite: "Monkey_Gunner_Sprite", Range: 130, AlarmBase: 18, Camo: 1, Lead: 1},
	"Dart_Tank":           {Sprite: "Monkey_Tank_Sprite", Range: 130, AlarmBase: 18, Camo: 1, Lead: 1},
	"Spike_o_Pult":        {Sprite: "Spike_o_Pult_Sprite", Range: 127, AlarmBase: 27, Camo: 1, Lead: 0},
	"Triple_Pult":         {Sprite: "Triple_Pult_spr", Range: 129, AlarmBase: 24, Camo: 1, Lead: 0},
	"Juggernaut":          {Sprite: "Juggernaut_Sprite", Range: 145, AlarmBase: 17, Camo: 1, Lead: 1},
	"Spike_o_Pult_Plus":   {Sprite: "Spike_o_Pult_Sprite", Range: 127, AlarmBase: 21, Camo: 1, Lead: 0},
	"Spike_Assault_Rifle": {Sprite: "Spike_ball_Gun_spr", Range: 134, AlarmBase: 21, Camo: 1, Lead: 1},
	"Spike_Mini_Gun":      {Sprite: "Spike_ball_mini_gun", Range: 139, AlarmBase: 12, Camo: 1, Lead: 1},
	"Triple_Dart_Monkey":  {Sprite: "Triple_Darts_Sprite", Range: 125, AlarmBase: 24, Camo: 1, Lead: 0},
	"Dart_Forest_Ranger":  {Sprite: "Dart_Ranger_Sprite", Range: 132, AlarmBase: 24, Camo: 1, Lead: 0},
	"SMFC_Aficionado":     {Sprite: "Super_Monkey_Fan_Club_Leader_Sprite", Range: 132, AlarmBase: 11, Camo: 1, Lead: 0},
}

func dartUpgradeName(inst *engine.Instance) string {
	return upgradeName(inst, "Dart_Monkey_3")
}

func dartMonkeyFormFor(inst *engine.Instance) dartMonkeyForm {
	if form, ok := dartMonkeyForms[dartUpgradeName(inst)]; ok {
		return form
	}
	return dartMonkeyForms["Dart_Monkey_3"]
}

func fireProjectileAt(g *engine.Game, source, target *engine.Instance, objectName string, speed, lp, pp, leadpop, camopop, life, angleOffset float64) {
	if source == nil || target == nil {
		return
	}
	if life < 1 {
		life = 1
	}
	proj := g.InstanceMgr.Create(objectName, source.X, source.Y)
	if proj == nil {
		return
	}
	dx := target.X - source.X
	dy := target.Y - source.Y
	angle := math.Atan2(-dy, dx) + angleOffset
	proj.HSpeed = math.Cos(angle) * speed
	proj.VSpeed = -math.Sin(angle) * speed
	proj.Direction = angle * 180 / math.Pi
	proj.ImageAngle = proj.Direction
	proj.Vars["LP"] = lp
	proj.Vars["PP"] = pp
	proj.Vars["leadpop"] = leadpop
	proj.Vars["camopop"] = camopop
	proj.Vars["range"] = life
	proj.Alarms[0] = int(math.Round(life))
}

// dart_Monkey_3 — tier 3+ Dart Monkey with all path upgrades
type DartMonkey3Behavior struct {
	engine.DefaultBehavior
	attackRate float64
	rng        float64
	camoDetect float64
	leadDetect float64
}

func (b *DartMonkey3Behavior) refreshForm(inst *engine.Instance, g *engine.Game) dartMonkeyForm {
	form := dartMonkeyFormFor(inst)
	b.rng = form.Range
	b.camoDetect = form.Camo
	b.leadDetect = form.Lead
	inst.Vars["range"] = b.rng
	if form.Sprite != "" && g.AssetManager.GetSprite(form.Sprite) != nil {
		inst.SpriteName = form.Sprite
	}
	return form
}

func (b *DartMonkey3Behavior) Create(inst *engine.Instance, g *engine.Game) {
	b.attackRate = 1.0
	inst.Vars["select"] = 0.0
	inst.Vars["stun"] = 0.0
	if _, ok := inst.Vars["ppbuff"]; !ok {
		inst.Vars["ppbuff"] = 0.0
	}
	if _, ok := inst.Vars["legacy_object"]; !ok {
		inst.Vars["legacy_object"] = "Dart_Monkey_3"
	}
	if getVar(inst, "tower_code") <= 0 {
		inst.Vars["tower_code"] = 1.20
	}
	// tier 3 is the first branchable Dart stage.
	if getVar(inst, "tier") < 2 {
		inst.Vars["tier"] = 2.0
	}
	form := b.refreshForm(inst, g)
	inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
}

func (b *DartMonkey3Behavior) fireAtTarget(inst, target *engine.Instance, g *engine.Game) {
	ppbuff := getVar(inst, "ppbuff")
	lead := b.leadDetect
	camo := b.camoDetect

	switch dartUpgradeName(inst) {
	case "Bloontonium_Darts":
		fireProjectileAt(g, inst, target, "Bloontonium_Dart", 16, 1, 5+ppbuff, lead, camo, 20, 0)
	case "Dart_Monkey_Gunner":
		fireProjectileAt(g, inst, target, "Bloontonium_Dart", 18, 2, 5+ppbuff, lead, camo, 25, 0)
	case "Dart_Tank":
		// dart Tank combines a heavy cannon shot with a dart-gunner shot.
		fireProjectileAt(g, inst, target, "TankDart", 27, 3, 15+ppbuff, lead, camo, 25, 0)
		fireProjectileAt(g, inst, target, "Bloontonium_Dart", 18, 2, 5+ppbuff, lead, camo, 25, 0)
	case "Spike_o_Pult":
		fireProjectileAt(g, inst, target, "Spikeball", 20, 1, 21+ppbuff, lead, camo, 25, 0)
	case "Triple_Pult":
		for i := 0; i < 3; i++ {
			off := (float64(i) - 1.0) * (3.0 * math.Pi / 180.0)
			fireProjectileAt(g, inst, target, "Spikeball", 20, 1, 21+ppbuff, lead, camo, 25, off)
		}
	case "Juggernaut":
		fireProjectileAt(g, inst, target, "Juggernaut_Ball", 24, 3, 75+ppbuff, lead, camo, 25, 0)
	case "Spike_o_Pult_Plus":
		fireProjectileAt(g, inst, target, "Spike_Plus", 22, 1, 16+ppbuff, lead, camo, 25, 0)
	case "Spike_Assault_Rifle":
		speeds := []float64{11, 19, 27, 35}
		life := []float64{30, 25, 20, 20}
		for i := range speeds {
			fireProjectileAt(g, inst, target, "Spike_Plus", speeds[i], 1, 16+ppbuff, lead, camo, life[i], 0)
		}
	case "Spike_Mini_Gun":
		speeds := []float64{11, 19, 27, 35}
		life := []float64{30, 25, 20, 20}
		for i := range speeds {
			fireProjectileAt(g, inst, target, "Spike_Plus", speeds[i], 1, 16+ppbuff, lead, camo, life[i], 0)
		}
	case "Triple_Dart_Monkey":
		for i := 0; i < 3; i++ {
			off := (float64(i) - 1.0) * (14.0 * math.Pi / 180.0)
			fireProjectileAt(g, inst, target, "Dart", 16, 1, 4+ppbuff, lead, camo, 12, off)
		}
	case "Dart_Forest_Ranger":
		for i := 0; i < 3; i++ {
			off := (float64(i) - 1.0) * (8.0 * math.Pi / 180.0)
			fireProjectileAt(g, inst, target, "Dart", 17, 1, 4+ppbuff, lead, camo, 13, off)
		}
	case "SMFC_Aficionado":
		for i := 0; i < 3; i++ {
			off := (float64(i) - 1.0) * (8.0 * math.Pi / 180.0)
			fireProjectileAt(g, inst, target, "Dart", 17, 1, 4+ppbuff, lead, camo, 15, off)
		}
	default:
		fireProjectileAt(g, inst, target, "Dart", 16, 1, 4+ppbuff, lead, camo, 12, 0)
	}
}

func (b *DartMonkey3Behavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
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
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
		b.fireAtTarget(inst, target, g)
	}
	inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
}

func (b *DartMonkey3Behavior) Step(inst *engine.Instance, g *engine.Game) {
	form := b.refreshForm(inst, g)
	if applyPathUpgrade(inst, g) {
		form = b.refreshForm(inst, g)
		inst.Alarms[0] = int(math.Round(form.AlarmBase / b.attackRate))
		return
	}
	target := findNearestBloon(inst, g, b.rng, b.camoDetect == 1)
	if target != nil {
		inst.ImageAngle = math.Atan2(-(target.Y-inst.Y), target.X-inst.X) * 180 / math.Pi
	}
}

func (b *DartMonkey3Behavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	val := getVar(inst, "tower_code")
	if val <= 0 {
		val = 1.20
	}
	towerClickSelect(inst, g, val)
}

func (b *DartMonkey3Behavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := int(inst.ImageIndex) % len(spr.Frames)
	if frame < 0 {
		frame = 0
	}
	engine.DrawSpriteExt(screen, spr.Frames[frame], spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

type LinearProjectileBehavior struct {
	engine.DefaultBehavior
	hitRadius float64
}

func (b *LinearProjectileBehavior) Create(inst *engine.Instance, g *engine.Game) {
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
		inst.Vars["camopop"] = 1.0
	}
	if _, ok := inst.Vars["range"]; !ok {
		inst.Vars["range"] = 20.0
	}
}

func (b *LinearProjectileBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle = inst.Direction
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	r := b.hitRadius
	if r <= 0 {
		r = 20
	}
	projectileHitBloons(inst, g, r)
}

func (b *LinearProjectileBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

type PultBallBehavior struct {
	engine.DefaultBehavior
}

func (b *PultBallBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Vars["LP"] = 1.0
	inst.Vars["PP"] = 21.0
	inst.Vars["leadpop"] = 0.0
	inst.Vars["camopop"] = 1.0
}

func (b *PultBallBehavior) Step(inst *engine.Instance, g *engine.Game) {
	inst.ImageAngle += 12
	if getVar(inst, "PP") <= 0 {
		g.InstanceMgr.Destroy(inst.ID)
		return
	}
	// spikeball/Juggernaut balls are chunky projectiles with larger hit radius.
	projectileHitBloons(inst, g, 24)
}

func (b *PultBallBehavior) Alarm(inst *engine.Instance, idx int, g *engine.Game) {
	if idx == 0 {
		g.InstanceMgr.Destroy(inst.ID)
	}
}

// registerUpgradeBehaviors registers upgrade panels and upgraded towers
func RegisterUpgradeBehaviors(im *engine.InstanceManager) {
	// sell panel button.
	im.RegisterBehavior("Sell", func() engine.InstanceBehavior { return &SellBehavior{} })
	im.RegisterBehavior("sell_tower_butt", func() engine.InstanceBehavior { return &SellBehavior{} })

	// upgrade panels
	im.RegisterBehavior("Upgrade_Panel0", func() engine.InstanceBehavior { return &UpgradePanel0Behavior{} })
	im.RegisterBehavior("Upgrade_PanelLeft", func() engine.InstanceBehavior { return &UpgradePanelLeftBehavior{} })
	im.RegisterBehavior("Upgrade_PanelMiddle", func() engine.InstanceBehavior { return &UpgradePanel0Behavior{} }) // same as Panel0 (handles middle mode)
	im.RegisterBehavior("Upgrade_PanelRight", func() engine.InstanceBehavior { return &UpgradePanelRightBehavior{} })
	im.RegisterBehavior("X_sign_bar", func() engine.InstanceBehavior { return &CancelUpgradeBarBehavior{} })

	// dart Monkey upgrades
	im.RegisterBehavior("Dart_Monkey_2", func() engine.InstanceBehavior { return &DartMonkey2Behavior{} })
	im.RegisterBehavior("Dart_Monkey_3", func() engine.InstanceBehavior { return &DartMonkey3Behavior{} })
	im.RegisterBehavior("Bloontonium_Dart", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 20} })
	im.RegisterBehavior("Spike_Plus", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 20} })
	im.RegisterBehavior("TankDart", func() engine.InstanceBehavior { return &LinearProjectileBehavior{hitRadius: 24} })
	im.RegisterBehavior("Spikeball", func() engine.InstanceBehavior { return &PultBallBehavior{} })
	im.RegisterBehavior("Juggernaut_Ball", func() engine.InstanceBehavior { return &PultBallBehavior{} })
}
