package common

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

func OpenDatabaseVip() (*gorm.DB, error) {
	var (
		prefix  = "database."
		driver  = viper.GetString(prefix + "driver")
		dialect = prefix + driver
	)

	switch driver {
	case "mysql":
		return openMySQLVip(viper.Sub(dialect))
	case "postgres":
		return openPostgresVip(viper.Sub(dialect))
	case "sqlite":
		return openSQLiteVip(viper.Sub(dialect))
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}

func openMySQLVip(vip *viper.Viper) (*gorm.DB, error) {
	var (
		user        = vip.GetString("user")
		pass        = vip.GetString("pass")
		host        = vip.GetString("host")
		port        = vip.GetInt("port")
		dbname      = vip.GetString("database")
		charset     = vip.GetString("charset")
		parseTime   = vip.GetBool("parseTime")
		local       = vip.GetString("local")
		loglevel    = vip.GetString("loglevel")
		migrateWarn = vip.GetBool("disableMigrateWarn")
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

func openPostgresVip(vip *viper.Viper) (*gorm.DB, error) {
	var (
		user     = vip.GetString("user")
		pass     = vip.GetString("pass")
		host     = vip.GetString("host")
		port     = vip.GetInt("port")
		dbname   = vip.GetString("database")
		timezone = vip.GetString("timezone")
		sslmode  = vip.GetString("sslmode")
		dsn      = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=%s TimeZone=%s", host, port, user, dbname, pass, sslmode, timezone)
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func openSQLiteVip(vip *viper.Viper) (*gorm.DB, error) {
	var (
		file = vip.GetString("file")
	)

	return gorm.Open(sqlite.Open(file), &gorm.Config{})
}

func OpenRedisVip() *redis.Client {
	return redis.NewClient(&redis.Options{
		Network: viper.GetString("redis.network"),
		Addr:    viper.GetString("redis.addr"),
		DB:      viper.GetInt("redis.db"),
		// Username: cfg.Str("redis.username"),
		Password: viper.GetString("redis.password"),
	})
}
