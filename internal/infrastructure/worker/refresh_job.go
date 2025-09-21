package worker

import (
	"log"
	"time"

	"github.com/AhmadNizar/cata-dtc/internal/usecase/pokemon"
	"github.com/cenkalti/backoff/v4"
	"github.com/sony/gobreaker"
)

type RefreshJob struct {
	pokemonService pokemon.Service
	circuitBreaker *gobreaker.CircuitBreaker
	logger         *log.Logger
}

func NewRefreshJob(pokemonService pokemon.Service) *RefreshJob {
	// Circuit breaker settings
	cbSettings := gobreaker.Settings{
		Name:        "pokemon-sync",
		MaxRequests: 3,
		Interval:    60 * time.Second,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			// Trip if 3 failures in a row
			return counts.ConsecutiveFailures >= 3
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit breaker '%s' changed from %s to %s", name, from, to)
		},
	}

	return &RefreshJob{
		pokemonService: pokemonService,
		circuitBreaker: gobreaker.NewCircuitBreaker(cbSettings),
		logger:         log.Default(),
	}
}

func (j *RefreshJob) Execute() {
	startTime := time.Now()
	j.logger.Printf("🔄 [CRON] Starting scheduled Pokemon data refresh at %s", startTime.Format("2006-01-02 15:04:05"))
	j.logger.Printf("🛡️  Circuit breaker state: %v (failures: %d)", j.circuitBreaker.State(), j.circuitBreaker.Counts().ConsecutiveFailures)

	// Execute with circuit breaker protection
	_, err := j.circuitBreaker.Execute(func() (interface{}, error) {
		return nil, j.executeWithRetry()
	})

	duration := time.Since(startTime)
	if err != nil {
		j.logger.Printf("❌ [CRON] Scheduled Pokemon refresh FAILED after %v: %v", duration, err)
		j.logger.Printf("🛡️  Circuit breaker state after failure: %v", j.circuitBreaker.State())
	} else {
		j.logger.Printf("✅ [CRON] Scheduled Pokemon refresh COMPLETED successfully in %v", duration)
		j.logger.Println("📅 Next refresh scheduled in 15 minutes")
	}
}

func (j *RefreshJob) executeWithRetry() error {
	operation := func() error {
		return j.pokemonService.SyncPokemonData()
	}

	// Exponential backoff configuration
	backoffConfig := backoff.NewExponentialBackOff()
	backoffConfig.InitialInterval = 5 * time.Second
	backoffConfig.MaxInterval = 2 * time.Minute
	backoffConfig.MaxElapsedTime = 10 * time.Minute
	backoffConfig.Multiplier = 2
	backoffConfig.RandomizationFactor = 0.1

	// Retry with backoff
	return backoff.Retry(operation, backoffConfig)
}

func (j *RefreshJob) GetCircuitBreakerState() gobreaker.State {
	return j.circuitBreaker.State()
}

func (j *RefreshJob) GetCircuitBreakerCounts() gobreaker.Counts {
	return j.circuitBreaker.Counts()
}