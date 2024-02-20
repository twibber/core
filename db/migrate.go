package db

import (
	"github.com/twibber/core/app/models"
	"log/slog"
	"reflect"
)

// MigrateDB migrates models into the database.
func MigrateDB() {
	// spread the models into the AutoMigrate function, so that all models are migrated.
	if err := DB.Migrator().AutoMigrate(models.Models...); err != nil {
		panic(err)
	}

	// collect the names of the models that were migrated.
	modelNames := make([]string, 0)
	for _, n := range models.Models {
		modelNames = append(modelNames, reflect.TypeOf(n).Elem().Name())
	}

	// log the models that were migrated.
	slog.With("models", modelNames).Info("database migrated successfully")
}
