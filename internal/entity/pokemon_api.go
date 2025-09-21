package entity

type PokemonTypeAPI struct {
	Slot int `json:"slot"`
	Type struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"type"`
}

type PokemonAbilityAPI struct {
	IsHidden bool `json:"is_hidden"`
	Slot     int  `json:"slot"`
	Ability  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"ability"`
}

type PokemonSprites struct {
	FrontDefault string `json:"front_default"`
	BackDefault  string `json:"back_default"`
	FrontShiny   string `json:"front_shiny"`
	BackShiny    string `json:"back_shiny"`
}

type PokemonAPIResponse struct {
	ID              int              `json:"id"`
	Name            string           `json:"name"`
	Height          int              `json:"height"`
	Weight          int              `json:"weight"`
	BaseExperience  int              `json:"base_experience"`
	Order           int              `json:"order"`
	Types           []PokemonTypeAPI    `json:"types"`
	Abilities       []PokemonAbilityAPI `json:"abilities"`
	Sprites         PokemonSprites   `json:"sprites"`
}