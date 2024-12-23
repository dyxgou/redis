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

	// A GetDelCommand represents the "GETDEL <key>"
	GetDelCommand struct {
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

	// A IncrCommand represents "INCR <key>" and increments the key by 1.
	IncrCommand struct {
		token token.Token
		key   string
	}

	// A IncrByCommand represents "INCRBY <key> <increment>" and increments the key by increment.
	IncrByCommand struct {
		token     token.Token
		key       string
		increment int
	}

	// A DercCommand represents "DECR <key>" and decrements the key by 1.
	DercCommand struct {
		token token.Token
		key   string
	}

	// A DercByCommand represents "DECRBY <key> <decrement>" and decrements the key by decrement.
	DercByCommand struct {
		token     token.Token
		key       string
		decrement int
	}

	// A MGetCommand stands for Multiple Get and represents "MGET <key> [<key> ...]"
	MGetCommand struct {
		token token.Token
		keys  []string
	}

	// MSetCommand stands for Multiple Set and represents "MSET <key> <value> [<key> <value> ...]"
	MSetCommand struct {
		token token.Token
		pairs []struct {
			key   string
			value string
		}
	}

	// AppendCommand represents "APPEND <key> <value>"
	AppendCommand struct {
		token token.Token
		key   string
		value string
	}

	// ExistsCommand represents "EXISTS <key> <value>"
	ExistsCommand struct {
		token token.Token
		key   string
	}
)
