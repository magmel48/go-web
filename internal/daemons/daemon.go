package daemons

import (
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db/userlinks"
)

const maxBatchSizeToProcess = 5

type QueryItem struct {
	UserID   auth.UserID
	ShortIDs []string
}

// Daemon is simple daemon that can make a job that specified in Run override.
//go:generate mockery --name=Daemon
type Daemon interface {
	Run(userLinksRepository userlinks.Repository)
	EnqueueJob(item QueryItem)
}
