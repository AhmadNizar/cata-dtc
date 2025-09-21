package handler

import (
    "net/http"

    "github.com/AhmadNizar/cata-dtc/internal/dto"
    "github.com/AhmadNizar/cata-dtc/internal/presenter"
    "github.com/AhmadNizar/cata-dtc/internal/usecase/pokemon"
    "github.com/gin-gonic/gin"
)

// ApiHandler handles API integration HTTP requests
type ApiHandler struct {
    pokemonService pokemon.Service
}

// NewApiHandler returns a new ApiHandler
func NewApiHandler(pokemonService pokemon.Service) *ApiHandler {
    return &ApiHandler{pokemonService: pokemonService}
}

func (ah *ApiHandler) Sync(c *gin.Context) {
    err := ah.pokemonService.SyncPokemonData()
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.GeneralResponseDTO{
            OK:      false,
            Message: "failed to sync pokemon data",
        })
        return
    }
    c.JSON(http.StatusOK, dto.GeneralResponseDTO{
        OK:      true,
        Message: "Successfully synced pokemon data",
    })
}

func (ah *ApiHandler) GetItems(c *gin.Context) {
    pokemons, total, err := ah.pokemonService.GetPokemonItems()
    if err != nil {
        c.JSON(http.StatusInternalServerError, dto.GeneralResponseDTO{
            OK:      false,
            Message: "failed to fetch pokemon data",
        })
        return
    }

    // Convert entities to presenter format
    items := make([]presenter.Pokemon, len(pokemons))
    for i, pokemon := range pokemons {
        // Convert types
        types := make([]presenter.PokemonType, len(pokemon.Types))
        for j, pokemonType := range pokemon.Types {
            types[j] = presenter.PokemonType{
                Name: pokemonType.TypeName,
            }
        }

        // Convert abilities
        abilities := make([]presenter.PokemonAbility, len(pokemon.Abilities))
        for j, pokemonAbility := range pokemon.Abilities {
            abilities[j] = presenter.PokemonAbility{
                Name:     pokemonAbility.AbilityName,
                IsHidden: pokemonAbility.IsHidden,
            }
        }

        items[i] = presenter.Pokemon{
            ID:        pokemon.ID,
            Name:      pokemon.Name,
            Height:    pokemon.Height,
            Weight:    pokemon.Weight,
            BaseExp:   pokemon.BaseExp,
            Order:     pokemon.OrderNum,
            Types:     types,
            Abilities: abilities,
            CreatedAt: pokemon.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
            UpdatedAt: pokemon.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
        }
    }

    result := presenter.PokemonList{
        Items: items,
        Total: total,
        Page:  1,
        Limit: len(items),
    }

    c.JSON(http.StatusOK, dto.GeneralResponseDTO{
        OK:      true,
        Message: "Successfully get pokemon data",
        Data:    result,
    })
}