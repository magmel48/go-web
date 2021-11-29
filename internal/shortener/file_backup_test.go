package shortener

import (
	"github.com/magmel48/go-web/internal/shortener/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNewFileBackup(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "opens proper file by specified path",
			args: args{filePath: "some_long_path_to_file.txt"},
			want: "some_long_path_to_file.txt",
		},
	}

	mockOpenFile := &mocks.OpenFile{}
	mockOpenFile.On(
		"Execute",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("int"),
		mock.Anything).Return(nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileBackup(tt.args.filePath, mockOpenFile.Execute); !assert.Equal(t, len(mockOpenFile.Calls), 1) {
				t.Errorf("NewFileBackup() = %v, want %v", got, tt.want)
			} else {
				assert.Equal(t, mockOpenFile.Calls[0].Arguments[0], tt.want)
			}
		})
	}
}
