package main

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatesSuite struct {
	BaseSuite
}

func (suite *StatesSuite) getAndLockJob(id primitive.ObjectID) *Job {
	job, err := suite.service.getAndLockJob(context.Background(), id)
	require.NoError(suite.T(), err)
	return job
}

func (suite *StatesSuite) TestChangeState() {
	t := suite.T()

	// Act
	id := suite.mustInsertJob(bson.M{"name": "magnus"})

	job := suite.getAndLockJob(id)

	// Assert
	require.Equal(t, job.State, "new")
	require.Equal(t, job.Locked, true)
}

func (suite *StatesSuite) TestErrorWhenLocked() {
	t := suite.T()

	// Arrange
	id := suite.mustInsertJob(bson.M{"name": "magnus"})
	_ = suite.getAndLockJob(id)

	// Act
	_, err := suite.service.getAndLockJob(context.Background(), id)

	require.ErrorIs(t, err, mongo.ErrNoDocuments)
}
