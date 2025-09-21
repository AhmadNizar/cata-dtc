package presenter

type PokemonType struct {
	Name string `json:"name"`
}

type PokemonAbility struct {
	Name     string `json:"name"`
	IsHidden bool   `json:"is_hidden"`
}

type Pokemon struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Height    int              `json:"height"`
	Weight    int              `json:"weight"`
	BaseExp   int              `json:"base_experience"`
	Order     int              `json:"order"`
	Types     []PokemonType    `json:"types"`
	Abilities []PokemonAbility `json:"abilities"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
}

type PokemonList struct {
	Items []Pokemon `json:"items"`
	Total int64     `json:"total"`
	Page  int       `json:"page"`
	Limit int       `json:"limit"`
}