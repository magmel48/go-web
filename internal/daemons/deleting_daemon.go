package daemons

import (
	"context"
	"github.com/magmel48/go-web/internal/db/userlinks"
	"log"
	"time"
)

// DeletingRecordsDaemon deletes link records periodically.
type DeletingRecordsDaemon struct {
	ctx        context.Context
	repository userlinks.Repository
	items      chan QueryItem
}

// NewDeletingRecordsDaemon create new daemon that deletes links asynchronously.
func NewDeletingRecordsDaemon(ctx context.Context, repository userlinks.Repository) *DeletingRecordsDaemon {
	return &DeletingRecordsDaemon{
		ctx:        ctx,
		repository: repository,
		items:      make(chan QueryItem, 100),
	}
}

// EnqueueJob enqueues new items (links) for deletion.
func (daemon *DeletingRecordsDaemon) EnqueueJob(item QueryItem) {
	daemon.items <- item
}

// Run runs links deletion.
func (daemon *DeletingRecordsDaemon) Run() {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-daemon.ctx.Done():
			log.Println("stopped links deletion processing")

			close(daemon.items)
			return

		case <-ticker.C:
			log.Println("processing new requests for links deletion")
			daemon.deleteLinks()
		}
	}
}

// deleteLinks deletes links by batches.
func (daemon *DeletingRecordsDaemon) deleteLinks() {
	items := make([]userlinks.DeleteQueryItem, 0)
	for i := 0; i < maxBatchSizeToProcess; i++ {
		select {
		case item := <-daemon.items:
			items = append(items, userlinks.DeleteQueryItem{
				UserID:   item.UserID,
				ShortIDs: item.ShortIDs,
			})
		default:
		}
	}

	if len(items) > 0 {
		if err := daemon.repository.DeleteLinks(daemon.ctx, items); err != nil {
			log.Println("the error occurred while link deletion", err)
		}
	}
}
