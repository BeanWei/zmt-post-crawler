package main

import (
	"log"
	"os"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("cookie", "这里直接填入从浏览器中拿到的cookie即可")
		r.Headers.Set("referer", "https://mp.dayu.com/dashboard/video/write?spm=a2s0i.db_messages_comment.menu.4.52ad3caaCO3AVH")
		r.Headers.Set("user-agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.79 Safari/537.36")
	})
	c.OnError(func(r *colly.Response, e error) {
		log.Println("Line_18-error:", e, r.Request.URL, string(r.Body))
	})
	c.OnResponse(func(r *colly.Response) {
		log.Println("Response received:", r.StatusCode)
		html := string(r.Body)
		f, err := os.OpenFile("test.html", os.O_WRONLY|os.O_CREATE, 0777)
		defer f.Close()
		if err != nil {
			log.Println("Line_26-Error:", err)
		} else {
			_, err := f.Write([]byte(html))
			if err != nil {
				log.Println("Line_30-Error:", err)
			}
		}

	})

	c.Visit("https://mp.dayu.com/dashboard/contents")
}
