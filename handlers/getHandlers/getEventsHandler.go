package getHandlers

import (
	"encoding/json"
	"log"

	"github.com/Dayasagara/meeting-scheduler/helpers"
	"github.com/Dayasagara/meeting-scheduler/interfaces"
	"github.com/Dayasagara/meeting-scheduler/model"

	"github.com/labstack/echo"
)

func (g *GetHandler) GetEventsHandler(ctx echo.Context) error {
	var slots []model.ScheduleEvent
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()
	mapClaims, tokenErr := helpers.ValidateToken(ctx)
	if tokenErr != nil {
		log.Println(tokenErr)
		return helpers.CommonResponseHandler(400, "Token Error", ctx)
	}

	slots, dbErr := interfaces.DBEngine.GetMyEvents(int(mapClaims["userID"].(float64)))

	var formattedSlots = []model.ScheduleEvent{}

	for _, slot := range slots {
		slot.StartingFrom = slot.StartingFrom[11:16]
		slot.EndingTill = slot.EndingTill[11:16]
		slot.Date = slot.Date[0:10]
		formattedSlots = append(formattedSlots, slot)
	}

	if dbErr != nil {
		log.Println(dbErr)
		return helpers.CommonResponseHandler(400, "DB Error", ctx)
	}
	ctx.Response().WriteHeader(200)
	return json.NewEncoder(ctx.Response()).Encode(formattedSlots)

}
