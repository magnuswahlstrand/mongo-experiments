package main

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BaseSuite struct {
	suite.Suite
	service Service
}

type mongoContainer struct {
	testcontainers.Container
	URI string
}

func setupMongo(ctx context.Context) (*mongoContainer, error) {
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

	return &mongoContainer{Container: container, URI: uri}, nil
}

const useTestContainers = false

func (suite *BaseSuite) SetupSuite() {
	ctx := context.Background()

	uri := testMongoURI
	if useTestContainers {
		c, err := setupMongo(ctx)
		require.NoError(suite.T(), err)
		uri = c.URI
	}

	service := NewService(ctx, uri)
	suite.service = service
}

func (suite *BaseSuite) TearDownSuite() {
	suite.service.shutdown()
}

func (suite *BaseSuite) SetupTest() {
	err := suite.service.db.DropCollection(context.Background())
	require.NoError(suite.T(), err)
}

func (suite *BaseSuite) mustListJobs() []Job {
	users, err := suite.service.listJobs(context.Background(), bson.M{})
	require.NoError(suite.T(), err)
	return users
}

func (suite *BaseSuite) mustInsertJob(data bson.M) primitive.ObjectID {
	id, err := suite.service.insertJob(context.Background(), data)
	require.NoError(suite.T(), err)
	return id
}
