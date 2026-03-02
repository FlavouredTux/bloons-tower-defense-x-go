package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TimelineGMX represents a GameMaker .timeline.gmx file
type TimelineGMX struct {
	XMLName xml.Name       `xml:"timeline"`
	Entries []TimelineEntry `xml:"entry"`
}

type TimelineEntry struct {
	Step  int      `xml:"step"`
	Event TLEvent  `xml:"event"`
}

type TLEvent struct {
	Actions []ActionGMX `xml:"action"`
}

// TimelineData is our intermediate representation
type TimelineData struct {
	Name    string
	Entries []TimelineStepData
}

type TimelineStepData struct {
	Step    int
	Actions []ActionData
}

// ParseTimeline parses a .timeline.gmx file
func ParseTimeline(path string) (*TimelineData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading timeline gmx %s: %w", path, err)
	}

	var gmx TimelineGMX
	if err := xml.Unmarshal(data, &gmx); err != nil {
		return nil, fmt.Errorf("parsing timeline gmx %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".timeline.gmx")

	tl := &TimelineData{
		Name: name,
	}

	for _, entry := range gmx.Entries {
		step := TimelineStepData{
			Step: entry.Step,
		}

		for _, act := range entry.Event.Actions {
			ad := ActionData{
				Kind:         act.Kind,
				ExeType:      act.ExeType,
				FunctionName: act.FunctionName,
				WhoName:      act.WhoName,
				IsQuestion:   act.IsQuestion != 0,
				UseRelative:  act.UseRelative != 0,
			}

			// Extract code/function info
			switch act.Kind {
			case ActionKindCode:
				if len(act.Arguments.Arguments) > 0 {
					ad.Code = act.Arguments.Arguments[0].String
				}
			case ActionKindVar:
				if len(act.Arguments.Arguments) >= 2 {
					ad.Code = fmt.Sprintf("%s = %s", act.Arguments.Arguments[0].String, act.Arguments.Arguments[1].String)
				}
			case ActionKindNormal:
				ad.FunctionName = act.FunctionName
			}

			for _, arg := range act.Arguments.Arguments {
				ad.Arguments = append(ad.Arguments, ArgData{
					Kind:   arg.Kind,
					Value:  arg.String,
					Object: arg.Object,
				})
			}

			step.Actions = append(step.Actions, ad)
		}

		tl.Entries = append(tl.Entries, step)
	}

	return tl, nil
}

// ParseAllTimelines parses all timeline.gmx files in a directory
func ParseAllTimelines(tlDir string) ([]*TimelineData, error) {
	entries, err := os.ReadDir(tlDir)
	if err != nil {
		return nil, fmt.Errorf("reading timeline directory: %w", err)
	}

	var timelines []*TimelineData
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".timeline.gmx") {
			continue
		}
		t, err := ParseTimeline(filepath.Join(tlDir, e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping timeline %s: %v\n", e.Name(), err)
			continue
		}
		timelines = append(timelines, t)
	}

	return timelines, nil
}

// SoundGMX represents a GameMaker .sound.gmx file
type SoundGMX struct {
	XMLName           xml.Name `xml:"sound"`
	Kind              int      `xml:"kind"`
	Extension         string   `xml:"extension"`
	OrigName          string   `xml:"origname"`
	Effects           int      `xml:"effects"`
	Volume            struct {
		Volume float64 `xml:"volume"`
	} `xml:"volume"`
	Pan               float64  `xml:"pan"`
	Preload           int      `xml:"preload"`
	Data              string   `xml:"data"`
	Compressed        int      `xml:"compressed"`
	Streamed          int      `xml:"streamed"`
	UncompressOnLoad  int      `xml:"uncompressOnLoad"`
	AudioGroup        int      `xml:"audioGroup"`
	BitRates          struct {
		BitRate int `xml:"bitRate"`
	} `xml:"bitRates"`
	SampleRates struct {
		SampleRate int `xml:"sampleRate"`
	} `xml:"sampleRates"`
	Types struct {
		Type int `xml:"type"`
	} `xml:"types"`
	BitDepths struct {
		BitDepth int `xml:"bitDepth"`
	} `xml:"bitDepths"`
}

