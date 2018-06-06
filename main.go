package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {

	MainWindow{
		Title:   "文章采集器V1.0",
		MinSize: Size{600, 400},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Rows: 3, Spacing: 0},
				Children: []Widget{
					VSplitter{
						Children: []Widget{
							Label{
								Text: "作者主页链接",
							},
						},
					},
					VSplitter{
						Children: []Widget{
							LineEdit{
								MinSize: Size{160, 0},
							},
						},
					},
					VSplitter{
						Children: []Widget{
							Label{
								Text: "发文时间段",
							},
						},
					},
					VSplitter{
						Children: []Widget{
							LineEdit{
								MinSize: Size{160, 0},
							},
						},
					},
					VSplitter{
						Children: []Widget{
							Label{
								Text: "阅读量",
							},
						},
					},
					VSplitter{
						Children: []Widget{
							LineEdit{
								MinSize: Size{160, 0},
							},
						},
					},
				},
			},
			PushButton{
				Text:    "开始抓取",
				MinSize: Size{120, 50},
				OnClicked: func() {
					var tmp walk.Form
					walk.MsgBox(tmp, "hello Bean", "", walk.MsgBoxIconInformation)
				},
			},
		},
	}.Run()
}
