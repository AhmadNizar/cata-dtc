package entity

import (
	"time"
)

type PokemonAbility struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	PokemonID   uint      `json:"pokemon_id" gorm:"not null;index:idx_pokemon_ability_pokemon_id"`
	AbilityName string    `json:"ability_name" gorm:"size:100;not null"`
	IsHidden    bool      `json:"is_hidden" gorm:"default:false"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Foreign key relationship
	Pokemon Pokemon `json:"pokemon" gorm:"foreignKey:PokemonID;constraint:OnDelete:CASCADE"`
}

func (PokemonAbility) TableName() string {
	return "pokemon_ability"
}
