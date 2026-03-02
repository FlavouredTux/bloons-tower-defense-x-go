package engine

import (
	"math"
	"sort"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// instance represents a game object instance
type Instance struct {
	// identity
	ID         int
	ObjectName string

	// position and motion
	X, Y               float64
	XPrevious, YPrevious float64
	XStart, YStart     float64
	HSpeed, VSpeed     float64
	Speed, Direction   float64
	Gravity            float64
	GravityDirection   float64
	Friction           float64

	// rendering
	Visible        bool
	Depth          int
	SpriteName     string
	ImageIndex     float64
	ImageSpeed     float64
	ImageXScale    float64
	ImageYScale    float64
	ImageAngle     float64
	ImageAlpha     float64
	ImageBlend     uint32

	// collision
	Solid      bool
	MaskName   string

	// state
	Active     bool
	Persistent bool
	Destroyed  bool
	Alarms     [12]int // 12 alarm slots

	// custom instance variables
	Vars map[string]interface{}

	// reference to the game engine (set when created)
	GameRef *Game

	// path following
	PathIndex    int
	PathPosition float64
	PathSpeed    float64
	PathScale    float64
	PathOrientation float64
	PathEndAction   int // 0=stop, 1=restart, 2=reverse, 3=continue

	// event handlers (set by the object type)
	Behavior InstanceBehavior
}

// instanceBehavior defines the event handlers for an object type.
// each object type implements this interface.
// all methods receive the Game pointer for access to engine systems.
type InstanceBehavior interface {
	Create(inst *Instance, g *Game)
	Destroy(inst *Instance, g *Game)
	Step(inst *Instance, g *Game)
	StepBegin(inst *Instance, g *Game)
	StepEnd(inst *Instance, g *Game)
	Alarm(inst *Instance, alarmIndex int, g *Game)
	Draw(inst *Instance, screen *ebiten.Image, g *Game)
	Collision(inst *Instance, other *Instance, g *Game)
	// mouse events
	MouseLeftPressed(inst *Instance, g *Game)
	MouseRightPressed(inst *Instance, g *Game)
	MouseGlobalLeftPressed(inst *Instance, g *Game)
	MouseGlobalLeftReleased(inst *Instance, g *Game)
	// keyboard events
	KeyPress(inst *Instance, g *Game)
}

// defaultBehavior provides no-op implementations for all events
type DefaultBehavior struct{}

func (d *DefaultBehavior) Create(inst *Instance, g *Game)                              {}
func (d *DefaultBehavior) Destroy(inst *Instance, g *Game)                             {}
func (d *DefaultBehavior) Step(inst *Instance, g *Game)                                {}
func (d *DefaultBehavior) StepBegin(inst *Instance, g *Game)                           {}
func (d *DefaultBehavior) StepEnd(inst *Instance, g *Game)                             {}
func (d *DefaultBehavior) Alarm(inst *Instance, alarmIndex int, g *Game)               {}
func (d *DefaultBehavior) Collision(inst *Instance, other *Instance, g *Game)          {}
func (d *DefaultBehavior) MouseLeftPressed(inst *Instance, g *Game)                    {}
func (d *DefaultBehavior) MouseRightPressed(inst *Instance, g *Game)                   {}
func (d *DefaultBehavior) MouseGlobalLeftPressed(inst *Instance, g *Game)              {}
func (d *DefaultBehavior) MouseGlobalLeftReleased(inst *Instance, g *Game)             {}
func (d *DefaultBehavior) KeyPress(inst *Instance, g *Game)                            {}
func (d *DefaultBehavior) Draw(inst *Instance, screen *ebiten.Image, g *Game) {
	// default draw: render the sprite at the instance position
	if inst.SpriteName == "" {
		return
	}
	spr := g.AssetManager.GetSprite(inst.SpriteName)
	if spr == nil {
		return
	}
	if len(spr.Frames) == 0 {
		return
	}
	frameIdx := int(inst.ImageIndex) % len(spr.Frames)
	if frameIdx < 0 {
		frameIdx = 0
	}
	frame := spr.Frames[frameIdx]
	if frame == nil {
		return
	}

	DrawSpriteExt(screen, frame, spr.XOrigin, spr.YOrigin, inst.X, inst.Y,
		inst.ImageXScale, inst.ImageYScale, inst.ImageAngle, inst.ImageAlpha)
}

// drawSpriteExt draws a sprite frame with full transform (used by default draw and custom draws).
func DrawSpriteExt(screen, frame *ebiten.Image, xOrigin, yOrigin int,
	x, y, xScale, yScale, angle, alpha float64) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(xOrigin), -float64(yOrigin))
	op.GeoM.Scale(xScale, yScale)
	if angle != 0 {
		op.GeoM.Rotate(-angle * math.Pi / 180)
	}
	// round to integer coordinates to prevent "fractional aliasing" / shimmering
	// during nearest-neighbor subpixel rendering
	op.GeoM.Translate(math.Round(x), math.Round(y))
	op.ColorScale.ScaleAlpha(float32(alpha))
	screen.DrawImage(frame, op)
}

