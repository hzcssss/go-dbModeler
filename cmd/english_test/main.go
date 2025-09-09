package main

import (
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	// Create application
	a := app.New()
	w := a.NewWindow("English Interface Test")

	// Create labels
	label1 := widget.NewLabel("This is a test for English display")
	label2 := widget.NewLabel("If you can see this text clearly, the interface works correctly")
	label3 := widget.NewLabel("No more encoding issues!")

	// Create button
	button := widget.NewButton("Click Me", func() {
		label1.SetText("Button was clicked!")
	})

	// Create layout
	content := container.NewVBox(
		label1,
		label2,
		label3,
		button,
	)

	// Set window content
	w.SetContent(content)

	// Show window
	w.ShowAndRun()
}
