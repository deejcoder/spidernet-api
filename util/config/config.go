package config

import (
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
	Host     string
	Port     int
	Db       string
	User     string
	Password string
	SSLMode  string
}

type secrets struct {
	CSRFKey   string
	JWTSecret string
	ApiLogin  string
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
	cfg.Database.Db = viper.GetString("database.database")
	cfg.Database.User = viper.GetString("database.user")
	cfg.Database.Password = viper.GetString("database.password")
	cfg.Database.SSLMode = viper.GetString("database.sslmode")

	cfg.Keys.CSRFKey = viper.GetString("secrets.csrf")
	cfg.Keys.JWTSecret = viper.GetString("secrets.jwtsecret")
	cfg.Keys.ApiLogin = viper.GetString("secrets.api_login")
}

// GetConfig loads config data into a Config struct
func GetConfig() *Config {
	cfg := new(Config)
	cfg.init()

	return cfg
}

// InitConfig sets up the config file
func InitConfig(path string) *Config {

	viper.SetConfigName("config")
	viper.AddConfigPath(path)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	cfg := GetConfig()

	// setup logging
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	log.SetFormatter(formatter)
	log.SetLevel(cfg.LogLevel)
	return cfg
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
