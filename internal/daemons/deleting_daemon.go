package daemons

import (
	"context"
	"github.com/magmel48/go-web/internal/db/userlinks"
	"log"
	"time"
)

// DeletingRecordsDaemon deletes link records periodically.
type DeletingRecordsDaemon struct {
	ctx   context.Context
	items chan QueryItem
}

func NewDeletingRecordsDaemon(ctx context.Context) *DeletingRecordsDaemon {
	return &DeletingRecordsDaemon{
		ctx:   ctx,
		items: make(chan QueryItem, 100),
	}
}

func (daemon *DeletingRecordsDaemon) EnqueueJob(item QueryItem) {
	daemon.items <- item
}

func (daemon *DeletingRecordsDaemon) Run(userLinksRepository userlinks.Repository) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-daemon.ctx.Done():
			log.Println("stopped links deletion processing")

			close(daemon.items)
			return

		case <-ticker.C:
			log.Println("processing new requests for links deletion")

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
				if err := userLinksRepository.DeleteLinks(daemon.ctx, items); err != nil {
					log.Println("the error occurred while link deletion", err)
				}
			}
		}
	}
}
