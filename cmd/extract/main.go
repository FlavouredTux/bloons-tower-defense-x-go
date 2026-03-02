// package main - asset extraction and data conversion tool for btdx.
// parses .gmx XML files and copies/organizes assets into the project structure.
// generates JSON manifests for runtime loading.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"btdx/tools/gmx2go/parser"
)

func main() {
	gmxDir := flag.String("gmx", "", "Path to the .gmx project directory")
	outDir := flag.String("out", "", "Output directory for extracted assets and data")
	flag.Parse()

	if *gmxDir == "" || *outDir == "" {
		fmt.Fprintln(os.Stderr, "Usage: extract -gmx <path/to/Bloons TDX.gmx> -out <path/to/btdx-go/assets>")
		os.Exit(1)
	}

	start := time.Now()
	fmt.Println("=== BTDX Asset Extraction Pipeline ===")
	fmt.Printf("Source: %s\n", *gmxDir)
	fmt.Printf("Output: %s\n\n", *outDir)

	// ensure output directories exist
	dirs := []string{
		filepath.Join(*outDir, "sprites"),
		filepath.Join(*outDir, "sounds"),
		filepath.Join(*outDir, "backgrounds"),
		filepath.Join(*outDir, "fonts"),
		filepath.Join(*outDir, "data"),
	}
	for _, d := range dirs {
		os.MkdirAll(d, 0o755)
	}

	var totalErrors int

	// 1. Parse and extract sprites
	fmt.Println("[1/7] Parsing sprites...")
	sprites, err := parser.ParseAllSprites(filepath.Join(*gmxDir, "sprites"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing sprites: %v\n", err)
		totalErrors++
	} else {
		fmt.Printf("  Parsed %d sprites\n", len(sprites))
		// copy sprite images
		srcImgDir := filepath.Join(*gmxDir, "sprites", "images")
		dstImgDir := filepath.Join(*outDir, "sprites")
		copied := 0
		for _, s := range sprites {
			for _, fp := range s.FramePaths {
				src := filepath.Join(srcImgDir, filepath.Base(fp))
				dst := filepath.Join(dstImgDir, filepath.Base(fp))
				if err := copyFile(src, dst); err != nil {
					// try with the path as-is from the gmx
					src2 := filepath.Join(*gmxDir, "sprites", fp)
					if err2 := copyFile(src2, dst); err2 != nil {
						fmt.Fprintf(os.Stderr, "  WARN: missing frame %s\n", fp)
						continue
					}
				}
				copied++
			}
		}
		fmt.Printf("  Copied %d frame images\n", copied)
		writeJSON(filepath.Join(*outDir, "data", "sprites.json"), sprites)
	}

	// 2. Parse and extract sounds
	fmt.Println("[2/7] Parsing sounds...")
	sounds, err := parser.ParseAllSounds(filepath.Join(*gmxDir, "sound"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing sounds: %v\n", err)
		totalErrors++
	} else {
		fmt.Printf("  Parsed %d sounds\n", len(sounds))
		srcAudioDir := filepath.Join(*gmxDir, "sound", "audio")
		dstAudioDir := filepath.Join(*outDir, "sounds")
		copied := 0
		for _, s := range sounds {
			if s.FileName == "" {
				continue
			}
			src := filepath.Join(srcAudioDir, s.FileName)
			dst := filepath.Join(dstAudioDir, s.FileName)
			if err := copyFile(src, dst); err != nil {
				fmt.Fprintf(os.Stderr, "  WARN: missing audio %s\n", s.FileName)
				continue
			}
			copied++
		}
		fmt.Printf("  Copied %d audio files\n", copied)
		writeJSON(filepath.Join(*outDir, "data", "sounds.json"), sounds)
	}

	// 3. Parse and extract backgrounds
	fmt.Println("[3/7] Parsing backgrounds...")
	backgrounds, err := parser.ParseAllBackgrounds(filepath.Join(*gmxDir, "background"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing backgrounds: %v\n", err)
		totalErrors++
	} else {
		fmt.Printf("  Parsed %d backgrounds\n", len(backgrounds))
		srcBGDir := filepath.Join(*gmxDir, "background", "images")
		dstBGDir := filepath.Join(*outDir, "backgrounds")
		copied := 0
		for _, b := range backgrounds {
			if b.ImagePath == "" {
				continue
			}
			imgFile := filepath.Base(b.ImagePath)
			src := filepath.Join(srcBGDir, imgFile)
			dst := filepath.Join(dstBGDir, imgFile)
			if err := copyFile(src, dst); err != nil {
				fmt.Fprintf(os.Stderr, "  WARN: missing bg image %s\n", imgFile)
				continue
			}
			copied++
		}
		fmt.Printf("  Copied %d background images\n", copied)
		writeJSON(filepath.Join(*outDir, "data", "backgrounds.json"), backgrounds)
	}

	// 4. Parse rooms (in project file order so the first room = starting room)
	fmt.Println("[4/7] Parsing rooms...")
	// read room order from project file
	projectFile := findProjectFile(*gmxDir)
	var rooms []*parser.RoomData
	if projectFile != "" {
		roomOrder, orderErr := parser.ParseRoomOrderFromProject(projectFile)
		if orderErr != nil {
			fmt.Fprintf(os.Stderr, "  WARN: could not read room order from project file: %v\n", orderErr)
			fmt.Println("  Falling back to filesystem order...")
			rooms, err = parser.ParseAllRooms(filepath.Join(*gmxDir, "rooms"))
		} else {
			fmt.Printf("  Room order from project: %d entries (first: %s)\n", len(roomOrder), roomOrder[0])
			rooms, err = parser.ParseAllRoomsOrdered(filepath.Join(*gmxDir, "rooms"), roomOrder)
		}
	} else {
		fmt.Println("  WARN: project file not found, using filesystem order")
		rooms, err = parser.ParseAllRooms(filepath.Join(*gmxDir, "rooms"))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing rooms: %v\n", err)
		totalErrors++
	} else {
		fmt.Printf("  Parsed %d rooms\n", len(rooms))
		writeJSON(filepath.Join(*outDir, "data", "rooms.json"), rooms)
	}

	// 5. Parse paths
	fmt.Println("[5/7] Parsing paths...")
	paths, err := parser.ParseAllPaths(filepath.Join(*gmxDir, "paths"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing paths: %v\n", err)
		totalErrors++
	} else {
		fmt.Printf("  Parsed %d paths\n", len(paths))
		writeJSON(filepath.Join(*outDir, "data", "paths.json"), paths)
	}

	// 6. Parse timelines
	fmt.Println("[6/7] Parsing timelines...")
	timelines, err := parser.ParseAllTimelines(filepath.Join(*gmxDir, "timelines"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing timelines: %v\n", err)
		totalErrors++
	} else {
		fmt.Printf("  Parsed %d timelines\n", len(timelines))
		writeJSON(filepath.Join(*outDir, "data", "timelines.json"), timelines)
	}

	// 7. Parse objects (the big one - game logic)
	fmt.Println("[7/7] Parsing objects...")
	objects, err := parser.ParseAllObjects(filepath.Join(*gmxDir, "objects"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing objects: %v\n", err)
		totalErrors++
	} else {
		fmt.Printf("  Parsed %d objects\n", len(objects))
		writeJSON(filepath.Join(*outDir, "data", "objects.json"), objects)

		// build parent hierarchy map
		parentMap := make(map[string]string)
		spriteMap := make(map[string]string)
		for _, o := range objects {
			if o.ParentName != "" {
				parentMap[o.Name] = o.ParentName
			}
			if o.SpriteName != "" {
				spriteMap[o.Name] = o.SpriteName
			}
		}
		writeJSON(filepath.Join(*outDir, "data", "object_parents.json"), parentMap)
		writeJSON(filepath.Join(*outDir, "data", "object_sprites.json"), spriteMap)

		// categorize objects
		categories := categorizeObjects(objects)
		writeJSON(filepath.Join(*outDir, "data", "object_categories.json"), categories)
	}

	// parse scripts
	fmt.Println("\n[+] Parsing scripts...")
	scripts, err := parser.ParseAllScripts(filepath.Join(*gmxDir, "scripts"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing scripts: %v\n", err)
	} else {
		fmt.Printf("  Parsed %d scripts\n", len(scripts))
		writeJSON(filepath.Join(*outDir, "data", "scripts.json"), scripts)
	}

	// parse fonts
	fmt.Println("[+] Parsing fonts...")
	fonts, err := parser.ParseAllFonts(filepath.Join(*gmxDir, "fonts"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing fonts: %v\n", err)
	} else {
		fmt.Printf("  Parsed %d fonts\n", len(fonts))
		// copy font atlas images
		for _, f := range fonts {
			if f.ImageFile != "" {
				src := filepath.Join(*gmxDir, "fonts", f.ImageFile)
				dst := filepath.Join(*outDir, "fonts", f.ImageFile)
				copyFile(src, dst)
			}
		}
		writeJSON(filepath.Join(*outDir, "data", "fonts.json"), fonts)
	}

	// summary
	elapsed := time.Since(start)
	fmt.Printf("\n=== Extraction Complete ===\n")
	fmt.Printf("Time: %v\n", elapsed.Round(time.Millisecond))
	if totalErrors > 0 {
		fmt.Printf("Errors: %d (check warnings above)\n", totalErrors)
	} else {
		fmt.Println("No errors!")
	}

	// print asset summary
	fmt.Printf("\nAsset Summary:\n")
	if sprites != nil {
		totalFrames := 0
		for _, s := range sprites {
			totalFrames += s.FrameCount
		}
		fmt.Printf("  Sprites:      %d (%d total frames)\n", len(sprites), totalFrames)
	}
	if sounds != nil {
		fmt.Printf("  Sounds:       %d\n", len(sounds))
	}
	if backgrounds != nil {
		fmt.Printf("  Backgrounds:  %d\n", len(backgrounds))
	}
	if rooms != nil {
		fmt.Printf("  Rooms:        %d\n", len(rooms))
	}
	if paths != nil {
		fmt.Printf("  Paths:        %d\n", len(paths))
	}
	if timelines != nil {
		fmt.Printf("  Timelines:    %d\n", len(timelines))
	}
	if objects != nil {
		fmt.Printf("  Objects:      %d\n", len(objects))
	}
	if scripts != nil {
		fmt.Printf("  Scripts:      %d\n", len(scripts))
	}
	if fonts != nil {
		fmt.Printf("  Fonts:        %d\n", len(fonts))
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	os.MkdirAll(filepath.Dir(dst), 0o755)

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func writeJSON(path string, data interface{}) {
	os.MkdirAll(filepath.Dir(path), 0o755)
	f, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR writing %s: %v\n", path, err)
		return
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR encoding JSON %s: %v\n", path, err)
	}
}

// findProjectFile locates the .project.gmx file inside the gmx directory.
func findProjectFile(gmxDir string) string {
	entries, err := os.ReadDir(gmxDir)
	if err != nil {
		return ""
	}
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".project.gmx") {
			return filepath.Join(gmxDir, e.Name())
		}
	}
	return ""
}

// objectCategory groups objects by their role in the game
type ObjectCategory struct {
	Bloons      []string `json:"bloons"`
	Towers      []string `json:"towers"`
	Projectiles []string `json:"projectiles"`
	UI          []string `json:"ui"`
	Maps        []string `json:"maps"`
	Effects     []string `json:"effects"`
	Control     []string `json:"control"`
	Other       []string `json:"other"`
}

func categorizeObjects(objects []*parser.ObjectData) *ObjectCategory {
	cat := &ObjectCategory{}

	// build parent chain lookup
	parentOf := make(map[string]string)
	for _, o := range objects {
		if o.ParentName != "" {
			parentOf[o.Name] = o.ParentName
		}
	}

	// check ancestry
	hasAncestor := func(name, ancestor string) bool {
		visited := make(map[string]bool)
		current := name
		for {
			p, ok := parentOf[current]
			if !ok || visited[current] {
				return false
			}
			if p == ancestor {
				return true
			}
			visited[current] = true
			current = p
		}
	}

	for _, o := range objects {
		n := o.Name
		nl := strings.ToLower(n)
		pl := strings.ToLower(o.ParentName)

		switch {
		// bloons
		case strings.Contains(nl, "bloon") && !strings.Contains(nl, "spawn") && !strings.Contains(nl, "path"),
			hasAncestor(n, "Normal"),
			pl == "normal" || pl == "moab_class" || pl == "boss_class":
			cat.Bloons = append(cat.Bloons, n)

		// towers
		case strings.Contains(nl, "tower") && !strings.Contains(nl, "panel"),
			strings.Contains(nl, "_upg") || strings.Contains(nl, "_ups"),
			hasAncestor(n, "Tower"):
			cat.Towers = append(cat.Towers, n)

		// projectiles
		case strings.Contains(nl, "projectile") || strings.Contains(nl, "dart") ||
			strings.Contains(nl, "bullet") || strings.Contains(nl, "bomb") ||
			strings.Contains(nl, "bolt") || strings.Contains(nl, "shot"),
			hasAncestor(n, "Normal_Projectile"):
			cat.Projectiles = append(cat.Projectiles, n)

		// uI elements
		case strings.Contains(nl, "button") || strings.Contains(nl, "butt") ||
			strings.Contains(nl, "panel") || strings.Contains(nl, "menu") ||
			strings.Contains(nl, "bar") || strings.Contains(nl, "hud") ||
			strings.Contains(nl, "achieve") || strings.Contains(nl, "buy") ||
			strings.Contains(nl, "upgrade") || strings.Contains(nl, "select"):
			cat.UI = append(cat.UI, n)

		// map/room related
		case strings.Contains(nl, "water") || strings.Contains(nl, "path") ||
			strings.Contains(nl, "spawn") || strings.Contains(nl, "track"):
			cat.Maps = append(cat.Maps, n)

		// effects
		case strings.Contains(nl, "effect") || strings.Contains(nl, "particle") ||
			strings.Contains(nl, "explosion") || strings.Contains(nl, "aura") ||
			strings.Contains(nl, "flash"):
			cat.Effects = append(cat.Effects, n)

		// control objects
		case strings.Contains(nl, "control") || strings.Contains(nl, "manager") ||
			strings.Contains(nl, "global") || strings.Contains(nl, "save") ||
			strings.Contains(nl, "load"):
			cat.Control = append(cat.Control, n)

		default:
			cat.Other = append(cat.Other, n)
		}
	}

	return cat
}
