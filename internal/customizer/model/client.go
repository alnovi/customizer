package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Client struct {
	Id     primitive.ObjectID     `bson:"_id" json:"id"`
	Config map[string]interface{} `bson:"config" json:"config"`
}

func NewClient() *Client {
	return &Client{
		Config: make(map[string]interface{}),
	}
}
