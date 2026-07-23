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

func Receiver() *fyne.Container {
	title := widget.NewLabel("> Receiver:")

	generateKeyPairButton := widget.NewButton("Generate Key Pair", func() {
		internal.GenerateKeypair()
		Log("Generate Key Pair.")
	})

	var hasPrivKey, hasMessage bool

	decryptCheckText := widget.NewMultiLineEntry()

	go func() {
		for range time.NewTicker(time.Second).C {
			fyne.Do(func() {
				hasPrivKey, hasMessage = readyForDecrypt()

				text := fmt.Sprintf("Check 1 - Has Private Key. ('./priv.key' file): %t\n"+
					"Check 2 - Has Cipher File. ('./CIPHER.XXX' file): %t", hasPrivKey, hasMessage)
				decryptCheckText.SetText(text)
			})
		}
	}()

	decryptButton := widget.NewButton("Decrypt", func() {
		if !hasPrivKey || !hasMessage {
			Log("Not Ready for Decrypt.")
			return
		}

		isSuccess := internal.Decrypt()
		Log("Decrypt Message " + boolToString(isSuccess) + ".")
	})

	blank40 := canvas.NewRectangle(color.White)
	blank40.SetMinSize(fyne.NewSquareSize(40))

	return container.NewVBox(title, blank40, generateKeyPairButton, blank40, decryptCheckText, decryptButton)
}

func readyForDecrypt() (hasPrivKey bool, hasMessage bool) {
	entry, err := os.ReadDir("./")
	if err != nil {
		Log("Check for Decrypt Failed.")
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
		case fileInfo.Name() == "priv.key":
			hasPrivKey = true
		case strings.HasPrefix(fileInfo.Name(), "CIPHER"):
			hasMessage = true
		}

		if hasPrivKey && hasMessage {
			break // all check passed, break and return
		}
	}

	return
}
