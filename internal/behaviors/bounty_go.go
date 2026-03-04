package behaviors

import (
	"fmt"
	"math"

	"btdx/internal/engine"
	"btdx/internal/savedata"

	"github.com/hajimehoshi/ebiten/v2"
)

// ---------------------------------------------------------------------------
// BountyGoBehavior — wave controller for bounty/boss challenge rooms.
// Replaces GoBehavior with configurable starting values per boss/difficulty.
// Reuses the existing timeline system for wave spawning.
// ---------------------------------------------------------------------------

type BountyGoConfig struct {
	StartMoney float64
	StartLife  float64
	BPower     float64 // bloon power multiplier
	BSpeed     float64 // bloon speed multiplier
}

type BountyGoBehavior struct {
	engine.DefaultBehavior
	Config            BountyGoConfig
	shiftpress        int
	afterwave         int
	autoButtonChecked bool
}

func (b *BountyGoBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["autostart"] = 0.0
	g.GlobalVars["cashinflate"] = 0.0
	g.GlobalVars["freeplay"] = 0.0
	g.GlobalVars["cashflow"] = 1.0
	g.GlobalVars["cashwavereward"] = 0.0
	g.GlobalVars["challenge"] = 0.0

	bpower := b.Config.BPower
	if bpower == 0 {
		bpower = 1.0
	}
	g.GlobalVars["bpower"] = bpower

	bspeed := b.Config.BSpeed
	if bspeed == 0 {
		bspeed = 1.0
	}
	g.GlobalVars["bspeed"] = bspeed

	g.GlobalVars["wave"] = 1.0
	g.GlobalVars["wavenow"] = 0.0
	g.GlobalVars["cycle"] = 0.0
	g.GlobalVars["endsequence"] = 0.0
	g.GlobalVars["points"] = 0.0
	g.GlobalVars["gamespeed"] = 30.0
	g.SetGameSpeed(30)

	money := b.Config.StartMoney
	if money == 0 {
		money = 944
	}
	g.GlobalVars["money"] = money

	life := b.Config.StartLife
	if life == 0 {
		life = 200
	}
	g.GlobalVars["life"] = life

	// PlayBar normally sets these before entering a gameplay room.
	// The bounty flow skips PlayBar, so initialise them here.
	if getGlobal(g, "pointmultiplier") == 0 {
		g.GlobalVars["pointmultiplier"] = 1.0
	}
	if getGlobal(g, "towerlimit") == 0 {
		g.GlobalVars["towerlimit"] = 1000000.0
	}
	// reset modifier flags (bounty fights don't use track-setup modifiers)
	g.GlobalVars["sixtowers"] = 0.0
	g.GlobalVars["randomtowers"] = 0.0
	g.GlobalVars["wavesqueeze"] = 0.0
	g.GlobalVars["waveskip"] = 0.0
	g.GlobalVars["strongerbloons"] = 0.0
	g.GlobalVars["fasterbloons"] = 0.0
	g.GlobalVars["noliveslost"] = 0.0

	// Unlock tower buy panels based on rank (same as PlayBar.unlockTowers).
	// Without this the panels stay hidden because all lock globals default to 0.
	unlockTowersForRank(g)

	b.shiftpress = 0
	b.afterwave = -1
	b.autoButtonChecked = false
}

func (b *BountyGoBehavior) Step(inst *engine.Instance, g *engine.Game) {
	// right-click near the Go button toggles 10x speed
	if g.InputMgr.MouseRightPressed() {
		mx, my := g.GetMouseRoomPos()
		dx := mx - inst.X
		dy := my - inst.Y
		if dx*dx+dy*dy <= 48*48 {
			b.toggleHyperSpeed(inst, g)
			return
		}
	}

	if !b.autoButtonChecked {
		b.autoButtonChecked = true
		if g.InstanceMgr.InstanceCount("auto_start_button") == 0 {
			g.InstanceMgr.Create("auto_start_button", 944, 544)
		}
	}

	wavenow := getGlobal(g, "wavenow")
	wave := getGlobal(g, "wave")
	bloonCount := float64(countBloons(g))

	// when timeline finishes spawning, mark wavenow=0
	if wavenow == 1 && g.ActiveTimeline != nil && !g.ActiveTimeline.Running {
		g.GlobalVars["wavenow"] = 0.0
		wavenow = 0
	}

	// switch sprite
	if wavenow == 1 || bloonCount > 0 {
		inst.SpriteName = "Going"
	} else {
		inst.SpriteName = "sprite278"
	}

	// wave-end logic
	if wavenow == 0 && bloonCount == 0 && b.afterwave == 0 {
		inst.SpriteName = "sprite278"

		cashReward := getGlobal(g, "cashwavereward")
		cashFlow := getGlobal(g, "cashflow")
		cashInflate := getGlobal(g, "cashinflate")
		money := getGlobal(g, "money")

		reward := math.Round(cashReward * cashFlow * (1.0 + cashInflate*0.1))
		money += reward

		pm := getGlobal(g, "pointmultiplier")
		if pm == 0 {
			pm = 1
		}
		pts := getGlobal(g, "points")
		pts += (100.0 + wave*3.0) * pm

		xpGain := (10.0 + wave*2.0) * pm
		g.GlobalVars["XP"] = getGlobal(g, "XP") + xpGain

		g.GlobalVars["money"] = money
		g.GlobalVars["points"] = pts
		g.GlobalVars["cashwavereward"] = 0.0
		b.afterwave = 1

		if err := savedata.Save(g); err != nil {
			fmt.Printf("WARNING: bounty auto-save failed: %v\n", err)
		}
	}

	// auto-start
	if getGlobal(g, "autostart") == 1 && bloonCount == 0 && wavenow == 0 && b.afterwave == 1 {
		b.startNextWave(inst, g)
	}
}

