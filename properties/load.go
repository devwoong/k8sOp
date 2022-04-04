package properties

import (
	"github.com/magiconair/properties"
)

var MongoConfig mongoConfig

func (config *mongoConfig) Load() {
	p := properties.MustLoadFile("properties/config.properties", properties.UTF8)
	config.Host = p.MustGetString("mongo.host")
	config.Port = p.GetInt("mongo.port", 8080)
	config.Database = p.MustGetString("mongo.database")
}

type mongoConfig struct {
	Host     string `properties:"mongo.host"`
	Port     int    `properties:"mongo.port"`
	Database string `properties:"mongo.database"`
}
