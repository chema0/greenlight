package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

type Environment = string

const (
	Test        = "test"
	Development = "development"
	Production  = "production"
)

type dbConfig struct {
	Database     string
	Username     string
	Password     string
	Host         string
	Port         int
	SSLMode      string
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

type limiterConfig struct {
	RPS     float64
	Burst   int
	Enabled bool
}

type smtpConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Sender   string
}

type corsConfig struct {
	TrustedOrigins []string
}

type Config struct {
	Port    int
	Env     string
	DB      dbConfig
	Limiter limiterConfig
	SMTP    smtpConfig
	CORS    corsConfig
}

func NewConfig(env Environment) Config {
	cfg := Config{
		Env: env,
	}

	cfg.parse()

	return cfg
}

func (cfg *Config) parse() {
	var filename string

	switch cfg.Env {
	case Test:
		filename = "test.toml"
	case Development:
		filename = "dev.toml"
	case Production:
		filename = "prod.toml"
	}

	filename = filepath.Join("config", filename)

	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AddConfigPath("./../../")

	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	cfg.Port = v.GetInt("api.port")

	maxIdleTime, err := time.ParseDuration(v.GetString("database.max_idle_time"))
	if err != nil {
		panic(fmt.Errorf("invalid max iddle time: %w", err))
	}

	cfg.DB = dbConfig{
		Database:     v.GetString("database.dbname"),
		Username:     v.GetString("database.username"),
		Password:     v.GetString("database.password"),
		Host:         v.GetString("database.host"),
		Port:         v.GetInt("database.port"),
		SSLMode:      v.GetString("database.sslmode"),
		DSN:          v.GetString("database.dsn"),
		MaxOpenConns: v.GetInt("database.max_open_conns"),
		MaxIdleConns: v.GetInt("database.max_idle_conns"),
		MaxIdleTime:  maxIdleTime,
	}

	cfg.Limiter = limiterConfig{
		RPS:     float64(v.GetInt("limiter.rps")),
		Burst:   v.GetInt("limiter.burst"),
		Enabled: v.GetBool("limiter.enabled"),
	}

	cfg.SMTP = smtpConfig{
		Host:     v.GetString("smtp.host"),
		Port:     v.GetInt("smtp.port"),
		Username: v.GetString("smtp.username"),
		Password: v.GetString("smtp.password"),
		Sender:   v.GetString("smtp.sender"),
	}

	cfg.CORS = corsConfig{
		TrustedOrigins: v.GetStringSlice("cors.trusted_origins"),
	}

	// fmt.Fprintf(os.Stdout, "cfg: %+v\n", cfg)
}
