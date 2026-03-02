package parser

import (
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FontGMX represents a GameMaker font .gmx file
type FontGMX struct {
	XMLName    xml.Name     `xml:"font"`
	Name       string       `xml:"name"`
	Size       int          `xml:"size"`
	Bold       int          `xml:"bold"`
	Italic     int          `xml:"italic"`
	RenderHQ   int          `xml:"renderhq"`
	Charset    int          `xml:"charset"`
	AA         int          `xml:"aa"`
	IncludeTTF int          `xml:"includeTTF"`
	TTFName    string       `xml:"TTFName"`
	Glyphs     FontGlyphs   `xml:"glyphs"`
	Image      string       `xml:"image"`
}

type FontGlyphs struct {
	Glyphs []FontGlyph `xml:"glyph"`
}

type FontGlyph struct {
	Character int `xml:"character,attr"`
	X         int `xml:"x,attr"`
	Y         int `xml:"y,attr"`
	W         int `xml:"w,attr"`
	H         int `xml:"h,attr"`
	Shift     int `xml:"shift,attr"`
	Offset    int `xml:"offset,attr"`
}

type FontData struct {
	Name      string
	FontName  string // actual font family name
	Size      int
	Bold      bool
	Italic    bool
	AA        int // antialiasing level
	ImageFile string
	Glyphs    []GlyphData
}

type GlyphData struct {
	Character rune
	X, Y      int // position in atlas
	W, H      int // size in atlas
	Shift     int // advance width
	Offset    int // x offset when drawing
}

// ParseFont parses a font .gmx file
func ParseFont(path string) (*FontData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading font gmx %s: %w", path, err)
	}

	var gmx FontGMX
	if err := xml.Unmarshal(data, &gmx); err != nil {
		return nil, fmt.Errorf("parsing font gmx %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".font.gmx")

	fd := &FontData{
		Name:      name,
		FontName:  gmx.Name,
		Size:      gmx.Size,
		Bold:      gmx.Bold != 0,
		Italic:    gmx.Italic != 0,
		AA:        gmx.AA,
		ImageFile: gmx.Image,
	}

	for _, g := range gmx.Glyphs.Glyphs {
		fd.Glyphs = append(fd.Glyphs, GlyphData{
			Character: rune(g.Character),
			X:         g.X,
			Y:         g.Y,
			W:         g.W,
			H:         g.H,
			Shift:     g.Shift,
			Offset:    g.Offset,
		})
	}

	return fd, nil
}

// ParseAllFonts parses all font files in a directory
func ParseAllFonts(fontDir string) ([]*FontData, error) {
	entries, err := os.ReadDir(fontDir)
	if err != nil {
		return nil, fmt.Errorf("reading font directory: %w", err)
	}

	var fonts []*FontData
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".font.gmx") {
			continue
		}
		f, err := ParseFont(filepath.Join(fontDir, e.Name()))
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping font %s: %v\n", e.Name(), err)
			continue
		}
		fonts = append(fonts, f)
	}

	return fonts, nil
}

// ProjectGMX represents the top-level .project.gmx file
type ProjectGMX struct {
	XMLName xml.Name `xml:"assets"`
	Sounds  struct {
		Sounds []SoundRef `xml:"sound"`
		Groups []struct {
			Name   string     `xml:"name,attr"`
			Sounds []SoundRef `xml:"sound"`
		} `xml:"sounds"`
	} `xml:"sounds"`
}

type SoundRef struct {
	Path string `xml:",chardata"`
}

// ScriptData holds a raw GML script
type ScriptData struct {
	Name string
	Code string
}

// ParseScript reads a .gml script file
func ParseScript(path string) (*ScriptData, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading script %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), ".gml")

	return &ScriptData{
		Name: name,
		Code: string(data),
	}, nil
}

// ParseAllScripts parses all .gml files in a directory (recursively)
func ParseAllScripts(scriptDir string) ([]*ScriptData, error) {
	var scripts []*ScriptData

	err := filepath.Walk(scriptDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".gml") {
			return nil
		}
		s, err := ParseScript(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WARNING: skipping script %s: %v\n", info.Name(), err)
			return nil
		}
		scripts = append(scripts, s)
		return nil
	})

	return scripts, err
}