func (b *BountyGoBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	b.handleWaveStart(inst, g)
}

func (b *BountyGoBehavior) MouseRightPressed(inst *engine.Instance, g *engine.Game) {
	b.toggleHyperSpeed(inst, g)
}

func (b *BountyGoBehavior) KeyPress(inst *engine.Instance, g *engine.Game) {
	if g.InputMgr.KeyPressed(ebiten.KeySpace) || g.InputMgr.KeyPressed(ebiten.KeyEnter) {
		b.handleWaveStart(inst, g)
	}
}

func (b *BountyGoBehavior) toggleHyperSpeed(inst *engine.Instance, g *engine.Game) {
	if getGlobal(g, "gamespeed") == 300 {
		g.SetGameSpeed(30)
		g.GlobalVars["gamespeed"] = 30.0
		b.shiftpress = 0
	} else {
		g.SetGameSpeed(300)
		g.GlobalVars["gamespeed"] = 300.0
	}
}

func (b *BountyGoBehavior) handleWaveStart(inst *engine.Instance, g *engine.Game) {
	wavenow := getGlobal(g, "wavenow")
	bloonCount := float64(countBloons(g))

	if wavenow == 1 || bloonCount > 0 {
		switch b.shiftpress {
		case 0:
			g.SetGameSpeed(60)
			g.GlobalVars["gamespeed"] = 60.0
			b.shiftpress = 1
		case 1:
			g.SetGameSpeed(90)
			g.GlobalVars["gamespeed"] = 90.0
			b.shiftpress = 2
		case 2:
			g.SetGameSpeed(30)
			g.GlobalVars["gamespeed"] = 30.0
			b.shiftpress = 0
		}
		return
	}

	if bloonCount == 0 && wavenow == 0 {
		b.startNextWave(inst, g)
	}
}

func (b *BountyGoBehavior) startNextWave(inst *engine.Instance, g *engine.Game) {
	wave := int(getGlobal(g, "wave"))
	if wave > 90 {
		return
	}

	inst.SpriteName = "Going"

	tlName := fmt.Sprintf("N%d", wave)
	tl := g.TimelineMgr.Get(tlName)
	if tl != nil {
		wsqueeze := getGlobal(g, "wavesqueeze")
		runner := engine.NewTimelineRunner(tl)
		runner.Speed = 1.0 + wsqueeze
		runner.OnAction = func(action engine.TimelineAction) {
			if action.WhoName == "self" {
				if val := extractCodeVar(action.Code, "cashwavereward"); val > 0 {
					g.GlobalVars["cashwavereward"] = val
				}
				if val := extractCodeVar(action.Code, "cashflow"); val > 0 {
					g.GlobalVars["cashflow"] = val
				}
				return
			}
			executeBloonSpawn(g, action)
		}
		g.ActiveTimeline = runner
		fmt.Printf("Bounty wave %d started (timeline %s)\n", wave, tlName)
	} else {
		fmt.Printf("WARNING: Timeline %s not found for bounty wave %d\n", tlName, wave)
		inst.SpriteName = "sprite278"
		return
	}

	g.GlobalVars["wavenow"] = 1.0
	g.GlobalVars["cycle"] = math.Mod(getGlobal(g, "cycle")+1, 4)
	if getGlobal(g, "cycle") == 0 {
		g.GlobalVars["cycle"] = 1.0
	}
	b.afterwave = 0
	g.GlobalVars["wave"] = float64(wave + 1)
}

// Draw renders the Go button with speed indicators (same as GoBehavior)
func (b *BountyGoBehavior) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr != nil && len(spr.Frames) > 0 {
		engine.DrawSpriteExt(screen, spr.Frames[0], spr.XOrigin, spr.YOrigin,
			inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
	}

	gs := getGlobal(g, "gamespeed")
	fill := gs / 90.0
	if fill > 1 {
		fill = 1
	}
	barH := 44.0
	filledH := barH * fill
	barBottom := inst.Y + 54

	drawRect(screen, inst.X+8, barBottom-filledH, 2, filledH, speedBarColor(fill))
	drawRect(screen, inst.X+8, inst.Y+10, 2, barH-filledH, color0x000000)
	drawRect(screen, inst.X+53, barBottom-filledH, 2, filledH, speedBarColor(fill))
	drawRect(screen, inst.X+53, inst.Y+10, 2, barH-filledH, color0x000000)
}

