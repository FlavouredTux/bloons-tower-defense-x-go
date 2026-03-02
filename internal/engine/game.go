// package engine provides the core game engine.
// handles rooms, instances, sprites, events, and the main loop.
package engine

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	// game window dimensions
	ScreenWidth  = 1024
	ScreenHeight = 576

	// default game speed (frames per second) from room settings
	DefaultGameSpeed = 30
)

// game is the top-level Ebitengine game struct
type Game struct {
	// core systems
	AssetManager *AssetManager
	RoomManager  *RoomManager
	InstanceMgr  *InstanceManager
	AudioMgr     *AudioManager
	InputMgr     *InputManager
	PathMgr      *PathManager
	TimelineMgr  *TimelineManager

	// game state
	CurrentRoom string
	GameSpeed   int // target FPS for game logic (room_speed)
	GlobalVars  map[string]interface{}

	// bitmap font for HUD text rendering
	BMFont *BMFont

	// wave/timeline execution
	ActiveTimeline *TimelineRunner

	// timing
	tickAccum    float64
	tickInterval float64 // seconds per game tick
	lastUpdate   time.Time
	paused       bool
	pendingRoom  string // deferred room change
}

// newGame creates and initializes the game engine
func NewGame() *Game {
	g := &Game{
		AssetManager: NewAssetManager(),
		RoomManager:  NewRoomManager(),
		InstanceMgr:  NewInstanceManager(),
		AudioMgr:     NewAudioManager(),
		InputMgr:     NewInputManager(),
		PathMgr:      NewPathManager(),
		TimelineMgr:  NewTimelineManager(),
		GameSpeed:    DefaultGameSpeed,
		GlobalVars:   make(map[string]interface{}),
	}
	g.tickInterval = 1.0 / float64(g.GameSpeed)
	g.InstanceMgr.SetGameRef(g)
	return g
}

// update implements ebiten.Game interface - called every frame (60fps)
func (g *Game) Update() error {
	// update input state
	g.InputMgr.Update()
	g.AudioMgr.Update()

	// eSC to quit (backup for window X button)
	if g.InputMgr.KeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	if g.paused {
		return nil
	}

	// delta time: measure actual elapsed time between Update() calls
	now := time.Now()
	if g.lastUpdate.IsZero() {
		g.lastUpdate = now
	}
	dt := now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	// cap delta to prevent lag spikes
	// max ~3 catch-up ticks at 30fps game speed, ~6 at 60fps
	if dt > 0.1 {
		dt = 0.1
	}

	g.tickAccum += dt

	// limit ticks per frame to avoid runaway updates
	ticksThisFrame := 0
	const maxTicksPerFrame = 4

	for g.tickAccum >= g.tickInterval && ticksThisFrame < maxTicksPerFrame {
		g.tickAccum -= g.tickInterval
		ticksThisFrame++

		// run one game tick
		g.gameTick()

		// consume single-frame input events so they don't fire again
		g.InputMgr.ConsumeEvents()
	}

	return nil
}

