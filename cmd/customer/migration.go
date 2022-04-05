package main

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	logOption "github.com/nbs-go/nlogger/v2/option"
	"path/filepath"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/customer"
	"repo.pegadaian.co.id/ms-pds/srv-customer/internal/pkg/nucleo/nval"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func bootMigration(workDir string, config *customer.Config) error {
	// If not enabled, then skip
	isEnabled := nval.ParseBooleanFallback(config.DatabaseBootMigration, false)
	if !isEnabled {
		return nil
	}

	// Load database url
	dbUri := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		config.DatabaseDriver,
		config.DatabaseUser, config.DatabasePass,
		config.DatabaseHost, config.DatabasePort,
		config.DatabaseName,
	)

	// Source dir
	sourcePath := filepath.Join(workDir, "/migrations/sql")
	sourceUri := "file://" + sourcePath

	// Init database
	log.Info("migration: Connecting to database...")
	m, err := migrate.New(sourceUri, dbUri)
	if err != nil {
		log.Error("migration: Failed to connect database", logOption.Error(err))
		return err
	}

	err = m.Up()
	if err != nil {
		if err.Error() == "no change" {
			log.Infof("migration: No changes")
			return nil
		}
		log.Error("migration: Failed to run up migration scripts", logOption.Error(err))
		return err
	}

	// Get status
	version, dirty, err := m.Version()
	if err != nil {
		log.Error("migration: Failed to get database migration version")
		return err
	}
	log.Infof("migration: Database version = %d, Forced = %v", version, dirty)

	return nil
}
