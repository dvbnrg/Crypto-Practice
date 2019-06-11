package main

import (
	"crypto/rsa"
	"os"
	"reflect"
	"testing"
)

func TestPublicKey(t *testing.T) {
	type args struct {
		cab []byte
	}
	tests := []struct {
		name string
		args args
		want *os.File
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PublicKey(tt.args.cab); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PublicKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrivateKey(t *testing.T) {
	type args struct {
		priv *rsa.PrivateKey
	}
	tests := []struct {
		name string
		args args
		want *os.File
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrivateKey(tt.args.priv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrivateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
