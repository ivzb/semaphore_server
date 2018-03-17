package config

import (
	"encoding/json"
	"fmt"

	"github.com/ivzb/semaphore_server/app/shared/database"
	"github.com/ivzb/semaphore_server/app/shared/server"
	"github.com/ivzb/semaphore_server/app/shared/token"
)

// *****************************************************************************
// Application Settings
// *****************************************************************************

// Config contains the application settings
type Config struct {
	Database database.Info `json:"Database"`
	// Email    email.SMTPInfo  `json:"Email"`
	Server server.Info `json:"Server"`
	Token  token.Info  `json:"Token"`
}

// New config instance
func New(jsonConf []byte) (*Config, error) {
	conf := &Config{}

	if err := json.Unmarshal(jsonConf, &conf); err != nil {
		return nil, fmt.Errorf("Could not parse config: %v", err)
	}

	return conf, nil
}
