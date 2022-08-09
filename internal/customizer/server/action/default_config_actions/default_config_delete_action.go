package default_config_actions

import (
	"alnovi/customizer/internal/customizer/storage/repository"
	"alnovi/customizer/pkg/logger"
	"alnovi/customizer/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type DeleteDefaultConfigAction struct {
	defaultConfigRepository *repository.DefaultConfigRepository
}

func NewDeleteDefaultConfigAction(defaultConfigRepository *repository.DefaultConfigRepository) DeleteDefaultConfigAction {
	return DeleteDefaultConfigAction{
		defaultConfigRepository: defaultConfigRepository,
	}
}

func (s DeleteDefaultConfigAction) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	muxVars := mux.Vars(request)

	collection := utils.NewV(muxVars["collection"])
	category := utils.NewV(muxVars["category"])
	namespace := utils.NewV(muxVars["namespace"])

	err := s.
		defaultConfigRepository.
		Delete("default_config_"+collection.String(), category.String(), namespace.String())

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "s.defaultConfigRepository.Delete").
			Error("client_config_action.static_delete_action.DeleteDefaultConfigAction.ServeHTTP")

		utils.ResponseCode(writer, http.StatusInternalServerError)

		return
	}

	err = utils.ResponseJson(writer, namespace, http.StatusOK)

	if err != nil {
		logger.
			WithError(err).
			WithField("flag", "utils.ResponseJson").
			Error("client_config_action.static_delete_action.DeleteDefaultConfigAction.ServeHTTP")

		return
	}
}