// gameTick runs one game logic step
func (g *Game) gameTick() {
	instances := g.InstanceMgr.GetAll()

	// step Begin events
	for _, inst := range instances {
		if inst.Active && !inst.Destroyed {
			inst.OnStepBegin()
		}
	}

	// alarm events
	for _, inst := range instances {
		if inst.Active && !inst.Destroyed {
			for i := range inst.Alarms {
				if inst.Alarms[i] > 0 {
					inst.Alarms[i]--
					if inst.Alarms[i] == 0 {
						inst.OnAlarm(i)
					}
				}
			}
		}
	}

	// tick active timeline (wave spawning)
	if g.ActiveTimeline != nil {
		g.ActiveTimeline.Tick()
	}

	// key Press events (any key)
	if g.InputMgr.AnyKeyPressed() {
		for _, inst := range g.InstanceMgr.GetAll() {
			if inst.Active && !inst.Destroyed {
				inst.Behavior.KeyPress(inst, g)
			}
		}
	}

	// mouse events
	if g.InputMgr.MouseLeftPressed() {
		// global left pressed → fire on all instances
		for _, inst := range g.InstanceMgr.GetAll() {
			if inst.Active && !inst.Destroyed {
				inst.Behavior.MouseGlobalLeftPressed(inst, g)
			}
		}
		// per-instance left pressed → only on instances under cursor
		for _, inst := range g.InstanceMgr.GetAll() {
			if inst.Active && !inst.Destroyed {
				if g.IsMouseOverInstance(inst) {
					inst.Behavior.MouseLeftPressed(inst, g)
				}
			}
		}
	}
	if g.InputMgr.MouseRightPressed() {
		for _, inst := range g.InstanceMgr.GetAll() {
			if inst.Active && !inst.Destroyed {
				if g.IsMouseOverInstance(inst) {
					inst.Behavior.MouseRightPressed(inst, g)
				}
			}
		}
	}
	if g.InputMgr.MouseLeftReleased() {
		// global left released → fire on all instances
		for _, inst := range g.InstanceMgr.GetAll() {
			if inst.Active && !inst.Destroyed {
				inst.Behavior.MouseGlobalLeftReleased(inst, g)
			}
		}
	}

	// step events
	for _, inst := range instances {
		if inst.Active && !inst.Destroyed {
			inst.OnStep()
		}
	}

	// collision events
	g.InstanceMgr.CheckCollisions()

	// step End events
	for _, inst := range instances {
		if inst.Active && !inst.Destroyed {
			inst.OnStepEnd()
		}
	}

	// remove destroyed instances
	g.InstanceMgr.FlushDestroyed()

	// advance sprite animations
	for _, inst := range g.InstanceMgr.GetAll() {
		if inst.Active {
			inst.AdvanceAnimation()
		}
	}

	// update view following
	g.updateViewFollowing()

	// process pending room change
	if g.pendingRoom != "" {
		room := g.pendingRoom
		g.pendingRoom = ""
		g.GotoRoom(room)
	}
}

// getActiveView returns a pointer to the first enabled view for the current room, or a default.
func (g *Game) getActiveView(room *Room) *RoomView {
	if room != nil && room.EnableViews {
		for i := range room.Views {
			if room.Views[i].Visible {
				return &room.Views[i]
			}
		}
	}
	// default: show top-left of room at screen size
	return &RoomView{
		Visible: true,
		XView:   0, YView: 0,
		WView: ScreenWidth, HView: ScreenHeight,
		XPort: 0, YPort: 0,
		WPort: ScreenWidth, HPort: ScreenHeight,
	}
}

// updateViewFollowing handles view-follows-object logic.
// if a view has ObjName set, the view tracks that object within HBorder/VBorder.
func (g *Game) updateViewFollowing() {
	room := g.RoomManager.GetCurrent()
	if room == nil || !room.EnableViews {
		return
	}
	for i := range room.Views {
		v := &room.Views[i]
		if !v.Visible || v.ObjName == "" {
			continue
		}
		// find the first instance of the target object
		targets := g.InstanceMgr.FindByObject(v.ObjName)
		if len(targets) == 0 {
			continue
		}
		target := targets[0]

		// keep the object within HBorder/VBorder of the view edges
		// (Not the center)
		hb := float64(v.HBorder)
		vb := float64(v.VBorder)

		// a value of -1 for borders often means
		// "center the object". If border is very large, clamp it.
		if hb > float64(v.WView)/2 {
			hb = float64(v.WView) / 2
		}
		if vb > float64(v.HView)/2 {
			vb = float64(v.HView) / 2
		}

		// horizontal follow
		if target.X < float64(v.XView)+hb {
			v.XView = int(target.X - hb)
		} else if target.X > float64(v.XView)+float64(v.WView)-hb {
			v.XView = int(target.X - float64(v.WView) + hb)
		}

		// vertical follow
		if target.Y < float64(v.YView)+vb {
			v.YView = int(target.Y - vb)
		} else if target.Y > float64(v.YView)+float64(v.HView)-vb {
			v.YView = int(target.Y - float64(v.HView) + vb)
		}

		// clamp to room bounds
		if v.XView < 0 {
			v.XView = 0
		}
		if v.YView < 0 {
			v.YView = 0
		}
		if room.Width > 0 && v.XView+v.WView > room.Width {
			v.XView = room.Width - v.WView
		}
		if room.Height > 0 && v.YView+v.HView > room.Height {
			v.YView = room.Height - v.HView
		}
	}
}

