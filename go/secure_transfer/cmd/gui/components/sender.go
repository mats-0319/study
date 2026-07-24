package components

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mats0319/secure_transfer/internal"
)

func makeSenderContent() *fyne.Container {
	titleText := widget.NewLabel("> Sender:")

	initializeFileButton := widget.NewButton("Initialize Message File", func() {
		internal.InitMessageFile()
		Log("Initialize Message File.")
	})
	initializeFileButton.Resize(fyne.NewSize(200, 100))

	var hasPubKey, hasMessage bool

	// 检查：如果将该组件设置为不可编辑，则组件内文字将使用灰色。
	// 该组件纯作为展示使用且内容每秒刷新，所以我认为该组件可编辑是可以接受的。接收者部分同理
	encryptCheckText := widget.NewMultiLineEntry()

	go func() {
		for range time.NewTicker(time.Second).C {
			fyne.Do(func() {
				hasPubKey, hasMessage = isFileExist("PUB.KEY", "message")

				text := fmt.Sprintf("Check 1 - Has Receiver Public Key. ('./PUB.KEY' file): %t\n"+
					"Check 2 - Has Message File. ('./message.xxx' file): %t", hasPubKey, hasMessage)
				encryptCheckText.SetText(text)
			})
		}
	}()

	encryptButton := widget.NewButton("Encrypt", func() {
		if !hasPubKey || !hasMessage {
			Log("Not Ready for Encrypt.")
			return
		}

		isSuccess := internal.Encrypt()
		Log("Encrypt Message " + boolToString(isSuccess) + ".")
	})

	blank40 := canvas.NewRectangle(color_MainDark)
	blank40.SetMinSize(fyne.NewSquareSize(40))

	return container.NewVBox(titleText, blank40, initializeFileButton, blank40, encryptCheckText, encryptButton)
}
