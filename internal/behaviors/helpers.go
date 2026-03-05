package behaviors

import "btdx/internal/engine"

// ── general convenience ─────────────────────────────────────────────

// getGlobal retrieves a float64 from the game's global variable map.
func getGlobal(g *engine.Game, key string) float64 {
	v, _ := g.GlobalVars[key].(float64)
	return v
}

// getVar retrieves a float64 from an instance's variable map.
func getVar(inst *engine.Instance, key string) float64 {
	v, _ := inst.Vars[key].(float64)
	return v
}

// ── projectile defaults ─────────────────────────────────────────────

// initProjDefaults sets LP, PP, leadpop and camopop on a projectile
// instance only if they have NOT already been pre-set by the spawning
// tower. This is the "safe" variant used by most projectiles so that
// upgraded towers can override defaults before Create is called.
func initProjDefaults(inst *engine.Instance, lp, pp, lead, camo float64) {
	if _, ok := inst.Vars["LP"]; !ok {
		inst.Vars["LP"] = lp
	}
	if _, ok := inst.Vars["PP"]; !ok {
		inst.Vars["PP"] = pp
	}
	if _, ok := inst.Vars["leadpop"]; !ok {
		inst.Vars["leadpop"] = lead
	}
	if _, ok := inst.Vars["camopop"]; !ok {
		inst.Vars["camopop"] = camo
	}
}

// setProjDefaults unconditionally sets LP, PP, leadpop and camopop.
// Used by projectiles whose stats are never overridden by the spawner.
func setProjDefaults(inst *engine.Instance, lp, pp, lead, camo float64) {
	inst.Vars["LP"] = lp
	inst.Vars["PP"] = pp
	inst.Vars["leadpop"] = lead
	inst.Vars["camopop"] = camo
}
