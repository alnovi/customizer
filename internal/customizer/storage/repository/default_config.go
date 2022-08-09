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

type DefaultConfigRepository struct {
	storage *drivers.Mongo
}

func NewDefaultConfigRepository(storage *drivers.Mongo) *DefaultConfigRepository {
	return &DefaultConfigRepository{
		storage: storage,
	}
}

func (dcr *DefaultConfigRepository) Get(collection string, category string) (*flatten.Flatten, error) {
	var err error

	if collection == "" || category == "" {
		return nil, errors.New("collection or category is empty")
	}

	sourceResult := new(dto.DefaultConfigDTO)

	err = dcr.
		storage.
		CurrentDB().
		Collection(collection).
		FindOne(
			context.Background(),
			bson.M{"_id": category},
		).
		Decode(sourceResult)

	if err == mongo.ErrNoDocuments {
		return flatten.NewFlatten(), nil
	} else if err != nil {
		return nil, err
	}

	result, err := flatten.NewFlattenFromMap(utils.NormalizeMongoPrimitive(sourceResult.Data), "#")

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (dcr *DefaultConfigRepository) List(collection string) ([]map[string]interface{}, error) {
	var err error

	if collection == "" {
		return nil, errors.New("collection or category is empty")
	}

	cur, err := dcr.
		storage.
		CurrentDB().
		Collection(collection).
		Aggregate(
			context.Background(),
			mongo.Pipeline{
				bson.D{{
					"$project", bson.D{{
						"fields", bson.D{{
							"$size", bson.D{{
								"$objectToArray", "$data",
							}},
						}},
					}},
				}},
			},
		)

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

func (dcr *DefaultConfigRepository) Store(collection string, category string, dto *flatten.Flatten, replace bool) error {
	var err error

	if collection == "" || category == "" {
		return errors.New("collection or category is empty")
	}

	dto = utilsInternal.PreparedStoreData(dto)

	if dcr.IsConfigExist(collection, category) {
		err = dcr.update(collection, category, dto, replace)
	} else {
		err = dcr.create(collection, category, dto)
	}

	return err
}

func (dcr *DefaultConfigRepository) IsConfigExist(collection string, categoryId string) bool {
	filterOptions := options.
		FindOne().
		SetProjection(bson.M{"_id": 1})

	var result interface{}

	err := dcr.storage.
		CurrentDB().
		Collection(collection).
		FindOne(
			context.Background(),
			bson.M{"_id": categoryId},
			filterOptions,
		).
		Decode(&result)

	if err != nil || result == nil {
		return false
	}

	return true
}

func (dcr *DefaultConfigRepository) Delete(collection string, categoryId string, namespace string) error {
	config, err := dcr.Get(collection, categoryId)

	if err != nil {
		return err
	}

	for _, key := range config.Keys(namespace) {
		config.Delete(key)
	}

	if len(config.Keys("")) == 0 {
		_, err = dcr.
			storage.
			CurrentDB().
			Collection(collection).
			DeleteOne(
				context.Background(),
				bson.M{"_id": categoryId},
			)

		return err
	}

	_, err = dcr.
		storage.
		CurrentDB().
		Collection(collection).
		UpdateOne(
			context.Background(),
			bson.M{"_id": categoryId},
			bson.M{"$set": bson.M{"data": config.ToNested(true)}},
		)

	if err != nil {
		return err
	}

	return nil
}

func (dcr *DefaultConfigRepository) create(collection string, category string, dto *flatten.Flatten) error {
	var err error
	var storeData interface{}

	dto = utilsInternal.DataMongoNormalize(dto)

	storeData = dto.ToNested(true)

	if storeData == nil {
		storeData = bson.M{}
	}

	_, err = dcr.
		storage.
		CurrentDB().
		Collection(collection).
		InsertOne(
			context.Background(),
			bson.M{"_id": category, "data": storeData},
		)

	if err != nil {
		return err
	}

	return nil
}

func (dcr *DefaultConfigRepository) update(
	collection string,
	category string,
	dto *flatten.Flatten,
	replace bool,
) error {
	var err error
	var storeData interface{}

	sourceData, err := dcr.Get(collection, category)

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

	storeData = dataStore.ToNested(true)

	if storeData == nil {
		storeData = bson.M{}
	}

	_, err = dcr.
		storage.
		CurrentDB().
		Collection(collection).
		UpdateOne(
			context.Background(),
			bson.M{"_id": category},
			bson.M{"$set": bson.M{"data": storeData}},
		)

	if err != nil {
		return err
	}

	return err
}
