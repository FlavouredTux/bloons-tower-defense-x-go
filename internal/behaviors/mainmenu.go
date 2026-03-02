package behaviors

import (
	"fmt"

	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

// play_Statue_butt — clicking goes to Track_Select_I
type PlayStatueButt struct {
	engine.DefaultBehavior
}

func (b *PlayStatueButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["challenge"] = 0.0
	g.GlobalVars["normalmodeselect"] = 1.0
	g.GlobalVars["impoppablemodeselect"] = 0.0
	g.GlobalVars["nightmaremodeselect"] = 0.0
	g.GotoNextRoom() // main_Menu -> Track_Select_I
}

// tower_Statue_butt — goes to Tower_Upgrades
type TowerStatueButt struct {
	engine.DefaultBehavior
}

func (b *TowerStatueButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.RequestRoomGoto("Tower_Upgrades")
}

// mission_Statue_butt — goes to Special_Missions0
type MissionStatueButt struct {
	engine.DefaultBehavior
}

func (b *MissionStatueButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["challenge"] = 0.0
	g.GlobalVars["normalmodeselect"] = 1.0
	g.GlobalVars["impoppablemodeselect"] = 0.0
	g.GlobalVars["nightmaremodeselect"] = 0.0
	g.RequestRoomGoto("Special_Missions0")
}

// achievements_Statue_butt — goes to Achievement_Room
type AchievementsStatueButt struct {
	engine.DefaultBehavior
}

func (b *AchievementsStatueButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.RequestRoomGoto("Achievement_Room")
}

// agent_Statue_butt — goes to Agents_and_other_goods
type AgentStatueButt struct {
	engine.DefaultBehavior
}

func (b *AgentStatueButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.RequestRoomGoto("Agents_and_other_goods")
}

// bounty_Statue_butt — goes to Bloons_Bounty_Center
type BountyStatueButt struct {
	engine.DefaultBehavior
}

func (b *BountyStatueButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.RequestRoomGoto("Bloons_Bounty_Center")
}

// challenges_statue_butt — goes to Challenge_Room
type ChallengesStatueButt struct {
	engine.DefaultBehavior
}

func (b *ChallengesStatueButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["challenge"] = 0.0
	g.GlobalVars["normalmodeselect"] = 1.0
	g.GlobalVars["impoppablemodeselect"] = 0.0
	g.GlobalVars["nightmaremodeselect"] = 0.0
	g.RequestRoomGoto("Challenge_Room")
}

// info_Butt — creates Game_Info panel
type InfoButt struct {
	engine.DefaultBehavior
}

func (b *InfoButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	info := g.InstanceMgr.Create("Game_Info", 512, 800)
	info.Direction = 90
	info.Speed = 75
	info.Friction = 5
	info.MotionSet(90, 75)
}

// credits_Butt — creates Credits_Info panel
type CreditsButt struct {
	engine.DefaultBehavior
}

func (b *CreditsButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	info := g.InstanceMgr.Create("Credits_Info", 512, 800)
	info.Direction = 90
	info.Speed = 75
	info.Friction = 5
	info.MotionSet(90, 75)
}

// main_Menu_Settings — settings gear icon; draws itself
type MainMenuSettings struct {
	engine.DefaultBehavior
}

func (b *MainMenuSettings) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["sandbox"] = 0.0
}

func (b *MainMenuSettings) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	// draw the settings sprite
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

// variable_set_upper — persistent object that maps Career_Control instance
// vars into globals. Since we already set globals in CareerControl,
// this is mostly a no-op in our port.
type VariableSetUpper struct {
	engine.DefaultBehavior
}

func (b *VariableSetUpper) Create(inst *engine.Instance, g *engine.Game) {
	inst.Persistent = true
	// variables are already set as globals by CareerControl
}

// save_Control — calls save every end step (no-op for now)
type SaveControl struct {
	engine.DefaultBehavior
}

// bP_Display — draws currency/rank HUD (simplified for now)
type BPDisplay struct {
	engine.DefaultBehavior
}

func (b *BPDisplay) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	// tODO: draw BP, monkey money, bsouls, trophies, rank
	// requires font rendering implementation
}

// level_uper — handles XP/rank progression
type LevelUper struct {
	engine.DefaultBehavior
}

func (b *LevelUper) Create(inst *engine.Instance, g *engine.Game) {
	inst.Persistent = true
	if _, ok := g.GlobalVars["rank"]; !ok {
		g.GlobalVars["rank"] = 0.0
	}
	if _, ok := g.GlobalVars["criteria"]; !ok {
		g.GlobalVars["criteria"] = 50.0
	}
	if _, ok := g.GlobalVars["XP"]; !ok {
		g.GlobalVars["XP"] = 0.0
	}
}

