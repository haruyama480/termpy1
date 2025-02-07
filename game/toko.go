package game

import (
	"context"
	"io"
	"os"
	"sync"
	"syscall"
	"time"

	"github.com/haruyama480/termpy1/ansi"
	"github.com/haruyama480/termpy1/pu2"
	"golang.org/x/term"
)

const (
	initTax  = 2
	initTay  = pu2.Y12 // 実際の回し練習ができるように
	initTdir = pu2.DirUp
)

type Toko struct {
	rec pu2.SoloRecord

	// state
	tind  int     // tsumo index
	tax   int     // tsumo ax
	tay   int     // tsumo ay
	tdir  pu2.Dir // tsumo dir
	tswap int     // 0 or 1
}

func NewToko(width int, height int, seed int64) *Toko {
	yama := pu2.NewYama(pu2.TsumoLoop, seed)
	rec := pu2.NewSoloRecord(yama)
	t := Toko{rec: rec}
	t.SetTind(0)
	return &t
}

func (g *Toko) SetTind(i int) {
	g.tind = i
	g.tax = initTax
	g.tay = initTay
	g.tdir = initTdir
	g.tswap = 0
}

type TokoConsole struct {
	Toko
	con      *ansi.Console
	oldstate *term.State
}

func NewTokoConsole() {}

func (g *TokoConsole) GetTerm() {
	tty, err := os.OpenFile("/dev/tty", syscall.O_RDWR, 0)
	if err != nil {
		panic(err)
	}
	g.con = ansi.NewConsole(tty)

	oldstate, err := term.MakeRaw(int(tty.Fd()))
	if err != nil {
		panic(err)
	}
	g.oldstate = oldstate

	g.con.HideCursor()
}

func (g *TokoConsole) RestoreTerm() {
	g.con.ShowCursor()

	err := term.Restore(int(g.con.Fd()), g.oldstate)
	if err != nil {
		panic(err)
	}
}

type Action int

const (
	ActionNothing Action = iota
	ActionQuit

	ActionPlayStart // for iteration
	ActionRight
	ActionLeft
	ActionQuickDrop
	ActionRotateR
	ActionRotateL
	ActionUndo
	ActionRedo
	ActionPlayEnd // for iteration
)

const sizePerUnit = 2

type DrawType int

const (
	DrawTypeNone DrawType = iota
	DrawTypeCell
	DrawTypeCellY13
	DrawTypeTsumoA
	DrawTypeTsumoB
	DrawTypeNextTsumo
)

var DrawMap = map[DrawType]func(*ansi.Console){
	DrawTypeNone:      func(w *ansi.Console) {},
	DrawTypeCell:      func(w *ansi.Console) { w.FontBold() },
	DrawTypeCellY13:   func(w *ansi.Console) { w.FontFaint() },
	DrawTypeTsumoA:    func(w *ansi.Console) { w.FontBold(); w.FontUnderline() },
	DrawTypeTsumoB:    func(w *ansi.Console) { w.FontBold() },
	DrawTypeNextTsumo: func(w *ansi.Console) { w.FontBold() },
}

func PrintCell(w *ansi.Console, p pu2.Cell, t DrawType) {
	colors := []ansi.FontColor{ansi.FontColorRed, ansi.FontColorGreen, ansi.FontColorYellow, ansi.FontColorBlue}

	if p == 0 {
		w.WriteString("﹍")
		return
	}
	DrawMap[t](w)
	w.FontColor(colors[p-1])
	w.WriteString("⬤")
	w.FontColorReset()
	w.WriteString(" ")
}

func (g *TokoConsole) Render() {
	f := g.rec.Field()
	w := g.con

	// field
	{
		for y := pu2.Y14; y >= 0; y-- {
			if y == pu2.Y14 {
				w.EraseEndOfLine()
				for x := 0; x < pu2.Width; x++ {
					if f[x][y] != 0 {
						w.WriteString("⬤ ")
					} else {
						w.WriteString("  ")
					}
				}
				w.NewLine()
				continue
			}
			if y == pu2.Y13 {
				w.EraseEndOfLine()
				for x := 0; x < pu2.Width; x++ {
					if f[x][y] != 0 {
						PrintCell(w, f[x][y], DrawTypeCellY13)
					} else {
						w.WriteString("  ")
					}
				}
				w.NewLine()
				continue
			}
			for x := 0; x < pu2.Width; x++ {
				PrintCell(w, f[x][y], DrawTypeCell)
			}
			w.NewLine()
		}

		w.MoveToHead()           // x=0
		w.MoveTo(-pu2.Height, 0) // y=0
	}

	// tsumo
	{
		tsumo := g.rec.Yama.Get(g.tind)

		ax := g.tax
		ay := g.tay
		bx, by := pu2.BXY(ax, ay, g.tdir)

		day := pu2.Y14 - ay
		dax := ax * sizePerUnit
		dby := pu2.Y14 - by
		dbx := bx * sizePerUnit

		w.MoveTo(day, dax)
		PrintCell(w, pu2.Cell(tsumo.A), DrawTypeTsumoA)
		w.MoveTo(dby-day, dbx-dax-sizePerUnit)
		PrintCell(w, pu2.Cell(tsumo.B), DrawTypeTsumoB)

		w.MoveToHead()    // x=0
		w.MoveTo(-dby, 0) // y=0
	}

	// next tsumos
	{
		nex1 := g.rec.Yama.Get(g.tind + 1)
		nex2 := g.rec.Yama.Get(g.tind + 2)

		w.MoveTo(0, pu2.Width*sizePerUnit+1)
		PrintCell(w, pu2.Cell(nex1.B), DrawTypeNextTsumo)
		w.MoveTo(0, 1)
		PrintCell(w, pu2.Cell(nex2.B), DrawTypeNextTsumo)
		w.MoveTo(1, -sizePerUnit*2-1)
		PrintCell(w, pu2.Cell(nex1.A), DrawTypeNextTsumo)
		w.MoveTo(0, 1)
		PrintCell(w, pu2.Cell(nex2.A), DrawTypeNextTsumo)

		w.MoveToHead()  // x=0
		w.MoveTo(-1, 0) // y=0
	}

	w.Flush()
}

