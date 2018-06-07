package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {

	MainWindow{
		Title:   "文章采集器V1.0",
		MinSize: Size{1000, 400},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					RadioButtonGroup{
						DataMember: "ZMTSite",
						Buttons: []RadioButton{
							RadioButton{
								Name:  "dayu",
								Text:  "大鱼号",
								Value: "dayu",
							},
							RadioButton{
								Name:  "baijia",
								Text:  "百家号",
								Value: "baijia",
							},
							RadioButton{
								Name:  "qie",
								Text:  "企鹅号",
								Value: "qie",
							},
						},
					},
				},
			},
			Composite{
				MaxSize: Size{0, 50},
				Layout:  HBox{},
				Children: []Widget{
					Label{Text: "ID: "},
					LineEdit{

						Text: "请输入作者ID",
					},
					Label{Text: "阅读量(下限): "},
					LineEdit{},
					Label{Text: "发布时间: "},
					DateEdit{
						Format: "2006-01-02 15:04:05",
					},
					Label{Text: "-- "},
					DateEdit{},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					PushButton{
						Text:    "开始抓取",
						MinSize: Size{120, 30},
					},
					TableView{
						AlternatingRowBGColor: walk.RGB(239, 239, 239),
						ColumnsOrderable:      true,
						Columns: []TableViewColumn{
							{Title: "#"},
							{Title: "类型"},
							{Title: "分类"},
							{Title: "标题"},
							{Title: "封面"},
							{Title: "时间", Format: "2006-01-02 15:04:05", Width: 150},
							{Title: "阅读量"},
							{Title: "链接"},
						},
					},
				},
			},
		},
	}.Run()
}
