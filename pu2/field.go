package pu2

import (
	"strings"
)

type Field [Width][Height]Cell

func NewField() Field {
	return Field{}
}

func (f Field) String() string {
	var sb strings.Builder
	for y := Height - 1; y >= 0; y-- {
		for x := 0; x < Width; x++ {
			p := f[x][y]
			sb.WriteString(p.String())
		}
		if y != 0 {
			sb.WriteRune('\n')
		}
	}
	return sb.String()
}

func (f Field) IsDead() bool {
	return f[2][Y12] != 0
}

// ツモをその列に配置できるかどうか
// incxは-1, 0, 1のみ許容
// tokoでは軸ぷよが13段目にあるように見え、12段ある列を壁越えできそうに見えるが実際に壁越えも判定する
// 14段目も考慮にいれる
func (f Field) MoveTsumo(ax, ay int, dir Dir, incx int) (ok bool, ax_ int) {
	if incx > 1 || incx < -1 {
		panic("invalid incx")
	}
	if incx == 0 {
		return true, ax
	}

	ax_ = ax + incx
	bx_, by_ := BXY(ax_, ay, dir)
	if ax_ < 0 || ax_ >= Width || bx_ < 0 || bx_ >= Width {
		return false, 0
	}
	if f.columnHeight(ax_) > ay || f.columnHeight(bx_) > by_ {
		return false, 0
	}
	return true, ax_
}

// RotateTsumo ローテート後の位置を返す
// incdirは通常-1/1だが、swapのために-2/2も許容する。0は何もしない
// 回転可能な場合、trueと新しいax, ay, dirを返す。ayはliftされるかもしれないしswapで下がる可能性もある
// 回転不可能の場合、falseと0,0,0を返す
func (f Field) RotateTsumo(ax int, ay int, dir Dir, incdir int) (ok bool, ax_, ay_ int, dir_ Dir) {
	if incdir == 0 {
		return true, ax, ay, dir
	}

	dir_ = dir.Rotete(incdir)
	bx, by := BXY(ax, ay, dir_)

	if incdir == 2 || incdir == -2 {
		// swap
		if dir == DirUp {
			// ay_ := ay + 1
			// if ay <= Y12 {
			// 	ay++ // axisが上にくる場合半マス浮く
			// }
			return true, ax, ay + 1, dir_
		} else if dir == DirDown {
			return true, ax, ay - 1, dir_
		} else {
			panic("incdir==2 must occur only in DirUp or DirDown")
		}
	}

	// indir := 1 or -1
	switch dir_ {
	case DirUp:
		if f[bx][by] == 0 { // axisは14段目にはならないため、byのチェックは不要
			return true, ax, ay, dir_
		}
		return false, 0, 0, 0
	case DirDown:
		if by >= 0 && f[bx][by] == 0 {
			return true, ax, ay, dir_
		}
		// lift ay
		if ay >= Y14 {
			panic("invalid ay in RotateTsumo")
		}
		if ay == Y13 { // ay must be less than 14
			return false, 0, 0, 0
		}
		return true, ax, ay + 1, dir_
	case DirLeft:
		if bx >= 0 && f[bx][by] == 0 {
			return true, ax, ay, dir_
		}
		// move right
		if ax+1 == Width || f[ax+1][ay] != 0 {
			return false, 0, 0, 0
		}
		return true, ax + 1, ay, dir_
	case DirRight:
		if bx < Width && f[bx][by] == 0 {
			return true, ax, ay, dir_
		}
		// move left
		if ax-1 == -1 || f[ax-1][ay] != 0 {
			return false, 0, 0, 0
		}
		return true, ax - 1, ay, dir_
	}
	panic("invalid")
}

// 設置した場合の位置を返す
func (f Field) GhostTsumo(h Handle) (ay, ax, by, bx int) {
	x := h.AX
	dir := h.Dir
	ax = x
	ay = f.columnHeight(x)
	switch dir {
	case DirUp:
		bx = ax
		by = ay + 1
	case DirDown:
		bx = ax
		by = ay
		ay++
	case DirLeft:
		bx = ax - 1
		by = f.columnHeight(bx)
	case DirRight:
		bx = ax + 1
		by = f.columnHeight(bx)
	}
	return
}

// 設置可能な高さを返す。列に存在する数に一致する
func (b Field) columnHeight(col int) int {
	if col < 0 || col >= Width {
		panic("invalid col")
	}
	for y := 0; y < Height; y++ {
		if b[col][y] == 0 {
			return y
		}
	}
	return Height
}

// 設置可能か判定
func (f Field) WillDead(h Handle) bool {
	ay, ax, by, bx := f.GhostTsumo(h)
	if ax == 2 && ay == Y12 || bx == 2 && by == Y12 {
		return true
	}
	return false
}

// 設置する。設置可能かは判定しない
func (f *Field) AddHandle(pp Pair, h Handle) {
	x := h.AX
	if x < 0 || x >= Width {
		panic("invalid x")
	}
	ay, ax, by, bx := f.GhostTsumo(h)
	f[ax][ay] = Cell(pp.A)
	f[bx][by] = Cell(pp.B)
}

func (f *Field) Fall() (changed bool) {
	for col := range Width {
		cur := 0
		for y := range H13 { // 14段目は無視
			if f[col][y] == CellNone {
				continue
			}
			if cur != y {
				changed = true
				f[col][cur] = f[col][y]
			}
			cur++
		}
		for cur < H13 {
			f[col][cur] = CellNone
			cur++
		}
	}
	return
}

func (f *Field) Vanish() (changed bool, bonusConColor int) {
	visited := [Width][Height]byte{}
	usedColor := [ColorSize]byte{}

	for x := range Width {
		for y := range H13 {
			if !f[x][y].IsPlain() {
				continue
			}
			if visited[x][y] != 0 {
				continue
			}
			connection := 0
			targetColor := f[x][y]
			group := [][2]int{{x, y}}
			for i := 0; i < len(group); i++ {
				posi := group[i]
				x := posi[0]
				y := posi[1]
				connection++
				visited[x][y] = 1

				if x > 0 && f[x-1][y] == targetColor && visited[x-1][y] == 0 {
					group = append(group, [2]int{x - 1, y})
				}
				if x < Width-1 && f[x+1][y] == targetColor && visited[x+1][y] == 0 {
					group = append(group, [2]int{x + 1, y})
				}
				if y > 0 && f[x][y-1] == targetColor && visited[x][y-1] == 0 {
					group = append(group, [2]int{x, y - 1})
				}
				if y < H13-1 && f[x][y+1] == targetColor && visited[x][y+1] == 0 {
					group = append(group, [2]int{x, y + 1})
				}
			}
			if connection < ConnectionThreshold {
				continue
			}
			usedColor[targetColor-1] = 1
			changed = true
			bonusConColor += BonusConnection(connection)

			for _, posi := range group {
				x := posi[0]
				y := posi[1]
				f[x][y] = CellNone // vanish
			}
		}
	}
	colorNum := 0
	for i := range ColorSize {
		if usedColor[i] != 0 {
			colorNum++
		}
	}
	bonusConColor += BonusColor(colorNum)
	return
}
