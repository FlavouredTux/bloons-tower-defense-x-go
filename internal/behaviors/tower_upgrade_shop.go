package behaviors

import (
	"fmt"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

// ═══════════════════════════════════════════════════════════════════════
// Tower Upgrade Shop — data-driven port of all 24 *_Upg and *_Ups objects
// ═══════════════════════════════════════════════════════════════════════

// ── Upgrade rule: one purchasable upgrade slot ──────────────────────

type upgradeRule struct {
	PanelID       int     // which panel slot (1-based counter from Create)
	PathVar       string  // global var (e.g. "DML", "DMM", "DMR")
	RequiredValue float64 // path var must equal this to be purchasable
	NewValue      float64 // value set after purchase
	BPCost        float64 // cost in Bloon Points
	RankRequired  float64 // minimum rank to unlock this tier
	DrawText      string  // text drawn on panel when available (e.g. "x2")
}

// ── Tower definition: all data for one tower's shop entry ───────────

type towerUpgDef struct {
	// *_Upg bar
	UpgObjectName string // e.g. "Dart_Monkey_Upg"
	PanelVar      string // e.g. "DMpanel"
	MenuRoom      string // e.g. "DartMenu"

	// *_Ups panel
	UpsObjectName string // e.g. "Dart_Monkey_Ups"
	InstVar       string // e.g. "DMup" (not used directly—derived from PanelVar)
	PathVars      [3]string // [Left, Middle, Right] global path vars

	// All upgrade rules for this tower
	Upgrades []upgradeRule
}

// allTowerUpgDefs — the complete table of 24 tower upgrade shop definitions.
// Ported exactly from the original GML objects.
var allTowerUpgDefs = []towerUpgDef{
	// 1. Dart Monkey
	{
		UpgObjectName: "Dart_Monkey_Upg", PanelVar: "DMpanel", MenuRoom: "DartMenu",
		UpsObjectName: "Dart_Monkey_Ups", PathVars: [3]string{"DML", "DMM", "DMR"},
		Upgrades: []upgradeRule{
			{3, "DML", 1, 2, 2, 16, "x2"},
			{6, "DMM", 1, 2, 3, 16, "x3"},
			{9, "DMR", 1, 2, 2, 16, "x2"},
			{4, "DML", 2, 3, 5, 40, "x5"},
			{7, "DMM", 2, 3, 7, 40, "x7"},
			{10, "DMR", 2, 3, 7, 40, "x7"},
			{14, "DMM", 11, 12, 15, 40, "x15"},
			{15, "DMM", 12, 13, 25, 40, "x25"},
		},
	},
	// 2. Tack Shooter
	{
		UpgObjectName: "Tack_Shooter_Upg", PanelVar: "TSpanel", MenuRoom: "Tack_Menu",
		UpsObjectName: "Tack_Shooter_Ups", PathVars: [3]string{"TSL", "TSM", "TSR"},
		Upgrades: []upgradeRule{
			{3, "TSL", 1, 2, 2, 17, "x2"},
			{6, "TSM", 1, 2, 3, 17, "x3"},
			{9, "TSR", 1, 2, 3, 17, "x3"},
			{4, "TSL", 2, 3, 5, 41, "x5"},
			{7, "TSM", 2, 3, 7, 41, "x7"},
			{10, "TSR", 2, 3, 7, 41, "x7"},
			{14, "TSM", 11, 12, 15, 41, "x15"},
			{15, "TSM", 12, 13, 25, 41, "x25"},
		},
	},
	// 3. Boomerang Thrower
	{
		UpgObjectName: "Boomerang_Upg", PanelVar: "BMpanel", MenuRoom: "Boomerang_Menu",
		UpsObjectName: "Boomerang_Thrower_Ups", PathVars: [3]string{"BML", "BMM", "BMR"},
		Upgrades: []upgradeRule{
			{3, "BML", 1, 2, 3, 18, "x3"},
			{6, "BMM", 1, 2, 3, 18, "x3"},
			{9, "BMR", 1, 2, 3, 18, "x3"},
			{4, "BML", 2, 3, 7, 42, "x7"},
			{7, "BMM", 2, 3, 9, 42, "x9"},
			{10, "BMR", 2, 3, 7, 42, "x7"},
			{14, "BML", 11, 12, 15, 42, "x15"},
			{15, "BML", 12, 13, 25, 42, "x25"},
		},
	},
	// 4. Sniper Monkey
	{
		UpgObjectName: "Sniper_Upg", PanelVar: "SnMpanel", MenuRoom: "Sniper_Menu",
		UpsObjectName: "Sniper_Ups", PathVars: [3]string{"SnML", "SnMM", "SnMR"},
		Upgrades: []upgradeRule{
			{3, "SnML", 1, 2, 4, 19, "x4"},
			{6, "SnMM", 1, 2, 5, 19, "x5"},
			{9, "SnMR", 1, 2, 4, 19, "x4"},
			{4, "SnML", 2, 3, 7, 43, "x7"},
			{7, "SnMM", 2, 3, 11, 43, "x11"},
			{10, "SnMR", 2, 3, 11, 43, "x11"},
			{14, "SnML", 11, 12, 15, 43, "x15"},
			{15, "SnML", 12, 13, 25, 43, "x25"},
		},
	},
	// 5. Ninja Monkey
	{
		UpgObjectName: "Ninja_Upg", PanelVar: "NMpanel", MenuRoom: "Ninja_Menu",
		UpsObjectName: "Ninja_Monkey_Ups", PathVars: [3]string{"NML", "NMM", "NMR"},
		Upgrades: []upgradeRule{
			{3, "NML", 1, 2, 3, 20, "x3"},
			{6, "NMM", 1, 2, 4, 20, "x4"},
			{9, "NMR", 1, 2, 3, 20, "x3"},
			{4, "NML", 2, 3, 7, 44, "x7"},
			{7, "NMM", 2, 3, 9, 44, "x9"},
			{10, "NMR", 2, 3, 7, 44, "x7"},
			{14, "NMR", 11, 12, 15, 44, "x15"},
			{15, "NMR", 12, 13, 25, 44, "x25"},
		},
	},
	// 6. Bomb Cannon
	{
		UpgObjectName: "Bomb_Upg", PanelVar: "BCpanel", MenuRoom: "Bomb_Menu",
		UpsObjectName: "Bomb_Ups", PathVars: [3]string{"BCL", "BCM", "BCR"},
		Upgrades: []upgradeRule{
			{3, "BCL", 1, 2, 2, 21, "x2"},
			{6, "BCM", 1, 2, 4, 21, "x4"},
			{9, "BCR", 1, 2, 4, 21, "x4"},
			{4, "BCL", 2, 3, 7, 45, "x7"},
			{7, "BCM", 2, 3, 7, 45, "x7"},
			{10, "BCR", 2, 3, 9, 45, "x9"},
			{14, "BCR", 11, 12, 15, 45, "x15"},
			{15, "BCR", 12, 13, 25, 45, "x25"},
		},
	},
	// 7. Monkey Sub
	{
		UpgObjectName: "Sub_Upg", PanelVar: "MSpanel", MenuRoom: "Sub_Menu",
		UpsObjectName: "Submarine_Ups", PathVars: [3]string{"MSL", "MSM", "MSR"},
		Upgrades: []upgradeRule{
			{3, "MSL", 1, 2, 3, 22, "x3"},
			{6, "MSM", 1, 2, 4, 22, "x4"},
			{9, "MSR", 1, 2, 3, 22, "x3"},
			{4, "MSL", 2, 3, 9, 46, "x9"},
			{7, "MSM", 2, 3, 9, 46, "x9"},
			{10, "MSR", 2, 3, 11, 46, "x11"},
			{14, "MSL", 11, 12, 15, 46, "x15"},
			{15, "MSL", 12, 13, 25, 46, "x25"},
		},
	},
	// 8. Charge Tower
	{
		UpgObjectName: "Charge_Upg", PanelVar: "CTpanel", MenuRoom: "Charge_Menu",
		UpsObjectName: "Charge_Ups", PathVars: [3]string{"CTL", "CTM", "CTR"},
		Upgrades: []upgradeRule{
			{3, "CTL", 1, 2, 3, 23, "x3"},
			{6, "CTM", 1, 2, 4, 23, "x4"},
			{9, "CTR", 1, 2, 4, 23, "x4"},
			{4, "CTL", 2, 3, 7, 47, "x7"},
			{7, "CTM", 2, 3, 9, 47, "x9"},
			{10, "CTR", 2, 3, 9, 47, "x9"},
			{14, "CTL", 11, 12, 15, 47, "x15"},
			{15, "CTL", 12, 13, 25, 47, "x25"},
		},
	},
	// 9. Glue Gunner
	{
		UpgObjectName: "Glue_Upg", PanelVar: "GGpanel", MenuRoom: "Glue_Menu",
		UpsObjectName: "Glue_Ups", PathVars: [3]string{"GGL", "GGM", "GGR"},
		Upgrades: []upgradeRule{
			{3, "GGL", 1, 2, 4, 24, "x4"},
			{6, "GGM", 1, 2, 5, 24, "x5"},
			{9, "GGR", 1, 2, 4, 24, "x4"},
			{4, "GGL", 2, 3, 9, 48, "x9"},
			{7, "GGM", 2, 3, 11, 48, "x11"},
			{10, "GGR", 2, 3, 9, 48, "x9"},
			{14, "GGM", 11, 12, 15, 48, "x15"},
			{15, "GGM", 12, 13, 25, 48, "x25"},
		},
	},
	// 10. Ice Monkey
	{
		UpgObjectName: "Ice_Upg", PanelVar: "IMpanel", MenuRoom: "Ice_Menu",
		UpsObjectName: "Ice_Ups", PathVars: [3]string{"IML", "IMM", "IMR"},
		Upgrades: []upgradeRule{
			{3, "IML", 1, 2, 3, 25, "x3"},
			{6, "IMM", 1, 2, 5, 25, "x5"},
			{9, "IMR", 1, 2, 3, 25, "x3"},
			{4, "IML", 2, 3, 9, 49, "x9"},
			{7, "IMM", 2, 3, 13, 49, "x13"},
			{10, "IMR", 2, 3, 11, 49, "x11"},
			{14, "IMR", 11, 12, 15, 49, "x15"},
			{15, "IMR", 12, 13, 25, 49, "x25"},
		},
	},
	// 11. Monkey Buccaneer
	{
		UpgObjectName: "Buccaneer_Upg", PanelVar: "MBpanel", MenuRoom: "Bucc_Menu",
		UpsObjectName: "Buccaneer_Ups", PathVars: [3]string{"MBL", "MBM", "MBR"},
		Upgrades: []upgradeRule{
			{3, "MBL", 1, 2, 3, 26, "x3"},
			{6, "MBM", 1, 2, 4, 26, "x4"},
			{9, "MBR", 1, 2, 3, 26, "x3"},
			{4, "MBL", 2, 3, 7, 50, "x7"},
			{7, "MBM", 2, 3, 11, 50, "x11"},
			{10, "MBR", 2, 3, 9, 50, "x9"},
			{14, "MBL", 11, 12, 15, 50, "x15"},
			{15, "MBL", 12, 13, 25, 50, "x25"},
		},
	},
	// 12. Monkey Engineer
	{
		UpgObjectName: "Engineer_Upg", PanelVar: "MEpanel", MenuRoom: "Engineer_Menu",
		UpsObjectName: "Engineer_Ups", PathVars: [3]string{"MEL", "MEM", "MER"},
		Upgrades: []upgradeRule{
			{3, "MEL", 1, 2, 3, 27, "x3"},
			{6, "MEM", 1, 2, 2, 27, "x2"},
			{9, "MER", 1, 2, 4, 27, "x3"}, // draw text "x3" but cost is 4 in GML
			{4, "MEL", 2, 3, 7, 51, "x7"},
			{7, "MEM", 2, 3, 9, 51, "x9"},
			{10, "MER", 2, 3, 9, 51, "x9"},
			{14, "MEM", 11, 12, 15, 51, "x15"},
			{15, "MEM", 12, 13, 25, 51, "x25"},
		},
	},
	// 13. Monkey Ace
	{
		UpgObjectName: "Ace_Upg", PanelVar: "MApanel", MenuRoom: "Ace_Menu",
		UpsObjectName: "Ace_Ups", PathVars: [3]string{"MAL", "MAM", "MAR"},
		Upgrades: []upgradeRule{
			{3, "MAL", 1, 2, 3, 28, "x3"},
			{6, "MAM", 1, 2, 5, 28, "x5"},
			{9, "MAR", 1, 2, 4, 28, "x4"},
			{4, "MAL", 2, 3, 9, 52, "x9"},
			{7, "MAM", 2, 3, 11, 52, "x11"},
			{10, "MAR", 2, 3, 11, 52, "x11"},
			{14, "MAR", 11, 12, 15, 52, "x15"},
			{15, "MAR", 12, 13, 25, 52, "x25"},
		},
	},
	// 14. Bloonchipper
	{
		UpgObjectName: "Chipper_Upg", PanelVar: "BChpanel", MenuRoom: "Chipper_Menu",
		UpsObjectName: "Chipper_Ups", PathVars: [3]string{"BChL", "BChM", "BChR"},
		Upgrades: []upgradeRule{
			{3, "BChL", 1, 2, 4, 29, "x4"},
			{6, "BChM", 1, 2, 4, 29, "x4"},
			{9, "BChR", 1, 2, 4, 29, "x4"},
			{4, "BChL", 2, 3, 9, 53, "x9"},
			{7, "BChM", 2, 3, 9, 53, "x11"}, // draw "x11" but cost 9 in GML
			{10, "BChR", 2, 3, 9, 53, "x9"},
			{14, "BChM", 11, 12, 15, 53, "x15"},
			{15, "BChM", 12, 13, 25, 53, "x25"},
		},
	},
	// 15. Monkey Alchemist
	{
		UpgObjectName: "Alchemist_Upg", PanelVar: "MAlpanel", MenuRoom: "Alchemist_Menu",
		UpsObjectName: "Alchemist_Ups", PathVars: [3]string{"MAlL", "MAlM", "MAlR"},
		Upgrades: []upgradeRule{
			{3, "MAlL", 1, 2, 4, 30, "x4"},
			{6, "MAlM", 1, 2, 4, 30, "x4"},
			{9, "MAlR", 1, 2, 3, 30, "x3"},
			{4, "MAlL", 2, 3, 7, 54, "x7"},
			{7, "MAlM", 2, 3, 9, 54, "x9"},
			{10, "MAlR", 2, 3, 9, 54, "x9"},
			{14, "MAlR", 11, 12, 15, 54, "x15"},
			{15, "MAlR", 12, 13, 25, 54, "x25"},
		},
	},
	// 16. Monkey Apprentice
	{
		UpgObjectName: "Apprentice_Upg", PanelVar: "MAppanel", MenuRoom: "Apprentice_Menu",
		UpsObjectName: "Apprentice_Ups", PathVars: [3]string{"MApL", "MApM", "MApR"},
		Upgrades: []upgradeRule{
			{3, "MApL", 1, 2, 4, 31, "x4"},
			{6, "MApM", 1, 2, 4, 31, "x4"},
			{9, "MApR", 1, 2, 5, 31, "x5"},
			{4, "MApL", 2, 3, 11, 55, "x11"},
			{7, "MApM", 2, 3, 11, 55, "x11"},
			{10, "MApR", 2, 3, 9, 55, "x9"},
			{14, "MApM", 11, 12, 15, 55, "x15"},
			{15, "MApM", 12, 13, 25, 55, "x25"},
		},
	},
	// 17. Banana Farm
	{
		UpgObjectName: "Farm_Upg", PanelVar: "BTpanel", MenuRoom: "Farm_Menu",
		UpsObjectName: "Farm_Ups", PathVars: [3]string{"BTL", "BTM", "BTR"},
		Upgrades: []upgradeRule{
			{3, "BTL", 1, 2, 4, 33, "x4"},
			{6, "BTM", 1, 2, 5, 33, "x5"},
			{9, "BTR", 1, 2, 4, 33, "x4"},
			{4, "BTL", 2, 3, 9, 57, "x9"},
			{7, "BTM", 2, 3, 13, 57, "x13"},
			{10, "BTR", 2, 3, 9, 57, "x9"},
			{14, "BTM", 11, 12, 15, 57, "x15"},
			{15, "BTM", 12, 13, 25, 57, "x25"},
		},
	},
	// 18. Monkey Village
	{
		UpgObjectName: "Village_Upg", PanelVar: "MVpanel", MenuRoom: "Village_Menu",
		UpsObjectName: "Village_Ups", PathVars: [3]string{"MVL", "MVM", "MVR"},
		Upgrades: []upgradeRule{
			{3, "MVL", 1, 2, 4, 32, "x4"},
			{6, "MVM", 1, 2, 5, 32, "x5"},
			{9, "MVR", 1, 2, 5, 32, "x5"},
			{4, "MVL", 2, 3, 11, 56, "x11"},
			{7, "MVM", 2, 3, 11, 56, "x11"},
			{10, "MVR", 2, 3, 11, 56, "x11"},
			{14, "MVR", 11, 12, 15, 56, "x15"},
			{15, "MVR", 12, 13, 25, 56, "x25"},
		},
	},
	// 19. Mortar Launcher
	{
		UpgObjectName: "Mortar_Upg", PanelVar: "MLpanel", MenuRoom: "Mortar_Menu",
		UpsObjectName: "Mortar_Ups", PathVars: [3]string{"MLL", "MLM", "MLR"},
		Upgrades: []upgradeRule{
			{3, "MLL", 1, 2, 4, 34, "x4"},
			{6, "MLM", 1, 2, 5, 34, "x5"},
			{9, "MLR", 1, 2, 4, 34, "x4"},
			{4, "MLL", 2, 3, 9, 58, "x9"},
			{7, "MLM", 2, 3, 13, 58, "x13"},
			{10, "MLR", 2, 3, 9, 58, "x9"},
			{14, "MLR", 11, 12, 15, 58, "x15"},
			{15, "MLR", 12, 13, 25, 58, "x25"},
		},
	},
	// 20. Dartling Gunner
	{
		UpgObjectName: "Dartling_Upg", PanelVar: "DGpanel", MenuRoom: "Dartling_Menu",
		UpsObjectName: "Dartling_Ups", PathVars: [3]string{"DGL", "DGM", "DGR"},
		Upgrades: []upgradeRule{
			{3, "DGL", 1, 2, 5, 35, "x5"},
			{6, "DGM", 1, 2, 7, 35, "x7"},
			{9, "DGR", 1, 2, 5, 35, "x5"},
			{4, "DGL", 2, 3, 11, 59, "x11"},
			{7, "DGM", 2, 3, 15, 59, "x15"},
			{10, "DGR", 2, 3, 11, 59, "x11"},
			{14, "DGL", 11, 12, 15, 59, "x15"},
			{15, "DGL", 12, 13, 25, 59, "x25"},
		},
	},
	// 21. Spike Factory
	{
		UpgObjectName: "Spike_Upg", PanelVar: "SFpanel", MenuRoom: "Spike_Menu",
		UpsObjectName: "Spike_Ups", PathVars: [3]string{"SFL", "SFM", "SFR"},
		Upgrades: []upgradeRule{
			{3, "SFL", 1, 2, 5, 36, "x5"},
			{6, "SFM", 1, 2, 5, 36, "x5"},
			{9, "SFR", 1, 2, 5, 36, "x5"},
			{4, "SFL", 2, 3, 9, 60, "x9"},
			{7, "SFM", 2, 3, 13, 60, "x13"},
			{10, "SFR", 2, 3, 11, 60, "x11"},
			{14, "SFM", 11, 12, 15, 60, "x15"},
			{15, "SFM", 12, 13, 25, 60, "x25"},
		},
	},
	// 22. Heli Pilot
	{
		UpgObjectName: "Heli_Upg", PanelVar: "HPpanel", MenuRoom: "Heli_Menu",
		UpsObjectName: "Heli_Ups", PathVars: [3]string{"HPL", "HPM", "HPR"},
		Upgrades: []upgradeRule{
			{3, "HPL", 1, 2, 5, 36, "x5"},
			{6, "HPM", 1, 2, 5, 36, "x5"},
			{9, "HPR", 1, 2, 5, 36, "x5"},
			{4, "HPL", 2, 3, 9, 60, "x9"},
			{7, "HPM", 2, 3, 13, 60, "x13"},
			{10, "HPR", 2, 3, 11, 60, "x11"},
			{14, "HPL", 11, 12, 15, 60, "x15"},
			{15, "HPL", 12, 13, 25, 60, "x25"},
		},
	},
	// 23. Plasma Monkey (note: Upg object is "Plamsa_Upg" — typo preserved from original)
	{
		UpgObjectName: "Plamsa_Upg", PanelVar: "PMpanel", MenuRoom: "Plasma_Menu",
		UpsObjectName: "Plasma_Ups", PathVars: [3]string{"PML", "PMM", "PMR"},
		Upgrades: []upgradeRule{
			{3, "PML", 1, 2, 4, 38, "x4"},
			{6, "PMM", 1, 2, 5, 38, "x5"},
			{9, "PMR", 1, 2, 5, 38, "x5"},
			{4, "PML", 2, 3, 7, 62, "x7"},
			{7, "PMM", 2, 3, 13, 62, "x13"},
			{10, "PMR", 2, 3, 11, 62, "x11"},
			{14, "PML", 11, 12, 15, 62, "x15"},
			{15, "PML", 12, 13, 25, 62, "x25"},
		},
	},
	// 24. Super Monkey
	{
		UpgObjectName: "Super_Upg", PanelVar: "SuMpanel", MenuRoom: "Super_Menu",
		UpsObjectName: "Super_Ups", PathVars: [3]string{"SuML", "SuMM", "SuMR"},
		Upgrades: []upgradeRule{
			{3, "SuML", 1, 2, 6, 39, "x6"},
			{6, "SuMM", 1, 2, 8, 39, "x8"},
			{9, "SuMR", 1, 2, 7, 39, "x7"},
			{4, "SuML", 2, 3, 13, 63, "x13"},
			{7, "SuMM", 2, 3, 17, 63, "x17"},
			{10, "SuMR", 2, 3, 13, 63, "x13"},
			{14, "SuMR", 11, 12, 15, 63, "x15"},
			{15, "SuMR", 12, 13, 25, 63, "x25"},
		},
	},
}

// ═══════════════════════════════════════════════════════════════════════
// TowerUpgBar — generic behavior for all *_Upg bar objects
// ═══════════════════════════════════════════════════════════════════════
// On click: sets the panel counter to -1, then navigates to the per-tower menu room.

type TowerUpgBar struct {
	engine.DefaultBehavior
	def *towerUpgDef
}

func (b *TowerUpgBar) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars[b.def.PanelVar] = -1.0
	g.RequestRoomGoto(b.def.MenuRoom)
}

