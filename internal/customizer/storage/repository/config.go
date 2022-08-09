package repository

import (
	"alnovi/customizer/internal/customizer/storage/dto"
	utilsInternal "alnovi/customizer/internal/customizer/utils"
	"alnovi/customizer/pkg/flatten"
	"alnovi/customizer/pkg/storage/drivers"
	"alnovi/customizer/pkg/utils"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConfigRepository struct {
	storage *drivers.Mongo
}

func NewConfigRepository(storage *drivers.Mongo) *ConfigRepository {
	return &ConfigRepository{
		storage: storage,
	}
}

func (cr *ConfigRepository) Get(clientId string, collection string, category string) (*flatten.Flatten, string, error) {
	var err error

	if clientId == "" || collection == "" || category == "" {
		return nil, "", errors.New("client id, collection or category is empty")
	}

	sourceResult := new(dto.ConfigDTO)

	err = cr.
		storage.
		CurrentDB().
		Collection(collection).
		FindOne(
			context.Background(),
			bson.M{
				"_id":      clientId + ":" + category,
				"category": category,
			}).
		Decode(sourceResult)

	if err == mongo.ErrNoDocuments {
		return flatten.NewFlatten(), "", nil
	} else if err != nil {
		return nil, "", err
	}

	result, err := flatten.NewFlattenFromMap(utils.NormalizeMongoPrimitive(sourceResult.Data), "#")

	if err != nil {
		return nil, "", err
	}

	return result, sourceResult.Parent, nil
}

func (cr *ConfigRepository) List(clientId string, collection string) ([]map[string]interface{}, error) {
	var err error

	if clientId == "" || collection == "" {
		return nil, errors.New("client id or category is empty")
	}

	cur, err := cr.
		storage.
		CurrentDB().
		Collection(collection).
		Aggregate(
			context.Background(),
			mongo.Pipeline{
				bson.D{{"$match", bson.M{"client_id": clientId}}},
				bson.D{
					{"$project", bson.D{
						{"client_id", 1},
						{"category", 1},
						{"fields", bson.D{
							{"$size", bson.D{{"$objectToArray", "$data"}}},
						}},
					}},
				},
			})

	if err == mongo.ErrNoDocuments {
		return []map[string]interface{}{}, nil
	} else if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0)

	err = cur.All(context.Background(), &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (cr *ConfigRepository) Store(
	clientId string,
	collection string,
	category string,
	parent string,
	dto *flatten.Flatten,
	replace bool,
) error {
	var err error

	if clientId == "" || collection == "" || category == "" {
		return errors.New("client id, collection or category is empty")
	}

	dto = utilsInternal.PreparedStoreData(dto)

	if cr.IsConfigExist(clientId, collection, category) {
		err = cr.update(clientId, collection, category, dto, replace)
	} else {
		err = cr.create(clientId, collection, category, parent, dto)
	}

	return err
}

func (cr *ConfigRepository) Delete(clientId string, collection string, category string, namespace string) error {
	config, _, err := cr.Get(clientId, collection, category)

	if err != nil {
		return err
	}

	for _, key := range config.Keys(namespace) {
		config.Delete(key)
	}

	if len(config.Keys("")) == 0 {
		_, err = cr.
			storage.
			CurrentDB().
			Collection(collection).
			UpdateOne(
				context.Background(),
				bson.M{"_id": clientId + ":" + category},
				bson.M{"$set": bson.M{"data": bson.M{}}},
			)

		return err
	}

	_, err = cr.
		storage.
		CurrentDB().
		Collection(collection).
		UpdateOne(
			context.Background(),
			bson.M{"_id": clientId + ":" + category},
			bson.M{"$set": bson.M{"data": config.ToNested(true)}},
		)

	if err != nil {
		return err
	}

	return err
}

func (cr *ConfigRepository) DeleteCategory(clientId string, collection string, category string) error {
	if !cr.IsConfigExist(clientId, collection, category) {
		return errors.New("there is no such category")
	}

	_, err := cr.
		storage.
		CurrentDB().
		Collection(collection).
		DeleteOne(
			context.Background(),
			bson.M{"_id": clientId + ":" + category},
		)

	return err
}

func (cr *ConfigRepository) IsConfigExist(clientId string, collection string, category string) bool {
	var err error

	filterOptions := options.
		FindOne().
		SetProjection(bson.M{"_id": 1, "category": 1})

	var result interface{}

	err = cr.
		storage.
		CurrentDB().
		Collection(collection).
		FindOne(
			context.Background(),
			bson.M{
				"_id":       clientId + ":" + category,
				"client_id": clientId,
				"category":  category,
			},
			filterOptions,
		).
		Decode(&result)

	if err != nil || result == nil {
		return false
	}

	return true
}

func (cr *ConfigRepository) create(
	clientId string,
	collection string,
	category string,
	parent string,
	dto *flatten.Flatten,
) error {
	var err error

	dto = utilsInternal.DataMongoNormalize(dto)

	storeData := dto.ToNested(true)

	if storeData == nil {
		storeData = bson.M{}
	}

	_, err = cr.
		storage.
		CurrentDB().
		Collection(collection).
		InsertOne(
			context.Background(),
			bson.M{
				"_id":       clientId + ":" + category,
				"client_id": clientId,
				"category":  category,
				"data":      storeData,
				"parent":    parent,
			},
		)

	if err != nil {
		return err
	}

	return nil
}

func (cr *ConfigRepository) update(
	clientId string,
	collection string,
	category string,
	dto *flatten.Flatten,
	replace bool,
) error {
	var err error

	sourceData, _, err := cr.Get(clientId, collection, category)

	if err != nil {
		return err
	}

	var dataStore *flatten.Flatten

	if replace {
		dataStore = utilsInternal.DataMongoNormalize(sourceData)
	} else {
		sourceData = utilsInternal.DataMongoDenormalize(sourceData)

		dataStore = flatten.NewMerge(sourceData, dto).Compile()

		dataStore = utilsInternal.ClearArtifactFromStruct(dataStore)

		dataStore = utilsInternal.DataMongoNormalize(dataStore)
	}

	_, err = cr.
		storage.
		CurrentDB().
		Collection(collection).
		UpdateOne(
			context.Background(),
			bson.M{"_id": clientId + ":" + category},
			bson.M{"$set": bson.M{"data": dataStore.ToNested(true)}},
		)

	if err != nil {
		return err
	}

	return err
}
