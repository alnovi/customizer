package default_config_actions

import (
	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/flatten"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

type StoreDefaultConfigAction struct {
	defaultConfigRepository *repository.DefaultConfigRepository
}

func NewStoreDefaultConfigAction(defaultConfigRepository *repository.DefaultConfigRepository) StoreDefaultConfigAction {
	return StoreDefaultConfigAction{
		defaultConfigRepository: defaultConfigRepository,
	}
}

func (s StoreDefaultConfigAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error

	muxVars := mux.Vars(request)

	collection := utils.NewV(muxVars["collection"])
	category := utils.NewV(muxVars["category"])
	namespace := utils.QueryValue(request, "namespace")
	replace := false

	if val := utils.QueryValue(request, "replace"); val.Bool() {
		replace = true
	}

	if collection.IsEmpty() || category.IsEmpty() {
		logger.
			WithError(errors.New("collection or category is empty")).
			WithField("flag", "mux.Vars").
			Warn("default_config_actions.default_config_store_action.StoreDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	config := new(interface{})

	err = json.NewDecoder(request.Body).Decode(config)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "json.NewDecoder.Decode").
			Warn("default_config_actions.default_config_store_action.StoreDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	newFlattenConfig, err := flatten.NewFlattenFromMap(*config, flatten.DEFAULT_DELIMITER)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "flatten.NewFlattenFromMap").
			Error("default_config_actions.default_config_store_action.StoreDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	newFlattenConfig.SetNamespace(namespace.String())

	err = s.
		defaultConfigRepository.
		Store("default_config_"+collection.String(), category.String(), newFlattenConfig, replace)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "s.defaultConfigRepository.Store").
			Error("default_config_actions.default_config_store_action.StoreDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	storedConfig, err := s.defaultConfigRepository.Get("default_config_"+collection.String(), category.String())

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "s.defaultConfigRepository.Get").
			Error("default_config_actions.default_config_store_action.StoreDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	err = utils.ResponseJson(writer, storedConfig.ToNested(true), http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("default_config_actions.default_config_store_action.StoreDefaultConfigAction.ServeHTTP")

		return
	}
}
