package auth

import (
	"fmt"

	eapp "github.com/ivzb/semaphore_server/app"
	"github.com/ivzb/semaphore_server/app/middleware/app"
	"github.com/ivzb/semaphore_server/app/shared/consts"
	"github.com/ivzb/semaphore_server/app/shared/request"
	"github.com/ivzb/semaphore_server/app/shared/response"
)

// Handler will authorize HTTP requests
func Handler(app app.App) app.App {
	prevHandler := app.Handler

	app.Handler = func(env *eapp.Env) *response.Message {
		at, err := request.HeaderValue(app.Env.Request, consts.AuthToken)

		if err != nil {
			return response.Unauthorized(fmt.Sprintf(consts.FormatMissing, consts.AuthToken))
		}

		userID, err := app.Env.Token.Decrypt(at)

		if err != nil {
			return response.Unauthorized(fmt.Sprintf(consts.FormatInvalid, consts.AuthToken))
		}

		exists, err := app.Env.DB.User().Exists(userID)

		if err != nil {
			return response.InternalServerError()
		}

		if exists == false {
			return response.Unauthorized(fmt.Sprintf(consts.FormatInvalid, consts.AuthToken))
		}

		(*app.Env).UserID = userID

		return prevHandler(env)
	}

	return app
}
