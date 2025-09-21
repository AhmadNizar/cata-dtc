package repository

import (
	"context"

	"github.com/AhmadNizar/cata-dtc/internal/entity"
)

type PokemonAPIRepository interface {
	GetPokemon(ctx context.Context, id int) (*entity.PokemonAPIResponse, error)
}