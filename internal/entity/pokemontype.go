package entity

import (
	"time"
)

type PokemonType struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	PokemonID uint      `json:"pokemon_id" gorm:"not null;index:idx_pokemon_type_pokemon_id"`
	TypeName  string    `json:"type_name" gorm:"size:100;not null;index:idx_pokemon_type_type_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Foreign key relationship
	Pokemon Pokemon `json:"pokemon" gorm:"foreignKey:PokemonID;constraint:OnDelete:CASCADE"`
}

func (PokemonType) TableName() string {
	return "pokemon_type"
}
