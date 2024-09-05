package common

import (
	"fmt"
	"reflect"

	"github.com/hysios/mx/config"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
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

func DbDialet(cfg *config.Config) string {
	var (
		prefix  = "database."
		driver  = cfg.Str(prefix + "driver")
		dialect = prefix + driver + "."
	)

	return dialect
}

func DbDialetVip() string {
	var (
		prefix  = "database."
		driver  = viper.GetString(prefix + "driver")
		dialect = prefix + driver + "."
	)

	return dialect
}

func OpenDBScope(scope string, cfg *config.Config) (*gorm.DB, error) {
	var (
		prefix  = scope + ".database."
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
		user        = cfg.Str(dialect + "user")
		pass        = cfg.Str(dialect + "pass")
		host        = cfg.Str(dialect + "host")
		port        = cfg.Int(dialect + "port")
		dbname      = cfg.Str(dialect + "database")
		charset     = cfg.Str(dialect + "charset")
		parseTime   = cfg.Bool(dialect + "parseTime")
		local       = cfg.Str(dialect + "local")
		loglevel    = cfg.Str(dialect + "loglevel")
		migrateWarn = cfg.Bool(dialect + "disableMigrateWarn")
		dsn         = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?", user, pass, host, port, dbname)
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

	var gormcfg = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: migrateWarn,
	}
	if loglevel != "" {
		gormcfg.Logger = glogger.Default.LogMode(log2level(loglevel))
	}

	return gorm.Open(mysql.Open(dsn), gormcfg)
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

var gormloglevel = map[string]glogger.LogLevel{
	"silent": 1,
	"error":  2,
	"warn":   3,
	"info":   4,
}

func log2level(s string) glogger.LogLevel {
	if lv, ok := gormloglevel[s]; ok {
		return lv
	}

	return glogger.Warn
}

func AutoIncrementStart(db *gorm.DB, model interface{}, column string, start uint) (err error) {
	var (
		naming = db.Config.NamingStrategy
		v      = reflect.ValueOf(model)
	)

	v = reflect.Indirect(v)

	switch db.Dialector.Name() {
	case "postgres":
		var (
			table     = naming.TableName(v.Type().Name())
			updateSql = fmt.Sprintf("ALTER SEQUENCE %s_%s_seq RESTART WITH %d", table, column, start)
		)
		if err = db.Debug().Exec(updateSql).Error; err != nil {
			return err
		}
	case "mysql":
		var (
			table     = naming.TableName(v.Type().Name())
			updateSql = fmt.Sprintf("ALTER TABLE %s AUTO_INCREMENT = %d", table, start)
		)
		if err = db.Debug().Exec(updateSql).Error; err != nil {
			return err
		}
	case "sqlite":
	}

	return nil
}
