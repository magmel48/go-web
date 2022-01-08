package shortener

import (
	"context"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db/links"
	linksmocks "github.com/magmel48/go-web/internal/db/links/mocks"
	"github.com/magmel48/go-web/internal/db/userlinks"
	userlinksmocks "github.com/magmel48/go-web/internal/db/userlinks/mocks"
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
		&links.Link{ShortID: "1"}, false, nil)

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

			if got, _, err := s.MakeShorter(context.Background(), tt.args.url, auth.NewUserID()); got != tt.want || err != nil {
				t.Errorf("MakeShorter() = %v, want %v, err %v", got, tt.want, err)
			}
		})
	}
}

func TestShortener_RestoreLong(t *testing.T) {
	type fields struct {
		prefix              string
		linksRepository     links.Repository
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
				prefix: "http://localhost:8080",
				linksRepository: &withLinksRepository,
			},
			args:    args{id: "1"},
			want:    "https://google.com",
			wantErr: false,
		},
		{
			name: "unhappy path",
			fields: fields{
				prefix: "http://localhost:8080",
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
				prefix: tt.fields.prefix,
				linksRepository: tt.fields.linksRepository,
			}

			got, err := s.RestoreLong(context.Background(), tt.args.id)
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
