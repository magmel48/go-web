package daemons

import (
	"context"
	"github.com/magmel48/go-web/internal/db/userlinks"
	"github.com/magmel48/go-web/internal/db/userlinks/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"strconv"
	"testing"
)

func TestDeletingRecordsDaemon_deleteLinks(t *testing.T) {
	type fields struct {
		ctx        context.Context
		repository userlinks.Repository
		items      chan QueryItem
	}

	ctx := context.TODO()
	userID := "test_user_id"
	shortID := "1"

	repositoryMock := &mocks.Repository{}
	repositoryMock.On("DeleteLinks", mock.Anything, mock.Anything).Return(nil)

	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "should call proper methods with proper payload",
			fields: fields{ctx: ctx, repository: repositoryMock, items: make(chan QueryItem, 10)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			daemon := &DeletingRecordsDaemon{
				ctx:        tt.fields.ctx,
				repository: tt.fields.repository,
				items:      tt.fields.items,
			}

			daemon.EnqueueJob(QueryItem{UserID: &userID, ShortIDs: []string{shortID}})
			daemon.deleteLinks()

			assert.Equal(t, len(repositoryMock.Calls), 1)
			assert.Equal(t, repositoryMock.Calls[0].Method, "DeleteLinks")
			assert.Equal(t, len(repositoryMock.Calls[0].Arguments), 2)

			args := []userlinks.DeleteQueryItem{{UserID: &userID, ShortIDs: []string{shortID}}}
			if !reflect.DeepEqual(repositoryMock.Calls[0].Arguments[1], args) {
				t.Errorf("got = %v, want = %v", repositoryMock.Calls[0].Arguments[1], args)
			}
		})
	}
}

func BenchmarkDeletingRecordsDaemon_DeleteLinks(b *testing.B) {
	repository := &mocks.Repository{}
	repository.On("DeleteLinks", mock.Anything, mock.Anything).Return(nil)

	daemon := DeletingRecordsDaemon{
		ctx:        context.TODO(),
		repository: repository,
		items:      make(chan QueryItem, 100),
	}

	userID := "test_user_id"
	shortIDs := make([]string, 1000000)
	for i := 0; i < 1000000; i++ {
		shortIDs[i] = strconv.Itoa(i + 1)
	}

	for i := 0; i < maxBatchSizeToProcess; i++ {
		daemon.items <- QueryItem{UserID: &userID, ShortIDs: shortIDs}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		daemon.deleteLinks()
	}
}
