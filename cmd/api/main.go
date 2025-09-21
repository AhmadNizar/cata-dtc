package main

import (
	"log"
	"os"

	api "github.com/AhmadNizar/cata-dtc/cmd/api/http"
	"github.com/AhmadNizar/cata-dtc/internal/config"
	"github.com/subosito/gotenv"
	"github.com/urfave/cli"
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:   "log, l",
		Value:  "info",
		Usage:  "logging level (info, debug, error)",
		EnvVar: "APP_ENV",
	},
	cli.StringFlag{
		Name:   "port, p",
		Value:  "8080",
		Usage:  "app http port",
		EnvVar: "APP_PORT",
	},
	cli.StringFlag{
		Name:   "db-host",
		Value:  "mysql",
		Usage:  "MySQL host address",
		EnvVar: "MYSQL_HOST",
	},
	cli.StringFlag{
		Name:   "db-port",
		Value:  "3306",
		Usage:  "MySQL port",
		EnvVar: "MYSQL_PORT",
	},
	cli.StringFlag{
		Name:   "db-user",
		Value:  "root",
		Usage:  "MySQL database username",
		EnvVar: "MYSQL_USER",
	},
	cli.StringFlag{
		Name:   "db-password",
		Value:  "password",
		Usage:  "MySQL database password",
		EnvVar: "MYSQL_ROOT_PASSWORD",
	},
	cli.StringFlag{
		Name:   "db-name",
		Value:  "pokemon_db",
		Usage:  "MySQL database name",
		EnvVar: "MYSQL_DATABASE",
	},
	cli.StringFlag{
		Name:   "redis-host",
		Value:  "redis",
		Usage:  "Redis host address",
		EnvVar: "REDIS_HOST",
	},
	cli.StringFlag{
		Name:   "redis-port",
		Value:  "6379",
		Usage:  "Redis port",
		EnvVar: "REDIS_PORT",
	},
}

func action(c *cli.Context) {
	cfg := config.LoadConfig()

	log.Printf("Starting Pokemon API service on port %s...", cfg.App.Port)
	api.Start(cfg)
}

func main() {
	gotenv.OverLoad("/workspace/.env")

	app := cli.NewApp()
	app.Name = "pokemon-api"
	app.Usage = "Pokemon API integration with MySQL and Redis"
	app.Version = "1.0.0"
	app.Flags = flags
	app.Action = action

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}
