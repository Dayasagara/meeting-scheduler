package postHandlers

import (
	"encoding/json"
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

	_, tokenErr := helpers.ValidateToken(ctx)
	if tokenErr != nil {
		log.Println(tokenErr)
		return helpers.CommonResponseHandler(400, "Token Error", ctx)
	}

	reqErr := json.NewDecoder(ctx.Request().Body).Decode(&event)
	if reqErr != nil || !helpers.ValidateDate(event.Date) || !helpers.ValidateTime(event.StartingFrom) || !helpers.ValidateTime(event.EndingTill) {
		log.Println(reqErr)
		return helpers.CommonResponseHandler(400, "Request Error", ctx)
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
