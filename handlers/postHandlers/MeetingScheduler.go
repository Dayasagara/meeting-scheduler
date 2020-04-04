package postHandlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Dayasagara/meeting-scheduler/helpers"
	"github.com/Dayasagara/meeting-scheduler/model"

	"github.com/Dayasagara/meeting-scheduler/interfaces"
	"github.com/labstack/echo"
)

func (p *PostHandler) MeetingScheduler(ctx echo.Context) error {
	var event model.ScheduleEvent
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()

	req := ctx.Request().Header
	token := req.Get("token")

	//decrypt the token and get the jwt claims
	mapClaims, tokenErr := helpers.DecryptToken(token)
	if tokenErr != nil {
		log.Println(tokenErr)
		return helpers.CommonResponseHandler(400, "Invalid token", ctx)
	}

	reqErr := json.NewDecoder(ctx.Request().Body).Decode(&event)
	if reqErr != nil {
		log.Println(reqErr)
		return helpers.CommonResponseHandler(400, "Request Error", ctx)
	}

	//Authenticate the token claims
	userExists, _ := interfaces.DBEngine.Authenticate(fmt.Sprintf("%v", mapClaims["email"]), fmt.Sprintf("%v", mapClaims["password"]))
	if userExists != nil {
		log.Println(userExists)
		return helpers.CommonResponseHandler(400, "Invalid user", ctx)
	}

	if !interfaces.DBEngine.CheckAvailability(event.UserID, event.Date, event.StartingFrom) {
		return helpers.CommonResponseHandler(400, "Slot not available", ctx)
	}

	err := interfaces.DBEngine.ScheduleEvent(event)

	if err != nil {
		return helpers.CommonResponseHandler(400, "Internal Error", ctx)
	}

	return helpers.CommonResponseHandler(200, "Event Created successfully", ctx)
}
