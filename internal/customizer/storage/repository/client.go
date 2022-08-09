package repository

import (
	"alnovi/customizer/pkg/utils"
	"context"
	"errors"

	"alnovi/customizer/internal/customizer/storage/dto"
	"alnovi/customizer/pkg/storage/drivers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	clientCollection = "clients"
)

type ClientRepository struct {
	storage *drivers.Mongo
}

func NewClientRepository(storage *drivers.Mongo) *ClientRepository {
	return &ClientRepository{
		storage: storage,
	}
}

func (repository *ClientRepository) GetList(page int64, limit int64) ([]*dto.ClientDTO, error) {
	var err error

	skip := int64(0)

	if limit == 0 {
		limit = 10
	}

	if page > 0 {
		skip = (page - 1) * limit
	}

	findOptions := options.Find().
		SetSkip(skip).
		SetLimit(limit).
		SetSort(bson.M{"_id": 1})

	res, err := repository.storage.CurrentDB().
		Collection(clientCollection).
		Find(context.Background(), bson.M{}, findOptions)

	if err != nil {
		return nil, err
	}

	var result []*dto.ClientDTO

	err = res.All(context.Background(), &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repository *ClientRepository) GetClientInfo(id string) (*dto.ClientDTO, error) {
	var err error

	filter := bson.M{"_id": id}

	result := &dto.ClientDTO{}

	res := repository.storage.CurrentDB().
		Collection(clientCollection).
		FindOne(context.Background(), filter)

	err = res.Decode(&result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (repository *ClientRepository) IsClientExist(id string) bool {
	var err error

	filterOptions := options.FindOne().SetProjection(bson.M{"_id": id})

	var result interface{}

	err = repository.storage.CurrentDB().
		Collection(clientCollection).
		FindOne(
			context.Background(),
			bson.M{"_id": id},
			filterOptions,
		).Decode(&result)

	if err != nil || result == nil {
		return false
	}

	return true
}

func (repository *ClientRepository) CreateClientInfo(dto dto.ClientDTO) (string, error) {
	var err error

	if dto.Id == "" {
		dto.Id = utils.StringRand(10)
	}

	_, err = repository.storage.CurrentDB().
		Collection(clientCollection).
		InsertOne(context.Background(), dto)

	if err != nil {
		return "", err
	}

	return dto.Id, nil
}

func (repository *ClientRepository) DeleteClient(id string) error {

	if repository.IsClientExist(id) == false {
		return errors.New("there is no such client")
	}

	_, err := repository.storage.CurrentDB().
		Collection(clientCollection).
		DeleteOne(
			context.Background(),
			bson.M{"_id": id},
		)

	return err
}
