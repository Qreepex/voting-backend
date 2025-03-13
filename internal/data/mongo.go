package data

import (
	ctx "context"
	"fmt"
	"log"

	"github.com/qreepex/voting-backend/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database struct {
	Client *mongo.Client
}

func InitDatabase() (*Database, error) {
	clientOptions := options.
		Client().
		ApplyURI("mongodb://localhost:27017/vote")

	client, err := mongo.Connect(ctx.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	return &Database{Client: client}, nil
}

func (db *Database) GetDb() *mongo.Database {
	return db.Client.Database("vote")
}

func (db *Database) CreateVote(vote model.Vote) (*mongo.InsertOneResult, error) {
	coll := db.GetVotesCollection()

	return coll.InsertOne(ctx.TODO(), vote)
}

func (db *Database) GetVotesCollection() *mongo.Collection {
	return db.GetDb().Collection("votes")
}

func (db *Database) GetCandidatesCollection() *mongo.Collection {
	return db.GetDb().Collection("candidates")
}

func (db *Database) GetCampaignsCollection() *mongo.Collection {
	return db.GetDb().Collection("campaigns")
}

func (db *Database) GetCandidate(id string) (*model.Candidate, error) {
	coll := db.GetCandidatesCollection()

	var candidate model.Candidate

	err := coll.FindOne(ctx.TODO(), map[string]string{"id": id}, &options.FindOneOptions{Projection: map[string]int{"_id": 0}}).Decode(&candidate)

	if err != nil {
		return nil, err
	}

	return &candidate, nil
}

func (db *Database) GetCandidates() ([]model.Candidate, error) {
	coll := db.GetCandidatesCollection()

	cursor, err := coll.Find(ctx.TODO(), map[string]string{}, &options.FindOptions{Projection: map[string]int{"_id": 0}})
	if err != nil {
		return nil, err
	}

	var candidates []model.Candidate

	err = cursor.All(ctx.TODO(), &candidates)
	if err != nil {
		return nil, err
	}

	return candidates, nil
}

func (db *Database) GetCampaign(id string) (*model.Campaign, error) {
	coll := db.GetCampaignsCollection()

	var campaign model.Campaign

	err := coll.FindOne(ctx.TODO(), map[string]string{"id": id}, &options.FindOneOptions{Projection: map[string]int{"_id": 0}}).Decode(&campaign)

	if err != nil {
		return nil, err
	}

	return &campaign, nil
}

func (db *Database) GetCandidateVotes(candidateId string) ([]model.Vote, error) {
	coll := db.GetVotesCollection()

	cursor, err := coll.Find(ctx.TODO(), map[string]string{"candidate": candidateId}, &options.FindOptions{Projection: map[string]int{"_id": 0}})
	if err != nil {
		return nil, err
	}

	var votes []model.Vote

	err = cursor.All(ctx.TODO(), &votes)
	if err != nil {
		return nil, err
	}

	return votes, nil
}
