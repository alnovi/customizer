package client_actions

import (
	"net/http"

	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
)

type ClientsResponse []clientData

type clientData struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

type ListClientAction struct {
	clientRepository *repository.ClientRepository
}

func NewListClientAction(clientRepository *repository.ClientRepository) ListClientAction {
	return ListClientAction{
		clientRepository: clientRepository,
	}
}

func (a ListClientAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	page := utils.QueryValue(request, "page")
	limit := utils.QueryValue(request, "limit")

	data, err := a.clientRepository.GetList(page.Int64(), limit.Int64())

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "l.clientRepository.GetList").
			Error("actions.store.ListClientAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)
	}

	result := make(ClientsResponse, len(data))

	for index, value := range data {
		result[index] = clientData{
			Id:     value.Id,
			Name:   value.Name,
			Secret: value.Secret,
		}
	}

	err = utils.ResponseJson(writer, result, http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("actions.store.ListClientAction.ServeHTTP")
	}
}
