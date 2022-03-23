package main

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"time"
)

const (
	mongoDatabase   = "job-tracker"
	mongoCollection = "jobs"
)

type Job struct {
	ID        primitive.ObjectID `bson:"_id"`
	State     string             `bson:"state"`
	Locked    bool               `bson:"locked"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Data      bson.M             `bson:"data"`
}

type Service struct {
	db       *qmgo.QmgoClient
	shutdown func()
}

func (s *Service) insertJob(ctx context.Context, data bson.M) (primitive.ObjectID, error) {
	t := time.Now()
	user := Job{
		ID:        primitive.NewObjectIDFromTimestamp(t),
		State:     "new",
		Locked:    false,
		UpdatedAt: t,
		Data:      data,
	}
	res, err := s.db.InsertOne(ctx, user)

	if err != nil {
		return primitive.ObjectID{}, err
	}
	return res.InsertedID.(primitive.ObjectID), nil
}

func (s *Service) listJobs(ctx context.Context, filter bson.M) ([]Job, error) {
	var batch []Job
	if err := s.db.Find(ctx, filter).All(&batch); err != nil {
		return nil, err
	}

	return batch, nil
}

func (s *Service) getAndLockJob(ctx context.Context, id primitive.ObjectID) (*Job, error) {
	filter := bson.M{"_id": id, "locked": false}
	change := qmgo.Change{
		Update: bson.M{
			"$set": bson.M{
				"locked": true,
			},
		},
		ReturnNew: true,
	}

	var job Job
	if err := s.db.Find(ctx, filter).Apply(change, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

func NewService(ctx context.Context, uri string) Service {
	mongoConfig := &qmgo.Config{
		Uri:      uri,
		Database: mongoDatabase,
		Coll:     mongoCollection,
	}

	db, err := qmgo.Open(ctx, mongoConfig)
	if err != nil {
		log.Fatalln(db)
	}

	shutdown := func() {
		ctx := context.Background()
		if err := db.Close(ctx); err != nil {
			fmt.Println("an error occurred, ignoring", err)
			return
		}
	}

	return Service{
		db:       db,
		shutdown: shutdown,
	}
}

func main() {
	uri := os.Getenv("MONGODB_URI")
	ctx := context.Background()
	service := NewService(ctx, uri)
	defer service.shutdown()

	users, err := service.listJobs(ctx, bson.M{})
	if err != nil {
		log.Println("failed to get users", err)
		return
	}
	for _, user := range users {
		fmt.Println(user.ID)
	}

	//err := service.insertJob(ctx, uuid.MustParse("7d90007a-1668-43a2-8c42-d675268af20c"), "Magnus")
	//if err != nil {
	//	log.Println("failed to insert user", err)
	//	return
	//}
}
