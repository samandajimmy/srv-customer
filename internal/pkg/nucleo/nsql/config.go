package nsql

import (
	"fmt"
)

type Config struct {
	Driver          string `validate:"required"`
	Host            string `validate:"required"`
	Port            uint16 `validate:"required"`
	Username        string `validate:"required"`
	Password        string `validate:"required"`
	Database        string `validate:"required"`
	MaxIdleConn     *int   `validate:"gte=0"`
	MaxOpenConn     *int   `validate:"gte=0"`
	MaxConnLifetime *int   `validate:"gte=0"`
}

func (c *Config) normalizeValue() {
	// Check for optional values, set values if unset

	// If max idle connection is unset, set to 10
	if c.MaxIdleConn == nil {
		c.MaxIdleConn = newInt(10)
	}
	// If max open connection is unset, set to 10
	if c.MaxOpenConn == nil {
		c.MaxOpenConn = newInt(10)
	}
	// If max idle connection is unset, set to 1 second
	if c.MaxConnLifetime == nil {
		c.MaxConnLifetime = newInt(1)
	}

	// Normalize driver
	switch c.Driver {
	case "postgresql", "pg":
		c.Driver = DriverPostgreSQL
	}
}

func (c *Config) getDSN() (dsn string, err error) {
	switch c.Driver {
	case DriverMySQL:
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", c.Username, c.Password, c.Host, c.Port,
			c.Database)
	case DriverPostgreSQL:
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", c.Host, c.Port,
			c.Username, c.Password, c.Database)
	default:
		err = fmt.Errorf("nsql: unsupported database driver %s", c.Driver)
	}
	return
}
