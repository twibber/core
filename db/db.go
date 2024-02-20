package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "gorm.io/gorm/logger"
	"log/slog"

	"github.com/twibber/core/cfg"
)

// DB is the global variable used to use GORM
var DB *gorm.DB

// init the database connection
func init() {
	// Create the connection URL
	connUrl := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		cfg.Config.DBUsername,
		cfg.Config.DBPassword,
		cfg.Config.DBHost,
		cfg.Config.DBPort,
		cfg.Config.DBDatabase,
	)

	// Connect to the database via GORM
	if conn, err := gorm.Open(postgres.Open(connUrl), &gorm.Config{
		FullSaveAssociations: true,
	}); err != nil {
		// if an error occurs panic, which will cause the application to crash before the webserver is started
		panic(err)
	} else {
		// Log the database connection
		slog.With(slog.String("host", cfg.Config.DBHost),
			slog.String("port", cfg.Config.DBPort),
			slog.String("username", cfg.Config.DBUsername),
			slog.String("database", cfg.Config.DBDatabase),
		).Info("initiated database connection")

		DB = conn
	}

	if cfg.Config.Debug {
		DB.Debug()
	}

	// Post connection we can migrate the database.
	MigrateDB()
}