// ═══════════════════════════════════════════════════════════════════════
// TowerUpsPanel — generic behavior for all *_Ups panel objects
// ═══════════════════════════════════════════════════════════════════════
// Each room has multiple instances of *_Ups. On Create, each increments a
// global panel counter and captures its own slot index. This determines
// which upgrade this panel represents.

type TowerUpsPanel struct {
	engine.DefaultBehavior
	def *towerUpgDef
}

func (b *TowerUpsPanel) Create(inst *engine.Instance, g *engine.Game) {
	// increment the shared panel counter
	g.GlobalVars[b.def.PanelVar] = getGlobal(g, b.def.PanelVar) + 1
	// capture this instance's panel slot
	inst.Vars["panelSlot"] = getGlobal(g, b.def.PanelVar)
	inst.Depth = 1
}

// isTier2Panel returns true for panels that use frame 12 (locked) instead of frame 11 (highlighted).
// GML: tier 1 panels (3/6/9/14) show frame 11 when purchasable; tier 2 (4/7/10/15) show frame 12 when not yet purchased.
func isTier2Panel(panelID int) bool {
	return panelID == 4 || panelID == 7 || panelID == 10 || panelID == 15
}

// secretPathVarFor returns the path var that controls the secret path for this tower.
// Derived from the PanelID 14 upgrade rule.
func secretPathVarFor(def *towerUpgDef) string {
	for _, u := range def.Upgrades {
		if u.PanelID == 14 {
			return u.PathVar
		}
	}
	return ""
}

