package getHandlers

import (
	"encoding/json"
	"fmt"
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
	req := ctx.Request().Header
	token := req.Get("token")
	targetUserID, _ := strconv.Atoi(req.Get("targetUserID"))
	date := ctx.Param("date")
	if !helpers.ValidateDate(date) {
		return helpers.CommonResponseHandler(400, "Invalid date", ctx)
	}
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
