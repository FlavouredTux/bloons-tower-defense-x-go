package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// RoomGMX represents a GameMaker .room.gmx file
type RoomGMX struct {
	XMLName             xml.Name         `xml:"room"`
	Caption             string           `xml:"caption"`
	Width               int              `xml:"width"`
	Height              int              `xml:"height"`
	VSnap               int              `xml:"vsnap"`
	HSnap               int              `xml:"hsnap"`
	Isometric           int              `xml:"isometric"`
	Speed               int              `xml:"speed"`
	Persistent          int              `xml:"persistent"`
	Colour              uint32           `xml:"colour"`
	ShowColour          int              `xml:"showcolour"`
	Code                string           `xml:"code"`
	EnableViews         int              `xml:"enableViews"`
	ClearViewBackground int              `xml:"clearViewBackground"`
	ClearDisplayBuffer  int              `xml:"clearDisplayBuffer"`
	Backgrounds         RoomBackgrounds  `xml:"backgrounds"`
	Views               RoomViews        `xml:"views"`
	Instances           RoomInstances    `xml:"instances"`
	Tiles               RoomTiles        `xml:"tiles"`
}

type RoomBackgrounds struct {
	Backgrounds []RoomBG `xml:"background"`
}

type RoomBG struct {
	Visible    int    `xml:"visible,attr"`
	Foreground int    `xml:"foreground,attr"`
	Name       string `xml:"name,attr"`
	X          int    `xml:"x,attr"`
	Y          int    `xml:"y,attr"`
	HTiled     int    `xml:"htiled,attr"`
	VTiled     int    `xml:"vtiled,attr"`
	HSpeed     int    `xml:"hspeed,attr"`
	VSpeed     int    `xml:"vspeed,attr"`
	Stretch    int    `xml:"stretch,attr"`
}

type RoomViews struct {
	Views []RoomView `xml:"view"`
}

type RoomView struct {
	Visible  int    `xml:"visible,attr"`
	ObjName  string `xml:"objName,attr"`
	XView    int    `xml:"xview,attr"`
	YView    int    `xml:"yview,attr"`
	WView    int    `xml:"wview,attr"`
	HView    int    `xml:"hview,attr"`
	XPort    int    `xml:"xport,attr"`
	YPort    int    `xml:"yport,attr"`
	WPort    int    `xml:"wport,attr"`
	HPort    int    `xml:"hport,attr"`
	HBorder  int    `xml:"hborder,attr"`
	VBorder  int    `xml:"vborder,attr"`
	HSpeed   int    `xml:"hspeed,attr"`
	VSpeed   int    `xml:"vspeed,attr"`
}

type RoomInstances struct {
	Instances []RoomInstance `xml:"instance"`
}

type RoomInstance struct {
	ObjName  string  `xml:"objName,attr"`
	X        float64 `xml:"x,attr"`
	Y        float64 `xml:"y,attr"`
	Name     string  `xml:"name,attr"`
	Locked   int     `xml:"locked,attr"`
	Code     string  `xml:"code,attr"`
	ScaleX   float64 `xml:"scaleX,attr"`
	ScaleY   float64 `xml:"scaleY,attr"`
	Colour   uint32  `xml:"colour,attr"`
	Rotation float64 `xml:"rotation,attr"`
}

type RoomTiles struct {
	Tiles []RoomTile `xml:"tile"`
}

type RoomTile struct {
	BGName string  `xml:"bgName,attr"`
	X      float64 `xml:"x,attr"`
	Y      float64 `xml:"y,attr"`
	W      int     `xml:"w,attr"`
	H      int     `xml:"h,attr"`
	XO     int     `xml:"xo,attr"`
	YO     int     `xml:"yo,attr"`
	ID     int     `xml:"id,attr"`
	Name   string  `xml:"name,attr"`
	Depth  int     `xml:"depth,attr"`
	Locked int     `xml:"locked,attr"`
	ScaleX float64 `xml:"scaleX,attr"`
	ScaleY float64 `xml:"scaleY,attr"`
	Colour uint32  `xml:"colour,attr"`
}

// RoomData is our intermediate representation
type RoomData struct {
	Name        string
	Width       int
	Height      int
	Speed       int
	Persistent  bool
	Code        string
	BgColor     uint32
	ShowBgColor bool
	EnableViews bool
	Backgrounds []RoomBGData
	Views       []RoomViewData
	Instances   []RoomInstanceData
	Tiles       []RoomTileData
}

