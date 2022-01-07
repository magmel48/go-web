package shortener

import (
	"context"
	"github.com/magmel48/go-web/internal/auth"
	"testing"
)

func TestShortener_MakeShorter(t *testing.T) {
	type fields struct {
		prefix string
	}
	type args struct {
		url string
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "happy path",
			fields: fields{
				prefix: "http://localhost:8080",
			},
			args: args{url: "https://google.com"},
			want: "http://localhost:8080/1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shortener{
				prefix: tt.fields.prefix,
			}

			if got, _, err := s.MakeShorter(context.Background(), tt.args.url, auth.NewUserID()); got != tt.want || err != nil {
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
				links: map[string]string{
					"https://google.com": "1",
				},
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
