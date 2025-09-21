package mysql

import (
	"context"
	"fmt"

	"github.com/AhmadNizar/cata-dtc/internal/entity"
	"github.com/AhmadNizar/cata-dtc/internal/repository"
	"gorm.io/gorm"
)

type pokemonTypeRepository struct {
	db *gorm.DB
}

func NewPokemonTypeRepository(db *gorm.DB) repository.PokemonTypeRepository {
	return &pokemonTypeRepository{
		db: db,
	}
}

func (r *pokemonTypeRepository) Create(ctx context.Context, pokemonType *entity.PokemonType) error {
	if err := r.db.WithContext(ctx).Create(pokemonType).Error; err != nil {
		return fmt.Errorf("creating pokemon type: %w", err)
	}
	return nil
}

func (r *pokemonTypeRepository) GetByID(ctx context.Context, id uint) (*entity.PokemonType, error) {
	var pokemonType entity.PokemonType
	if err := r.db.WithContext(ctx).First(&pokemonType, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("getting pokemon type by id: %w", err)
	}
	return &pokemonType, nil
}

func (r *pokemonTypeRepository) GetByPokemonID(ctx context.Context, pokemonID uint) ([]*entity.PokemonType, error) {
	var pokemonTypes []*entity.PokemonType
	if err := r.db.WithContext(ctx).Where("pokemon_id = ?", pokemonID).Find(&pokemonTypes).Error; err != nil {
		return nil, fmt.Errorf("getting pokemon types by pokemon id: %w", err)
	}
	return pokemonTypes, nil
}

func (r *pokemonTypeRepository) List(ctx context.Context, limit, offset int) ([]*entity.PokemonType, error) {
	var pokemonTypes []*entity.PokemonType
	query := r.db.WithContext(ctx).Order("id ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&pokemonTypes).Error; err != nil {
		return nil, fmt.Errorf("listing pokemon types: %w", err)
	}

	return pokemonTypes, nil
}

func (r *pokemonTypeRepository) Update(ctx context.Context, pokemonType *entity.PokemonType) error {
	if err := r.db.WithContext(ctx).Save(pokemonType).Error; err != nil {
		return fmt.Errorf("updating pokemon type: %w", err)
	}
	return nil
}

func (r *pokemonTypeRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.PokemonType{}, id).Error; err != nil {
		return fmt.Errorf("deleting pokemon type: %w", err)
	}
	return nil
}

func (r *pokemonTypeRepository) DeleteByPokemonID(ctx context.Context, pokemonID uint) error {
	if err := r.db.WithContext(ctx).Where("pokemon_id = ?", pokemonID).Delete(&entity.PokemonType{}).Error; err != nil {
		return fmt.Errorf("deleting pokemon types by pokemon id: %w", err)
	}
	return nil
}

func (r *pokemonTypeRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.PokemonType{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("counting pokemon types: %w", err)
	}
	return count, nil
}