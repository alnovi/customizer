package storage

import "alnovi/customizer/internal/customizer/model"

type Cache interface {
	Connect() error
	Disconnect() error
	Get(key string) (string, error)
	Set(key string, value interface{}) error
	GetStruct(key string, value interface{}) error
	SetStruct(key string, value interface{}) error
	Scan(template string) ([]string, error)
	DeleteByPattern(pattern string) error
	GetClientConfig(id, namespace string) (map[string]interface{}, error)
	SetClientConfig(id string, config map[string]interface{}) error
	GetClinicConfig(id int, clientId, namespace string) (map[string]interface{}, error)
	SetClinicConfig(id int, clientId string, config map[string]interface{}) error
	InvalidateClientConfig(id string) error
	InvalidateClinicConfig(id int, clientId string) error
}

type HealthChecker interface {
	HealthCheck() int
}

type DefaultConfigRepositoryInterface interface {
	Get(id string) (*model.DefaultConfig, error)
	Update(id string, defaultConfig *model.DefaultConfig) error
}
