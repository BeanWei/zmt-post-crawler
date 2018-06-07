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
					Label{Text: "阅读量: "},
					LineEdit{
						AssignTo: &mw.hot,
						Text:     "请输入阅读量下限",
					},
					Label{Text: "发布时间: "},
					LineEdit{
						AssignTo: &mw.timeFrom,
						Text:     "请输入发布时间范围",
					},
					Label{Text: "-- "},
					LineEdit{
						AssignTo: &mw.timeTo,
						Text:     "请输入发布时间范围",
					},
				},
			},
			PushButton{
				Text:     "开始抓取",
				MinSize:  Size{120, 30},
				AssignTo: &mw.Start,
			},
			Composite{
				Layout: Grid{Columns: 2, Spacing: 10},
				Children: []Widget{
					TableView{
						AssignTo:              &mw.tv,
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
						//TODO:加入表格风格，根据值的大小高亮显示
						Model: model,
						OnCurrentIndexChanged: func() {
							i := mw.tv.CurrentIndex()
							if 0 <= i {
								fmt.Printf("OnCurrentIndexChanged: %v\n", model.items[i].Title)
							}
						},
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

	// 处理抓取状态
	mw.Start.Clicked().Attach(func() {
		go func() {
			mw.Start.SetText("正在努力抓取中···")
			mw.Start.SetEnabled(false)
			mw.Crawler()
			mw.Start.SetText("开始抓取")
			mw.Start.SetEnabled(true)
		}()
	})

	//初始化
	mw.dayu.SetChecked(true)
	mw.Run()
}

// MyMainWindow 整体界面结构
type MyMainWindow struct {
	*walk.MainWindow
	tv *walk.TableView
	// wv       *walk.WebView
	Start    *walk.PushButton
	dayu     *walk.RadioButton
	baijia   *walk.RadioButton
	qie      *walk.RadioButton
	id       *walk.LineEdit
	hot      *walk.LineEdit
	timeFrom *walk.LineEdit
	timeTo   *walk.LineEdit
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
}

// ResultsTableModel 表格模型
type ResultsTableModel struct {
	walk.TableModelBase
	walk.SorterBase
	sortColumn int
	sortOrder  walk.SortOrder
	items      []*ResultsTable
}

// 定义全局数据集给主函数表格
var model ResultsTableModel

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
	Publishtime string
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
		id := resp.Request.Ctx.Get("id")
		result.Url = fmt.Sprintf("html: http://a.mp.uc.cn/article.html?uc_param_str=frdnsnpfvecpntnwprdssskt&from=media#!wm_cid=%s", id)
		result.Type = resp.Request.Ctx.Get("type")
		result.Category = resp.Request.Ctx.Get("category")
		result.Title = resp.Request.Ctx.Get("title")
		result.Coverlink = resp.Request.Ctx.Get("coverlink")
		result.Publishtime = resp.Request.Ctx.Get("publishtime")

		results = append(results, result)
	})

	c.OnResponse(func(resp *colly.Response) {
		productionlist := &Resp{}
		if err := json.Unmarshal([]byte(string(resp.Body)), productionlist); err != nil {
			log.Fatal(err)
		}
		for _, oneproduction := range productionlist.Data {
			id := oneproduction.ID
			targetURL := fmt.Sprintf("http://ff.dayu.com/contents/%s?biz_id=1002&_fetch_author=1&_fetch_incrs=1", id)
			ctx := colly.NewContext()
			ctx.Put("id", oneproduction.ID)
			ctx.Put("type", strconv.Itoa(oneproduction.Type))
			ctx.Put("category", oneproduction.Category)
			ctx.Put("title", oneproduction.Title)
			ctx.Put("coverlink", oneproduction.Coverlink)
			ctx.Put("publishtime", oneproduction.Publishtime)
			getHotCollector.Request("GET", targetURL, nil, ctx, nil)
			log.Println("Visiting: ", targetURL)
		}

	})

	c.Visit("http://ff.dayu.com/contents/author/a2c99b15af2b413ea29c6ebf40b9750c?biz_id=1002")

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

// Sort	表格排序
func (m *ResultsTableModel) Sort(col int, order walk.SortOrder) error {
	m.sortColumn, m.sortOrder = col, order
	sort.SliceStable(m.items, func(i, j int) bool {
		a, b := m.items[i], m.items[j]
		c := func(ls bool) bool {
			if m.sortOrder == walk.SortAscending {
				return ls
			}
			return !ls
		}

		switch m.sortColumn {
		case 5:
			return c(a.Publishtime.Before(b.Publishtime))
		case 6:
			return c(a.Hot < b.Hot)
		}
		panic("Unreachable")
	})

	return m.SorterBase.Sort(col, order)
}

// TODO:Swap 交换列

// ResetRows 重置表格
// func (m *ResultsTableModel) ResetRows() {

// }

// NewResultsTableModel	数据源(平台文章爬虫)
func NewResultsTableModel(site string, results []Result) *ResultsTableModel {
	m := new(ResultsTableModel)
	m.items = make([]*ResultsTable, 7)
	for i, result := range results {
		var (
			protype  string
			timeline time.Time
		)

		switch site {
		case "dayu":
			// 格式转换成界面表格的格式
			if result.Type == "1001" {
				protype := "图文"
			} else if result.Type == "1002" {
				protype := "视频"
			} else if result.Type == "1005" {
				protype := "图集"
			} else {
				protype := "未知类型" + result.Type
			}
			changetime := strings.Replace(strings.Split(result.Publishtime, ".")[0], "T", " ", -1)
			timeline, _ := time.Parse("2018-06-07 15:00:00", changetime)
		}

		m.items[i] = &ResultsTable{
			Index:       i,
			Type:        protype,
			Category:    result.Category,
			Title:       result.Title,
			Coverlink:   result.Coverlink,
			Publishtime: timeline,
			Hot:         result.Hot,
			Srcurl:      result.Url,
		}
	}
	return m
}

// Crawler 开启爬虫后的主函数入口
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
		model = NewResultsTableModel("dayu", Dayu())

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

/*================辅助工具======================*/

// OpenBrowser 打开默认浏览器
func OpenBrowser(uri string) error {
	cmd := exec.Command("rundll32 url.dll,FileProtocolHandler", uri)
	return cmd.Start()
}
