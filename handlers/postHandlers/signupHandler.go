package postHandlers

import (
	"encoding/json"
	"log"

	"github.com/Dayasagara/meeting-scheduler/helpers"
	"github.com/Dayasagara/meeting-scheduler/model"

	"github.com/Dayasagara/meeting-scheduler/interfaces"
	"github.com/labstack/echo"
)

type PostHandler struct{}

func (p *PostHandler) SignUpHandler(ctx echo.Context) error {
	var user model.User
	reqErr := json.NewDecoder(ctx.Request().Body).Decode(&user)
	ctx.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	defer ctx.Request().Body.Close()
	if reqErr != nil || user.Email == "" || user.Password == "" {
		log.Println(reqErr)
		return helpers.CommonResponseHandler(400, "Request Error", ctx)
	}

	//To encrypt password instead of storing as plaintext
	user.Password = helpers.Encrypt(user.Password)

	userExists := interfaces.DBEngine.CheckUser(user.Email)

	if userExists != nil {
		log.Println(userExists)
		return helpers.CommonResponseHandler(400, "User already exists", ctx)
	}

	err := interfaces.DBEngine.CreateUser(user)

	if err != nil {
		return helpers.CommonResponseHandler(400, "Internal Error", ctx)
	}

	return helpers.CommonResponseHandler(200, "User Created successfully", ctx)
}
