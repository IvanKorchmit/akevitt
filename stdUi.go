package akevitt

import (
	"fmt"
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func AppendText(message string, chatlog io.Writer) {
	fmt.Fprintln(chatlog, message)
}

func ErrorBox(message string, app *tview.Application, back tview.Primitive) {
	result := tview.NewModal().SetText("Error!").SetText(message).SetTextColor(tcell.ColorRed).
		SetBackgroundColor(tcell.ColorBlack).
		AddButtons([]string{"Close"}).SetDoneFunc(func(buttonIndex int, buttonLabel string) {
		app.SetRoot(back, true)
	})

	result.SetBorder(true).SetBorderColor(tcell.ColorDarkRed)
	app.SetRoot(result, true)
}
