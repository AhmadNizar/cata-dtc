package mysql

import (
	"context"
	"fmt"
	"log"

	"github.com/AhmadNizar/cata-dtc/internal/entity"
	"github.com/AhmadNizar/cata-dtc/internal/repository"
	"gorm.io/gorm"
)

type pokemonRepository struct {
	db *gorm.DB
}

func NewPokemonRepository(db *gorm.DB) repository.PokemonRepository {
	return &pokemonRepository{
		db: db,
	}
}

func (r *pokemonRepository) Create(ctx context.Context, pokemon *entity.Pokemon) error {
	if err := r.db.WithContext(ctx).Create(pokemon).Error; err != nil {
		return fmt.Errorf("creating pokemon: %w", err)
	}
	log.Printf("✅ SQL CREATE SUCCESS: Pokemon ID %d (%s) inserted into database", pokemon.ID, pokemon.Name)
	return nil
}

func (r *pokemonRepository) GetByID(ctx context.Context, id uint) (*entity.Pokemon, error) {
	var pokemon entity.Pokemon
	if err := r.db.WithContext(ctx).First(&pokemon, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("getting pokemon by id: %w", err)
	}
	return &pokemon, nil
}

func (r *pokemonRepository) GetByIDWithRelations(ctx context.Context, id uint) (*entity.Pokemon, error) {
	var pokemon entity.Pokemon
	if err := r.db.WithContext(ctx).Preload("Types").Preload("Abilities").First(&pokemon, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("getting pokemon by id with relations: %w", err)
	}
	return &pokemon, nil
}

func (r *pokemonRepository) GetByName(ctx context.Context, name string) (*entity.Pokemon, error) {
	var pokemon entity.Pokemon
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&pokemon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("getting pokemon by name: %w", err)
	}
	return &pokemon, nil
}

func (r *pokemonRepository) GetByNameWithRelations(ctx context.Context, name string) (*entity.Pokemon, error) {
	var pokemon entity.Pokemon
	if err := r.db.WithContext(ctx).Preload("Types").Preload("Abilities").Where("name = ?", name).First(&pokemon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("getting pokemon by name with relations: %w", err)
	}
	return &pokemon, nil
}

func (r *pokemonRepository) List(ctx context.Context, limit, offset int) ([]*entity.Pokemon, error) {
	var pokemons []*entity.Pokemon
	query := r.db.WithContext(ctx).Order("id ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&pokemons).Error; err != nil {
		return nil, fmt.Errorf("listing pokemons: %w", err)
	}

	return pokemons, nil
}

func (r *pokemonRepository) ListWithRelations(ctx context.Context, limit, offset int) ([]*entity.Pokemon, error) {
	var pokemons []*entity.Pokemon
	query := r.db.WithContext(ctx).Preload("Types").Preload("Abilities").Order("id ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&pokemons).Error; err != nil {
		return nil, fmt.Errorf("listing pokemons with relations: %w", err)
	}

	return pokemons, nil
}

func (r *pokemonRepository) Update(ctx context.Context, pokemon *entity.Pokemon) error {
	if err := r.db.WithContext(ctx).Save(pokemon).Error; err != nil {
		return fmt.Errorf("updating pokemon: %w", err)
	}
	return nil
}

func (r *pokemonRepository) Delete(ctx context.Context, id uint) error {
	if err := r.db.WithContext(ctx).Delete(&entity.Pokemon{}, id).Error; err != nil {
		return fmt.Errorf("deleting pokemon: %w", err)
	}
	return nil
}

func (r *pokemonRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Pokemon{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("counting pokemons: %w", err)
	}
	return count, nil
}

func (r *pokemonRepository) CreateOrUpdate(ctx context.Context, pokemon *entity.Pokemon) error {
	// Use transaction to ensure atomicity and idempotency
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		existing, err := r.getByNameInTx(ctx, tx, pokemon.Name)
		if err != nil {
			return fmt.Errorf("checking existing pokemon: %w", err)
		}

		if existing != nil {
			// Check if data actually changed to avoid unnecessary updates
			if r.isDataUnchanged(existing, pokemon) {
				log.Printf("⚡ SQL SKIP: Pokemon ID %d (%s) unchanged, skipping update", existing.ID, existing.Name)
				return nil
			}

			pokemon.ID = existing.ID
			pokemon.CreatedAt = existing.CreatedAt

			// Delete existing types and abilities first
			if err := tx.Where("pokemon_id = ?", pokemon.ID).Delete(&entity.PokemonType{}).Error; err != nil {
				return fmt.Errorf("deleting existing pokemon types: %w", err)
			}
			if err := tx.Where("pokemon_id = ?", pokemon.ID).Delete(&entity.PokemonAbility{}).Error; err != nil {
				return fmt.Errorf("deleting existing pokemon abilities: %w", err)
			}

			// Update the pokemon with new relationships
			if err := tx.Save(pokemon).Error; err != nil {
				return fmt.Errorf("updating pokemon with relationships: %w", err)
			}

			log.Printf("✅ SQL UPDATE SUCCESS: Pokemon ID %d (%s) updated with %d types and %d abilities",
				pokemon.ID, pokemon.Name, len(pokemon.Types), len(pokemon.Abilities))
			return nil
		}

		// Create new pokemon - use ON CONFLICT for extra safety
		if err := tx.Create(pokemon).Error; err != nil {
			return fmt.Errorf("creating pokemon with relationships: %w", err)
		}

		log.Printf("✅ SQL CREATE SUCCESS: Pokemon ID %d (%s) created with %d types and %d abilities",
			pokemon.ID, pokemon.Name, len(pokemon.Types), len(pokemon.Abilities))
		return nil
	})
}

// Helper function to get pokemon by name within a transaction
func (r *pokemonRepository) getByNameInTx(ctx context.Context, tx *gorm.DB, name string) (*entity.Pokemon, error) {
	var pokemon entity.Pokemon
	if err := tx.Preload("Types").Preload("Abilities").Where("name = ?", name).First(&pokemon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("getting pokemon by name in transaction: %w", err)
	}
	return &pokemon, nil
}

// Helper function to check if Pokemon data has actually changed
func (r *pokemonRepository) isDataUnchanged(existing, new *entity.Pokemon) bool {
	// Compare basic fields
	if existing.Height != new.Height ||
		existing.Weight != new.Weight ||
		existing.BaseExp != new.BaseExp ||
		existing.OrderNum != new.OrderNum {
		return false
	}

	// Compare types
	if len(existing.Types) != len(new.Types) {
		return false
	}
	existingTypes := make(map[string]bool)
	for _, t := range existing.Types {
		existingTypes[t.TypeName] = true
	}
	for _, t := range new.Types {
		if !existingTypes[t.TypeName] {
			return false
		}
	}

	// Compare abilities
	if len(existing.Abilities) != len(new.Abilities) {
		return false
	}
	existingAbilities := make(map[string]bool)
	for _, a := range existing.Abilities {
		key := fmt.Sprintf("%s_%t", a.AbilityName, a.IsHidden)
		existingAbilities[key] = true
	}
	for _, a := range new.Abilities {
		key := fmt.Sprintf("%s_%t", a.AbilityName, a.IsHidden)
		if !existingAbilities[key] {
			return false
		}
	}

	return true
}
