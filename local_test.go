package common

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// mockDialector implements gorm.Dialector interface
type mockDialector struct {
	DSN string
}

func (m *mockDialector) Name() string {
	return "clickhouse"
}

func (m *mockDialector) Initialize(*gorm.DB) error {
	return nil
}

func (m *mockDialector) Migrator(db *gorm.DB) gorm.Migrator {
	return nil
}

func (m *mockDialector) DataTypeOf(*schema.Field) string {
	return ""
}

func (m *mockDialector) DefaultValueOf(*schema.Field) clause.Expression {
	return nil
}

func (m *mockDialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	// Do nothing
}

func (m *mockDialector) QuoteTo(writer clause.Writer, str string) {
	// Do nothing
}

func (m *mockDialector) Explain(sql string, vars ...interface{}) string {
	return ""
}

func (m *mockDialector) Connect() gorm.ConnPool {
	return nil
}

func TestOpenClickhouseVip(t *testing.T) {

	// Setup test viper config
	vip := viper.New()
	vip.Set("user", os.Getenv("CLICKHOUSE_USER"))
	vip.Set("pass", os.Getenv("CLICKHOUSE_PASSWORD"))
	vip.Set("host", os.Getenv("CLICKHOUSE_HOST"))
	vip.Set("port", os.Getenv("CLICKHOUSE_PORT"))
	vip.Set("database", os.Getenv("CLICKHOUSE_DATABASE"))
	vip.Set("protocol", os.Getenv("CLICKHOUSE_PROTOCOL"))
	vip.Set("skiptls", os.Getenv("CLICKHOUSE_SKIP_TLS"))
	vip.Set("timeout", 30*time.Second)
	vip.Set("read", 60*time.Second)

	// Call the function under test
	db, err := OpenClickhouseVip(vip)

	// Our mock successfully connects but won't be able to do real operations
	assert.NoError(t, err)
	assert.NotNil(t, db)
}

func TestOpenDatabaseVipClickhouse(t *testing.T) {
	// Save the original Open function and restore after test
	originalOpenFunc := clickhouseOpenFunc
	defer func() { clickhouseOpenFunc = originalOpenFunc }()

	// Clear any existing viper settings
	viper.Reset()

	// Setup viper config for clickhouse
	viper.Set("database.driver", "clickhouse")
	viper.Set("database.clickhouse.user", "testuser")
	viper.Set("database.clickhouse.pass", "testpass")
	viper.Set("database.clickhouse.host", "testhost")
	viper.Set("database.clickhouse.port", 9000)
	viper.Set("database.clickhouse.database", "testdb")
	viper.Set("database.clickhouse.timeout", 10*time.Second)
	viper.Set("database.clickhouse.read", 30*time.Second)

	// Mock the clickhouse.Open function to capture the DSN
	var capturedDSN string
	clickhouseOpenFunc = func(dsn string) gorm.Dialector {
		capturedDSN = dsn
		return &mockDialector{DSN: dsn}
	}

	// Call the function under test
	db, err := OpenDatabaseVip()

	// Verify the results
	expectedDSN := "clickhouse://testuser:testpass@testhost:9000?testdb&dial_timeout=10s&read_timeout=30s"
	assert.NoError(t, err)
	assert.NotNil(t, db)
	assert.Equal(t, expectedDSN, capturedDSN)
}
