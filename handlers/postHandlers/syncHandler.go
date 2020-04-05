package postHandlers

import (
	"fmt"
	"log"

	"github.com/Dayasagara/meeting-scheduler/helpers"
	"github.com/Dayasagara/meeting-scheduler/interfaces"

	"github.com/labstack/echo"
)

func (p *PostHandler) SyncHandler(ctx echo.Context) error {
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	req := ctx.Request().Header
	token := req.Get("token")

	//decrypt the token and get the jwt claims
	mapClaims, tokenErr := helpers.DecryptToken(token)
	if tokenErr != nil {
		log.Println(tokenErr)
		return helpers.CommonResponseHandler(400, "Invalid token", ctx)
	}
	defer ctx.Request().Body.Close()

	//Authenticate the token claims
	userExists, _ := interfaces.DBEngine.Authenticate(fmt.Sprintf("%v", mapClaims["email"]), fmt.Sprintf("%v", mapClaims["password"]))
	if userExists != nil {
		log.Println(userExists)
		return helpers.CommonResponseHandler(400, "Invalid user", ctx)
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
