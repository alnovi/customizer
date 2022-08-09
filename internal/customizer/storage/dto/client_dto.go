package dto

type ClientDTO struct {
	Id     string `bson:"_id"`
	Secret string `bson:"secret"`
	Name   string `bson:"name"`
}
