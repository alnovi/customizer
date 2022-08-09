package client_config_action

import (
	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type DeleteStaticCategoryConfigAction struct {
	clientRepository *repository.ClientRepository
	configRepository *repository.ConfigRepository
}

func NewDeleteStaticCategoryConfigAction(
	clientRepository *repository.ClientRepository,
	configRepository *repository.ConfigRepository,
) DeleteStaticCategoryConfigAction {
	return DeleteStaticCategoryConfigAction{
		clientRepository: clientRepository,
		configRepository: configRepository,
	}
}

func (s DeleteStaticCategoryConfigAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	muxVars := mux.Vars(request)

	clientId := utils.NewV(muxVars["client_id"])
	category := utils.NewV(muxVars["category"])

	err := s.
		configRepository.
		DeleteCategory(
			clientId.String(),
			"client_config_static",
			category.String(),
		)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "s.configRepository.DeleteCategory").
			Error("client_config_action.static_delete_category_action.DeleteStaticCategoryConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	err = utils.ResponseJson(writer, nil, http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("client_config_action.static_delete_category_action.DeleteStaticCategoryConfigAction.ServeHTTP")

		return
	}
}
