package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

func main() {

	mw := &MyMainWindow{}
	model := NewResultsTableModel()
	if err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "文章采集器V1.0",
		MinSize:  Size{1000, 400},
		Layout:   VBox{},
		Children: []Widget{
			Composite{
				Layout: HBox{},
				Children: []Widget{
					RadioButtonGroup{
						DataMember: "ZMTSite",
						Buttons: []RadioButton{
							RadioButton{
								Name:     "dayu",
								Text:     "大鱼号",
								Value:    "dayu",
								AssignTo: &mw.dayu,
							},
							RadioButton{
								Name:     "baijia",
								Text:     "百家号",
								Value:    "baijia",
								AssignTo: &mw.baijia,
							},
							RadioButton{
								Name:     "qie",
								Text:     "企鹅号",
								Value:    "qie",
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
						Text:     "请输入作者ID",
					},
					Label{Text: "阅读量(下限): "},
					// NumberEdit{
					// 	AssignTo: &mw.hot,
					// 	Value:    0,
					// },
					LineEdit{
						AssignTo: &mw.hot,
						Text:     "0",
					},
					Label{Text: "发布时间: "},
					DateEdit{
						AssignTo: &mw.timeFrom,
					},
					Label{Text: "-- "},
					DateEdit{
						AssignTo: &mw.timeTo,
					},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					PushButton{
						Text: "开始抓取",
						//MinSize:  Size{120, 30},
						OnClicked: model.ResetRows,
					},
					TableView{
						AssignTo:              &mw.tv,
						AlternatingRowBGColor: walk.RGB(239, 239, 239),
						CheckBoxes:            true,
						MultiSelection:        true,
						ColumnsOrderable:      true,
						Columns: []TableViewColumn{
							{Title: "#"},
							{Title: "类型"},
							{Title: "分类"},
							{Title: "标题"},
							{Title: "封面"},
							{Title: "时间"},
							{Title: "阅读量"},
							{Title: "链接"},
						},
						//TODO:加入表格风格，根据值的大小高亮显示
						Model: model,
						OnCurrentIndexChanged: func() {
							i := mw.tv.CurrentIndex()
							if 0 <= i {
								fmt.Printf("OnCurrentIndexChanged: %v\n", model.items[i].Title)
							}
						},
						// OnItemActivated: mw.tv_ItemActivated,
					},
					//TODO:暂时不添加网页浏览功能
					// WebView{
					// 	AssignTo: &mw.wv,
					// },
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	// // 处理抓取状态
	// mw.start.Clicked().Attach(func() {
	// 	go func() {
	// 		mw.start.SetText("正在努力抓取中···")
	// 		mw.start.SetEnabled(false)
	// 		mw.ResetRows()
	// 		mw.start.SetText("开始抓取")
	// 		mw.start.SetEnabled(true)
	// 	}()
	// })

	//初始化
	mw.dayu.SetChecked(true)
	mw.Run()
}

// MyMainWindow 整体界面结构
type MyMainWindow struct {
	*walk.MainWindow
	tv *walk.TableView
	//model *ResultsTableModel
	// wv       *walk.WebView
	//start    *walk.PushButton
	dayu     *walk.RadioButton
	baijia   *walk.RadioButton
	qie      *walk.RadioButton
	id       *walk.LineEdit
	hot      *walk.LineEdit
	timeFrom *walk.DateEdit
	timeTo   *walk.DateEdit
}

// ResultsTable 表格结构
type ResultsTable struct {
	Index       int
	Type        string
	Category    string
	Title       string
	Coverlink   string
	Publishtime time.Time
	Hot         int
	Srcurl      string
	checked     bool
}

// ResultsTableModel 表格模型
type ResultsTableModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*ResultsTable
}

/*====================全===局===变===量======================*/

// 定义全局数据集给主函数表格
// var model *ResultsTableModel

//作者ID
var AuthorID string

//阅读量下限
var Hotvalue string

//时间段from
var Timefrom time.Time

//时间段Timeto
var Timeto time.Time

/*================爬===虫===入===口==========================*/

// Resp 响应结构
type Resp struct {
	Data []Production `json:"data"`
}

// Production 作者创作列表json结构
type Production struct {
	ID          string `json:"content_id"`   //作品ID
	Type        int    `json:"format_type"`  //作品类型
	Category    string `json:"category"`     //作品分类
	Title       string `json:"title"`        //作品标题
	Coverlink   string `json:"cover_url"`    //封面链接
	Publishtime string `json:"published_at"` //发表时间
}

// ProductionHot 获取作品的热度(阅读量)需要构造的json结构
// type ProductionHot struct {
// 	//data. _incrs. click_total
// 	Data struct {
// 		Info struct {
// 			HotValue int `json:"click_total"`
// 		} `json:"_incrs"`
// 	} `json:"data"`
// }

// Result 自己需要的结构
type Result struct {
	Url         string
	Type        string
	Category    string
	Title       string
	Coverlink   string
	Publishtime time.Time
	Hot         int
}

// Dayu 大鱼号
func Dayu() (results []Result) {

	results = []Result{}

	c := colly.NewCollector(
		colly.Async(true),
	)
	//设置请求超时
	c.SetRequestTimeout(200 * time.Second)
	extensions.RandomUserAgent(c)
	extensions.Referrer(c)
	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting: ", r.URL.String())
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
		r.Request.Retry()
	})

	getHotCollector := c.Clone()
	getHotCollector.OnResponse(func(resp *colly.Response) {
		result := Result{}
		reg := regexp.MustCompile(`"click_total":(.*?),`)
		result.Hot, _ = strconv.Atoi(reg.FindString(string(resp.Body)))

		// 根据阅读量筛选
		hv, _ := strconv.Atoi(Hotvalue)
		if result.Hot < hv {
			return
		}

		id := resp.Request.Ctx.Get("id")
		result.Url = fmt.Sprintf("html: http://a.mp.uc.cn/article.html?uc_param_str=frdnsnpfvecpntnwprdssskt&from=media#!wm_cid=%s", id)
		result.Type = resp.Request.Ctx.Get("type")
		result.Category = resp.Request.Ctx.Get("category")
		result.Title = resp.Request.Ctx.Get("title")
		result.Coverlink = resp.Request.Ctx.Get("coverlink")
		result.Publishtime, _ = time.Parse("2006-01-02 15:04:05", resp.Request.Ctx.Get("publishtime"))

		results = append(results, result)
	})

	c.OnResponse(func(resp *colly.Response) {
		productionlist := &Resp{}
		if err := json.Unmarshal([]byte(string(resp.Body)), productionlist); err != nil {
			log.Fatal(err)
		}
		for _, oneproduction := range productionlist.Data {

			//根据发布的时间段筛选
			changetime := strings.Replace(strings.Split(oneproduction.Publishtime, ".")[0], "T", " ", -1)
			timeline, _ := time.Parse("2006-01-02 15:04:05", changetime)

			// if timeline.Before(Timefrom) && Timeto.Before(timeline) {
			// 	continue
			// }
			if (Timefrom.Before(timeline) && timeline.Before(Timeto)) || (Timeto.Before(timeline) && timeline.Before(Timefrom)) {
				id := oneproduction.ID
				targetURL := fmt.Sprintf("http://ff.dayu.com/contents/%s?biz_id=1002&_fetch_author=1&_fetch_incrs=1", id)
				ctx := colly.NewContext()
				ctx.Put("id", oneproduction.ID)
				ctx.Put("type", strconv.Itoa(oneproduction.Type))
				ctx.Put("category", oneproduction.Category)
				ctx.Put("title", oneproduction.Title)
				ctx.Put("coverlink", oneproduction.Coverlink)
				ctx.Put("publishtime", timeline)
				getHotCollector.Request("GET", targetURL, nil, ctx, nil)
				log.Println("Visiting: ", targetURL)
			}
		}

	})

	visitUrl := fmt.Sprintf("http://ff.dayu.com/contents/author/%s?biz_id=1002", AuthorID)
	c.Visit(visitUrl)

	getHotCollector.Wait()

	return results

}

/*================ENDING==========================*/

// RowCount 刷新列表
func (m *ResultsTableModel) RowCount() int {
	return len(m.items)
}

// Value 结果渲染在table中
func (m *ResultsTableModel) Value(row, col int) interface{} {
	item := m.items[row]

	switch col {
	case 0:
		return item.Index
	case 1:
		return item.Type
	case 2:
		return item.Category
	case 3:
		return item.Title
	case 4:
		return item.Coverlink
	case 5:
		return item.Publishtime
	case 6:
		return item.Hot
	case 7:
		return item.Srcurl
	}
	panic("Unexpected col")
}

func (m *ResultsTableModel) Checked(row int) bool {
	return m.items[row].checked
}
func (m *ResultsTableModel) SetChecked(row int, checked bool) error {
	m.items[row].checked = checked
	return nil
}
func (m *ResultsTableModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.Stable(m)
	return m.SorterBase.Sort(col, order)
}
func (m *ResultsTableModel) Len() int {
	return len(m.items)
}
func (m *ResultsTableModel) Swap(i, j int) {
	m.items[i], m.items[j] = m.items[j], m.items[i]
}
func (m *ResultsTableModel) Less(i, j int) bool {
	a, b := m.items[i], m.items[j]
	c := func(ls bool) bool {
		if m.sortOrder == walk.SortAscending {
			return ls
		}
		return !ls
	}

	switch m.sortColumn {
	case 0:
		return c(a.Index < b.Index)
	case 1:
		return c(a.Type < b.Type)
	case 2:
		return c(a.Category < b.Category)
	case 3:
		return c(a.Title < b.Title)
	case 4:
		return c(a.Coverlink < b.Coverlink)
	case 5:
		return c(a.Publishtime.Before(b.Publishtime))
	case 6:
		return c(a.Hot < b.Hot)
	}
	panic("Unreachable")
}

func NewResultsTableModel() *ResultsTableModel {
	m := new(ResultsTableModel)
	return m
}

func (m *ResultsTableModel) ResetRows() {
	//进行参数校验
	mw := &MyMainWindow{}
	log.Println(*(&mw.id))
	AuthorID = "a2c99b15af2b413ea29c6ebf40b9750c"
	if len(AuthorID) != 32 {
		walk.MsgBox(mw, "ID错误", "请填写正确的作者ID(*大鱼号的作者ID默认为32位*),作者ID需要与目标平台匹配", walk.MsgBoxIconWarning)
		return
	}
	Hotvalue := mw.hot.Text()
	//hot := mw.hot.Text()
	if _, err := strconv.Atoi(Hotvalue); err != nil {
		walk.MsgBox(mw, "阅读量设置错误", "请填写正确的数字", walk.MsgBoxIconWarning)
		return
	}
	Timefrom = mw.timeFrom.Date()
	Timeto = mw.timeTo.Date()

	if mw.dayu.Checked() == true {
		//大鱼号平台爬取
		// mw.model.items = NewResultsTableModel("dayu", Dayu())
		// // m := new(ResultsTableModel)
		results := Dayu()
		m.items = make([]*ResultsTable, 50000)
		for i, result := range results {

			// 格式转换成界面表格的格式
			if result.Type == "1001" {
				result.Type = "图文"
			} else if result.Type == "1002" {
				result.Type = "视频"
			} else if result.Type == "1005" {
				result.Type = "图集"
			} else {
				result.Type = "未知类型" + result.Type
			}

			m.items[i] = &ResultsTable{
				Index:       i,
				Type:        result.Type,
				Category:    result.Category,
				Title:       result.Title,
				Coverlink:   result.Coverlink,
				Publishtime: result.Publishtime,
				Hot:         result.Hot,
				Srcurl:      result.Url,
			}
		}
		m.PublishRowsReset()
		m.Sort(m.sortColumn, m.sortOrder)
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

// Crawler 开启爬虫后的主函数入口
// func (mw *MyMainWindow) Crawler() {
// 	//进行参数校验
// 	id := mw.id.Text()
// 	if len(id) != 32 {
// 		walk.MsgBox(mw, "ID错误", "请填写正确的作者ID(*大鱼号的作者ID默认为32位*),作者ID需要与目标平台匹配", walk.MsgBoxIconWarning)
// 		return
// 	}
// 	Hotvalue := mw.hot.Text()
// 	//hot := mw.hot.Text()
// 	if _, err := strconv.Atoi(Hotvalue); err != nil {
// 		walk.MsgBox(mw, "阅读量设置错误", "请填写正确的数字", walk.MsgBoxIconWarning)
// 		return
// 	}
// 	// if mw.timeFrom.Text() == "" {
// 	// 	Timefrom = nil
// 	// } else {
// 	// 	Timefrom, err := time.Parse("2006-01-02 15:04:05", mw.timeFrom.Text())
// 	// 	if err != nil {
// 	// 		walk.MsgBox(mw, "Time Error", "请填写正确格式的时间段(2006-01-02 15:04:05)", walk.MsgBoxIconWarning)
// 	// 		return
// 	// 	}
// 	// }
// 	// if mw.timeTo.Text() == "" {
// 	// 	Timeto = nil
// 	// } else {
// 	// 	Timeto, err := time.Parse("2006-01-02 15:04:05", mw.timeTo.Text())
// 	// 	if err != nil {
// 	// 		walk.MsgBox(mw, "Time Error", "请填写正确格式的时间段(2006-01-02 15:04:05)", walk.MsgBoxIconWarning)
// 	// 		return
// 	// 	}
// 	// }
// 	Timefrom = mw.timeFrom.Date()
// 	Timeto = mw.timeTo.Date()

// 	if mw.dayu.Checked() == true {
// 		//大鱼号平台爬取
// 		// mw.model.items = NewResultsTableModel("dayu", Dayu())
// 		// // m := new(ResultsTableModel)
// 		results := Dayu()
// 		mw.model.items = make([]*ResultsTable, 7)
// 		for i, result := range results {

// 			// 格式转换成界面表格的格式
// 			if result.Type == "1001" {
// 				result.Type = "图文"
// 			} else if result.Type == "1002" {
// 				result.Type = "视频"
// 			} else if result.Type == "1005" {
// 				result.Type = "图集"
// 			} else {
// 				result.Type = "未知类型" + result.Type
// 			}

// 			mw.model.items[i] = &ResultsTable{
// 				Index:       i,
// 				Type:        result.Type,
// 				Category:    result.Category,
// 				Title:       result.Title,
// 				Coverlink:   result.Coverlink,
// 				Publishtime: result.Publishtime,
// 				Hot:         result.Hot,
// 				Srcurl:      result.Url,
// 			}
// 		}
// 		mw.model.PublishRowsReset()
// 		mw.tv.SetSelectedIndexes([]int{})
// 	}

// 	if mw.baijia.Checked() == true {
// 		//百家号平台爬虫(暂不支持)
// 		walk.MsgBox(mw, "Sorry", "百家号平台爬虫暂不支持,请选择大鱼号平台", walk.MsgBoxIconInformation)
// 		return
// 	}

// 	if mw.qie.Checked() == true {
// 		//企鹅号平台爬虫(暂不支持)
// 		walk.MsgBox(mw, "Sorry", "企鹅号平台爬虫暂不支持,请选择大鱼号平台", walk.MsgBoxIconInformation)
// 		return
// 	}
// }

/*================辅助工具======================*/

// OpenBrowser 打开默认浏览器
func OpenBrowser(uri string) error {
	cmd := exec.Command("rundll32 url.dll,FileProtocolHandler", uri)
	return cmd.Start()
}

// StringtimeToTime 字符串时间转换成时间类型
// func StringtimeToTime(t string) time.Time {
// 	timeline, _ = time.Parse("2018-06-07 15:00:00", t)
// 	return timeline
// }
