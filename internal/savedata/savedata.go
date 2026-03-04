// Package savedata handles persistent career progression saves.
// Career data is written to a JSON file in the OS user-config directory
// (~/.config/btdx/career.json on Linux) and loaded once at game startup.
package savedata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"btdx/internal/engine"
)

// careerKeys is the complete set of GlobalVars keys that are career-persistent.
// Everything else (wave, money, life, etc.) is in-game only and not saved.
var careerKeys []string

func init() {
	// core career stats + settings
	careerKeys = append(careerKeys,
		"BP", "monkeymoney", "bsouls", "trophies", "XP", "rank", "criteria",
		"totalachievements",
		"soundandmusic", "mute", "soundmute",
		// agents
		"angrysquirrel", "bloonburybush", "sprinkler", "monkeynurse", "bananamobile",
		// bloon toggles
		"bullyenable", "mmoabenable", "horrorenable", "ufoenable",
		"superenable", "motherenable", "lolenable", "clownenable",
		"flowerenable", "crawlerenable", "destroyerenable",
	)

	// tower path unlock counters
	for _, p := range []string{
		"DML", "DMM", "DMR", "TSL", "TSM", "TSR", "BML", "BMM", "BMR",
		"SnML", "SnMM", "SnMR", "NML", "NMM", "NMR", "BCL", "BCM", "BCR",
		"MSL", "MSM", "MSR", "CTL", "CTM", "CTR",
		"GGL", "GGM", "GGR", "IML", "IMM", "IMR", "MBL", "MBM", "MBR",
		"MEL", "MEM", "MER", "MAL", "MAM", "MAR", "BChL", "BChM", "BChR",
		"MApL", "MApM", "MApR", "MAlL", "MAlM", "MAlR",
		"MVL", "MVM", "MVR", "BTL", "BTM", "BTR", "DGL", "DGM", "DGR",
		"MLL", "MLM", "MLR", "HPL", "HPM", "HPR", "SFL", "SFM", "SFR",
		"PML", "PMM", "PMR", "SuML", "SuMM", "SuMR",
	} {
		careerKeys = append(careerKeys, p)
	}

	// per-track stats (32 tracks)
	for i := 1; i <= 32; i++ {
		careerKeys = append(careerKeys,
			fmt.Sprintf("track%dmilestone", i),
			fmt.Sprintf("track%dhardstone", i),
			fmt.Sprintf("track%dbestwave", i),
			fmt.Sprintf("track%dbesthardwave", i),
			fmt.Sprintf("track%dnightstone", i),
			fmt.Sprintf("x%d", i),
			fmt.Sprintf("xx%d", i),
			fmt.Sprintf("t%d", i),
			fmt.Sprintf("n%d", i),
		)
	}

	// special missions
	for i := 1; i <= 16; i++ {
		careerKeys = append(careerKeys, fmt.Sprintf("specialmission%d", i))
	}

	// bounties
	for i := 1; i <= 12; i++ {
		careerKeys = append(careerKeys, fmt.Sprintf("b%d", i))
	}

	// challenges
	for i := 1; i <= 6; i++ {
		careerKeys = append(careerKeys,
			fmt.Sprintf("c%d", i),
			fmt.Sprintf("c%d", i+100),
		)
	}
}

// savePath returns the path to the career.json save file,
// creating the directory if needed.
func savePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		dir, err = os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot locate config dir: %w", err)
		}
	}
	dir = filepath.Join(dir, "btdx")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("cannot create save dir %s: %w", dir, err)
	}
	return filepath.Join(dir, "career.json"), nil
}

// Save writes all career GlobalVars to disk as JSON.
// Uses write-to-temp-then-rename to avoid corruption on crash.
func Save(g *engine.Game) error {
	data := make(map[string]float64, len(careerKeys))
	for _, k := range careerKeys {
		if v, ok := g.GlobalVars[k]; ok {
			if f, ok := v.(float64); ok {
				data[k] = f
			}
		}
	}

	raw, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("savedata marshal: %w", err)
	}

	path, err := savePath()
	if err != nil {
		return err
	}

	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return fmt.Errorf("savedata write: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return fmt.Errorf("savedata rename: %w", err)
	}
	return nil
}

// Load reads career globals from disk and applies them over the defaults already
// set by CareerControl.Create. Keys missing from the file keep their defaults.
// Returns nil on first run (no file yet).
func Load(g *engine.Game) error {
	path, err := savePath()
	if err != nil {
		return err
	}

	raw, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("savedata read: %w", err)
	}

	var data map[string]float64
	if err := json.Unmarshal(raw, &data); err != nil {
		return fmt.Errorf("savedata unmarshal: %w", err)
	}

	for k, v := range data {
		g.GlobalVars[k] = v
	}
	fmt.Printf("[savedata] loaded %d career keys from %s\n", len(data), path)
	return nil
}
