package common

import (
	"crypto/tls"
	"fmt"

	std_ck "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
)

// For testing purposes
var clickhouseOpenFunc = clickhouse.Open

func OpenDatabaseVip() (*gorm.DB, error) {
	// Load database configuration from environment variables
	var (
		prefix  = "database."
		driver  = viper.GetString(prefix + "driver")
		dialect = prefix + driver
	)

	switch driver {
	case "mysql":
		return OpenMySQLVip(viper.Sub(dialect))
	case "postgres":
		return OpenPostgresVip(viper.Sub(dialect))
	case "sqlite":
		return OpenSQLiteVip(viper.Sub(dialect))
	case "clickhouse":
		return OpenClickhouseVip(viper.Sub(dialect))
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", driver)
	}
}

func OpenMySQLVip(vip *viper.Viper) (*gorm.DB, error) {

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

func OpenPostgresVip(vip *viper.Viper) (*gorm.DB, error) {
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

func OpenSQLiteVip(vip *viper.Viper) (*gorm.DB, error) {

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

func OpenClickhouseVip(vip *viper.Viper) (*gorm.DB, error) {
	var (
		user     = vip.GetString("user")
		pass     = vip.GetString("pass")
		host     = vip.GetString("host")
		port     = vip.GetInt("port")
		dbname   = vip.GetString("database")
		timeout  = vip.GetDuration("timeout")
		skipTLS  = vip.GetBool("skiptls")
		debug    = vip.GetBool("debug")
		protocol = vip.GetString("protocol")
		proto    = std_ck.Native
	)

	if protocol == "http" {
		proto = std_ck.HTTP
	}

	var tlsConfig *tls.Config
	if protocol == "https" {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: skipTLS,
		}
	}

	sqlDB := std_ck.OpenDB(&std_ck.Options{
		Protocol: proto,
		Addr:     []string{fmt.Sprintf("%s:%d", host, port)},
		Auth: std_ck.Auth{
			Database: dbname,
			Username: user,
			Password: pass,
		},
		TLS: tlsConfig,
		Settings: std_ck.Settings{
			"max_execution_time": 60,
		},
		DialTimeout: timeout,
		Debug:       debug,
	})
	if sqlDB == nil {
		return nil, fmt.Errorf("failed to connect to ClickHouse")
	}

	return gorm.Open(clickhouse.New(clickhouse.Config{
		Conn: sqlDB, // initialize with existing database conn
	}))
}
