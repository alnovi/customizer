package default_config_actions

import (
	"errors"
	"net/http"

	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
	"github.com/gorilla/mux"
)

type ListDefaultConfigAction struct {
	defaultConfigRepository *repository.DefaultConfigRepository
}

func NewListDefaultConfigAction(defaultConfigRepository *repository.DefaultConfigRepository) ListDefaultConfigAction {
	return ListDefaultConfigAction{
		defaultConfigRepository: defaultConfigRepository,
	}
}

func (g ListDefaultConfigAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	muxVars := mux.Vars(request)

	collection := utils.NewV(muxVars["collection"])

	if collection.IsEmpty() {
		logger.
			WithError(errors.New("collection or category is empty")).
			WithField("flag", "mux.Vars").
			Warn("default_config_actions.default_config_List_action.ListDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	configs, err := g.defaultConfigRepository.List("default_config_" + collection.String())

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "g.defaultConfigRepository.Get").
			Error("default_config_actions.default_config_List_action.ListDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	err = utils.ResponseJson(writer, configs, http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("default_config_actions.default_config_List_action.ListDefaultConfigAction.ServeHTTP")

		return
	}
}
