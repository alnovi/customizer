package client_config_action

import (
	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
)

type StaticListAction struct {
	clientRepository *repository.ClientRepository
	configRepository *repository.ConfigRepository
}

func NewStaticListAction(
	clientRepository *repository.ClientRepository,
	defaultConfigRepository *repository.ConfigRepository,
) StaticListAction {
	return StaticListAction{
		clientRepository: clientRepository,
		configRepository: defaultConfigRepository,
	}
}

func (g StaticListAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	muxVars := mux.Vars(request)

	clientId := utils.NewV(muxVars["client_id"])

	if clientId.IsEmpty() {
		logger.
			WithError(errors.New("client id is empty")).
			WithField("flag", "mux.Vars").
			Warn("client_config_action.static_list_action.StaticListAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	if !g.clientRepository.IsClientExist(clientId.String()) {
		logger.
			WithError(errors.New("client not found")).
			WithField("flag", "g.clientRepository.IsClientExist").
			Warn("client_config_action.static_list_action.StaticListAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusNotFound)

		return
	}

	statics, err := g.configRepository.List(clientId.String(), "client_config_static")

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "s.configRepository.List").
			Error("default_config_actions.default_config_store_action.StaticListAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	err = utils.ResponseJson(writer, statics, http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("client_config_action.static_list_action.StaticListAction.ServeHTTP")

		return
	}
}
