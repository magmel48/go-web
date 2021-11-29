package shortener

import (
	"github.com/magmel48/go-web/internal/shortener/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestShortener_MakeShorter(t *testing.T) {
	type fields struct {
		prefix string
		links  map[string]string
	}
	type args struct {
		url string
	}

	MockBackup := &mocks.Backup{}
	MockBackup.On("Append", mock.AnythingOfType("string")).Return()

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "happy path",
			fields: fields{prefix: "http://localhost:8080", links: make(map[string]string)},
			args:   args{url: "https://google.com"},
			want:   "http://localhost:8080/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				prefix: tt.fields.prefix,
				links:  tt.fields.links,
				backup: MockBackup,
			}

			if got, err := s.MakeShorter(tt.args.url); got != tt.want || err != nil {
				t.Errorf("MakeShorter() = %v, want %v, err %v", got, tt.want, err)
			}
		})
	}
}

func TestShortener_RestoreLong(t *testing.T) {
	type fields struct {
		prefix string
		links  map[string]string
	}
	type args struct {
		id string
	}

	happyPathMap := make(map[string]string)
	happyPathMap["https://google.com"] = "1"

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "happy path",
			fields:  fields{prefix: "http://localhost:8080", links: happyPathMap},
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
