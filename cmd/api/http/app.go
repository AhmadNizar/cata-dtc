package api

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/AhmadNizar/cata-dtc/cmd/api/http/handler"
    "github.com/AhmadNizar/cata-dtc/cmd/api/http/router"
    "github.com/AhmadNizar/cata-dtc/internal/config"
    "github.com/AhmadNizar/cata-dtc/internal/entity"
    worker "github.com/AhmadNizar/cata-dtc/internal/infrastructure/worker"
    "github.com/AhmadNizar/cata-dtc/internal/infrastructure/cache"
    infrahttp "github.com/AhmadNizar/cata-dtc/internal/infrastructure/http"
    "github.com/AhmadNizar/cata-dtc/internal/infrastructure/db/mysql"
    httprepo "github.com/AhmadNizar/cata-dtc/internal/repository/http"
    mysqlrepo "github.com/AhmadNizar/cata-dtc/internal/repository/mysql"
    redisrepo "github.com/AhmadNizar/cata-dtc/internal/repository/redis"
    "github.com/AhmadNizar/cata-dtc/internal/usecase/pokemon"
)

func Start(cfg *config.Config) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
        cfg.Database.User,
        cfg.Database.Password,
        cfg.Database.Host,
        cfg.Database.Port,
        cfg.Database.Name,
    )

    dbOption := &entity.MysqlDBConnOption{
        URL:                 dsn,
        MaxIdleConn:         "10",
        MaxOpenConn:         "100",
        MaxLifetimeInMinute: "5",
    }

    db := mysql.NewMysqlRepository(dbOption, cfg.App.Env)

    redisClient, err := cache.NewRedisClient(cache.Config{
        Host:     cfg.Redis.Host,
        Port:     cfg.Redis.Port,
        Password: cfg.Redis.Password,
        DB:       cfg.Redis.DB,
    })
    if err != nil {
        log.Fatalf("Failed to connect to Redis: %v", err)
    }

    httpClient := infrahttp.NewHTTPClient(infrahttp.Config{
        BaseURL:    cfg.Pokemon.BaseURL,
        Timeout:    cfg.Pokemon.Timeout,
        MaxRetries: cfg.Pokemon.MaxRetries,
    })

    pokemonRepo := mysqlrepo.NewPokemonRepository(db)
    pokemonAPIRepo := httprepo.NewPokemonAPIRepository(httpClient, httprepo.Config{
        BaseURL:    cfg.Pokemon.BaseURL,
        MaxRetries: cfg.Pokemon.MaxRetries,
    })
    cacheRepo := redisrepo.NewCacheRepository(redisClient, "pokemon_api")
    pokemonUseCase := pokemon.NewUsecase(pokemonRepo, pokemonAPIRepo, cacheRepo, cfg.Pokemon.CacheTTL)
    apiHandler := handler.NewApiHandler(pokemonUseCase)

    // Initialize background scheduler
    log.Println("📋 Initializing background job scheduler...")
    scheduler := worker.NewScheduler()
    refreshJob := worker.NewRefreshJob(pokemonUseCase)

    // Schedule data refresh every 15 minutes
    log.Println("⏰ Setting up Pokemon data refresh job (every 15 minutes)...")
    if err := scheduler.AddJob("pokemon-refresh", "0 */15 * * * *", refreshJob.Execute); err != nil {
        log.Fatalf("❌ Failed to schedule Pokemon refresh job: %v", err)
    }

    // Start the scheduler
    log.Println("🚀 Starting background scheduler...")
    scheduler.Start()

    // Show active jobs
    activeJobs := scheduler.ListJobs()
    log.Printf("✅ Background scheduler started successfully with %d active job(s): %v", len(activeJobs), activeJobs)
    log.Println("📅 Next Pokemon data refresh will occur within 15 minutes")

    // Log circuit breaker status
    log.Printf("🛡️  Circuit breaker initialized - State: %v", refreshJob.GetCircuitBreakerState())

    // Setup graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

    r := router.NewRouter(apiHandler)

    // Start server in a goroutine
    go func() {
        log.Printf("Starting server on %s:%s", cfg.App.Host, cfg.App.Port)
        if err := r.Run("0.0.0.0:" + cfg.App.Port); err != nil {
            log.Fatalf("could not start server: %v", err)
        }
    }()

    // Wait for interrupt signal to gracefully shutdown
    <-stop
    log.Println("🛑 Shutting down gracefully...")

    // Stop the scheduler
    scheduler.Stop()
    log.Println("✅ Background scheduler stopped")
    log.Println("✅ Application shutdown complete")
}