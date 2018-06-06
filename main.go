package main

import (
	"strconv"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {

	mw := &MyMainWindow{}
	
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
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
								AssignTo: &mw.dayu,
							},
							RadioButton{
								Name:  "baijia",
								Text:  "百家号",
								Value: "baijia",
								AssignTo: &mw.baijia,
							},
							RadioButton{
								Name:  "qie",
								Text:  "企鹅号",
								Value: "qie",
								AssignTo: &mw.qie,
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
						AssignTo: &mw.id,
						Text: "请输入作者ID",
					},
					Label{Text: "阅读量: "},
					LineEdit{
						AssignTo: &mw.hot,
						Text: "请输入阅读量下限",
					},
					Label{Text: "发布时间: "},
					LineEdit{
						AssignTo: &mw.timeFrom,
						Text: "请输入发布时间范围",
					},
					Label{Text: "-- "},
					LineEdit{
						AssignTo: &mw.timeTo,
						Text: "请输入发布时间范围",
					},
				},
			},
			PushButton{
				Text:    "开始抓取",
				MinSize: Size{120, 30},
				AssignTo: &mw.Crawler,
				},
			},
			Composite{
				Layout: Grid{Columns: 2, Spacing: 10},
				Children: []Widget{
					ListBox{
						MaxSize:               Size{200, 0},
						AssignTo:              &mw.lb,
						OnCurrentIndexChanged: mw.lb_CurrentIndexChanged,
						OnItemActivated:       mw.lb_ItemActivated,
					},
					WebView{
						AssignTo: &mw.wv,
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err),
	}

	mw.Crawler.Clicked().Attach(func() {
		//TODO:添加爬虫入口
		return
	})

	//初始化
	mw.dayu.SetChecked(true)
	mw.Run()
}

type MyMainWindow struct {
	*walk.MainWindow
	lb				*walk.ListBox
	wv				*walk.WebView
	dayu			*walk.RadioButton
	baijia			*walk.RadioButton
	qie				*walk.RadioButton
	id				*walk.LineEdit
	hot				*walk.LineEdit
	timeFrom		*walk.LineEdit
	timeTo			*walk.LineEdit
}


func (mw *MyMainWindow) Crawler() {
	//进行参数校验
	id := mw.id.Text()
	if len(id) != 32 {
		walk.MsgBox(mw, "ID错误", "请填写正确的作者ID(*大鱼号的作者ID默认为32位*),作者ID需要与目标平台匹配", walk.MsgBoxIconWarning)
		return
	}
	hot := mw.hot.Text()
	if hotvalue, err := strconv.Atoi(hot); err != nil {
		walk.MsgBox(mw, "阅读量设置错误", "请填写正确的数字", walk.MsgBoxIconWarning)
		return
	}
	timeFrom := mw.timeFrom.Text()
	timeTo := mw.timeTo.Text()

	if mw.dayu.Checked() == true {
		//大鱼号平台爬取
	}

	if mw.baijia.Checked() == true {
		//百家号平台爬虫(暂不支持)
		walk.MsgBox(mw, "Sorry", "百家号平台爬虫暂不支持,请选择大鱼号平台", walk.MsgBoxIconInformation)
		return
	}

	if mw.qie.Checked() == true {
		//企鹅号平台爬虫(暂不支持)
		walk.MsgBox(mw, "Sorry", "企鹅号平台爬虫暂不支持,请选择大鱼号平台", walk.MsgBoxIconInformation)
		return
	}
}