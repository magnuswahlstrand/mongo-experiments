package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

func (suite *ServiceSuite) TestUpsert() {
	t := suite.T()

	// Arrange
	id := uuid.New()

	// Act
	suite.mustUpsertUser(id, "Magnxs")
	suite.mustUpsertUser(id, "Magnus")

	// Assert
	users := suite.mustGetUsers()
	require.Len(t, users, 1)
	require.Equal(t, users[0].Name, "Magnus")
	require.Equal(t, users[0].ID, id.String())
}

func (suite *ServiceSuite) TestInsertMultiple() {
	t := suite.T()

	// Arrange
	id1 := uuid.New()
	id2 := uuid.New()

	// Act
	suite.mustUpsertUser(id1, "Magnus")
	suite.mustUpsertUser(id2, "Adam")

	// Assert
	users := suite.mustGetUsers()
	require.Len(t, users, 2)

	assert.Equal(t, users[0].Name, "Magnus")
	assert.Equal(t, users[0].ID, id1.String())

	assert.Equal(t, users[1].Name, "Adam")
	assert.Equal(t, users[1].ID, id2.String())
}

func (suite *ServiceSuite) mustGetUsers() []UserInfo {
	ctx := context.Background()
	users, err := suite.service.findUsers(ctx)
	require.NoError(suite.T(), err)
	return users
}

func (suite *ServiceSuite) mustUpsertUser(id uuid.UUID, name string) {
	err := suite.service.upsertUser(context.Background(), id, name)
	require.NoError(suite.T(), err)
}
