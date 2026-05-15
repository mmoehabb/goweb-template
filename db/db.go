package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"goweb/constants"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Connection *gorm.DB

func Init() (*gorm.DB, error) {
	var err error
	Connection, err = gorm.Open(postgres.Open(constants.AppConfig.DatabaseUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v\n", err)
	}
	return Connection, nil
}

func GetConnection() (*DB, error) {
	if Connection == nil {
		return nil, fmt.Errorf("database not initialized. Call Init() first")
	}
	return &DB{conn: Connection}, nil
}

type DB struct {
	conn *gorm.DB
}

func (d *DB) Close() {
	sqlDB, err := d.conn.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func (d *DB) Create(value any) *gorm.DB {
	return d.conn.Create(value)
}

func (d *DB) Where(query any, args ...any) *gorm.DB {
	return d.conn.Where(query, args...)
}

func (d *DB) First(dest any, conditions ...any) *gorm.DB {
	return d.conn.First(dest, conditions...)
}

func (d *DB) Delete(value any, conditions ...any) *gorm.DB {
	return d.conn.Delete(value, conditions...)
}

func (d *DB) Updates(value any) *gorm.DB {
	return d.conn.Updates(value)
}

func RunMigrations(models ...any) error {
	if Connection == nil {
		return fmt.Errorf("database not initialized. Call Init() first")
	}
	return Connection.AutoMigrate(models...)
}

func RunGooseMigrations() error {
	if Connection == nil {
		return fmt.Errorf("database not initialized. Call Init() first")
	}

	sqlDB, err := Connection.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working dir: %w", err)
	}

	migrationsDir := filepath.Join(wd, "db", "goose_migrations")
	goose.SetTableName("goose_db_version")
	if err := goose.RunContext(context.Background(), "up", sqlDB, migrationsDir); err != nil {
		return fmt.Errorf("failed to run goose migrations: %w", err)
	}
	return nil
}

func Seed() error {
	return nil
}
