package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ObjectGMX represents a GameMaker .object.gmx file
type ObjectGMX struct {
	XMLName    xml.Name    `xml:"object"`
	SpriteName string      `xml:"spriteName"`
	Solid      int         `xml:"solid"`
	Visible    int         `xml:"visible"`
	Depth      int         `xml:"depth"`
	Persistent int         `xml:"persistent"`
	ParentName string      `xml:"parentName"`
	MaskName   string      `xml:"maskName"`
	Events     EventsBlock `xml:"events"`
}

type EventsBlock struct {
	Events []EventGMX `xml:"event"`
}

type EventGMX struct {
	EventType int         `xml:"eventtype,attr"`
	ENum      int         `xml:"enumb,attr"`
	Actions   []ActionGMX `xml:"action"`
}

type ActionGMX struct {
	LibID        int           `xml:"libid"`
	ID           int           `xml:"id"`
	Kind         int           `xml:"kind"`
	UseRelative  int           `xml:"userelative"`
	IsQuestion   int           `xml:"isquestion"`
	UseApplyTo   int           `xml:"useapplyto"`
	ExeType      int           `xml:"exetype"`
	FunctionName string        `xml:"functionname"`
	CodeString   string        `xml:"codestring"`
	WhoName      string        `xml:"whoName"`
	Relative     int           `xml:"relative"`
	IsNot        int           `xml:"isnot"`
	Arguments    ArgumentsGMX  `xml:"arguments"`
}

type ArgumentsGMX struct {
	Arguments []ArgumentGMX `xml:"argument"`
}

type ArgumentGMX struct {
	Kind   int    `xml:"kind"`
	String string `xml:"string"`
	Object string `xml:"object"`
}

// GM Event Types
const (
	EventCreate       = 0
	EventDestroy      = 1
	EventAlarm        = 2
	EventStep         = 3
	EventCollision    = 4
	EventKeyboard     = 5
	EventMouse        = 6
	EventOther        = 7
	EventDraw         = 8
	EventKeyPress     = 9
	EventKeyRelease   = 10
	EventTrigger      = 11
	EventCleanUp      = 12
)

// GM Action Kinds
const (
	ActionKindNormal  = 0
	ActionKindBegin   = 1
	ActionKindEnd     = 2
	ActionKindElse    = 3
	ActionKindExit    = 4
	ActionKindRepeat  = 5
	ActionKindVar     = 6
	ActionKindCode    = 7
)

// ObjectData is our intermediate representation
type ObjectData struct {
	Name       string
	SpriteName string
	Solid      bool
	Visible    bool
	Depth      int
	Persistent bool
	ParentName string
	MaskName   string
	Events     []EventData
}

type EventData struct {
	Type    int    // EventCreate, EventStep, etc.
	SubType int    // alarm number, key code, collision object index, etc.
	Actions []ActionData
}

type ActionData struct {
	Kind         int    // ActionKindCode, ActionKindVar, etc.
	ExeType      int    // 1=function, 2=code
	FunctionName string
	WhoName      string
	Code         string // extracted GML code
	Arguments    []ArgData
	IsQuestion   bool
	UseRelative  bool
}

type ArgData struct {
	Kind   int
	Value  string
	Object string
}

// ParseObject parses a .object.gmx file
func ParseObject(path string) (*ObjectData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading object gmx %s: %w", path, err)
	}

	var gmx ObjectGMX
	if err := xml.Unmarshal(data, &gmx); err != nil {
		return nil, fmt.Errorf("parsing object gmx %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".object.gmx")

	obj := &ObjectData{
		Name:       name,
		SpriteName: cleanGMXString(gmx.SpriteName),
		Solid:      gmx.Solid != 0,
		Visible:    gmx.Visible != 0,
		Depth:      gmx.Depth,
		Persistent: gmx.Persistent != 0,
		ParentName: cleanGMXString(gmx.ParentName),
		MaskName:   cleanGMXString(gmx.MaskName),
	}

	for _, ev := range gmx.Events.Events {
		evData := EventData{
			Type:    ev.EventType,
			SubType: ev.ENum,
		}

		for _, act := range ev.Actions {
			ad := ActionData{
				Kind:         act.Kind,
				ExeType:      act.ExeType,
				FunctionName: act.FunctionName,
				WhoName:      act.WhoName,
				IsQuestion:   act.IsQuestion != 0,
				UseRelative:  act.UseRelative != 0,
			}

			// Extract code from actions
			switch act.Kind {
			case ActionKindCode: // kind=7: inline GML code
				if len(act.Arguments.Arguments) > 0 {
					ad.Code = act.Arguments.Arguments[0].String
				}
			case ActionKindVar: // kind=6: variable assignment
				if len(act.Arguments.Arguments) >= 2 {
					ad.Code = fmt.Sprintf("%s = %s", act.Arguments.Arguments[0].String, act.Arguments.Arguments[1].String)
				}
			case ActionKindNormal: // kind=0: drag-and-drop action (function call)
				ad.FunctionName = act.FunctionName
			}

			for _, arg := range act.Arguments.Arguments {
				ad.Arguments = append(ad.Arguments, ArgData{
					Kind:   arg.Kind,
					Value:  arg.String,
					Object: arg.Object,
				})
			}

			evData.Actions = append(evData.Actions, ad)
		}

		obj.Events = append(obj.Events, evData)
	}

	return obj, nil
}

func cleanGMXString(s string) string {
	s = strings.TrimSpace(s)
	if s == "<undefined>" || s == "&lt;undefined&gt;" {
		return ""
	}
	return s
}

// ParseAllObjects parses all object.gmx files in a directory
func ParseAllObjects(objDir string) ([]*ObjectData, error) {
	entries, err := os.ReadDir(objDir)
	if err != nil {
		return nil, fmt.Errorf("reading object directory: %w", err)
	}

	var objects []*ObjectData
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".object.gmx") {
			continue
		}
		o, err := ParseObject(filepath.Join(objDir, e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping object %s: %v\n", e.Name(), err)
			continue
		}
		objects = append(objects, o)
	}

	return objects, nil
}