// draw implements ebiten.Game interface
func (g *Game) Draw(screen *ebiten.Image) {
	room := g.RoomManager.GetCurrent()
	view := g.getActiveView(room)

	// clear with room background color
	if room != nil && room.ShowBgColor {
		r, gc, b := colorFromBGR(room.BgColor)
		screen.Fill(color.RGBA{r, gc, b, 255})
	} else {
		screen.Fill(color.Black)
	}

	// view transform: maps room coords to screen coords
	// scaleX/Y handles the case where view size != port size
	viewScaleX := float64(view.WPort) / float64(view.WView)
	viewScaleY := float64(view.HPort) / float64(view.HView)
	viewX := float64(view.XView)
	viewY := float64(view.YView)
	portX := float64(view.XPort)
	portY := float64(view.YPort)

	// draw room backgrounds (non-foreground)
	if room != nil {
		for _, bg := range room.Backgrounds {
			if bg.Visible && !bg.Foreground {
				g.drawBackground(screen, &bg, viewX, viewY, viewScaleX, viewScaleY, portX, portY, view.WPort, view.HPort)
			}
		}
	}

	// draw room tiles behind instances (non-negative depth = behind)
	if room != nil {
		g.drawTiles(screen, room, viewX, viewY, viewScaleX, viewScaleY, portX, portY)
	}

	// draw instances sorted by depth (higher depth = drawn first / further back)
	instances := g.InstanceMgr.GetSortedByDepth()
	for _, inst := range instances {
		if inst.Visible && !inst.Destroyed {
			// translate instance drawing by view offset
			oldX, oldY := inst.X, inst.Y
			inst.X = portX + (inst.X-viewX)*viewScaleX
			inst.Y = portY + (inst.Y-viewY)*viewScaleY
			inst.OnDraw(screen)
			inst.X, inst.Y = oldX, oldY
		}
	}

	// draw tower range overlays (hovered or selected towers).
	g.drawTowerRanges(screen, instances, viewX, viewY, viewScaleX, viewScaleY, portX, portY)

	// draw foreground backgrounds
	if room != nil {
		for _, bg := range room.Backgrounds {
			if bg.Visible && bg.Foreground {
				g.drawBackground(screen, &bg, viewX, viewY, viewScaleX, viewScaleY, portX, portY, view.WPort, view.HPort)
			}
		}
	}
}

func varToFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case float32:
		return float64(n), true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	default:
		return 0, false
	}
}

func instanceVarFloat(inst *Instance, key string) (float64, bool) {
	if inst == nil || inst.Vars == nil {
		return 0, false
	}
	v, ok := inst.Vars[key]
	if !ok {
		return 0, false
	}
	return varToFloat(v)
}

