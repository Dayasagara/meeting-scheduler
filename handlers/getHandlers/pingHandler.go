package getHandlers

import (
	"github.com/Dayasagara/meeting-scheduler/helpers"

	"github.com/labstack/echo"
)

type GetHandler struct{}

//Simple ping API
func (g *GetHandler) PingHandler(ctx echo.Context) error {
	//To set the response type
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return helpers.CommonResponseHandler(200, "OK", ctx)
}
