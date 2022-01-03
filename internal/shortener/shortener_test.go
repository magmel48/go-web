package shortener

import (
	"fmt"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/shortener/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

func TestShortener_MakeShorter(t *testing.T) {
	type fields struct {
		prefix    string
		links     map[string]string
		userLinks map[string][]string
		backup    Backup
	}
	type args struct {
		url string
	}

	backup := &mocks.Backup{}
	backup.On("Append", mock.AnythingOfType("string")).Return(nil)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "happy path",
			fields: fields{
				prefix:    "http://localhost:8080",
				links:     make(map[string]string),
				userLinks: make(map[string][]string),
				backup:    backup,
			},
			args: args{url: "https://google.com"},
			want: "http://localhost:8080/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				prefix:    tt.fields.prefix,
				links:     tt.fields.links,
				userLinks: tt.fields.userLinks,
				backup:    backup,
			}

			if got, err := s.MakeShorter(tt.args.url, auth.NewUserID()); got != tt.want || err != nil {
				t.Errorf("MakeShorter() = %v, want %v, err %v", got, tt.want, err)
			} else {
				assert.Equal(t, len(backup.Calls), 2)
			}
		})
	}
}

func TestShortener_RestoreLong(t *testing.T) {
	type fields struct {
		prefix string
		links  map[string]string
		backup Backup
	}
	type args struct {
		id string
	}

	backup := &mocks.Backup{}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "happy path",
			fields:  fields{
				prefix: "http://localhost:8080",
				links: map[string]string{
					"https://google.com": "1",
				},
				backup: backup,
			},
			args:    args{id: "1"},
			want:    "https://google.com",
			wantErr: false,
		},
		{
			name:    "unhappy path",
			fields:  fields{prefix: "http://localhost:8080", links: make(map[string]string)},
			args:    args{id: "1"},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				prefix: tt.fields.prefix,
				links:  tt.fields.links,
				backup: &mocks.Backup{},
			}

			got, err := s.RestoreLong(tt.args.id)
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

func TestShortener_storeLink(t *testing.T) {
	type fields struct {
		prefix string
		links  map[string]string
		backup Backup
	}
	type args struct {
		link string
		id   string
	}

	backup := &mocks.Backup{}
	backup.On("Append", mock.AnythingOfType("string")).Return(nil)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name:   "should call Append method of backup field",
			fields: fields{backup: backup},
			args:   args{link: "https://google.com/", id: "2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				prefix: tt.fields.prefix,
				links:  tt.fields.links,
				backup: tt.fields.backup,
			}

			s.storeLink(tt.args.link, tt.args.id)

			assert.Equal(t, fmt.Sprintf("%s|%s\n", tt.args.link, tt.args.id), backup.Calls[0].Arguments[0])
		})
	}
}

func TestShortener_GetUserLinks(t *testing.T) {
	type fields struct {
		prefix    string
		links     map[string]string
		userLinks map[string][]string
		backup    Backup
	}
	type args struct {
		userID auth.UserID
	}

	userID := auth.NewUserID()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []UrlsMap
	}{
		{
			name:    "happy path",
			fields:  fields{
				prefix: "http://localhost:8080",
				links: map[string]string{
					"https://google.com": "1",
				},
				userLinks: map[string][]string{
					userID.String(): {"1"},
				},
			},
			args:    args{userID: userID},
			want:    []UrlsMap{{ShortURL: "http://localhost:8080/1", OriginalURL: "https://google.com"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				prefix:    tt.fields.prefix,
				links:     tt.fields.links,
				userLinks: tt.fields.userLinks,
				backup:    tt.fields.backup,
			}
			if got := s.GetUserLinks(tt.args.userID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUserLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
