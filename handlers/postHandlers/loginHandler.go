package postHandlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Dayasagara/meeting-scheduler/helpers"
	"github.com/Dayasagara/meeting-scheduler/model"

	"github.com/Dayasagara/meeting-scheduler/interfaces"
	"github.com/labstack/echo"
)

func (p *PostHandler) LoginHandler(ctx echo.Context) error {
	var user model.User
	reqErr := json.NewDecoder(ctx.Request().Body).Decode(&user)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()
	if reqErr != nil || user.Email == "" || user.Password == "" {
		return helpers.CommonResponseHandler(400, "Req Error", ctx)
	}

	user.Password = helpers.Encrypt(user.Password)

	userExists, userID := interfaces.DBEngine.Authenticate(user.Email, user.Password)
	if userExists != nil {
		log.Println(userExists)
		return helpers.CommonResponseHandler(400, "Invalid user", ctx)
	}

	user.UserID = userID
	token, tokenErr := helpers.CreateToken(user)

	if tokenErr != nil {
		return helpers.CommonResponseHandler(400, "Internal error", ctx)
	}

	apiResponse := helpers.LoginResponse(200, "User Authenticated", token)
	ctx.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(ctx.Response()).Encode(apiResponse)
}
