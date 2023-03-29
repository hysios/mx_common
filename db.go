package common

import (
	"fmt"

	"github.com/hysios/mx/config"
	"github.com/hysios/mx/logger"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func OpenDatabase(cfg *config.Config) (*gorm.DB, error) {
	var (
		prefix  = "database."
		driver  = cfg.Str(prefix + "driver")
		dialect = prefix + driver + "."
	)

	switch driver {
	case "mysql":
		return openMySQL(cfg, dialect)
	case "postgres":
		return openPostgres(cfg, dialect)
	case "sqlite":
		return openSQLite(cfg, dialect)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}

func openMySQL(cfg *config.Config, dialect string) (*gorm.DB, error) {
	var (
		user      = cfg.Str(dialect + "user")
		pass      = cfg.Str(dialect + "pass")
		host      = cfg.Str(dialect + "host")
		port      = cfg.Int(dialect + "port")
		dbname    = cfg.Str(dialect + "database")
		charset   = cfg.Str(dialect + "charset")
		parseTime = cfg.Bool(dialect + "parseTime")
		local     = cfg.Str(dialect + "local")
		dsn       = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?", user, pass, host, port, dbname)
	)

	if charset != "" {
		dsn += "&charset=" + charset
	}

	if local != "" {
		dsn += "&loc=" + local
	}

	if parseTime {
		dsn += "&parseTime=True"
	}

	logger.Logger.Info("open database", zap.String("dsn", dsn))
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func openPostgres(cfg *config.Config, dialect string) (*gorm.DB, error) {
	var (
		user     = cfg.Str(dialect + "user")
		pass     = cfg.Str(dialect + "pass")
		host     = cfg.Str(dialect + "host")
		port     = cfg.Int(dialect + "port")
		dbname   = cfg.Str(dialect + "database")
		timezone = cfg.Str(dialect + "timezone")
		sslmode  = cfg.Str(dialect + "sslmode")
		dsn      = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s TimeZone=%s", host, port, user, dbname, pass, sslmode, timezone)
	)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func openSQLite(cfg *config.Config, dialect string) (*gorm.DB, error) {
	var (
		file = cfg.Str(dialect + "file")
	)

	return gorm.Open(sqlite.Open(file), &gorm.Config{})
}
