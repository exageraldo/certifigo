package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/exageraldo/certifigo"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// Current Git tag or the name of the snapshot
	// (https://goreleaser.com/cookbooks/using-main.version/)
	version = "dev"
)

type config struct {
	port              int
	dsn               string
	configFileFromCLI string
	credentials       certifigo.EnvCredentials
}

func (c *config) checkCertificateConfigFile() error {
	if c.configFileFromCLI != "" {
		cfgFilePath, err := filepath.Abs(c.configFileFromCLI)
		if err != nil {
			return err
		}
		cfg.configFileFromCLI = cfgFilePath
		if _, err := os.Stat(cfg.configFileFromCLI); os.IsNotExist(err) {
			return err
		}
	}
	return nil
}

func (c *config) getCertificateConfigFile(event certifigo.Event) (*certifigo.CertificateConfigFile, error) {
	if err := cfg.checkCertificateConfigFile(); err != nil {
		return nil, err
	}

	var certificateConfigFile certifigo.CertificateConfigFile
	defaultCfgFile, err := certifigo.LoadDefaultCertificateConfigFile(
		map[string]any{"Event": event},
	)
	if err != nil {
		return nil, err
	}
	if c.configFileFromCLI != "" {
		cliCfgFile, err := certifigo.LoadCertificateConfigFile(
			c.configFileFromCLI,
			map[string]any{
				"Event":  event,
				"Config": defaultCfgFile,
			},
		)
		if err != nil {
			return nil, err
		}
		certificateConfigFile = certifigo.Merge(*defaultCfgFile, *cliCfgFile)
	} else {
		certificateConfigFile = *defaultCfgFile
	}

	return &certificateConfigFile, nil
}

var cfg config

type application struct {
	logger          *slog.Logger
	validationModel ValidationRecordModelInterface
}

func main() {
	// don't use ports 0 ~ 1023 as it used by OS
	flag.IntVar(&cfg.port, "port", 4000, "HTTP network address")
	flag.StringVar(&cfg.dsn, "dsn", "web:web_pwd@/snippetbox?parseTime=true", "MySQL data source name")
	flag.StringVar(&cfg.configFileFromCLI, "config", "", "config file")

	// you need to call this *before* you use the addr variable
	// otherwise it will always be the default value ":4000"
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))

	if err := cfg.checkCertificateConfigFile(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	credConfig, err := certifigo.NewEnvCredentials()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	credConfig.JWTToken = "abc"
	cfg.credentials = *credConfig

	db, err := openDB()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	app := &application{
		logger: logger,
		validationModel: &ValidationRecordModel{
			DB: db,
		},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	logger.Info("starting server", slog.Int("port", cfg.port))
	if err := srv.ListenAndServe(); err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "certifigo_database.db")
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
