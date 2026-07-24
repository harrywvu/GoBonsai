package main

func main() {
    w := NewWindow(80, 25)
    DrawBase(w, 40)

    trunkCol, trunkRow := 40, 22
    topCol, topRow := getEnd(trunkCol, trunkRow, 8, 0)

    for r := trunkRow; r >= topRow; r-- {
        w.SetChar(r, topCol, '|')
    }

    w.Render()
}