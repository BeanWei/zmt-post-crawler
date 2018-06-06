package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

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

func main() {

	results := []Result{}

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

}