// panelFrameAndText computes the frame index and draw-text for a given panel slot.
// Returns (frameIdx, drawText). drawText is "" when no cost text should be drawn.
func panelFrameAndText(def *towerUpgDef, g *engine.Game, slot int) (int, string) {
	// Secret path button (slot 1): hidden if secret path not unlocked
	if slot == 1 {
		sv := secretPathVarFor(def)
		if sv != "" && getGlobal(g, sv) < 11 {
			return 12, ""
		}
		return 1, ""
	}

	for _, u := range def.Upgrades {
		if u.PanelID != slot {
			continue
		}
		val := getGlobal(g, u.PathVar)

		if isTier2Panel(u.PanelID) {
			// Tier 2 panels: frame 12 when not yet purchased (val < NewValue).
			// Cost text drawn when val < NewValue (both locked and purchasable states).
			if val < u.NewValue {
				return 12, u.DrawText
			}
		} else {
			// Tier 1 panels: frame 11 when purchasable (val == RequiredValue).
			// Cost text drawn only when purchasable.
			if val == u.RequiredValue {
				return 11, u.DrawText
			}
		}
	}

	// Default: show normal slot frame (already purchased or no matching rule), no text.
	return slot, ""
}

// step updates image_index based on lock state + hover depth toggle.
// GML uses MouseEnter/MouseLeave for depth; we simulate in Step since engine lacks those events.
func (b *TowerUpsPanel) Step(inst *engine.Instance, g *engine.Game) {
	// Hover depth toggle (GML: MouseEnter → depth=-1, MouseLeave → depth=1)
	if g.IsMouseOverInstance(inst) {
		inst.Depth = -1
	} else {
		inst.Depth = 1
	}

	slot := int(getVar(inst, "panelSlot"))
	frameIdx, _ := panelFrameAndText(b.def, g, slot)
	inst.ImageIndex = float64(frameIdx)
}

