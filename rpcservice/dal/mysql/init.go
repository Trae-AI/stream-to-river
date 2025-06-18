// Copyright (c) 2025 Bytedance Ltd. and/or its affiliates
// SPDX-License-Identifier: MIT

package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// db is a global variable that holds the GORM database connection instance.
var db *gorm.DB

// DataBaseConfig defines the configuration structure for the database.
// It supports both MySQL and SQLite databases.
type DataBaseConfig struct {
	DBType       string // The type of the database, either "mysql" or other values for SQLite.
	MySQLDSN     string // The Data Source Name for MySQL connections.
	SqliteDBPath string // The file path for SQLite database files.
}

// InitDBWithConfig initializes the database connection based on the provided configuration.
// It supports two types of databases: MySQL and SQLite.
//
// Parameters:
//   - config: A pointer to the DataBaseConfig struct containing database configuration information.
//
// Returns:
//   - error: An error object if an unexpected error occurs during the database initialization process.
func InitDBWithConfig(config *DataBaseConfig) error {
	var err error
	if config.DBType == "mysql" {
		// Initialize MySQL database connection
		db, err = gorm.Open(mysql.Open(config.MySQLDSN), &gorm.Config{
			SkipDefaultTransaction: true,                                // Skip the default transaction for each operation
			PrepareStmt:            true,                                // Cache prepared statements
			Logger:                 logger.Default.LogMode(logger.Info), // Set the logging level to Info
		})
	} else {
		// Initialize SQLite database connection
		db, err = gorm.Open(sqlite.Open(config.SqliteDBPath), &gorm.Config{
			SkipDefaultTransaction: true,
			PrepareStmt:            true,
			Logger:                 logger.Default.LogMode(logger.Info),
		})
	}
	return err
}

// InitMockDB initializes a mock database connection using SQLite.
// This function is typically used for testing purposes.
// It creates a new SQLite database named "stream_mock.db".
//
// Returns:
//   - *gorm.DB: A pointer to the initialized GORM database connection.
//   - error: An error object if an unexpected error occurs during the database initialization process.
func InitMockDB() (*gorm.DB, error) {
	mdb, err := gorm.Open(sqlite.Open("stream_mock.db"), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	db = mdb
	return db, err
}

// GetDB retrieves the global GORM database connection instance.
//
// Returns:
//   - *gorm.DB: A pointer to the global GORM database connection.
func GetDB() *gorm.DB {
	return db
}
