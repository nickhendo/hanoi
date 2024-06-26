package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)


func repeated(value string, number int) []string {
    arr := make([]string, number)
    for i := 0; i < number; i++ {
        arr[i] = value
    }
    return arr
}
    

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
    line2 = append(line2, repeated(" ", width * 2 + 2)...)
    line2 = append(line2, "│")
    line2 = append(line2, repeated(" ", ((maxWidth - width) * 2))...)

    disk = append(disk, line1)
    disk = append(disk, line2)
    
    return disk
}

func buildRod(rodDef []int, maxWidth int) [][]string {
    var rodImage = [][]string{}

    // Rod length is 4 higher than a full stacked rod
    for i := 0; i < 4 + 2 * (maxWidth - len(rodDef)); i++ {
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

func printBoard(screen tcell.Screen, style tcell.Style, left, middle, right [][]string, numLinesToRender int, numDisks int) {
    offset := 2
    for i := 0; i < numLinesToRender; i++ {
        draw(screen, 0, offset + i, style, strings.Join(left[i], ""))
        draw(screen, len(left[i]), offset + i, style, strings.Join(middle[i], ""))
        draw(screen, len(left[i]) + len(right[i]), offset + i, style, strings.Join(right[i], ""))
    }
    draw(screen, 0, numLinesToRender + offset, style, strings.Repeat("─", 3 * (numDisks * 4 + 7)))
}

func draw(screen tcell.Screen, x, y int, style tcell.Style, text string) {
    col := x
    row := y

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

    numDisks := 9
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

    printBoard(screen, boxStyle, left, middle, right, numLinesToRender, numDisks)

    disks := map[int]*[]int{
        1: &leftDisks,
        2: &middleDisks,
        3: &rightDisks,
    }

    for {
 
        draw(screen, 0, 0, boxStyle, "Welcome to Towers of Hanoi!")

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
            log.Printf("input: %s", input)
            if fromSelection {
                from, err = strconv.Atoi(input)
                log.Printf("firstFrom: %d", from)

                if err != nil {
                    log.Fatalf("%+v", err)
                }
                fromSelection = false
                continue
            } else {
                to, err = strconv.Atoi(input)
                if err != nil {
                    log.Fatalf("%+v", err)
                }
                fromSelection = true
            }
            log.Printf("from: %d", from)
            log.Printf("to: %d", to)
            log.Printf("disks: %+v", disks)

            log.Printf("leftDisk: %+v", *disks[1])
            log.Printf("middleDisk: %+v", *disks[2])
            log.Printf("rightDisk: %+v", *disks[3])
            var left [][]string
            var middle[][]string
            var right[][]string

            if len(*disks[from]) == 0 {
                continue
            }

            selectedDisk := (*disks[from])[0]

            // Ensure there is either no disk on the 'to' rod, or that there is a smaller disk on the rod
            // if len(*disks[to]) == 0 || selectedDisk < (*disks[to])[0] {
            if len(*disks[to]) == 0 || selectedDisk < (*disks[to])[0] {
                log.Println("here")

                *disks[from] = (*disks[from])[1:]
                *disks[to] = append([]int{selectedDisk}, *disks[to]...)

                log.Printf("disks: %+v", disks)
                log.Printf("leftDisk: %+v", *disks[1])
                log.Printf("middleDisk: %+v", *disks[2])
                log.Printf("rightDisk: %+v", *disks[3])

                left = buildRod(leftDisks, numDisks)
                middle = buildRod(middleDisks, numDisks)
                right = buildRod(rightDisks, numDisks)

                // fmt.Print(strings.Repeat("\033[F", numLinesToRender + 3))
                printBoard(screen, boxStyle, left, middle, right, numLinesToRender, numDisks)
            }

            
        }
    }
}
