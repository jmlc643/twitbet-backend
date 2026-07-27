package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDB(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, fmt.Errorf("error al abrir conexión con PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error al obtener sql.DB %w", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("error al hacer ping a PostgreSQL: %w", err)
	}

	log.Println("Conexión exitosa a PostgreSQL")
	return db, nil
}

func AutoMigrate(db *gorm.DB, dst ...interface{}) error {
	log.Println("Ejecutando migraciones de base de datos...")
	if err := db.AutoMigrate(dst...); err != nil {
		return fmt.Errorf("error ejecutando automigraciones: %w", err)
	}
	log.Println("Migraciones completadas correctamente")
	return nil
}