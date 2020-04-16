package getHandlers

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Dayasagara/meeting-scheduler/helpers"
	"github.com/Dayasagara/meeting-scheduler/interfaces"
	"github.com/Dayasagara/meeting-scheduler/model"

	"github.com/labstack/echo"
)

func (g *GetHandler) GetSlotsHandler(ctx echo.Context) error {
	var slots []model.AvailabilitySlots
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()

	req := ctx.Request().Header
	targetUserID, idConvErr := strconv.Atoi(req.Get("targetUserID"))
	if idConvErr != nil {
		return helpers.CommonResponseHandler(400, "Invalid Target User ID", ctx)
	}
	date := ctx.Param("date")

	if !helpers.ValidateDate(date) || !helpers.CheckPastDate(date) {
		return helpers.CommonResponseHandler(400, "Invalid date", ctx)
	}

	_, tokenErr := helpers.ValidateToken(ctx)
	if tokenErr != nil {
		log.Println(tokenErr)
		return helpers.CommonResponseHandler(400, "Token Error", ctx)
	}

	slots, dbErr := interfaces.DBEngine.GetAvailability(targetUserID, date)

	var formattedSlots = []model.AvailabilitySlots{}
	var formattedSlot model.AvailabilitySlots

	for _, slot := range slots {
		formattedSlot.StartSlot = slot.StartTime.String()[11:16]
		formattedSlot.EndSlot = slot.StartTime.Add(time.Hour).String()[11:16]
		formattedSlot.Availability = slot.Availability
		formattedSlots = append(formattedSlots, formattedSlot)
	}

	if dbErr != nil {
		log.Println(dbErr)
		return helpers.CommonResponseHandler(400, "DB Error", ctx)
	}
	ctx.Response().WriteHeader(200)
	return json.NewEncoder(ctx.Response()).Encode(formattedSlots)

}
