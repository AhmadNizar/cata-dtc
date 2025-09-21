package mysql

import (
	"context"
	"fmt"

	"github.com/AhmadNizar/cata-dtc/internal/entity"
	"github.com/AhmadNizar/cata-dtc/internal/repository"
	"gorm.io/gorm"
)

type pokemonAbilityRepository struct {
	db *gorm.DB
}

func NewPokemonAbilityRepository(db *gorm.DB) repository.PokemonAbilityRepository {
	return &pokemonAbilityRepository{
		db: db,
	}
}

func (r *pokemonAbilityRepository) Create(ctx context.Context, pokemonAbility *entity.PokemonAbility) error {
	if err := r.db.WithContext(ctx).Create(pokemonAbility).Error; err != nil {
		return fmt.Errorf("creating pokemon ability: %w", err)
	}
	return nil
}

func (r *pokemonAbilityRepository) GetByID(ctx context.Context, id uint) (*entity.PokemonAbility, error) {
	var pokemonAbility entity.PokemonAbility
	if err := r.db.WithContext(ctx).First(&pokemonAbility, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("getting pokemon ability by id: %w", err)
	}
	return &pokemonAbility, nil
}

func (r *pokemonAbilityRepository) GetByPokemonID(ctx context.Context, pokemonID uint) ([]*entity.PokemonAbility, error) {
	var pokemonAbilities []*entity.PokemonAbility
	if err := r.db.WithContext(ctx).Where("pokemon_id = ?", pokemonID).Find(&pokemonAbilities).Error; err != nil {
		return nil, fmt.Errorf("getting pokemon abilities by pokemon id: %w", err)
	}
	return pokemonAbilities, nil
}

func (r *pokemonAbilityRepository) List(ctx context.Context, limit, offset int) ([]*entity.PokemonAbility, error) {
	var pokemonAbilities []*entity.PokemonAbility
	query := r.db.WithContext(ctx).Order("id ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&pokemonAbilities).Error; err != nil {
		return nil, fmt.Errorf("listing pokemon abilities: %w", err)
	}

	return pokemonAbilities, nil
}

func (r *pokemonAbilityRepository) Update(ctx context.Context, pokemonAbility *entity.PokemonAbility) error {
	if err := r.db.WithContext(ctx).Save(pokemonAbility).Error; err != nil {
		return fmt.Errorf("updating pokemon ability: %w", err)
	}
	return nil
}

func (r *pokemonAbilityRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.PokemonAbility{}, id).Error; err != nil {
		return fmt.Errorf("deleting pokemon ability: %w", err)
	}
	return nil
}

func (r *pokemonAbilityRepository) DeleteByPokemonID(ctx context.Context, pokemonID uint) error {
	if err := r.db.WithContext(ctx).Where("pokemon_id = ?", pokemonID).Delete(&entity.PokemonAbility{}).Error; err != nil {
		return fmt.Errorf("deleting pokemon abilities by pokemon id: %w", err)
	}
	return nil
}

func (r *pokemonAbilityRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.PokemonAbility{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("counting pokemon abilities: %w", err)
	}
	return count, nil
}