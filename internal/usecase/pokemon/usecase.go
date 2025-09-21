package pokemon

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AhmadNizar/cata-dtc/internal/entity"
	"github.com/AhmadNizar/cata-dtc/internal/repository"
)

type usecase struct {
	pokemonRepo   repository.PokemonRepository
	pokemonAPIRepo repository.PokemonAPIRepository
	cache         repository.CacheRepository
	cacheTTL      time.Duration
}

func NewUsecase(
	pokemonRepo repository.PokemonRepository,
	pokemonAPIRepo repository.PokemonAPIRepository,
	cache repository.CacheRepository,
	cacheTTL time.Duration,
) Service {
	return &usecase{
		pokemonRepo:   pokemonRepo,
		pokemonAPIRepo: pokemonAPIRepo,
		cache:         cache,
		cacheTTL:      cacheTTL,
	}
}

func (u *usecase) SyncPokemonData() error {
	ctx := context.Background()

	log.Println("Starting Pokemon data sync...")

	successCount := 0
	errorCount := 0

	for i := 1; i <= 20; i++ {
		log.Printf("Fetching Pokemon ID: %d", i)

		pokemonData, err := u.pokemonAPIRepo.GetPokemon(ctx, i)
		if err != nil {
			log.Printf("❌ Error fetching Pokemon ID %d: %v", i, err)
			errorCount++
			continue
		}

		pokemon := convertAPIResponseToPokemon(pokemonData)

		if err := u.pokemonRepo.CreateOrUpdate(ctx, pokemon); err != nil {
			log.Printf("❌ Error saving Pokemon ID %d: %v", i, err)
			errorCount++
			continue
		}

		log.Printf("✅ Successfully synced Pokemon: %s", pokemon.Name)
		successCount++
	}

	if err := u.cache.DeleteByPattern(ctx, "pokemon:*"); err != nil {
		log.Printf("Warning: failed to invalidate cache: %v", err)
	}

	if successCount == 0 {
		log.Printf("❌ Pokemon data sync FAILED: 0 success, %d errors", errorCount)
		return fmt.Errorf("all pokemon sync attempts failed")
	} else if errorCount > 0 {
		log.Printf("⚠️ Pokemon data sync PARTIAL: %d success, %d errors", successCount, errorCount)
	} else {
		log.Printf("✅ Pokemon data sync COMPLETED: %d success, 0 errors", successCount)
	}

	return nil
}

func (u *usecase) GetPokemonItems() ([]*entity.Pokemon, int64, error) {
	ctx := context.Background()
	cacheKey := "pokemon:list"

	var cachedPokemons []*entity.Pokemon
	if err := u.cache.Get(ctx, cacheKey, &cachedPokemons); err == nil {
		log.Println("Returning cached Pokemon list")
		total, err := u.pokemonRepo.Count(ctx)
		if err != nil {
			return nil, 0, fmt.Errorf("counting pokemons: %w", err)
		}
		return cachedPokemons, total, nil
	} else if err.Error() != "cache miss" {
		log.Printf("Warning: cache error: %v", err)
	}

	log.Println("Fetching Pokemon list from database")

	pokemons, err := u.pokemonRepo.ListWithRelations(ctx, 0, 0)
	if err != nil {
		return nil, 0, fmt.Errorf("fetching pokemons: %w", err)
	}

	total, err := u.pokemonRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("counting pokemons: %w", err)
	}

	if err := u.cache.Set(ctx, cacheKey, pokemons, u.cacheTTL); err != nil {
		log.Printf("Warning: failed to cache result: %v", err)
	}

	return pokemons, total, nil
}

func convertAPIResponseToPokemon(apiResponse *entity.PokemonAPIResponse) *entity.Pokemon {
	pokemon := &entity.Pokemon{
		ID:       uint(apiResponse.ID),
		Name:     apiResponse.Name,
		Height:   apiResponse.Height,
		Weight:   apiResponse.Weight,
		BaseExp:  apiResponse.BaseExperience,
		OrderNum: apiResponse.Order,
	}

	// Convert types
	for _, typeAPI := range apiResponse.Types {
		pokemon.Types = append(pokemon.Types, entity.PokemonType{
			PokemonID: pokemon.ID,
			TypeName:  typeAPI.Type.Name,
		})
	}

	// Convert abilities
	for _, abilityAPI := range apiResponse.Abilities {
		pokemon.Abilities = append(pokemon.Abilities, entity.PokemonAbility{
			PokemonID:   pokemon.ID,
			AbilityName: abilityAPI.Ability.Name,
			IsHidden:    abilityAPI.IsHidden,
		})
	}

	return pokemon
}
