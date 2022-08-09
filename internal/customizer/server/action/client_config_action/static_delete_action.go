package client_config_action

import (
	"net/http"

	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
	"github.com/gorilla/mux"
)

type DeleteStaticConfigAction struct {
	clientRepository *repository.ClientRepository
	configRepository *repository.ConfigRepository
}

func NewDeleteStaticConfigAction(
	clientRepository *repository.ClientRepository,
	configRepository *repository.ConfigRepository,
) DeleteStaticConfigAction {
	return DeleteStaticConfigAction{
		clientRepository: clientRepository,
		configRepository: configRepository,
	}
}

func (s DeleteStaticConfigAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	muxVars := mux.Vars(request)

	clientId := utils.NewV(muxVars["client_id"])
	category := utils.NewV(muxVars["category"])
	namespace := utils.NewV(muxVars["namespace"])

	err := s.
		configRepository.
		Delete(
			clientId.String(),
			"client_config_static",
			category.String(),
			namespace.String(),
		)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "s.configRepository.Delete").
			Error("client_config_action.static_delete_action.DeleteStaticConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	err = utils.ResponseJson(writer, namespace, http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("client_config_action.static_delete_action.DeleteStaticConfigAction.ServeHTTP")

		return
	}
}
