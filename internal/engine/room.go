package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// room represents a loaded game room
type Room struct {
	Name        string
	Width       int
	Height      int
	Speed       int
	Persistent  bool
	Code        string
	BgColor     uint32
	ShowBgColor bool
	EnableViews bool
	Backgrounds []RoomBackground
	Views       []RoomView
	Instances   []RoomInstanceDef
	Tiles       []RoomTileDef
}

type RoomBackground struct {
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

type RoomView struct {
	Visible                            bool
	ObjName                            string
	XView, YView, WView, HView         int
	XPort, YPort, WPort, HPort         int
	HBorder, VBorder, HSpeed, VSpeed   int
}

type RoomTileDef struct {
	BGName       string
	X, Y         float64
	W, H         int
	XO, YO       int
	ID           int
	Depth        int
	ScaleX       float64
	ScaleY       float64
}

// roomManager handles room loading and transitions
type RoomManager struct {
	rooms       map[string]*Room
	roomOrder   []string
	currentRoom string
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]*Room),
	}
}

// loadRoomsFromJSON loads all rooms from the extracted JSON manifest
func (rm *RoomManager) LoadRoomsFromJSON(assetsDir string) error {
	path := filepath.Join(assetsDir, "data", "rooms.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading rooms.json: %w", err)
	}

	var roomDefs []RoomJSON
	if err := json.Unmarshal(data, &roomDefs); err != nil {
		return fmt.Errorf("parsing rooms.json: %w", err)
	}

	for _, rd := range roomDefs {
		room := &Room{
			Name:        rd.Name,
			Width:       rd.Width,
			Height:      rd.Height,
			Speed:       rd.Speed,
			Persistent:  rd.Persistent,
			Code:        rd.Code,
			BgColor:     rd.BgColor,
			ShowBgColor: rd.ShowBgColor,
			EnableViews: rd.EnableViews,
		}

		for _, bg := range rd.Backgrounds {
			room.Backgrounds = append(room.Backgrounds, RoomBackground{
				Visible:    bg.Visible,
				Foreground: bg.Foreground,
				Name:       bg.Name,
				X:          bg.X,
				Y:          bg.Y,
				HTiled:     bg.HTiled,
				VTiled:     bg.VTiled,
				HSpeed:     bg.HSpeed,
				VSpeed:     bg.VSpeed,
				Stretch:    bg.Stretch,
			})
		}

		for _, v := range rd.Views {
			room.Views = append(room.Views, RoomView{
				Visible:  v.Visible,
				ObjName:  v.ObjName,
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

		for _, inst := range rd.Instances {
			room.Instances = append(room.Instances, RoomInstanceDef{
				ObjName:  inst.ObjName,
				X:        inst.X,
				Y:        inst.Y,
				ScaleX:   inst.ScaleX,
				ScaleY:   inst.ScaleY,
				Rotation: inst.Rotation,
				Code:     inst.Code,
			})
		}

		for _, t := range rd.Tiles {
			room.Tiles = append(room.Tiles, RoomTileDef{
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
			})
		}

		rm.rooms[room.Name] = room
		rm.roomOrder = append(rm.roomOrder, room.Name)
	}

	fmt.Printf("Loaded %d rooms\n", len(rm.rooms))
	return nil
}

// get returns a room by name
func (rm *RoomManager) Get(name string) *Room {
	return rm.rooms[name]
}

// getCurrent returns the current room
func (rm *RoomManager) GetCurrent() *Room {
	return rm.rooms[rm.currentRoom]
}

// setCurrent sets the current room
func (rm *RoomManager) SetCurrent(name string) {
	rm.currentRoom = name
}

// getRoomOrder returns the list of rooms in order
func (rm *RoomManager) GetRoomOrder() []string {
	return rm.roomOrder
}

// getRoomNames returns all room names
func (rm *RoomManager) GetRoomNames() []string {
	names := make([]string, 0, len(rm.rooms))
	for name := range rm.rooms {
		names = append(names, name)
	}
	return names
}

// jSON structures for deserialization
type RoomJSON struct {
	Name        string          `json:"Name"`
	Width       int             `json:"Width"`
	Height      int             `json:"Height"`
	Speed       int             `json:"Speed"`
	Persistent  bool            `json:"Persistent"`
	Code        string          `json:"Code"`
	BgColor     uint32          `json:"BgColor"`
	ShowBgColor bool            `json:"ShowBgColor"`
	EnableViews bool            `json:"EnableViews"`
	Backgrounds []RoomBGJSON    `json:"Backgrounds"`
	Views       []RoomViewJSON  `json:"Views"`
	Instances   []RoomInstJSON  `json:"Instances"`
	Tiles       []RoomTileJSON  `json:"Tiles"`
}

type RoomBGJSON struct {
	Visible    bool   `json:"Visible"`
	Foreground bool   `json:"Foreground"`
	Name       string `json:"Name"`
	X          int    `json:"X"`
	Y          int    `json:"Y"`
	HTiled     bool   `json:"HTiled"`
	VTiled     bool   `json:"VTiled"`
	HSpeed     int    `json:"HSpeed"`
	VSpeed     int    `json:"VSpeed"`
	Stretch    bool   `json:"Stretch"`
}

type RoomViewJSON struct {
	Visible  bool   `json:"Visible"`
	ObjName  string `json:"ObjName"`
	XView    int    `json:"XView"`
	YView    int    `json:"YView"`
	WView    int    `json:"WView"`
	HView    int    `json:"HView"`
	XPort    int    `json:"XPort"`
	YPort    int    `json:"YPort"`
	WPort    int    `json:"WPort"`
	HPort    int    `json:"HPort"`
	HBorder  int    `json:"HBorder"`
	VBorder  int    `json:"VBorder"`
	HSpeed   int    `json:"HSpeed"`
	VSpeed   int    `json:"VSpeed"`
}

type RoomInstJSON struct {
	ObjName  string  `json:"ObjName"`
	X        float64 `json:"X"`
	Y        float64 `json:"Y"`
	Name     string  `json:"Name"`
	Code     string  `json:"Code"`
	ScaleX   float64 `json:"ScaleX"`
	ScaleY   float64 `json:"ScaleY"`
	Colour   uint32  `json:"Colour"`
	Rotation float64 `json:"Rotation"`
}

type RoomTileJSON struct {
	BGName string  `json:"BGName"`
	X      float64 `json:"X"`
	Y      float64 `json:"Y"`
	W      int     `json:"W"`
	H      int     `json:"H"`
	XO     int     `json:"XO"`
	YO     int     `json:"YO"`
	ID     int     `json:"ID"`
	Depth  int     `json:"Depth"`
	ScaleX float64 `json:"ScaleX"`
	ScaleY float64 `json:"ScaleY"`
}
