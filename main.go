package main

import (
	"strings"

	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

type TextArea struct {
	rawText    string
	words      []string
	currentPos int
}

func DrawText(x, y int, word string, fg, bg termbox.Attribute) {
	xoff := 0
	for _, r := range word {
		termbox.SetCell(x+xoff, y, r, fg, bg)
		xoff += runewidth.RuneWidth(r)
	}
}

func (ta *TextArea) SetText(text string) {
	ta.rawText = text
	ta.words = strings.Fields(text)
}

func (ta *TextArea) CurrentWord() string {
	if ta.currentPos < len(ta.words) {
		return ta.words[ta.currentPos]
	}
	return ""
}

func (ta *TextArea) NextWord() (string, bool) {
	ta.currentPos += 1
	if ta.currentPos < len(ta.words) {
		return ta.words[ta.currentPos], true
	}
	return "", false
}

type InputArea struct {
	words        []string
	CurrentInput string
}

type Typo struct {
	ia *InputArea
	ta *TextArea
}

func NewTypo(text string) Typo {
	ta := TextArea{}
	ta.SetText(text)
	ia := InputArea{}
	return Typo{
		ta: &ta,
		ia: &ia,
	}
}

// r is not space
func (ty *Typo) GetRune(r rune) {
	ty.ia.CurrentInput = ty.ia.CurrentInput + string(r)
}

func (ty *Typo) BackSpace() {
	if ty.ia.CurrentInput != "" {
		ty.ia.CurrentInput = ty.ia.CurrentInput[0 : len(ty.ia.CurrentInput)-1]
	}
}

func (ty *Typo) NextWord() (string, bool) {
	word, ok := ty.ta.NextWord()
	ty.ia.CurrentInput = ""
	return word, ok
}

func (ty *Typo) IsMatch() bool {
	return strings.HasPrefix(ty.ta.CurrentWord(), ty.ia.CurrentInput)
}

func (ty *Typo) IsFinish() bool {
	return (ty.ta.CurrentWord() == ty.ia.CurrentInput) && (ty.ta.currentPos == len(ty.ta.words)-1)
}

func (ty *Typo) DrawTextArea(x, y, width, height int) {
	xoff := 0
	yoff := 0

	for i, word := range ty.ta.words {
		colfg := termbox.ColorDefault
		colbg := termbox.ColorDefault

		if i == ty.ta.currentPos {
			if ty.IsMatch() {
				colfg = termbox.ColorGreen | termbox.AttrUnderline
			} else {
				colfg = termbox.ColorRed | termbox.AttrUnderline
			}
		}
		// word len
		if xoff+runewidth.StringWidth(word)+1 > width {
			xoff = 0
			yoff += 1
		}
		// draw word
		for _, r := range word {
			termbox.SetCell(x+xoff, y+yoff, r, colfg, colbg)
			xoff += runewidth.RuneWidth(r)
		}
		termbox.SetCell(x+xoff, y+yoff, ' ', termbox.ColorDefault, termbox.ColorDefault)
		xoff += 1
	}
}

func (ty *Typo) DrawInputArea(x, y, width, height int) {
	// xoff := 0
	// yoff := 0
	colfg := termbox.ColorDefault
	colbg := termbox.ColorDefault
	if ty.IsMatch() {
		colfg = termbox.ColorGreen
	} else {
		colfg = termbox.ColorRed
	}
	termbox.SetCell(x, y-1, '┌', colfg, colbg)
	termbox.SetCell(x, y, '│', colfg, colbg)
	termbox.SetCell(x, y+1, '└', colfg, colbg)
	termbox.SetCell(x+width, y-1, '┐', colfg, colbg)
	termbox.SetCell(x+width, y, '│', colfg, colbg)
	termbox.SetCell(x+width, y+1, '┘', colfg, colbg)
	for i := 1; i < width; i++ {
		termbox.SetCell(x+i, y-1, '─', colfg, colbg)
		termbox.SetCell(x+i, y+1, '─', colfg, colbg)
	}
	DrawText(x+2, y, ty.ia.CurrentInput, colfg, termbox.ColorDefault)
	ty.ia.DrawCursor(x, y)
}

func (ia *InputArea) DrawCursor(x, y int) {
	termbox.SetCursor(x+runewidth.StringWidth(ia.CurrentInput)+2, y)
}

// main
var ty Typo = NewTypo("When on board H.M.S.")

func drawAll() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	// text area
	ty.DrawTextArea(5, 2, w-10, h-10)
	// input area
	ty.DrawInputArea(5, h-5, w-10, 3)
	DrawText(5, h-2, ty.ia.CurrentInput, termbox.ColorDefault, termbox.ColorDefault)
	DrawText(5, h-1, ty.ta.CurrentWord(), termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	drawAll()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeySpace:
				if ty.ia.CurrentInput == ty.ta.CurrentWord() {
					ty.NextWord()
				} else {
					ty.GetRune(' ')
				}
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				ty.BackSpace()
			default:
				if ev.Ch != 0 {
					ty.GetRune(ev.Ch)
				}

			}
		case termbox.EventResize:
		case termbox.EventError:
			panic(ev.Err)
		}
		if ty.IsFinish() {
			break
		}
		drawAll()
	}
}
