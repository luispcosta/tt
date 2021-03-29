package persistence

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/luispcosta/go-tt/configuration"
	"github.com/luispcosta/go-tt/core"
	"github.com/luispcosta/go-tt/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepository struct {
	Client        *mongo.Client
	Database      *mongo.Database
	Ctx           context.Context
	CancelContext context.CancelFunc
	Clock         utils.Clock
}

func NewMongoRepository() (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		cancel()
		return nil, err
	}

	return &MongoRepository{Client: client, Ctx: ctx, CancelContext: cancel, Clock: utils.NewLiveClock()}, nil
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

func (repo *MongoRepository) Find(activityNameOrAlias string) (*core.Activity, error) {
	var activity *core.Activity
	repo.activityCollection().FindOne(repo.Ctx, bson.D{{"Name", activityNameOrAlias}}).Decode(&activity)
	if activity == nil {
		err := repo.activityCollection().FindOne(repo.Ctx, bson.D{{"Alias", activityNameOrAlias}}).Decode(&activity)
		if err != nil {
			return nil, err
		}
	}

	return activity, nil
}

func (repo *MongoRepository) Update(activity core.Activity) error {
	return nil
}

func (repo *MongoRepository) Delete(activityNameOrAlias string) error {
	res, err := repo.activityCollection().DeleteOne(repo.Ctx, bson.D{{"Name", activityNameOrAlias}})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		res, err = repo.activityCollection().DeleteOne(repo.Ctx, bson.D{{"Alias", activityNameOrAlias}})
		if err != nil {
			return err
		}
		if res.DeletedCount == 0 {
			return errors.New("activity not found")
		}
		return nil
	} else {
		return nil
	}
}

func (repo *MongoRepository) List() ([]core.Activity, error) {
	cursor, err := repo.activityCollection().Find(repo.Ctx, bson.D{})
	if err != nil {
		fmt.Println(err.Error())
		return []core.Activity{}, err
	}
	defer cursor.Close(repo.Ctx)
	var result []core.Activity
	for cursor.Next(repo.Ctx) {
		var act core.Activity
		errDecode := cursor.Decode(&act)
		if errDecode != nil {
			return result, errDecode
		}
		result = append(result, act)
	}

	if err := cursor.Err(); err != nil {
		return result, err
	}
	return result, nil
}

func (repo *MongoRepository) Start(activity core.Activity) error {
	instant := repo.Clock.Now()
	year := instant.Year()
	month := instant.Month()
	day := instant.Day()
	date := fmt.Sprintf("%s-%s-%s", strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day))
	_, err := repo.activityLogCollection().InsertOne(repo.Ctx, bson.D{{"ActivityID", activity.ID}, {"Date", date}, {"Start", int64(repo.Clock.Now().Unix())}})
	if err != nil {
		return err
	}
	return nil
}

func (repo *MongoRepository) Stop(activity core.Activity) error {
	instant := repo.Clock.Now()
	year := instant.Year()
	month := instant.Month()
	day := instant.Day()
	date := fmt.Sprintf("%s-%s-%s", strconv.Itoa(year), strconv.Itoa(int(month)), strconv.Itoa(day))
	var activityLog core.ActivityLog
	repo.activityLogCollection().FindOne(repo.Ctx, bson.D{{"ActivityID", activity.ID}, {"Date", date}}).Decode(&activityLog)
	activityLog.End = repo.Clock.Now().Unix()
	updateStmt := bson.M{"$set": bson.D{{"End", repo.Clock.Now().Unix()}}}
	_, err := repo.activityLogCollection().UpdateByID(repo.Ctx, activityLog.ID, updateStmt)
	if err != nil {
		return err
	}
	return nil
}

func (repo *MongoRepository) activityCollection() *mongo.Collection {
	return repo.Database.Collection("activities")
}

func (repo *MongoRepository) activityLogCollection() *mongo.Collection {
	return repo.Database.Collection("activity_logs")
}

func (repo *MongoRepository) activityToBson(activity core.Activity) bson.D {
	fmt.Println("Would like to add")
	fmt.Println(bson.D{{"Name", activity.Name}, {"Alias", activity.Alias}, {"Description", activity.Description}})
	return bson.D{{"Name", activity.Name}, {"Alias", activity.Alias}, {"Description", activity.Description}}
}
