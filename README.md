# Database Connection Library

A Go library for connecting to various database backends using configuration from both files and environment variables.

## Supported Databases

- MySQL
- PostgreSQL
- SQLite
- ClickHouse
- Redis

## Usage

### Basic Usage

```go
import (
    "github.com/hysios/mx_common"
    "github.com/spf13/viper"
)

func main() {
    // Load config from file (optional)
    viper.SetConfigFile("config.yaml")
    viper.ReadInConfig()
    
    // Connect to database based on config
    db, err := common.OpenDatabaseVip()
    if err != nil {
        panic(err)
    }
    
    // Use the database...
}
```

### Environment Variables

You can configure the database connections using environment variables. 
The library will automatically load these variables and override any file-based configuration:

#### General Database Configuration

```bash
# Set the database driver
export DATABASE_DRIVER=clickhouse

# ClickHouse specific settings
export DATABASE_CLICKHOUSE_HOST=localhost
export DATABASE_CLICKHOUSE_PORT=9000
export DATABASE_CLICKHOUSE_USER=default
export DATABASE_CLICKHOUSE_PASS=password
export DATABASE_CLICKHOUSE_DATABASE=default
export DATABASE_CLICKHOUSE_TIMEOUT=10s
export DATABASE_CLICKHOUSE_READ=30s
```

#### Direct ClickHouse Connection

```bash
# When connecting directly to ClickHouse
export CLICKHOUSE_HOST=localhost
export CLICKHOUSE_PORT=9000
export CLICKHOUSE_USER=default
export CLICKHOUSE_PASS=password
```

### Config File Format

Example YAML configuration:

```yaml
database:
  driver: clickhouse
  clickhouse:
    host: localhost
    port: 9000
    user: default
    pass: password
    database: default
    timeout: 10s
    read: 30s
  
  # MySQL configuration
  mysql:
    host: localhost
    port: 3306
    user: root
    pass: password
    database: mydb
    charset: utf8mb4
    parseTime: true
    local: Local
    
  # Other database configs...
```

## Advanced Usage

### Clickhouse Connection

```go
import (
    "github.com/hysios/mx_common"
    "github.com/spf13/viper"
)

func main() {
    // Set up viper configuration
    v := viper.New()
    v.Set("host", "localhost")
    v.Set("port", 9000)
    v.Set("user", "default")
    v.Set("pass", "")
    v.Set("database", "default")
    v.Set("timeout", "10s")
    v.Set("read", "30s")
    
    // Connect to ClickHouse
    db, err := common.OpenClickhouseVip(v)
    if err != nil {
        panic(err)
    }
    
    // Use the database...
}
```

## Environment Variable Transformation

Environment variables are automatically transformed from `UPPERCASE_FORMAT` to `lowercase.nested.format` to match the expected viper configuration patterns:

- `DATABASE_CLICKHOUSE_HOST` → `database.clickhouse.host`
- `DATABASE_DRIVER` → `database.driver`
- `CLICKHOUSE_USER` → `user` (when directly using ClickHouse connections) 