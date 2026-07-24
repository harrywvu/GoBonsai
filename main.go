package main

func main() {
	w := NewWindow(80, 25)
    DrawBase(w, 40) // 40 = width/2, never hardcode once you add sin/cos
    w.Render()
}