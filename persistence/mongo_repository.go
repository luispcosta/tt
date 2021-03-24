package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/luispcosta/go-tt/configuration"
	"github.com/luispcosta/go-tt/core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client        *mongo.Client
	Database      *mongo.Database
	Ctx           context.Context
	CancelContext context.CancelFunc
}

func NewMongoRepository() (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		cancel()
		return nil, err
	}

	return &MongoRepository{Client: client, Ctx: ctx, CancelContext: cancel}, nil
}

func (repo *MongoRepository) Initialize() error {
	repo.Database = repo.Client.Database("gott")
	return nil
}

func (repo *MongoRepository) Shutdown() error {
	repo.CancelContext()
	return repo.Client.Disconnect(repo.Ctx)
}

func (repo *MongoRepository) SchemaMigrate(direction string) error {
	mongoMigrations, err := configuration.NewMongoMigrations()
	if err != nil {
		return err
	}

	if direction == "up" {
		return mongoMigrations.MigrateAll()
	} else if direction == "down" {
		return mongoMigrations.DownAll()
	} else {
		return fmt.Errorf("don't know how to handle migrations for direction %s", direction)
	}
}

func (repo *MongoRepository) Add(activity core.Activity) error {
	_, err := repo.activityCollection().InsertOne(repo.Ctx, repo.activityToBson(activity))
	return err
}

func (repo *MongoRepository) FindLogsForDay(day time.Time) (core.ActivityDayLog, error) {
	return core.ActivityDayLog{}, nil
}

func (repo *MongoRepository) Find(activityNameOrAlias string) (*core.Activity, error) {
	return &core.Activity{}, nil
}

func (repo *MongoRepository) Update(activity core.Activity) error {
	return nil
}

func (repo *MongoRepository) Delete(activityName string) error {
	return nil
}

func (repo *MongoRepository) List() []core.Activity {
	return []core.Activity{}
}

func (repo *MongoRepository) Start(activity core.Activity) error {
	return nil
}

func (repo *MongoRepository) Stop(activity core.Activity) error {
	return nil
}

func (repo *MongoRepository) Purge() error {
	return nil
}

func (repo *MongoRepository) Backup(destination string) (string, error) {
	return "", nil
}

func (repo *MongoRepository) Restore(restoreFilePath string) error {
	return nil
}

func (repo *MongoRepository) activityCollection() *mongo.Collection {
	return repo.Database.Collection("activities")
}

func (repo *MongoRepository) activityToBson(activity core.Activity) bson.D {
	return bson.D{{"Name", activity.Name}, {"Alias", activity.Alias}, {"Description", activity.Description}}
}
