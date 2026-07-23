package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/mats0319/secure_transfer/cmd/gui/components"
)

func main() {
	a := app.NewWithID("Secure Transfer") // 不链式调用下去，因为a可以创建多个window

	window := a.NewWindow("Secure Transfer")
	window.SetContent(container.NewVBox(makeUI(a)))
	window.SetMaster() // 主窗口，关闭它即视为关闭程序；窗口大小由内部组件撑开、不主动设置
	window.Show()

	a.Run()
}

func makeUI(a fyne.App) fyne.CanvasObject {
	top := components.TopBar(a)
	content := components.Content()
	bottom := components.Bottom()

	return container.NewBorder(top, bottom, nil, nil, content)
}
