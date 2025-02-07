// pu2 is a package for game logic, but it doesn't include UI.

package pu2

import (
	"math/rand"
	"strings"
	"time"
)

const (
	Width  = 6
	Height = 14

	// マジックナンバー回避目的
	H14 = 14 // 1-indexな14段目。高さを表現する
	H13 = 13
	H12 = 12

	Y14 = 13 // 0-indexな14段目。インデックスを表現する
	Y13 = 12
	Y12 = 11

	ColorSize           = 4
	ConnectionThreshold = 4

	TsumoLoop = 128
)

// Cell := None | Pu | Ojama | Filled
// Plain := 1-4
// None := 0
// Ojama := 11
// Filled := 22 for 14th rows

type Cell int

const (
	CellNone Cell = 0

	CellPuStart Cell = 1
	CellPuEnd   Cell = 1 + ColorSize

	CellOjama  Cell = 11
	CellFilled Cell = 22 // for 14th rows
)

func (u Cell) IsPlain() bool {
	return CellPuStart <= u && u < CellPuEnd
}

func (u Cell) IsOjama() bool {
	return u == 9
}

func (u Cell) String() string {
	if u == 0 {
		return "_"
	}
	if u.IsPlain() {
		return Color(u).String()
	}
	return "*"
}

// Color is subject to CellPuStart <= Color < CellPuEnd
type Color int

func (c Color) ToCell() Cell {
	return Cell(c)
}

func (c Color) String() string {
	PRune := []rune{'X', 'O', '0', 'I'}
	// PRune := []rune{'●', '●', '●', '●'}
	// PRune := []rune{'⬤', '⬤', '⬤', '⬤'}
	return string(PRune[c-1])
}

func NewColorFromRune(r rune) Color {
	s2r := map[rune]Color{
		'X': 1,
		'O': 2,
		'0': 3,
		'I': 4,
	}
	return s2r[r]
}

type Pair struct {
	A Color // A is for Axis. int codes color
	B Color
}

func (p Pair) String() string {
	PRune := []rune{'X', 'O', '0', 'I'}
	return string(PRune[p.A-1]) + string(PRune[p.B-1])
}

type Dir int

const (
	DirUp Dir = iota
	DirRight
	DirDown
	DirLeft
)

// -2 <= i <= 2
func (d Dir) Rotete(i int) Dir {
	return (d + Dir(i) + 4) % 4
}

func (d Dir) Inc() Dir {
	return (d + 1) % 4
}

func (d Dir) Dec() Dir {
	return (d + 3) % 4
}

func (d Dir) IsVertical() bool {
	return d&1 == 0
}

func (d Dir) IsHorizontal() bool {
	return d&1 == 0
}

// bx, by may be out of range
func BXY(ax int, ay int, dir Dir) (bx int, by int) {
	switch dir {
	case DirUp:
		return ax, ay + 1
	case DirRight:
		return ax + 1, ay
	case DirDown:
		return ax, ay - 1
	case DirLeft:
		return ax - 1, ay
	}
	panic("invalid Dir")
}

type Handle struct {
	AX  int
	Dir Dir
}

func NewHandle(ax int, dir Dir) Handle {
	return Handle{AX: ax, Dir: dir}
}

func (h Handle) Valid() bool {
	return (h.Dir.IsVertical() && 0 <= h.AX && h.AX < Width) ||
		(h.Dir == DirRight && 0 <= h.AX && h.AX < Width-1) ||
		(h.Dir == DirLeft && 1 <= h.AX && h.AX < Width)
}

type Yama struct {
	ps []Pair
}

func NewYama(n int, seed int64) Yama {
	var r *rand.Rand
	if seed != 0 {
		r = rand.New(rand.NewSource(seed))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	ps := make([]Pair, n)
	for i := range ps {
		var a, b int
		if i < 2 { // 初手2手3色のツモ補正
			a = 1 + r.Intn(3)
			b = 1 + r.Intn(3)
		} else {
			a = 1 + r.Intn(ColorSize)
			b = 1 + r.Intn(ColorSize)
		}
		ps[i] = Pair{A: Color(a), B: Color(b)}
	}
	return Yama{
		ps: ps,
	}
}

func (y *Yama) Get(i int) Pair {
	if i < 0 {
		panic("invalid Yama.Get()")
	}
	n := len(y.ps)
	if i >= n {
		i = i % n
	}
	return y.ps[i]
}

func (y *Yama) String() string {
	var sb strings.Builder
	for i, p := range y.ps {
		sb.WriteString(p.String())
		if i != len(y.ps)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
