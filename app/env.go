package app

import (
	"database/sql"
	"net/http"

	//"github.com/ivzb/semaphore_server/app/db"
	"github.com/ivzb/semaphore_server/app/shared/config"
	"github.com/ivzb/semaphore_server/app/shared/logger"
	"github.com/ivzb/semaphore_server/app/shared/token"
	"github.com/ivzb/semaphore_server/app/shared/uuid"
)

type Env struct {
	Request *http.Request
	//DB      db.DBSourcer
	DB     *sql.DB
	Log    logger.Loggerer
	Token  token.Tokener
	Config *config.Config
	UserID string
	UUID   uuid.UUIDer
}
