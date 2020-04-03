package getHandlers

import (
	"encoding/json"

	"github.com/Dayasagara/meeting-scheduler/helpers"

	"github.com/Dayasagara/meeting-scheduler/model"
	"github.com/labstack/echo"
)

type GetHandler struct{}

//Simple ping API
func (g *GetHandler) PingHandler(ctx echo.Context) error {
	var response model.APIResponse
	//To set the response type
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	response = helpers.ResponseMapper(200, "OK")
	return json.NewEncoder(ctx.Response()).Encode(response)
}