// mouseLeftPressed — attempt to purchase the upgrade
func (b *TowerUpsPanel) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	slot := int(getVar(inst, "panelSlot"))
	rank := getGlobal(g, "rank")
	bp := getGlobal(g, "BP")

	for _, u := range b.def.Upgrades {
		if u.PanelID != slot {
			continue
		}
		if rank < u.RankRequired {
			continue
		}
		if getGlobal(g, u.PathVar) != u.RequiredValue {
			continue
		}
		if bp < u.BPCost {
			continue
		}
		// purchase!
		g.GlobalVars[u.PathVar] = u.NewValue
		g.GlobalVars["BP"] = bp - u.BPCost
		return
	}
}

// draw — render the panel sprite at the correct frame + cost text
func (b *TowerUpsPanel) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}

	slot := int(getVar(inst, "panelSlot"))
	frameIdx, drawText := panelFrameAndText(b.def, g, slot)

	// clamp frame index
	if frameIdx < 0 {
		frameIdx = 0
	}
	if frameIdx >= len(spr.Frames) {
		frameIdx = len(spr.Frames) - 1
	}

	frame := spr.Frames[frameIdx]
	engine.DrawSpriteExt(screen, frame, spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)

	// draw BP cost text when applicable
	if drawText != "" {
		drawUpgCostText(screen, g, drawText, inst.X+120, inst.Y+40)
	}
}

