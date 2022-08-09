package client_actions

import (
	"alnovi/customizer/internal/customizer/storage/dto"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"net/http"

	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
)

type CreateClientRequest struct {
	Id     string `json:"id" valid:"required,minstringlength(10),maxstringlength(10)"`
	Name   string `json:"name" valid:"required,maxstringlength(60)"`
	Secret string `json:"secret" valid:"alphanum,maxstringlength(36)"`
}

type CreateClientResponse struct {
	CreateClientRequest
}

type CreateClientAction struct {
	clientRepository *repository.ClientRepository
}

func NewCreateClientAction(clientRepository *repository.ClientRepository) CreateClientAction {
	return CreateClientAction{
		clientRepository: clientRepository,
	}
}

func (a CreateClientAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	var err error

	newClient := new(CreateClientRequest)

	err = json.NewDecoder(request.Body).Decode(newClient)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "json.NewDecoder").
			Warn("client_action.client_create.CreateClientAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusBadRequest)

		return
	}

	isValid, err := govalidator.ValidateStruct(newClient)

	if !isValid {
		logger.
			WithError(err).
			WithField("flag", "govalidator.ValidateStruct").
			Warn("client_action.client_create.CreateClientAction.ServeHTTP")

		_ = utils.ResponseJson(
			writer,
			utils.H{
				"error": err.Error(),
			},
			http.StatusUnprocessableEntity,
		)

		return
	}

	newClientSecret := utils.StringRand(36)

	if newClient.Secret != "" {
		newClientSecret = newClient.Secret
	}

	clientDto := dto.ClientDTO{
		Id:     newClient.Id,
		Name:   newClient.Name,
		Secret: newClientSecret,
	}

	newClientId, err := a.clientRepository.CreateClientInfo(clientDto)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "c.storage.CreateClientInfo").
			Error("client_action.client_create.CreateClientAction.ServeHTTP")

		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	data, err := a.clientRepository.GetClientInfo(newClientId)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "c.clientRepository.GetClientInfo").
			Error("client_action.client_create.CreateClientAction.ServeHTTP")

		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	response := new(CreateClientResponse)
	response.Id = data.Id
	response.Name = data.Name
	response.Secret = data.Secret

	err = utils.ResponseJson(writer, response, http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("client_action.client_create.CreateClientAction.ServeHTTP")

		return
	}
}
