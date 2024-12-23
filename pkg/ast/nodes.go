package ast

import (
	"github/dyxgou/redis/pkg/token"
	"time"
)

type (
	// A GetCommand represents the command "GET <key>"
	GetCommand struct {
		token token.Token
		key   string
	}

	// A SetCommand represents the command "SET <key> <value>" if Ex is not provided it is assumed that the key won't be expired
	SetCommand struct {
		token token.Token
		key   string
		value string
		Ex    time.Duration

		// Sets the key if not exists
		Nx bool
		// Sets the key if it exists
		Xx bool
	}

	// A GetSetCommand represents the "GETSET <key> <value>" and acts like a SET followed by a GET
	GetSetCommand struct {
		token token.Token
		key   string
		value string
	}

	// A GetExCommand represents the "GETEX <key> <value> <EX> <time>" if the Ex is not provided, the parser throws an error
	GetExCommand struct {
		token token.Token
		key   string
		value string
		Ex    time.Duration
	}
)
