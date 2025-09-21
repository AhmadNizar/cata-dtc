package entity

import (
	"time"
)

type Pokemon struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"uniqueIndex:idx_pokemon_name;size:255;not null"`
	Height    int       `json:"height" gorm:"default:0"`
	Weight    int       `json:"weight" gorm:"default:0"`
	BaseExp   int       `json:"base_experience" gorm:"column:base_experience;default:0"`
	OrderNum  int       `json:"order_num" gorm:"column:order_num;default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Types     []PokemonType    `json:"types" gorm:"foreignKey:PokemonID"`
	Abilities []PokemonAbility `json:"abilities" gorm:"foreignKey:PokemonID"`
}

func (Pokemon) TableName() string {
	return "pokemon"
}
