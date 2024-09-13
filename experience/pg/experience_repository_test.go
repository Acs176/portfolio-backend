package pg

import (
	"localbe/pkg/testdocker"
	"os"
	"testing"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/stretchr/testify/assert"
)

const (
	dbMigrationURL = "file://../../database/migration"
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	dockerSetup := testdocker.New(dbMigrationURL)
	dockerSetup.SetupDocker()
	defer dockerSetup.PurgePool()

	return m.Run()
}

func TestCreateExperienceEntry(t *testing.T) {
	assert.Equal(t, 1, 1)
}
