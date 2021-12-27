package auth

import (
	"crypto/cipher"
	"github.com/google/uuid"
	"testing"
)

type TestAEAD struct{}

func (t TestAEAD) NonceSize() int {
	return 1
}

func (t TestAEAD) Overhead() int {
	return 0
}

func (t TestAEAD) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	return nil
}

func (t TestAEAD) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	return nil, nil
}

//func TestCustomAuth_Decode(t *testing.T) {
//	type fields struct {
//		algo cipher.AEAD
//	}
//	type args struct {
//		sequence []byte
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    uuid.UUID
//		wantErr bool
//	}{
//		{},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			auth := CustomAuth{
//				algo: tt.fields.algo,
//			}
//			got, err := auth.Decode(tt.args.sequence)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("Decode() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

func TestCustomAuth_Encode(t *testing.T) {
	type fields struct {
		algo cipher.AEAD
	}
	type args struct {
		id uuid.UUID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "happy path",
			fields:  fields{algo: TestAEAD{}},
			args:    args{id: uuid.New()},
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := CustomAuth{
				algo: tt.fields.algo,
			}

			got, err := auth.Encode(tt.args.id, func() ([]byte, error) { return nil, nil })
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestNewCustomAuth(t *testing.T) {
//	tests := []struct {
//		name    string
//		want    *CustomAuth
//		wantErr bool
//	}{
//		{},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := NewCustomAuth()
//			if (err != nil) != tt.wantErr {
//				t.Errorf("NewCustomAuth() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("NewCustomAuth() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
