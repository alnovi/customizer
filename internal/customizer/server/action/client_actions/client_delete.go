package client_actions

import (
	"github.com/gorilla/mux"
	"net/http"

	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
)

type DeleteClientAction struct {
	clientRepository *repository.ClientRepository
}

func NewDeleteClientAction(clientRepository *repository.ClientRepository) DeleteClientAction {
	return DeleteClientAction{
		clientRepository: clientRepository,
	}
}

func (d DeleteClientAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	muxVars := mux.Vars(request)

	clientId := utils.NewV(muxVars["client_id"])

	err := d.clientRepository.DeleteClient(clientId.String())

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "s.configRepository.DeleteClient").
			Error("client_actions.client_delete_action.DeleteClientAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	err = utils.ResponseJson(writer, nil, http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("client_actions.client_delete_action.DeleteClientAction.ServeHTTP")

		return
	}
}
