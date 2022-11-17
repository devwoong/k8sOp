package main

import (
	"k8sOp/channel"
	"k8sOp/service"
	"time"
)

func main() {
	// properties.MongoConfig.Load()
	// mongo.InitConnection(properties.MongoConfig.Host, properties.MongoConfig.Port)
	// mongo.MongoConnection.SwtichDatabase(properties.MongoConfig.Database)

	// TODO logic
	// service.ApiLogService.InsertLog(model.ApiLog{Url: "localhost:8080/test", LogDatetime: time.Now()})
	service.K8sService.Init()

	defer channel.CommonChannel.Destroy()

	go service.RpcService.Start()

	go service.K8sService.CheckAndBootDeployment()

	for {
		time.Sleep(10000 * time.Millisecond)
		service.K8sService.Run()
	}

	// Close
	// mongo.CloseConnection()
}
