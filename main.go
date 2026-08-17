package main

func main() {
    w := NewWindow(80, 25)
    DrawBase(w, 40)

    startX, startY := 40.0, 8.0
    endX, endY := getEnd(startX, startY, 8, 0)

    for y := startY; y <= endY; y++ {
        w.SetCharPlane(startX, y, '|')
    }

    w.Render()

}