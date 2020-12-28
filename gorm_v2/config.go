package gorm

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// Config defines config for gorm dialer
type Config struct {
	Driver               string        `yaml:"driver"`
	DSN                  string        `yaml:"dsn"`
	DialTimeout          time.Duration `yaml:"dial_timeout"`
	ReadTimeout          time.Duration `yaml:"read_timeout"`
	WriteTimeout         time.Duration `yaml:"write_timeout"`
	MaxOpenConns         int           `yaml:"max_open_conns"`
	MaxIdleConns         int           `yaml:"max_idle_conns"`
	MaxLifeConns         int           `yaml:"max_life_conns"`
	TraceIncludeNotFound bool          `yaml:"trace_include_not_found"`
	DebugSQL             bool          `yaml:"debug_sql"`

	// internal
	mycfg *mysql.Config `yaml:"-"`
}

// Name returns name of the gorm dialer for Manager.
func (c *Config) Name() string {
	dsn, err := mysql.ParseDSN(c.DSN)
	if err != nil {
		return c.Driver
	}

	return fmt.Sprintf("%s(%s/%s)", c.Driver, dsn.Addr, dsn.DBName)
}

// FillWithDefaults apply default values for field with invalid value.
func (c *Config) FillWithDefaults() {
	maxCPU := runtime.NumCPU()

	if c.DialTimeout <= 0 || c.DialTimeout > time.Duration(MaxDialTimeout*maxCPU) {
		c.DialTimeout = MaxDialTimeout
	}

	if c.ReadTimeout <= 0 || c.ReadTimeout > time.Duration(MaxReadTimeout*maxCPU) {
		c.ReadTimeout = MaxReadTimeout
	}

	if c.WriteTimeout <= 0 || c.WriteTimeout > time.Duration(MaxWriteTimeout*maxCPU) {
		c.WriteTimeout = MaxWriteTimeout
	}

	if c.MaxOpenConns <= 0 || c.MaxOpenConns > MaxOpenConn*maxCPU {
		c.MaxOpenConns = MaxOpenConn
	}

	if c.MaxIdleConns <= 0 || c.MaxIdleConns > MaxIdleConn*maxCPU {
		c.MaxIdleConns = MaxIdleConn
	}

	if c.MaxLifeConns <= 0 || c.MaxLifeConns > MaxLifecycleConn*maxCPU {
		c.MaxLifeConns = MaxLifecycleConn
	}
}

// NewMycfg returns a *mysql.Config with timeout settings
func (c *Config) NewMycfg() (dsn *mysql.Config, err error) {
	dsn, err = mysql.ParseDSN(c.DSN)
	if err != nil {
		return
	}

	// adjust timeout of DSN
	if dsn.Timeout <= 0 {
		dsn.Timeout = c.DialTimeout * time.Millisecond
	}
	if dsn.ReadTimeout <= 0 {
		dsn.ReadTimeout = c.ReadTimeout * time.Millisecond
	}
	if dsn.WriteTimeout <= 0 {
		dsn.WriteTimeout = c.WriteTimeout * time.Millisecond
	}

	// sync
	c.DSN = dsn.FormatDSN()

	return
}

// NewWithDB creates a new config with the database name given for gorm dialer
func (c *Config) NewWithDB(dbname string) (*Config, error) {
	mycfg, err := mysql.ParseDSN(c.DSN)
	if err != nil {
		return nil, err
	}

	mycfg.DBName = dbname

	copied := *c
	copied.DSN = mycfg.FormatDSN()
	copied.mycfg = mycfg

	return &copied, nil
}

// IsEqualDB returns true if database specified by dsn is equal to dbname given.
func (c *Config) IsEqualDB(dbname string) bool {
	dsn, err := mysql.ParseDSN(c.DSN)
	if err != nil {
		return false
	}

	return strings.Compare(dsn.DBName, dbname) == 0
}

// A ManagerConfig defines a list of gorm dialer config with its name
type ManagerConfig map[string]*Config
