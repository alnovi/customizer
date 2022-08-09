package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type DefaultConfigDTO struct {
	Id   string      `bson:"_id"`
	Data primitive.D `bson:"data"`
}
