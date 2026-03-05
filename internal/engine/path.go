package engine

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

// pathAsset represents a loaded bloon path
type PathAsset struct {
	Name      string
	Kind      int // 0=straight lines, 1=smooth curve
	Closed    bool
	Precision int
	Points    []PathPointDef
	// precomputed smooth curve points (if Kind==1)
	CurvePoints []PathPointDef
	TotalLength float64
}

type PathPointDef struct {
	X, Y  float64
	Speed float64
}

// pathManager handles loading and querying of bloon paths
type PathManager struct {
	paths map[string]*PathAsset
}

func NewPathManager() *PathManager {
	return &PathManager{
		paths: make(map[string]*PathAsset),
	}
}

// loadPathsFromJSON loads all paths from the extracted JSON manifest
func (pm *PathManager) LoadPathsFromJSON(assetsDir string) error {
	path := filepath.Join(assetsDir, "data", "paths.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading paths.json: %w", err)
	}

	var pathDefs []PathJSON
	if err := json.Unmarshal(data, &pathDefs); err != nil {
		return fmt.Errorf("parsing paths.json: %w", err)
	}

	for _, pd := range pathDefs {
		pa := &PathAsset{
			Name:      pd.Name,
			Kind:      pd.Kind,
			Closed:    pd.Closed,
			Precision: pd.Precision,
		}

		for _, pt := range pd.Points {
			pa.Points = append(pa.Points, PathPointDef{
				X:     pt.X,
				Y:     pt.Y,
				Speed: pt.Speed,
			})
		}

		// use raw path points for traversal. The previous Catmull-Rom smoothing
		// could overshoot tight corners and produce loop/circle artifacts.
		pa.CurvePoints = pa.Points

		// compute total length
		pa.TotalLength = computePathLength(pa.CurvePoints)

		pm.paths[pa.Name] = pa
	}

	fmt.Printf("Loaded %d paths\n", len(pm.paths))
	return nil
}

// get returns a path by name
func (pm *PathManager) Get(name string) *PathAsset {
	return pm.paths[name]
}

// getPositionAtProgress returns the (x,y) position along a path at a given progress [0..1]
func (pm *PathManager) GetPositionAtProgress(name string, progress float64) (float64, float64) {
	pa := pm.paths[name]
	if pa == nil || len(pa.CurvePoints) < 2 {
		return 0, 0
	}

	// clamp progress
	if progress < 0 {
		progress = 0
	}
	if progress > 1 {
		progress = 1
	}

	targetDist := progress * pa.TotalLength
	currentDist := 0.0

	for i := 1; i < len(pa.CurvePoints); i++ {
		p0 := pa.CurvePoints[i-1]
		p1 := pa.CurvePoints[i]
		segLen := dist(p0.X, p0.Y, p1.X, p1.Y)

		if currentDist+segLen >= targetDist {
			// interpolate within this segment
			t := (targetDist - currentDist) / segLen
			if segLen == 0 {
				t = 0
			}
			return p0.X + t*(p1.X-p0.X), p0.Y + t*(p1.Y-p0.Y)
		}
		currentDist += segLen
	}

	// return last point
	last := pa.CurvePoints[len(pa.CurvePoints)-1]
	return last.X, last.Y
}

// getDirectionAtProgress returns the direction (degrees) at a given progress
func (pm *PathManager) GetDirectionAtProgress(name string, progress float64) float64 {
	pa := pm.paths[name]
	if pa == nil || len(pa.CurvePoints) < 2 {
		return 0
	}

	const epsilon = 0.001
	x1, y1 := pm.GetPositionAtProgress(name, progress)
	x2, y2 := pm.GetPositionAtProgress(name, progress+epsilon)

	dx := x2 - x1
	dy := y2 - y1
	if dx == 0 && dy == 0 {
		return 0
	}

	return math.Atan2(-dy, dx) * 180 / math.Pi
}

func computePathLength(points []PathPointDef) float64 {
	length := 0.0
	for i := 1; i < len(points); i++ {
		length += dist(points[i-1].X, points[i-1].Y, points[i].X, points[i].Y)
	}
	return length
}

func dist(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// pathJSON for deserialization
type PathJSON struct {
	Name      string       `json:"Name"`
	Kind      int          `json:"Kind"`
	Closed    bool         `json:"Closed"`
	Precision int          `json:"Precision"`
	Points    []PathPtJSON `json:"Points"`
}

type PathPtJSON struct {
	X     float64 `json:"X"`
	Y     float64 `json:"Y"`
	Speed float64 `json:"Speed"`
}
