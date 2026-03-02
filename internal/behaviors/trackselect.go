package behaviors

import (
	"btdx/internal/engine"
)

// track_Select_I room behaviors
// contains: 30 track goto buttons, track_up, track_down scrolling

// trackGotoBehavior is a generic behavior for all *_goto track buttons.
// each sets global.trackselect to its track number and goes to Track_Setup_II.
type trackGotoBehavior struct {
	engine.DefaultBehavior
	TrackNum float64
}

func (b *trackGotoBehavior) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["trackselect"] = b.TrackNum
	g.RequestRoomGoto("Track_Setup_II")
}

// track_up — scroll track list up (y += 272 on all track_goto)
type TrackUp struct {
	engine.DefaultBehavior
}

func (b *TrackUp) Create(inst *engine.Instance, g *engine.Game) {
	g.GlobalVars["trackpanel"] = 0.0
}

func (b *TrackUp) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	tp := getGlobal(g, "trackpanel")
	if tp > 0 {
		g.GlobalVars["trackpanel"] = tp - 1
		for _, t := range g.InstanceMgr.GetAll() {
			if isTrackGoto(t.ObjectName) {
				t.Y += 272
			}
		}
	}
}

// track_down — scroll track list down (y -= 272 on all track_goto)
type TrackDown struct {
	engine.DefaultBehavior
}

func (b *TrackDown) MouseLeftPressed(inst *engine.Instance, g *engine.Game) {
	tp := getGlobal(g, "trackpanel")
	if tp < 9 {
		g.GlobalVars["trackpanel"] = tp + 1
		for _, t := range g.InstanceMgr.GetAll() {
			if isTrackGoto(t.ObjectName) {
				t.Y -= 272
			}
		}
	}
}

// isTrackGoto checks if an object name is a track goto button
func isTrackGoto(name string) bool {
	for _, n := range trackGotoNames {
		if name == n {
			return true
		}
	}
	return false
}

// trackGotoNames lists all track goto object names
var trackGotoNames = []string{
	"Monkey_Meadows_goto", "Bloon_Oasis_goto", "Diagonal_Swamp_goto",
	"Monkey_Fort_goto", "Monkey_Town_Docks_goto", "Conveyor_Belt_goto",
	"The_Depths_goto", "Sun_Dial_goto", "Shade_Woods_goto",
	"Minecarts_goto", "Crimson_Creek_goto", "Xtreme_Park_goto",
	"Portal_Lab_goto", "Omega_River_goto", "Space_Portals_goto",
	"Bloon_Light_Throwback_goto", "Bloon_Circles_X_goto",
	"Autumn_Acres_goto", "Graveyard_goto", "Village_Defense_goto",
	"Circuit_goto", "Grand_Canyon_goto", "Bloonside_River_goto",
	"Cotton_goto", "Rubber_Rug_goto", "Frozen_Lake_goto",
	"Sky_Battles_goto", "Lava_Stream_goto", "Ravine_River_goto",
	"Peaceful_Bridge_goto",
}

// registerTrackSelectBehaviors registers all Track_Select_I behaviors
func RegisterTrackSelectBehaviors(im *engine.InstanceManager) {
	// track goto buttons — each sets trackselect and goes to Track_Setup_II
	trackGotos := map[string]float64{
		"Monkey_Meadows_goto":        1,
		"Bloon_Oasis_goto":           2,
		"Diagonal_Swamp_goto":        3,
		"Monkey_Fort_goto":           4,
		"Monkey_Town_Docks_goto":     5,
		"Conveyor_Belt_goto":         6,
		"The_Depths_goto":            7,
		"Sun_Dial_goto":              8,
		"Shade_Woods_goto":           9,
		"Minecarts_goto":             10,
		"Crimson_Creek_goto":         11,
		"Xtreme_Park_goto":           12,
		"Portal_Lab_goto":            13,
		"Omega_River_goto":           14,
		"Space_Portals_goto":         15,
		"Bloon_Light_Throwback_goto": 17,
		"Bloon_Circles_X_goto":       18,
		"Autumn_Acres_goto":          19,
		"Graveyard_goto":             20,
		"Village_Defense_goto":       21,
		"Circuit_goto":               22,
		"Grand_Canyon_goto":          23,
		"Bloonside_River_goto":       24,
		"Cotton_goto":                25,
		"Rubber_Rug_goto":            27,
		"Frozen_Lake_goto":           28,
		"Sky_Battles_goto":           29,
		"Lava_Stream_goto":           30,
		"Ravine_River_goto":          31,
		"Peaceful_Bridge_goto":       32,
	}

	for name, trackNum := range trackGotos {
		tn := trackNum // capture for closure
		im.RegisterBehavior(name, func() engine.InstanceBehavior {
			return &trackGotoBehavior{TrackNum: tn}
		})
	}

	// scroll buttons
	im.RegisterBehavior("track_up", func() engine.InstanceBehavior { return &TrackUp{} })
	im.RegisterBehavior("track_down", func() engine.InstanceBehavior { return &TrackDown{} })
}
