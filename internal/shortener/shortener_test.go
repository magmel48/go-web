package shortener

import (
	"context"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/daemons"
	"github.com/magmel48/go-web/internal/daemons/mocks"
	"github.com/magmel48/go-web/internal/db"
	"github.com/magmel48/go-web/internal/db/links"
	linksmocks "github.com/magmel48/go-web/internal/db/links/mocks"
	"github.com/magmel48/go-web/internal/db/userlinks"
	userlinksmocks "github.com/magmel48/go-web/internal/db/userlinks/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestShortener_MakeShorter(t *testing.T) {
	type fields struct {
		prefix              string
		linksRepository     links.Repository
		userLinksRepository userlinks.Repository
	}
	type args struct {
		url string
	}

	linksRepository := linksmocks.Repository{}
	linksRepository.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(
		&links.Link{ShortID: "1"}, nil)

	userLinksRepository := userlinksmocks.Repository{}
	userLinksRepository.On(
		"FindByLinkID", mock.Anything, mock.Anything, mock.Anything).Return(&userlinks.UserLink{}, nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "happy path",
			fields: fields{
				prefix:              "http://localhost:8080",
				linksRepository:     &linksRepository,
				userLinksRepository: &userLinksRepository,
			},
			args: args{url: "https://google.com"},
			want: "http://localhost:8080/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				prefix:              tt.fields.prefix,
				linksRepository:     tt.fields.linksRepository,
				userLinksRepository: tt.fields.userLinksRepository,
			}

			if got, err := s.MakeShorter(context.TODO(), tt.args.url, auth.NewUserID()); got != tt.want {
				t.Errorf("MakeShorter() = %v, want %v, err %v", got, tt.want, err)
			}
		})
	}
}

func TestShortener_RestoreLong(t *testing.T) {
	type fields struct {
		prefix          string
		linksRepository links.Repository
	}
	type args struct {
		id string
	}

	withLinksRepository := linksmocks.Repository{}
	withLinksRepository.On("FindByShortID", mock.Anything, mock.Anything).Return(
		&links.Link{ShortID: "1", OriginalURL: "https://google.com"}, nil)

	withoutLinksRepository := linksmocks.Repository{}
	withoutLinksRepository.On("FindByShortID", mock.Anything, mock.Anything).Return(nil, nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "happy path",
			fields: fields{
				prefix:          "http://localhost:8080",
				linksRepository: &withLinksRepository,
			},
			args:    args{id: "1"},
			want:    "https://google.com",
			wantErr: false,
		},
		{
			name: "unhappy path",
			fields: fields{
				prefix:          "http://localhost:8080",
				linksRepository: &withoutLinksRepository,
			},
			args:    args{id: "1"},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				prefix:          tt.fields.prefix,
				linksRepository: tt.fields.linksRepository,
			}

			got, err := s.RestoreLong(context.TODO(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RestoreLong() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("RestoreLong() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortener_DeleteURLs(t *testing.T) {
	type fields struct {
		ctx                 context.Context
		prefix              string
		database            db.DB
		linksRepository     links.Repository
		userLinksRepository userlinks.Repository
		daemon              daemons.Daemon
	}
	type args struct {
		userID   auth.UserID
		shortIDs []string
	}

	mockDaemon := mocks.Daemon{}
	mockDaemon.On("EnqueueJob", mock.Anything)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "should call EnqueueJob of daemon",
			fields: fields{
				daemon: &mockDaemon,
			},
			args: args{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				ctx:                 tt.fields.ctx,
				prefix:              tt.fields.prefix,
				database:            tt.fields.database,
				linksRepository:     tt.fields.linksRepository,
				userLinksRepository: tt.fields.userLinksRepository,
				daemon:              tt.fields.daemon,
			}

			s.DeleteURLs(tt.args.userID, tt.args.shortIDs)
			assert.Equal(t, len(mockDaemon.Calls), 1)
			assert.Equal(t, mockDaemon.Calls[0].Method, "EnqueueJob")
		})
	}
}
