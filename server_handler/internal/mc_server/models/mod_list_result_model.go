package models

// ModListResult represents the list of installed mods
type ModListResult struct {
	Mods  []string `json:"mods"`
	Total int      `json:"total"`
}

// NewModListResult creates a new ModListResult
func NewModListResult(mods []string) *ModListResult {
	return &ModListResult{
		Mods:  mods,
		Total: len(mods),
	}
}
