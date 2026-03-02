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

// computeSmoothPath generates interpolated points for smooth paths using Catmull-Rom
func computeSmoothPath(points []PathPointDef, precision int, closed bool) []PathPointDef {
	if precision < 1 {
		precision = 4
	}
	if len(points) < 2 {
		return points
	}

	stepsPerSegment := precision * 4 // more steps = smoother
	var result []PathPointDef

	n := len(points)
	for i := 0; i < n-1; i++ {
		// get 4 control points for Catmull-Rom
		var p0, p1, p2, p3 PathPointDef
		p1 = points[i]
		p2 = points[i+1]

		if i > 0 {
			p0 = points[i-1]
		} else if closed {
			p0 = points[n-1]
		} else {
			// reflect p1 across itself
			p0 = PathPointDef{X: 2*p1.X - p2.X, Y: 2*p1.Y - p2.Y, Speed: p1.Speed}
		}

		if i < n-2 {
			p3 = points[i+2]
		} else if closed {
			p3 = points[0]
		} else {
			p3 = PathPointDef{X: 2*p2.X - p1.X, Y: 2*p2.Y - p1.Y, Speed: p2.Speed}
		}

		for step := 0; step < stepsPerSegment; step++ {
			t := float64(step) / float64(stepsPerSegment)
			x := catmullRom(p0.X, p1.X, p2.X, p3.X, t)
			y := catmullRom(p0.Y, p1.Y, p2.Y, p3.Y, t)
			spd := p1.Speed + t*(p2.Speed-p1.Speed) // linear interpolate speed
			result = append(result, PathPointDef{X: x, Y: y, Speed: spd})
		}
	}

	// add the last point
	result = append(result, points[n-1])

	return result
}

// catmullRom evaluates a Catmull-Rom spline at parameter t
func catmullRom(p0, p1, p2, p3, t float64) float64 {
	t2 := t * t
	t3 := t2 * t
	return 0.5 * (
		(2 * p1) +
		(-p0+p2)*t +
		(2*p0-5*p1+4*p2-p3)*t2 +
		(-p0+3*p1-3*p2+p3)*t3)
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
	Name      string         `json:"Name"`
	Kind      int            `json:"Kind"`
	Closed    bool           `json:"Closed"`
	Precision int            `json:"Precision"`
	Points    []PathPtJSON   `json:"Points"`
}

type PathPtJSON struct {
	X     float64 `json:"X"`
	Y     float64 `json:"Y"`
	Speed float64 `json:"Speed"`
}
