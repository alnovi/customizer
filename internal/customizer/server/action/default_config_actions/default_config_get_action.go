package default_config_actions

import (
	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/flatten"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type GetDefaultConfigAction struct {
	defaultConfigRepository *repository.DefaultConfigRepository
}

func NewGetDefaultConfigAction(defaultConfigRepository *repository.DefaultConfigRepository) GetDefaultConfigAction {
	return GetDefaultConfigAction{
		defaultConfigRepository: defaultConfigRepository,
	}
}

func (g GetDefaultConfigAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	muxVars := mux.Vars(request)

	collection := utils.NewV(muxVars["collection"])
	category := utils.NewV(muxVars["category"])
	namespace := utils.QueryValue(request, "namespace")

	if collection.IsEmpty() || category.IsEmpty() {
		logger.
			WithError(errors.New("collection or category is empty")).
			WithField("flag", "mux.Vars").
			Warn("default_config_actions.default_config_store_action.GetDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	config, err := g.defaultConfigRepository.Get("default_config_"+collection.String(), category.String())

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "g.defaultConfigRepository.Get").
			Error("default_config_actions.default_config_store_action.GetDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	if !namespace.IsEmpty() {
		// TODO: If namespace contain ended path then be error
		config = config.Get(namespace.String()).(*flatten.Flatten)
	}

	err = utils.ResponseJson(writer, config.ToNested(true), http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("default_config_actions.default_config_store_action.GetDefaultConfigAction.ServeHTTP")

		return
	}
}