// newInstance creates a new instance with default values
func NewInstance(id int, objectName string, x, y float64) *Instance {
	return &Instance{
		ID:          id,
		ObjectName:  objectName,
		X:           x,
		Y:           y,
		XStart:      x,
		YStart:      y,
		XPrevious:   x,
		YPrevious:   y,
		Visible:     true,
		Active:      true,
		ImageSpeed:  0,
		ImageXScale: 1.0,
		ImageYScale: 1.0,
		ImageAlpha:  1.0,
		ImageBlend:  0xFFFFFF,
		PathScale:   1.0,
		Vars:        make(map[string]interface{}),
		Behavior:    &DefaultBehavior{},
	}
}

// onCreate fires the Create event
func (inst *Instance) OnCreate() {
	if inst.Behavior != nil {
		inst.Behavior.Create(inst, inst.GameRef)
	}
}

// onDestroy fires the Destroy event
func (inst *Instance) OnDestroy() {
	if inst.Behavior != nil {
		inst.Behavior.Destroy(inst, inst.GameRef)
	}
}

// onStep fires the Step event
func (inst *Instance) OnStep() {
	if inst.Behavior != nil {
		inst.Behavior.Step(inst, inst.GameRef)
	}
	// apply built-in motion
	inst.applyMotion()
}

// onStepBegin fires the Begin Step event
func (inst *Instance) OnStepBegin() {
	inst.XPrevious = inst.X
	inst.YPrevious = inst.Y
	if inst.Behavior != nil {
		inst.Behavior.StepBegin(inst, inst.GameRef)
	}
}

// onStepEnd fires the End Step event
func (inst *Instance) OnStepEnd() {
	if inst.Behavior != nil {
		inst.Behavior.StepEnd(inst, inst.GameRef)
	}
}

// onAlarm fires an Alarm event
func (inst *Instance) OnAlarm(index int) {
	if inst.Behavior != nil {
		inst.Behavior.Alarm(inst, index, inst.GameRef)
	}
}

// onDraw fires the Draw event
func (inst *Instance) OnDraw(screen *ebiten.Image) {
	if inst.Behavior != nil {
		inst.Behavior.Draw(inst, screen, inst.GameRef)
	}
}

// onCollision fires a Collision event with another instance
func (inst *Instance) OnCollision(other *Instance) {
	if inst.Behavior != nil {
		inst.Behavior.Collision(inst, other, inst.GameRef)
	}
}

// advanceAnimation advances the sprite animation by image_speed
func (inst *Instance) AdvanceAnimation() {
	inst.ImageIndex += inst.ImageSpeed
}

