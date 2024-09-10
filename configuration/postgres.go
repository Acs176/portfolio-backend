package configuration

import (
	"errors"
	"fmt"
)

type Postgres struct {
	Host      string
	Port      uint16
	User      string
	Password  string
	DbName    string
	SSLMode   string
	DbVersion uint
}

func (p Postgres) Validate() error {
	if p.Host == "" {
		return errors.Join(ErrorValidatingPostgresConfig, fmt.Errorf("host is required"))
	}
	if p.Port == 0 {
		return errors.Join(ErrorValidatingPostgresConfig, fmt.Errorf("port is required"))
	}
	if p.User == "" {
		return errors.Join(ErrorValidatingPostgresConfig, fmt.Errorf("user is required"))
	}
	if p.Password == "" {
		return errors.Join(ErrorValidatingPostgresConfig, fmt.Errorf("password is required"))
	}
	if p.DbName == "" {
		return errors.Join(ErrorValidatingPostgresConfig, fmt.Errorf("database name is required"))
	}
	if p.SSLMode != "disable" && p.SSLMode != "require" && p.SSLMode != "verify-ca" && p.SSLMode != "verify-full" && p.SSLMode != "prefer" {
		return errors.Join(ErrorValidatingPostgresConfig, fmt.Errorf("ssl_mode must be one of: disable, require, verify-ca, verify-full, prefer"))
	}
	if p.DbVersion == 0 {
		return errors.Join(ErrorValidatingPostgresConfig, fmt.Errorf("database_version is required"))
	}
	return nil
}
