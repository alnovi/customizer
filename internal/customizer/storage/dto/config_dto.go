package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type ConfigDTO struct {
	Id       string      `bson:"_id"`
	Category string      `bson:"category"`
	Parent   string      `bson:"parent"`
	Data     primitive.D `bson:"data"`
}
