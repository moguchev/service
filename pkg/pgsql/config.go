package pgsql

import (
	"database/sql"
	"database/sql/driver"
	"time"

	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/opencensus-integrations/ocsql"
)

type Config struct {
	Connection      string        `yaml:"postgresql"`
	MaxOpenConn     int           `yaml:"max_open_conn"`
	MaxIdleConn     int           `yaml:"max_idle_conn"`
	MaxConnLifetime time.Duration `yaml:"max_conn_lifetime"`
	opts            options
}

type options struct {
	Wrapper func(driver.Connector) driver.Connector
}

func TelemetryWrapper(drv driver.Connector) driver.Connector {
	return ocsql.WrapConnector(drv, ocsql.WithOptions(ocsql.TraceOptions{
		AllowRoot:    true,
		Ping:         true,
		RowsClose:    true,
		RowsAffected: true,
		LastInsertID: true,
		Query:        true,
	}))
}

func WithWrapper(cfg Config, wrapper func(drv driver.Connector) driver.Connector) Config {
	if cfg.opts.Wrapper == nil {
		cfg.opts.Wrapper = wrapper
	} else {
		cfg.opts.Wrapper = func(drv driver.Connector) driver.Connector {
			return wrapper(cfg.opts.Wrapper(drv))
		}
	}
	return cfg
}

var DefaultWrapper = TelemetryWrapper

func (cfg Config) CreateDB() (*sqlx.DB, error) {
	var (
		ctor driver.Connector
	)

	drv := stdlib.GetDefaultDriver().(*stdlib.Driver)

	ctor, err := drv.OpenConnector(cfg.Connection)
	if err != nil {
		return nil, err
	}

	if cfg.opts.Wrapper == nil {
		cfg.opts.Wrapper = DefaultWrapper
	}

	ctor = cfg.opts.Wrapper(ctor)

	db := sql.OpenDB(ctor)

	if cfg.MaxConnLifetime != 0 {
		db.SetConnMaxLifetime(cfg.MaxConnLifetime)
	}

	if cfg.MaxIdleConn != 0 {
		db.SetMaxIdleConns(cfg.MaxIdleConn)
	}

	if cfg.MaxOpenConn != 0 {
		db.SetMaxOpenConns(cfg.MaxOpenConn)
	}

	return sqlx.NewDb(db, "pgx"), nil
}
