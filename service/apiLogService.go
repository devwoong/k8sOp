package service

import (
	"context"
	"k8sOp/model"
	"k8sOp/mongo"
	"time"
)

const COLLECTION = "api_log"

type service struct{}

var ApiLogService service = service{}

func (s service) InsertLog(apiLog model.ApiLog) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo.MongoConnection.SwtichCollection(COLLECTION)
	mongo.MongoConnection.CurrentCollection.InsertOne(ctx, apiLog)
}
