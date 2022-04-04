package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConnection struct {
	Client            *mongo.Client
	CurrentDatabase   *mongo.Database
	CurrentCollection *mongo.Collection
}

var MongoConnection mongoConnection

func InitConnection(host string, port int) {
	MongoConnection = mongoConnection{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoUrl := fmt.Sprintf("mongodb://%s:%d", host, port)
	var err error
	MongoConnection.Client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return
	}
}

func (conn *mongoConnection) SwtichDatabase(database string) {
	conn.CurrentDatabase = conn.Client.Database(database)
}

func (conn *mongoConnection) SwtichCollection(collection string) {
	conn.CurrentCollection = conn.CurrentDatabase.Collection(collection)
}

func (conn *mongoConnection) SwtichDatabaseAndCollection(database string, collection string) {
	conn.CurrentDatabase = conn.Client.Database(database)
	conn.CurrentCollection = conn.CurrentDatabase.Collection(collection)
}

func CloseConnection() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	MongoConnection.Client.Disconnect(ctx)
}