// sound_Control — sets sound volumes
type SoundControl struct {
	engine.DefaultBehavior
}

func (b *SoundControl) Create(inst *engine.Instance, g *engine.Game) {
	// set default volumes
	g.AudioMgr.SetVolume("Popping", 0.12)
	g.AudioMgr.SetVolume("Blimp_Hit", 0.5)
	g.AudioMgr.SetVolume("Blimp_Destroyed", 0.5)
	g.AudioMgr.SetVolume("Popp", 0.3)
	g.AudioMgr.SetVolume("Tower_Select", 0.5)
	g.AudioMgr.SetVolume("Bounty_BTFO", 0.5)
	g.AudioMgr.SetVolume("Small_Boom", 0.5)
	g.AudioMgr.SetVolume("Large_Boom", 0.5)
	g.AudioMgr.SetVolume("Ceramic_Hit", 0.5)
	g.AudioMgr.SetVolume("Winning", 0.5)
	g.AudioMgr.SetVolume("Lightning_Sound", 0.5)
	g.AudioMgr.SetVolume("Tower_Place", 0.5)
	g.AudioMgr.SetVolume("Temple_sound", 0.5)
	g.AudioMgr.SetVolume("Upgrade", 0.5)
	g.AudioMgr.SetVolume("Lead_Hit", 0.5)

	if _, ok := g.GlobalVars["mute"]; !ok {
		g.GlobalVars["mute"] = 0.0
	}
	if _, ok := g.GlobalVars["soundmute"]; !ok {
		g.GlobalVars["soundmute"] = 0.0
	}
}

// mM_Version — draws version text (simplified)
type MMVersion struct {
	engine.DefaultBehavior
}

func (b *MMVersion) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	// tODO: draw "Version 1.4\nMade by Ramaf Party" - needs font rendering
}

// moon_Temple_Launch_Pad — nightmare mode entrance
type MoonTempleLaunchPad struct {
	engine.DefaultBehavior
}

func (b *MoonTempleLaunchPad) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["challenge"] = 0.0
	g.GlobalVars["normalmodeselect"] = 0.0
	g.GlobalVars["impoppablemodeselect"] = 0.0
	g.GlobalVars["nightmaremodeselect"] = 1.0
	g.GotoNextRoom()
}

func (b *MoonTempleLaunchPad) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	sm12, _ := g.GlobalVars["specialmission12"].(float64)
	frameIdx := int(sm12) % len(spr.Frames)
	frame := spr.Frames[frameIdx]
	engine.DrawSpriteExt(screen, frame, spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// achieve_Control — calculates bsouls/trophies every step
// (simplified — full calculation)
type AchieveControl struct {
	engine.DefaultBehavior
}

func (b *AchieveControl) Create(inst *engine.Instance, g *engine.Game) {
	inst.Persistent = true
	inst.Depth = -10
}

func getGlobal(g *engine.Game, key string) float64 {
	v, _ := g.GlobalVars[key].(float64)
	return v
}

func (b *AchieveControl) Step(inst *engine.Instance, g *engine.Game) {
	// calculate bsouls
	bsouls := 0.0
	for i := 1; i <= 11; i++ {
		bsouls += getGlobal(g, fmt.Sprintf("b%d", i))
	}
	bsouls += getGlobal(g, "specialmission12")
	for i := 1; i <= 32; i++ {
		if getGlobal(g, fmt.Sprintf("track%dnightstone", i)) >= 4 {
			bsouls += 1
		}
	}
	g.GlobalVars["bsouls"] = bsouls

	// calculate trophies
	trophies := 0.0
	towerPaths := [][3]string{
		{"DML", "DMM", "DMR"}, {"TSL", "TSM", "TSR"}, {"BML", "BMM", "BMR"},
		{"SnML", "SnMM", "SnMR"}, {"NML", "NMM", "NMR"}, {"BCL", "BCM", "BCR"},
		{"MSL", "MSM", "MSR"}, {"CTL", "CTM", "CTR"},
		{"GGL", "GGM", "GGR"}, {"IML", "IMM", "IMR"}, {"MBL", "MBM", "MBR"},
		{"MEL", "MEM", "MER"}, {"MAL", "MAM", "MAR"}, {"BChL", "BChM", "BChR"},
		{"MApL", "MApM", "MApR"}, {"MAlL", "MAlM", "MAlR"},
		{"MVL", "MVM", "MVR"}, {"BTL", "BTM", "BTR"}, {"DGL", "DGM", "DGR"},
		{"MLL", "MLM", "MLR"}, {"HPL", "HPM", "HPR"}, {"SFL", "SFM", "SFR"},
		{"PML", "PMM", "PMR"}, {"SuML", "SuMM", "SuMR"},
	}
	for _, tp := range towerPaths {
		if getGlobal(g, tp[0]) >= 3 && getGlobal(g, tp[1]) >= 3 && getGlobal(g, tp[2]) >= 3 {
			trophies++
		}
		// secret path trophy (middle path >= 13)
		if getGlobal(g, tp[1]) >= 13 {
			trophies++
		}
	}
	for i := 1; i <= 32; i++ {
		if getGlobal(g, fmt.Sprintf("track%dmilestone", i)) >= 4 {
			trophies++
		}
		if getGlobal(g, fmt.Sprintf("track%dhardstone", i)) >= 4 {
			trophies++
		}
		if getGlobal(g, fmt.Sprintf("track%dnightstone", i)) >= 4 {
			trophies++
		}
		if getGlobal(g, fmt.Sprintf("track%dmilestone", i)) >= 6 {
			trophies++
		}
		if getGlobal(g, fmt.Sprintf("t%d", i)) <= 300 && getGlobal(g, fmt.Sprintf("t%d", i)) > 0 {
			trophies++
		}
		if getGlobal(g, fmt.Sprintf("x%d", i)) >= 1000000 {
			trophies++
		}
	}
	g.GlobalVars["trophies"] = trophies
}

// go_Back_to_Main_butt — the red X button in sub-menus, returns to Main_Menu
type GoBackToMainButt struct {
	engine.DefaultBehavior
}

func (b *GoBackToMainButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.RequestRoomGoto("Main_Menu")
}

// right_Main_butt — another red X that returns to Main_Menu
// (used in Track_Select_I, Achievement_Room, Special_Missions0, Challenge_Room)
type RightMainButt struct {
	engine.DefaultBehavior
}

func (b *RightMainButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.RequestRoomGoto("Main_Menu")
}

// go_to_Track_Select_butt — red X that returns to Track_Select_I
// (used in Track_Setup_II, Bloons_Bounty_Center)
type GoToTrackSelectButt struct {
	engine.DefaultBehavior
}

func (b *GoToTrackSelectButt) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.RequestRoomGoto("Track_Select_I")
}