type RoomBGData struct {
	Visible    bool
	Foreground bool
	Name       string
	X, Y       int
	HTiled     bool
	VTiled     bool
	HSpeed     int
	VSpeed     int
	Stretch    bool
}

type RoomViewData struct {
	Visible                            bool
	ObjName                            string
	XView, YView, WView, HView         int
	XPort, YPort, WPort, HPort         int
	HBorder, VBorder, HSpeed, VSpeed   int
}

type RoomInstanceData struct {
	ObjName  string
	X, Y     float64
	Name     string
	Code     string
	ScaleX   float64
	ScaleY   float64
	Colour   uint32
	Rotation float64
}

type RoomTileData struct {
	BGName       string
	X, Y         float64
	W, H         int
	XO, YO       int
	ID           int
	Depth        int
	ScaleX       float64
	ScaleY       float64
	Colour       uint32
}

// ParseRoom parses a .room.gmx file
func ParseRoom(path string) (*RoomData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading room gmx %s: %w", path, err)
	}

	var gmx RoomGMX
	if err := xml.Unmarshal(data, &gmx); err != nil {
		return nil, fmt.Errorf("parsing room gmx %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".room.gmx")

	room := &RoomData{
		Name:        name,
		Width:       gmx.Width,
		Height:      gmx.Height,
		Speed:       gmx.Speed,
		Persistent:  gmx.Persistent != 0,
		Code:        gmx.Code,
		BgColor:     gmx.Colour,
		ShowBgColor: gmx.ShowColour != 0,
		EnableViews: gmx.EnableViews != 0,
	}

	for _, bg := range gmx.Backgrounds.Backgrounds {
		room.Backgrounds = append(room.Backgrounds, RoomBGData{
			Visible:    bg.Visible != 0,
			Foreground: bg.Foreground != 0,
			Name:       bg.Name,
			X:          bg.X,
			Y:          bg.Y,
			HTiled:     bg.HTiled != 0,
			VTiled:     bg.VTiled != 0,
			HSpeed:     bg.HSpeed,
			VSpeed:     bg.VSpeed,
			Stretch:    bg.Stretch != 0,
		})
	}

	for _, v := range gmx.Views.Views {
		room.Views = append(room.Views, RoomViewData{
			Visible:  v.Visible != 0,
			ObjName:  cleanGMXString(v.ObjName),
			XView:    v.XView,
			YView:    v.YView,
			WView:    v.WView,
			HView:    v.HView,
			XPort:    v.XPort,
			YPort:    v.YPort,
			WPort:    v.WPort,
			HPort:    v.HPort,
			HBorder:  v.HBorder,
			VBorder:  v.VBorder,
			HSpeed:   v.HSpeed,
			VSpeed:   v.VSpeed,
		})
	}

	for _, inst := range gmx.Instances.Instances {
		room.Instances = append(room.Instances, RoomInstanceData{
			ObjName:  inst.ObjName,
			X:        inst.X,
			Y:        inst.Y,
			Name:     inst.Name,
			Code:     inst.Code,
			ScaleX:   inst.ScaleX,
			ScaleY:   inst.ScaleY,
			Colour:   inst.Colour,
			Rotation: inst.Rotation,
		})
	}

	for _, t := range gmx.Tiles.Tiles {
		room.Tiles = append(room.Tiles, RoomTileData{
			BGName: t.BGName,
			X:      t.X,
			Y:      t.Y,
			W:      t.W,
			H:      t.H,
			XO:     t.XO,
			YO:     t.YO,
			ID:     t.ID,
			Depth:  t.Depth,
			ScaleX: t.ScaleX,
			ScaleY: t.ScaleY,
			Colour: t.Colour,
		})
	}

	return room, nil
}

// ParseRoomOrderFromProject reads the .project.gmx file and returns room names in project order.
// In GameMaker, the first room listed is the starting room.
func ParseRoomOrderFromProject(projectFile string) ([]string, error) {
	data, err := os.ReadFile(projectFile)
	if err != nil {
		return nil, fmt.Errorf("reading project file: %w", err)
	}

	type ProjectGMX struct {
		XMLName xml.Name `xml:"assets"`
		Rooms   []struct {
			Room  []string `xml:"room"`
			Inner []struct {
				Room []string `xml:"room"`
			} `xml:"rooms"`
		} `xml:"rooms"`
	}

	// The project XML has nested <rooms> groups. We need ALL <room> elements in order.
	// Easiest: just regex/scan for <room>rooms\Name</room> entries.
	var order []string
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "<room>") && strings.HasSuffix(line, "</room>") {
			val := strings.TrimPrefix(line, "<room>")
			val = strings.TrimSuffix(val, "</room>")
			// Value is like "rooms\Start_Screen" - extract just the name
			val = strings.ReplaceAll(val, "\\", "/")
			parts := strings.Split(val, "/")
			name := parts[len(parts)-1]
			order = append(order, name)
		}
	}

	return order, nil
}

