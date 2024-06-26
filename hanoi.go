package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
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

func printBoard(left, middle, right [][]string, numLinesToRender int, numDisks int) {
    for i := 0; i < numLinesToRender; i++ {
        fmt.Print(strings.Join(left[i], ""))
        fmt.Print(strings.Join(middle[i], ""))
        fmt.Print(strings.Join(right[i], ""))
        fmt.Println()
    }
    fmt.Print(strings.Repeat("─", 3 * (numDisks * 4 + 7)))
    fmt.Println()
}

func main() {
    fmt.Println("Welcome to Towers of Hanoi!")
    fmt.Println()

    numDisks := 9
    numLinesToRender := 2 * numDisks + 4

    leftDisks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
    middleDisks := []int{}
    rightDisks := []int{}

    left := buildRod(leftDisks, numDisks)
    middle := buildRod(middleDisks, numDisks)
    right := buildRod(rightDisks, numDisks)

    printBoard(left, middle, right, numLinesToRender, numDisks)

    var input string

    disks := map[int]*[]int{
        1: &leftDisks,
        2: &middleDisks,
        3: &rightDisks,
    }

    for {
        fmt.Print("Select Rod: ")
        fmt.Scanln(&input)
        fmt.Println()

        from, err := strconv.Atoi(string(input[0]))
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        to, err := strconv.Atoi(string(input[1]))
        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        
        selectedDisk := (*disks[from])[0]
        *disks[from] = (*disks[from])[1:]

        *disks[to] = append([]int{selectedDisk}, *disks[to]...)

        left := buildRod(leftDisks, numDisks)
        middle := buildRod(middleDisks, numDisks)
        right := buildRod(rightDisks, numDisks)
        
        fmt.Print(strings.Repeat("\033[F", numLinesToRender + 3))
        printBoard(left, middle, right, numLinesToRender, numDisks)

    }
}