func (g *Game) drawTowerRanges(screen *ebiten.Image, instances []*Instance,
	viewX, viewY, viewScaleX, viewScaleY, portX, portY float64) {
	scale := (viewScaleX + viewScaleY) * 0.5
	for _, inst := range instances {
		if inst == nil || inst.Destroyed || !inst.Visible {
			continue
		}
		if _, ok := inst.Vars["select"]; !ok {
			continue
		}
		rng, ok := instanceVarFloat(inst, "range")
		if !ok || rng <= 0 {
			continue
		}
		// skip range circle for global-range towers (e.g. Sniper, range=2000).
		globalRange := rng >= 1500
		sel, _ := instanceVarFloat(inst, "select")
		if sel != 1 && !g.IsMouseOverInstance(inst) {
			continue
		}

		sx := portX + (inst.X-viewX)*viewScaleX
		sy := portY + (inst.Y-viewY)*viewScaleY

		if !globalRange {
			rr := rng * scale
			if rr > 0 {
				vector.DrawFilledCircle(screen, float32(sx), float32(sy), float32(rr), color.RGBA{60, 60, 60, 140}, true)
				vector.StrokeCircle(screen, float32(sx), float32(sy), float32(rr), 1.5, color.RGBA{35, 35, 35, 220}, true)
			}
		}

		// ability charge bar (legacy style): shown under towers that have an ability.
		abilityMax, hasAbility := instanceVarFloat(inst, "ability_max")
		if !hasAbility || abilityMax <= 0 {
			continue
		}
		abilityCur, _ := instanceVarFloat(inst, "ability")
		if abilityCur < 0 {
			abilityCur = 0
		}
		if abilityCur > abilityMax {
			abilityCur = abilityMax
		}
		pct := abilityCur / abilityMax

		barW := 30.0 * scale
		if barW < 20 {
			barW = 20
		}
		barH := 5.0 * scale
		if barH < 3 {
			barH = 3
		}
		bx := sx - barW*0.5
		by := sy + 18.0*scale

		// outer and inner background.
		vector.DrawFilledRect(screen, float32(bx-1), float32(by-1), float32(barW+2), float32(barH+2), color.RGBA{20, 20, 20, 220}, false)
		vector.DrawFilledRect(screen, float32(bx), float32(by), float32(barW), float32(barH), color.RGBA{55, 55, 55, 220}, false)
		if pct > 0 {
			fill := barW * pct
			fillColor := color.RGBA{90, 190, 255, 230}
			if pct >= 1 {
				fillColor = color.RGBA{255, 220, 70, 245}
			}
			vector.DrawFilledRect(screen, float32(bx), float32(by), float32(fill), float32(barH), fillColor, false)
		}
	}
}

func (g *Game) drawBackground(screen *ebiten.Image, bg *RoomBackground,
	viewX, viewY, scaleX, scaleY, portX, portY float64, portW, portH int) {
	if bg.Name == "" {
		return
	}
	img := g.AssetManager.GetBackground(bg.Name)
	if img == nil {
		return
	}
	bw, bh := img.Bounds().Dx(), img.Bounds().Dy()
	if bw == 0 || bh == 0 {
		return
	}

	if bg.Stretch {
		// stretch to fill port
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(portW)/float64(bw), float64(portH)/float64(bh))
		op.GeoM.Translate(portX, portY)
		screen.DrawImage(img, op)
		return
	}

	// calculate tiling
	startX := float64(bg.X) - viewX
	startY := float64(bg.Y) - viewY

	if bg.HTiled {
		// wrap startX into [-bw, 0)
		for startX > 0 {
			startX -= float64(bw)
		}
		for startX < -float64(bw) {
			startX += float64(bw)
		}
	}
	if bg.VTiled {
		for startY > 0 {
			startY -= float64(bh)
		}
		for startY < -float64(bh) {
			startY += float64(bh)
		}
	}

	xEnd := float64(portW)
	yEnd := float64(portH)

	for dy := startY; dy < yEnd; dy += float64(bh) {
		for dx := startX; dx < xEnd; dx += float64(bw) {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(scaleX, scaleY)
			tx := math.Round(portX + dx*scaleX)
			ty := math.Round(portY + dy*scaleY)
			op.GeoM.Translate(tx, ty)
			screen.DrawImage(img, op)
			if !bg.HTiled {
				break
			}
		}
		if !bg.VTiled {
			break
		}
	}
}

