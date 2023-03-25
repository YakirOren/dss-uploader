package db

import (
	"DSS-uploader/config"
	"DSS-uploader/models"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/readpref"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDataStore struct {
	FilesCollection *mongo.Collection
	Client          *mongo.Client
}

func createConnectionString(username string, password string, address string) string {
	return fmt.Sprintf("mongodb://%s:%s@%s", username, password, address)
}

func NewMongoDataStore(config *config.Config) (*MongoDataStore, error) {
	connection := createConnectionString(
		config.MongoDbUsername,
		config.MongoDbPassword,
		config.MongoURL)

	client, err := mongo.NewClient(options.Client().ApplyURI(connection))
	if err != nil {
		return nil, err
	}

	if err := client.Connect(context.TODO()); err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return &MongoDataStore{
		FilesCollection: client.Database(config.DBName).Collection(config.FileCollection),
		Client:          client,
	}, nil
}

func (db *MongoDataStore) WriteFile(ctx context.Context, file models.FileMetadata) (string, error) {
	result, err := db.FilesCollection.InsertOne(ctx, file)
	if err != nil {
		return "", err
	}
	fileID := result.InsertedID.(primitive.ObjectID).Hex()

	return fileID, nil
}

func (db *MongoDataStore) AppendFragment(ctx context.Context, filename string, fragment models.Fragment) error {
	_, err := db.FilesCollection.UpdateOne(
		ctx,
		bson.M{"name": filename},
		bson.M{"$push": bson.M{"fragments": fragment}, "$inc": bson.M{"size": fragment.Size}},
	)

	return err
}

func (db *MongoDataStore) GetMetadataByName(ctx context.Context, name string) (*models.FileMetadata, bool) {
	output := models.FileMetadata{}
	filter := bson.D{{"name", name}}

	if err := db.FilesCollection.FindOne(ctx, filter).Decode(&output); err != nil {
		return nil, false
	}

	return &output, true
}

func (db *MongoDataStore) ListFiles(ctx context.Context) ([]models.FileMetadata, error) {
	var output []models.FileMetadata

	cursor, err := db.FilesCollection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	if err := cursor.All(ctx, &output); err != nil {
		return nil, err
	}

	return output, nil
}
