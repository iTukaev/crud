//go:build integration
// +build integration

package integration

import (
	"context"
	"log"
	"testing"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/suite"

	apiDataPkg "gitlab.ozon.dev/iTukaev/homework/internal/api/data"
	userPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user"
	postgresPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/postgres"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
	"gitlab.ozon.dev/iTukaev/homework/tests/integration/tdb"
)

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(repositorySuite))
}

type repositorySuite struct {
	suite.Suite
	ctx context.Context

	db   *pgxpool.Pool
	user pb.UserServer

	dockerPool *dockertest.Pool
	resource   *dockertest.Resource
}

func (s *repositorySuite) SetupSuite() {
	s.ctx = context.Background()
	var err error
	s.dockerPool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	s.resource, err = s.dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.4",
		Env: []string{
			"POSTGRES_USER=" + tdb.User,
			"POSTGRES_PASSWORD=" + tdb.Password,
			"POSTGRES_DB=" + tdb.DBName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: tdb.Host, HostPort: tdb.Port},
			},
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start docker resource: %s", err)
	}

	if err = s.dockerPool.Retry(func() error {
		s.db, err = tdb.NewTestDB(s.ctx)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	_, err = s.db.Exec(s.ctx, tableCreate)
	if err != nil {
		log.Fatalf("Could not create table: %s", err)
	}
	logger := loggerPkg.NewFatal()
	data := postgresPkg.New(s.db, logger)
	user := userPkg.New(data, logger)
	s.user = apiDataPkg.New(user, logger)
}

func (s *repositorySuite) TearDownSuite() {
	if err := s.dockerPool.Purge(s.resource); err != nil {
		log.Fatalf("Could not purge docker resource: %s", err)
	}

	s.db.Close()
}

func (s *repositorySuite) SetupTest() {
	_, err := s.db.Exec(s.ctx, insertUsers)
	s.Require().NoError(err, "INSERT data to table")
}

func (s *repositorySuite) TearDownTest() {
	_, err := s.db.Exec(s.ctx, deleteUsers)
	s.Require().NoError(err, "DELETE data from table")
}
