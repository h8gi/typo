package main

import (
	"fmt"
	"strings"
	"time"

	runewidth "github.com/mattn/go-runewidth"
	termbox "github.com/nsf/termbox-go"
)

// show text appropriately
func DrawText(x, y int, word string, fg, bg termbox.Attribute) {
	xoff := 0
	for _, r := range word {
		termbox.SetCell(x+xoff, y, r, fg, bg)
		xoff += runewidth.RuneWidth(r)
	}
}

type TextArea struct {
	rawText        string
	words          []string
	currentWordPos int
	currentCharPos int
}

func (ta *TextArea) SetText(text string) {
	ta.rawText = text
	ta.words = strings.Fields(text)
}

func (ta *TextArea) CurrentWord() string {
	if ta.currentWordPos < len(ta.words) {
		return ta.words[ta.currentWordPos]
	}
	return ""
}

func (ta *TextArea) NextWord() (string, bool) {
	ta.currentWordPos += 1
	if ta.currentWordPos < len(ta.words) {
		return ta.words[ta.currentWordPos], true
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
	ok bool
	dr time.Duration
}

func NewTypo(text string) Typo {
	ta := TextArea{}
	ta.SetText(text)
	ia := InputArea{}
	return Typo{
		ta: &ta,
		ia: &ia,
		ok: true,
	}
}

// r is not space
func (ty *Typo) GetRune(r rune) {
	ty.ia.CurrentInput = ty.ia.CurrentInput + string(r)
	ty.ok = ty.IsMatch()
}

func (ty *Typo) BackSpace() {
	if ty.ia.CurrentInput != "" {
		ty.ia.CurrentInput = ty.ia.CurrentInput[0 : len(ty.ia.CurrentInput)-1]
	}
	ty.ok = ty.IsMatch()
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
	return (ty.ta.CurrentWord() == ty.ia.CurrentInput) && (ty.ta.currentWordPos == len(ty.ta.words)-1)
}

func (ty *Typo) DrawTextArea(x, y, width, height int) {
	xoff := 0
	yoff := 0

	for i, word := range ty.ta.words {
		colfg := termbox.ColorDefault
		colbg := termbox.ColorDefault

		if i == ty.ta.currentWordPos {
			if ty.ok {
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
	if ty.ok {
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
	var displayText string
	if len(ty.ia.CurrentInput) > width-3 {
		displayText = ty.ia.CurrentInput[len(ty.ia.CurrentInput)-width+2:]
	} else {
		displayText = ty.ia.CurrentInput
	}
	DrawText(x+2, y, displayText, colfg, termbox.ColorDefault)
	termbox.SetCursor(x+runewidth.StringWidth(displayText)+2, y)
}

func (ty *Typo) Draw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	// text area
	ty.DrawTextArea(5, 2, w-10, h-10)
	// input area
	ty.DrawInputArea((w-50)/2, h-5, 50, 3)
	DrawText(5, h-2, ty.ia.CurrentInput, termbox.ColorDefault, termbox.ColorDefault)
	DrawText(5, h-1, ty.ta.CurrentWord(), termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
}

func (ty *Typo) Start() {
	start := time.Now()
	ty.Draw()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				return
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
		case termbox.EventError:
			panic(ev.Err)
		}

		if ty.IsFinish() {
			break mainloop
		}

		ty.Draw()
	}
	ty.dr = time.Since(start)
	ty.Result()
}

func (ty *Typo) Result() {
	resultStr := fmt.Sprintf("WPM: %f", 12*float64(len(ty.ta.rawText))/ty.dr.Seconds())
	termbox.HideCursor()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	DrawText(w/2-10, h/2, resultStr,
		termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break mainloop
		default:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			w, h := termbox.Size()
			DrawText(w/2-10, h/2, resultStr,
				termbox.ColorDefault, termbox.ColorDefault)
			termbox.Flush()
		}
	}
}

// main
func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()

	var ty Typo = NewTypo("When on board H.M.S.")
	ty.Start()
}
