package db

import "context"

// BaseRepository is high level implementation
type BaseRepository interface {
	// Ping checks connection to database
	Ping(ctx context.Context) error
	Disconnect(ctx context.Context) error
}
