package licenser

import (
	"encoding/base64"
	"fmt"

	"github.com/mechiko/utility"
	"github.com/mechiko/walk"
	dcl "github.com/mechiko/walk/declarative"
	"golang.design/x/clipboard"
)

func (l *Licenser) startDialog() (out string, err error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var request, responce *walk.TextEdit

	toServer := fmt.Sprintf("%s@%s", l.box.LocalPub, l.ParameterValue)
	encoded64 := base64.StdEncoding.EncodeToString([]byte(toServer))
	icon, err := walk.Resources.Icon("3")
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	if err := (dcl.Dialog{
		AssignTo:      &dlg,
		Title:         "Активация",
		Size:          dcl.Size{Width: 700, Height: 400},
		Icon:          icon,
		Layout:        dcl.VBox{Spacing: 10, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 5}},
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		Children: []dcl.Widget{
			dcl.GroupBox{
				Title:  "Параметры",
				Layout: dcl.VBox{MarginsZero: false, SpacingZero: false, Margins: dcl.Margins{Left: 10, Top: 10, Right: 10, Bottom: 10}},
				Children: []dcl.Widget{
					dcl.Composite{
						Layout:    dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
						Border:    false,
						Alignment: dcl.Alignment2D(walk.AlignHNearVNear),
						Children: []dcl.Widget{
							dcl.Label{
								Text: "Ключ:",
							},
							dcl.TextEdit{
								AssignTo: &request,
								Enabled:  false,
								MinSize:  dcl.Size{Width: 500},
								Text:     encoded64,
							},
							dcl.PushButton{
								AssignTo: &acceptPB,
								Text:     "Copy",
								OnClicked: func() {
									clipboard.Write(clipboard.FmtText, []byte(encoded64))
									// walk.Clipboard().SetText(encoded64)
								},
							},
							dcl.HSpacer{},
						},
					},
					dcl.Composite{
						Layout:    dcl.HBox{MarginsZero: true, SpacingZero: false, Margins: dcl.Margins{Left: 0, Top: 0, Right: 0, Bottom: 0}},
						Border:    false,
						Alignment: dcl.Alignment2D(walk.AlignHNearVNear),
						Children: []dcl.Widget{
							dcl.Label{
								Text: "Лицензия:",
							},
							dcl.TextEdit{
								AssignTo: &responce,
								Text:     "",
							},
							dcl.PushButton{
								AssignTo: &acceptPB,
								Text:     "Paste",
								OnClicked: func() {
									data := clipboard.Read(clipboard.FmtText)
									if data == nil {
										utility.MessageBox("ошибка", "Clipboard is empty or does not contain text data")
										return
									}
									responce.SetText(string(data))
									// txt, _ := walk.Clipboard().Text()
									// responce.SetText(txt)
								},
							},
							dcl.HSpacer{},
						},
					},
				}},
			dcl.Composite{
				Border: false,
				Layout: dcl.HBox{MarginsZero: true, SpacingZero: true, Margins: dcl.Margins{Left: 5, Top: 5, Right: 5, Bottom: 0}},
				Children: []dcl.Widget{
					dcl.PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							out = responce.Text()
							dlg.Accept()
						},
					},
					dcl.PushButton{
						AssignTo:  &cancelPB,
						Text:      "Выход",
						OnClicked: func() { dlg.Cancel() },
					},
				},
			},
			dcl.VSpacer{},
		},
	}).Create(nil); err != nil {
		return "", fmt.Errorf("%w", err)
	}

	if ret := dlg.Run(); ret != 1 {
		return "", fmt.Errorf("отмена ввода лицензии")
	}
	if out == "" {
		return "", fmt.Errorf("пустое значение лицензии")
	}
	return out, nil
}