func (g *Game) drawTiles(screen *ebiten.Image, room *Room,
	viewX, viewY, scaleX, scaleY, portX, portY float64) {
	for _, tile := range room.Tiles {
		bgImg := g.AssetManager.GetBackground(tile.BGName)
		if bgImg == nil {
			continue
		}

		// get the sub-image region from the background tileset
		subRect := image.Rect(tile.XO, tile.YO, tile.XO+tile.W, tile.YO+tile.H)
		subImg := bgImg.SubImage(subRect).(*ebiten.Image)

		// position in room -> screen via view transform (rounded to fix pixel aliasing)
		scrX := math.Round(portX + (tile.X-viewX)*scaleX)
		scrY := math.Round(portY + (tile.Y-viewY)*scaleY)

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(tile.ScaleX*scaleX, tile.ScaleY*scaleY)
		op.GeoM.Translate(scrX, scrY)
		screen.DrawImage(subImg, op)
	}
}

// layout implements ebiten.Game interface
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

// drawFinalScreen allows custom filtering when Ebitengine scales the logical screen to the window.
// using FilterLinear fixes "shimmering" and uneven pixel sizes (aliasing) on high-DPI/4K monitors
// where the screen scale is not a perfect integer.
func (g *Game) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	screen.Clear()
	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.GeoM = geoM
	screen.DrawImage(offscreen, op)
}

// gotoRoom changes the current room
func (g *Game) GotoRoom(name string) {
	// destroy all non-persistent instances
	g.InstanceMgr.DestroyNonPersistent()

	room := g.RoomManager.Get(name)
	if room == nil {
		fmt.Printf("WARNING: room '%s' not found\n", name)
		return
	}

	g.CurrentRoom = name
	g.RoomManager.SetCurrent(name)

	if room.Speed > 0 {
		g.GameSpeed = room.Speed
		g.tickInterval = 1.0 / float64(g.GameSpeed)
	}

	// spawn room instances
	for _, inst := range room.Instances {
		g.InstanceMgr.CreateFromRoom(inst)
	}

	// room-specific tune objects control bgm.
	if music := roomMusicFromRoom(room); music != "" {
		g.AudioMgr.PlayMusic(music)
	}

	fmt.Printf("Entered room: %s (%dx%d, %d instances)\n",
		name, room.Width, room.Height, len(room.Instances))
}

func roomMusicFromRoom(room *Room) string {
	if room == nil {
		return ""
	}
	tuneToSound := map[string]string{
		"Main_Menu_Music":      "Main_Menu0",
		"Real_MM_Music":        "Main_Menu0",
		"Meadows_Tune":         "Meadows_Music",
		"Intense_Desert_music": "Desert_Music",
		"Monkey_Docks_music":   "Docks_Music",
		"Swamp_tune":           "Swamp_Music",
		"Fort_Tune":            "Fort_Music",
		"Depths_tube":          "Depths_Music",
		"Sun_Tune":             "Sun_Music",
		"Woods_Tune":           "Woods_Music",
		"Minecart_Tune":        "Minecarts_Music",
		"Crimson_Tune":         "Crimson_Music",
		"Xtreme_Tune":          "Xtreme_Music",
		"Space_Tune":           "Space_Music",
		"River_Tune":           "River_Music",
		"Factory_Tune":         "Factory_Music",
		"Sky_Tune":             "Sky_Theme",
		"Lava_Tune":            "Lava_Theme",
		"Boss_Tune":            "Boss_Music",
		"Tim_Tune":             "Birthday_Theme",
	}
	for _, inst := range room.Instances {
		if s, ok := tuneToSound[inst.ObjName]; ok {
			return s
		}
	}

	// menu/front-end fallback.
	name := room.Name
	if name == "Start_Screen" || name == "Main_Menu" ||
		name == "Track_Select_I" || name == "Track_Setup_II" ||
		strings.Contains(name, "Bloons_Bounty") {
		return "Main_Menu0"
	}
	return ""
}

// requestRoomGoto schedules a room change at end of tick (safe to call mid-event).
// if a room change is already pending, the first request wins (prevents
// multiple buttons firing in the same tick from overwriting each other).
func (g *Game) RequestRoomGoto(name string) {
	if g.pendingRoom == "" {
		g.pendingRoom = name
	}
}

