package client_config_action

import (
	"errors"
	"net/http"

	"alnovi/customizer/internal/customizer/storage/repository"
	utilsInternal "alnovi/customizer/internal/customizer/utils"
	"alnovi/customizer/pkg/flatten"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
	"github.com/gorilla/mux"
)

type GetStaticConfigNormalizeAction struct {
	clientRepository        *repository.ClientRepository
	configRepository        *repository.ConfigRepository
	defaultConfigRepository *repository.DefaultConfigRepository
}

func NewGetStaticConfigNormalizeAction(clientRepository *repository.ClientRepository,
	configRepository *repository.ConfigRepository,
	defaultConfigRepository *repository.DefaultConfigRepository) GetStaticConfigNormalizeAction {
	return GetStaticConfigNormalizeAction{
		clientRepository:        clientRepository,
		configRepository:        configRepository,
		defaultConfigRepository: defaultConfigRepository,
	}
}

func (g GetStaticConfigNormalizeAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	muxVars := mux.Vars(request)

	clientId := utils.NewV(muxVars["client_id"])
	category := utils.NewV(muxVars["category"])
	namespace := utils.QueryValue(request, "namespace")

	if clientId.IsEmpty() || category.IsEmpty() {
		logger.
			WithError(errors.New("client id or category is empty")).
			WithField("flag", "mux.Vars").
			Warn("client_config_action.static_get_action.GetStaticConfigNormalizeAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	if !g.clientRepository.IsClientExist(clientId.String()) {
		logger.
			WithError(errors.New("client not found")).
			WithField("flag", "g.clientRepository.IsClientExist").
			Warn("client_config_action.static_get_action.GetStaticConfigNormalizeAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusNotFound)

		return
	}

	config, parent, err := g.configRepository.Get(clientId.String(), "client_config_static", category.String())

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "g.configRepository.Get").
			Error("client_config_action.static_get_action.GetStaticConfigNormalizeAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	defaultDataStatic, err := g.defaultConfigRepository.Get("default_config_statics", parent)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "g.defaultConfigRepository.Get").
			Error("client_config_action.static_get_action.GetStaticConfigNormalizeAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	config = flatten.NewMerge(defaultDataStatic, config).Compile().SetDelimiter("#")

	if !namespace.IsEmpty() {
		// TODO: If namespace contain ended path then be error
		config = config.Get(namespace.String()).(*flatten.Flatten)
	}

	normalize := utilsInternal.PreparedViewData(config)

	err = utils.ResponseJson(writer, normalize.ToNested(true), http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("default_config_actions.default_config_store_action.GetStaticConfigNormalizeAction.ServeHTTP")

		return
	}
}