type SoundData struct {
	Name             string
	FileName         string // filename in audio/ directory
	Extension        string
	Volume           float64
	Pan              float64
	Kind             int  // 0=normal, 1=background, 2=3d, 3=multimedia
	Preload          bool
	Compressed       bool
	Streamed         bool
	UncompressOnLoad bool
	SampleRate       int
	BitRate          int
	BitDepth         int
}

// ParseSound parses a .sound.gmx file
func ParseSound(path string) (*SoundData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading sound gmx %s: %w", path, err)
	}

	var gmx SoundGMX
	if err := xml.Unmarshal(data, &gmx); err != nil {
		return nil, fmt.Errorf("parsing sound gmx %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".sound.gmx")

	return &SoundData{
		Name:             name,
		FileName:         gmx.Data,
		Extension:        gmx.Extension,
		Volume:           gmx.Volume.Volume,
		Pan:              gmx.Pan,
		Kind:             gmx.Kind,
		Preload:          gmx.Preload != 0,
		Compressed:       gmx.Compressed != 0,
		Streamed:         gmx.Streamed != 0,
		UncompressOnLoad: gmx.UncompressOnLoad != 0,
		SampleRate:       gmx.SampleRates.SampleRate,
		BitRate:          gmx.BitRates.BitRate,
		BitDepth:         gmx.BitDepths.BitDepth,
	}, nil
}

// ParseAllSounds parses all sound.gmx files in a directory
func ParseAllSounds(soundDir string) ([]*SoundData, error) {
	entries, err := os.ReadDir(soundDir)
	if err != nil {
		return nil, fmt.Errorf("reading sound directory: %w", err)
	}

	var sounds []*SoundData
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".sound.gmx") {
			continue
		}
		s, err := ParseSound(filepath.Join(soundDir, e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping sound %s: %v\n", e.Name(), err)
			continue
		}
		sounds = append(sounds, s)
	}

	return sounds, nil
}

// BackgroundGMX represents a GameMaker .background.gmx file
type BackgroundGMX struct {
	XMLName    xml.Name `xml:"background"`
	IsTileset  int      `xml:"istileset"`
	TileWidth  int      `xml:"tilewidth"`
	TileHeight int      `xml:"tileheight"`
	TileXOff   int      `xml:"tilexoff"`
	TileYOff   int      `xml:"tileyoff"`
	TileHSep   int      `xml:"tilehsep"`
	TileVSep   int      `xml:"tilevsep"`
	HTile      int      `xml:"HTile"`
	VTile      int      `xml:"VTile"`
	For3D      int      `xml:"For3D"`
	Width      int      `xml:"width"`
	Height     int      `xml:"height"`
	Data       string   `xml:"data"`
}

type BackgroundData struct {
	Name       string
	ImagePath  string
	IsTileset  bool
	TileWidth  int
	TileHeight int
	TileXOff   int
	TileYOff   int
	TileHSep   int
	TileVSep   int
	Width      int
	Height     int
}

// ParseBackground parses a .background.gmx file
func ParseBackground(path string) (*BackgroundData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading background gmx %s: %w", path, err)
	}

	var gmx BackgroundGMX
	if err := xml.Unmarshal(data, &gmx); err != nil {
		return nil, fmt.Errorf("parsing background gmx %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".background.gmx")

	return &BackgroundData{
		Name:       name,
		ImagePath:  strings.ReplaceAll(gmx.Data, "\\", "/"),
		IsTileset:  gmx.IsTileset != 0,
		TileWidth:  gmx.TileWidth,
		TileHeight: gmx.TileHeight,
		TileXOff:   gmx.TileXOff,
		TileYOff:   gmx.TileYOff,
		TileHSep:   gmx.TileHSep,
		TileVSep:   gmx.TileVSep,
		Width:      gmx.Width,
		Height:     gmx.Height,
	}, nil
}

// ParseAllBackgrounds parses all background.gmx files in a directory
func ParseAllBackgrounds(bgDir string) ([]*BackgroundData, error) {
	entries, err := os.ReadDir(bgDir)
	if err != nil {
		return nil, fmt.Errorf("reading background directory: %w", err)
	}

	var bgs []*BackgroundData
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".background.gmx") {
			continue
		}
		b, err := ParseBackground(filepath.Join(bgDir, e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping background %s: %v\n", e.Name(), err)
			continue
		}
		bgs = append(bgs, b)
	}

	return bgs, nil
}
