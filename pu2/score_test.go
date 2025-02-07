package pu2

import (
	"testing"
)

func Test_chainBonus(t *testing.T) {
	type args struct {
		chain int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "chain 1",
			args: args{chain: 1},
			want: 0,
		},
		{
			name: "chain 7",
			args: args{chain: 7},
			want: 128,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BonusChain(tt.args.chain); got != tt.want {
				t.Errorf("BonusChain() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBonusConnection(t *testing.T) {
	type args struct {
		connect int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "connect 4",
			args: args{connect: 4},
			want: 0,
		},
		{
			name: "connect 5",
			args: args{connect: 5},
			want: 2,
		},
		{
			name: "connect 10",
			args: args{connect: 10},
			want: 7,
		},
		{
			name: "connect 11",
			args: args{connect: 11},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BonusConnection(tt.args.connect); got != tt.want {
				t.Errorf("BonusConnection() = %v, want %v", got, tt.want)
			}
		})
	}
}
