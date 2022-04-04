package main

import (
	"k8sOp/model"
	"k8sOp/mongo"
	"k8sOp/properties"
	"k8sOp/service"
	"time"
)

func main() {
	properties.MongoConfig.Load()
	mongo.InitConnection(properties.MongoConfig.Host, properties.MongoConfig.Port)
	mongo.MongoConnection.SwtichDatabase(properties.MongoConfig.Database)

	// TODO logic
	service.ApiLogService.InsertLog(model.ApiLog{Url: "localhost:8080/test", LogDatetime: time.Now()})

	// Close
	mongo.CloseConnection()
}