// applyMotion applies speed/direction, gravity, friction
func (inst *Instance) applyMotion() {
	// apply gravity
	if inst.Gravity != 0 {
		gRad := inst.GravityDirection * math.Pi / 180
		inst.HSpeed += inst.Gravity * math.Cos(gRad)
		inst.VSpeed -= inst.Gravity * math.Sin(gRad)
	}

	// apply friction
	if inst.Friction != 0 && (inst.HSpeed != 0 || inst.VSpeed != 0) {
		spd := math.Sqrt(inst.HSpeed*inst.HSpeed + inst.VSpeed*inst.VSpeed)
		newSpd := spd - inst.Friction
		if newSpd < 0 {
			newSpd = 0
		}
		if spd > 0 {
			ratio := newSpd / spd
			inst.HSpeed *= ratio
			inst.VSpeed *= ratio
		}
	}

	// update speed/direction from hspeed/vspeed
	inst.Speed = math.Sqrt(inst.HSpeed*inst.HSpeed + inst.VSpeed*inst.VSpeed)
	if inst.Speed > 0 {
		inst.Direction = math.Atan2(-inst.VSpeed, inst.HSpeed) * 180 / math.Pi
	}

	// move
	inst.X += inst.HSpeed
	inst.Y += inst.VSpeed
}

// motionSet sets speed and direction
func (inst *Instance) MotionSet(dir, spd float64) {
	inst.Direction = dir
	inst.Speed = spd
	rad := dir * math.Pi / 180
	inst.HSpeed = spd * math.Cos(rad)
	inst.VSpeed = -spd * math.Sin(rad)
}

// instanceManager manages all active instances
type InstanceManager struct {
	mu         sync.RWMutex
	instances  map[int]*Instance
	nextID     int
	toDestroy  []int
	gameRef    *Game

	// object type registry - maps object name to behavior factory
	behaviors  map[string]func() InstanceBehavior

	// object sprite mapping (from objects.json)
	objectSprites map[string]string
	objectDepths  map[string]int
	objectVisible map[string]bool
	objectPersist map[string]bool
}

func NewInstanceManager() *InstanceManager {
	return &InstanceManager{
		instances:     make(map[int]*Instance),
		nextID:        100000,
		behaviors:     make(map[string]func() InstanceBehavior),
		objectSprites: make(map[string]string),
		objectDepths:  make(map[string]int),
		objectVisible: make(map[string]bool),
		objectPersist: make(map[string]bool),
	}
}

// setGameRef sets the game reference for all new instances
func (im *InstanceManager) SetGameRef(g *Game) {
	im.gameRef = g
}

// registerBehavior registers a behavior factory for an object type
func (im *InstanceManager) RegisterBehavior(objectName string, factory func() InstanceBehavior) {
	im.behaviors[objectName] = factory
}

// create creates a new instance at the given position
func (im *InstanceManager) Create(objectName string, x, y float64) *Instance {
	im.mu.Lock()
	inst := NewInstance(im.nextID, objectName, x, y)
	inst.GameRef = im.gameRef
	im.nextID++

	// set sprite/depth/visible from object metadata
	if spr, ok := im.objectSprites[objectName]; ok && spr != "" {
		inst.SpriteName = spr
	}
	if d, ok := im.objectDepths[objectName]; ok {
		inst.Depth = d
	}
	if v, ok := im.objectVisible[objectName]; ok {
		inst.Visible = v
	}
	if p, ok := im.objectPersist[objectName]; ok {
		inst.Persistent = p
	}

	// assign behavior from registry
	if factory, ok := im.behaviors[objectName]; ok {
		inst.Behavior = factory()
	}

	im.instances[inst.ID] = inst
	im.mu.Unlock()

	// run Create after releasing the manager lock so behaviors can safely
	// query/create instances during their Create event without deadlocking.
	inst.OnCreate()
	return inst
}

