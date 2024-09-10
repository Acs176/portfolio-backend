package main

import (
	"context"
	"database/sql"
	"fmt"
	"localbe/configuration"
	"time"

	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	config configuration.Configuration
)

func init() {
	config = configuration.Get()
	if err := config.Validate(); err != nil {
		panic(fmt.Errorf("configuration failed; err = %w", err))
	}
}

func main() {
	ctx := context.Background()

	// create pg pool
	conf, err := pgxpool.ParseConfig("postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		fmt.Println(err)
	}
	conf.ConnConfig.User = "postgres"
	conf.ConnConfig.Password = "postgres"
	conf.ConnConfig.Database = "test_db"

	pgpool, err := pgxpool.NewWithConfig(ctx, conf)
	if err != nil {
		panic(fmt.Errorf("pgx connection error; err = %w", err))
	}

	// validate connection
	connectionStablished := false
	for i := 0; i < 5; i++ {
		if err := pgpool.Ping(ctx); err != nil {
			fmt.Printf("could not ping database %s; error=%s\n", pgpool.Config().ConnConfig.Database, err.Error())
			time.Sleep(5 * time.Second)
			continue
		}
		connectionStablished = true
		break
	}
	if !connectionStablished {
		panic(fmt.Errorf("could not ping pool with url: ; got error %w", err))
	}

	//db migration

	db, err := sql.Open("pgx", "postgres://postgres:postgres@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		panic(fmt.Errorf("failed to open database; error %v", err))
	}
	defer db.Close()

	driver, err := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to create driver; error %v", err))
	}

	sourceURL := "file://" + "database/migration" // migration path
	m, err := migrate.NewWithDatabaseInstance(sourceURL, config.Postgres.DbName, driver)
	if err != nil {
		panic(fmt.Errorf("failed to create migrate instance; error %v", err))
	}

	// check if already migrated
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		panic(fmt.Errorf("failed to get database version; error %v", err))
	}
	if dirty {
		fmt.Println("Current DB version is dirty, re-migrating")
		err = m.Force(int(version))
		if err != nil && err != migrate.ErrNoChange {
			panic(fmt.Errorf("failed to migrate database; error %v", err))
		}
	}
	if err == migrate.ErrNilVersion {
		// version is 0
		version = 0
	}
	if version == config.Postgres.DbVersion {
		fmt.Println("database is already migrated")
	} else {
		err = m.Migrate(config.Postgres.DbVersion)
		if err != nil && err != migrate.ErrNoChange {
			panic(fmt.Errorf("failed to migrate database; error %v", err))
		}
	}
	fmt.Println("Grande jefe, has migrado la DB")

}
