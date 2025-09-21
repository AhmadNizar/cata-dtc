package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AhmadNizar/cata-dtc/internal/entity"
	"github.com/AhmadNizar/cata-dtc/internal/repository"
)

type pokemonAPIRepository struct {
	baseURL    string
	httpClient *http.Client
	maxRetries int
}

type Config struct {
	BaseURL    string
	MaxRetries int
}

func NewPokemonAPIRepository(httpClient *http.Client, config Config) repository.PokemonAPIRepository {
	if config.BaseURL == "" {
		config.BaseURL = "https://pokeapi.co/api/v2"
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}

	return &pokemonAPIRepository{
		baseURL:    config.BaseURL,
		httpClient: httpClient,
		maxRetries: config.MaxRetries,
	}
}

func (r *pokemonAPIRepository) GetPokemon(ctx context.Context, pokemonID int) (*entity.PokemonAPIResponse, error) {
	url := fmt.Sprintf("%s/pokemon/%d", r.baseURL, pokemonID)

	var lastErr error
	for attempt := 0; attempt <= r.maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(attempt) * time.Second):
			}
		}

		pokemon, err := r.fetchPokemon(ctx, url)
		if err == nil {
			return pokemon, nil
		}

		lastErr = err
		if !r.shouldRetry(err) {
			break
		}
	}

	return nil, fmt.Errorf("failed to fetch pokemon after %d retries: %w", r.maxRetries, lastErr)
}

func (r *pokemonAPIRepository) fetchPokemon(ctx context.Context, url string) (*entity.PokemonAPIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("User-Agent", "github.com/AhmadNizar/cata-dtc/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var pokemon entity.PokemonAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &pokemon, nil
}

func (r *pokemonAPIRepository) shouldRetry(err error) bool {
	if err == nil {
		return false
	}
	return true
}