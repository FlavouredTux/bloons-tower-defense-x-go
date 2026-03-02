package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SpriteGMX represents a GameMaker .sprite.gmx file
type SpriteGMX struct {
	XMLName       xml.Name `xml:"sprite"`
	Type          int      `xml:"type"`
	XOrigin       int      `xml:"xorig"`
	YOrigin       int      `xml:"yorigin"`
	ColKind       int      `xml:"colkind"`
	ColTolerance  int      `xml:"coltolerance"`
	SepMasks      int      `xml:"sepmasks"`
	BBoxMode      int      `xml:"bboxmode"`
	BBoxLeft      int      `xml:"bbox_left"`
	BBoxRight     int      `xml:"bbox_right"`
	BBoxTop       int      `xml:"bbox_top"`
	BBoxBottom    int      `xml:"bbox_bottom"`
	HTile         int      `xml:"HTile"`
	VTile         int      `xml:"VTile"`
	For3D         int      `xml:"For3D"`
	Width         int      `xml:"width"`
	Height        int      `xml:"height"`
	Frames        []Frame  `xml:"frames>frame"`
	TextureGroups struct {
		TextureGroup0 int `xml:"TextureGroup0"`
	} `xml:"TextureGroups"`
}

type Frame struct {
	Index int    `xml:"index,attr"`
	Path  string `xml:",chardata"`
}

// SpriteData is our intermediate representation for the Go port
type SpriteData struct {
	Name       string
	Width      int
	Height     int
	XOrigin    int
	YOrigin    int
	BBox       BBox
	FrameCount int
	FramePaths []string // relative paths to PNG files
	ColKind    int      // collision kind: 0=precise, 1=rect, 2=disk, 3=diamond
}

type BBox struct {
	Left   int
	Right  int
	Top    int
	Bottom int
}

// ParseSprite parses a .sprite.gmx file
func ParseSprite(path string) (*SpriteData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading sprite gmx %s: %w", path, err)
	}

	var gmx SpriteGMX
	if err := xml.Unmarshal(data, &gmx); err != nil {
		return nil, fmt.Errorf("parsing sprite gmx %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".sprite.gmx")

	framePaths := make([]string, len(gmx.Frames))
	for i, f := range gmx.Frames {
		// GM uses backslash paths, normalize to forward slash
		framePaths[i] = strings.ReplaceAll(strings.TrimSpace(f.Path), "\\", "/")
	}

	return &SpriteData{
		Name:       name,
		Width:      gmx.Width,
		Height:     gmx.Height,
		XOrigin:    gmx.XOrigin,
		YOrigin:    gmx.YOrigin,
		BBox: BBox{
			Left:   gmx.BBoxLeft,
			Right:  gmx.BBoxRight,
			Top:    gmx.BBoxTop,
			Bottom: gmx.BBoxBottom,
		},
		FrameCount: len(gmx.Frames),
		FramePaths: framePaths,
		ColKind:    gmx.ColKind,
	}, nil
}

// ParseAllSprites parses all sprite.gmx files in a directory
func ParseAllSprites(spriteDir string) ([]*SpriteData, error) {
	entries, err := os.ReadDir(spriteDir)
	if err != nil {
		return nil, fmt.Errorf("reading sprite directory: %w", err)
	}

	var sprites []*SpriteData
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".sprite.gmx") {
			continue
		}
		s, err := ParseSprite(filepath.Join(spriteDir, e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping sprite %s: %v\n", e.Name(), err)
			continue
		}
		sprites = append(sprites, s)
	}

	return sprites, nil
}
