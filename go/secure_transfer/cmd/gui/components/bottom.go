package components

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type OperateLog struct {
	Time    string
	Details string
}

func (l *OperateLog) String() string {
	return fmt.Sprintf("> %s, %s", l.Time, l.Details)
}

var logData = make([]*OperateLog, 0)
var logList *widget.List

func Bottom() *fyne.Container {
	title := widget.NewLabel("> Operate Log: ")

	logList = widget.NewList(
		func() int { return len(logData) },
		func() fyne.CanvasObject { return widget.NewLabel("some text") },
		func(index widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(logData[index].String())
		},
	)

	background := canvas.NewRectangle(color.White)
	background.SetMinSize(fyne.NewSize(800, 200))

	logListWrapper := container.NewStack(background, logList)

	return container.NewVBox(widget.NewSeparator(), title, logListWrapper)
}

func Log(details string) {
	logData = append(logData, &OperateLog{
		Time:    time.Now().Format("2006-01-02 15:04:05.000"),
		Details: details,
	})

	if len(logData) > 1000 {
		logData = logData[len(logData)-1000:] // 保留最近1000条记录
	}

	logList.Refresh()
	logList.ScrollToBottom()
}
