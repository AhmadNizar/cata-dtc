package pokemon

import (
	"github.com/AhmadNizar/cata-dtc/internal/entity"
)

type Service interface {
	SyncPokemonData() error
	GetPokemonItems() ([]*entity.Pokemon, int64, error)
}