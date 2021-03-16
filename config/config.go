package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {

	// Database Config

	// SSLMode to enable/disable SSL connection
	MysqlSSLMode bool `envconfig:"MYSQL_SSL_MODE" default:"true"`
	// MaxIdleConnection to set max idle connection pooling
	MysqlMaxIdleConnection int `envconfig:"MYSQL_MAX_IDLE_CONNECTION" default:"5"`
	// MaxOpenConnection to set max open connection pooling
	MysqlMaxOpenConnection int `envconfig:"MYSQL_MAX_OPEN_CONNECTION" default:"10"`
	// MaxLifetimeConnectionn to set max lifetime of pooling | minutes unit
	MysqlMaxLifetimeConnection int `envconfig:"MYSQL_MAX_LIFETIME_CONNECTION" default:"10"`
	// Host is host of mysql service
	MysqlHost string `envconfig:"MYSQL_HOST" default:""`
	// Port is port of mysql service
	MysqlPort string `envconfig:"MYSQL_PORT" default:""`
	// Username is name of registered user in mysql service
	MysqlUsername string `envconfig:"MYSQL_USERNAME" default:""`
	// DBName is name of registered database in mysql service
	MysqlDBName string `envconfig:"MYSQL_DB_NAME" default:""`
	// Password is password of used Username in mysql service
	MysqlPassword string `envconfig:"MYSQL_PASSWORD" default:""`
	// LogMode is toggle to enable/disable log query in your service by default false
	MysqlLogMode int `envconfig:"MYSQL_LOG_MODE" default:"1"`
	// SingularTable to activate singular table if you are using eloquent query
	MysqlSingularTable bool `envconfig:"MYSQL_SINGULAR_TABLE" default:"true"`
	// ParseTime to parse to local time
	MysqlParseTime bool `envconfig:"MYSQL_PARSE_TIME" default:"true"`
	// Charset to define charset of database
	MysqlCharset string `envconfig:"MYSQL_CHARSET" default:"utf8mb4"`
	// Charset to define charset of database
	MysqlLoc string `envconfig:"MYSQL_LOC" default:"Local"`

	// Redis name
	RedisName string `envconfig:"REDIS_NAME" default:""`
	// Host is host of redis
	RedisHost string `envconfig:"REDIS_HOST" default:""`
	// Post of redis service
	RedisPort int `envconfig:"REDIS_PORT" default:"6379"`
	// configured password of redis server
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	// MaxIdle is max idle connection of redis pooling
	RedisMaxIdle int `envconfig:"REDIS_MAX_IDLE" default:"10"`
	// MaxActive is max active connection of redis pooling
	RedisMaxActive int `envconfig:"REDIS_MAX_ACTIVE" default:"320"`
	// MaxIdelTimeout is max idle timeout of worker in pooling
	RedisMaxIdleTimeOut int `envconfig:"REDIS_MAX_IDLE_TIMEOUT" default:"5"`
	// MaxConnectionLifetime is max life time connection in pool
	RedisMaxConnLifetime int `envconfig:"REDIS_MAX_CONN_LIFETIME" default:"10"`
	// Wait to disable/enable redis using only connection from pooling
	RedisWait bool `envconfig:"REDIS_WAIT" default:"true"`
}

// Get to get defined configuration
func Get() Config {
	cfg := Config{}

	envconfig.MustProcess("", &cfg)

	return cfg
}
