package gui

import (
	"encoding/base64"
	"fmt"

	"github.com/mechiko/walk"
	dcl "github.com/mechiko/walk/declarative"
)

func StartDialog(pub string, identity string) (out string, err error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var request, responce *walk.TextEdit

	toServer := fmt.Sprintf("%s@%s", pub, identity)
	encoded := base64.StdEncoding.EncodeToString([]byte(toServer))
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
								Text:     encoded,
							},
							dcl.PushButton{
								AssignTo: &acceptPB,
								Text:     "Copy",
								OnClicked: func() {
									walk.Clipboard().SetText(encoded)
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
									txt, _ := walk.Clipboard().Text()
									responce.SetText(txt)
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
		return "", fmt.Errorf("dialog return %d", ret)
	}
	return out, nil
}
