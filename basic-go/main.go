package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"os"
	"time"
)

type UserInfo struct {
	ID        string    `bson:"id"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type Service struct {
	db       *qmgo.QmgoClient
	shutdown func()
}

func (s *Service) upsertUser(ctx context.Context, id uuid.UUID, name string) error {
	user := UserInfo{
		ID:        id.String(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := s.db.Upsert(ctx, bson.M{"id": user.ID}, user)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) findUsers(ctx context.Context) ([]UserInfo, error) {
	var batch []UserInfo
	if err := s.db.Find(ctx, bson.M{}).All(&batch); err != nil {
		return nil, err
	}

	return batch, nil
}

func NewService(ctx context.Context, uri string) Service {
	mongoConfig := &qmgo.Config{
		Uri:             uri,
		Database:        "service",
		Coll:            "user",
		MaxPoolSize:     nil,
		MinPoolSize:     nil,
		SocketTimeoutMS: nil,
		ReadPreference:  nil,
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

	users, err := service.findUsers(ctx)
	if err != nil {
		log.Println("failed to get users", err)
		return
	}
	for _, user := range users {
		fmt.Println(user.ID, user.Name)
	}

	//err := service.upsertUser(ctx, uuid.MustParse("7d90007a-1668-43a2-8c42-d675268af20c"), "Magnus")
	//if err != nil {
	//	log.Println("failed to insert user", err)
	//	return
	//}
}
