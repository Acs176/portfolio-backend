package testdocker

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

type TestDocker struct {
	pool           *dockertest.Pool
	resource       *dockertest.Resource
	dbMigrationURL string
}

func New(migrationURL string) *TestDocker {
	return &TestDocker{dbMigrationURL: migrationURL}
}

func (t *TestDocker) PurgePool() {
	t.pool.Purge(t.resource)
}

func (t *TestDocker) SetupDocker() *pgxpool.Pool {
	var err error
	t.pool, err = dockertest.NewPool("")
	if err != nil {
		panic(fmt.Errorf("could not construct pool: %s", err))
	}

	// uses pool to try to connect to Docker
	err = t.pool.Client.Ping()
	if err != nil {
		panic(fmt.Errorf("could not connect to Docker: %s", err))
	}

	// pulls an image, creates a container based on it and runs it
	t.resource, err = t.pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "12",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=postgres",
			"POSTGRES_DB=test-db",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		panic(err)
	}
	pgPool, err := t.ConnectAndSetupDb(t.pool, t.resource, t.dbMigrationURL)
	if err != nil {
		panic(err)
	}
	return pgPool

}

func (t *TestDocker) ConnectAndSetupDb(pool *dockertest.Pool, resource *dockertest.Resource, migrationURL string) (*pgxpool.Pool, error) {
	host, port := GetHostAndPort(resource, "5432/tcp")
	dbURL := fmt.Sprintf("postgres://postgres:postgres@%s:%d/test-db?sslmode=disable", host, port)

	var conn *pgx.Conn
	err := pool.Retry(func() (err error) {
		conn, err = pgx.Connect(context.Background(), dbURL)
		return err
	})
	if err != nil {
		panic(err)
	}

	migration, err := migrate.New(migrationURL, dbURL)

	if err != nil {
		panic(err)
	}
	err = migration.Up()
	if err != nil {
		panic(err)
	}
	// OJO: Si estas to atascao, aqui falta una linea de codigo
	defer conn.Close(context.Background())

	conf, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}
	pgpool, err := pgxpool.NewWithConfig(context.Background(), conf)
	if err != nil {
		err := fmt.Errorf("pgx connection error: %w", err)
		return nil, err
	}

	return pgpool, nil
}

func GetHostAndPort(resource *dockertest.Resource, id string) (string, uint16) {
	hostPort := strings.Split(resource.GetHostPort(id), ":")
	port, err := strconv.Atoi(hostPort[1])
	if err != nil {
		panic(err)
	}
	return hostPort[0], uint16(port)
}
