package helpers

import (
	"encoding/json"

	"github.com/labstack/echo"

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

func LoginResponse(code int, message string, token string) model.CreateResponse {
	var response model.CreateResponse
	response = model.CreateResponse{
		Code:    code,
		Message: message,
		Token:   token,
	}
	return response
}
func CommonResponseHandler(code int, message string, ctx echo.Context) error {
	response := ResponseMapper(code, message)
	ctx.Response().WriteHeader(code)
	return json.NewEncoder(ctx.Response()).Encode(response)
}
