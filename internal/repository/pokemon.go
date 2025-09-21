package repository

import (
	"context"

	"github.com/AhmadNizar/cata-dtc/internal/entity"
)

type PokemonRepository interface {
	Create(ctx context.Context, pokemon *entity.Pokemon) error
	GetByID(ctx context.Context, id uint) (*entity.Pokemon, error)
	GetByIDWithRelations(ctx context.Context, id uint) (*entity.Pokemon, error)
	GetByName(ctx context.Context, name string) (*entity.Pokemon, error)
	GetByNameWithRelations(ctx context.Context, name string) (*entity.Pokemon, error)
	List(ctx context.Context, limit, offset int) ([]*entity.Pokemon, error)
	ListWithRelations(ctx context.Context, limit, offset int) ([]*entity.Pokemon, error)
	Update(ctx context.Context, pokemon *entity.Pokemon) error
	Delete(ctx context.Context, id uint) error
	Count(ctx context.Context) (int64, error)
	CreateOrUpdate(ctx context.Context, pokemon *entity.Pokemon) error
}