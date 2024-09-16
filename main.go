package main

import (
	"context"
	"database/sql"
	"fmt"
	"localbe/configuration"
	"localbe/experience"
	"localbe/experience/pg"
	"localbe/gen/experience/v1/v1connect"
	"localbe/port/connect"
	"net/http"
	"time"

	"connectrpc.com/grpcreflect"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

var (
	getExperienceEntryFunc    experience.GetExperienceEntryFunc
	createExperienceEntryFunc experience.CreateExperienceEntryFunc
	getExperienceFunc         experience.GetExperienceFunc
)

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
	m, err := migrate.NewWithDatabaseInstance(sourceURL, "test_db", driver)
	if err != nil {
		panic(fmt.Errorf("failed to create migrate instance; error %v", err))
	}

	// check if already migrated
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		panic(fmt.Errorf("failed to get database version; error %v", err))
	}
	if dirty {
		panic(fmt.Errorf("current DB version is dirty, FIX MANUALLY BEFORE MIGRATING AGAIN"))
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

	experienceRepo := pg.NewExperienceRepository(pgpool)
	e, err := experienceRepo.CreateExperienceEntry(
		ctx,
		"Zeekr Tech. EU",
		"Software Engineer",
		"Sept 2023",
		"",
		"I have done some stuff here",
	)
	if err != nil {
		panic(fmt.Errorf("failed to create an experience entry; error=%v", err))
	}
	fmt.Printf("Dale, has creado tu primer objeto en la DB: %v", e)

	// setup repository functions
	getExperienceEntryFunc = experienceRepo.GetExperienceEntry
	createExperienceEntryFunc = experienceRepo.CreateExperienceEntry
	getExperienceFunc = experienceRepo.GetExperience
	// SETUP SERVICES

	mux := http.NewServeMux()
	// grpc reflector
	reflector := grpcreflect.NewStaticReflector(
		v1connect.ExperienceServiceName,
	)

	experienceService, err := connect.NewExperienceService(createExperienceEntryFunc, getExperienceEntryFunc, getExperienceFunc)
	if err != nil {
		panic(fmt.Errorf("failed to create experience service; error=%w", err))
	}
	path, handler := v1connect.NewExperienceServiceHandler(experienceService)
	mux.Handle(path, handler)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	http.ListenAndServe(
		"localhost:8080",
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
