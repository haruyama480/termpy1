// not support: windows, concurrency
// https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797

package ansi

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type Console struct {
	w   *os.File
	buf strings.Builder
}

func NewConsole(tty *os.File) *Console {
	return &Console{
		w:   tty,
		buf: strings.Builder{},
	}
}

func (c *Console) Fd() uintptr {
	return c.w.Fd()
}

func (c *Console) Read(b []byte) (int, error) {
	return c.w.Read(b)
}

func (c *Console) Write(b []byte) {
	_, err := c.buf.Write(b)
	if err != nil {
		panic(err)
	}
}

func (c *Console) WriteString(s string) {
	_, err := c.buf.WriteString(s)
	if err != nil {
		panic(err)
	}
}

func (c *Console) Flush() {
	buf := c.buf.String()
	_, err := io.Copy(c.w, strings.NewReader(buf))
	if err != nil {
		panic(err)
	}
	c.buf.Reset()
}

// // Note: Some sequences, like saving and restoring cursors, are private sequences and
// // are not standardized. While some terminal emulators (i.e. xterm and derived) support
// // both SCO and DEC sequences, they are likely to have different functionality.
// // It is therefore recommended to use DEC sequences.
// func (c *AnsiConsole) SaveCursor() {
// 	c.WriteString("\033 7") // DEC
// }

// func (c *AnsiConsole) LoadCursor() {
// 	c.WriteString("\033 8") // DEC
// }

func (c *Console) EraseEndOfLine() {
	c.WriteString("\033[0K")
}

func (c *Console) NewLine() {
	c.WriteString("\r\n")
}

func (c *Console) MoveToHead() {
	c.WriteString("\r")
}

func (c *Console) MoveTo(y int, x int) {
	if y < 0 {
		c.WriteString(fmt.Sprintf("\033[%dA", -y))
	} else if y > 0 {
		c.WriteString(fmt.Sprintf("\033[%dB", y))
	}
	if x > 0 {
		c.WriteString(fmt.Sprintf("\033[%dC", x))
	} else if x < 0 {
		c.WriteString(fmt.Sprintf("\033[%dD", -x))
	}
}

type FontColor int

const (
	FontColorReset FontColor = 0
	FontColorBlack FontColor = 30 + iota
	FontColorRed
	FontColorGreen
	FontColorYellow
	FontColorBlue
	FontColorPurple
	FontColorCyan
	FontColorWhite
)

func (c *Console) FontBold() {
	c.WriteString("\033[1m")
}

func (c *Console) FontFaint() {
	c.WriteString("\033[2m")
}

func (c *Console) FontItalic() {
	c.WriteString("\033[3m")
}

func (c *Console) FontUnderline() {
	c.WriteString("\033[4m")
}

func (c *Console) FontBlinking() {
	c.WriteString("\033[5m")
}

func (c *Console) FontInverse() {
	c.WriteString("\033[7m")
}

func (c *Console) FontStrikethrough() {
	c.WriteString("\033[9m")
}

func (c *Console) FontColor(color FontColor) {
	c.WriteString(fmt.Sprintf("\033[%dm", color))
}

func (c *Console) FontColorReset() {
	c.FontColor(FontColorReset)
}

func (c *Console) HideCursor() {
	c.WriteString("\033[?25l")
}

func (c *Console) ShowCursor() {
	c.WriteString("\033[?25h")
}