// ParseAllRooms parses all room.gmx files in a directory
func ParseAllRooms(roomDir string) ([]*RoomData, error) {
	entries, err := os.ReadDir(roomDir)
	if err != nil {
		return nil, fmt.Errorf("reading room directory: %w", err)
	}

	var rooms []*RoomData
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".room.gmx") {
			continue
		}
		r, err := ParseRoom(filepath.Join(roomDir, e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping room %s: %v\n", e.Name(), err)
			continue
		}
		rooms = append(rooms, r)
	}

	return rooms, nil
}

// ParseAllRoomsOrdered parses all rooms and returns them in the order specified.
// Rooms not in the order list are appended at the end.
func ParseAllRoomsOrdered(roomDir string, order []string) ([]*RoomData, error) {
	rooms, err := ParseAllRooms(roomDir)
	if err != nil {
		return nil, err
	}

	// Build a map for lookup
	roomMap := make(map[string]*RoomData, len(rooms))
	for _, r := range rooms {
		roomMap[r.Name] = r
	}

	// Build ordered result
	var ordered []*RoomData
	seen := make(map[string]bool)
	for _, name := range order {
		if r, ok := roomMap[name]; ok {
			ordered = append(ordered, r)
			seen[name] = true
		}
	}

	// Append any rooms not in the order list (shouldn't happen but be safe)
	for _, r := range rooms {
		if !seen[r.Name] {
			ordered = append(ordered, r)
		}
	}

	return ordered, nil
}

// PathGMX represents a GameMaker .path.gmx file
type PathGMX struct {
	XMLName   xml.Name    `xml:"path"`
	Kind      int         `xml:"kind"`
	Closed    int         `xml:"closed"`
	Precision int         `xml:"precision"`
	BackRoom  int         `xml:"backroom"`
	HSnap     int         `xml:"hsnap"`
	VSnap     int         `xml:"vsnap"`
	Points    PathPoints  `xml:"points"`
}

type PathPoints struct {
	Points []string `xml:"point"`
}

type PathData struct {
	Name      string
	Kind      int // 0=straight, 1=smooth
	Closed    bool
	Precision int
	Points    []PathPoint
}

type PathPoint struct {
	X, Y  float64
	Speed float64
}

// ParsePath parses a .path.gmx file
func ParsePath(path string) (*PathData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading path gmx %s: %w", path, err)
	}

	var gmx PathGMX
	if err := xml.Unmarshal(data, &gmx); err != nil {
		return nil, fmt.Errorf("parsing path gmx %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".path.gmx")

	pd := &PathData{
		Name:      name,
		Kind:      gmx.Kind,
		Closed:    gmx.Closed != 0,
		Precision: gmx.Precision,
	}

	for _, pt := range gmx.Points.Points {
		parts := strings.Split(strings.TrimSpace(pt), ",")
		if len(parts) < 3 {
			continue
		}
		x, _ := strconv.ParseFloat(parts[0], 64)
		y, _ := strconv.ParseFloat(parts[1], 64)
		spd, _ := strconv.ParseFloat(parts[2], 64)
		pd.Points = append(pd.Points, PathPoint{X: x, Y: y, Speed: spd})
	}

	return pd, nil
}

// ParseAllPaths parses all path.gmx files in a directory
func ParseAllPaths(pathDir string) ([]*PathData, error) {
	entries, err := os.ReadDir(pathDir)
	if err != nil {
		return nil, fmt.Errorf("reading path directory: %w", err)
	}

	var paths []*PathData
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".path.gmx") {
			continue
		}
		p, err := ParsePath(filepath.Join(pathDir, e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping path %s: %v\n", e.Name(), err)
			continue
		}
		paths = append(paths, p)
	}

	return paths, nil
}