func (g *TokoConsole) Play(act Action) (err error) {
	fraw := g.rec.Field()
	f := &fraw
	switch act {
	case ActionRight:
		ok, ax_ := f.MoveTsumo(g.tax, g.tay, g.tdir, 1)
		if ok {
			g.tax = ax_
		}
	case ActionLeft:
		ok, ax_ := f.MoveTsumo(g.tax, g.tay, g.tdir, -1)
		if ok {
			g.tax = ax_
		}
	case ActionRotateR:
		ok, ax_, ay_, adir_ := f.RotateTsumo(g.tax, g.tay, g.tdir, 1+g.tswap)
		if ok {
			g.tax = ax_
			g.tay = ay_
			g.tdir = adir_
			g.tswap = 0
		} else if g.tdir.IsVertical() {
			g.tswap = 1
		}
	case ActionRotateL:
		ok, ax_, ay_, adir_ := f.RotateTsumo(g.tax, g.tay, g.tdir, -1-g.tswap)
		if ok {
			g.tax = ax_
			g.tay = ay_
			g.tdir = adir_
			g.tswap = 0
		} else if g.tdir.IsVertical() {
			g.tswap = 1
		}
	case ActionUndo:
		if g.rec.Undo() {
			g.SetTind(g.tind - 1)
		}
	case ActionRedo:
		if g.rec.Redo() {
			g.SetTind(g.tind + 1)
		}
	case ActionQuickDrop:
		h := pu2.Handle{AX: g.tax, Dir: g.tdir}
		maydead := f.WillDead(h)

		g.rec.Push(h)
		g.SetTind(g.tind + 1) // A

		vanish := false
		for {
			if !g.rec.Vanish() {
				break
			}
			vanish = true
			g.rec.Fall()
		}

		if maydead && !vanish {
			// game overさせない
			g.rec.Pop()
			g.SetTind(g.tind - 1) // NOTE: Aを復元できたらいいかも
		}
	}
	return nil
}

func keyStream(ctx context.Context, in io.Reader, key chan<- []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buf := make([]byte, 1024)
			n, err := in.Read(buf)
			if err != nil {
				panic(err)
			}
			key <- buf[:n]
			time.Sleep(1000 * time.Microsecond) // 1000fps. 適当
		}
	}
}

func (g *TokoConsole) Run(seed int64) {
	g.GetTerm()
	defer g.RestoreTerm()

	g.con.WriteString("Let's play toko!\n\r")
	g.con.Flush()

	g.Toko = *NewToko(6, 13, seed)
	g.Render()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	keyst := make(chan []byte)
	go keyStream(ctx, g.con, keyst)

	ticker := time.NewTicker(33333 * time.Microsecond) // 30fps

	keystreams := make([][]byte, 0, 100)
	mu := &sync.Mutex{}

Loop:
	for {
		select {
		case data := <-keyst:
			// g.con.WriteString(fmt.Sprintf("data: %d %v\n\r", len(data), string(data)))
			mu.Lock()
			keystreams = append(keystreams, data)
			mu.Unlock()
		case <-ticker.C:
			act := ActionNothing

			mu.Lock()
			if len(keystreams) != 0 {
				// ある意味ratelimitな実装
				data := keystreams[len(keystreams)-1]
				act = b2action(data[0])
			}
			keystreams = keystreams[:0]
			mu.Unlock()

			if act == ActionQuit {
				break Loop
			}
			if act <= ActionPlayStart && act >= ActionPlayEnd {
				continue
			}

			err := g.Play(act)
			if err != nil {
				break Loop
			}

			g.Render()
		}
	}
}

func b2action(b byte) Action {
	a, ok := byte2action[b]
	if !ok {
		return ActionNothing
	}
	return a
}

var action2byte = map[Action]byte{
	ActionLeft:      'd',
	ActionQuickDrop: ' ',
	ActionRight:     'f',
	ActionRotateR:   'k',
	ActionRotateL:   'j',
	ActionUndo:      'z',
	ActionRedo:      'r',
	ActionQuit:      'q',
}
var byte2action = map[byte]Action{}

func init() {
	for k, v := range action2byte {
		byte2action[v] = k
	}
}