// registerMainMenuBehaviors registers all Main_Menu room behaviors
func RegisterMainMenuBehaviors(im *engine.InstanceManager) {
	im.RegisterBehavior("Play_Statue_butt", func() engine.InstanceBehavior { return &PlayStatueButt{} })
	im.RegisterBehavior("Tower_Statue_butt", func() engine.InstanceBehavior { return &TowerStatueButt{} })
	im.RegisterBehavior("Mission_Statue_butt", func() engine.InstanceBehavior { return &MissionStatueButt{} })
	im.RegisterBehavior("Achievements_Statue_butt", func() engine.InstanceBehavior { return &AchievementsStatueButt{} })
	im.RegisterBehavior("Agent_Statue_butt", func() engine.InstanceBehavior { return &AgentStatueButt{} })
	im.RegisterBehavior("Bounty_Statue_butt", func() engine.InstanceBehavior { return &BountyStatueButt{} })
	im.RegisterBehavior("Challenges_statue_butt", func() engine.InstanceBehavior { return &ChallengesStatueButt{} })
	im.RegisterBehavior("Info_Butt", func() engine.InstanceBehavior { return &InfoButt{} })
	im.RegisterBehavior("Credits_Butt", func() engine.InstanceBehavior { return &CreditsButt{} })
	im.RegisterBehavior("Main_Menu_Settings", func() engine.InstanceBehavior { return &MainMenuSettings{} })
	im.RegisterBehavior("Variable_set_upper", func() engine.InstanceBehavior { return &VariableSetUpper{} })
	im.RegisterBehavior("Save_Control", func() engine.InstanceBehavior { return &SaveControl{} })
	im.RegisterBehavior("BP_Display", func() engine.InstanceBehavior { return &BPDisplay{} })
	im.RegisterBehavior("Level_uper", func() engine.InstanceBehavior { return &LevelUper{} })
	im.RegisterBehavior("Sound_Control", func() engine.InstanceBehavior { return &SoundControl{} })
	im.RegisterBehavior("MM_Version", func() engine.InstanceBehavior { return &MMVersion{} })
	im.RegisterBehavior("Moon_Temple_Launch_Pad", func() engine.InstanceBehavior { return &MoonTempleLaunchPad{} })
	im.RegisterBehavior("Achieve_Control", func() engine.InstanceBehavior { return &AchieveControl{} })
	im.RegisterBehavior("Go_Back_to_Main_butt", func() engine.InstanceBehavior { return &GoBackToMainButt{} })
	im.RegisterBehavior("Right_Main_butt", func() engine.InstanceBehavior { return &RightMainButt{} })
	im.RegisterBehavior("Go_to_Track_Select_butt", func() engine.InstanceBehavior { return &GoToTrackSelectButt{} })
}
