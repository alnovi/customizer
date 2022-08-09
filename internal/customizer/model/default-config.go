package model

type DefaultConfig struct {
	Id     string                 `bson:"_id"`
	Config map[string]interface{} `bson:"config"`
}

func NewDefaultConfig(id string) *DefaultConfig {
	return &DefaultConfig{
		Id:     id,
		Config: make(map[string]interface{}),
	}
}
