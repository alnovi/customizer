package model

type Clinic struct {
	Id       int                    `bson:"_id" json:"id"`
	ClientId string                 `bson:"client_id" json:"client_id"`
	Config   map[string]interface{} `bson:"config" json:"config"`
}

func NewClinic(id int, clientId string) *Clinic {
	return &Clinic{
		Id:       id,
		ClientId: clientId,
		Config:   make(map[string]interface{}),
	}
}
