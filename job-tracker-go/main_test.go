package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

const testMongoURI = "mongodb://localhost:27017"

type ServiceSuite struct {
	suite.Suite
	service Service
}

func (suite *ServiceSuite) SetupSuite() {
	service := NewService(context.Background(), testMongoURI)
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
