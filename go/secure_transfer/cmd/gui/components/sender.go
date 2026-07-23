package components

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/mats0319/secure_transfer/internal"
)

func Sender() *fyne.Container {
	titleText := widget.NewLabel("> Sender:")

	initializeFileButton := widget.NewButton("Initialize Message File", func() {
		internal.InitMessageFile()
		Log("Initialize Message File.")
	})
	initializeFileButton.Resize(fyne.NewSize(200, 100))

	var hasPubKey, hasMessage bool

	// 检查：如果将该组件设置为不可编辑，则组件内文字将使用灰色。
	// 该组件纯作为展示使用且内容每秒刷新，所以我认为该组件可编辑是可以接受的。
	encryptCheckText := widget.NewMultiLineEntry()

	go func() {
		for range time.NewTicker(time.Second).C {
			fyne.Do(func() {
				hasPubKey, hasMessage = readyForEncrypt()

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

	blank40 := canvas.NewRectangle(color.White)
	blank40.SetMinSize(fyne.NewSquareSize(40))

	return container.NewVBox(titleText, blank40, initializeFileButton, blank40, encryptCheckText, encryptButton)
}

func readyForEncrypt() (hasPubKey bool, hasMessage bool) {
	entry, err := os.ReadDir("./")
	if err != nil {
		Log("Check for Encrypt Failed.")
		return
	}

	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		fileInfo, err := entry[i].Info()
		if err != nil {
			Log("Get File Info Failed.")
			continue
		}

		switch {
		case fileInfo.Name() == "PUB.KEY":
			hasPubKey = true
		case strings.HasPrefix(fileInfo.Name(), "message"):
			hasMessage = true
		}

		if hasPubKey && hasMessage {
			break // all check passed, break and return
		}
	}

	return
}
