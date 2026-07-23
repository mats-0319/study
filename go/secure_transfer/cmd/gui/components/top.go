package components

import (
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func TopBar(a fyne.App) fyne.CanvasObject {
	return container.NewVBox(topButtons(a), widget.NewSeparator(), workDir(a), widget.NewSeparator())
}

func topButtons(a fyne.App) *fyne.Container {
	aboutButton := widget.NewButton("About", func() {
		window := a.NewWindow("About")
		window.Resize(fyne.NewSize(600, 400))
		window.SetContent(widget.NewLabel(aboutText))
		window.Show()
	})

	helpButton := widget.NewButton("Help", func() {
		window := a.NewWindow("Help")
		window.Resize(fyne.NewSize(600, 400))
		window.SetContent(widget.NewLabel(helpText))
		window.Show()
	})

	return container.NewHBox(aboutButton, helpButton, container.NewHBox()) // 添加占位符，避免按钮占满行
}

func workDir(a fyne.App) *fyne.Container {
	path, err := os.Getwd()
	if err != nil {
		path = "./"
	}
	if !strings.HasSuffix(path, "/") { // must dir
		path += "/"
	}

	pathText := widget.NewLabel(fmt.Sprintf("Work Dir: %s", path))

	copyButton := widget.NewButtonWithIcon("", theme.ContentCopyIcon(), func() {
		a.Clipboard().SetContent(path)
		Log("Copied.")
	})

	return container.NewHBox(pathText, copyButton)
}
