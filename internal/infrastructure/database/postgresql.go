package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/configuration"
	"github.com/arisdolanan/demo-gofiber-clean-architecture/pkg/utils"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
)

var (
	dbx  *sqlx.DB
	once sync.Once
)

// ConnectPostgresqlx returns a singleton sqlx.DB connection pool
func ConnectPostgresqlx() *sqlx.DB {
	once.Do(func() {
		pgConfig := configuration.GetPostgresConfig()
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
			pgConfig.Host,
			pgConfig.Username,
			pgConfig.Password,
			pgConfig.DBName,
			pgConfig.Port,
		)

		var err error
		dbx, err = sqlx.Connect("postgres", dsn)
		if err != nil {
			log.Fatalf("failed to connect db: %v", err)
		}

		dbx.SetMaxIdleConns(pgConfig.Pool.Idle)
		dbx.SetMaxOpenConns(pgConfig.Pool.Max)
		dbx.SetConnMaxLifetime(time.Duration(pgConfig.Pool.Lifetime) * time.Second)

		if err = dbx.Ping(); err != nil {
			log.Fatalf("failed to ping db: %v", err)
		}

		if pgConfig.IsMigrate {
			utils.LogFiber("Running database migrations")
			if err := RunMigrations(dbx); err != nil {
				utils.LogFiber("Warning: migration failed: %v" + err.Error())
				utils.LogFiber("Continuing without migrations...")
			}
		}
	})

	return dbx
}

// TestConnection tests if database is accessible
func TestConnection() error {
	db := ConnectPostgresqlx()
	return db.Ping()
}

// GetConnectionInfo returns database connection info
func GetConnectionInfo() string {
	pgConfig := configuration.GetPostgresConfig()
	return fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s",
		pgConfig.Host,
		pgConfig.Username,
		pgConfig.DBName,
		pgConfig.Port,
	)
}

// RunMigrations executes database migrations
func RunMigrations(dbx *sqlx.DB) error {
	driver, err := postgres.WithInstance(dbx.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	fileSource, err := (&file.File{}).Open("migrations")
	if err != nil {
		return err
	}

	pgConfig := configuration.GetPostgresConfig()
	m, err := migrate.NewWithInstance("file", fileSource, pgConfig.DBName, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
