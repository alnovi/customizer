package drivers

import (
	"alnovi/customizer/pkg/logger"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-redis/redis"
)

const (
	clientKeyTemplate = "configs.client_%s.%s"
	clinicKeyTemplate = "configs.client_%s.clinic_%d.%s"
)

var ErrKeysNotFound = errors.New("keys not found")

type RedisConfig struct {
	Host     string
	Port     uint
	Database int
	Password string
}

type Redis struct {
	client *redis.Client
	config RedisConfig
}

func NewRedis(config RedisConfig) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Database,
	})

	_, err := client.Ping().Result()

	if err != nil {
		return nil, err
	}

	return &Redis{
		client: client,
		config: config,
	}, nil
}

func (r *Redis) HealthCheck() int {
	err := r.client.Ping().Err()

	if err != nil {
		logger.WithError(err).
			Error("Health check for Redis has failed")

		return 0
	}

	return 1
}

func (r *Redis) Disconnect() error {
	return r.client.Close()
}

func (r *Redis) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *Redis) Set(key string, value interface{}) error {
	return r.client.Set(key, value, 0).Err()
}

func (r *Redis) GetStruct(key string, value interface{}) error {
	data, err := r.Get(key)

	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), value)
}

func (r *Redis) SetStruct(key string, value interface{}) error {
	result, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return r.Set(key, result)
}

func (r *Redis) Scan(pattern string) ([]string, error) {
	var keys []string
	var cursor uint64
	var err error

	for {
		var fetchedKeys []string

		scanCmd := r.client.Scan(cursor, pattern+"*", 50)
		fetchedKeys, cursor, err = scanCmd.Result()

		if err != nil {
			return nil, err
		}

		keys = append(keys, fetchedKeys...)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

func (r *Redis) DeleteByPattern(pattern string) error {
	keys, err := r.Scan(pattern)

	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	err = r.client.Del(keys...).Err()

	if err != nil {
		return err
	}

	return nil
}

func (r *Redis) GetClientConfig(id, namespace string) (map[string]interface{}, error) {
	template := fmt.Sprintf(clientKeyTemplate, id, namespace)
	keyPartToTrim := fmt.Sprintf(clientKeyTemplate, id, "")

	return r.getConfigByTemplate(template, keyPartToTrim)
}

func (r *Redis) SetClientConfig(id string, config map[string]interface{}) error {
	for key, value := range config {
		resultKey := fmt.Sprintf(clientKeyTemplate, id, key)

		err := r.Set(resultKey, value)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Redis) GetClinicConfig(id int, clientId, namespace string) (map[string]interface{}, error) {
	template := fmt.Sprintf(clinicKeyTemplate, clientId, id, namespace)
	keyPartToTrim := fmt.Sprintf(clinicKeyTemplate, clientId, id, "")

	return r.getConfigByTemplate(template, keyPartToTrim)
}

func (r *Redis) SetClinicConfig(id int, clientId string, config map[string]interface{}) error {
	for key, value := range config {
		resultKey := fmt.Sprintf(clinicKeyTemplate, clientId, id, key)

		err := r.Set(resultKey, value)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Redis) InvalidateClientConfig(id string) error {
	return r.DeleteByPattern(fmt.Sprintf(clientKeyTemplate, id, ""))
}

func (r *Redis) InvalidateClinicConfig(id int, clientId string) error {
	return r.DeleteByPattern(fmt.Sprintf(clinicKeyTemplate, clientId, id, ""))
}

func (r *Redis) getConfigByTemplate(template, keyPartToTrim string) (map[string]interface{}, error) {
	keys, err := r.Scan(template)

	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return nil, ErrKeysNotFound
	}

	result := make(map[string]interface{})
	for _, key := range keys {
		value, err := r.client.Get(key).Result()

		if err != nil {
			return nil, err
		}

		result[strings.TrimPrefix(key, keyPartToTrim)] = value
	}

	return result, nil
}
