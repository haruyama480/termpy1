package pu2

import (
	"strings"
	"testing"
)

func newS2Pair(s string) Pair {
	if len(s) != 2 {
		panic("invalid PPair string")
	}
	s2r := map[rune]Color{
		'X': 1,
		'O': 2,
		'0': 3,
		'I': 4,
	}
	return Pair{
		A: s2r[rune(s[0])],
		B: s2r[rune(s[1])],
	}
}

func TestBoard_FallTsumo(t *testing.T) {
	type args struct {
		pp Pair
		h  Handle
	}
	tests := []struct {
		name  string
		field Field
		argss []args
		want  string
	}{
		{
			field: NewField(),
			argss: []args{
				{newS2Pair("XO"), Handle{0, DirRight}},
				{newS2Pair("0I"), Handle{2, DirLeft}},
				{newS2Pair("XO"), Handle{4, DirUp}},
				{newS2Pair("0I"), Handle{5, DirDown}},
			},
			want: strings.Join([]string{
				"______",
				"______",
				"______",
				"______",
				"______",
				"______",
				"______",
				"______",
				"______",
				"______",
				"______",
				"______",
				"_I__O0",
				"XO0_XI",
			}, "\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.field
			for _, args := range tt.argss {
				b.AddHandle(args.pp, args.h)
			}
			got := b.String()
			if got != tt.want {
				t.Errorf("FallTsumo Result = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestField_RotateTsumo(t *testing.T) {
	type args struct {
		ax     int
		ay     int
		dir    Dir
		incdir int
	}
	tests := []struct {
		name     string
		f        Field
		args     args
		wantOk   bool
		wantAx_  int
		wantAy_  int
		wantDir_ Dir
	}{
		{
			name: "simple rotate r",
			f:    NewField(),
			args: args{
				ax:     1,
				ay:     2,
				dir:    DirLeft,
				incdir: 1,
			},
			wantOk:   true,
			wantAx_:  1,
			wantAy_:  2,
			wantDir_: DirUp,
		},
		{
			name: "simple rotate l",
			f:    NewField(),
			args: args{
				ax:     2,
				ay:     2,
				dir:    DirUp,
				incdir: -1,
			},
			wantOk:   true,
			wantAx_:  2,
			wantAy_:  2,
			wantDir_: DirLeft,
		},
		{
			name: "simple rotate failed",
			f: func() Field {
				f := NewField()
				f[2][Y14] = 1
				return f
			}(),
			args: args{
				ax:     2,
				ay:     Y13,
				dir:    DirRight,
				incdir: -1,
			},
			wantOk: false,
		},
		{
			name: "swap",
			f:    NewField(),
			args: args{
				ax:     2,
				ay:     0,
				dir:    DirUp,
				incdir: 2,
			},
			wantOk:   true,
			wantAx_:  2,
			wantAy_:  1,
			wantDir_: DirDown,
		},
		{
			name: "lift y",
			f:    NewField(),
			args: args{
				ax:     2,
				ay:     0,
				dir:    DirRight,
				incdir: 1,
			},
			wantOk:   true,
			wantAx_:  2,
			wantAy_:  1,
			wantDir_: DirDown,
		},
		{
			name: "lift y failed",
			f: func() Field {
				f := NewField()
				f[2][Y12] = 1
				return f
			}(),
			args: args{
				ax:     2,
				ay:     Y13,
				dir:    DirRight,
				incdir: 1,
			},
			wantOk: false,
		},
		{
			name: "inc x",
			f:    NewField(),
			args: args{
				ax:     0,
				ay:     3,
				dir:    DirDown,
				incdir: 1,
			},
			wantOk:   true,
			wantAx_:  1,
			wantAy_:  3,
			wantDir_: DirLeft,
		},
		{
			name: "dec x",
			f:    NewField(),
			args: args{
				ax:     5,
				ay:     3,
				dir:    DirDown,
				incdir: -1,
			},
			wantOk:   true,
			wantAx_:  4,
			wantAy_:  3,
			wantDir_: DirRight,
		},
		{
			name: "inc x failed",
			f: func() Field {
				f := NewField()
				f[1][3] = 1
				f[3][3] = 1
				return f
			}(),
			args: args{
				ax:     2,
				ay:     3,
				dir:    DirUp,
				incdir: 1,
			},
			wantOk: false,
		},
		{
			name: "dec x failed",
			f: func() Field {
				f := NewField()
				f[1][3] = 1
				f[3][3] = 1
				return f
			}(),
			args: args{
				ax:     2,
				ay:     3,
				dir:    DirUp,
				incdir: -1,
			},
			wantOk: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotAx_, gotAy_, gotDir_ := tt.f.RotateTsumo(tt.args.ax, tt.args.ay, tt.args.dir, tt.args.incdir)
			if gotOk != tt.wantOk {
				t.Errorf("Field.RotateTsumo() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
			if gotAx_ != tt.wantAx_ {
				t.Errorf("Field.RotateTsumo() gotAx_ = %v, want %v", gotAx_, tt.wantAx_)
			}
			if gotAy_ != tt.wantAy_ {
				t.Errorf("Field.RotateTsumo() gotAy_ = %v, want %v", gotAy_, tt.wantAy_)
			}
			if gotDir_ != tt.wantDir_ {
				t.Errorf("Field.RotateTsumo() gotDir_ = %v, want %v", gotDir_, tt.wantDir_)
			}
		})
	}
}

func TestField_Fall(t *testing.T) {
	tests := []struct {
		name        string
		f           Field
		want        Field
		wantChanged bool
	}{
		{
			name: "no fall",
			f: func() Field {
				f := NewField()
				f[2][0] = 1
				f[3][Y14] = 1
				return f
			}(),
			want: func() Field {
				f := NewField()
				f[2][0] = 1
				f[3][Y14] = 1
				return f
			}(),
		},
		{
			name: "fall",
			f: func() Field {
				f := NewField()
				f[2][3] = 3
				f[2][5] = 2
				f[2][6] = 3
				f[2][Y13] = 1
				return f
			}(),
			want: func() Field {
				f := NewField()
				f[2][0] = 3
				f[2][1] = 2
				f[2][2] = 3
				f[2][3] = 1
				return f
			}(),
			wantChanged: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotChanged := tt.f.Fall(); gotChanged != tt.wantChanged {
				t.Errorf("Field.Fall() = %v, want %v", gotChanged, tt.wantChanged)
			}
		})
	}
}
