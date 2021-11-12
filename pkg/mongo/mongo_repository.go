package repository

import (
	"github.com/Kamva/mgm/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

// MongoRepository is the interface to MongoDB
type MongoRepository struct {
	Model mgm.Model
}

// NewMongoRepository creates a new MongoRepository
func NewMongoRepository(model mgm.Model) Repository {
	return &MongoRepository{model}
}

var readPref int

// Read Preferences
const (
	Primary = iota
	PrimaryPreferred
	Secondary
	SecondaryPreferred
)

// InitMongoWithReadPref configures Mongo with custom read preference
func InitMongoWithReadPref(
	dbName string,
	url string,
	pref int,
) (
	error,
) {
	readPrefMap := map[int]*readpref.ReadPref{
		Primary: readpref.Primary(),
		PrimaryPreferred: readpref.PrimaryPreferred(),
		Secondary: readpref.Secondary(),
		SecondaryPreferred: readpref.SecondaryPreferred(),
	}

	readPref = pref
	return mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(url).SetReadPreference(readPrefMap[pref]))
}

// InitMongo configures MongoDB with default primary read preference
func InitMongo(
	dbName string,
	url string,
) (
	error,
) {
	readPref = 0
	return mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(url))
}

// Create document
func (repo *MongoRepository) Create() error {
	return mgm.Coll(repo.Model).Create(repo.Model)
}

// FindOne finds a document by filter condition
func (repo *MongoRepository) FindOne(
	filter interface{},
) (
	error,
) {
	model := repo.Model
	coll := mgm.Coll(model)
	err := coll.First(filter, model)
	return err
}

// Find returns a cursor pointer to all documents that match the filter condition
func (repo *MongoRepository) Find(
	filter interface{},
) (
	*mongo.Cursor,
	error,
) {
	coll := mgm.Coll(repo.Model)
	ctx := mgm.NewCtx(10 * time.Second)

	findOneErr := repo.FindOne(filter)
	if findOneErr != nil {
		return nil, findOneErr
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	return cursor, err
}

