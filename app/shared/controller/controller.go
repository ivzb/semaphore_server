package controller

import (
	"fmt"

	"github.com/ivzb/semaphore_server"
	"github.com/ivzb/semaphore_server/db"
	"github.com/ivzb/semaphore_server/shared/consts"
	"github.com/ivzb/semaphore_server/shared/form"
	"github.com/ivzb/semaphore_server/shared/response"
)

func GetFormString(env *app.Env, key string, exister db.Exister) (string, *response.Message) {
	value, err := form.StringValue(env.Request, key)

	if err != nil {
		return "", response.BadRequest(fmt.Sprintf(consts.FormatMissing, key))
	}

	exists, err := exister.Exists(value)

	if err != nil {
		env.Log.Error(err)
		return "", response.InternalServerError()
	}

	if !exists {
		return "", response.NotFound(key)
	}

	return value, nil
}
