package postHandlers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Dayasagara/meeting-scheduler/helpers"
	"github.com/Dayasagara/meeting-scheduler/interfaces"
	"github.com/Dayasagara/meeting-scheduler/model"

	"github.com/labstack/echo"
)

func (p *PostHandler) DefineAvHandler(ctx echo.Context) error {
	var availablility model.Availability
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	req := ctx.Request().Header
	token := req.Get("token")

	//decrypt the token and get the jwt claims
	mapClaims, tokenErr := helpers.DecryptToken(token)
	if tokenErr != nil {
		log.Println(tokenErr)
		return helpers.CommonResponseHandler(400, "Invalid token", ctx)
	}

	reqErr := json.NewDecoder(ctx.Request().Body).Decode(&availablility)
	defer ctx.Request().Body.Close()
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
