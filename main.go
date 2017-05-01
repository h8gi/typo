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

func (ta *TextArea) SetText(text string) {
	ta.rawText = text
	ta.words = strings.Fields(text)
}

func (ta *TextArea) Draw(x, y, width, height int) {
	xoff := 0
	yoff := 0
	for _, word := range ta.words {
		// word len
		if x+xoff+runewidth.StringWidth(word)+1 > width {
			xoff = 0
			yoff += 1
		}
		// draw word
		for _, r := range word {
			termbox.SetCell(x+xoff, y+yoff, r, termbox.ColorDefault, termbox.ColorDefault)
			xoff += runewidth.RuneWidth(r)
		}
		termbox.SetCell(x+xoff, y+yoff, ' ', termbox.ColorDefault, termbox.ColorDefault)
		xoff += 1
	}
}

type InputArea struct {
	words       []string
	currentWord string
}

func (ia *InputArea) GetRune(r rune) {

}

var ta TextArea

func drawAll() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	w, h := termbox.Size()
	ta.Draw(5, 2, w-5, h-10)
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}

	defer termbox.Close()
	ta.SetText("When on board H.M.S. 'ビーグル,' as naturalist, I was much struck with certain facts in the distribution of the inhabitants of South America, and in the geological relations of the present to the past inhabitants of that continent. These facts seemed to me to throw some light on the origin of species—that mystery of mysteries, as it has been called by one of our greatest philosophers. On my return home, it occurred to me, in 1837, that something might perhaps be made out on this question by patiently accumulating and reflecting on all sorts of facts which could possibly have any bearing on it. After five years' work I allowed myself to speculate on the subject, and drew up some short notes; these I enlarged in 1844 into a sketch of the conclusions, which then seemed to me probable: from that period to the present day I have steadily pursued the same object. I hope that I may be excused for entering on these personal details, as I give them to show that I have not been hasty in coming to a decision.")

	drawAll()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			}
		case termbox.EventResize:
			drawAll()
		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
