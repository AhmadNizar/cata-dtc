package repository

import (
	"context"

	"github.com/AhmadNizar/cata-dtc/internal/entity"
)

type PokemonTypeRepository interface {
	Create(ctx context.Context, pokemonType *entity.PokemonType) error
	GetByID(ctx context.Context, id uint) (*entity.PokemonType, error)
	GetByPokemonID(ctx context.Context, pokemonID uint) ([]*entity.PokemonType, error)
	List(ctx context.Context, limit, offset int) ([]*entity.PokemonType, error)
	Update(ctx context.Context, pokemonType *entity.PokemonType) error
	Delete(ctx context.Context, id uint) error
	DeleteByPokemonID(ctx context.Context, pokemonID uint) error
	Count(ctx context.Context) (int64, error)
}