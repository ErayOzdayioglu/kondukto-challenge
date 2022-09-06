package repository

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"kondukto.com/challenge/domain"
)

type ScanResultsRepositoryDB struct {
	ScanResultsCollection *mongo.Collection
}

type ScanResultsRepository interface {
	Insert(scanResults domain.ScanResults) (string, error)
	Update(id string, scanResults domain.ScanResults)
	Find(id string) (domain.ScanResults, int)
}

func (s ScanResultsRepositoryDB) Insert(scanResults domain.ScanResults) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := s.ScanResultsCollection.InsertOne(ctx, scanResults)

	if result.InsertedID == nil || err != nil {
		log.Println(err)
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s ScanResultsRepositoryDB) Find(id string) (domain.ScanResults, int) {
	var scanResult domain.ScanResults

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	Id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": Id}
	err := s.ScanResultsCollection.FindOne(ctx, filter).Decode(&scanResult)
	if err != nil {
		log.Println(err.Error())
		return scanResult, -1

	}
	scanResult.Id = Id
	return scanResult, 0

}

func (s ScanResultsRepositoryDB) Update(id string, scanResults domain.ScanResults) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": Id}

	_, err := s.ScanResultsCollection.ReplaceOne(ctx, filter, scanResults)

	if err != nil {
		log.Println(err.Error())
	}

}

func NewScanResultsRepositoryDB(dbClient *mongo.Collection) ScanResultsRepositoryDB {
	return ScanResultsRepositoryDB{ScanResultsCollection: dbClient}
}
