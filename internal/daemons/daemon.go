package daemons

import (
	"github.com/magmel48/go-web/internal/auth"
)

const maxBatchSizeToProcess = 5

// QueryItem is what daemon will accept as payload for a new job.
type QueryItem struct {
	UserID   auth.UserID
	ShortIDs []string
}

// Daemon is simple daemon that can make a job that specified in Run override.
//go:generate mockery --name=Daemon
type Daemon interface {
	Run()
	EnqueueJob(item QueryItem)
}
