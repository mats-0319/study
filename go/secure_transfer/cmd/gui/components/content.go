package components

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func Content() fyne.CanvasObject {
	displayContent := container.NewStack()

	var left *container.Split
	{
		background := canvas.NewRectangle(color.White)
		background.SetMinSize(fyne.NewSquareSize(200))

		var sender *fyne.Container
		{
			senderContent := Sender()
			displayContent.Add(senderContent) // default

			senderButton := widget.NewButton("Sender", func() {
				displayContent.Objects = []fyne.CanvasObject{senderContent}
				displayContent.Refresh()
			})

			sender = container.NewStack(background, container.NewCenter(senderButton))
		}

		var receiver *fyne.Container
		{
			receiverContent := Receiver()

			receiverButton := widget.NewButton("Receiver", func() {
				displayContent.Objects = []fyne.CanvasObject{receiverContent}
				displayContent.Refresh()
			})

			receiver = container.NewStack(background, container.NewCenter(receiverButton))
		}

		left = container.NewVSplit(sender, receiver)
		left.SetOffset(0.5)
	}

	mainContent := container.NewHSplit(left, displayContent)
	mainContent.SetOffset(0.25)

	return mainContent
}
