package main

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"testing"
)

const testMongoURI = "mongodb://localhost:27017"

type BasicSuite struct {
	BaseSuite
}

func TestBasic(t *testing.T) {
	suite.Run(t, new(BasicSuite))
}

func TestStates(t *testing.T) {
	suite.Run(t, new(StatesSuite))
}

func (suite *BasicSuite) TestInsert() {
	t := suite.T()

	// Act
	id := suite.mustInsertJob(bson.M{"name": "magnus"})

	// Assert
	users := suite.mustListJobs()
	require.Len(t, users, 1)
	require.Equal(t, users[0].ID, id)
	require.Equal(t, users[0].Data, bson.M{"name": "magnus"})
}
