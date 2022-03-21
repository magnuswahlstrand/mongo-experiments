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

func (s *Service) insertJob(ctx context.Context, data bson.M) error {
	t := time.Now()
	user := Job{
		ID:        primitive.NewObjectIDFromTimestamp(t),
		State:     "new",
		Locked:    false,
		UpdatedAt: t,
		Data:      data,
	}
	_, err := s.db.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) listJobs(ctx context.Context) ([]Job, error) {
	var batch []Job
	if err := s.db.Find(ctx, bson.M{}).All(&batch); err != nil {
		return nil, err
	}

	return batch, nil
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

	users, err := service.listJobs(ctx)
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
