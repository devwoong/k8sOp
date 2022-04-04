package model

import "time"

type Metric struct {
	Cpu          string    `bson:"cpu"`
	Memory       string    `bson:"memory"`
	BaseDatetime time.Time `bson:"baseDatetime"`
}
