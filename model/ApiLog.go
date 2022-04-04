package model

import "time"

type ApiLog struct {
	Url         string    `bson:"url"`
	LogDatetime time.Time `bson:"logDatetime,string"`
}
