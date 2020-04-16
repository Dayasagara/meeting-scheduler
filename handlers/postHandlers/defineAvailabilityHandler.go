package postHandlers

import (
	"encoding/json"
	"github.com/Dayasagara/meeting-scheduler/helpers"
	"github.com/Dayasagara/meeting-scheduler/interfaces"
	"github.com/Dayasagara/meeting-scheduler/model"
	"github.com/labstack/echo"
	"log"
)

func (p *PostHandler) DefineAvHandler(ctx echo.Context) error {
	var availablility model.Availability
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()

	reqErr := json.NewDecoder(ctx.Request().Body).Decode(&availablility)
	if !helpers.ValidateDate(availablility.Date) || !helpers.ValidateTime(availablility.StartSlot) || !helpers.ValidateTime(availablility.EndSlot) {
		return helpers.CommonResponseHandler(400, "Invalid date or time", ctx)
	}

	if !helpers.CheckPastDate(availablility.Date) {
		return helpers.CommonResponseHandler(400, "Past Date", ctx)
	}

	if reqErr != nil {
		log.Println(reqErr)
		return helpers.CommonResponseHandler(400, "Request Error", ctx)
	}

	mapClaims, tokenErr := helpers.ValidateToken(ctx)
	if tokenErr != nil {
		log.Println(tokenErr)
		return helpers.CommonResponseHandler(400, "Token Error", ctx)
	}

	availablility.UserID = int(mapClaims["userID"].(float64))
	if interfaces.DBEngine.CheckForDuplicate(availablility.UserID, availablility.Date) {
		return helpers.CommonResponseHandler(400, "Slots already defined", ctx)
	}

	dbErr := interfaces.DBEngine.DefineAvailability(availablility)
	if dbErr != nil {
		log.Println(dbErr)
		return helpers.CommonResponseHandler(400, "DB Error", ctx)
	}
	return helpers.CommonResponseHandler(200, "OK", ctx)
}
