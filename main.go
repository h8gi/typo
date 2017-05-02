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

func DrawText(x, y int, word string) {
	xoff := 0
	for _, r := range word {
		termbox.SetCell(x+xoff, y, r, termbox.ColorDefault, termbox.ColorDefault)
		xoff += runewidth.RuneWidth(r)
	}
}

func (ta *TextArea) SetText(text string) {
	ta.rawText = text
	ta.words = strings.Fields(text)
}

func (ta *TextArea) CurrentWord() string {
	if ta.currentPos > len(ta.words) {
		return ""
	}
	return ta.words[ta.currentPos]
}

func (ta *TextArea) NextWord() (string, bool) {
	ta.currentPos += 1
	if ta.currentPos > len(ta.words) {
		return "", false
	}
	return ta.words[ta.currentPos], true
}

type InputArea struct {
	words        []string
	CurrentInput string
}

func (ia *InputArea) Draw(x, y, width, height int) {
	// xoff := 0
	// yoff := 0
	termbox.SetCell(x, y-1, '┌', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x, y, '│', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x, y+1, '└', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x+width, y-1, '┐', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x+width, y, '│', termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x+width, y+1, '┘', termbox.ColorDefault, termbox.ColorDefault)
	for i := 1; i < width; i++ {
		termbox.SetCell(x+i, y-1, '─', termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x+i, y+1, '─', termbox.ColorDefault, termbox.ColorDefault)
	}
	DrawText(x+2, y, ia.CurrentInput)
	ia.DrawCursor(x, y)
}

func (ia *InputArea) DrawCursor(x, y int) {
	termbox.SetCursor(x+runewidth.StringWidth(ia.CurrentInput)+2, y)
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

func (ty *Typo) DrawTextArea(x, y, width, height int) {
	xoff := 0
	yoff := 0

	for i, word := range ty.ta.words {
		colfg := termbox.ColorDefault
		colbg := termbox.ColorDefault

		if i == ty.ta.currentPos {
			colbg = termbox.ColorWhite
			if strings.HasPrefix(ty.ta.CurrentWord(), ty.ia.CurrentInput) {
				colfg = termbox.ColorGreen
			} else {
				colfg = termbox.ColorRed
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

var ty Typo = Typo{
	ia: &InputArea{},
	ta: &TextArea{},
}

func drawAll() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	// text area
	ty.DrawTextArea(5, 2, w-10, h-10)
	// input area
	ty.ia.Draw(5, h-5, w-10, 3)
	DrawText(5, h-2, ty.ia.CurrentInput)
	DrawText(5, h-1, ty.ta.CurrentWord())
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	ty.ta.SetText("When on board H.M.S. 'ビーグル,' as naturalist, I was much struck with certain facts in the distribution of the inhabitants of South America, and in the geological relations of the present to the past inhabitants of that continent. These facts seemed to me to throw some light on the origin of species--that mystery of mysteries, as it has been called by one of our greatest philosophers. On my return home, it occurred to me, in 1837, that something might perhaps be made out on this question by patiently accumulating and reflecting on all sorts of facts which could possibly have any bearing on it. After five years' work I allowed myself to speculate on the subject, and drew up some short notes; these I enlarged in 1844 into a sketch of the conclusions, which then seemed to me probable: from that period to the present day I have steadily pursued the same object. I hope that I may be excused for entering on these personal details, as I give them to show that I have not been hasty in coming to a decision.")

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
		drawAll()
	}
}
