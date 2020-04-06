package postHandlers

import (
	"log"

	"github.com/Dayasagara/meeting-scheduler/helpers"
	"github.com/Dayasagara/meeting-scheduler/interfaces"

	"github.com/labstack/echo"
)

func (p *PostHandler) SyncHandler(ctx echo.Context) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()

	mapClaims, tokenErr := helpers.ValidateToken(ctx)
	if tokenErr != nil {
		log.Println(tokenErr)
		return helpers.CommonResponseHandler(400, "Token Error", ctx)
	}

	//Get all unsync events
	events, dbErr := interfaces.DBEngine.GetEvents(int(mapClaims["userID"].(float64)))
	log.Println(events)
	if dbErr != nil {
		log.Println(dbErr)
		return helpers.CommonResponseHandler(400, "DB Error", ctx)
	}
	//if there are any unsync events, sync them with the google calendar
	if len(events) > 0 {
		syncErr := helpers.SyncWithGCalendar(events)
		if syncErr != nil {
			return helpers.CommonResponseHandler(400, "Sync Error", ctx)
		}
	}

	return helpers.CommonResponseHandler(200, "OK", ctx)
}
