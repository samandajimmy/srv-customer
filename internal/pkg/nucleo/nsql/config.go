package nsql

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"os"
)

type Config struct {
	Driver          string
	Host            string
	Port            string
	Username        string
	Password        string
	Database        string
	MaxIdleConn     *int
	MaxOpenConn     *int
	MaxConnLifetime *int
}

func LoadFromEnv(c Config, prefixKey string) Config {

	// Default load key
	defaultLoadKey := map[string]string{
		"DB_DRIVER": "DB_DRIVER",
		"DB_HOST":   "DB_HOST",
		"DB_PORT":   "DB_PORT",
		"DB_USER":   "DB_USER",
		"DB_PASS":   "DB_PASS",
		"DB_NAME":   "DB_NAME",
	}

	if prefixKey != "" {
		for i, val := range defaultLoadKey {
			defaultLoadKey[i] = fmt.Sprintf("%s_%s", prefixKey, val)
		}
	}

	// If driver is unset set driver from env
	if c.Driver == "" {
		c.Driver = os.Getenv(defaultLoadKey["DB_DRIVER"])
	}

	// Normalize driver
	switch c.Driver {
	case "postgresql", "pg":
		c.Driver = DriverPostgreSQL
	case "mysql", "mariadb":
		c.Driver = DriverMySQL
	}

	// If host is unset set host from env
	if c.Host == "" {
		c.Host = os.Getenv(defaultLoadKey["DB_HOST"])
	}

	// If port is unset set port from env
	if c.Port == "" {
		c.Port = os.Getenv(defaultLoadKey["DB_PORT"])
	}

	// If username is unset set username from env
	if c.Username == "" {
		c.Username = os.Getenv(defaultLoadKey["DB_USER"])
	}

	// If password is unset set password from env
	if c.Password == "" {
		c.Password = os.Getenv(defaultLoadKey["DB_PASS"])
	}

	// If database name is unset set database name from env
	if c.Database == "" {
		c.Database = os.Getenv(defaultLoadKey["DB_NAME"])
	}

	// If max idle connection is unset, set to 10
	if c.MaxIdleConn == nil {
		c.MaxIdleConn = NewInt(10)
	}
	// If max open connection is unset, set to 10
	if c.MaxOpenConn == nil {
		c.MaxOpenConn = NewInt(10)
	}
	// If max idle connection is unset, set to 1 second
	if c.MaxConnLifetime == nil {
		c.MaxConnLifetime = NewInt(1)
	}
	return c
}

func (c *Config) getDSN() (dsn string, err error) {
	switch c.Driver {
	case DriverMySQL:
		dsn = fmt.Sprintf(`%v:%v@tcp(%v:%v)/%v?parseTime=true`, c.Username, c.Password, c.Host, c.Port,
			c.Database)
	case DriverPostgreSQL:
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port,
			c.Username, c.Password, c.Database)
	default:
		err = fmt.Errorf("nsql: unsupported database driver %s", c.Driver)
	}
	return
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Driver, validation.Required),
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.Database, validation.Required),
		validation.Field(&c.MaxIdleConn, validation.Min(0)),
		validation.Field(&c.MaxOpenConn, validation.Min(0)),
		validation.Field(&c.MaxConnLifetime, validation.Min(0)),
	)
}