// ---------------------------------------------------------------------------
// SuperBountyGoBehavior — the 4 "prison break" special bounty fights.
// Locks ALL towers and sets towerlimit to 1000.
// ---------------------------------------------------------------------------
type SuperBountyGoBehavior struct {
	BountyGoBehavior
	MissionIndex int // 1-4 → maps to specialmission
}

func (b *SuperBountyGoBehavior) Create(inst *engine.Instance, g *engine.Game) {
	// lock all tower types
	towerLocks := []string{
		"DMlock", "TSlock", "BMlock", "SnMlock", "NMlock", "BClock",
		"MSlock", "CTlock", "GGlock", "IMlock", "MBlock", "MElock",
		"MAlock", "BChlock", "MAplock", "MAllock", "MVlock", "BTlock",
		"DGlock", "MLlock", "HPlock", "SFlock", "PMlock", "SuMlock",
	}
	for _, lk := range towerLocks {
		g.GlobalVars[lk] = 1.0
	}
	g.GlobalVars["towerlimit"] = 1000.0

	// call parent Create
	b.BountyGoBehavior.Create(inst, g)
}

// ---------------------------------------------------------------------------
// BossTuneBehavior — plays Boss_Music in loop.
// ---------------------------------------------------------------------------
type BossTuneBehavior struct {
	engine.DefaultBehavior
}

func (b *BossTuneBehavior) Create(inst *engine.Instance, g *engine.Game) {
	g.AudioMgr.PlayMusic("Boss_Music")
}

func (b *BossTuneBehavior) Destroy(inst *engine.Instance, g *engine.Game) {
	g.AudioMgr.StopMusic()
}

// ---------------------------------------------------------------------------
// Registration helpers — bounty _Go configs by boss type and difficulty.
// ---------------------------------------------------------------------------

// bountyGoConfigs maps _Go object names to their starting configs.
var bountyGoConfigs = map[string]BountyGoConfig{
	// Bully
	"Bully_Go":  {StartMoney: 35000, StartLife: 200, BPower: 0.6, BSpeed: 1},
	"Bully_Go2": {StartMoney: 944, StartLife: 369, BPower: 1, BSpeed: 1},
	"Bully_Go3": {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"Bully_Go4": {StartMoney: 944, StartLife: 1, BPower: 1, BSpeed: 1},

	// Mother
	"Mother_Go":  {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"Mother_Go2": {StartMoney: 944, StartLife: 369, BPower: 1, BSpeed: 1},
	"Mother_Go3": {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"Mother_Go4": {StartMoney: 944, StartLife: 1, BPower: 1, BSpeed: 1},

	// Clown
	"Clown_Go":  {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"Clown_Go2": {StartMoney: 944, StartLife: 369, BPower: 1, BSpeed: 1},
	"Clown_Go3": {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"Clown_Go4": {StartMoney: 944, StartLife: 1, BPower: 1, BSpeed: 1},

	// LUL
	"LUL_Go":  {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"LUL_Go2": {StartMoney: 944, StartLife: 369, BPower: 1, BSpeed: 1},
	"LUL_Go3": {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"LUL_Go4": {StartMoney: 944, StartLife: 1, BPower: 1, BSpeed: 1},

	// Blooming
	"Blooming_Go":  {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"Blooming_Go2": {StartMoney: 944, StartLife: 369, BPower: 1, BSpeed: 1},
	"Blooming_Go3": {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"Blooming_Go4": {StartMoney: 944, StartLife: 1, BPower: 1, BSpeed: 1},

	// Crawler
	"Crawler_Go":  {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"Crawler_Go2": {StartMoney: 944, StartLife: 369, BPower: 1, BSpeed: 1},
	"Crawler_Go3": {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
	"Crawler_Go4": {StartMoney: 944, StartLife: 1, BPower: 1, BSpeed: 1},

	// Generic fallback for any unmapped Go
	"Bounty_Go": {StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
}

func RegisterBountyGoBehaviors(im *engine.InstanceManager) {
	// register all known _Go controllers
	for name, cfg := range bountyGoConfigs {
		c := cfg // capture
		im.RegisterBehavior(name, func() engine.InstanceBehavior {
			return &BountyGoBehavior{Config: c}
		})
	}

	// Super Bounty _Go (Prison Break)
	for i := 1; i <= 4; i++ {
		idx := i
		names := []string{"Super_Bounty_Go_I", "Super_Bounty_Go_II", "Super_Bounty_Go_III", "Super_Bounty_Go_IV"}
		im.RegisterBehavior(names[idx-1], func() engine.InstanceBehavior {
			return &SuperBountyGoBehavior{
				BountyGoBehavior: BountyGoBehavior{
					Config: BountyGoConfig{StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
				},
				MissionIndex: idx,
			}
		})
	}

	// Boss audio
	im.RegisterBehavior("Boss_Tune", func() engine.InstanceBehavior { return &BossTuneBehavior{} })

	// Boss Bash (tutorial)
	im.RegisterBehavior("Boss_Bash_go", func() engine.InstanceBehavior {
		return &BountyGoBehavior{
			Config: BountyGoConfig{StartMoney: 944, StartLife: 200, BPower: 1, BSpeed: 1},
		}
	})
}
