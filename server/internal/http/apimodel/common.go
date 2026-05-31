package apimodel

type Reward struct {
	OpenPower int           `json:"openPower" example:"120"`
	Sitones   []SitoneGrant `json:"sitones,omitempty"`
	Items     []ItemGrant   `json:"items,omitempty"`
}

type SitoneGrant struct {
	DefinitionID string `json:"definitionId" example:"sitone-engineering"`
	Quantity     int    `json:"quantity" example:"1"`
}

type ItemGrant struct {
	DefinitionID string `json:"definitionId" example:"item-camp-sticker"`
	Quantity     int    `json:"quantity" example:"2"`
}

type SitoneSummary struct {
	ID           string `json:"id" example:"sitone_01HR9Z7E2Z2VJ2QZ4P4Z"`
	DefinitionID string `json:"definitionId" example:"sitone-engineering"`
	Name         string `json:"name" example:"Engineering Sitone"`
	Type         string `json:"type" example:"engineering"`
	Rarity       string `json:"rarity" example:"rare"`
	Style        string `json:"style,omitempty" example:"default"`
}

type ItemSummary struct {
	ID           string `json:"id" example:"pit_01HR9Z7E2Z2VJ2QZ4P4Z"`
	DefinitionID string `json:"definitionId" example:"item-camp-sticker"`
	Name         string `json:"name" example:"Camp Sticker"`
	Quantity     int    `json:"quantity" example:"3"`
}
