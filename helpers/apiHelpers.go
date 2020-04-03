package helpers

import (
	"github.com/Dayasagara/meeting-scheduler/model"
)

//ResponseMapper maps http code and message to response struct
func ResponseMapper(code int, message string) model.APIResponse {
	response := model.APIResponse{
		Code:    code,
		Message: message,
	}
	return response
}
