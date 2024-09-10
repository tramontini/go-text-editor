package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
)

var (
	CursorX, CursorY int
	Content          [][]rune
	insertMode       bool
)

func updateScreen(screen tcell.Screen, style tcell.Style) {
	for y, line := range Content {
		for x, char := range line {
			screen.SetContent(x, y, char, nil, style)
		}
	}

	screen.ShowCursor(CursorX, CursorY)
	screen.Show()
}

func saveFile(path string) error {
	file, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, lines := range Content {
		var lineText string
		for _, runeValue := range lines {
			lineText += string(runeValue)
		}
		fmt.Fprintln(writer, lineText)
	}

	return writer.Flush()

}

func runEditor() {
	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("Erro ao iniciar a tela: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Erro ao inicializar a tela: %v", err)
	}
	defer screen.Fini()

	getFileValues("exemplo.txt")

	_, height := screen.Size()
	style := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	screen.Clear()

	setFooter(height, screen, style)

	CursorX, CursorY = 0, 0

	updateScreen(screen, style)

	for {
		event := screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyUp:
				if CursorY > 0 {
					CursorY -= 1
				}
			case tcell.KeyDown:
				if CursorY < len(Content)-1 {
					CursorY += 1
				}
			case tcell.KeyLeft:
				if CursorX > 0 {
					CursorX -= 1
				}
			case tcell.KeyRight:
				// Impedindo de ir além do que tem no vetor atual
				if CursorX < len(Content[CursorY]) {
					CursorX += 1
				}
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(Content) == 0 {
					return
				}

				Content[CursorY] = append(Content[CursorY][:CursorX-1], Content[CursorY][CursorX:]...)
			case tcell.KeyCtrlC:
				return
			case tcell.KeyCtrlI:
				insertMode = true
				setFooter(height, screen, style)
			case tcell.KeyCtrlS:
				err := saveFile("exemplo.txt")

				if err != nil {
					log.Fatal(err)
				}
			case tcell.KeyEnter:
				Content = append(Content, []rune{})
				CursorY += 1
				CursorX = 0
			default:
				if event.Rune() != 0 && insertMode {
					if len(Content) == 0 {
						Content = append(Content, []rune{})
					}

					line := Content[CursorY]
					// Verifica se o cursor vai estourar a posição dos caracteres
					if CursorX < len(line) {
						// Apendando caracter digitado no meio da linha a partir do que tudo que tinha antes e appendando o restante depois
						// line[:CursorX] -> Valores na linha até onde estava o cursor
						// []rune{event.Rune()} -> Colocando no vertor de rune o evento capturado
						//  line[CursorX:] -> Adicionando o restante dos caracteres
						// ... -> Nos appends significa pra adicionar tudo do slice 2 no slice 1
						Content[CursorY] = append(line[:CursorX], append([]rune{event.Rune()}, line[CursorX:]...)...)
					} else {
						Content[CursorY] = append(line, event.Rune())
					}

					CursorX++
				}
			}
		}
		updateScreen(screen, style)
	}

}

func getFileValues(path string) {
	file, err := os.Open(path)

	if err != nil {
		log.Panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scanText := scanner.Text()
		Content = append(Content, []rune(scanText))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func setFooter(height int, screen tcell.Screen, style tcell.Style) {

	var footerText string
	if insertMode {
		footerText = "GO Editor - Insert Mode"
	} else {
		footerText = "GO Editor - Read Mode"
	}

	footerX := 1
	footerY := height - 1

	for i, ch := range footerText {
		screen.SetContent(footerX+i, footerY, ch, nil, style)
	}
}

func main() {
	runEditor()
}
