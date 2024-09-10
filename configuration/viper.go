package configuration

import (
	"strings"

	"github.com/spf13/viper"
)

// postgres
const (
	postgresHost      = "Postgres.host"
	postgresPort      = "Postgres.port"
	postgresUser      = "Postgres.user"
	postgresPassword  = "Postgres.password"
	postgresDbName    = "Postgres.database"
	postgresSSLMode   = "Postgres.ssl_mode"
	postgresDbVersion = "Postgres.database_version"
)

func viperSetup(configPath, configName, configType string) (*viper.Viper, error) {
	root := viper.New()

	root.AddConfigPath(configPath)
	root.SetConfigName(configName)
	root.SetConfigType(configType)

	root.SetEnvPrefix("ACS")
	root.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	root.AutomaticEnv()

	// Default values
	root.SetDefault(postgresHost, "localhost")
	root.SetDefault(postgresPort, 5432)
	root.SetDefault(postgresUser, "postgres")
	root.SetDefault(postgresPassword, "postgres")
	root.SetDefault(postgresDbName, "test_db")
	root.SetDefault(postgresSSLMode, "disable")
	root.SetDefault(postgresDbVersion, 1)

	return root, nil
}

func setValues(root *viper.Viper, c *Configuration) error {
	c.Postgres.Host = root.GetString(postgresHost)
	c.Postgres.Port = root.GetUint16(postgresPort)
	c.Postgres.DbName = root.GetString(postgresDbName)
	c.Postgres.User = root.GetString(postgresUser)
	c.Postgres.Password = root.GetString(postgresPassword)
	c.Postgres.SSLMode = root.GetString(postgresSSLMode)
	c.Postgres.DbVersion = root.GetUint(postgresDbVersion)

	return nil
}
