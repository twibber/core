package cfg

import (
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"reflect"
)

// Configuration holds all the configuration settings for the application.
// It uses struct tags to map environment variables to struct fields.
//
// The reason why database is not separated from the Configuration struct
// is because there would need to be extra logic to load the database configuration from the environment variables.
type Configuration struct {
	// Application
	Debug  bool   `env:"DEBUG"`
	Port   string `env:"PORT"`
	Name   string `env:"NAME"`
	Domain string `env:"DOMAIN"`

	// Database
	DBHost     string `env:"DB_HOST"`     // Database host address
	DBPort     string `env:"DB_PORT"`     // Database port
	DBUsername string `env:"DB_USERNAME"` // Database username
	DBPassword string `env:"DB_PASSWORD"` // Database password
	DBDatabase string `env:"DB_DATABASE"` // Database name

	// Only required if DEBUG is false
	MailHost     string `env:"MAIL_HOST"`
	MailPort     string `env:"MAIL_PORT"`
	MailSecure   bool   `env:"MAIL_SECURE"`
	MailUsername string `env:"MAIL_AUTH_USERNAME"`
	MailPassword string `env:"MAIL_AUTH_PASSWORD"`
	MailSender   string `env:"MAIL_SENDER"`
	MailReply    string `env:"MAIL_REPLY"`
}

// Config is the global configuration variable
var Config = &Configuration{}

// LoadConfiguration loads the configuration from environment variables
func LoadConfiguration(config *Configuration) {
	val := reflect.ValueOf(config).Elem()

	for i := 0; i < val.NumField(); i++ {
		typeField := val.Type().Field(i)
		env := typeField.Tag.Get("env")

		// Support for boolean fields
		if typeField.Type.Kind() == reflect.Bool {
			val.Field(i).SetBool(os.Getenv(env) == "true")
		} else {
			val.Field(i).SetString(os.Getenv(env))
		}
	}
}

// init loads the configuration from .env and environment variables
func init() {
	// load the .env file into the environment variables
	err := godotenv.Load()
	if err != nil {
		slog.Warn(".env file not loaded, resorting to environment variables alone.")
	}

	// Use the LoadConfiguration function to load the configuration from environment variables
	LoadConfiguration(Config)

	// Set log/slog to use the debug setting
	if Config.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	} else {
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}
}
