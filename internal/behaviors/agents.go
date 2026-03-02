package behaviors

import (
	"btdx/internal/engine"

	"github.com/hajimehoshi/ebiten/v2"
)

// agent panel helpers
// all agent panels are 4-frame sprites (L1/L2/L3/L4).
// image_speed=0.01, draws sub-image 0
// (always frame 0). The level text is baked into the sprite frames.

// agentPanelDraw draws the panel sprite at frame 0, plus the owned count.
func agentPanelDraw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil || len(spr.Frames) == 0 {
		return
	}
	// always draw frame 0
	frame := spr.Frames[0]
	engine.DrawSpriteExt(screen, frame, spr.XOrigin, spr.YOrigin,
		inst.X, inst.Y, inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// angry_Squirrel_Panel — costs 60 MM
type AngrySquirrelPanel struct {
	engine.DefaultBehavior
}

func (b *AngrySquirrelPanel) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0.01
}

func (b *AngrySquirrelPanel) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	agentPanelDraw(inst, screen, g)
}

func (b *AngrySquirrelPanel) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	mm := getGlobal(g, "monkeymoney")
	if mm >= 60 {
		g.GlobalVars["angrysquirrel"] = getGlobal(g, "angrysquirrel") + 1
		g.GlobalVars["monkeymoney"] = mm - 60
	}
}

// bloonbury_Bush_Panel — costs 70 MM
type BloonburyBushPanel struct {
	engine.DefaultBehavior
}

func (b *BloonburyBushPanel) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0.01
}

func (b *BloonburyBushPanel) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	agentPanelDraw(inst, screen, g)
}

func (b *BloonburyBushPanel) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	mm := getGlobal(g, "monkeymoney")
	if mm >= 70 {
		g.GlobalVars["bloonburybush"] = getGlobal(g, "bloonburybush") + 1
		g.GlobalVars["monkeymoney"] = mm - 70
	}
}

// sprinkler_Panel — costs 100 MM
type SprinklerPanel struct {
	engine.DefaultBehavior
}

func (b *SprinklerPanel) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0.01
}

func (b *SprinklerPanel) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	agentPanelDraw(inst, screen, g)
}

func (b *SprinklerPanel) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	mm := getGlobal(g, "monkeymoney")
	if mm >= 100 {
		g.GlobalVars["sprinkler"] = getGlobal(g, "sprinkler") + 1
		g.GlobalVars["monkeymoney"] = mm - 100
	}
}

// monkey_Nurse_Panel — costs 90 MM
type MonkeyNursePanel struct {
	engine.DefaultBehavior
}

func (b *MonkeyNursePanel) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0.01
}

func (b *MonkeyNursePanel) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	agentPanelDraw(inst, screen, g)
}

func (b *MonkeyNursePanel) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	mm := getGlobal(g, "monkeymoney")
	if mm >= 90 {
		g.GlobalVars["monkeynurse"] = getGlobal(g, "monkeynurse") + 1
		g.GlobalVars["monkeymoney"] = mm - 90
	}
}

// banana_Mobile_Panel — costs 100 MM
type BananaMobilePanel struct {
	engine.DefaultBehavior
}

func (b *BananaMobilePanel) Create(inst *engine.Instance, g *engine.Game) {
	inst.ImageSpeed = 0.01
	if _, ok := g.GlobalVars["bananamobile"]; !ok {
		g.GlobalVars["bananamobile"] = 0.0
	}
}

func (b *BananaMobilePanel) Draw(inst *engine.Instance, screen *ebiten.Image, g *engine.Game) {
	agentPanelDraw(inst, screen, g)
}

func (b *BananaMobilePanel) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	mm := getGlobal(g, "monkeymoney")
	if mm >= 100 {
		g.GlobalVars["bananamobile"] = getGlobal(g, "bananamobile") + 1
		g.GlobalVars["monkeymoney"] = mm - 100
	}
}

// bloon_Points — costs 400 MM for 4 BP
type BloonPoints struct {
	engine.DefaultBehavior
}

func (b *BloonPoints) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	mm := getGlobal(g, "monkeymoney")
	if mm >= 400 {
		g.GlobalVars["BP"] = getGlobal(g, "BP") + 4
		g.GlobalVars["monkeymoney"] = mm - 400
	}
}

// registerAgentsBehaviors registers all Agents_and_other_goods room behaviors
func RegisterAgentsBehaviors(im *engine.InstanceManager) {
	im.RegisterBehavior("Angry_Squirrel_Panel", func() engine.InstanceBehavior { return &AngrySquirrelPanel{} })
	im.RegisterBehavior("Bloonbury_Bush_Panel", func() engine.InstanceBehavior { return &BloonburyBushPanel{} })
	im.RegisterBehavior("Sprinkler_Panel", func() engine.InstanceBehavior { return &SprinklerPanel{} })
	im.RegisterBehavior("Monkey_Nurse_Panel", func() engine.InstanceBehavior { return &MonkeyNursePanel{} })
	im.RegisterBehavior("Banana_Mobile_Panel", func() engine.InstanceBehavior { return &BananaMobilePanel{} })
	im.RegisterBehavior("Bloon_Points", func() engine.InstanceBehavior { return &BloonPoints{} })
}
