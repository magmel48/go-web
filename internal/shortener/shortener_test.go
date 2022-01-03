package shortener

import (
	"fmt"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/shortener/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestShortener_MakeShorter(t *testing.T) {
	type fields struct {
		prefix string
		links  map[string]string
		backup Backup
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
			name:   "happy path",
			fields: fields{prefix: "http://localhost:8080", links: make(map[string]string), backup: backup},
			args:   args{url: "https://google.com"},
			want:   "http://localhost:8080/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				prefix: tt.fields.prefix,
				links:  tt.fields.links,
				backup: backup,
			}

			if got, err := s.MakeShorter(tt.args.url, auth.NewUserID()); got != tt.want || err != nil {
				t.Errorf("MakeShorter() = %v, want %v, err %v", got, tt.want, err)
			} else {
				assert.Equal(t, len(backup.Calls), 1)
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

	happyPathMap := make(map[string]string)
	happyPathMap["https://google.com"] = "1"

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
			fields:  fields{prefix: "http://localhost:8080", links: happyPathMap, backup: backup},
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
