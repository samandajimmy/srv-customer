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

func (c *Config) loadFromEnv() {
	// If driver is unset set driver from env
	if c.Driver == "" {
		c.Driver = os.Getenv("DB_DRIVER")
	} else {
		// Normalize driver
		switch c.Driver {
		case "postgresql", "pg":
			c.Driver = DriverPostgreSQL
		}
	}

	// If host is unset set host from env
	if c.Host == "" {
		c.Host = os.Getenv("DB_HOST")
	}

	// If port is unset set port from env
	if c.Port == "" {
		c.Port = os.Getenv("DB_PORT")
	}

	// If username is unset set username from env
	if c.Username == "" {
		c.Username = os.Getenv("DB_USER")
	}

	// If password is unset set password from env
	if c.Password == "" {
		c.Password = os.Getenv("DB_PASS")
	}

	// If database name is unset set database name from env
	if c.Database == "" {
		c.Database = os.Getenv("DB_NAME")
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
}

func (c *Config) getDSN() (dsn string, err error) {
	switch c.Driver {
	case DriverMySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port,
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
