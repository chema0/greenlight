package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/chema0/greenlight/config"
	"github.com/chema0/greenlight/internal/data"
	"github.com/chema0/greenlight/internal/mailer"
	"github.com/chema0/greenlight/internal/vcs"
	_ "github.com/lib/pq"
)

var (
	version = vcs.Version()
)

type application struct {
	config config.Config
	logger *slog.Logger
	models data.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

func main() {
	var env string

	// flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&env, "env", "development", "Environment (development|staging|production)")

	// flag.StringVar(&cfg.db.dsn, "db-dsn", "", "PostgreSQL DSN")

	// flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	// flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	// flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	// flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	// flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	// flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	// flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	// flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	// flag.StringVar(&cfg.smtp.username, "smtp-username", "405d8fc6459cbc", "SMTP username")
	// flag.StringVar(&cfg.smtp.password, "smtp-password", "394388cf9ccdeb", "SMTP password")
	// flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-reply@greenlight.io>", "SMTP sender")

	// flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(s string) error {
	// 	cfg.cors.trustedOrigins = strings.Fields(s)
	// 	return nil
	// })

	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	cfg := config.NewConfig(env)

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		os.Exit(0)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection pool established")

	expvar.NewString("version").Set(version)

	expvar.Publish("goroutines", expvar.Func((func() any {
		return runtime.NumGoroutine()
	})))

	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)

	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	db.SetConnMaxIdleTime(cfg.DB.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}
