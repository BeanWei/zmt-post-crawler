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

// MyMainWindow 整体界面结构
type MyMainWindow struct {
	*walk.MainWindow
	tv        *walk.TableView
	model     *ResultsTableModel
	start     *walk.PushButton
	dayu      *walk.RadioButton
	baijia    *walk.RadioButton
	qie       *walk.RadioButton
	sharelink *walk.LineEdit
	getID     *walk.PushButton
	idvalue   *walk.LineEdit
	id        *walk.LineEdit
	hot       *walk.LineEdit
	timeFrom  *walk.DateEdit
	timeTo    *walk.DateEdit
}

func NewResultsTableModel() *ResultsTableModel {
	m := new(ResultsTableModel)
	return m
}

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
	case 7:
		return c(a.Srcurl < b.Srcurl)
	}
	panic("Unreachable")
}

func main() {

	mw := &MyMainWindow{model: NewResultsTableModel()}
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
				Layout: HBox{},
				Children: []Widget{
					LineEdit{
						AssignTo: &mw.sharelink,
						Text:     "输入分享链接",
					},
					PushButton{
						Text:     "解析作者ID",
						AssignTo: &mw.getID,
					},
					LineEdit{
						AssignTo: &mw.idvalue,
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
						//Text:     "请输入作者ID",
						Text: "a2c99b15af2b413ea29c6ebf40b9750c",
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
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						Text: "开始抓取",
						//MinSize:  Size{120, 30},
						AssignTo: &mw.start,
					},
					PushButton{
						Text: "查看(勾选列表索引)",
						OnClicked: func() {
							viewindex := mw.tv.SelectedIndexes()
							log.Println(viewindex)
							if len(viewindex) == 0 {
								walk.MsgBox(mw, "错误操作", "请先勾选索引然后点击查看,程序会调用默认浏览器直接跳转对应的页面", walk.MsgBoxIconWarning)
								return
							}
							viewUrl := mw.model.items[viewindex[0]].Srcurl
							cmd := exec.Command("rundll32 url.dll,FileProtocolHandler", viewUrl)
							cmd.Start()
						},
					},
				},
			},
			Composite{
				Layout: VBox{MarginsZero: true},
				Children: []Widget{
					TableView{
						AssignTo:              &mw.tv,
						AlternatingRowBGColor: walk.RGB(239, 239, 239),
						CheckBoxes:            true,
						ColumnsOrderable:      true,
						LastColumnStretched:   true,
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
						Model: mw.model,
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	// 处理抓取状态
	mw.start.Clicked().Attach(func() {

		mw.start.SetText("正在努力抓取中···")
		mw.start.SetEnabled(false)
		mw.Spider()
		mw.start.SetText("开始抓取")
		mw.start.SetEnabled(true)

	})

	// 处理作者ID获取状态
	mw.getID.Clicked().Attach(func() {

		mw.idvalue.SetText("")
		mw.getID.SetText("正在努力解析中···")
		mw.getID.SetEnabled(false)
		mw.getAuthorID()
		mw.getID.SetText("解析作者ID")
		mw.getID.SetEnabled(true)

	})

	//初始化
	mw.dayu.SetChecked(true)
	mw.Run()
}

func (mw *MyMainWindow) Spider() {
	AuthorID := mw.id.Text()
	log.Println(AuthorID)
	if len(AuthorID) != 32 {
		walk.MsgBox(mw, "ID错误", "请填写正确的作者ID(*大鱼号的作者ID默认为32位*),作者ID需要与目标平台匹配", walk.MsgBoxIconWarning)
		return
	}

	Hotvalue := mw.hot.Text()
	if _, err := strconv.Atoi(Hotvalue); err != nil {
		walk.MsgBox(mw, "阅读量设置错误", "请填写正确的数字", walk.MsgBoxIconWarning)
		return
	}

	t1 := mw.timeFrom.Date()
	Timefrom := t1.Format("2006-01-02 15:04:05")

	t2 := mw.timeFrom.Date()
	Timeto := t2.Format("2006-01-02 15:04:05")

	if mw.dayu.Checked() == true {
		results := Dayu(AuthorID, Hotvalue, Timefrom, Timeto)
		if len(results) == 0 {
			walk.MsgBox(mw, "结束", "没有找到符合要求的数据，请重新选择", walk.MsgBoxIconInformation)
			return
		} else {
			walk.MsgBox(mw, "结束", "数据已获取成功", walk.MsgBoxIconInformation)
		}

		m := new(ResultsTableModel)
		m.items = make([]*ResultsTable, len(results))
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

			pt, _ := time.Parse("2006-01-02 15:04:05", result.Publishtime)
			log.Println(pt)
			// nt, _ := strconv.Atoi(result.Hot)

			mw.model.items = append(mw.model.items, &ResultsTable{
				Index:       i,
				Type:        result.Type,
				Category:    result.Category,
				Title:       result.Title,
				Coverlink:   result.Coverlink,
				Publishtime: pt,
				Hot:         result.Hot,
				Srcurl:      result.Url,
			})

			mw.model.PublishRowsReset()
			mw.model.Sort(mw.model.sortColumn, mw.model.sortOrder)
			mw.tv.SetSelectedIndexes([]int{})
		}

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
func Dayu(AuthorID, Hotvalue, Timefrom, Timeto string) (results []Result) {

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
		//TODO:分析最符合的阅读量
		// reg := regexp.MustCompile(`"click_total":(.*?),`)
		// result.Hot, _ = strconv.Atoi(reg.FindString(string(resp.Body)))
		reg1 := regexp.MustCompile(`"click1":(.*?)`)
		reg2 := regexp.MustCompile(`"click2":(.*?),`)
		reg3 := regexp.MustCompile(`"click3":(.*?),`)
		c1, err := strconv.Atoi(reg1.FindStringSubmatch(string(resp.Body))[1])
		c2, err := strconv.Atoi(reg2.FindStringSubmatch(string(resp.Body))[1])
		c3, err := strconv.Atoi(reg3.FindStringSubmatch(string(resp.Body))[1])
		if err != nil {
			log.Fatal(err)
		}
		result.Hot = c1 + c2 + c3

		id := resp.Request.Ctx.Get("id")
		result.Url = fmt.Sprintf("html: http://a.mp.uc.cn/article.html?uc_param_str=frdnsnpfvecpntnwprdssskt&from=media#!wm_cid=%s", id)
		result.Type = resp.Request.Ctx.Get("type")
		result.Category = resp.Request.Ctx.Get("category")
		result.Title = resp.Request.Ctx.Get("title")
		result.Coverlink = resp.Request.Ctx.Get("coverlink")
		result.Publishtime = resp.Request.Ctx.Get("publishtime")

		// 根据阅读量筛选
		hv, _ := strconv.Atoi(Hotvalue)
		if hv <= result.Hot {
			results = append(results, result)
		}

	})

	c.OnResponse(func(resp *colly.Response) {
		productionlist := &Resp{}
		if err := json.Unmarshal([]byte(string(resp.Body)), productionlist); err != nil {
			log.Fatal(err)
		}
		for _, oneproduction := range productionlist.Data {

			//根据发布的时间段筛选
			changetime := strings.Replace(strings.Split(oneproduction.Publishtime, ".")[0], "T", " ", -1)
			timeline, err := time.Parse("2006-01-02 15:04:05", changetime)

			tfrom, err := time.Parse("2006-01-02 15:04:05", Timefrom)
			tto, err := time.Parse("2006-01-02 15:04:05", Timeto)

			if err != nil {
				log.Fatal(err)
			}
			if tfrom.Before(timeline) && tto.Before(timeline) {
				id := oneproduction.ID
				timel := timeline.Format("2006-01-02 15:04:05")

				targetURL := fmt.Sprintf("http://ff.dayu.com/contents/%s?biz_id=1002&_fetch_author=1&_fetch_incrs=1", id)
				ctx := colly.NewContext()
				ctx.Put("id", oneproduction.ID)
				ctx.Put("type", strconv.Itoa(oneproduction.Type))
				ctx.Put("category", oneproduction.Category)
				ctx.Put("title", oneproduction.Title)
				ctx.Put("coverlink", oneproduction.Coverlink)
				ctx.Put("publishtime", timel)
				getHotCollector.Request("GET", targetURL, nil, ctx, nil)
				log.Println("Visiting: ", targetURL)
			}
		}

	})

	visitUrl := fmt.Sprintf("http://ff.dayu.com/contents/author/%s?biz_id=1002&_size=1000000", AuthorID)
	c.Visit(visitUrl)

	getHotCollector.Wait()

	log.Println(len(results))
	return

}

/*============提===取===作====者=====ID==========*/

func (mw *MyMainWindow) getAuthorID() {
	sharelink := mw.sharelink.Text()
	if len(sharelink) == 0 {
		walk.MsgBox(mw, "错误", "请填写图文分享连接后再进行操作", walk.MsgBoxIconWarning)
		return
	}

	if mw.dayu.Checked() == true {

		reg := regexp.MustCompile(`wm_id=(.*?)&title_type`)
		rF := reg.FindStringSubmatch(sharelink)
		if len(rF) == 2 {
			theid := rF[1]
			mw.idvalue.SetText(theid)
		} else {
			walk.MsgBox(mw, "错误", "请填写对应平台的图文分享链接", walk.MsgBoxIconWarning)
			return
		}

	} else if mw.baijia.Checked() == true {

		//百家号平台爬虫(暂不支持)
		walk.MsgBox(mw, "Sorry", "百家号平台爬虫暂不支持,请选择大鱼号平台", walk.MsgBoxIconInformation)
		return

	} else if mw.qie.Checked() == true {

		//企鹅号平台爬虫(暂不支持)
		walk.MsgBox(mw, "Sorry", "企鹅号平台爬虫暂不支持,请选择大鱼号平台", walk.MsgBoxIconInformation)
		return

	}
}

/*============辅===助===函===数==============*/
