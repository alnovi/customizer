package client_config_action

import (
	"errors"
	"net/http"

	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/flatten"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
	"encoding/json"
	"github.com/gorilla/mux"
)

type StoreStaticConfigAction struct {
	clientRepository *repository.ClientRepository
	configRepository *repository.ConfigRepository
}

func NewStoreStaticConfigAction(
	clientRepository *repository.ClientRepository,
	defaultConfigRepository *repository.ConfigRepository,
) StoreStaticConfigAction {
	return StoreStaticConfigAction{
		clientRepository: clientRepository,
		configRepository: defaultConfigRepository,
	}
}

func (s StoreStaticConfigAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var err error

	replace := false
	muxVars := mux.Vars(request)

	clientId := utils.NewV(muxVars["client_id"])
	category := utils.NewV(muxVars["category"])
	namespace := utils.QueryValue(request, "namespace")
	parent := ""

	if clientId.IsEmpty() || category.IsEmpty() {
		logger.
			WithError(errors.New("client id or category is empty")).
			WithField("flag", "mux.Vars").
			Warn("client_config_action.static_get_action.StoreStaticConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	if val := utils.QueryValue(request, "replace"); !val.IsEmpty() {
		replace = true
	}

	config := make(map[string]interface{})

	err = json.NewDecoder(request.Body).Decode(&config)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "json.NewDecoder.Decode").
			Warn("client_config_action.static_get_action.StoreStaticConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	if _, ok := config["data"]; !ok {
		logger.
			WithError(errors.New("config data is empty")).
			WithField("flag", "config").
			Warn("client_config_action.static_get_action.StoreStaticConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	if val, ok := config["parent"]; ok {
		parent = val.(string)
	}

	newFlattenConfig, err := flatten.NewFlattenFromMap(config["data"], flatten.DEFAULT_DELIMITER)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "flatten.NewFlattenFromMap").
			Error("client_config_action.static_get_action.StoreStaticConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	newFlattenConfig.SetNamespace(namespace.String())

	err = s.
		configRepository.
		Store(
			clientId.String(),
			"client_config_static",
			category.String(),
			parent,
			newFlattenConfig,
			replace,
		)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "s.configRepository.Store").
			Error("client_config_action.static_get_action.StoreStaticConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	storedConfig, _, err := s.
		configRepository.
		Get(
			clientId.String(),
			"client_config_static",
			category.String(),
		)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "s.configRepository.Get").
			Error("default_config_actions.default_config_store_action.StoreStaticConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	err = utils.ResponseJson(writer, storedConfig.ToNested(true), http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("default_config_actions.default_config_store_action.StoreStaticConfigAction.ServeHTTP")

		return
	}
}