// createFromRoom creates an instance from room data
func (im *InstanceManager) CreateFromRoom(ri RoomInstanceDef) *Instance {
	inst := im.Create(ri.ObjName, ri.X, ri.Y)
	inst.ImageXScale = ri.ScaleX
	inst.ImageYScale = ri.ScaleY
	inst.ImageAngle = ri.Rotation
	if ri.ScaleX == 0 {
		inst.ImageXScale = 1
	}
	if ri.ScaleY == 0 {
		inst.ImageYScale = 1
	}
	return inst
}

// destroy marks an instance for destruction
func (im *InstanceManager) Destroy(id int) {
	im.mu.Lock()
	defer im.mu.Unlock()

	if inst, ok := im.instances[id]; ok {
		inst.Destroyed = true
		inst.OnDestroy()
		im.toDestroy = append(im.toDestroy, id)
	}
}

// flushDestroyed removes all destroyed instances
func (im *InstanceManager) FlushDestroyed() {
	im.mu.Lock()
	defer im.mu.Unlock()

	for _, id := range im.toDestroy {
		delete(im.instances, id)
	}
	im.toDestroy = im.toDestroy[:0]
}

// getAll returns all active instances
func (im *InstanceManager) GetAll() []*Instance {
	im.mu.RLock()
	defer im.mu.RUnlock()

	result := make([]*Instance, 0, len(im.instances))
	for _, inst := range im.instances {
		result = append(result, inst)
	}
	return result
}

// getSortedByDepth returns instances sorted by depth (descending - higher depth drawn first)
func (im *InstanceManager) GetSortedByDepth() []*Instance {
	all := im.GetAll()
	sort.Slice(all, func(i, j int) bool {
		return all[i].Depth > all[j].Depth
	})
	return all
}

// getByID returns an instance by ID, or nil if not found
func (im *InstanceManager) GetByID(id int) *Instance {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.instances[id]
}

// findByObject returns all instances of a given object type
func (im *InstanceManager) FindByObject(objectName string) []*Instance {
	im.mu.RLock()
	defer im.mu.RUnlock()

	var result []*Instance
	for _, inst := range im.instances {
		if inst.ObjectName == objectName && !inst.Destroyed {
			result = append(result, inst)
		}
	}
	return result
}

// instanceCount returns the number of instances of an object type
func (im *InstanceManager) InstanceCount(objectName string) int {
	return len(im.FindByObject(objectName))
}

// instanceExists checks if any instance of the given object type exists
func (im *InstanceManager) InstanceExists(objectName string) bool {
	return im.InstanceCount(objectName) > 0
}

// objectSpriteName returns the default sprite mapped for an object name.
func (im *InstanceManager) ObjectSpriteName(objectName string) string {
	im.mu.RLock()
	defer im.mu.RUnlock()
	return im.objectSprites[objectName]
}

// destroyNonPersistent destroys all non-persistent instances (for room transitions)
func (im *InstanceManager) DestroyNonPersistent() {
	im.mu.Lock()
	defer im.mu.Unlock()

	for id, inst := range im.instances {
		if !inst.Persistent {
			inst.Destroyed = true
			inst.OnDestroy()
			im.toDestroy = append(im.toDestroy, id)
		}
	}
}

// checkCollisions checks collisions between instances using bounding boxes
func (im *InstanceManager) CheckCollisions() {
	all := im.GetAll()
	// simple O(n²) for now - can optimize with spatial partitioning later
	for i := 0; i < len(all); i++ {
		if !all[i].Active || all[i].Destroyed {
			continue
		}
		for j := i + 1; j < len(all); j++ {
			if !all[j].Active || all[j].Destroyed {
				continue
			}
			// tODO: implement proper bbox collision using sprite data
			// for now this is a placeholder
		}
	}
}

// roomInstanceDef defines an instance placement in a room
type RoomInstanceDef struct {
	ObjName  string
	X, Y     float64
	ScaleX   float64
	ScaleY   float64
	Rotation float64
	Code     string
}
