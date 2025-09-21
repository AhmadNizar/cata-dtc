package repository

import (
	"context"

	"github.com/AhmadNizar/cata-dtc/internal/entity"
)

type PokemonAbilityRepository interface {
	Create(ctx context.Context, pokemonAbility *entity.PokemonAbility) error
	GetByID(ctx context.Context, id uint) (*entity.PokemonAbility, error)
	GetByPokemonID(ctx context.Context, pokemonID uint) ([]*entity.PokemonAbility, error)
	List(ctx context.Context, limit, offset int) ([]*entity.PokemonAbility, error)
	Update(ctx context.Context, pokemonAbility *entity.PokemonAbility) error
	Delete(ctx context.Context, id uint) error
	DeleteByPokemonID(ctx context.Context, pokemonID uint) error
	Count(ctx context.Context) (int64, error)
}