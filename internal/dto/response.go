package dto

type GeneralResponseDTO struct {
	OK      bool        `json:"ok"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PokemonTypeDTO struct {
	Name string `json:"name"`
}

type PokemonAbilityDTO struct {
	Name     string `json:"name"`
	IsHidden bool   `json:"is_hidden"`
}

type PokemonResponseDTO struct {
	ID        uint                 `json:"id"`
	Name      string               `json:"name"`
	Height    int                  `json:"height"`
	Weight    int                  `json:"weight"`
	BaseExp   int                  `json:"base_experience"`
	Order     int                  `json:"order"`
	Types     []PokemonTypeDTO     `json:"types"`
	Abilities []PokemonAbilityDTO `json:"abilities"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
}

type PokemonListResponseDTO struct {
	Items []PokemonResponseDTO `json:"items"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}