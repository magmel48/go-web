package shortener

import (
	"reflect"
	"testing"
)

func TestNewFileBackup(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name string
		args args
		want FileBackup
	}{
		{
			name: "opens proper file by specified path",
			args: args{filePath: "some_long_path_to_file.txt"},
			want: FileBackup{file: nil}, // TODO
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFileBackup(tt.args.filePath); !reflect.DeepEqual(got.file, tt.want.file) {
				t.Errorf("NewFileBackup() = %v, want %v", got, tt.want)
			}
		})
	}
}
