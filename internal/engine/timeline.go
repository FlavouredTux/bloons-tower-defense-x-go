package engine

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// timelineAsset represents a loaded wave timeline
type TimelineAsset struct {
	Name    string
	Steps   []TimelineStep
	MaxStep int
}

type TimelineStep struct {
	Step    int
	Actions []TimelineAction
}

type TimelineAction struct {
	Kind         int
	ExeType      int
	FunctionName string
	WhoName      string
	Code         string
	Arguments    []TimelineArg
}

type TimelineArg struct {
	Kind   int
	Value  string
	Object string
}

// timelineManager handles wave timeline loading and execution
type TimelineManager struct {
	timelines map[string]*TimelineAsset
}

func NewTimelineManager() *TimelineManager {
	return &TimelineManager{
		timelines: make(map[string]*TimelineAsset),
	}
}

// loadTimelinesFromJSON loads all timelines from the extracted JSON manifest
func (tm *TimelineManager) LoadTimelinesFromJSON(assetsDir string) error {
	path := filepath.Join(assetsDir, "data", "timelines.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("reading timelines.json: %w", err)
	}

	var tlDefs []TimelineJSON
	if err := json.Unmarshal(data, &tlDefs); err != nil {
		return fmt.Errorf("parsing timelines.json: %w", err)
	}

	for _, td := range tlDefs {
		tl := &TimelineAsset{
			Name: td.Name,
		}

		maxStep := 0
		for _, entry := range td.Entries {
			step := TimelineStep{
				Step: entry.Step,
			}
			if entry.Step > maxStep {
				maxStep = entry.Step
			}

			for _, act := range entry.Actions {
				ta := TimelineAction{
					Kind:         act.Kind,
					ExeType:      act.ExeType,
					FunctionName: act.FunctionName,
					WhoName:      act.WhoName,
					Code:         act.Code,
				}
				for _, arg := range act.Arguments {
					ta.Arguments = append(ta.Arguments, TimelineArg{
						Kind:   arg.Kind,
						Value:  arg.Value,
						Object: arg.Object,
					})
				}
				step.Actions = append(step.Actions, ta)
			}

			tl.Steps = append(tl.Steps, step)
		}

		tl.MaxStep = maxStep
		tm.timelines[tl.Name] = tl
	}

	fmt.Printf("Loaded %d timelines\n", len(tm.timelines))
	return nil
}

// get returns a timeline by name
func (tm *TimelineManager) Get(name string) *TimelineAsset {
	return tm.timelines[name]
}

// getTimelineNames returns all timeline names
func (tm *TimelineManager) GetTimelineNames() []string {
	names := make([]string, 0, len(tm.timelines))
	for name := range tm.timelines {
		names = append(names, name)
	}
	return names
}

// timelineRunner manages the execution of a timeline
type TimelineRunner struct {
	Timeline    *TimelineAsset
	CurrentStep int
	Running     bool
	Speed       float64 // steps per game tick
	Looping     bool

	// callback for executing timeline actions
	OnAction func(action TimelineAction)
}

func NewTimelineRunner(tl *TimelineAsset) *TimelineRunner {
	return &TimelineRunner{
		Timeline: tl,
		Speed:    1.0,
		Running:  true,
	}
}

// tick advances the timeline by one game step
func (tr *TimelineRunner) Tick() {
	if !tr.Running || tr.Timeline == nil {
		return
	}

	// check if any steps should fire at the current step
	for _, step := range tr.Timeline.Steps {
		if step.Step == tr.CurrentStep {
			for _, act := range step.Actions {
				if tr.OnAction != nil {
					tr.OnAction(act)
				}
			}
		}
	}

	tr.CurrentStep++

	// check if we've reached the end
	if tr.CurrentStep > tr.Timeline.MaxStep {
		if tr.Looping {
			tr.CurrentStep = 0
		} else {
			tr.Running = false
		}
	}
}

// reset resets the timeline to the beginning
func (tr *TimelineRunner) Reset() {
	tr.CurrentStep = 0
	tr.Running = true
}

// jSON types for deserialization
type TimelineJSON struct {
	Name    string             `json:"Name"`
	Entries []TimelineEntryJSON `json:"Entries"`
}

type TimelineEntryJSON struct {
	Step    int                  `json:"Step"`
	Actions []TimelineActionJSON `json:"Actions"`
}

type TimelineActionJSON struct {
	Kind         int              `json:"Kind"`
	ExeType      int              `json:"ExeType"`
	FunctionName string           `json:"FunctionName"`
	WhoName      string           `json:"WhoName"`
	Code         string           `json:"Code"`
	Arguments    []TimelineArgJSON `json:"Arguments"`
}

type TimelineArgJSON struct {
	Kind   int    `json:"Kind"`
	Value  string `json:"Value"`
	Object string `json:"Object"`
}
