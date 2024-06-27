package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Create and return a slice of length number, with all elements set to value.
func repeated(value string, number int) []string {
    arr := make([]string, number)
    for i := 0; i < number; i++ {
        arr[i] = value
    }
    return arr
}

// Create a slice of length 2, with each element being another slice. The combined grid, creates a 
// visual representation of the disk (represented by width) and is returned.
func buildDisk(width, maxWidth int) [][]string {
    renderWidth := width * 4 + 5

    var disk [][]string

    var line1 = []string{}
    line1 = append(line1, repeated(" ", ((maxWidth - width) * 2))...)
    line1 = append(line1, "┌")
    line1 = append(line1, (repeated("─", renderWidth))...)
    line1 = append(line1, ("┐"))
    line1 = append(line1, repeated(" ", ((maxWidth - width) * 2))...)

    var line2 = []string{}
    line2 = append(line2, repeated(" ", ((maxWidth - width) * 2))...)
    line2 = append(line2, "│")
    line2 = append(line2, repeated(" ", width * 2 + 2)...)
    line2 = append(line2, strconv.Itoa(width))
    if width > 9 {
        line2 = append(line2, repeated(" ", width * 2 + 1)...)
    } else {
        line2 = append(line2, repeated(" ", width * 2 + 2)...)
    }
    line2 = append(line2, "│")
    if width > 9 {
        line2 = append(line2, " ")
    }

    line2 = append(line2, repeated(" ", ((maxWidth - width) * 2))...)

    disk = append(disk, line1)
    disk = append(disk, line2)
    
    return disk
}

// rodDef is an array of disks. These are then rendered, to create a visual representation of the stack
// of disks, in combination with the rod sticking out the top of the stack and return the slice.
func buildRod(rodDef []int, maxWidth int) [][]string {
    var rodImage = [][]string{}

    // Rod length is 4 higher than a full stacked rod
    var line = []string{}
    line = append(line, repeated(" ", maxWidth * 2 + 2)...)
    line = append(line, "┌")
    line = append(line, "─")
    line = append(line, "┐")
    line = append(line, repeated(" ", maxWidth * 2 + 2)...)
    rodImage = append(rodImage, line)
    for i := 0; i < 3 + 2 * (maxWidth - len(rodDef)); i++ {
        var line = []string{}
        line = append(line, repeated(" ", maxWidth * 2 + 2)...)
        line = append(line, "│")
        line = append(line, " ")
        line = append(line, "│")
        line = append(line, repeated(" ", maxWidth * 2 + 2)...)
        rodImage = append(rodImage, line)
    }

    for _, disk := range rodDef {
        rodImage = append(rodImage, buildDisk(disk, maxWidth)...)
    }

    return rodImage

}

// Create the visual representations of each rod and print them to the terminal
func printBoard(screen tcell.Screen, style tcell.Style, left, middle, right [][]string, numLinesToRender, numDisks, fromSelection int) {
    offset := 6
    for i := 0; i < numLinesToRender; i++ {
        draw(screen, 0, offset + i, style, fromSelection == 1, strings.Join(left[i], ""))
        draw(screen, len(left[i]), offset + i, style, fromSelection == 2, strings.Join(middle[i], ""))
        draw(screen, len(left[i]) + len(right[i]), offset + i, style, fromSelection == 3, strings.Join(right[i], ""))
    }
    draw(screen, 0, numLinesToRender + offset, style, false, strings.Repeat("─", 3 * (numDisks * 4 + 7)))
}

// Draw the given text beginning at the location x, y. If selected is true, we change the styling to purple
// so that we can represent the selected rod.
func draw(screen tcell.Screen, x, y int, style tcell.Style, selected bool, text string) {
    col := x
    row := y

    if selected {
        style = tcell.StyleDefault.Foreground(tcell.ColorPurple)
    }

    for _, r := range []rune(text) {
        screen.SetContent(col, row, r, nil, style)
        col++
    }
}


func main() {
    f, err := os.OpenFile("hanoi.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()

    log.SetOutput(f)

    screen, err := tcell.NewScreen()

    defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
    boxStyle := tcell.StyleDefault.Foreground(tcell.ColorReset).Background(tcell.ColorReset)

    if err != nil {
        log.Fatalf("%+v", err)
    }

    if err := screen.Init(); err != nil {
        log.Fatalf("%+v", err)
    }

    screen.SetStyle(defStyle)
    screen.Clear()

    quit := func() {
		// You have to catch panics in a defer, clean up, and
		// re-raise them - otherwise your application can
		// die without leaving any diagnostic trace.
		maybePanic := recover()
		screen.Fini()
		if maybePanic != nil {
			panic(maybePanic)
		}
	}
	defer quit()
    fromSelection := true
    var from int
    var to int

    numDisks := 5
    numLinesToRender := 2 * numDisks + 4
    
    leftDisks := []int{}
    for i := 0; i < numDisks; i++ {
        leftDisks = append(leftDisks, i + 1)
    }
    middleDisks := []int{}
    rightDisks := []int{}

    left := buildRod(leftDisks, numDisks)
    middle := buildRod(middleDisks, numDisks)
    right := buildRod(rightDisks, numDisks)

    printBoard(screen, boxStyle, left, middle, right, numLinesToRender, numDisks, 0)

    disks := map[int]*[]int{
        1: &leftDisks,
        2: &middleDisks,
        3: &rightDisks,
    }

    numMoves := 0

    for {
        draw(screen, 0, 0, boxStyle, false, "Welcome to Towers of Hanoi!")
        draw(screen, 0, 1, boxStyle, false, fmt.Sprintf("Optimal number of moves: %d", int(math.Pow(2, float64(numDisks)) - 1)))
        draw(screen, 0, 2, boxStyle, false, fmt.Sprintf("Number of Moves: %d", numMoves))
        draw(screen, 0, 3, boxStyle, false, "<Esc> or <Ctrl>+C to exit.")
        draw(screen, 0, 4, boxStyle, false, "Use the '1', '2', and '3' keys to select a disk and to then place the disk")

        screen.Show()

        ev := screen.PollEvent()

        switch ev := ev.(type) {
        case *tcell.EventKey:
            if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
                return
            }

            input := string(ev.Rune())
            if input != "1" && input != "2" && input != "3" {
                continue
            }

            if fromSelection {
                from, err = strconv.Atoi(input)

                if err != nil {
                    log.Fatalf("%+v", err)
                }
                fromSelection = false
                printBoard(screen, boxStyle, left, middle, right, numLinesToRender, numDisks, from)
                continue
            } else {
                to, err = strconv.Atoi(input)
                if err != nil {
                    log.Fatalf("%+v", err)
                }
                fromSelection = true
            }

            if len(*disks[from]) == 0 {
                printBoard(screen, boxStyle, left, middle, right, numLinesToRender, numDisks, 0)
                continue
            }

            selectedDisk := (*disks[from])[0]

            // Ensure there is either no disk on the 'to' rod, or that there is a smaller disk on the rod
            if len(*disks[to]) == 0 || selectedDisk < (*disks[to])[0] {
                *disks[from] = (*disks[from])[1:]
                *disks[to] = append([]int{selectedDisk}, *disks[to]...)

                left = buildRod(leftDisks, numDisks)
                middle = buildRod(middleDisks, numDisks)
                right = buildRod(rightDisks, numDisks)

                numMoves++
            }

            printBoard(screen, boxStyle, left, middle, right, numLinesToRender, numDisks, 0)
        }
    }
}
