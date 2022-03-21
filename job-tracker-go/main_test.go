package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

const testMongoURI = "mongodb://localhost:27017"

type ServiceSuite struct {
	suite.Suite
	service Service
}

type nginxContainer struct {
	testcontainers.Container
	URI string
}

func setupMongo(ctx context.Context) (*nginxContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	ip, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("mongodb://%s:%s", ip, mappedPort.Port())

	return &nginxContainer{Container: container, URI: uri}, nil
}

func (suite *ServiceSuite) SetupSuite() {
	ctx := context.Background()

	c, err := setupMongo(ctx)
	require.NoError(suite.T(), err)

	service := NewService(ctx, c.URI)
	suite.service = service
}

func (suite *ServiceSuite) TearDownSuite() {
	suite.service.shutdown()
}

func (suite *ServiceSuite) SetupTest() {
	err := suite.service.db.DropCollection(context.Background())
	require.NoError(suite.T(), err)
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (suite *ServiceSuite) TestInsert() {
	t := suite.T()

	// Act
	suite.mustInsertJob(bson.M{"name": "magnus"})

	// Assert
	users := suite.mustListJobs()
	require.Len(t, users, 1)
	require.Equal(t, users[0].Data, bson.M{"name": "magnus"})
}

func (suite *ServiceSuite) mustListJobs() []Job {
	users, err := suite.service.listJobs(context.Background())
	require.NoError(suite.T(), err)
	return users
}

func (suite *ServiceSuite) mustInsertJob(data bson.M) {
	err := suite.service.insertJob(context.Background(), data)
	require.NoError(suite.T(), err)
}