// drawUpgCostText draws the BP cost string in black at the given position using BMFont.
func drawUpgCostText(screen *ebiten.Image, g *engine.Game, text string, x, y float64) {
	drawShopText(screen, g, text, x, y, [3]uint8{0, 0, 0})
}

// drawShopText draws text using the game's BMFont (matching font0 from GML).
func drawShopText(screen *ebiten.Image, g *engine.Game, text string, x, y float64, clr [3]uint8) {
	if g.BMFont != nil && len(g.BMFont.Glyphs) > 0 {
		g.BMFont.DrawText(screen, text, x, y, clr)
		return
	}
	// fallback: basic font (shouldn't happen if font is loaded)
	drawHUDTextSmall(screen, g, text, x, y, clr)
}

// ═══════════════════════════════════════════════════════════════════════
// Panel_Block — decorative background tile, purely visual
// ═══════════════════════════════════════════════════════════════════════

type PanelBlockBehavior struct {
	engine.DefaultBehavior
}

// draw Panel_Block sprite at its position
func (b *PanelBlockBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	frame := spr.Frames[0]
	engine.DrawSpriteExt(screen, frame, spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// ═══════════════════════════════════════════════════════════════════════
// BP_Display — draws BP count, monkey money, rank, bsouls, trophies
// in the Agents and Tower Upgrades rooms
// ═══════════════════════════════════════════════════════════════════════

type BPDisplayBehavior struct {
	engine.DefaultBehavior
}

func (b *BPDisplayBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	bp := int(getGlobal(g, "BP"))
	mm := int(getGlobal(g, "monkeymoney"))
	rank := int(getGlobal(g, "rank"))
	bsouls := int(getGlobal(g, "bsouls"))
	trophies := int(getGlobal(g, "trophies"))

	black := [3]uint8{0, 0, 0}
	white := [3]uint8{255, 255, 255}

	// Row of icons at y=500 with values at y=548 (matching original GML exactly)
	// BP icon at (700, 500), value at (740, 548)
	drawBPIcon(screen, g, "BP_Icon_Spr", 700, 500)
	drawShopText(screen, g, fmt.Sprintf("%d", bp), 740, 548, black)

	// MM icon at (775, 500), value at (810, 548)
	drawBPIcon(screen, g, "MM_Icon_Spr", 775, 500)
	drawShopText(screen, g, fmt.Sprintf("%d", mm), 810, 548, black)

	// Evil Bloon Soul icon at (850, 500), value at (895, 548)
	drawBPIcon(screen, g, "Evil_Bloon_Soul_Icon_Spr", 850, 500)
	drawShopText(screen, g, fmt.Sprintf("%d", bsouls), 895, 548, black)

	// Trophy icon at (920, 500), value at (970, 548)
	drawBPIcon(screen, g, "Trophy_Icon_Spr", 920, 500)
	drawShopText(screen, g, fmt.Sprintf("%d", trophies), 970, 548, black)

	// "Rank:" label at (870, 10) in white, rank value at (930, 10) in white
	drawShopText(screen, g, "Rank:", 870, 10, white)
	drawShopText(screen, g, fmt.Sprintf("%d", rank), 930, 10, white)
}

// drawBPIcon draws a sprite at the exact position (used by BP_Display)
func drawBPIcon(screen *ebiten.Image, g *engine.Game, spriteName string, x, y float64) {
	spr := g.AssetManager.GetSprite(spriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	engine.DrawSpriteExt(screen, spr.Frames[0], spr.XOrigin, spr.YOrigin,
		x, y, 1, 1, 0, 1)
}

// ═══════════════════════════════════════════════════════════════════════
// *_Challenge_goto — stub for per-tower challenge buttons
// ═══════════════════════════════════════════════════════════════════════

type ChallengeGotoBehavior struct {
	engine.DefaultBehavior
}

func (b *ChallengeGotoBehavior) Create(inst *engine.Instance, g *engine.Game) {
	inst.Visible = false
}

// ═══════════════════════════════════════════════════════════════════════
// Go_Back_to_Main_butt — navigates back to Main_Menu
// ═══════════════════════════════════════════════════════════════════════

type GoBackToMainBehavior struct {
	engine.DefaultBehavior
}

func (b *GoBackToMainBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.RequestRoomGoto("Main_Menu")
}

// ═══════════════════════════════════════════════════════════════════════
// Up_agent / Down_Agent — scroll buttons for agent shop
// ═══════════════════════════════════════════════════════════════════════

type UpAgentBehavior struct {
	engine.DefaultBehavior
}

func (b *UpAgentBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["totalpan"] = getGlobal(g, "totalpan") - 1
}

type DownAgentBehavior struct {
	engine.DefaultBehavior
}

func (b *DownAgentBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["totalpan"] = getGlobal(g, "totalpan") + 1
}

// ═══════════════════════════════════════════════════════════════════════
// Registration
// ═══════════════════════════════════════════════════════════════════════

func RegisterTowerUpgradeShopBehaviors(im *engine.InstanceManager) {
	// register all 24 tower _Upg bars and _Ups panels
	for i := range allTowerUpgDefs {
		def := &allTowerUpgDefs[i]

		// *_Upg bar
		im.RegisterBehavior(def.UpgObjectName, func() engine.InstanceBehavior {
			return &TowerUpgBar{def: def}
		})

		// *_Ups panel
		im.RegisterBehavior(def.UpsObjectName, func() engine.InstanceBehavior {
			return &TowerUpsPanel{def: def}
		})
	}

	// Panel_Block
	im.RegisterBehavior("Panel_Block", func() engine.InstanceBehavior { return &PanelBlockBehavior{} })

	// BP_Display — override the stub from mainmenu.go
	im.RegisterBehavior("BP_Display", func() engine.InstanceBehavior { return &BPDisplayBehavior{} })

	// Challenge goto buttons (stubs)
	challengeNames := []string{
		"Dart_Challenge_goto", "Tack_Challenge_goto", "Boomerang_Challenge_goto",
		"Sniper_Challenge_goto", "Ninja_Challenge_goto", "Bomb_Challenge_goto",
		"Sub_Challenge_goto", "Charge_Challenge_goto", "Glue_Challenge_goto",
		"Ice_Challenge_goto", "Buccaneer_Challenge_goto", "Engineer_Challenge_goto",
		"Ace_Challenge_goto", "Chipper_Challenge_goto", "Alchemist_Challenge_goto",
		"Apprentice_Challenge_goto", "Farm_Challenge_goto", "Village_Challenge_goto",
		"Mortar_Challenge_goto", "Dartling_Challenge_goto", "Spike_Challenge_goto",
		"Heli_Challenge_goto", "Plasma_Challenge_goto", "Super_Challenge_goto",
	}
	for _, name := range challengeNames {
		im.RegisterBehavior(name, func() engine.InstanceBehavior { return &ChallengeGotoBehavior{} })
	}

	// Go Back to Main Menu button
	im.RegisterBehavior("Go_Back_to_Main_butt", func() engine.InstanceBehavior { return &GoBackToMainBehavior{} })

	// Scroll buttons
	im.RegisterBehavior("Up_agent", func() engine.InstanceBehavior { return &UpAgentBehavior{} })
	im.RegisterBehavior("Down_Agent", func() engine.InstanceBehavior { return &DownAgentBehavior{} })
}
