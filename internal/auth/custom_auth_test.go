package auth

import (
	"crypto/cipher"
	"reflect"
	"testing"
)

// TestAEAD represents algorithm that does not perform a mutation on provided bytes for encoding or decoding.
// Created only for testing purposes.
type TestAEAD struct{}

func (t TestAEAD) NonceSize() int {
	return 1
}

func (t TestAEAD) Overhead() int {
	return 0
}

func (t TestAEAD) Seal(dst, _, plaintext, _ []byte) []byte {
	dst = append(dst, plaintext...)

	return dst
}

func (t TestAEAD) Open(dst, _, ciphertext, _ []byte) ([]byte, error) {
	dst = append(dst, ciphertext...)

	return dst, nil
}

func TestCustomAuth_Decode(t *testing.T) {
	type fields struct {
		algo cipher.AEAD
	}
	type args struct {
		sequence []byte
	}

	id := "26d1ac21-57d5-43ba-b2f7-08d36310aa07"

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    UserID
		wantErr bool
	}{
		{
			name:   "happy path",
			fields: fields{algo: TestAEAD{}},
			args: args{sequence: []byte{
				77, 106, 90, 107, 77, 87, 70, 106, 77, 106, 69, 116, 78, 84, 100, 107, 78, 83, 48, 48, 77, 50, 74, 104,
				76, 87, 73, 121, 90, 106, 99, 116, 77, 68, 104, 107, 77, 122, 89, 122, 77, 84, 66, 104, 89, 84, 65, 51,
				65, 81}},
			want:    &id,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := CustomAuth{
				algo: tt.fields.algo,
			}
			got, err := auth.Decode(tt.args.sequence)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCustomAuth_Encode(t *testing.T) {
	type fields struct {
		algo      cipher.AEAD
		NonceFunc NonceFunc
	}
	type args struct {
		id UserID
	}

	id := "26d1ac21-57d5-43ba-b2f7-08d36310aa07"

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name:   "happy path",
			fields: fields{algo: TestAEAD{}, NonceFunc: func(_ int) ([]byte, error) { return []byte{1}, nil }},
			args:   args{id: &id},
			want: []byte{
				77, 106, 90, 107, 77, 87, 70, 106, 77, 106, 69, 116, 78, 84, 100, 107, 78, 83, 48, 48, 77, 50, 74, 104,
				76, 87, 73, 121, 90, 106, 99, 116, 77, 68, 104, 107, 77, 122, 89, 122, 77, 84, 66, 104, 89, 84, 65, 51,
				65, 81},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			auth := CustomAuth{
				algo:      tt.fields.algo,
				NonceFunc: tt.fields.NonceFunc,
			}

			got, err := auth.Encode(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCustomAuth_Decode(b *testing.B) {
	auth, _ := NewCustomAuth()
	sequence := []byte{
		77, 106, 90, 107, 77, 87, 70, 106, 77, 106, 69, 116, 78, 84, 100, 107, 78, 83, 48, 48, 77, 50, 74, 104,
		76, 87, 73, 121, 90, 106, 99, 116, 77, 68, 104, 107, 77, 122, 89, 122, 77, 84, 66, 104, 89, 84, 65, 51,
		65, 81}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auth.Decode(sequence)
	}
}

func BenchmarkCustomAuth_Encode(b *testing.B) {
	auth, _ := NewCustomAuth()
	userID := "26d1ac21-57d5-43ba-b2f7-08d36310aa07"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auth.Encode(&userID)
	}
}