// setGameSpeed updates the game speed and recalculates the tick interval.
func (g *Game) SetGameSpeed(speed int) {
	g.GameSpeed = speed
	g.tickInterval = 1.0 / float64(speed)
}

// gotoNextRoom goes to the next room in project order
func (g *Game) GotoNextRoom() {
	order := g.RoomManager.GetRoomOrder()
	for i, name := range order {
		if name == g.CurrentRoom && i+1 < len(order) {
			g.RequestRoomGoto(order[i+1])
			return
		}
	}
}

// isMouseOverInstance checks if the mouse is over an instance's sprite bbox
func (g *Game) IsMouseOverInstance(inst *Instance) bool {
	room := g.RoomManager.GetCurrent()
	view := g.getActiveView(room)
	mx := float64(g.InputMgr.MouseX)
	my := float64(g.InputMgr.MouseY)
	roomX := float64(view.XView) + (mx-float64(view.XPort))*float64(view.WView)/float64(view.WPort)
	roomY := float64(view.YView) + (my-float64(view.YPort))*float64(view.HView)/float64(view.HPort)

	if inst.SpriteName == "" {
		return false // no sprite = not clickable
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil {
		return false // sprite not loaded = not clickable
	}
	left := inst.X - float64(spr.XOrigin)*inst.ImageXScale + float64(spr.BBox.Left)*inst.ImageXScale
	top := inst.Y - float64(spr.YOrigin)*inst.ImageYScale + float64(spr.BBox.Top)*inst.ImageYScale
	right := inst.X - float64(spr.XOrigin)*inst.ImageXScale + float64(spr.BBox.Right+1)*inst.ImageXScale
	bottom := inst.Y - float64(spr.YOrigin)*inst.ImageYScale + float64(spr.BBox.Bottom+1)*inst.ImageYScale
	return roomX >= left && roomX <= right && roomY >= top && roomY <= bottom
}

// getMouseRoomPos returns the mouse position in room coordinates
func (g *Game) GetMouseRoomPos() (float64, float64) {
	room := g.RoomManager.GetCurrent()
	view := g.getActiveView(room)
	mx := float64(g.InputMgr.MouseX)
	my := float64(g.InputMgr.MouseY)
	roomX := float64(view.XView) + (mx-float64(view.XPort))*float64(view.WView)/float64(view.WPort)
	roomY := float64(view.YView) + (my-float64(view.YPort))*float64(view.HView)/float64(view.HPort)
	return roomX, roomY
}

// loadObjectMetadata loads sprite/depth/visible/persistent from objects.json
func (g *Game) LoadObjectMetadata(assetsDir string) error {
	path := filepath.Join(assetsDir, "data", "objects.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var objects []struct {
		Name       string `json:"Name"`
		SpriteName string `json:"SpriteName"`
		Depth      int    `json:"Depth"`
		Visible    bool   `json:"Visible"`
		Persistent bool   `json:"Persistent"`
	}
	if err := json.Unmarshal(data, &objects); err != nil {
		return err
	}
	for _, o := range objects {
		if o.SpriteName != "" {
			g.InstanceMgr.objectSprites[o.Name] = o.SpriteName
		}
		g.InstanceMgr.objectDepths[o.Name] = o.Depth
		g.InstanceMgr.objectVisible[o.Name] = o.Visible
		g.InstanceMgr.objectPersist[o.Name] = o.Persistent
	}
	fmt.Printf("Loaded metadata for %d objects\n", len(objects))
	return nil
}

// colorFromBGR converts a BGR color integer to RGB
// stored as BGR (blue * 65536 + green * 256 + red)
func colorFromBGR(c uint32) (r, g, b uint8) {
	return uint8(c & 0xFF), uint8((c >> 8) & 0xFF), uint8((c >> 16) & 0xFF)
}

// run starts the game
func (g *Game) Run() error {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Bloons Tower Defense X")
	ebiten.SetWindowClosingHandled(false) // ensure X button closes window
	ebiten.SetTPS(60)
	return ebiten.RunGame(g)
}
