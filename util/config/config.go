package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// Config wraps the config file as a struct
type Config struct {
	Version      string
	LogLevel     logrus.Level
	AppName      string
	AppShortName string
	API          apiConfig
	Database     databaseConfig
	Keys         secrets
}

type apiConfig struct {
	UsingHttps     bool
	Port           int
	AllowedMethods []string
	AllowedHeaders []string
	AllowedOrigins []string
}

type databaseConfig struct {
	Host string
	Port int
}

type secrets struct {
	CSRFKey   string
	JWTSecret string
}

// init serializes YAML into a Config struct
func (cfg *Config) init() {
	cfg.Version = viper.GetString("version")
	cfg.setLogLevel(viper.GetString("log_level"))
	cfg.AppName = viper.GetString("app_name")
	cfg.AppShortName = viper.GetString("app_short_name")
	cfg.API.UsingHttps = viper.GetBool("api.usingHttps")
	cfg.API.Port = viper.GetInt("api.port")
	cfg.API.AllowedMethods = viper.GetStringSlice("api.allowed_methods")
	cfg.API.AllowedHeaders = viper.GetStringSlice("api.allowed_headers")
	cfg.API.AllowedOrigins = viper.GetStringSlice("api.allowed_origins")
	cfg.Database.Host = viper.GetString("database.host")
	cfg.Database.Port = viper.GetInt("database.port")
	cfg.Keys.CSRFKey = viper.GetString("secrets.csrf")
	cfg.Keys.JWTSecret = viper.GetString("secrets.jwtsecret")
}

// GetConfig loads config data into a Config struct
func GetConfig() *Config {
	cfg := new(Config)
	cfg.init()

	return cfg
}

// InitConfig sets up the config file
func InitConfig() (*Config, error) {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("%v", err)
		os.Exit(1)
	}

	cfg := GetConfig()

	// setup logging
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)
	log.SetLevel(cfg.LogLevel)
	return cfg, nil
}

func (cfg *Config) setLogLevel(loglevel string) {
	switch loglevel {
	case "debug":
		cfg.LogLevel = logrus.DebugLevel
	case "info":
		cfg.LogLevel = logrus.InfoLevel
	case "warn":
		cfg.LogLevel = logrus.WarnLevel
	case "fatal":
		cfg.LogLevel = logrus.FatalLevel
	}
}
