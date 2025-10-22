package server

import (
	"context"
)

// Server contract
type Server interface {
	Run(context.Context) error
}
